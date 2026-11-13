package websocket

import (
	"log"
	"net/http"
	"sync"
	"time"

	gorilla "github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 50 * time.Second
	maxMessageSize = 512
)

type Client struct {
	hub  *Hub
	conn *gorilla.Conn
	send chan []byte
}

type Hub struct {
	mu            sync.Mutex
	clients       map[*Client]bool
	upgrader      gorilla.Upgrader
	allowedOrigin string
}

func NewHub(allowedOrigin string) *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		allowedOrigin: allowedOrigin,
		upgrader: gorilla.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				return origin == "" || origin == allowedOrigin
			},
		},
	}
}

// RegisterRoutes registers the real-time alert WebSocket endpoint.
func (h *Hub) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /ws/alerts", h.ServeWS)
}

// ServeWS upgrades an HTTP connection to a WebSocket connection.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 32),
	}

	h.mu.Lock()
	h.clients[client] = true
	clientCount := len(h.clients)
	h.mu.Unlock()

	log.Printf("WebSocket client connected. Active clients: %d", clientCount)

	go client.writePump()
	go client.readPump()
}

// Broadcast sends a JSON message to all connected clients.
// Messages are queued without blocking the caller.
func (h *Hub) Broadcast(message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			delete(h.clients, client)
			close(client.send)
			_ = client.conn.Close()
			log.Println("Removed WebSocket client: outgoing queue full")
		}
	}
}

func (h *Hub) unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[client]; exists {
		delete(h.clients, client)
		close(client.send)
		log.Printf("WebSocket client disconnected. Active clients: %d", len(h.clients))
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister(c)
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		// The server currently accepts no client commands.
		// Reading keeps the connection alive and detects disconnects.
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.unregister(c)
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				_ = c.conn.WriteMessage(gorilla.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(gorilla.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(gorilla.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

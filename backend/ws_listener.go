package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(
		"ws://localhost:8080/ws/alerts",
		nil,
	)
	if err != nil {
		log.Fatalf("WebSocket connection failed: %v", err)
	}
	defer conn.Close()

	log.Println("CONNECTED: Listening for GeoTrack alerts...")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Connection closed: %v", err)
			return
		}

		log.Printf("REAL-TIME ALERT RECEIVED: %s", message)
		_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	}
}

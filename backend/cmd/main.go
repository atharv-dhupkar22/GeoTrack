package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"geotrack/backend/internal/config"
	"geotrack/backend/internal/database"
	"geotrack/backend/internal/handlers"
	"geotrack/backend/internal/middleware"
	"geotrack/backend/internal/repository"
	"geotrack/backend/internal/services"
	ws "geotrack/backend/internal/websocket"
)

func main() {
	// Load application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to PostgreSQL.
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	log.Println("PostgreSQL connection established successfully")

	// Initialize repository.
	repo := repository.New(db)

	// Initialize WebSocket hub.
	websocketHub := ws.NewHub(cfg.FrontendURL)

	// Initialize geofence service and handler.
	geofenceService := services.NewGeofenceService(repo)
	geofenceHandler := handlers.NewGeofenceHandler(geofenceService)

	// Initialize vehicle service and handler.
	vehicleService := services.NewVehicleService(repo)
	vehicleService.SetAlertPublisher(websocketHub)

	vehicleHandler := handlers.NewVehicleHandler(vehicleService)

	// Initialize GPS location handler.
	locationHandler := handlers.NewLocationHandler(vehicleService)

	// Initialize alert service and handler.
	alertService := services.NewAlertService(repo)
	alertHandler := handlers.NewAlertHandler(alertService)

	// Initialize violation history service and handler.
	violationService := services.NewViolationService(repo)
	violationHandler := handlers.NewViolationHandler(violationService)

	// Initialize HTTP router.
	mux := http.NewServeMux()

	// Health endpoint.
	mux.HandleFunc("GET /health", healthHandler(db))

	// Register API routes.
	geofenceHandler.RegisterRoutes(mux)
	vehicleHandler.RegisterRoutes(mux)
	locationHandler.RegisterRoutes(mux)
	alertHandler.RegisterRoutes(mux)
	violationHandler.RegisterRoutes(mux)

	// Register WebSocket route.
	websocketHub.RegisterRoutes(mux)

	// Apply middleware.
	var handler http.Handler = mux
	handler = corsMiddleware(cfg.FrontendURL, handler)
	handler = middleware.Timing(handler)

	// Configure HTTP server.
	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start HTTP server.
	go func() {
		log.Printf("GeoTrack backend starting on port %s", cfg.AppPort)
		log.Printf("Environment: %s", cfg.AppEnv)
		log.Printf("Health endpoint: http://localhost:%s/health", cfg.AppPort)
		log.Printf("Geofence endpoint: http://localhost:%s/geofences", cfg.AppPort)
		log.Printf("Vehicle endpoint: http://localhost:%s/vehicles", cfg.AppPort)
		log.Printf("GPS endpoint: http://localhost:%s/vehicles/location", cfg.AppPort)
		log.Printf("Alert endpoint: http://localhost:%s/alerts", cfg.AppPort)
		log.Printf("Alert configurations: http://localhost:%s/alerts/configurations", cfg.AppPort)
		log.Printf("Violation history: http://localhost:%s/violations/history", cfg.AppPort)
		log.Printf("WebSocket endpoint: ws://localhost:%s/ws/alerts", cfg.AppPort)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Wait for termination signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutdown signal received. Stopping GeoTrack backend...")

	// Gracefully shut down HTTP server.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)

		if closeErr := server.Close(); closeErr != nil {
			log.Printf("HTTP server close failed: %v", closeErr)
		}
	}

	log.Println("GeoTrack backend stopped successfully")
}

func healthHandler(db interface {
	Ping(context.Context) error
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{
				"status":   "error",
				"service":  "backend",
				"database": "disconnected",
				"message":  "Database connection unavailable",
				"time_ns":  "0",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"status":   "success",
			"service":  "backend",
			"database": "connected",
			"message":  "GeoTrack API is running",
			"time_ns":  "0",
		})
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func corsMiddleware(
	allowedOrigin string,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS",
		)

		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization",
		)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

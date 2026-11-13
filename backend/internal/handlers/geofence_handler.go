package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"geotrack/backend/internal/models"
	"geotrack/backend/internal/services"
)

// GeofenceHandler handles HTTP requests for geofences.
type GeofenceHandler struct {
	service *services.GeofenceService
}

// NewGeofenceHandler creates a new geofence handler.
func NewGeofenceHandler(service *services.GeofenceService) *GeofenceHandler {
	return &GeofenceHandler{
		service: service,
	}
}

// RegisterRoutes registers geofence API endpoints.
func (h *GeofenceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /geofences", h.CreateGeofence)
	mux.HandleFunc("GET /geofences", h.ListGeofences)
}

// CreateGeofence handles POST /geofences.
func (h *GeofenceHandler) CreateGeofence(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		writeGeofenceJSON(w, http.StatusInternalServerError, "error",
			"Geofence service is unavailable", nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var request models.CreateGeofenceRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeGeofenceJSON(w, http.StatusBadRequest, "error",
			"Invalid JSON request: "+err.Error(), nil)
		return
	}

	// Reject additional JSON values after the first object.
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeGeofenceJSON(w, http.StatusBadRequest, "error",
			"Request must contain exactly one JSON object", nil)
		return
	}

	geofence, err := h.service.Create(r.Context(), request)
	if err != nil {
		message := err.Error()

		// Validation errors are client errors; database failures are server errors.
		if strings.HasPrefix(message, "failed to create geofence:") ||
			strings.HasPrefix(message, "failed to encode geofence coordinates:") {
			writeGeofenceJSON(w, http.StatusInternalServerError, "error",
				"Unable to save geofence", nil)
			return
		}

		writeGeofenceJSON(w, http.StatusBadRequest, "error", message, nil)
		return
	}

	writeGeofenceJSON(w, http.StatusCreated, "success",
		"Geofence created successfully", geofence)
}

// ListGeofences handles GET /geofences.
func (h *GeofenceHandler) ListGeofences(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		writeGeofenceJSON(w, http.StatusInternalServerError, "error",
			"Geofence service is unavailable", nil)
		return
	}

	category := r.URL.Query().Get("category")

	geofences, err := h.service.List(r.Context(), category)
	if err != nil {
		writeGeofenceJSON(w, http.StatusBadRequest, "error",
			err.Error(), nil)
		return
	}

	writeGeofenceJSON(w, http.StatusOK, "success",
		"Geofences retrieved successfully", geofences)
}

// writeGeofenceJSON writes a consistent JSON API response.
// The Timing middleware adds the time_ns field.
func writeGeofenceJSON(
	w http.ResponseWriter,
	statusCode int,
	status string,
	message string,
	data any,
) {
	response := models.APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
		TimeNS:  time.Duration(0).String(),
	}

	// The Timing middleware replaces this placeholder with actual elapsed nanoseconds.
	response.TimeNS = "0"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// The response may already have been written; log-free handler.
		return
	}
}

package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"

	"geotrack/backend/internal/models"
	"geotrack/backend/internal/services"
)

// LocationHandler handles vehicle GPS location requests.
type LocationHandler struct {
	service *services.VehicleService
}

// NewLocationHandler creates a location handler.
func NewLocationHandler(service *services.VehicleService) *LocationHandler {
	return &LocationHandler{
		service: service,
	}
}

// RegisterRoutes registers GPS location endpoints.
func (h *LocationHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /vehicles/location", h.RecordLocation)
	mux.HandleFunc("GET /vehicles/location/{vehicle_id}", h.GetLatestLocation)
}

// RecordLocation accepts and stores a vehicle GPS location.
func (h *LocationHandler) RecordLocation(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	var request models.CreateVehicleLocationRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeLocationJSON(
			w,
			http.StatusBadRequest,
			"error",
			"Invalid request body",
			nil,
		)
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err == nil {
		writeLocationJSON(
			w,
			http.StatusBadRequest,
			"error",
			"Request must contain only one JSON object",
			nil,
		)
		return
	}

	location, err := h.service.RecordLocation(r.Context(), request)
	if err != nil {
		statusCode := http.StatusBadRequest

		if strings.Contains(err.Error(), "failed to save vehicle location") {
			statusCode = http.StatusNotFound
		}

		if strings.Contains(err.Error(), "database repository") ||
			strings.Contains(err.Error(), "vehicle service") {
			statusCode = http.StatusInternalServerError
		}

		writeLocationJSON(
			w,
			statusCode,
			"error",
			err.Error(),
			nil,
		)
		return
	}

	writeLocationJSON(
		w,
		http.StatusCreated,
		"success",
		"Vehicle location recorded successfully",
		location,
	)
}

// GetLatestLocation retrieves the latest location for a vehicle.
func (h *LocationHandler) GetLatestLocation(
	w http.ResponseWriter,
	r *http.Request,
) {
	vehicleID := r.PathValue("vehicle_id")

	location, err := h.service.GetLatestLocation(r.Context(), vehicleID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to retrieve vehicle location"

		if errors.Is(err, pgx.ErrNoRows) {
			statusCode = http.StatusNotFound
			message = "No location found for this vehicle"
		}

		if strings.Contains(err.Error(), "vehicle ID is required") {
			statusCode = http.StatusBadRequest
			message = err.Error()
		}

		writeLocationJSON(w, statusCode, "error", message, nil)
		return
	}

	writeLocationJSON(
		w,
		http.StatusOK,
		"success",
		"Latest vehicle location retrieved successfully",
		location,
	)
}

// writeLocationJSON writes a standard API response.
func writeLocationJSON(
	w http.ResponseWriter,
	statusCode int,
	status string,
	message string,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
		TimeNS:  "0",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

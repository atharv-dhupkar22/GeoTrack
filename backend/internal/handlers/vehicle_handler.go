package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"geotrack/backend/internal/models"
	"geotrack/backend/internal/services"
)

// VehicleHandler handles vehicle API requests.
type VehicleHandler struct {
	service *services.VehicleService
}

// NewVehicleHandler creates a vehicle handler.
func NewVehicleHandler(service *services.VehicleService) *VehicleHandler {
	return &VehicleHandler{
		service: service,
	}
}

// RegisterRoutes registers vehicle endpoints.
func (h *VehicleHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /vehicles", h.CreateVehicle)
	mux.HandleFunc("GET /vehicles", h.ListVehicles)
}

// CreateVehicle registers a new vehicle.
func (h *VehicleHandler) CreateVehicle(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	var request models.CreateVehicleRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeVehicleJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	if err := decoder.Decode(&struct{}{}); err == nil {
		writeVehicleJSON(w, http.StatusBadRequest, "error", "Request must contain only one JSON object", nil)
		return
	}

	vehicle, err := h.service.Create(r.Context(), request)
	if err != nil {
		statusCode := http.StatusInternalServerError

		if !strings.Contains(err.Error(), "failed to create vehicle") &&
			!strings.Contains(err.Error(), "database repository") &&
			!strings.Contains(err.Error(), "vehicle service") {
			statusCode = http.StatusBadRequest
		}

		if strings.Contains(err.Error(), "duplicate key") {
			statusCode = http.StatusConflict
		}

		writeVehicleJSON(w, statusCode, "error", err.Error(), nil)
		return
	}

	writeVehicleJSON(
		w,
		http.StatusCreated,
		"success",
		"Vehicle registered successfully",
		vehicle,
	)
}

// ListVehicles returns all registered vehicles.
func (h *VehicleHandler) ListVehicles(
	w http.ResponseWriter,
	r *http.Request,
) {
	vehicles, err := h.service.List(r.Context())
	if err != nil {
		writeVehicleJSON(
			w,
			http.StatusInternalServerError,
			"error",
			"Failed to retrieve vehicles",
			nil,
		)
		return
	}

	writeVehicleJSON(
		w,
		http.StatusOK,
		"success",
		"Vehicles retrieved successfully",
		vehicles,
	)
}

// writeVehicleJSON writes a consistent API response.
func writeVehicleJSON(
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
		_ = errors.New("failed to encode vehicle response")
	}
}

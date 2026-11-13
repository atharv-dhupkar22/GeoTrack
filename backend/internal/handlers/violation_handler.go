package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"geotrack/backend/internal/models"
	"geotrack/backend/internal/services"
)

type ViolationHandler struct {
	service *services.ViolationService
}

func NewViolationHandler(service *services.ViolationService) *ViolationHandler {
	return &ViolationHandler{
		service: service,
	}
}

func (h *ViolationHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /violations/history", h.ListHistory)
}

func (h *ViolationHandler) ListHistory(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := 1
	limit := 10

	if value := query.Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeViolationJSON(w, http.StatusBadRequest, "error", "page must be a valid integer", nil)
			return
		}
		page = parsed
	}

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeViolationJSON(w, http.StatusBadRequest, "error", "limit must be a valid integer", nil)
			return
		}
		limit = parsed
	}

	result, err := h.service.ListHistory(
		r.Context(),
		strings.TrimSpace(query.Get("vehicle_id")),
		strings.TrimSpace(query.Get("geofence_id")),
		strings.TrimSpace(query.Get("event_type")),
		page,
		limit,
	)

	if err != nil {
		status := http.StatusBadRequest

		if strings.Contains(err.Error(), "failed to count violations") ||
			strings.Contains(err.Error(), "failed to retrieve violations") ||
			strings.Contains(err.Error(), "failed to scan violation") ||
			strings.Contains(err.Error(), "failed while reading violations") {
			status = http.StatusInternalServerError
		}

		writeViolationJSON(w, status, "error", err.Error(), nil)
		return
	}

	writeViolationJSON(w, http.StatusOK, "success", "Violation history retrieved successfully", result)
}

func writeViolationJSON(
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
		TimeNS:  "0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		_ = errors.New("failed to encode violation response")
	}
}

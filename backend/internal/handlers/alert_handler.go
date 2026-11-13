package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"geotrack/backend/internal/models"
	"geotrack/backend/internal/services"
)

// AlertHandler handles alert configuration and notification requests.
type AlertHandler struct {
	service *services.AlertService
}

// NewAlertHandler creates an alert handler.
func NewAlertHandler(service *services.AlertService) *AlertHandler {
	return &AlertHandler{
		service: service,
	}
}

// RegisterRoutes registers alert endpoints.
func (h *AlertHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /alerts/configure", h.CreateAlert)
	mux.HandleFunc("GET /alerts/configurations", h.ListConfigurations)
	mux.HandleFunc("GET /alerts", h.ListAlertHistory)
}

// CreateAlert creates a new alert configuration.
func (h *AlertHandler) CreateAlert(
	w http.ResponseWriter,
	r *http.Request,
) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	var request models.CreateAlertConfigRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeAlertJSON(
			w,
			http.StatusBadRequest,
			"error",
			"Invalid request body",
			nil,
		)
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeAlertJSON(
			w,
			http.StatusBadRequest,
			"error",
			"Request must contain only one JSON object",
			nil,
		)
		return
	}

	alert, err := h.service.Create(r.Context(), request)
	if err != nil {
		statusCode := http.StatusBadRequest

		if strings.Contains(err.Error(), "service is not initialized") ||
			strings.Contains(err.Error(), "repository is not initialized") {
			statusCode = http.StatusInternalServerError
		}

		writeAlertJSON(
			w,
			statusCode,
			"error",
			err.Error(),
			nil,
		)
		return
	}

	writeAlertJSON(
		w,
		http.StatusCreated,
		"success",
		"Alert configuration created successfully",
		alert,
	)
}

// ListConfigurations retrieves configured alert rules.
func (h *AlertHandler) ListConfigurations(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	alerts, err := h.service.List(
		r.Context(),
		query.Get("geofence_id"),
		query.Get("vehicle_id"),
		query.Get("status"),
	)
	if err != nil {
		statusCode := http.StatusBadRequest

		if strings.Contains(err.Error(), "service is not initialized") ||
			strings.Contains(err.Error(), "repository is not initialized") {
			statusCode = http.StatusInternalServerError
		}

		writeAlertJSON(
			w,
			statusCode,
			"error",
			err.Error(),
			nil,
		)
		return
	}

	writeAlertJSON(
		w,
		http.StatusOK,
		"success",
		"Alert configurations retrieved successfully",
		alerts,
	)
}

// ListAlertHistory retrieves persisted notification events.
func (h *AlertHandler) ListAlertHistory(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	page := 1
	limit := 10

	if value := query.Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeAlertJSON(
				w,
				http.StatusBadRequest,
				"error",
				"page must be a valid integer",
				nil,
			)
			return
		}
		page = parsed
	}

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeAlertJSON(
				w,
				http.StatusBadRequest,
				"error",
				"limit must be a valid integer",
				nil,
			)
			return
		}
		limit = parsed
	}

	history, err := h.service.ListHistory(
		r.Context(),
		query.Get("vehicle_id"),
		query.Get("geofence_id"),
		query.Get("event_type"),
		page,
		limit,
	)
	if err != nil {
		statusCode := http.StatusBadRequest

		if strings.Contains(err.Error(), "service is not initialized") ||
			strings.Contains(err.Error(), "repository is not initialized") ||
			strings.Contains(err.Error(), "failed to retrieve") ||
			strings.Contains(err.Error(), "failed to count") {
			statusCode = http.StatusInternalServerError
		}

		writeAlertJSON(
			w,
			statusCode,
			"error",
			err.Error(),
			nil,
		)
		return
	}

	writeAlertJSON(
		w,
		http.StatusOK,
		"success",
		"Alert notification history retrieved successfully",
		history,
	)
}

// writeAlertJSON writes a consistent API response.
func writeAlertJSON(
	w http.ResponseWriter,
	statusCode int,
	status string,
	message string,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")

	response := models.APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
		TimeNS:  "0",
	}

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

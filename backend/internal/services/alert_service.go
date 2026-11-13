package services

import (
	"context"
	"errors"
	"strings"

	"geotrack/backend/internal/models"
	"geotrack/backend/internal/repository"
)

// AlertService handles alert configuration and notification history.
type AlertService struct {
	repository *repository.Repository
}

// NewAlertService creates an alert service.
func NewAlertService(repo *repository.Repository) *AlertService {
	return &AlertService{
		repository: repo,
	}
}

// Create validates and creates an alert configuration.
func (s *AlertService) Create(
	ctx context.Context,
	request models.CreateAlertConfigRequest,
) (*models.AlertConfig, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("alert service is not initialized")
	}

	request.GeofenceID = strings.TrimSpace(request.GeofenceID)
	request.EventType = strings.ToLower(strings.TrimSpace(request.EventType))

	if request.VehicleID != nil {
		vehicleID := strings.TrimSpace(*request.VehicleID)
		if vehicleID == "" {
			request.VehicleID = nil
		} else {
			request.VehicleID = &vehicleID
		}
	}

	if request.GeofenceID == "" {
		return nil, errors.New("geofence_id is required")
	}

	switch request.EventType {
	case "entry", "exit", "both":
	default:
		return nil, errors.New("event_type must be entry, exit, or both")
	}

	return s.repository.CreateAlertConfig(ctx, request)
}

// List retrieves alert configurations using optional filters.
func (s *AlertService) List(
	ctx context.Context,
	geofenceID string,
	vehicleID string,
	status string,
) ([]models.AlertConfig, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("alert service is not initialized")
	}

	geofenceID = strings.TrimSpace(geofenceID)
	vehicleID = strings.TrimSpace(vehicleID)
	status = strings.ToLower(strings.TrimSpace(status))

	if status != "" && status != "active" && status != "inactive" {
		return nil, errors.New("status must be active or inactive")
	}

	return s.repository.ListAlertConfigs(
		ctx,
		geofenceID,
		vehicleID,
		status,
	)
}

// ListHistory validates filters and retrieves notification history.
func (s *AlertService) ListHistory(
	ctx context.Context,
	vehicleID string,
	geofenceID string,
	eventType string,
	page int,
	limit int,
) (*models.AlertHistoryResult, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("alert service is not initialized")
	}

	vehicleID = strings.TrimSpace(vehicleID)
	geofenceID = strings.TrimSpace(geofenceID)
	eventType = strings.ToLower(strings.TrimSpace(eventType))

	if eventType != "" && eventType != "entry" && eventType != "exit" {
		return nil, errors.New("event_type must be entry or exit")
	}

	if page < 1 {
		return nil, errors.New("page must be greater than or equal to 1")
	}

	if limit < 1 || limit > 100 {
		return nil, errors.New("limit must be between 1 and 100")
	}

	alerts, total, err := s.repository.ListAlertEvents(
		ctx,
		vehicleID,
		geofenceID,
		eventType,
		page,
		limit,
	)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return &models.AlertHistoryResult{
		Alerts:     alerts,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

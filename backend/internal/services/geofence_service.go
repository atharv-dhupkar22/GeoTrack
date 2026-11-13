package services

import (
	"context"
	"errors"
	"strings"

	"geotrack/backend/internal/geofence"
	"geotrack/backend/internal/models"
	"geotrack/backend/internal/repository"
)

// GeofenceService handles geofence business logic.
type GeofenceService struct {
	repository *repository.Repository
}

// NewGeofenceService creates a geofence service.
func NewGeofenceService(repo *repository.Repository) *GeofenceService {
	return &GeofenceService{
		repository: repo,
	}
}

// Create validates and saves a new geofence.
func (s *GeofenceService) Create(
	ctx context.Context,
	request models.CreateGeofenceRequest,
) (*models.Geofence, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("geofence service is not initialized")
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)
	request.Category = strings.TrimSpace(request.Category)

	if request.Name == "" {
		return nil, errors.New("geofence name is required")
	}

	if len(request.Name) > 150 {
		return nil, errors.New("geofence name must not exceed 150 characters")
	}

	if !isValidGeofenceCategory(request.Category) {
		return nil, errors.New(
			"invalid geofence category: use delivery_zone, restricted_zone, toll_zone, or customer_area",
		)
	}

	if err := geofence.ValidatePolygon(request.Coordinates); err != nil {
		return nil, err
	}

	return s.repository.CreateGeofence(ctx, request)
}

// List retrieves all geofences or filters them by category.
func (s *GeofenceService) List(
	ctx context.Context,
	category string,
) ([]models.Geofence, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("geofence service is not initialized")
	}

	category = strings.TrimSpace(category)

	if category != "" && !isValidGeofenceCategory(category) {
		return nil, errors.New(
			"invalid geofence category: use delivery_zone, restricted_zone, toll_zone, or customer_area",
		)
	}

	return s.repository.ListGeofences(ctx, category)
}

// isValidGeofenceCategory checks the supported assessment categories.
func isValidGeofenceCategory(category string) bool {
	switch category {
	case models.CategoryDeliveryZone,
		models.CategoryRestrictedZone,
		models.CategoryTollZone,
		models.CategoryCustomerArea:
		return true
	default:
		return false
	}
}

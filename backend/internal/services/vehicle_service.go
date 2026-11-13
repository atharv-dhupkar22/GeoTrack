package services

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"

	"geotrack/backend/internal/models"
	"geotrack/backend/internal/repository"
)

// AlertPublisher defines the interface for publishing real-time alerts.
type AlertPublisher interface {
	Broadcast(message []byte)
}

// VehicleService handles vehicle-related business logic.
type VehicleService struct {
	repository     *repository.Repository
	alertPublisher AlertPublisher
}

// NewVehicleService creates a vehicle service.
func NewVehicleService(repo *repository.Repository) *VehicleService {
	return &VehicleService{
		repository: repo,
	}
}

// SetAlertPublisher connects the service to the WebSocket hub.
func (s *VehicleService) SetAlertPublisher(publisher AlertPublisher) {
	if s == nil {
		return
	}

	s.alertPublisher = publisher
}

// Create registers a new vehicle after validating its details.
func (s *VehicleService) Create(
	ctx context.Context,
	request models.CreateVehicleRequest,
) (*models.Vehicle, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("vehicle service is not initialized")
	}

	request.VehicleNumber = strings.TrimSpace(request.VehicleNumber)
	request.DriverName = strings.TrimSpace(request.DriverName)
	request.VehicleType = strings.TrimSpace(request.VehicleType)
	request.Phone = strings.TrimSpace(request.Phone)

	if request.VehicleNumber == "" {
		return nil, errors.New("vehicle number is required")
	}
	if len(request.VehicleNumber) > 50 {
		return nil, errors.New("vehicle number cannot exceed 50 characters")
	}
	if request.DriverName == "" {
		return nil, errors.New("driver name is required")
	}
	if len(request.DriverName) > 150 {
		return nil, errors.New("driver name cannot exceed 150 characters")
	}
	if request.VehicleType == "" {
		return nil, errors.New("vehicle type is required")
	}
	if len(request.VehicleType) > 50 {
		return nil, errors.New("vehicle type cannot exceed 50 characters")
	}
	if request.Phone == "" {
		return nil, errors.New("phone number is required")
	}
	if len(request.Phone) > 30 {
		return nil, errors.New("phone number cannot exceed 30 characters")
	}

	return s.repository.CreateVehicle(ctx, request)
}

// List retrieves all registered vehicles.
func (s *VehicleService) List(
	ctx context.Context,
) ([]models.Vehicle, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("vehicle service is not initialized")
	}

	return s.repository.ListVehicles(ctx)
}

// RecordLocation validates and saves a vehicle GPS location.
// Matching alert events are broadcast only after the database transaction commits.
func (s *VehicleService) RecordLocation(
	ctx context.Context,
	request models.CreateVehicleLocationRequest,
) (*models.VehicleLocation, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("vehicle service is not initialized")
	}

	request.VehicleID = strings.TrimSpace(request.VehicleID)

	if request.VehicleID == "" {
		return nil, errors.New("vehicle_id is required")
	}

	if math.IsNaN(request.Latitude) ||
		math.IsInf(request.Latitude, 0) ||
		request.Latitude < -90 ||
		request.Latitude > 90 {
		return nil, errors.New("latitude must be between -90 and 90")
	}

	if math.IsNaN(request.Longitude) ||
		math.IsInf(request.Longitude, 0) ||
		request.Longitude < -180 ||
		request.Longitude > 180 {
		return nil, errors.New("longitude must be between -180 and 180")
	}

	if request.RecordedAt.IsZero() {
		request.RecordedAt = time.Now().UTC()
	}

	location, events, err := s.repository.CreateVehicleLocationWithAlerts(
		ctx,
		request,
	)
	if err != nil {
		return nil, err
	}

	// Publish only events returned after a successful transaction commit.
	if s.alertPublisher != nil {
		for _, event := range events {
			message, marshalErr := json.Marshal(event)
			if marshalErr != nil {
				continue
			}

			s.alertPublisher.Broadcast(message)
		}
	}

	return location, nil
}

// GetLatestLocation retrieves the latest GPS location of a vehicle.
func (s *VehicleService) GetLatestLocation(
	ctx context.Context,
	vehicleID string,
) (*models.VehicleLocation, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("vehicle service is not initialized")
	}

	vehicleID = strings.TrimSpace(vehicleID)

	if vehicleID == "" {
		return nil, errors.New("vehicle ID is required")
	}

	return s.repository.GetLatestVehicleLocation(ctx, vehicleID)
}

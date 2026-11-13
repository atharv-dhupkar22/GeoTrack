package models

import "time"

// Supported geofence categories.
const (
	CategoryDeliveryZone   = "delivery_zone"
	CategoryRestrictedZone = "restricted_zone"
	CategoryTollZone       = "toll_zone"
	CategoryCustomerArea   = "customer_area"
)

// Geofence represents a geographical boundary.
type Geofence struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Coordinates [][]float64 `json:"coordinates"`
	Category    string      `json:"category"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
}

// CreateGeofenceRequest contains the data required to create a geofence.
type CreateGeofenceRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Coordinates [][]float64 `json:"coordinates"`
	Category    string      `json:"category"`
}

// Vehicle represents a registered vehicle.
type Vehicle struct {
	ID            string    `json:"id"`
	VehicleNumber string    `json:"vehicle_number"`
	DriverName    string    `json:"driver_name"`
	VehicleType   string    `json:"vehicle_type"`
	Phone         string    `json:"phone"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateVehicleRequest contains the data required to register a vehicle.
type CreateVehicleRequest struct {
	VehicleNumber string `json:"vehicle_number"`
	DriverName    string `json:"driver_name"`
	VehicleType   string `json:"vehicle_type"`
	Phone         string `json:"phone"`
}

// VehicleLocation represents a vehicle's recorded location.
type VehicleLocation struct {
	ID         int64     `json:"id"`
	VehicleID  string    `json:"vehicle_id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	RecordedAt time.Time `json:"recorded_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateVehicleLocationRequest contains a new vehicle location.
type CreateVehicleLocationRequest struct {
	VehicleID  string    `json:"vehicle_id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	RecordedAt time.Time `json:"recorded_at,omitempty"`
}

// AlertConfig represents an alert rule for a geofence.
type AlertConfig struct {
	ID         string    `json:"id"`
	GeofenceID string    `json:"geofence_id"`
	VehicleID  *string   `json:"vehicle_id,omitempty"`
	EventType  string    `json:"event_type"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateAlertConfigRequest contains the data required to configure alerts.
type CreateAlertConfigRequest struct {
	GeofenceID string  `json:"geofence_id"`
	VehicleID  *string `json:"vehicle_id,omitempty"`
	EventType  string  `json:"event_type"`
}

// Violation represents a vehicle entering or exiting a geofence.
type Violation struct {
	ID         string    `json:"id"`
	VehicleID  string    `json:"vehicle_id"`
	GeofenceID string    `json:"geofence_id"`
	EventType  string    `json:"event_type"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	RecordedAt time.Time `json:"recorded_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// APIResponse provides a consistent response structure for API endpoints.
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	TimeNS  string `json:"time_ns"`
}

type AlertEvent struct {
	ID            string    `json:"id"`
	ViolationID   string    `json:"violation_id"`
	AlertConfigID string    `json:"alert_config_id"`
	VehicleID     string    `json:"vehicle_id"`
	GeofenceID    string    `json:"geofence_id"`
	EventType     string    `json:"event_type"`
	Message       string    `json:"message"`
	CreatedAt     time.Time `json:"created_at"`
}

type AlertHistoryResult struct {
	Alerts     []AlertEvent `json:"alerts"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	Total      int64        `json:"total"`
	TotalPages int          `json:"total_pages"`
}

type ViolationHistoryResult struct {
	Violations []Violation `json:"violations"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"geotrack/backend/internal/models"
)

// CreateAlertConfig saves an alert configuration.
func (r *Repository) CreateAlertConfig(
	ctx context.Context,
	request models.CreateAlertConfigRequest,
) (*models.AlertConfig, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("database repository is not initialized")
	}

	id, err := generateID("alert")
	if err != nil {
		return nil, fmt.Errorf("failed to generate alert ID: %w", err)
	}

	query := `
		INSERT INTO alert_configs (
			id,
			geofence_id,
			vehicle_id,
			event_type,
			status
		)
		VALUES ($1, $2, $3, $4, 'active')
		RETURNING
			id,
			geofence_id,
			vehicle_id,
			event_type,
			status,
			created_at
	`

	var alert models.AlertConfig
	var vehicleID *string

	err = r.db.QueryRow(
		ctx,
		query,
		id,
		request.GeofenceID,
		request.VehicleID,
		request.EventType,
	).Scan(
		&alert.ID,
		&alert.GeofenceID,
		&vehicleID,
		&alert.EventType,
		&alert.Status,
		&alert.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create alert configuration: %w", err)
	}

	alert.VehicleID = vehicleID
	return &alert, nil
}

// ListAlertConfigs retrieves alert configurations with optional filters.
func (r *Repository) ListAlertConfigs(
	ctx context.Context,
	geofenceID string,
	vehicleID string,
	status string,
) ([]models.AlertConfig, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("database repository is not initialized")
	}

	query := `
		SELECT
			id,
			geofence_id,
			COALESCE(vehicle_id, ''),
			event_type,
			status,
			created_at
		FROM alert_configs
		WHERE
			($1 = '' OR geofence_id = $1)
			AND ($2 = '' OR vehicle_id = $2)
			AND ($3 = '' OR status = $3)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		geofenceID,
		vehicleID,
		status,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve alert configurations: %w", err)
	}
	defer rows.Close()

	alerts := make([]models.AlertConfig, 0)

	for rows.Next() {
		var alert models.AlertConfig
		var storedVehicleID string

		if err := rows.Scan(
			&alert.ID,
			&alert.GeofenceID,
			&storedVehicleID,
			&alert.EventType,
			&alert.Status,
			&alert.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan alert configuration: %w", err)
		}

		if storedVehicleID != "" {
			alert.VehicleID = &storedVehicleID
		}

		alerts = append(alerts, alert)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading alert configurations: %w", err)
	}

	return alerts, nil
}

// ListAlertEvents retrieves persisted notification events with filters and pagination.
func (r *Repository) ListAlertEvents(
	ctx context.Context,
	vehicleID string,
	geofenceID string,
	eventType string,
	page int,
	limit int,
) ([]models.AlertEvent, int64, error) {
	if r == nil || r.db == nil {
		return nil, 0, errors.New("database repository is not initialized")
	}

	var total int64

	countQuery := `
		SELECT COUNT(*)
		FROM alerts
		WHERE
			($1 = '' OR vehicle_id = $1)
			AND ($2 = '' OR geofence_id = $2)
			AND ($3 = '' OR event_type = $3)
	`

	err := r.db.QueryRow(
		ctx,
		countQuery,
		vehicleID,
		geofenceID,
		eventType,
	).Scan(&total)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to count alert events: %w", err)
	}

	offset := (page - 1) * limit

	query := `
		SELECT
			id,
			violation_id,
			alert_config_id,
			vehicle_id,
			geofence_id,
			event_type,
			message,
			created_at
		FROM alerts
		WHERE
			($1 = '' OR vehicle_id = $1)
			AND ($2 = '' OR geofence_id = $2)
			AND ($3 = '' OR event_type = $3)
		ORDER BY created_at DESC, id DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(
		ctx,
		query,
		vehicleID,
		geofenceID,
		eventType,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve alert events: %w", err)
	}
	defer rows.Close()

	alerts := make([]models.AlertEvent, 0)

	for rows.Next() {
		var alert models.AlertEvent

		if err := rows.Scan(
			&alert.ID,
			&alert.ViolationID,
			&alert.AlertConfigID,
			&alert.VehicleID,
			&alert.GeofenceID,
			&alert.EventType,
			&alert.Message,
			&alert.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan alert event: %w", err)
		}

		alerts = append(alerts, alert)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error while reading alert events: %w", err)
	}

	return alerts, total, nil
}

package repository

import (
	"context"
	"fmt"

	"geotrack/backend/internal/models"
)

func (r *Repository) ListViolations(
	ctx context.Context,
	vehicleID string,
	geofenceID string,
	eventType string,
	page int,
	limit int,
) ([]models.Violation, int64, error) {
	if r == nil || r.db == nil {
		return nil, 0, fmt.Errorf("repository database is not initialized")
	}

	var total int64

	countQuery := `
		SELECT COUNT(*)
		FROM violations
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
		return nil, 0, fmt.Errorf("failed to count violations: %w", err)
	}

	offset := (page - 1) * limit

	selectQuery := `
		SELECT
			id,
			vehicle_id,
			geofence_id,
			event_type,
			latitude,
			longitude,
			recorded_at,
			created_at
		FROM violations
		WHERE
			($1 = '' OR vehicle_id = $1)
			AND ($2 = '' OR geofence_id = $2)
			AND ($3 = '' OR event_type = $3)
		ORDER BY recorded_at DESC, id DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(
		ctx,
		selectQuery,
		vehicleID,
		geofenceID,
		eventType,
		limit,
		offset,
	)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve violations: %w", err)
	}
	defer rows.Close()

	violations := make([]models.Violation, 0)

	for rows.Next() {
		var violation models.Violation

		err := rows.Scan(
			&violation.ID,
			&violation.VehicleID,
			&violation.GeofenceID,
			&violation.EventType,
			&violation.Latitude,
			&violation.Longitude,
			&violation.RecordedAt,
			&violation.CreatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan violation: %w", err)
		}

		violations = append(violations, violation)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed while reading violations: %w", err)
	}

	return violations, total, nil
}

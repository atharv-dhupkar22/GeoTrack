package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"geotrack/backend/internal/models"
)

// CreateVehicleLocation preserves the original repository method.
// It stores the location and detects transitions, while discarding
// the returned notification events.
func (r *Repository) CreateVehicleLocation(
	ctx context.Context,
	request models.CreateVehicleLocationRequest,
) (*models.VehicleLocation, error) {
	location, _, err := r.CreateVehicleLocationWithAlerts(ctx, request)
	return location, err
}

// CreateVehicleLocationWithAlerts stores a GPS location, detects
// geofence transitions, records violations, matches alert configurations,
// and persists matching alerts in one database transaction.
func (r *Repository) CreateVehicleLocationWithAlerts(
	ctx context.Context,
	request models.CreateVehicleLocationRequest,
) (*models.VehicleLocation, []models.AlertEvent, error) {
	if r == nil || r.db == nil {
		return nil, nil, errors.New("database repository is not initialized")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin location transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Serialize location processing for this vehicle.
	var vehicleID string

	err = tx.QueryRow(
		ctx,
		`SELECT id FROM vehicles WHERE id = $1 FOR UPDATE`,
		request.VehicleID,
	).Scan(&vehicleID)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find vehicle: %w", err)
	}

	// Save the GPS location.
	locationQuery := `
		INSERT INTO vehicle_locations (
			vehicle_id,
			latitude,
			longitude,
			location,
			recorded_at
		)
		VALUES (
			$1,
			$2,
			$3,
			ST_SetSRID(ST_MakePoint($3, $2), 4326),
			$4
		)
		RETURNING
			id,
			vehicle_id,
			latitude,
			longitude,
			recorded_at,
			created_at
	`

	var location models.VehicleLocation

	err = tx.QueryRow(
		ctx,
		locationQuery,
		request.VehicleID,
		request.Latitude,
		request.Longitude,
		request.RecordedAt,
	).Scan(
		&location.ID,
		&location.VehicleID,
		&location.Latitude,
		&location.Longitude,
		&location.RecordedAt,
		&location.CreatedAt,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to save vehicle location: %w", err)
	}

	// Find active geofences containing the current point.
	// ST_Covers includes points on the polygon boundary.
	geofenceQuery := `
		SELECT
			id,
			ST_Covers(
				boundary,
				ST_SetSRID(ST_MakePoint($1, $2), 4326)
			) AS is_inside
		FROM geofences
		WHERE status = 'active'
	`

	rows, err := tx.Query(
		ctx,
		geofenceQuery,
		request.Longitude,
		request.Latitude,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check geofences: %w", err)
	}

	type geofenceResult struct {
		ID       string
		IsInside bool
	}

	results := make([]geofenceResult, 0)

	for rows.Next() {
		var result geofenceResult

		if err := rows.Scan(&result.ID, &result.IsInside); err != nil {
			rows.Close()
			return nil, nil, fmt.Errorf("failed to read geofence result: %w", err)
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, nil, fmt.Errorf("failed while reading geofences: %w", err)
	}

	rows.Close()

	// Collect matching notification events.
	alertEvents := make([]models.AlertEvent, 0)

	// Compare current position with the previous state.
	for _, result := range results {
		var previousInside bool

		stateQuery := `
			SELECT is_inside
			FROM vehicle_geofence_states
			WHERE vehicle_id = $1 AND geofence_id = $2
			FOR UPDATE
		`

		err := tx.QueryRow(
			ctx,
			stateQuery,
			request.VehicleID,
			result.ID,
		).Scan(&previousInside)

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf(
				"failed to retrieve previous geofence state: %w",
				err,
			)
		}

		// A missing state is treated as outside.
		if previousInside != result.IsInside {
			eventType := "exit"
			if result.IsInside {
				eventType = "entry"
			}

			// First persist the violation.
			violationID, err := generateID("vio")
			if err != nil {
				return nil, nil, fmt.Errorf(
					"failed to generate violation ID: %w",
					err,
				)
			}

			violationQuery := `
				INSERT INTO violations (
					id,
					vehicle_id,
					geofence_id,
					event_type,
					latitude,
					longitude,
					recorded_at
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`

			_, err = tx.Exec(
				ctx,
				violationQuery,
				violationID,
				request.VehicleID,
				result.ID,
				eventType,
				request.Latitude,
				request.Longitude,
				request.RecordedAt,
			)

			if err != nil {
				return nil, nil, fmt.Errorf(
					"failed to record geofence violation: %w",
					err,
				)
			}

			// Find active alert configurations matching this event.
			configQuery := `
				SELECT id
				FROM alert_configs
				WHERE geofence_id = $1
					AND (vehicle_id IS NULL OR vehicle_id = $2)
					AND (event_type = $3 OR event_type = 'both')
					AND status = 'active'
				ORDER BY created_at ASC
			`

			configRows, err := tx.Query(
				ctx,
				configQuery,
				result.ID,
				request.VehicleID,
				eventType,
			)

			if err != nil {
				return nil, nil, fmt.Errorf(
					"failed to find matching alert configurations: %w",
					err,
				)
			}

			configIDs := make([]string, 0)

			for configRows.Next() {
				var configID string

				if err := configRows.Scan(&configID); err != nil {
					configRows.Close()
					return nil, nil, fmt.Errorf(
						"failed to read alert configuration: %w",
						err,
					)
				}

				configIDs = append(configIDs, configID)
			}

			if err := configRows.Err(); err != nil {
				configRows.Close()
				return nil, nil, fmt.Errorf(
					"error while reading alert configurations: %w",
					err,
				)
			}

			configRows.Close()

			// Persist one notification per matching configuration.
			for _, configID := range configIDs {
				alertID, err := generateID("evt")
				if err != nil {
					return nil, nil, fmt.Errorf(
						"failed to generate alert event ID: %w",
						err,
					)
				}

				message := fmt.Sprintf(
					"Vehicle %s %sed geofence %s",
					request.VehicleID,
					eventType,
					result.ID,
				)

				alertQuery := `
					INSERT INTO alerts (
						id,
						violation_id,
						alert_config_id,
						vehicle_id,
						geofence_id,
						event_type,
						message
					)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
					RETURNING created_at
				`

				var createdAt models.AlertEvent

				err = tx.QueryRow(
					ctx,
					alertQuery,
					alertID,
					violationID,
					configID,
					request.VehicleID,
					result.ID,
					eventType,
					message,
				).Scan(&createdAt.CreatedAt)

				if err != nil {
					return nil, nil, fmt.Errorf(
						"failed to persist alert event: %w",
						err,
					)
				}

				alertEvents = append(alertEvents, models.AlertEvent{
					ID:            alertID,
					ViolationID:   violationID,
					AlertConfigID: configID,
					VehicleID:     request.VehicleID,
					GeofenceID:    result.ID,
					EventType:     eventType,
					Message:       message,
					CreatedAt:     createdAt.CreatedAt,
				})
			}
		}

		// Store the latest inside/outside state.
		stateUpsertQuery := `
			INSERT INTO vehicle_geofence_states (
				vehicle_id,
				geofence_id,
				is_inside,
				updated_at
			)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (vehicle_id, geofence_id)
			DO UPDATE SET
				is_inside = EXCLUDED.is_inside,
				updated_at = NOW()
		`

		_, err = tx.Exec(
			ctx,
			stateUpsertQuery,
			request.VehicleID,
			result.ID,
			result.IsInside,
		)

		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to update geofence state: %w",
				err,
			)
		}
	}

	// Commit location, states, violations, and alerts together.
	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf(
			"failed to commit location transaction: %w",
			err,
		)
	}

	return &location, alertEvents, nil
}

// GetLatestVehicleLocation retrieves the latest recorded GPS location.
func (r *Repository) GetLatestVehicleLocation(
	ctx context.Context,
	vehicleID string,
) (*models.VehicleLocation, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("database repository is not initialized")
	}

	query := `
		SELECT
			id,
			vehicle_id,
			latitude,
			longitude,
			recorded_at,
			created_at
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY recorded_at DESC, id DESC
		LIMIT 1
	`

	var location models.VehicleLocation

	err := r.db.QueryRow(
		ctx,
		query,
		vehicleID,
	).Scan(
		&location.ID,
		&location.VehicleID,
		&location.Latitude,
		&location.Longitude,
		&location.RecordedAt,
		&location.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to retrieve latest vehicle location: %w",
			err,
		)
	}

	return &location, nil
}

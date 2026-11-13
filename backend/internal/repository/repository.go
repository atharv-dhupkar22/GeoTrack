package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"geotrack/backend/internal/models"
)

// Repository provides database operations for GeoTrack.
type Repository struct {
	db *pgxpool.Pool
}

// New creates a repository using an existing PostgreSQL connection pool.
func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CreateGeofence stores a geofence and its PostGIS polygon boundary.
func (r *Repository) CreateGeofence(
	ctx context.Context,
	request models.CreateGeofenceRequest,
) (*models.Geofence, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("database repository is not initialized")
	}

	if len(request.Coordinates) < 4 {
		return nil, errors.New("polygon must contain at least four coordinate pairs")
	}

	coordinatesJSON, err := json.Marshal(request.Coordinates)
	if err != nil {
		return nil, fmt.Errorf("failed to encode geofence coordinates: %w", err)
	}

	polygonWKT, err := buildPolygonWKT(request.Coordinates)
	if err != nil {
		return nil, err
	}

	id, err := generateID("geo")
	if err != nil {
		return nil, fmt.Errorf("failed to generate geofence ID: %w", err)
	}

	query := `
		INSERT INTO geofences (
			id, name, description, coordinates,
			boundary, category, status
		)
		VALUES (
			$1, $2, $3, $4::jsonb,
			ST_GeomFromText($5, 4326), $6, 'active'
		)
		RETURNING id, name, description, coordinates,
		          category, status, created_at
	`

	var geofence models.Geofence
	var storedCoordinates []byte

	err = r.db.QueryRow(
		ctx,
		query,
		id,
		request.Name,
		request.Description,
		string(coordinatesJSON),
		polygonWKT,
		request.Category,
	).Scan(
		&geofence.ID,
		&geofence.Name,
		&geofence.Description,
		&storedCoordinates,
		&geofence.Category,
		&geofence.Status,
		&geofence.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create geofence: %w", err)
	}

	if err := json.Unmarshal(storedCoordinates, &geofence.Coordinates); err != nil {
		return nil, fmt.Errorf("failed to decode stored geofence coordinates: %w", err)
	}

	return &geofence, nil
}

// ListGeofences retrieves geofences.
// When category is empty, all geofences are returned.
func (r *Repository) ListGeofences(
	ctx context.Context,
	category string,
) ([]models.Geofence, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("database repository is not initialized")
	}

	query := `
		SELECT id, name, description, coordinates,
		       category, status, created_at
		FROM geofences
		WHERE ($1 = '' OR category = $1)
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve geofences: %w", err)
	}
	defer rows.Close()

	geofences := make([]models.Geofence, 0)

	for rows.Next() {
		var geofence models.Geofence
		var storedCoordinates []byte

		if err := rows.Scan(
			&geofence.ID,
			&geofence.Name,
			&geofence.Description,
			&storedCoordinates,
			&geofence.Category,
			&geofence.Status,
			&geofence.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan geofence: %w", err)
		}

		if err := json.Unmarshal(storedCoordinates, &geofence.Coordinates); err != nil {
			return nil, fmt.Errorf("failed to decode geofence coordinates: %w", err)
		}

		geofences = append(geofences, geofence)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading geofences: %w", err)
	}

	return geofences, nil
}

// CreateVehicle registers a vehicle in the database.
func (r *Repository) CreateVehicle(
	ctx context.Context,
	request models.CreateVehicleRequest,
) (*models.Vehicle, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("database repository is not initialized")
	}

	id, err := generateID("veh")
	if err != nil {
		return nil, fmt.Errorf("failed to generate vehicle ID: %w", err)
	}

	query := `
		INSERT INTO vehicles (
			id, vehicle_number, driver_name,
			vehicle_type, phone, status
		)
		VALUES ($1, $2, $3, $4, $5, 'active')
		RETURNING id, vehicle_number, driver_name,
		          vehicle_type, phone, status, created_at
	`

	var vehicle models.Vehicle

	err = r.db.QueryRow(
		ctx,
		query,
		id,
		request.VehicleNumber,
		request.DriverName,
		request.VehicleType,
		request.Phone,
	).Scan(
		&vehicle.ID,
		&vehicle.VehicleNumber,
		&vehicle.DriverName,
		&vehicle.VehicleType,
		&vehicle.Phone,
		&vehicle.Status,
		&vehicle.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle: %w", err)
	}

	return &vehicle, nil
}

// ListVehicles retrieves all registered vehicles.
func (r *Repository) ListVehicles(
	ctx context.Context,
) ([]models.Vehicle, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("database repository is not initialized")
	}

	query := `
		SELECT id, vehicle_number, driver_name,
		       vehicle_type, phone, status, created_at
		FROM vehicles
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve vehicles: %w", err)
	}
	defer rows.Close()

	vehicles := make([]models.Vehicle, 0)

	for rows.Next() {
		var vehicle models.Vehicle

		if err := rows.Scan(
			&vehicle.ID,
			&vehicle.VehicleNumber,
			&vehicle.DriverName,
			&vehicle.VehicleType,
			&vehicle.Phone,
			&vehicle.Status,
			&vehicle.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}

		vehicles = append(vehicles, vehicle)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading vehicles: %w", err)
	}

	return vehicles, nil
}

// buildPolygonWKT converts [latitude, longitude] coordinates
// into the WKT format expected by PostGIS: longitude latitude.
func buildPolygonWKT(coordinates [][]float64) (string, error) {
	if len(coordinates) < 4 {
		return "", errors.New("polygon must contain at least four coordinate pairs")
	}

	points := make([]string, 0, len(coordinates))

	for index, coordinate := range coordinates {
		if len(coordinate) != 2 {
			return "", fmt.Errorf("coordinate %d must contain latitude and longitude", index)
		}

		latitude := coordinate[0]
		longitude := coordinate[1]

		if latitude < -90 || latitude > 90 {
			return "", fmt.Errorf("coordinate %d has invalid latitude", index)
		}

		if longitude < -180 || longitude > 180 {
			return "", fmt.Errorf("coordinate %d has invalid longitude", index)
		}

		if !isFinite(latitude) || !isFinite(longitude) {
			return "", fmt.Errorf("coordinate %d contains a non-finite value", index)
		}

		lonText := strconv.FormatFloat(longitude, 'f', -1, 64)
		latText := strconv.FormatFloat(latitude, 'f', -1, 64)

		points = append(points, lonText+" "+latText)
	}

	first := coordinates[0]
	last := coordinates[len(coordinates)-1]

	if first[0] != last[0] || first[1] != last[1] {
		return "", errors.New("polygon must be closed: first and last coordinates must match")
	}

	return "POLYGON((" + strings.Join(points, ", ") + "))", nil
}

// generateID creates a random identifier with a resource prefix.
func generateID(prefix string) (string, error) {
	randomBytes := make([]byte, 16)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return prefix + "_" + hex.EncodeToString(randomBytes), nil
}

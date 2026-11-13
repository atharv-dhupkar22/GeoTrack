
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS geofences (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    description TEXT DEFAULT '',
    coordinates JSONB NOT NULL,
    boundary GEOMETRY(POLYGON, 4326) NOT NULL,
    category VARCHAR(30) NOT NULL CHECK (
        category IN (
            'delivery_zone',
            'restricted_zone',
            'toll_zone',
            'customer_area'
        )
    ),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_geofence CHECK (ST_IsValid(boundary))
);

CREATE INDEX IF NOT EXISTS idx_geofences_boundary
ON geofences USING GIST(boundary);

CREATE INDEX IF NOT EXISTS idx_geofences_category
ON geofences(category);

CREATE TABLE IF NOT EXISTS vehicles (
    id VARCHAR(50) PRIMARY KEY,
    vehicle_number VARCHAR(50) NOT NULL UNIQUE,
    driver_name VARCHAR(150) NOT NULL,
    vehicle_type VARCHAR(50) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS vehicle_locations (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id VARCHAR(50) NOT NULL REFERENCES vehicles(id),
    latitude DOUBLE PRECISION NOT NULL CHECK (latitude BETWEEN -90 AND 90),
    longitude DOUBLE PRECISION NOT NULL CHECK (longitude BETWEEN -180 AND 180),
    location GEOMETRY(POINT, 4326) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vehicle_locations_vehicle_time
ON vehicle_locations(vehicle_id, recorded_at DESC);

CREATE INDEX IF NOT EXISTS idx_vehicle_locations_point
ON vehicle_locations USING GIST(location);

CREATE TABLE IF NOT EXISTS vehicle_geofence_states (
    vehicle_id VARCHAR(50) NOT NULL REFERENCES vehicles(id),
    geofence_id VARCHAR(50) NOT NULL REFERENCES geofences(id),
    is_inside BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(vehicle_id, geofence_id)
);

CREATE TABLE IF NOT EXISTS alert_configs (
    id VARCHAR(50) PRIMARY KEY,
    geofence_id VARCHAR(50) NOT NULL REFERENCES geofences(id),
    vehicle_id VARCHAR(50) REFERENCES vehicles(id),
    event_type VARCHAR(10) NOT NULL CHECK (
        event_type IN ('entry', 'exit', 'both')
    ),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alert_configs_lookup
ON alert_configs(geofence_id, vehicle_id, status);

CREATE TABLE IF NOT EXISTS violations (
    id VARCHAR(50) PRIMARY KEY,
    vehicle_id VARCHAR(50) NOT NULL REFERENCES vehicles(id),
    geofence_id VARCHAR(50) NOT NULL REFERENCES geofences(id),
    event_type VARCHAR(10) NOT NULL CHECK (
        event_type IN ('entry', 'exit')
    ),
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_violations_vehicle_time
ON violations(vehicle_id, recorded_at DESC);

CREATE INDEX IF NOT EXISTS idx_violations_geofence_time
ON violations(geofence_id, recorded_at DESC);

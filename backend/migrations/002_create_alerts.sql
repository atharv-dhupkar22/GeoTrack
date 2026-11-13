CREATE TABLE IF NOT EXISTS alerts (
    id VARCHAR(50) PRIMARY KEY,
    violation_id VARCHAR(50) NOT NULL REFERENCES violations(id),
    alert_config_id VARCHAR(50) NOT NULL REFERENCES alert_configs(id),
    vehicle_id VARCHAR(50) NOT NULL REFERENCES vehicles(id),
    geofence_id VARCHAR(50) NOT NULL REFERENCES geofences(id),
    event_type VARCHAR(10) NOT NULL CHECK (
        event_type IN ('entry', 'exit')
    ),
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_alert_per_config_violation
        UNIQUE (violation_id, alert_config_id)
);

CREATE INDEX IF NOT EXISTS idx_alerts_vehicle_time
    ON alerts(vehicle_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_alerts_geofence_time
    ON alerts(geofence_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_alerts_created_at
    ON alerts(created_at DESC);

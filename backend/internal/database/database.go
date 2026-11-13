package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {

	// Validate database URL
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not configured")
	}

	// Parse PostgreSQL configuration
	config, err := pgxpool.ParseConfig(databaseURL)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid database configuration: %w",
			err,
		)
	}

	// Configure connection pool
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	// Create PostgreSQL connection pool
	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		config,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create database pool: %w",
			err,
		)
	}

	// Verify database connectivity
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {

		pool.Close()

		return nil, fmt.Errorf(
			"database connection failed: %w",
			err,
		)
	}

	return pool, nil
}

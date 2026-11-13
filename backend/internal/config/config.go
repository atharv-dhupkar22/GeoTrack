package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	AppName     string
	AppEnv      string
	AppPort     string
	FrontendURL string
	DatabaseURL string
}

func Load() (*Config, error) {

	// Load environment variables from .env
	if err := loadEnvFile(".env"); err != nil {
		return nil, err
	}

	cfg := &Config{
		AppName:     getEnv("APP_NAME", "GeoTrack"),
		AppEnv:      getEnv("APP_ENV", "development"),
		AppPort:     getEnv("APP_PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func loadEnvFile(filename string) error {

	file, err := os.Open(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")

		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Remove surrounding quotes
		value = strings.Trim(value, `"'`)

		if key != "" {

			// Preserve existing environment variables
			if _, exists := os.LookupEnv(key); !exists {

				if err := os.Setenv(key, value); err != nil {
					return err
				}
			}
		}
	}

	return scanner.Err()
}

package services

import (
	"context"
	"fmt"
	"strings"

	"geotrack/backend/internal/models"
	"geotrack/backend/internal/repository"
)

type ViolationService struct {
	repository *repository.Repository
}

func NewViolationService(repo *repository.Repository) *ViolationService {
	return &ViolationService{
		repository: repo,
	}
}

func (s *ViolationService) ListHistory(
	ctx context.Context,
	vehicleID string,
	geofenceID string,
	eventType string,
	page int,
	limit int,
) (models.ViolationHistoryResult, error) {
	var result models.ViolationHistoryResult

	if s == nil || s.repository == nil {
		return result, fmt.Errorf("violation service is not initialized")
	}

	vehicleID = strings.TrimSpace(vehicleID)
	geofenceID = strings.TrimSpace(geofenceID)
	eventType = strings.TrimSpace(strings.ToLower(eventType))

	if eventType != "" && eventType != "entry" && eventType != "exit" {
		return result, fmt.Errorf("event_type must be entry or exit")
	}

	if page < 1 {
		return result, fmt.Errorf("page must be at least 1")
	}

	if limit < 1 || limit > 100 {
		return result, fmt.Errorf("limit must be between 1 and 100")
	}

	violations, total, err := s.repository.ListViolations(
		ctx,
		vehicleID,
		geofenceID,
		eventType,
		page,
		limit,
	)

	if err != nil {
		return result, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	result = models.ViolationHistoryResult{
		Violations: violations,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return result, nil
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/stpnv0/SalesTracker/internal/domain"
)

type analyticsRepository interface {
	Aggregate(ctx context.Context, from, to time.Time, itemType string) (domain.AnalyticsResult, error)
	AggregateGrouped(ctx context.Context, from, to time.Time, groupBy, itemType string) ([]domain.GroupedAnalytics, error)
}

type AnalyticsService struct {
	repo analyticsRepository
}

func NewAnalyticsService(repo analyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) GetAnalytics(ctx context.Context, filter domain.AnalyticsFilter) (domain.AnalyticsResult, error) {
	if err := filter.Validate(); err != nil {
		return domain.AnalyticsResult{}, fmt.Errorf("validate analytics filter: %w", err)
	}

	result, err := s.repo.Aggregate(ctx, filter.From, filter.To, filter.Type)
	if err != nil {
		return domain.AnalyticsResult{}, err
	}

	if filter.GroupBy != "" {
		groups, err := s.repo.AggregateGrouped(ctx, filter.From, filter.To, filter.GroupBy, filter.Type)
		if err != nil {
			return domain.AnalyticsResult{}, err
		}
		result.Groups = groups
	}

	return result, nil
}

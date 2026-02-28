package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	analyticsFrom = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analyticsTo   = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
)

func newTestAnalyticsResult() domain.AnalyticsResult {
	return domain.AnalyticsResult{
		TotalSum: decimal.NewFromInt(1000),
		Avg:      decimal.NewFromInt(100),
		Count:    10,
		Median:   decimal.NewFromInt(90),
		P90:      decimal.NewFromInt(200),
	}
}

func TestAnalyticsService_GetAnalytics_Success(t *testing.T) {
	repo := newMockanalyticsRepository(t)
	svc := NewAnalyticsService(repo)

	filter := domain.AnalyticsFilter{From: analyticsFrom, To: analyticsTo}
	expected := newTestAnalyticsResult()

	repo.EXPECT().Aggregate(mock.Anything, analyticsFrom, analyticsTo, "").Return(expected, nil)

	result, err := svc.GetAnalytics(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, expected.Count, result.Count)
	assert.True(t, expected.TotalSum.Equal(result.TotalSum))
}

func TestAnalyticsService_GetAnalytics_WithGroupBy(t *testing.T) {
	repo := newMockanalyticsRepository(t)
	svc := NewAnalyticsService(repo)

	filter := domain.AnalyticsFilter{
		From:    analyticsFrom,
		To:      analyticsTo,
		GroupBy: domain.GroupByMonth,
	}

	aggregateResult := newTestAnalyticsResult()
	groups := []domain.GroupedAnalytics{
		{Key: "2024-01", TotalSum: decimal.NewFromInt(500), Count: 5},
		{Key: "2024-02", TotalSum: decimal.NewFromInt(500), Count: 5},
	}

	repo.EXPECT().Aggregate(mock.Anything, analyticsFrom, analyticsTo, "").Return(aggregateResult, nil)
	repo.EXPECT().AggregateGrouped(mock.Anything, analyticsFrom, analyticsTo, domain.GroupByMonth, "").Return(groups, nil)

	result, err := svc.GetAnalytics(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, aggregateResult.Count, result.Count)
	assert.Len(t, result.Groups, 2)
	assert.Equal(t, "2024-01", result.Groups[0].Key)
}

func TestAnalyticsService_GetAnalytics_InvalidFilter(t *testing.T) {
	repo := newMockanalyticsRepository(t)
	svc := NewAnalyticsService(repo)

	filter := domain.AnalyticsFilter{} // zero dates → invalid

	_, err := svc.GetAnalytics(context.Background(), filter)
	assert.Error(t, err)
	assert.True(t, domain.IsValidationError(err))
}

func TestAnalyticsService_GetAnalytics_AggregateError(t *testing.T) {
	repo := newMockanalyticsRepository(t)
	svc := NewAnalyticsService(repo)

	filter := domain.AnalyticsFilter{From: analyticsFrom, To: analyticsTo}
	dbErr := errors.New("database error")

	repo.EXPECT().Aggregate(mock.Anything, analyticsFrom, analyticsTo, "").Return(domain.AnalyticsResult{}, dbErr)

	_, err := svc.GetAnalytics(context.Background(), filter)
	assert.ErrorIs(t, err, dbErr)
}

func TestAnalyticsService_GetAnalytics_GroupedError(t *testing.T) {
	repo := newMockanalyticsRepository(t)
	svc := NewAnalyticsService(repo)

	filter := domain.AnalyticsFilter{
		From:    analyticsFrom,
		To:      analyticsTo,
		GroupBy: domain.GroupByDay,
	}
	dbErr := errors.New("grouped query failed")

	repo.EXPECT().Aggregate(mock.Anything, analyticsFrom, analyticsTo, "").Return(newTestAnalyticsResult(), nil)
	repo.EXPECT().AggregateGrouped(mock.Anything, analyticsFrom, analyticsTo, domain.GroupByDay, "").Return(nil, dbErr)

	_, err := svc.GetAnalytics(context.Background(), filter)
	assert.ErrorIs(t, err, dbErr)
}

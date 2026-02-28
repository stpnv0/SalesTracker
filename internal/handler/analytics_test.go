package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupAnalyticsRouter(h *AnalyticsHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/analytics", gin.HandlerFunc(h.Get))
	return r
}

func TestAnalyticsHandler_Get_Success(t *testing.T) {
	svc := newMockanalyticsService(t)
	h := NewAnalyticsHandler(svc, newTestLogger(t))
	router := setupAnalyticsRouter(h)

	result := domain.AnalyticsResult{
		TotalSum: decimal.NewFromInt(1000),
		Avg:      decimal.NewFromInt(100),
		Count:    10,
		Median:   decimal.NewFromInt(90),
		P90:      decimal.NewFromInt(200),
	}
	svc.EXPECT().GetAnalytics(mock.Anything, mock.Anything).Return(result, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics?from=2024-01-01&to=2024-12-31", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.AnalyticsResult
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int64(10), resp.Count)
}

func TestAnalyticsHandler_Get_MissingFrom(t *testing.T) {
	svc := newMockanalyticsService(t)
	h := NewAnalyticsHandler(svc, newTestLogger(t))
	router := setupAnalyticsRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics?to=2024-12-31", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyticsHandler_Get_MissingTo(t *testing.T) {
	svc := newMockanalyticsService(t)
	h := NewAnalyticsHandler(svc, newTestLogger(t))
	router := setupAnalyticsRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics?from=2024-01-01", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyticsHandler_Get_InvalidDateFormat(t *testing.T) {
	svc := newMockanalyticsService(t)
	h := NewAnalyticsHandler(svc, newTestLogger(t))
	router := setupAnalyticsRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics?from=bad&to=2024-12-31", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyticsHandler_Get_ValidationError(t *testing.T) {
	svc := newMockanalyticsService(t)
	h := NewAnalyticsHandler(svc, newTestLogger(t))
	router := setupAnalyticsRouter(h)

	valErr := fmt.Errorf("validate analytics filter: %w", domain.ErrInvalidType)
	svc.EXPECT().GetAnalytics(mock.Anything, mock.Anything).Return(domain.AnalyticsResult{}, valErr)

	req := httptest.NewRequest(http.MethodGet, "/api/analytics?from=2024-01-01&to=2024-12-31&type=bad", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyticsHandler_Get_ServiceError(t *testing.T) {
	svc := newMockanalyticsService(t)
	h := NewAnalyticsHandler(svc, newTestLogger(t))
	router := setupAnalyticsRouter(h)

	svc.EXPECT().GetAnalytics(mock.Anything, mock.Anything).Return(domain.AnalyticsResult{}, fmt.Errorf("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/analytics?from=2024-01-01&to=2024-12-31", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

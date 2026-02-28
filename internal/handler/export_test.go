package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupExportRouter(h *ExportHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/export/csv", gin.HandlerFunc(h.CSV))
	return r
}

func TestExportHandler_CSV_Success(t *testing.T) {
	svc := newMockexportItemService(t)
	h := NewExportHandler(svc, newTestLogger(t))
	router := setupExportRouter(h)

	items := []domain.Item{
		{
			ID:          "id-1",
			Type:        domain.TypeIncome,
			Amount:      decimal.NewFromInt(100),
			Category:    "salary",
			Description: "monthly",
			Date:        time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
		},
	}
	svc.EXPECT().List(mock.Anything, mock.Anything).Return(items, int64(1), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/export/csv", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "items.csv")

	body := w.Body.String()
	lines := strings.Split(strings.TrimSpace(body), "\n")
	assert.Len(t, lines, 2)
	assert.Equal(t, "id,type,amount,category,description,date,created_at,updated_at", lines[0])
	assert.Contains(t, lines[1], "id-1")
}

func TestExportHandler_CSV_WithFilters(t *testing.T) {
	svc := newMockexportItemService(t)
	h := NewExportHandler(svc, newTestLogger(t))
	router := setupExportRouter(h)

	svc.EXPECT().List(mock.Anything, mock.MatchedBy(func(f domain.ItemFilter) bool {
		return f.Category == "food" && f.Type == "expense" && f.NoLimit && f.From != nil && f.To != nil
	})).Return([]domain.Item{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/export/csv?from=2024-01-01&to=2024-12-31&category=food&type=expense", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestExportHandler_CSV_ServiceError(t *testing.T) {
	svc := newMockexportItemService(t)
	h := NewExportHandler(svc, newTestLogger(t))
	router := setupExportRouter(h)

	svc.EXPECT().List(mock.Anything, mock.Anything).Return(nil, int64(0), fmt.Errorf("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/export/csv", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExportHandler_CSV_EmptyResult(t *testing.T) {
	svc := newMockexportItemService(t)
	h := NewExportHandler(svc, newTestLogger(t))
	router := setupExportRouter(h)

	svc.EXPECT().List(mock.Anything, mock.Anything).Return([]domain.Item{}, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/export/csv", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := w.Body.String()
	lines := strings.Split(strings.TrimSpace(body), "\n")
	assert.Len(t, lines, 1) // header only
	assert.Equal(t, "id,type,amount,category,description,date,created_at,updated_at", lines[0])
}

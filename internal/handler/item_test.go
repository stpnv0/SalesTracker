package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wb-go/wbf/logger"
)

func setupItemRouter(h *ItemHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/items", gin.HandlerFunc(h.Create))
	r.GET("/api/items", gin.HandlerFunc(h.List))
	r.GET("/api/items/:id", gin.HandlerFunc(h.GetByID))
	r.PUT("/api/items/:id", gin.HandlerFunc(h.Update))
	r.DELETE("/api/items/:id", gin.HandlerFunc(h.Delete))
	return r
}

func newTestLogger(t *testing.T) logger.Logger {
	t.Helper()
	log, err := logger.InitLogger(logger.SlogEngine, "test", "test")
	require.NoError(t, err)
	return log
}

func testItemID() string {
	return "550e8400-e29b-41d4-a716-446655440000"
}

func testItem() domain.Item {
	return domain.Item{
		ID:          testItemID(),
		Type:        domain.TypeIncome,
		Amount:      decimal.NewFromInt(100),
		Category:    "salary",
		Description: "monthly salary",
		Date:        time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestItemHandler_Create_Success(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	created := testItem()
	svc.EXPECT().Create(mock.Anything, mock.Anything).Return(created, nil)

	body := `{"type":"income","amount":100,"category":"salary","description":"monthly salary","date":"2024-06-15"}`
	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp domain.Item
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, created.ID, resp.ID)
}

func TestItemHandler_Create_InvalidJSON(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_Create_ValidationError(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	body := `{"type":"bad","amount":100,"category":"salary","date":"2024-06-15"}`
	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_Create_ServiceError(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	svc.EXPECT().Create(mock.Anything, mock.Anything).Return(domain.Item{}, fmt.Errorf("internal error"))

	body := `{"type":"income","amount":100,"category":"salary","date":"2024-06-15"}`
	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestItemHandler_List_Success(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	items := []domain.Item{testItem()}
	svc.EXPECT().List(mock.Anything, mock.Anything).Return(items, int64(1), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(1), resp["total_count"])
	assert.Len(t, resp["items"], 1)
}

func TestItemHandler_List_EmptyResult(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	svc.EXPECT().List(mock.Anything, mock.Anything).Return(nil, int64(0), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["total_count"])
	items := resp["items"].([]interface{})
	assert.Empty(t, items)
}

func TestItemHandler_List_InvalidFilter(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/items?from=bad-date", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_List_ValidationError(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	valErr := fmt.Errorf("validate filter: %w", domain.ErrInvalidType)
	svc.EXPECT().List(mock.Anything, mock.Anything).Return(nil, int64(0), valErr)

	req := httptest.NewRequest(http.MethodGet, "/api/items?type=bad", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_GetByID_Success(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	item := testItem()
	svc.EXPECT().GetByID(mock.Anything, testItemID()).Return(item, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/items/"+testItemID(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.Item
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, item.ID, resp.ID)
}

func TestItemHandler_GetByID_NotFound(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	svc.EXPECT().GetByID(mock.Anything, testItemID()).Return(domain.Item{}, domain.ErrItemNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/items/"+testItemID(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestItemHandler_GetByID_InvalidID(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	svc.EXPECT().GetByID(mock.Anything, "bad-id").Return(domain.Item{}, domain.ErrInvalidID)

	req := httptest.NewRequest(http.MethodGet, "/api/items/bad-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_Update_Success(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	updated := testItem()
	updated.Amount = decimal.NewFromInt(200)
	svc.EXPECT().Update(mock.Anything, mock.Anything).Return(updated, nil)

	body := `{"type":"income","amount":200,"category":"salary","description":"updated","date":"2024-06-15"}`
	req := httptest.NewRequest(http.MethodPut, "/api/items/"+testItemID(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestItemHandler_Update_NotFound(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	svc.EXPECT().Update(mock.Anything, mock.Anything).Return(domain.Item{}, domain.ErrItemNotFound)

	body := `{"type":"income","amount":100,"category":"salary","date":"2024-06-15"}`
	req := httptest.NewRequest(http.MethodPut, "/api/items/"+testItemID(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestItemHandler_Update_InvalidJSON(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	req := httptest.NewRequest(http.MethodPut, "/api/items/"+testItemID(), bytes.NewBufferString(`{bad}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestItemHandler_Delete_Success(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	svc.EXPECT().Delete(mock.Anything, testItemID()).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/items/"+testItemID(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestItemHandler_Delete_NotFound(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	svc.EXPECT().Delete(mock.Anything, testItemID()).Return(domain.ErrItemNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/api/items/"+testItemID(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestItemHandler_Delete_InvalidID(t *testing.T) {
	svc := newMockitemService(t)
	h := NewItemHandler(svc, newTestLogger(t))
	router := setupItemRouter(h)

	svc.EXPECT().Delete(mock.Anything, "bad-id").Return(domain.ErrInvalidID)

	req := httptest.NewRequest(http.MethodDelete, "/api/items/bad-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

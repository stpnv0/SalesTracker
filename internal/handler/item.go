package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/logger"
)

type itemService interface {
	Create(ctx context.Context, item domain.Item) (domain.Item, error)
	GetByID(ctx context.Context, id string) (domain.Item, error)
	List(ctx context.Context, filter domain.ItemFilter) ([]domain.Item, int64, error)
	Update(ctx context.Context, item domain.Item) (domain.Item, error)
	Delete(ctx context.Context, id string) error
}

type ItemHandler struct {
	svc itemService
	log logger.Logger
}

func NewItemHandler(svc itemService, log logger.Logger) *ItemHandler {
	return &ItemHandler{
		svc: svc,
		log: log,
	}
}

// Create - POST /api/items.
func (h *ItemHandler) Create(c *ginext.Context) {
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := req.ToItem()
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.svc.Create(c.Request.Context(), item)
	if err != nil {
		h.log.LogAttrs(c.Request.Context(), logger.ErrorLevel, "create item",
			logger.String("error", err.Error()))
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(c, http.StatusCreated, created)
}

// List - GET /api/items.
func (h *ItemHandler) List(c *ginext.Context) {
	filter, err := parseItemFilter(c)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	items, totalCount, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		if domain.IsValidationError(err) {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		h.log.LogAttrs(c.Request.Context(), logger.ErrorLevel, "list items",
			logger.String("error", err.Error()))
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if items == nil {
		items = []domain.Item{}
	}
	response := map[string]interface{}{
		"items":       items,
		"total_count": totalCount,
	}
	respondJSON(c, http.StatusOK, response)
}

// GetByID — GET /api/items/:id.
func (h *ItemHandler) GetByID(c *ginext.Context) {
	id := c.Param("id")

	item, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			respondError(c, http.StatusNotFound, "item not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidID) {
			respondError(c, http.StatusBadRequest, "invalid item id")
			return
		}
		h.log.LogAttrs(c.Request.Context(), logger.ErrorLevel, "get item by id",
			logger.String("error", err.Error()))
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(c, http.StatusOK, item)
}

// Update - PUT /api/items/:id.
func (h *ItemHandler) Update(c *ginext.Context) {
	id := c.Param("id")

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := req.ToItem(id)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.svc.Update(c.Request.Context(), item)
	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			respondError(c, http.StatusNotFound, "item not found")
			return
		}
		if domain.IsValidationError(err) {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		h.log.LogAttrs(c.Request.Context(), logger.ErrorLevel, "update item",
			logger.String("error", err.Error()))
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(c, http.StatusOK, updated)
}

// Delete - DELETE /api/items/:id.
func (h *ItemHandler) Delete(c *ginext.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			respondError(c, http.StatusNotFound, "item not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidID) {
			respondError(c, http.StatusBadRequest, "invalid item id")
			return
		}
		h.log.LogAttrs(c.Request.Context(), logger.ErrorLevel, "delete item",
			logger.String("error", err.Error()))
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	respondNoContent(c)
}

func parseItemFilter(c *ginext.Context) (domain.ItemFilter, error) {
	var filter domain.ItemFilter

	if v := c.Query("from"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return filter, errors.New("invalid 'from' date format, expected YYYY-MM-DD")
		}
		filter.From = &t
	}
	if v := c.Query("to"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return filter, errors.New("invalid 'to' date format, expected YYYY-MM-DD")
		}
		filter.To = &t
	}
	filter.Category = c.Query("category")
	filter.Type = c.Query("type")
	filter.SortBy = c.Query("sort_by")
	filter.Order = c.Query("order")

	if v := c.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return filter, errors.New("invalid 'limit' parameter")
		}
		filter.Limit = n
	}
	if v := c.Query("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return filter, errors.New("invalid 'offset' parameter")
		}
		filter.Offset = n
	}

	return filter, nil
}

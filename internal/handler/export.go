package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/stpnv0/SalesTracker/internal/export"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/logger"
)

type exportItemService interface {
	List(ctx context.Context, filter domain.ItemFilter) ([]domain.Item, int64, error)
}

type ExportHandler struct {
	svc exportItemService
	log logger.Logger
}

func NewExportHandler(svc exportItemService, log logger.Logger) *ExportHandler {
	return &ExportHandler{
		svc: svc,
		log: log,
	}
}

// CSV - GET /api/export/csv.
func (h *ExportHandler) CSV(c *ginext.Context) {
	filter, err := parseExportFilter(c)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	items, _, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		h.log.LogAttrs(c.Request.Context(), logger.ErrorLevel, "export csv",
			logger.String("error", err.Error()))
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=items.csv")

	if err := export.WriteCSV(c.Writer, items); err != nil {
		h.log.LogAttrs(c.Request.Context(), logger.ErrorLevel, "write csv",
			logger.String("error", err.Error()))
	}
}
func parseExportFilter(c *ginext.Context) (domain.ItemFilter, error) {
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
	filter.SortBy = domain.SortByDate
	filter.Order = domain.OrderDesc
	filter.NoLimit = true

	return filter, nil
}

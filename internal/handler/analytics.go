package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/logger"
)

type analyticsService interface {
	GetAnalytics(ctx context.Context, filter domain.AnalyticsFilter) (domain.AnalyticsResult, error)
}

type AnalyticsHandler struct {
	svc analyticsService
	log logger.Logger
}

func NewAnalyticsHandler(svc analyticsService, log logger.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		svc: svc,
		log: log,
	}
}

// Get - GET /api/analytics.
func (h *AnalyticsHandler) Get(c *ginext.Context) {
	filter, err := parseAnalyticsFilter(c)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.svc.GetAnalytics(c.Request.Context(), filter)
	if err != nil {
		if domain.IsValidationError(err) {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		h.log.LogAttrs(c.Request.Context(), logger.ErrorLevel, "get analytics",
			logger.String("error", err.Error()))
		respondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(c, http.StatusOK, result)
}

func parseAnalyticsFilter(c *ginext.Context) (domain.AnalyticsFilter, error) {
	var filter domain.AnalyticsFilter

	fromStr := c.Query("from")
	if fromStr == "" {
		return filter, errors.New("'from' parameter is required")
	}
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return filter, errors.New("invalid 'from' date format, expected YYYY-MM-DD")
	}
	filter.From = from

	toStr := c.Query("to")
	if toStr == "" {
		return filter, errors.New("'to' parameter is required")
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return filter, errors.New("invalid 'to' date format, expected YYYY-MM-DD")
	}
	filter.To = to

	filter.GroupBy = c.Query("group_by")
	filter.Type = c.Query("type")

	return filter, nil
}

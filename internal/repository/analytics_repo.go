package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

type AnalyticsRepo struct {
	db       *dbpg.DB
	strategy retry.Strategy
}

// NewAnalyticsRepo creates a new analytics repository.
func NewAnalyticsRepo(db *dbpg.DB, strategy retry.Strategy) *AnalyticsRepo {
	return &AnalyticsRepo{
		db:       db,
		strategy: strategy,
	}
}

func buildAnalyticsWhere(from, to time.Time, itemType string) (string, []interface{}) {
	clauses := []string{"date >= $1", "date <= $2"}
	args := []interface{}{from, to}

	if itemType != "" {
		clauses = append(clauses, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, itemType)
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}

var allowedGroupBy = map[string]struct {
	selectExpr string
	groupExpr  string
	orderExpr  string
}{
	"day":      {"date::text", "date", "date"},
	"week":     {"DATE_TRUNC('week', date)::date::text", "DATE_TRUNC('week', date)", "DATE_TRUNC('week', date)"},
	"month":    {"DATE_TRUNC('month', date)::date::text", "DATE_TRUNC('month', date)", "DATE_TRUNC('month', date)"},
	"category": {"category", "category", "category"},
}

func (r *AnalyticsRepo) Aggregate(ctx context.Context, from, to time.Time, itemType string) (domain.AnalyticsResult, error) {
	where, args := buildAnalyticsWhere(from, to, itemType)

	query := fmt.Sprintf(`
		SELECT
			COUNT(*)                                                         AS count,
			COALESCE(SUM(amount), 0)                                         AS total_sum,
			COALESCE(AVG(amount), 0)                                         AS avg,
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0) AS median,
			COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0) AS p90
		FROM items %s`, where)

	row, err := r.db.QueryRowWithRetry(ctx, r.strategy, query, args...)
	if err != nil {
		return domain.AnalyticsResult{}, fmt.Errorf("aggregate analytics: %w", err)
	}

	var res domain.AnalyticsResult
	if err = row.Scan(&res.Count, &res.TotalSum, &res.Avg, &res.Median, &res.P90); err != nil {
		return domain.AnalyticsResult{}, fmt.Errorf("scan analytics: %w", err)
	}

	return res, nil
}

func (r *AnalyticsRepo) AggregateGrouped(ctx context.Context, from, to time.Time, groupBy, itemType string) ([]domain.GroupedAnalytics, error) {
	gb, ok := allowedGroupBy[groupBy]
	if !ok {
		return nil, fmt.Errorf("unsupported group_by value: %q", groupBy)
	}

	where, args := buildAnalyticsWhere(from, to, itemType)

	query := fmt.Sprintf(`
		SELECT
			%s                                                               AS key,
			COUNT(*)                                                         AS count,
			COALESCE(SUM(amount), 0)                                         AS total_sum,
			COALESCE(AVG(amount), 0)                                         AS avg,
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0) AS median,
			COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0) AS p90
		FROM items %s
		GROUP BY %s
		ORDER BY %s`,
		gb.selectExpr, where, gb.groupExpr, gb.orderExpr)

	rows, err := r.db.QueryWithRetry(ctx, r.strategy, query, args...)
	if err != nil {
		return nil, fmt.Errorf("aggregate grouped analytics: %w", err)
	}
	defer rows.Close()

	var res []domain.GroupedAnalytics
	for rows.Next() {
		var g domain.GroupedAnalytics
		if err = rows.Scan(&g.Key, &g.Count, &g.TotalSum, &g.Avg, &g.Median, &g.P90); err != nil {
			return nil, fmt.Errorf("scan grouped analytics: %w", err)
		}
		res = append(res, g)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return res, nil
}

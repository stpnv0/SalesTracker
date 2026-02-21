package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	TypeIncome  = "income"
	TypeExpense = "expense"
)

const (
	SortByDate     = "date"
	SortByAmount   = "amount"
	SortByCategory = "category"
	SortByType     = "type"
)

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

const (
	GroupByDay      = "day"
	GroupByWeek     = "week"
	GroupByMonth    = "month"
	GroupByCategory = "category"
)

type Item struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Amount      decimal.Decimal `json:"amount"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type AnalyticsResult struct {
	TotalSum decimal.Decimal    `json:"total_sum"`
	Avg      decimal.Decimal    `json:"avg"`
	Count    int64              `json:"count"`
	Median   decimal.Decimal    `json:"median"`
	P90      decimal.Decimal    `json:"p90"`
	Groups   []GroupedAnalytics `json:"groups,omitempty"`
}

type GroupedAnalytics struct {
	Key      string          `json:"key"`
	TotalSum decimal.Decimal `json:"total_sum"`
	Avg      decimal.Decimal `json:"avg"`
	Count    int64           `json:"count"`
	Median   decimal.Decimal `json:"median"`
	P90      decimal.Decimal `json:"p90"`
}

package domain

import "time"

type ItemFilter struct {
	From     *time.Time
	To       *time.Time
	Category string
	Type     string
	SortBy   string
	Order    string
	Limit    int
	Offset   int
	NoLimit  bool // true для экспорта — отключает пагинацию
}

func (f ItemFilter) Validate() error {
	if f.Type != "" && f.Type != TypeIncome && f.Type != TypeExpense {
		return ErrInvalidType
	}
	if f.SortBy != "" {
		switch f.SortBy {
		case SortByDate, SortByAmount, SortByCategory, SortByType:
		default:
			return ErrInvalidSortBy
		}
	}
	if f.Order != "" && f.Order != OrderAsc && f.Order != OrderDesc {
		return ErrInvalidOrder
	}
	if f.From != nil && f.To != nil && f.From.After(*f.To) {
		return ErrInvalidDateRange
	}
	return nil
}

type AnalyticsFilter struct {
	From    time.Time
	To      time.Time
	GroupBy string
	Type    string
}

func (f AnalyticsFilter) Validate() error {
	if f.From.IsZero() || f.To.IsZero() {
		return ErrInvalidDate
	}
	if f.From.After(f.To) {
		return ErrInvalidDateRange
	}
	if f.GroupBy != "" {
		switch f.GroupBy {
		case GroupByDay, GroupByWeek, GroupByMonth, GroupByCategory:
		default:
			return ErrInvalidGroupBy
		}
	}
	if f.Type != "" && f.Type != TypeIncome && f.Type != TypeExpense {
		return ErrInvalidType
	}
	return nil
}

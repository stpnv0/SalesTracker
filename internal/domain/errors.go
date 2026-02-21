package domain

import "errors"

var (
	ErrInvalidType      = errors.New("type must be 'income' or 'expense'")
	ErrInvalidAmount    = errors.New("amount must be greater than zero")
	ErrEmptyCategory    = errors.New("category must not be empty")
	ErrInvalidDate      = errors.New("date must not be zero")
	ErrInvalidID        = errors.New("id must be a valid UUID")
	ErrItemNotFound     = errors.New("item not found")
	ErrInvalidSortBy    = errors.New("sort_by must be one of: date, amount, category, type")
	ErrInvalidOrder     = errors.New("order must be 'asc' or 'desc'")
	ErrInvalidGroupBy   = errors.New("group_by must be one of: day, week, month, category")
	ErrInvalidDateRange = errors.New("'from' date must not be after 'to' date")
	ErrValidation       = errors.New("validation error")
)

var validationErrors = []error{
	ErrValidation,
	ErrInvalidType,
	ErrInvalidAmount,
	ErrEmptyCategory,
	ErrInvalidDate,
	ErrInvalidID,
	ErrInvalidSortBy,
	ErrInvalidOrder,
	ErrInvalidGroupBy,
	ErrInvalidDateRange,
}

func IsValidationError(err error) bool {
	for _, ve := range validationErrors {
		if errors.Is(err, ve) {
			return true
		}
	}
	return false
}

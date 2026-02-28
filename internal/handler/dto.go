package handler

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/stpnv0/SalesTracker/internal/domain"
)

var validate = validator.New()

func formatValidationErrors(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fmt.Errorf("%w: %s", domain.ErrValidation, err.Error())
	}

	for _, fe := range validationErrors {
		switch fe.Tag() {
		case "required":
			return fmt.Errorf("%w: %s is required", domain.ErrValidation, fe.Field())
		case "oneof":
			return fmt.Errorf("%w: %s must be one of [%s]", domain.ErrValidation, fe.Field(), fe.Param())
		case "gt":
			return fmt.Errorf("%w: %s must be greater than %s", domain.ErrValidation, fe.Field(), fe.Param())
		case "datetime":
			return fmt.Errorf("%w: %s must be a valid date in format %s", domain.ErrValidation, fe.Field(), fe.Param())
		case "max":
			return fmt.Errorf("%w: %s must be at most %s characters", domain.ErrValidation, fe.Field(), fe.Param())
		default:
			return fmt.Errorf("%w: %s failed on '%s' check", domain.ErrValidation, fe.Field(), fe.Tag())
		}
	}

	return fmt.Errorf("%w: unknown validation error", domain.ErrValidation)
}

type CreateItemRequest struct {
	Type        string          `json:"type"        validate:"required,oneof=income expense"`
	Amount      decimal.Decimal `json:"amount"`
	Category    string          `json:"category"    validate:"required,max=100"`
	Description string          `json:"description" validate:"max=1000"`
	Date        string          `json:"date"        validate:"required,datetime=2006-01-02"`
}

func (r CreateItemRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return formatValidationErrors(err)
	}
	if !r.Amount.IsPositive() {
		return fmt.Errorf("%w: Amount must be greater than 0", domain.ErrValidation)
	}
	return nil
}

func (r CreateItemRequest) ToItem() (domain.Item, error) {
	date, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return domain.Item{}, fmt.Errorf("%w: invalid date format", domain.ErrValidation)
	}

	now := time.Now().UTC()
	return domain.Item{
		Type:        r.Type,
		Amount:      r.Amount,
		Category:    r.Category,
		Description: r.Description,
		Date:        date,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

type UpdateItemRequest struct {
	Type        string          `json:"type"        validate:"required,oneof=income expense"`
	Amount      decimal.Decimal `json:"amount"`
	Category    string          `json:"category"    validate:"required,max=100"`
	Description string          `json:"description" validate:"max=1000"`
	Date        string          `json:"date"        validate:"required,datetime=2006-01-02"`
}

func (r UpdateItemRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return formatValidationErrors(err)
	}
	if !r.Amount.IsPositive() {
		return fmt.Errorf("%w: Amount must be greater than 0", domain.ErrValidation)
	}
	return nil
}

func (r UpdateItemRequest) ToItem(id string) (domain.Item, error) {
	date, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return domain.Item{}, fmt.Errorf("%w: invalid date format", domain.ErrValidation)
	}

	return domain.Item{
		ID:          id,
		Type:        r.Type,
		Amount:      r.Amount,
		Category:    r.Category,
		Description: r.Description,
		Date:        date,
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

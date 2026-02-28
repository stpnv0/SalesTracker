package handler

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateItemRequest_Validate_Valid(t *testing.T) {
	req := CreateItemRequest{
		Type:        "income",
		Amount:      decimal.NewFromInt(100),
		Category:    "salary",
		Description: "monthly salary",
		Date:        "2024-06-15",
	}
	err := req.Validate()
	assert.NoError(t, err)
}

func TestCreateItemRequest_Validate_MissingRequired(t *testing.T) {
	tests := []struct {
		name string
		req  CreateItemRequest
	}{
		{
			name: "missing type",
			req: CreateItemRequest{
				Amount:   decimal.NewFromInt(100),
				Category: "salary",
				Date:     "2024-06-15",
			},
		},
		{
			name: "missing category",
			req: CreateItemRequest{
				Type:   "income",
				Amount: decimal.NewFromInt(100),
				Date:   "2024-06-15",
			},
		},
		{
			name: "missing date",
			req: CreateItemRequest{
				Type:     "income",
				Amount:   decimal.NewFromInt(100),
				Category: "salary",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.Error(t, err)
			assert.ErrorIs(t, err, domain.ErrValidation)
		})
	}
}

func TestCreateItemRequest_Validate_InvalidType(t *testing.T) {
	req := CreateItemRequest{
		Type:     "bad",
		Amount:   decimal.NewFromInt(100),
		Category: "salary",
		Date:     "2024-06-15",
	}
	err := req.Validate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestCreateItemRequest_Validate_ZeroAmount(t *testing.T) {
	req := CreateItemRequest{
		Type:     "income",
		Amount:   decimal.Decimal{},
		Category: "salary",
		Date:     "2024-06-15",
	}
	err := req.Validate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestCreateItemRequest_Validate_NegativeAmount(t *testing.T) {
	req := CreateItemRequest{
		Type:     "income",
		Amount:   decimal.NewFromInt(-50),
		Category: "salary",
		Date:     "2024-06-15",
	}
	err := req.Validate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestCreateItemRequest_Validate_InvalidDate(t *testing.T) {
	req := CreateItemRequest{
		Type:     "income",
		Amount:   decimal.NewFromInt(100),
		Category: "salary",
		Date:     "15-06-2024",
	}
	err := req.Validate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
}

func TestCreateItemRequest_ToItem(t *testing.T) {
	req := CreateItemRequest{
		Type:        "expense",
		Amount:      decimal.NewFromFloat(55.50),
		Category:    "food",
		Description: "lunch",
		Date:        "2024-06-15",
	}

	item, err := req.ToItem()
	require.NoError(t, err)
	assert.Equal(t, "expense", item.Type)
	assert.True(t, decimal.NewFromFloat(55.50).Equal(item.Amount))
	assert.Equal(t, "food", item.Category)
	assert.Equal(t, "lunch", item.Description)
	assert.Equal(t, 2024, item.Date.Year())
	assert.Equal(t, 6, int(item.Date.Month()))
	assert.Equal(t, 15, item.Date.Day())
	assert.False(t, item.CreatedAt.IsZero())
	assert.False(t, item.UpdatedAt.IsZero())
}

func TestUpdateItemRequest_Validate_Valid(t *testing.T) {
	req := UpdateItemRequest{
		Type:        "expense",
		Amount:      decimal.NewFromInt(50),
		Category:    "food",
		Description: "dinner",
		Date:        "2024-07-20",
	}
	err := req.Validate()
	assert.NoError(t, err)
}

func TestUpdateItemRequest_ToItem(t *testing.T) {
	req := UpdateItemRequest{
		Type:        "income",
		Amount:      decimal.NewFromInt(200),
		Category:    "freelance",
		Description: "project",
		Date:        "2024-08-10",
	}
	id := "550e8400-e29b-41d4-a716-446655440000"

	item, err := req.ToItem(id)
	require.NoError(t, err)
	assert.Equal(t, id, item.ID)
	assert.Equal(t, "income", item.Type)
	assert.True(t, decimal.NewFromInt(200).Equal(item.Amount))
	assert.Equal(t, "freelance", item.Category)
	assert.False(t, item.UpdatedAt.IsZero())
}

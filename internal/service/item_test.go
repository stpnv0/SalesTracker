package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var validUUID = "550e8400-e29b-41d4-a716-446655440000"

func newTestItem() domain.Item {
	return domain.Item{
		ID:          validUUID,
		Type:        domain.TypeIncome,
		Amount:      decimal.NewFromInt(100),
		Category:    "salary",
		Description: "monthly salary",
		Date:        time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func TestItemService_Create_Success(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	input := newTestItem()
	input.ID = ""
	expected := newTestItem()

	repo.EXPECT().Create(mock.Anything, input).Return(expected, nil)

	result, err := svc.Create(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

func TestItemService_Create_RepoError(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	input := newTestItem()
	repoErr := errors.New("db connection failed")

	repo.EXPECT().Create(mock.Anything, input).Return(domain.Item{}, repoErr)

	_, err := svc.Create(context.Background(), input)
	assert.ErrorIs(t, err, repoErr)
}

func TestItemService_List_Success(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	filter := domain.ItemFilter{Type: domain.TypeIncome}
	items := []domain.Item{newTestItem()}

	repo.EXPECT().GetAll(mock.Anything, filter).Return(items, int64(1), nil)

	result, count, err := svc.List(context.Background(), filter)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), count)
}

func TestItemService_List_InvalidFilter(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	filter := domain.ItemFilter{Type: "invalid"}

	_, _, err := svc.List(context.Background(), filter)
	assert.Error(t, err)
	assert.True(t, domain.IsValidationError(err))
}

func TestItemService_GetByID_Success(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	expected := newTestItem()

	repo.EXPECT().GetByID(mock.Anything, validUUID).Return(expected, nil)

	result, err := svc.GetByID(context.Background(), validUUID)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

func TestItemService_GetByID_InvalidUUID(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	_, err := svc.GetByID(context.Background(), "not-a-uuid")
	assert.ErrorIs(t, err, domain.ErrInvalidID)
}

func TestItemService_GetByID_NotFound(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	repo.EXPECT().GetByID(mock.Anything, validUUID).Return(domain.Item{}, domain.ErrItemNotFound)

	_, err := svc.GetByID(context.Background(), validUUID)
	assert.ErrorIs(t, err, domain.ErrItemNotFound)
}

func TestItemService_Update_Success(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	input := newTestItem()
	expected := newTestItem()
	expected.Amount = decimal.NewFromInt(200)

	repo.EXPECT().Update(mock.Anything, input).Return(expected, nil)

	result, err := svc.Update(context.Background(), input)
	assert.NoError(t, err)
	assert.True(t, expected.Amount.Equal(result.Amount))
}

func TestItemService_Update_InvalidUUID(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	input := newTestItem()
	input.ID = "bad-id"

	_, err := svc.Update(context.Background(), input)
	assert.ErrorIs(t, err, domain.ErrInvalidID)
}

func TestItemService_Delete_Success(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	repo.EXPECT().Delete(mock.Anything, validUUID).Return(nil)

	err := svc.Delete(context.Background(), validUUID)
	assert.NoError(t, err)
}

func TestItemService_Delete_InvalidUUID(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	err := svc.Delete(context.Background(), "bad-id")
	assert.ErrorIs(t, err, domain.ErrInvalidID)
}

func TestItemService_Delete_NotFound(t *testing.T) {
	repo := newMockitemRepository(t)
	svc := NewItemService(repo)

	repo.EXPECT().Delete(mock.Anything, validUUID).Return(domain.ErrItemNotFound)

	err := svc.Delete(context.Background(), validUUID)
	assert.ErrorIs(t, err, domain.ErrItemNotFound)
}

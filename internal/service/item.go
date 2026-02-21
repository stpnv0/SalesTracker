package service

import (
	"context"
	"fmt"

	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/wb-go/wbf/helpers"
)

type itemRepository interface {
	Create(ctx context.Context, item domain.Item) (domain.Item, error)
	GetAll(ctx context.Context, filter domain.ItemFilter) ([]domain.Item, int64, error)
	GetByID(ctx context.Context, id string) (domain.Item, error)
	Update(ctx context.Context, item domain.Item) (domain.Item, error)
	Delete(ctx context.Context, id string) error
}

type ItemService struct {
	repo itemRepository
}

func NewItemService(repo itemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	created, err := s.repo.Create(ctx, item)
	if err != nil {
		return domain.Item{}, err
	}
	return created, nil
}

func (s *ItemService) List(ctx context.Context, filter domain.ItemFilter) ([]domain.Item, int64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, fmt.Errorf("validate filter: %w", err)
	}
	items, totalCount, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, totalCount, nil
}

func (s *ItemService) GetByID(ctx context.Context, id string) (domain.Item, error) {
	if err := helpers.ParseUUID(id); err != nil {
		return domain.Item{}, domain.ErrInvalidID
	}
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Item{}, err
	}
	return item, nil
}

func (s *ItemService) Update(ctx context.Context, item domain.Item) (domain.Item, error) {
	if err := helpers.ParseUUID(item.ID); err != nil {
		return domain.Item{}, domain.ErrInvalidID
	}

	updated, err := s.repo.Update(ctx, item)
	if err != nil {
		return domain.Item{}, err
	}
	return updated, nil
}

func (s *ItemService) Delete(ctx context.Context, id string) error {
	if err := helpers.ParseUUID(id); err != nil {
		return domain.ErrInvalidID
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

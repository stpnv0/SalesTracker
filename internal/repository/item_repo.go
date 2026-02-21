package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

const (
	defaultLimit = 50
	maxLimit     = 1000
)

type ItemRepo struct {
	db       *dbpg.DB
	strategy retry.Strategy
}

func NewItemRepo(db *dbpg.DB, strategy retry.Strategy) *ItemRepo {
	return &ItemRepo{
		db:       db,
		strategy: strategy,
	}
}

func (r *ItemRepo) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	query := `
		INSERT INTO items (type, amount, category, description, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, type, amount, category, description, date, created_at, updated_at`

	row, err := r.db.QueryRowWithRetry(
		ctx, r.strategy, query,
		item.Type, item.Amount, item.Category, item.Description,
		item.Date, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return domain.Item{}, fmt.Errorf("create item: %w", err)
	}

	var created domain.Item
	if err = row.Scan(
		&created.ID, &created.Type, &created.Amount, &created.Category,
		&created.Description, &created.Date, &created.CreatedAt, &created.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Item{}, domain.ErrItemNotFound
		}
		return domain.Item{}, fmt.Errorf("scan created item: %w", err)
	}

	return created, nil
}

func (r *ItemRepo) GetByID(ctx context.Context, id string) (domain.Item, error) {
	query := `
		SELECT id, type, amount, category, description, date, created_at, updated_at
		FROM items
		WHERE id = $1`

	row, err := r.db.QueryRowWithRetry(ctx, r.strategy, query, id)
	if err != nil {
		return domain.Item{}, fmt.Errorf("get item by id: %w", err)
	}

	var item domain.Item
	if err = row.Scan(
		&item.ID, &item.Type, &item.Amount, &item.Category,
		&item.Description, &item.Date, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Item{}, domain.ErrItemNotFound
		}
		return domain.Item{}, fmt.Errorf("scan item: %w", err)
	}

	return item, nil
}

var allowedSortColumns = map[string]string{
	domain.SortByDate:     "date",
	domain.SortByAmount:   "amount",
	domain.SortByCategory: "category",
	domain.SortByType:     "type",
}

func (r *ItemRepo) GetAll(ctx context.Context, filter domain.ItemFilter) ([]domain.Item, int64, error) {
	whereClauses := make([]string, 0, 4)
	args := make([]interface{}, 0, 4)
	argIdx := 1

	if filter.Type != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.Category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, filter.Category)
		argIdx++
	}
	if filter.From != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("date >= $%d", argIdx))
		args = append(args, *filter.From)
		argIdx++
	}
	if filter.To != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("date <= $%d", argIdx))
		args = append(args, *filter.To)
		argIdx++
	}

	where := ""
	if len(whereClauses) > 0 {
		where = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	sortCol := "date"
	if col, ok := allowedSortColumns[filter.SortBy]; ok {
		sortCol = col
	}
	sortDir := "DESC"
	if filter.Order == domain.OrderAsc {
		sortDir = "ASC"
	}

	var limitClause string
	if !filter.NoLimit {
		limit := defaultLimit
		if filter.Limit > 0 && filter.Limit <= maxLimit {
			limit = filter.Limit
		}
		offset := 0
		if filter.Offset > 0 {
			offset = filter.Offset
		}
		limitClause = fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}

	query := fmt.Sprintf(
		`SELECT 
					id, type, amount, category, description, date, created_at, updated_at,
					COUNT(*) OVER() AS total_count
			    FROM items %s
			    ORDER BY %s %s
			    %s`,
		where, sortCol, sortDir, limitClause,
	)

	rows, err := r.db.QueryWithRetry(ctx, r.strategy, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("get all items: %w", err)
	}
	defer rows.Close()

	var (
		items      []domain.Item
		totalCount int64
	)
	for rows.Next() {
		var i domain.Item
		if err = rows.Scan(
			&i.ID, &i.Type, &i.Amount, &i.Category, &i.Description,
			&i.Date, &i.CreatedAt, &i.UpdatedAt, &totalCount,
		); err != nil {
			return nil, 0, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, i)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration: %w", err)
	}

	return items, totalCount, nil
}

func (r *ItemRepo) Update(ctx context.Context, item domain.Item) (domain.Item, error) {
	query := `
		UPDATE items
		SET type = $2, amount = $3, category = $4, description = $5, date = $6, updated_at = $7
		WHERE id = $1
		RETURNING id, type, amount, category, description, date, created_at, updated_at`

	row, err := r.db.QueryRowWithRetry(ctx, r.strategy, query,
		item.ID, item.Type, item.Amount, item.Category,
		item.Description, item.Date, item.UpdatedAt,
	)
	if err != nil {
		return domain.Item{}, fmt.Errorf("update item: %w", err)
	}

	var updated domain.Item
	if err = row.Scan(
		&updated.ID, &updated.Type, &updated.Amount, &updated.Category,
		&updated.Description, &updated.Date, &updated.CreatedAt, &updated.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Item{}, domain.ErrItemNotFound
		}
		return domain.Item{}, fmt.Errorf("scan updated item: %w", err)
	}

	return updated, nil
}

func (r *ItemRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM items WHERE id = $1`

	res, err := r.db.ExecWithRetry(ctx, r.strategy, query, id)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrItemNotFound
	}

	return nil
}

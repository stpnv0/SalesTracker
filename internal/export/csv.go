package export

import (
	"encoding/csv"
	"io"
	"time"

	"github.com/stpnv0/SalesTracker/internal/domain"
)

func WriteCSV(w io.Writer, items []domain.Item) error {
	cw := csv.NewWriter(w)

	if err := cw.Write([]string{
		"id", "type", "amount", "category",
		"description", "date", "created_at", "updated_at",
	}); err != nil {
		return err
	}

	for _, item := range items {
		if err := cw.Write([]string{
			item.ID,
			item.Type,
			item.Amount.StringFixed(2),
			item.Category,
			item.Description,
			item.Date.Format("2006-01-02"),
			item.CreatedAt.Format(time.RFC3339),
			item.UpdatedAt.Format(time.RFC3339),
		}); err != nil {
			return err
		}
	}

	cw.Flush()
	return cw.Error()
}

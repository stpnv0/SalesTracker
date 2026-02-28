package export

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stpnv0/SalesTracker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDate = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
var testCreatedAt = time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)

func TestWriteCSV_Success(t *testing.T) {
	var buf bytes.Buffer
	items := []domain.Item{
		{
			ID:          "id-1",
			Type:        domain.TypeIncome,
			Amount:      decimal.NewFromFloat(150.50),
			Category:    "salary",
			Description: "monthly salary",
			Date:        testDate,
			CreatedAt:   testCreatedAt,
			UpdatedAt:   testCreatedAt,
		},
	}

	err := WriteCSV(&buf, items)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 2)
	assert.Equal(t, "id,type,amount,category,description,date,created_at,updated_at", lines[0])
	assert.Contains(t, lines[1], "id-1")
	assert.Contains(t, lines[1], "income")
	assert.Contains(t, lines[1], "150.50")
	assert.Contains(t, lines[1], "salary")
	assert.Contains(t, lines[1], "monthly salary")
	assert.Contains(t, lines[1], "2024-06-15")
}

func TestWriteCSV_EmptyItems(t *testing.T) {
	var buf bytes.Buffer

	err := WriteCSV(&buf, []domain.Item{})
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 1)
	assert.Equal(t, "id,type,amount,category,description,date,created_at,updated_at", lines[0])
}

func TestWriteCSV_MultipleItems(t *testing.T) {
	var buf bytes.Buffer
	items := []domain.Item{
		{
			ID:       "id-1",
			Type:     domain.TypeIncome,
			Amount:   decimal.NewFromInt(100),
			Category: "salary",
			Date:     testDate, CreatedAt: testCreatedAt, UpdatedAt: testCreatedAt,
		},
		{
			ID:       "id-2",
			Type:     domain.TypeExpense,
			Amount:   decimal.NewFromFloat(50.75),
			Category: "food",
			Date:     testDate, CreatedAt: testCreatedAt, UpdatedAt: testCreatedAt,
		},
		{
			ID:       "id-3",
			Type:     domain.TypeIncome,
			Amount:   decimal.NewFromInt(200),
			Category: "freelance",
			Date:     testDate, CreatedAt: testCreatedAt, UpdatedAt: testCreatedAt,
		},
	}

	err := WriteCSV(&buf, items)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 4) // header + 3 items
}

func TestWriteCSV_SpecialCharacters(t *testing.T) {
	var buf bytes.Buffer
	items := []domain.Item{
		{
			ID:          "id-1",
			Type:        domain.TypeExpense,
			Amount:      decimal.NewFromInt(50),
			Category:    "food",
			Description: `lunch with "friends", very expensive`,
			Date:        testDate, CreatedAt: testCreatedAt, UpdatedAt: testCreatedAt,
		},
	}

	err := WriteCSV(&buf, items)
	require.NoError(t, err)

	output := buf.String()
	// CSV should properly escape quotes and commas
	assert.Contains(t, output, `"lunch with ""friends"", very expensive"`)
}

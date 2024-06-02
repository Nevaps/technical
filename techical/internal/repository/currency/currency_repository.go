package repository

import (
	"context"
	"database/sql"
	"techical/internal/domain"
)

type CurrencyRepository struct {
	Conn *sql.DB
}

func NewCurrencyRepository(conn *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{conn}
}

func (m *CurrencyRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Currency, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	result = make([]domain.Currency, 0)
	for rows.Next() {
		t := domain.Currency{}

		err = rows.Scan(
			&t.Name,
			&t.Ticker,
			&t.Available,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *CurrencyRepository) GetAllAvailable(ctx context.Context) ([]domain.Currency, error) {
	query := `SELECT * FROM currencies WHERE available = TRUE`

	list, err := m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}

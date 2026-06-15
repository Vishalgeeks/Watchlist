package stock

import (
	"database/sql"

	"watchlist-backend/pkg/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(stock *models.Stock) error {
	query := `
		INSERT INTO stocks
		(symbol, company_name, exchange, current_price)
		VALUES ($1, $2, $3, $4)
		RETURNING id, last_updated
	`

	return r.db.QueryRow(
		query,
		stock.Symbol,
		stock.CompanyName,
		stock.Exchange,
		stock.LTP,
	).Scan(&stock.ID, &stock.LastUpdated)
}

func (r *Repository) GetAll() ([]models.Stock, error) {
	query := `
		SELECT id, symbol, company_name, exchange,
		       current_price, last_updated
		FROM stocks
		ORDER BY symbol
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock

	for rows.Next() {
		var s models.Stock

		err := rows.Scan(
			&s.ID,
			&s.Symbol,
			&s.CompanyName,
			&s.Exchange,
			&s.LTP,
			&s.LastUpdated,
		)

		if err != nil {
			return nil, err
		}

		stocks = append(stocks, s)
	}

	return stocks, nil
}

func (r *Repository) GetByID(id string) (*models.Stock, error) {
	query := `
		SELECT id, symbol, company_name,
		       exchange, current_price, last_updated
		FROM stocks
		WHERE id = $1
	`

	var stock models.Stock

	err := r.db.QueryRow(query, id).Scan(
		&stock.ID,
		&stock.Symbol,
		&stock.CompanyName,
		&stock.Exchange,
		&stock.LTP,
		&stock.LastUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &stock, nil
}

func (r *Repository) Update(id string, stock *models.Stock) error {
	query := `
		UPDATE stocks
		SET symbol = $1,
		    company_name = $2,
		    exchange = $3,
		    current_price = $4,
		    last_updated = NOW()
		WHERE id = $5
		RETURNING last_updated
	`

	return r.db.QueryRow(
		query,
		stock.Symbol,
		stock.CompanyName,
		stock.Exchange,
		stock.LTP,
		id,
	).Scan(&stock.LastUpdated)
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(
		"DELETE FROM stocks WHERE id = $1",
		id,
	)

	return err
}

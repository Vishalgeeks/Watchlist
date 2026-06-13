package watchlist

import (
	"database/sql"
	"errors"
	"watchlist-backend/pkg/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// jab new watchlisty bne
func (r *Repository) Create(w *models.Watchlist) error {
	query := `
		INSERT INTO watchlists (id, user_id, name, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := r.db.Exec(query, w.ID, w.UserID, w.Name)
	return err
}

// uss id ki sari watchlist showing
func (r *Repository) GetAllByUserID(userID string) ([]models.Watchlist, error) {
	query := `
		SELECT w.id, w.user_id, w.name, w.created_at,
		       COUNT(wi.id) as stock_count
		FROM watchlists w
		LEFT JOIN watchlist_items wi ON wi.watchlist_id = w.id
		WHERE w.user_id = $1
		GROUP BY w.id
		ORDER BY w.created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var watchlists []models.Watchlist
	for rows.Next() {
		var w models.Watchlist
		if err := rows.Scan(&w.ID, &w.UserID, &w.Name, &w.CreatedAt, &w.StockCount); err != nil {
			return nil, err
		}
		watchlists = append(watchlists, w)
	}
	return watchlists, nil
}

// ek specific wtchlist showing
func (r *Repository) GetByID(watchlistID string) (*models.Watchlist, error) {
	query := `SELECT id, user_id, name, created_at FROM watchlists WHERE id = $1`
	w := &models.Watchlist{}
	err := r.db.QueryRow(query, watchlistID).Scan(&w.ID, &w.UserID, &w.Name, &w.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("watchlist not found")
	}
	return w, err
}

// watchlist delete
func (r *Repository) Delete(watchlistID string) error {
	_, err := r.db.Exec(`DELETE FROM watchlists WHERE id = $1`, watchlistID)
	return err
}

// fetching stock of that watchlist
func (r *Repository) GetStocks(watchlistID string) ([]models.WatchlistItem, error) {
	query := `
		SELECT
			wi.id, wi.watchlist_id, wi.stock_id, wi.added_at,
			s.id, s.symbol, s.company_name, s.exchange, s.current_price, s.last_updated
		FROM watchlist_items wi
		JOIN stocks s ON s.id = wi.stock_id
		WHERE wi.watchlist_id = $1
		ORDER BY wi.added_at DESC
	`
	rows, err := r.db.Query(query, watchlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.WatchlistItem
	for rows.Next() {
		var item models.WatchlistItem
		var stock models.Stock
		err := rows.Scan(
			&item.ID, &item.WatchlistID, &item.StockID, &item.AddedAt,
			&stock.ID, &stock.Symbol, &stock.CompanyName,
			&stock.Exchange, &stock.CurrentPrice, &stock.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		item.Stock = &stock
		items = append(items, item)
	}
	return items, nil
}

// adding stocks to the watchlist
func (r *Repository) AddStock(watchlistID, stockID, itemID string) error {
	query := `
		INSERT INTO watchlist_items (id, watchlist_id, stock_id, added_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := r.db.Exec(query, itemID, watchlistID, stockID)
	return err
}

//deleting stocks to the watchlist

func (r *Repository) RemoveStock(watchlistID, stockID string) error {
	_, err := r.db.Exec(
		`DELETE FROM watchlist_items WHERE watchlist_id = $1 AND stock_id = $2`,
		watchlistID, stockID,
	)
	return err
}

//checking ki duplicate to ni h ni h n

func (r *Repository) StockExists(watchlistID, stockID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM watchlist_items WHERE watchlist_id = $1 AND stock_id = $2`,
		watchlistID, stockID,
	).Scan(&count)
	return count > 0, err
}

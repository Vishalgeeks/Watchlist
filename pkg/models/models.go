package models

import "time"

// ── User ──────────────────────────────────────────────
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ── Stock ─────────────────────────────────────────────
type Stock struct {
	ID           string    `json:"id"`
	Symbol       string    `json:"symbol"`
	CompanyName  string    `json:"company_name"`
	Exchange     string    `json:"exchange"`
	CurrentPrice float64   `json:"current_price"`
	LastUpdated  time.Time `json:"last_updated"`
}

// ── Watchlist ─────────────────────────────────────────
type Watchlist struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	StockCount int       `json:"stock_count"`
}

// ── Watchlist Item ────────────────────────────────────
type WatchlistItem struct {
	ID          string    `json:"id"`
	WatchlistID string    `json:"watchlist_id"`
	StockID     string    `json:"stock_id"`
	AddedAt     time.Time `json:"added_at"`
	Stock       *Stock    `json:"stock,omitempty"`
}

// ── Request DTOs ──────────────────────────────────────
type CreateWatchlistRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Name   string `json:"name"    binding:"required,min=1,max=100"`
}

type AddStockRequest struct {
	StockID string `json:"stock_id" binding:"required"`
}

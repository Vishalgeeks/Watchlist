package models

import "time"

// ── User ──────────────────────────────────────────────
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ── Stock ─────────────────────────────────────────────
type Stock struct {
	ID           int       `json:"id"`
	Symbol       string    `json:"symbol"`
	CompanyName  string    `json:"company_name"`
	Exchange     string    `json:"exchange"`
	CurrentPrice float64   `json:"current_price"`
	LastUpdated  time.Time `json:"last_updated"`
}

// ── Watchlist ─────────────────────────────────────────
type Watchlist struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	StockCount int       `json:"stock_count"`
}

// ── Watchlist Item ────────────────────────────────────
type WatchlistItem struct {
	ID          int       `json:"id"`
	WatchlistID int       `json:"watchlist_id"`
	StockID     int       `json:"stock_id"`
	AddedAt     time.Time `json:"added_at"`
	Stock       *Stock    `json:"stock,omitempty"`
}

// ── Request DTOs ──────────────────────────────────────
type CreateWatchlistRequest struct {
	UserID int    `json:"user_id" binding:"required"`
	Name   string `json:"name"    binding:"required,min=1,max=100"`
}

type AddStockRequest struct {
	StockID int `json:"stock_id" binding:"required"`
}

// ── Standard Response ─────────────────────────────────
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

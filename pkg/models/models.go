package models

import "time"

// ── User ──────────────────────────────────────────────
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
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

// ── Auth Request DTOs ─────────────────────────────────
type RegisterRequest struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ── Watchlist Request DTOs ────────────────────────────
type CreateWatchlistRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
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

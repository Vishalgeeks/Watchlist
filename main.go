package main

import (
	"log"
	"watchlist-backend/config"
	"watchlist-backend/internal/db"
	"watchlist-backend/internal/watchlist"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database := db.Connect(cfg.DBConnectionString())
	defer database.Close()

	// ── Repositories ──────────────────────────
	watchlistRepo := watchlist.NewRepository(database)

	// ── Services ──────────────────────────────
	watchlistService := watchlist.NewService(watchlistRepo)

	// ── Handlers ──────────────────────────────
	watchlistHandler := watchlist.NewHandler(watchlistService)

	// ── Routes ────────────────────────────────
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/watchlists", watchlistHandler.Create)
		api.GET("/watchlists", watchlistHandler.GetAll)
		api.DELETE("/watchlists/:id", watchlistHandler.Delete)
		api.GET("/watchlists/:id/stocks", watchlistHandler.GetStocks)
		api.POST("/watchlists/:id/stocks", watchlistHandler.AddStock)
		api.DELETE("/watchlists/:id/stocks/:stockId", watchlistHandler.RemoveStock)
	}

	log.Printf("Server running on port %s", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}

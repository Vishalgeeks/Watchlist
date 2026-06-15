package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"watchlist-backend/config"
	"watchlist-backend/internal/db"
	"watchlist-backend/internal/watchlist"
	"watchlist-backend/internal/stock"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DBConnectionString())
	defer database.Close()

	// Repositories
	watchlistRepo := watchlist.NewRepository(database)
	stockRepo := stock.NewRepository(database)

	// Services
	watchlistService := watchlist.NewService(watchlistRepo)
	stockService := stock.NewService(stockRepo)

	// Handlers
	watchlistHandler := watchlist.NewHandler(watchlistService)
	stockHandler := stock.NewHandler(stockService)

	// Router
	r := gin.Default()

	api := r.Group("/api")
	{
		// ✅ STOCK ROUTES (THIS WAS MISSING)
		stockHandler.RegisterRoutes(api)

		// Watchlist routes
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
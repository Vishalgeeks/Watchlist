package main

import (
	"log"
	"watchlist-backend/config"
	"watchlist-backend/internal/auth"
	csvhandler "watchlist-backend/internal/csv"
	"watchlist-backend/internal/db"
	"watchlist-backend/internal/middleware"
	"watchlist-backend/internal/watchlist"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database := db.Connect(cfg.DBConnectionString())
	defer database.Close()

	// ── Auth ──────────────────────────────────
	authRepo := auth.NewRepository(database)
	authService := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	// ── Watchlist ─────────────────────────────
	watchlistRepo := watchlist.NewRepository(database)
	watchlistService := watchlist.NewService(watchlistRepo)
	watchlistHandler := watchlist.NewHandler(watchlistService)

	// ── CSV ───────────────────────────────────
	csvRepo := csvhandler.NewRepository(database)
	csvHandler := csvhandler.NewHandler(csvRepo, "data/stocks.csv")

	// ── Router ────────────────────────────────
	r := gin.Default()

	api := r.Group("/api")
	{
		// Public routes
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		// CSV import — public (instructor use karega)
		api.POST("/stocks/import", csvHandler.ImportCSV)
		api.GET("/stocks", csvHandler.GetAllStocks)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.POST("/watchlists", watchlistHandler.Create)
			protected.GET("/watchlists", watchlistHandler.GetAll)
			protected.DELETE("/watchlists/:id", watchlistHandler.Delete)
			protected.GET("/watchlists/:id/stocks", watchlistHandler.GetStocks)
			protected.POST("/watchlists/:id/stocks", watchlistHandler.AddStock)
			protected.DELETE("/watchlists/:id/stocks/:stockId", watchlistHandler.RemoveStock)
		}
	}

	log.Printf("Server running on port %s", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}

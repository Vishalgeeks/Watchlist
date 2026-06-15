package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"watchlist-backend/config"
	"watchlist-backend/internal/auth"
	"watchlist-backend/internal/db"
	"watchlist-backend/internal/middleware"
	"watchlist-backend/internal/stock"
	"watchlist-backend/internal/watchlist"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DBConnectionString())
	defer database.Close()

	// Auth
	authRepo := auth.NewRepository(database)
	authService := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	// Watchlist
	watchlistRepo := watchlist.NewRepository(database)
	watchlistService := watchlist.NewService(watchlistRepo)

	// Stock
	stockRepo := stock.NewRepository(database)
	stockService := stock.NewService(stockRepo)

	// Handlers
	watchlistHandler := watchlist.NewHandler(watchlistService)
	stockHandler := stock.NewHandler(stockService)

	// Router
	r := gin.Default()

	// CORS Middleware
	r.Use(middleware.CORSMiddleware())

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	api := r.Group("/api")

	// Public Auth Routes
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// Stock Routes
	stockHandler.RegisterRoutes(api)

	// Protected Routes
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

	log.Printf("Server running on port %s", cfg.ServerPort)

	r.Run(":" + cfg.ServerPort)
}

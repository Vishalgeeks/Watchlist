package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"watchlist-backend/config"
	"watchlist-backend/internal/auth"
	csvhandler "watchlist-backend/internal/csv"
	"watchlist-backend/internal/db"
	"watchlist-backend/internal/middleware"
	"watchlist-backend/internal/stock"
	"watchlist-backend/internal/watchlist"
	"watchlist-backend/internal/watchlist/search"
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
	watchlistHandler := watchlist.NewHandler(watchlistService)

	// Stock
	stockRepo := stock.NewRepository(database)
	stockService := stock.NewService(stockRepo)
	stockHandler := stock.NewHandler(stockService)

	// CSV
	csvRepo := csvhandler.NewRepository(database)
	csvHandler := csvhandler.NewHandler(csvRepo, cfg.CSVURL)

	// Search
	searchRepo := search.NewRepository(database)
	searchService := search.NewService(searchRepo)
	searchHandler := search.NewHandler(searchService)

	// Server start hote hi CSV load karo
	go func() {
		log.Println("Loading CSV data from URL...")
		stocks, err := csvhandler.ParseCSV(cfg.CSVURL)
		if err != nil {
			log.Printf("CSV load error: %v", err)
			return
		}
		inserted := 0
		for _, s := range stocks {
			if err := csvRepo.UpsertStock(&s); err != nil {
				continue
			}
			inserted++
		}
		log.Printf("CSV loaded: %d stocks inserted/updated", inserted)
	}()

	// Router
	r := gin.Default()

	// CORS Middleware
	r.Use(middleware.CORSMiddleware())

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	// ── api group pehle banao ──────────────────
	api := r.Group("/api")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		if err := database.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"db":     "disconnected",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"db":        "connected",
			"timestamp": time.Now(),
		})
	})

	// Search
	api.GET("/search/stocks", searchHandler.SearchStocks)

	// Public Auth Routes
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// Stock Routes
	stockHandler.RegisterRoutes(api)

	// CSV Import Route
	api.POST("/stocks/import", csvHandler.ImportCSV)

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

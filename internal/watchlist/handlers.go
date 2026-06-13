package watchlist

import (
	"net/http"
	"watchlist-backend/pkg/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// POST /api/watchlists
func (h *Handler) Create(c *gin.Context) {
	var req models.CreateWatchlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	w, err := h.service.CreateWatchlist(req.UserID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create watchlist"})
		return
	}

	c.JSON(http.StatusCreated, w)
}

// GET /api/watchlists?user_id=xxx
func (h *Handler) GetAll(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	watchlists, err := h.service.GetWatchlists(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch watchlists"})
		return
	}

	if watchlists == nil {
		watchlists = []models.Watchlist{}
	}

	c.JSON(http.StatusOK, gin.H{"watchlists": watchlists})
}

// DELETE /api/watchlists/:id?user_id=xxx
func (h *Handler) Delete(c *gin.Context) {
	watchlistID := c.Param("id")
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	if err := h.service.DeleteWatchlist(userID, watchlistID); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		if err.Error() == "watchlist not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "watchlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "watchlist deleted"})
}

// GET /api/watchlists/:id/stocks?user_id=xxx
func (h *Handler) GetStocks(c *gin.Context) {
	watchlistID := c.Param("id")
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	items, err := h.service.GetStocks(userID, watchlistID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		if err.Error() == "watchlist not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "watchlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if items == nil {
		items = []models.WatchlistItem{}
	}

	c.JSON(http.StatusOK, gin.H{"stocks": items})
}

// POST /api/watchlists/:id/stocks?user_id=xxx
func (h *Handler) AddStock(c *gin.Context) {
	watchlistID := c.Param("id")
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	var req models.AddStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddStock(userID, watchlistID, req.StockID); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		if err.Error() == "watchlist not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "watchlist not found"})
			return
		}
		if err.Error() == "stock already in watchlist" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "stock already in watchlist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "stock added to watchlist"})
}

// DELETE /api/watchlists/:id/stocks/:stockId?user_id=xxx
func (h *Handler) RemoveStock(c *gin.Context) {
	watchlistID := c.Param("id")
	stockID := c.Param("stockId")
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	if err := h.service.RemoveStock(userID, watchlistID, stockID); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		if err.Error() == "watchlist not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "watchlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock removed"})
}

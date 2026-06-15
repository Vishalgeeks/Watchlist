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
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	w, err := h.service.CreateWatchlist(req.UserID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to create watchlist",
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "watchlist created successfully",
		Data:    w,
	})
}

// GET /api/watchlists?user_id=1
func (h *Handler) GetAll(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "user_id required",
		})
		return
	}

	userID, err := ParseInt(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid user_id",
		})
		return
	}

	watchlists, err := h.service.GetWatchlists(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch watchlists",
		})
		return
	}

	if watchlists == nil {
		watchlists = []models.Watchlist{}
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "watchlists fetched successfully",
		Data:    watchlists,
	})
}

// DELETE /api/watchlists/:id?user_id=1
func (h *Handler) Delete(c *gin.Context) {
	watchlistIDStr := c.Param("id")
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "user_id required",
		})
		return
	}

	watchlistID, err := ParseInt(watchlistIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid watchlist id",
		})
		return
	}

	userID, err := ParseInt(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid user_id",
		})
		return
	}

	if err := h.service.DeleteWatchlist(userID, watchlistID); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, models.Response{
				Success: false,
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "watchlist not found" {
			c.JSON(http.StatusNotFound, models.Response{
				Success: false,
				Message: "watchlist not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "watchlist deleted successfully",
	})
}

// GET /api/watchlists/:id/stocks?user_id=1
func (h *Handler) GetStocks(c *gin.Context) {
	watchlistIDStr := c.Param("id")
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "user_id required",
		})
		return
	}

	watchlistID, err := ParseInt(watchlistIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid watchlist id",
		})
		return
	}

	userID, err := ParseInt(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid user_id",
		})
		return
	}

	items, err := h.service.GetStocks(userID, watchlistID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, models.Response{
				Success: false,
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "watchlist not found" {
			c.JSON(http.StatusNotFound, models.Response{
				Success: false,
				Message: "watchlist not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if items == nil {
		items = []models.WatchlistItem{}
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "stocks fetched successfully",
		Data:    items,
	})
}

// POST /api/watchlists/:id/stocks?user_id=1
func (h *Handler) AddStock(c *gin.Context) {
	watchlistIDStr := c.Param("id")
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "user_id required",
		})
		return
	}

	watchlistID, err := ParseInt(watchlistIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid watchlist id",
		})
		return
	}

	userID, err := ParseInt(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid user_id",
		})
		return
	}

	var req models.AddStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := h.service.AddStock(userID, watchlistID, req.StockID); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, models.Response{
				Success: false,
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "watchlist not found" {
			c.JSON(http.StatusNotFound, models.Response{
				Success: false,
				Message: "watchlist not found",
			})
			return
		}
		if err.Error() == "stock already in watchlist" {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "stock already in watchlist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "stock added to watchlist successfully",
	})
}

// DELETE /api/watchlists/:id/stocks/:stockId?user_id=1
func (h *Handler) RemoveStock(c *gin.Context) {
	watchlistIDStr := c.Param("id")
	stockIDStr := c.Param("stockId")
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "user_id required",
		})
		return
	}

	watchlistID, err := ParseInt(watchlistIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid watchlist id",
		})
		return
	}

	userID, err := ParseInt(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid user_id",
		})
		return
	}

	stockID, err := ParseInt(stockIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid stock id",
		})
		return
	}

	if err := h.service.RemoveStock(userID, watchlistID, stockID); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, models.Response{
				Success: false,
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "watchlist not found" {
			c.JSON(http.StatusNotFound, models.Response{
				Success: false,
				Message: "watchlist not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "stock removed successfully",
	})
}

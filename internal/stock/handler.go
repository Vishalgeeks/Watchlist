package stock

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"watchlist-backend/pkg/models"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	stocks := rg.Group("/stocks")
	{
		stocks.POST("", h.CreateStock)
		stocks.GET("", h.GetAllStocks)
		stocks.GET("/:id", h.GetStockByID)
		stocks.PUT("/:id", h.UpdateStock)
		stocks.DELETE("/:id", h.DeleteStock)
	}
}

func (h *Handler) CreateStock(c *gin.Context) {
	var stock models.Stock

	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.CreateStock(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, stock)
}

func (h *Handler) GetAllStocks(c *gin.Context) {
	stocks, err := h.service.GetAllStocks()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (h *Handler) GetStockByID(c *gin.Context) {
	id := c.Param("id")

	stock, err := h.service.GetStockByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "stock not found",
		})
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (h *Handler) UpdateStock(c *gin.Context) {
	id := c.Param("id")

	var stock models.Stock

	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.UpdateStock(id, &stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (h *Handler) DeleteStock(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteStock(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "stock deleted successfully",
	})
}
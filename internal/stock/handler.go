package stock

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"watchlist-backend/pkg/models"
)

type StockResponse struct {
	ID           int     `json:"id"`
	Symbol       string  `json:"symbol"`
	CompanyName  string  `json:"company_name"`
	Exchange     string  `json:"exchange"`
	CurrentPrice float64 `json:"current_price"`
	Sector       string  `json:"sector,omitempty"`
	LastUpdated  string  `json:"last_updated"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func toStockResponse(s models.Stock) StockResponse {
	return StockResponse{
		ID:           s.ID,
		Symbol:       s.Symbol,
		CompanyName:  s.CompanyName,
		Exchange:     s.Exchange,
		CurrentPrice: s.CurrentPrice,
		LastUpdated:  s.LastUpdated.Format("2006-01-02 15:04:05"),
	}
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
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := h.service.CreateStock(&stock); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "stock created successfully",
		Data:    toStockResponse(stock),
	})
}

func (h *Handler) GetAllStocks(c *gin.Context) {
	stocks, err := h.service.GetAllStocks()

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var response []StockResponse

	for _, s := range stocks {
		response = append(response, toStockResponse(s))
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "stocks fetched successfully",
		Data:    response,
	})
}

func (h *Handler) GetStockByID(c *gin.Context) {
	id := c.Param("id")

	stock, err := h.service.GetStockByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Message: "stock not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "stock fetched successfully",
		Data:    toStockResponse(*stock),
	})
}

func (h *Handler) UpdateStock(c *gin.Context) {
	id := c.Param("id")

	var stock models.Stock

	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := h.service.UpdateStock(id, &stock); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "stock updated successfully",
		Data:    toStockResponse(stock),
	})
}

func (h *Handler) DeleteStock(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteStock(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "stock deleted successfully",
	})
}

package search

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

func (h *Handler) SearchStocks(c *gin.Context) {
	query := c.Query("query")

	if query == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "query is required",
		})
		return
	}

	stocks, err := h.service.SearchStocks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "search results",
		Data:    stocks,
	})
}
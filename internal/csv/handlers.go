package csv

import (
	"net/http"
	"watchlist-backend/pkg/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo   *Repository
	csvURL string
}

func NewHandler(repo *Repository, csvURL string) *Handler {
	return &Handler{repo: repo, csvURL: csvURL}
}

// POST /api/stocks/import
func (h *Handler) ImportCSV(c *gin.Context) {
	stocks, err := ParseCSV(h.csvURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to parse CSV: " + err.Error(),
		})
		return
	}

	inserted := 0
	failed := 0
	for _, stock := range stocks {
		if err := h.repo.UpsertStock(&stock); err != nil {
			failed++
			continue
		}
		inserted++
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "CSV imported successfully",
		Data: gin.H{
			"total_processed":  len(stocks),
			"inserted_updated": inserted,
			"failed":           failed,
		},
	})
}

package stock

import (
	"errors"

	"watchlist-backend/pkg/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateStock(stock *models.Stock) error {
	if stock.Symbol == "" {
		return errors.New("symbol is required")
	}

	if stock.CompanyName == "" {
		return errors.New("company name is required")
	}

	return s.repo.Create(stock)
}

func (s *Service) GetAllStocks() ([]models.Stock, error) {
	return s.repo.GetAll()
}

func (s *Service) GetStockByID(id string) (*models.Stock, error) {
	return s.repo.GetByID(id)
}

func (s *Service) UpdateStock(id string, stock *models.Stock) error {
	if stock.Symbol == "" {
		return errors.New("symbol is required")
	}

	if stock.CompanyName == "" {
		return errors.New("company name is required")
	}

	return s.repo.Update(id, stock)
}

func (s *Service) DeleteStock(id string) error {
	return s.repo.Delete(id)
}
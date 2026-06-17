package search

import "watchlist-backend/pkg/models"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SearchStocks(query string) ([]models.Stock, error) {
	return s.repo.SearchStocks(query)
}
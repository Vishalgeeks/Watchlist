package watchlist

import (
	"errors"
	"watchlist-backend/pkg/models"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateWatchlist(userID, name string) (*models.Watchlist, error) {
	w := &models.Watchlist{
		ID:     uuid.New().String(),
		UserID: userID,
		Name:   name,
	}
	if err := s.repo.Create(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) GetWatchlists(userID string) ([]models.Watchlist, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *Service) DeleteWatchlist(userID, watchlistID string) error {
	w, err := s.repo.GetByID(watchlistID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(watchlistID)
}

func (s *Service) GetStocks(userID, watchlistID string) ([]models.WatchlistItem, error) {
	w, err := s.repo.GetByID(watchlistID)
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return s.repo.GetStocks(watchlistID)
}

func (s *Service) AddStock(userID, watchlistID, stockID string) error {
	w, err := s.repo.GetByID(watchlistID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return errors.New("unauthorized")
	}
	exists, err := s.repo.StockExists(watchlistID, stockID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("stock already in watchlist")
	}
	return s.repo.AddStock(watchlistID, stockID, uuid.New().String())
}

func (s *Service) RemoveStock(userID, watchlistID, stockID string) error {
	w, err := s.repo.GetByID(watchlistID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.RemoveStock(watchlistID, stockID)
}

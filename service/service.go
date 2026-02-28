package service

import (
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/repository"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
)

type Service struct {
	btRepository *repository.BtRepository
	bqRepository *repository.BqRepository
}

func NewService(btRepository *repository.BtRepository, bqRepository *repository.BqRepository) *Service {
	return &Service{
		btRepository: btRepository,
		bqRepository: bqRepository,
	}
}

func (s *Service) CreateEvent(event *types.CreateEventRequest) error {
	return nil
}

func (s *Service) GetTopProductsFromStore(storeID string, windowHours int) (*types.GetTopProductsFromStoreResponse, error) {
	return nil, nil
}

func (s *Service) GetEventsFromUser(userID string, limit int) ([]*types.Event, error) {
	return nil, nil
}

func (s *Service) Ping() error {
	return nil
}

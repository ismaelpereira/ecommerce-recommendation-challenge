package service

import (
	"context"
	"fmt"

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

func (s *Service) CreateEvent(ctx context.Context, event types.CreateEventRequest) error {
	err := s.bqRepository.CreateEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("Error creating event on Big Query: %w", err)
	}
	err = s.btRepository.CreateEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("Error creating event on Big Query: %w", err)
	}
	return nil
}

func (s *Service) GetTopProductsFromStore(ctx context.Context, storeID string, windowHours int) (*types.GetTopProductsFromStoreResponse, error) {
	topProducts, err := s.bqRepository.GetTopProductsFromStore(ctx, storeID, windowHours)
	if err != nil {
		return nil, fmt.Errorf("Error Getting Top Products from store ID %s: %w", storeID, err)
	}

	return &types.GetTopProductsFromStoreResponse{
		StoreID:     storeID,
		WindowHours: windowHours,
		Products:    topProducts,
	}, nil
}

func (s *Service) GetEventsFromUser(ctx context.Context, userID string, limit int) ([]types.Event, error) {
	events, err := s.btRepository.GetEventsFromUser(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("Error Getting events from user %s: %w", userID, err)
	}

	return events, nil
}

func (s *Service) Ping(ctx context.Context) error {
	bqErr := s.bqRepository.Ping(ctx)
	btErr := s.btRepository.Ping(ctx)
	if bqErr != nil && btErr == nil {
		return fmt.Errorf("Error connecting on Big Query\nBigtable Connected!")
	}
	if bqErr == nil && btErr != nil {
		return fmt.Errorf("Error connecting on Big Table\n Bigquery Connected!")
	}
	if bqErr != nil && btErr != nil {
		return fmt.Errorf("Error connecting on Big Query and Big Table")
	}
	return nil
}

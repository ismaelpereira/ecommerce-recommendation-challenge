package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
)

type Service struct {
	btRepository BigTableRepository
	bqRepository BigQueryRepository
}

func NewService(btRepository BigTableRepository, bqRepository BigQueryRepository) *Service {
	return &Service{
		btRepository: btRepository,
		bqRepository: bqRepository,
	}
}

func (s *Service) CreateEvent(ctx context.Context, req types.CreateEventRequest) (*types.Event, error) {
	event := &types.Event{
		ID:        uuid.NewString(),
		UserID:    req.UserID,
		ProductID: req.ProductID,
		StoreID:   req.StoreID,
		EventType: req.EventType,
		Timestamp: req.Timestamp,
	}
	err := s.bqRepository.CreateEvent(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("Error creating event on Big Query: %w", err)
	}
	err = s.btRepository.CreateEvent(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("Error creating event on Big Table: %w", err)
	}
	return event, nil
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

func (s *Service) Ping(ctx context.Context) (*types.PingErrorResponse, error) {
	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	bqErr := s.bqRepository.Ping(healthCtx)
	btErr := s.btRepository.Ping(healthCtx)

	if bqErr != nil && btErr == nil {
		return &types.PingErrorResponse{
			Message:       "Big Query Connection Error",
			BigQueryError: bqErr.Error(),
		}, bqErr
	}
	if bqErr == nil && btErr != nil {
		return &types.PingErrorResponse{
			Message:       "Big Table Connection Error",
			BigTableError: btErr.Error(),
		}, btErr
	}
	if bqErr != nil && btErr != nil {
		return &types.PingErrorResponse{
			Message:       "Both Big Table and Bigtable services are down",
			BigTableError: btErr.Error(),
			BigQueryError: bqErr.Error(),
		}, bqErr
	}
	return &types.PingErrorResponse{
		Message: "All services connected successfully",
	}, nil
}

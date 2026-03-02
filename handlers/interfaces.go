package handlers

import (
	"context"

	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
)

type Service interface {
	CreateEvent(ctx context.Context, req types.CreateEventRequest) (*types.Event, error)
	GetTopProductsFromStore(ctx context.Context, storeID string, hours int) (*types.GetTopProductsFromStoreResponse, error)
	GetEventsFromUser(ctx context.Context, userID string, limit int) ([]types.Event, error)
	Ping(ctx context.Context) (*types.PingErrorResponse, error)
}

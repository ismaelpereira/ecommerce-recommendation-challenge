package service

import (
	"context"

	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
)

type BigTableRepository interface {
	CreateEvent(ctx context.Context, event types.CreateEventRequest) error
	GetEventsFromUser(ctx context.Context, userID string, limit int) ([]types.Event, error)
	Ping(ctx context.Context) error
}

type BigQueryRepository interface {
	CreateEvent(ctx context.Context, event types.CreateEventRequest) error
	GetTopProductsFromStore(ctx context.Context, storeID string, windowHours int) ([]types.Product, error)
	Ping(ctx context.Context) error
}

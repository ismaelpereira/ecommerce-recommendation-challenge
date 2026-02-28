package repository

import (
	cloudbq "cloud.google.com/go/bigquery"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
)

type BqRepository struct {
	client *cloudbq.Client
}

func NewBqRepository(client *cloudbq.Client) *BqRepository {
	return &BqRepository{
		client: client,
	}
}

func (r *BqRepository) CreateEvent(event *types.CreateEventRequest) error {
	return nil
}

func (r *BqRepository) GetTopProductsFromStore(storeID string, windowHours int) ([]*types.Product, error) {
	return nil, nil
}

func (r *BqRepository) Ping() error {
	return nil
}

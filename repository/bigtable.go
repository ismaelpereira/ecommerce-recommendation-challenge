package repository

import (
	cloudbt "cloud.google.com/go/bigtable"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
)

type BtRepository struct {
	client *cloudbt.Client
}

func NewBtRepository(client *cloudbt.Client) *BtRepository {
	return &BtRepository{
		client: client,
	}
}

func (r *BtRepository) CreateEvent(event *types.CreateEventRequest) error {
	return nil
}

func (r *BtRepository) GetEventsFromUser(userID string, limit int) ([]*types.Event, error) {
	return nil, nil
}

func (r *BtRepository) Ping() error {
	return nil
}

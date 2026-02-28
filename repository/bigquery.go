package repository

import (
	"context"
	"fmt"

	cloudbq "cloud.google.com/go/bigquery"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
	"google.golang.org/api/iterator"
)

type BqRepository struct {
	client *cloudbq.Client
}

func NewBqRepository(client *cloudbq.Client) *BqRepository {
	return &BqRepository{
		client: client,
	}
}

func (r *BqRepository) CreateEvent(ctx context.Context, event types.CreateEventRequest) error {
	query := r.client.Query(`
		INSERT INTO events(
				user_id,
				product_id,
				store_id,
				event_type,
				timestamp
		)
		VALUES(
			@user_id,
			@product_id,
			@store_id,
			@event_type,
			@timestamp
		)
	`)
	query.Parameters = []cloudbq.QueryParameter{
		{Name: "user_id", Value: event.UserID},
		{Name: "product_id", Value: event.ProductID},
		{Name: "store_id", Value: event.StoreID},
		{Name: "event_type", Value: event.EventType},
		{Name: "timestamp", Value: event.Timestamp},
	}

	job, err := query.Run(ctx)
	if err != nil {
		return fmt.Errorf("Error on Insert RUN into Biquery: %w", err)
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Error on Insert Wait into Biquery: %w", err)
	}

	if err := status.Err(); err != nil {
		return fmt.Errorf("Error inserting event into Biquery: %w", err)
	}

	return nil
}

func (r *BqRepository) GetTopProductsFromStore(ctx context.Context, storeID string, windowHours int) ([]types.Product, error) {
	dataset := r.client.Dataset("ecommerce_events")
	table := dataset.Table("events")

	queryStr := fmt.Sprintf(`
		SELECT product_id, COUNT(event_id) AS view_count
		FROM %s.%s
		WHERE store_id = @store_id
			AND event_type = "view"
			AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(),INTERVAL @hours HOUR)
		GROUP BY product_id
		GROUP BY view_count DESC
		LIMIT 10,
	`, dataset.DatasetID, table.TableID)

	query := r.client.Query(queryStr)

	query.Parameters = []cloudbq.QueryParameter{
		{Name: "store_id", Value: storeID},
		{Name: "hours", Value: windowHours},
	}

	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error reading top events: %w", err)
	}

	var topProducts []types.Product

	for {
		var row types.Product
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error iterating bigquery: %w", err)
		}

		topProducts = append(topProducts, row)
	}

	return topProducts, nil
}

func (r *BqRepository) Ping(ctx context.Context) error {
	query := r.client.Query("SELECT 1")

	it, err := query.Read(ctx)
	if err != nil {
		return fmt.Errorf("bigquery ping read: %w", err)
	}

	var row []cloudbq.Value
	err = it.Next(&row)
	if err != nil && err != iterator.Done {
		return fmt.Errorf("Error pinging Big Query: %w", err)
	}

	return nil
}

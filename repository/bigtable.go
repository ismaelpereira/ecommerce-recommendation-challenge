package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	cloudbt "cloud.google.com/go/bigtable"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
)

type BtRepository struct {
	client     *cloudbt.Client
	tableName  string
	familyName string
}

func NewBtRepository(client *cloudbt.Client, tableName string, familyName string) *BtRepository {
	return &BtRepository{
		client:     client,
		tableName:  tableName,
		familyName: familyName,
	}
}

func (r *BtRepository) CreateEvent(ctx context.Context, event types.CreateEventRequest) error {
	table := r.client.Open(r.tableName)

	rowKey := fmt.Sprintf("user#%s#revts#%d", event.UserID, math.MaxInt64-event.Timestamp.UnixNano())

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("Error marshalling event: %w", err)
	}

	mut := cloudbt.NewMutation()
	mut.Set("events", "data", cloudbt.Now(), data)

	return table.Apply(ctx, rowKey, mut)
}

func (r *BtRepository) GetEventsFromUser(ctx context.Context, userID string, limit int) ([]types.Event, error) {
	table := r.client.Open(r.tableName)

	prefix := fmt.Sprintf("user#%s#", userID)

	var events []types.Event

	err := table.ReadRows(ctx, cloudbt.PrefixRange(prefix),
		func(row cloudbt.Row) bool {

			item := row["events"][0]

			var e types.Event
			if err := json.Unmarshal(item.Value, &e); err == nil {
				events = append(events, e)
			}
			return len(events) < limit
		})

	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *BtRepository) Ping(ctx context.Context) error {
	table := r.client.Open(r.tableName)

	_, err := table.ReadRow(ctx, "healthcheck-row")

	if err != nil {
		return fmt.Errorf("bigtable ping failed: %w", err)
	}

	return nil
}

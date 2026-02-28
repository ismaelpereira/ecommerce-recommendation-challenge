package bigquery

import (
	"context"

	cloudbq "cloud.google.com/go/bigquery"
)

func NewClient(ctx context.Context, projectID string) (*cloudbq.Client, error) {
	client, err := cloudbq.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

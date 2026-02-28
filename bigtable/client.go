package bigtable

import (
	"context"

	cloudbt "cloud.google.com/go/bigtable"
)

func NewClient(ctx context.Context, projectID, instanceID string) (*cloudbt.Client, error) {
	client, err := cloudbt.NewClient(ctx, projectID, instanceID)
	if err != nil {
		return nil, err
	}

	return client, nil
}

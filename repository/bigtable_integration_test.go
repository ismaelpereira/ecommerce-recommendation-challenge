//go:build integration
// +build integration

package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	cloudbt "cloud.google.com/go/bigtable"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/repository"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
	"github.com/stretchr/testify/assert"
)

func TestBigtableIntegration_CreateAndRead(t *testing.T) {
	// Skip if emulator not running
	if os.Getenv("BIGTABLE_EMULATOR_HOST") == "" {
		t.Skip("BIGTABLE_EMULATOR_HOST not set")
	}

	ctx := context.Background()

	projectID := "test-project"
	instanceID := "test-instance"
	tableName := "events"
	familyName := "events"

	// Admin client (to create table)
	adminClient, err := cloudbt.NewAdminClient(ctx, projectID, instanceID)
	assert.NoError(t, err)
	defer adminClient.Close()

	// Create table if not exists
	_ = adminClient.DeleteTable(ctx, tableName) // clean state
	err = adminClient.CreateTable(ctx, tableName)
	assert.NoError(t, err)

	err = adminClient.CreateColumnFamily(ctx, tableName, familyName)
	assert.NoError(t, err)

	// Data client
	client, err := cloudbt.NewClient(ctx, projectID, instanceID)
	assert.NoError(t, err)
	defer client.Close()

	repo := repository.NewBtRepository(client, tableName, familyName)

	now := time.Now()

	event := &types.Event{
		ID: "event_1",
		UserID:    "user_1",
		ProductID: "product_1",
		StoreID:   "store_1",
		EventType: "view",
		Timestamp: now,
	}

	// Write
	err = repo.CreateEvent(ctx, event)
	assert.NoError(t, err)

	// Read
	events, err := repo.GetEventsFromUser(ctx, "user_1", 10)
	assert.NoError(t, err)

	assert.Len(t, events, 1)
	assert.Equal(t, "user_1", events[0].UserID)
	assert.Equal(t, "product_1", events[0].ProductID)
}

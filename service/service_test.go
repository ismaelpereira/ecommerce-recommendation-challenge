package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ismaelpereira/ecommerce-recommendation-challenge/service"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
	"github.com/stretchr/testify/assert"
)

type mockBT struct {
	createErr error
	events    []types.Event
	pingErr   error
}

func (m *mockBT) CreateEvent(ctx context.Context, event types.CreateEventRequest) error {
	return m.createErr
}

func (m *mockBT) GetEventsFromUser(ctx context.Context, userID string, limit int) ([]types.Event, error) {
	if m.events == nil {
		return nil, errors.New("bt error")
	}
	return m.events, nil
}

func (m *mockBT) Ping(ctx context.Context) error {
	return m.pingErr
}

type mockBQ struct {
	createErr error
	products  []types.Product
	pingErr   error
}

func (m *mockBQ) CreateEvent(ctx context.Context, event types.CreateEventRequest) error {
	return m.createErr
}

func (m *mockBQ) GetTopProductsFromStore(ctx context.Context, storeID string, windowHours int) ([]types.Product, error) {
	if m.products == nil {
		return nil, errors.New("bq error")
	}
	return m.products, nil
}

func (m *mockBQ) Ping(ctx context.Context) error {
	return m.pingErr
}

func TestCreateEvent_Success(t *testing.T) {
	bt := &mockBT{}
	bq := &mockBQ{}

	svc := service.NewService(bt, bq)

	err := svc.CreateEvent(context.Background(), types.CreateEventRequest{})

	assert.NoError(t, err)
}

func TestCreateEvent_BigQueryFails(t *testing.T) {
	bt := &mockBT{}
	bq := &mockBQ{createErr: errors.New("bq fail")}

	svc := service.NewService(bt, bq)

	err := svc.CreateEvent(context.Background(), types.CreateEventRequest{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Big Query")
}

func TestCreateEvent_BigTableFails(t *testing.T) {
	bt := &mockBT{createErr: errors.New("bt fail")}
	bq := &mockBQ{}

	svc := service.NewService(bt, bq)

	err := svc.CreateEvent(context.Background(), types.CreateEventRequest{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Big Table")
}

func TestGetTopProductsFromStore_Success(t *testing.T) {
	expectedProducts := []types.Product{
		{
			ProductID: "product_1",
			ViewCount: 10,
		},
	}

	bt := &mockBT{}
	bq := &mockBQ{products: expectedProducts}

	svc := service.NewService(bt, bq)

	resp, err := svc.GetTopProductsFromStore(context.Background(), "store_1", 24)

	assert.NoError(t, err)
	assert.Equal(t, "store_1", resp.StoreID)
	assert.Equal(t, 24, resp.WindowHours)
	assert.Len(t, resp.Products, 1)
	assert.Equal(t, "product_1", resp.Products[0].ProductID)
}

func TestGetTopProductsFromStore_Fails(t *testing.T) {
	bt := &mockBT{}
	bq := &mockBQ{products: nil}

	svc := service.NewService(bt, bq)

	resp, err := svc.GetTopProductsFromStore(context.Background(), "store_1", 24)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetEventsFromUser_Success(t *testing.T) {
	expectedEvents := []types.Event{
		{UserID: "user_1"},
	}

	bt := &mockBT{events: expectedEvents}
	bq := &mockBQ{}

	svc := service.NewService(bt, bq)

	events, err := svc.GetEventsFromUser(context.Background(), "user_1", 10)

	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "user_1", events[0].UserID)
}

func TestGetEventsFromUser_Fails(t *testing.T) {
	bt := &mockBT{events: nil}
	bq := &mockBQ{}

	svc := service.NewService(bt, bq)

	events, err := svc.GetEventsFromUser(context.Background(), "user_1", 10)

	assert.Error(t, err)
	assert.Nil(t, events)
}

func TestPing_Success(t *testing.T) {
	bt := &mockBT{}
	bq := &mockBQ{}

	svc := service.NewService(bt, bq)

	err := svc.Ping(context.Background())

	assert.NoError(t, err)
}

func TestPing_BigQueryFails(t *testing.T) {
	bt := &mockBT{}
	bq := &mockBQ{pingErr: errors.New("bq down")}

	svc := service.NewService(bt, bq)

	err := svc.Ping(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Big Query")
}

func TestPing_BigTableFails(t *testing.T) {
	bt := &mockBT{pingErr: errors.New("bt down")}
	bq := &mockBQ{}

	svc := service.NewService(bt, bq)

	err := svc.Ping(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Big Table")
}
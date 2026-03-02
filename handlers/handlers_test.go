package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/handlers"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	createEventResp *types.Event
	createEventErr  error

	topProductsResp *types.GetTopProductsFromStoreResponse
	topProductsErr  error

	eventsResp []types.Event
	eventsErr  error

	pingResp *types.PingErrorResponse
	pingErr  error
}

func (m *mockService) CreateEvent(ctx context.Context, req types.CreateEventRequest) (*types.Event, error) {
	return m.createEventResp, m.createEventErr
}

func (m *mockService) GetTopProductsFromStore(ctx context.Context, storeID string, hours int) (*types.GetTopProductsFromStoreResponse, error) {
	return m.topProductsResp, m.topProductsErr
}

func (m *mockService) GetEventsFromUser(ctx context.Context, userID string, limit int) ([]types.Event, error) {
	return m.eventsResp, m.eventsErr
}

func (m *mockService) Ping(ctx context.Context) (*types.PingErrorResponse, error) {
	return m.pingResp, m.pingErr
}

func TestCreateEvent_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockService{
		createEventResp: &types.Event{ID: "123"},
	}

	h := handlers.NewHandler(mockSvc)

	router := gin.New()
	router.POST("/events", h.CreateEvent)

	body := `{
		"user_id": "u1",
		"product_id": "p1",
		"store_id": "s1",
		"event_type": "view",
		"timestamp": "2024-01-01T00:00:00Z"
	}`

	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "123")
}

func TestCreateEvent_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockService{}
	h := handlers.NewHandler(mockSvc)

	router := gin.New()
	router.POST("/events", h.CreateEvent)

	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTopProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockService{
		topProductsResp: &types.GetTopProductsFromStoreResponse{
			StoreID: "store1",
		},
	}

	h := handlers.NewHandler(mockSvc)

	router := gin.New()
	router.GET("/top", h.GetTopProductsFromStore)

	req := httptest.NewRequest(http.MethodGet, "/top?store_id=store1&hours=24", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetEventsFromUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockService{
		eventsResp: []types.Event{{ID: "1"}},
	}

	h := handlers.NewHandler(mockSvc)

	router := gin.New()
	router.GET("/users/:user_id/events", h.GetEventsFromUser)

	req := httptest.NewRequest(http.MethodGet, "/users/u1/events?limit=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHealthCheck_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockService{
		pingResp: &types.PingErrorResponse{
			Message: "All services connected successfully",
		},
		pingErr: nil,
	}

	h := handlers.NewHandler(mockSvc)

	router := gin.New()
	router.GET("/health", h.HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "All services connected successfully")
}

func TestHealthCheck_ServiceUnavailable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockService{
		pingResp: &types.PingErrorResponse{
			Message:       "Big Query Connection Error",
			BigQueryError: "context deadline exceeded",
		},
		pingErr: errors.New("context deadline exceeded"),
	}

	h := handlers.NewHandler(mockSvc)

	router := gin.New()
	router.GET("/health", h.HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), "Big Query Connection Error")
}

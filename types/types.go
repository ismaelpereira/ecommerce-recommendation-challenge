package types

import "time"

type CreateEventRequest struct {
	UserID    string    `json:"user_id" binding:"required"`
	ProductID string    `json:"product_id" binding:"required"`
	StoreID   string    `json:"store_id" binding:"required"`
	EventType string    `json:"event_type" binding:"required,oneof='view' 'add_to_cart' 'purchase'"`
	Timestamp time.Time `json:"timestamp"`
}

type GetTopProductsFromStoreResponse struct {
	StoreID     string    `json:"store_id"`
	WindowHours int       `json:"window_hours"`
	Products    []Product `json:"products"`
}

type Product struct {
	ProductID string `json:"product_id" bigquery:"product_id"`
	ViewCount int    `json:"view_count" bigquery:"view_count"`
}

type Event struct {
	ID        string
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	StoreID   string    `json:"store_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

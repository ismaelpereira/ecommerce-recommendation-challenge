package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/handlers"
)

func main() {
	r := gin.Default()

	r.POST("/events", handlers.CreateEvent)
	r.GET("/analytics/top-products", handlers.GetTopProductsFromStore)
	r.GET("/events/user/:user_id", handlers.GetEventsFromUser)
	r.GET("/health", handlers.HealthCheck)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

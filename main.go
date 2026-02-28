package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/bigquery"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/bigtable"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/config"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/handlers"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/repository"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/service"
)

func main() {
	config := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	btClient, err := bigtable.NewClient(ctx, config.ProjectID, config.BigtableInstance)
	if err != nil {
		log.Fatalf("Unable to start Big Table Client")
	}

	bqClient, err := bigquery.NewClient(ctx, config.ProjectID)
	if err != nil {
		log.Fatalf("Unable to start Big Query Client")
	}

	btRepository := repository.NewBtRepository(btClient)

	bqRepository := repository.NewBqRepository(bqClient)

	service := service.NewService(btRepository, bqRepository)

	handlers := handlers.NewHandler(service)

	r := gin.Default()

	r.POST("/events", handlers.CreateEvent)
	r.GET("/analytics/top-products", handlers.GetTopProductsFromStore)
	r.GET("/events/user/:user_id", handlers.GetEventsFromUser)
	r.GET("/health", handlers.HealthCheck)

	go func() {
		if err := r.Run(); err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	}()

	log.Println("server started on :8080")

	<-ctx.Done()
	log.Println("shutdown signal received")

	btClient.Close()
	bqClient.Close()

}

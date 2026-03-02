package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	var bigtableIntance string

	if config.Env == "dev" {
		bigtableIntance = "local-instance"
	} else {
		bigtableIntance = config.BigtableInstance
	}

	btClient, err := bigtable.NewClient(ctx, config.ProjectID, bigtableIntance)
	if err != nil {
		log.Fatalf("Unable to start Big Table Client")
	}

	bqClient, err := bigquery.NewClient(ctx, config.ProjectID)
	if err != nil {
		log.Fatalf("Unable to start Big Query Client")
	}

	btRepository := repository.NewBtRepository(btClient, config.BigTableTable, config.BigTableFamily)

	bqRepository := repository.NewBqRepository(bqClient, config.BigQueryDataset, config.BigQueryTable)

	service := service.NewService(btRepository, bqRepository)

	handlers := handlers.NewHandler(service)

	r := gin.Default()

	r.POST("/events", handlers.CreateEvent)
	r.GET("/analytics/top-products", handlers.GetTopProductsFromStore)
	r.GET("/events/user/:user_id", handlers.GetEventsFromUser)
	r.GET("/health", handlers.HealthCheck)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run server: %v", err)
		}
	}()

	log.Println("server started on :8080")

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v\n", err)
	}

	if err := btClient.Close(); err != nil {
		log.Printf("Big Table close failed: %v\n", err)
	}
	if err := bqClient.Close(); err != nil {
		log.Printf("Big Query close failed: %v\n", err)
	}

	log.Println("server exited properly")

}

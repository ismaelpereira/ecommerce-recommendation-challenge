package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/service"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/types"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) CreateEvent(c *gin.Context) {
	ctx := c.Request.Context()

	var event types.CreateEventRequest
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	err := h.service.CreateEvent(ctx, event)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Event saved successfuly", "data": event})
}

func (h *Handler) GetTopProductsFromStore(c *gin.Context) {
	ctx := c.Request.Context()

	storeID := c.Query("store_id")
	if storeID == "" {
		c.JSON(400, gin.H{"error": "store_id query parameter missing"})
		return
	}
	windowHours := c.Query("hours")
	if windowHours == "" {
		c.JSON(400, gin.H{"error": "hours query parameter missing"})
		return
	}

	intWindowHours, err := strconv.Atoi(windowHours)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	topProducts, err := h.service.GetTopProductsFromStore(ctx, storeID, intWindowHours)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, topProducts)
}

func (h *Handler) GetEventsFromUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("user_id")
	log.Println(userID)
	if userID == ":user_id" {
		c.JSON(400, gin.H{"error": "mising user_id path parameter"})
		return
	}

	treshold := c.Query("limit")
	if treshold == "" {
		c.JSON(400, gin.H{"error": "limit query parameter missing"})
		return
	}

	intTreshold, err := strconv.Atoi(treshold)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	events, err := h.service.GetEventsFromUser(ctx, userID, intTreshold)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, events)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	err := h.service.Ping(ctx)
	if err != nil {
		c.JSON(503, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message": "All services connected successfuly",
	})
}

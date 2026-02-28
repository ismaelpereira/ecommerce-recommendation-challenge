package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ismaelpereira/ecommerce-recommendation-challenge/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) CreateEvent(c *gin.Context) {}

func (h *Handler) GetTopProductsFromStore(c *gin.Context) {}

func (h *Handler) GetEventsFromUser(c *gin.Context) {}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

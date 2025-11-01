package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxviazov/dolina-flower-order-backend/internal/services"
)

type FlowerHandler struct {
	orderService *services.OrderService
}

func NewFlowerHandler(orderService *services.OrderService) *FlowerHandler {
	return &FlowerHandler{
		orderService: orderService,
	}
}

func (h *FlowerHandler) GetAvailableFlowers(c *gin.Context) {
	flowers, err := h.orderService.GetAvailableFlowers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"flowers": flowers,
	})
}

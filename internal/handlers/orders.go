package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxviazov/dolina-flower-order-backend/internal/dto"
	"github.com/maxviazov/dolina-flower-order-backend/internal/services"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create order: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

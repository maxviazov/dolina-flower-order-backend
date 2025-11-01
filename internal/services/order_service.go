package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/maxviazov/dolina-flower-order-backend/internal/domain"
	"github.com/maxviazov/dolina-flower-order-backend/internal/repository/sqlite"
)

type OrderService struct {
	repo *sqlite.Repository
}

func NewOrderService(repo *sqlite.Repository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) GetAvailableFlowers(ctx context.Context) ([]domain.Item, error) {
	return s.repo.GetAvailableFlowers(ctx)
}

func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*domain.Order, error) {
	order := &domain.Order{
		ID:         uuid.New().String(),
		MarkBox:    req.MarkBox,
		CustomerID: req.CustomerID,
		Status:     domain.OrderStatusPending,
		CreatedAt:  time.Now(),
		Notes:      req.Notes,
	}

	for _, itemReq := range req.Items {
		item := domain.Item{
			ID:         uuid.New().String(),
			OrderID:    order.ID,
			Variety:    itemReq.Variety,
			Length:     itemReq.Length,
			BoxCount:   itemReq.BoxCount,
			PackRate:   itemReq.PackRate,
			TotalStems: itemReq.TotalStems,
			FarmName:   itemReq.FarmName,
			TruckName:  itemReq.TruckName,
			Comments:   itemReq.Comments,
			Price:      itemReq.Price,
		}
		order.Items = append(order.Items, item)
	}

	order.CalculateTotal()

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.GetByID(ctx, id)
}

type CreateOrderRequest struct {
	MarkBox    string                   `json:"mark_box" binding:"required,min=1,max=10"`
	CustomerID string                   `json:"customer_id" binding:"required,min=1"`
	Items      []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
	Notes      string                   `json:"notes,omitempty"`
}

type CreateOrderItemRequest struct {
	Variety    string  `json:"variety" binding:"required,min=1,max=100"`
	Length     int     `json:"length" binding:"required,min=1,max=200"`
	BoxCount   float64 `json:"box_count" binding:"required,gt=0"`
	PackRate   int     `json:"pack_rate" binding:"required,min=1"`
	TotalStems int     `json:"total_stems" binding:"required,min=1"`
	FarmName   string  `json:"farm_name" binding:"required,min=1,max=100"`
	TruckName  string  `json:"truck_name" binding:"required,min=1,max=100"`
	Comments   string  `json:"comments,omitempty"`
	Price      float64 `json:"price,omitempty" binding:"gte=0"`
}

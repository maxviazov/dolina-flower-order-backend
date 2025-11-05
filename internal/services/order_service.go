package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/maxviazov/dolina-flower-order-backend/internal/domain"
	"github.com/maxviazov/dolina-flower-order-backend/internal/dto"
)

type OrderService struct {
	repo domain.OrderRepository
}

func NewOrderService(repo domain.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) GetAvailableFlowers(ctx context.Context) ([]domain.Item, error) {
	return s.repo.GetAvailableFlowers(ctx)
}

func (s *OrderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*domain.Order, error) {
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

	order.TotalAmount = order.CalculateTotal()

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.GetByID(ctx, id)
}

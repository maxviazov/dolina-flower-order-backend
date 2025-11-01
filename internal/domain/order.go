package domain

import (
	"time"
)

// OrderStatus представляет статус заказа
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusFarmOrder  OrderStatus = "farm_order"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// Order представляет заказ клиента
type Order struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id"`
	Items       []Item      `json:"items"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	ProcessedAt *time.Time  `json:"processed_at,omitempty"`
	FarmOrderID *string     `json:"farm_order_id,omitempty"`
	Notes       string      `json:"notes,omitempty"`
	TotalAmount float64     `json:"total_amount"`
}

// Item представляет позицию в заказе
type Item struct {
	ID         string  `json:"id"`
	FlowerType string  `json:"flower_type"`
	Variety    string  `json:"variety,omitempty"`
	Color      string  `json:"color,omitempty"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	Notes      string  `json:"notes,omitempty"`
}

// Customer представляет клиента
type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Company   string    `json:"company,omitempty"`
	Address   Address   `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Address представляет адрес
type Address struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
}

// FarmOrderStatus представляет статус заказа на ферме
type FarmOrderStatus string

const (
	FarmOrderStatusSent      FarmOrderStatus = "sent"
	FarmOrderStatusConfirmed FarmOrderStatus = "confirmed"
	FarmOrderStatusDelivered FarmOrderStatus = "delivered"
	FarmOrderStatusCancelled FarmOrderStatus = "cancelled"
)

// FarmOrder представляет заказ для фермы
type FarmOrder struct {
	ID        string          `json:"id"`
	OrderID   string          `json:"order_id"`
	Items     []Item          `json:"items"`
	Status    FarmOrderStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Notes     string          `json:"notes,omitempty"`
}

// CalculateTotal вычисляет общую сумму заказа
func (o *Order) CalculateTotal() float64 {
	total := 0.0
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	o.TotalAmount = total
	return total
}

// IsValidStatus проверяет валидность статуса заказа
func (o *Order) IsValidStatus(status OrderStatus) bool {
	validStatuses := []OrderStatus{
		OrderStatusPending,
		OrderStatusProcessing,
		OrderStatusFarmOrder,
		OrderStatusCompleted,
		OrderStatusCancelled,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// CanTransitionTo проверяет возможность перехода к новому статусу
func (o *Order) CanTransitionTo(newStatus OrderStatus) bool {
	transitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending: {
			OrderStatusProcessing,
			OrderStatusCancelled,
		},
		OrderStatusProcessing: {
			OrderStatusFarmOrder,
			OrderStatusCancelled,
		},
		OrderStatusFarmOrder: {
			OrderStatusCompleted,
			OrderStatusCancelled,
		},
	}

	allowedTransitions, exists := transitions[o.Status]
	if !exists {
		return false
	}

	for _, allowedStatus := range allowedTransitions {
		if newStatus == allowedStatus {
			return true
		}
	}
	return false
}

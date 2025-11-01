package domain

import (
	"context"
)

// OrderRepository интерфейс для работы с заказами
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]*Order, error)
	GetByStatus(ctx context.Context, status OrderStatus) ([]*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Order, error)
}

// CustomerRepository интерфейс для работы с клиентами
type CustomerRepository interface {
	Create(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, id string) (*Customer, error)
	GetByEmail(ctx context.Context, email string) (*Customer, error)
	Update(ctx context.Context, customer *Customer) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Customer, error)
}

// FarmOrderRepository интерфейс для работы с заказами фермы
type FarmOrderRepository interface {
	Create(ctx context.Context, farmOrder *FarmOrder) error
	GetByID(ctx context.Context, id string) (*FarmOrder, error)
	GetByOrderID(ctx context.Context, orderID string) (*FarmOrder, error)
	GetByStatus(ctx context.Context, status FarmOrderStatus) ([]*FarmOrder, error)
	Update(ctx context.Context, farmOrder *FarmOrder) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*FarmOrder, error)
}

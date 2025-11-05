package domain

import "context"

// OrderRepository определяет интерфейс для работы с хранилищем заказов.
// Это позволяет абстрагироваться от конкретной реализации базы данных.
type OrderRepository interface {
	GetAvailableFlowers(ctx context.Context) ([]Item, error)
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByStatus(ctx context.Context, status OrderStatus) ([]*Order, error)
	Update(ctx context.Context, order *Order) error
	Close() error
}

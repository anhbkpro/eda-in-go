package application

import (
	"context"
	"eda-in-golang/search/internal/domain"
)

type (
	GetOrder struct {
		OrderID string
	}

	App interface {
		SearchOrders(ctx context.Context, search domain.SearchFilters) ([]*domain.Order, error)
		GetOrder(ctx context.Context, get GetOrder) (*domain.Order, error)
	}

	Application struct {
		orders domain.OrderRepository
	}
)

var _ App = (*Application)(nil)

func New(orders domain.OrderRepository) *Application {
	return &Application{
		orders: orders,
	}
}

func (a Application) SearchOrders(ctx context.Context, search domain.SearchFilters) ([]*domain.Order, error) {
	// TODO implement me
	panic("implement me")
}

func (a Application) GetOrder(ctx context.Context, get GetOrder) (*domain.Order, error) {
	// TODO implement me
	panic("implement me")
}

package application

import (
	"context"
	"eda-in-golang/search/internal/domain"
)

type (
	GetOrder struct {
		OrderID string
	}

	Application interface {
		SearchOrders(ctx context.Context, search domain.SearchFilters) ([]*domain.Order, error)
		GetOrder(ctx context.Context, get GetOrder) (*domain.Order, error)
	}

	app struct {
		orders domain.OrderRepository
	}
)

var _ Application = (*app)(nil)

func New(orders domain.OrderRepository) *app {
	return &app{
		orders: orders,
	}
}

func (a app) SearchOrders(ctx context.Context, search domain.SearchFilters) ([]*domain.Order, error) {
	// TODO implement me
	panic("implement me")
}

func (a app) GetOrder(ctx context.Context, get GetOrder) (*domain.Order, error) {
	// TODO implement me
	panic("implement me")
}

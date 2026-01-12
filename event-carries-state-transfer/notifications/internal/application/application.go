package application

import (
	"context"

	"eda-in-golang/notifications/internal/domain"
)

type (
	OrderCreated struct {
		OrderID    string
		CustomerID string
	}

	OrderCanceled struct {
		OrderID    string
		CustomerID string
	}

	OrderReady struct {
		OrderID    string
		CustomerID string
	}

	App interface {
		NotifyOrderCreated(context.Context, OrderCreated) error
		NotifyOrderCanceled(context.Context, OrderCanceled) error
		NotifyOrderReady(context.Context, OrderReady) error
	}

	Application struct {
		customers domain.CustomerRepository
	}
)

var _ App = (*Application)(nil)

func New(customers domain.CustomerRepository) *Application {
	return &Application{
		customers: customers,
	}
}

func (a Application) NotifyOrderCreated(ctx context.Context, notify OrderCreated) error {
	// not implemented

	return nil
}

func (a Application) NotifyOrderCanceled(ctx context.Context, notify OrderCanceled) error {
	// not implemented

	return nil
}

func (a Application) NotifyOrderReady(ctx context.Context, notify OrderReady) error {
	// not implemented

	return nil
}

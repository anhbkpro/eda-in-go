package domain

import "context"

type OrderRepository interface {
	Add(ctx context.Context, order *Order) error
	Update(ctx context.Context, orderID string, updater func(*Order) error) error
	Search(ctx context.Context, filters *SearchFilters) ([]*Order, error)
	Get(ctx context.Context, orderID string) (*Order, error)
}

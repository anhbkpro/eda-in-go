package domain

import "context"

type OrderRepository interface {
	Load(ctx context.Context, id string) (*Order, error)
	Save(ctx context.Context, order *Order) error
}

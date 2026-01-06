package models

import "context"

type CustomerCacheRepository interface {
	CustomerRepository
	Add(ctx context.Context, customerID, name string) error
}

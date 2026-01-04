package domain

import "context"

type StoreCacheRepository interface {
	StoreRepository
	Add(ctx context.Context, storeID, name, location string) error
	Rename(ctx context.Context, storeID, name string) error
}

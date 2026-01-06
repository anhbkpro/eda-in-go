package models

import "context"

type StoreCacheRepository interface {
	StoreRepository
	Add(ctx context.Context, storeID, name string) error
	Rename(ctx context.Context, storeID, name string) error
}

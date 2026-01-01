package domain

import "context"

type StoreRepository interface {
	Load(ctx context.Context, storeID string) error
	Save(ctx context.Context, store *Store) error
}

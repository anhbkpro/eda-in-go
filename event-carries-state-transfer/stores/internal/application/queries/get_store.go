package queries

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

// Create Query and Handler for get store
type (
	GetStore struct {
		ID string
	}

	GetStoreHandler struct {
		mall domain.MallRepository
	}
)

func NewGetStoreHandler(mall domain.MallRepository) GetStoreHandler {
	return GetStoreHandler{
		mall: mall,
	}
}

// Implement Handle method
func (h GetStoreHandler) GetStore(ctx context.Context, query GetStore) (*domain.MallStore, error) {
	return h.mall.Find(ctx, query.ID)
}

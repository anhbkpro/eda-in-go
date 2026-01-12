package queries

import (
	"context"

	"eda-in-golang/stores/internal/domain"
)

// Create Query and Handler for get stores
type (
	GetStores struct {
		StoreID string
	}

	GetStoresHandler struct {
		mall domain.MallRepository
	}
)

func NewGetStoresHandler(mall domain.MallRepository) GetStoresHandler {
	return GetStoresHandler{
		mall: mall,
	}
}

// Implement Handle method
func (h GetStoresHandler) GetStores(ctx context.Context, query GetStores) ([]*domain.MallStore, error) {
	return h.mall.All(ctx)
}

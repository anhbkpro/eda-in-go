package queries

import (
	"context"
	"eda-in-golang/stores/internal/domain"
)

// Create Query and Handler for get participating stores
type (
	GetParticipatingStores struct {
		StoreID string
	}

	GetParticipatingStoresHandler struct {
		mall domain.MallRepository
	}
)

func NewGetParticipatingStoresHandler(mall domain.MallRepository) GetParticipatingStoresHandler {
	return GetParticipatingStoresHandler{
		mall: mall,
	}
}

// Implement Handle method
func (h GetParticipatingStoresHandler) GetParticipatingStores(ctx context.Context, query GetParticipatingStores) ([]*domain.MallStore, error) {
	return h.mall.AllParticipating(ctx)
}

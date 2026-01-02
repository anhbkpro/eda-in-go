package queries

import (
	"context"
	"eda-in-golang/stores/internal/domain"
)

// Create Query and Handler for get product
type (
	GetProduct struct {
		ID string
	}

	GetProductHandler struct {
		catalog domain.CatalogRepository
	}
)

func NewGetProductHandler(catalog domain.CatalogRepository) GetProductHandler {
	return GetProductHandler{
		catalog: catalog,
	}
}

// Implement Handle method
func (h GetProductHandler) GetProduct(ctx context.Context, query GetProduct) (*domain.CatalogProduct, error) {
	return h.catalog.Find(ctx, query.ID)
}

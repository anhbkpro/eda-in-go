package commands

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/stores/internal/domain"
)

// Create Command and Handler for remove product
type (
	RemoveProduct struct {
		ID string
	}

	RemoveProductHandler struct {
		products domain.ProductRepository
	}
)

func NewRemoveProductHandler(products domain.ProductRepository) RemoveProductHandler {
	return RemoveProductHandler{
		products: products,
	}
}

// Implement Handle method
func (h RemoveProductHandler) RemoveProduct(ctx context.Context, cmd RemoveProduct) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}

	err = product.Remove()
	if err != nil {
		return errors.Wrap(err, "error removing product")
	}

	return errors.Wrap(h.products.Save(ctx, product), "error saving product")
}

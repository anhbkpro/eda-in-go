package commands

import (
	"context"
	"eda-in-golang/stores/internal/domain"

	"github.com/stackus/errors"
)

// Create Command and Handler for decrease product price
type (
	DecreaseProductPrice struct {
		ID    string
		Price float64
	}

	DecreaseProductPriceHandler struct {
		products domain.ProductRepository
	}
)

func NewDecreaseProductPriceHandler(products domain.ProductRepository) DecreaseProductPriceHandler {
	return DecreaseProductPriceHandler{
		products: products,
	}
}

// Implement Handle method
func (h DecreaseProductPriceHandler) Handle(ctx context.Context, cmd DecreaseProductPrice) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}

	err = product.DecreasePrice(cmd.Price)
	if err != nil {
		return errors.Wrap(err, "error decreasing product price")
	}

	return errors.Wrap(h.products.Save(ctx, product), "error saving product")
}

func (h DecreaseProductPriceHandler) DecreaseProductPrice(ctx context.Context, cmd DecreaseProductPrice) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}

	err = product.DecreasePrice(cmd.Price)
	if err != nil {
		return errors.Wrap(err, "error decreasing product price")
	}

	return errors.Wrap(h.products.Save(ctx, product), "error saving product")
}

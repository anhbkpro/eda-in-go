package commands

import (
	"context"
	"eda-in-golang/stores/internal/domain"

	"github.com/stackus/errors"
)

// Create Command and Handler for increase product price
type (
	IncreaseProductPrice struct {
		ID    string
		Price float64
	}

	IncreaseProductPriceHandler struct {
		products domain.ProductRepository
	}
)

func NewIncreaseProductPriceHandler(products domain.ProductRepository) IncreaseProductPriceHandler {
	return IncreaseProductPriceHandler{
		products: products,
	}
}

// Implement Handle method
func (h IncreaseProductPriceHandler) Handle(ctx context.Context, cmd IncreaseProductPrice) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}

	err = product.IncreasePrice(cmd.Price)
	if err != nil {
		return errors.Wrap(err, "error increasing product price")
	}

	return errors.Wrap(h.products.Save(ctx, product), "error saving product")
}

func (h IncreaseProductPriceHandler) IncreaseProductPrice(ctx context.Context, cmd IncreaseProductPrice) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}

	err = product.IncreasePrice(cmd.Price)
	if err != nil {
		return errors.Wrap(err, "error increasing product price")
	}

	return errors.Wrap(h.products.Save(ctx, product), "error saving product")
}

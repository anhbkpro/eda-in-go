package commands

import (
	"context"

	"github.com/stackus/errors"

	"eda-in-golang/stores/internal/domain"
)

// Create Command and Handler for rebrand product
type (
	RebrandProduct struct {
		ID          string
		Name        string
		Description string
	}

	RebrandProductHandler struct {
		products domain.ProductRepository
	}
)

func NewRebrandProductHandler(products domain.ProductRepository) RebrandProductHandler {
	return RebrandProductHandler{
		products: products,
	}
}

// Implement Handle method
func (h RebrandProductHandler) RebrandProduct(ctx context.Context, cmd RebrandProduct) error {
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}

	err = product.Rebrand(cmd.Name, cmd.Description)
	if err != nil {
		return errors.Wrap(err, "error rebranding product")
	}

	return errors.Wrap(h.products.Save(ctx, product), "error saving product")
}

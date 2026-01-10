package commands

import (
	"context"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/stores/internal/domain"
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
	fmt.Println("[Step 4] Handler.DecreaseProductPrice: received command")

	fmt.Println("[Step 5] Handler → Repo.Load: loading product aggregate")
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}
	fmt.Println("[Step 9] Repo → Handler: product loaded successfully")

	fmt.Println("[Step 10] Handler → Product.DecreasePrice: calling domain logic")
	err = product.DecreasePrice(cmd.Price)
	if err != nil {
		return errors.Wrap(err, "error decreasing product price")
	}
	fmt.Println("[Step 11] Product.AddEvent: ProductPriceDecreasedEvent added")

	fmt.Println("[Step 12] Handler → Repo.Save: saving product aggregate")
	err = h.products.Save(ctx, product)
	fmt.Println("[Step 25] Repo → Handler: save completed")

	return errors.Wrap(err, "error saving product")
}

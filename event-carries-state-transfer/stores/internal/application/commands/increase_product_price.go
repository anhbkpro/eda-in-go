package commands

import (
	"context"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/stores/internal/domain"
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
	fmt.Println("[Step 4] Handler.IncreaseProductPrice: received command")

	fmt.Println("[Step 5] Handler → Repo.Load: loading product aggregate")
	product, err := h.products.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}
	fmt.Println("[Step 9] Repo → Handler: product loaded successfully")

	fmt.Println("[Step 10] Handler → Product.IncreasePrice: calling domain logic")
	err = product.IncreasePrice(cmd.Price)
	if err != nil {
		return errors.Wrap(err, "error increasing product price")
	}
	fmt.Println("[Step 11] Product.AddEvent: ProductPriceIncreasedEvent added")

	fmt.Println("[Step 12] Handler → Repo.Save: saving product aggregate")
	err = h.products.Save(ctx, product)
	fmt.Println("[Step 25] Repo → Handler: save completed")

	return errors.Wrap(err, "error saving product")
}

package application

import (
	"context"
	"fmt"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

type CatalogHandlers[T ddd.AggregateEvent] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*CatalogHandlers[ddd.AggregateEvent])(nil)

func NewCatalogHandlers(catalog domain.CatalogRepository) CatalogHandlers[ddd.AggregateEvent] {
	return CatalogHandlers[ddd.AggregateEvent]{
		catalog: catalog,
	}
}

func (h CatalogHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	fmt.Printf("[Step 18] Dispatcher → CatalogHandler.HandleEvent: received %s\n", event.EventName())
	switch event.EventName() {
	case domain.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case domain.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case domain.ProductPriceIncreasedEvent:
		fmt.Println("[Step 19] CatalogHandler → onProductPriceIncreased")
		return h.onProductPriceIncreased(ctx, event)
	case domain.ProductPriceDecreasedEvent:
		fmt.Println("[Step 19] CatalogHandler → onProductPriceDecreased")
		return h.onProductPriceDecreased(ctx, event)
	case domain.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	}

	return nil
}

func (h CatalogHandlers[T]) onProductAdded(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ProductAdded)
	return h.catalog.AddProduct(ctx, event.AggregateID(), payload.StoreID, payload.Name, payload.Description, payload.SKU, payload.Price)
}

func (h CatalogHandlers[T]) onProductRebranded(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ProductRebranded)
	return h.catalog.Rebrand(ctx, event.AggregateID(), payload.Name, payload.Description)
}

func (h CatalogHandlers[T]) onProductPriceIncreased(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ProductPriceChanged)
	fmt.Printf("[Step 20] CatalogHandler → CatalogRepo.UpdatePrice: delta=+%.2f\n", payload.Delta)
	err := h.catalog.UpdatePrice(ctx, event.AggregateID(), payload.Delta)
	fmt.Println("[Step 22] CatalogRepo → CatalogHandler: price updated in read model")
	return err
}

func (h CatalogHandlers[T]) onProductPriceDecreased(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ProductPriceChanged)
	fmt.Printf("[Step 20] CatalogHandler → CatalogRepo.UpdatePrice: delta=%.2f\n", payload.Delta)
	err := h.catalog.UpdatePrice(ctx, event.AggregateID(), payload.Delta)
	fmt.Println("[Step 22] CatalogRepo → CatalogHandler: price updated in read model")
	return err
}

func (h CatalogHandlers[T]) onProductRemoved(ctx context.Context, event T) error {
	return h.catalog.RemoveProduct(ctx, event.AggregateID())
}

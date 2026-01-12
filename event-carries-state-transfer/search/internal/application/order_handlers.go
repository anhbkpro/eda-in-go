package application

import (
	"context"
	"time"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/orderingpb"
	"eda-in-golang/search/internal/domain"
)

type OrderHandlers[T ddd.Event] struct {
	orders    domain.OrderRepository
	customers domain.CustomerRepository
	stores    domain.StoreRepository
	products  domain.ProductRepository
}

var _ ddd.EventHandler[ddd.Event] = (*OrderHandlers[ddd.Event])(nil)

func NewOrderHandlers(orders domain.OrderRepository, customers domain.CustomerRepository, stores domain.StoreRepository, products domain.ProductRepository) OrderHandlers[ddd.Event] {
	return OrderHandlers[ddd.Event]{
		orders:    orders,
		customers: customers,
		stores:    stores,
		products:  products,
	}
}

func (h OrderHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case orderingpb.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case orderingpb.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	case orderingpb.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)
	case orderingpb.OrderCompletedEvent:
		return h.onOrderCompleted(ctx, event)
	}
	return nil
}

// [integration-event-flow.md] 5. Integration handler transforms and publishes (different service)
// Convert orderingpb.OrderCreated → domain.Order
func (h OrderHandlers[T]) onOrderCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCreated)

	customer, err := h.customers.Find(ctx, payload.GetCustomerId())
	if err != nil {
		return err
	}

	var total float64
	items := make([]domain.Item, 0, len(payload.GetItems()))
	seenStores := make(map[string]*domain.Store) // to avoid duplicate store lookups
	for _, item := range payload.GetItems() {
		product, err := h.products.Find(ctx, item.GetProductId())
		if err != nil {
			return err
		}
		var store *domain.Store
		var exists bool

		if store, exists = seenStores[product.StoreID]; !exists {
			store, err = h.stores.Find(ctx, product.StoreID)
			if err != nil {
				return err
			}
			seenStores[product.StoreID] = store
		}

		items = append(items, domain.Item{
			ProductID:   item.GetProductId(),
			StoreID:     store.ID,
			StoreName:   store.Name,
			ProductName: product.Name,
			Price:       item.GetPrice(),
			Quantity:    int(item.GetQuantity()),
		})
		total += float64(item.GetQuantity()) * item.GetPrice()
	}

	order := &domain.Order{
		ID:           payload.GetId(),
		CustomerID:   customer.ID,
		CustomerName: customer.Name,
		Items:        items,
		Total:        total,
		Status:       "New",
		CreatedAt:    time.Now(),
	}
	return h.orders.Add(ctx, order)
}

func (h OrderHandlers[T]) onOrderReadied(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderReadied)
	return h.orders.Update(ctx, payload.GetId(), func(o *domain.Order) error {
		o.Status = "Ready For Pickup"
		return nil
	})
}

func (h OrderHandlers[T]) onOrderCanceled(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCanceled)
	return h.orders.Update(ctx, payload.GetId(), func(o *domain.Order) error {
		o.Status = "Canceled"
		return nil
	})
}

func (h OrderHandlers[T]) onOrderCompleted(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCompleted)
	return h.orders.Update(ctx, payload.GetId(), func(o *domain.Order) error {
		o.Status = "Completed"
		return nil
	})
}

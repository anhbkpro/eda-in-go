package application

import (
	"context"
	"eda-in-golang/customers/customerspb"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/notifications/internal/domain"
)

type CustomerHandlers[T ddd.Event] struct {
	cache domain.CustomerCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*CustomerHandlers[ddd.Event])(nil)

func NewCustomerHandlers(cache domain.CustomerCacheRepository) CustomerHandlers[ddd.Event] {
	return CustomerHandlers[ddd.Event]{
		cache: cache,
	}
}

func (h CustomerHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case customerspb.CustomerRegisteredEvent:
		return h.onCustomerRegistered(ctx, event)
	case customerspb.CustomerSmsChangedEvent:
		return h.onCustomerSmsChanged(ctx, event)
	}

	return nil
}

func (h CustomerHandlers[T]) onCustomerRegistered(ctx context.Context, event T) error {
	payload := event.Payload().(*customerspb.CustomerRegistered)
	return h.cache.Add(ctx, payload.Id, payload.Name, payload.SmsNumber)
}

func (h CustomerHandlers[T]) onCustomerSmsChanged(ctx context.Context, event T) error {
	payload := event.Payload().(*customerspb.CustomerSmsChanged)
	return h.cache.UpdateSmsNumber(ctx, payload.Id, payload.SmsNumber)
}

package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

func RegisterIntegrationEventHandlers(handlers ddd.EventHandler[ddd.AggregateEvent], dispatcher *ddd.EventDispatcher[ddd.AggregateEvent]) {
	dispatcher.Subscribe(
		handlers,
		domain.OrderCreatedEvent,
		domain.OrderReadiedEvent,
		domain.OrderCanceledEvent,
		domain.OrderCompletedEvent,
	)
}

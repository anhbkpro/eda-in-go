package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/internal/domain"
)

func RegisterIntegrationEventHandlers(eventHandlers ddd.EventHandler[ddd.AggregateEvent], dispatcher *ddd.EventDispatcher[ddd.AggregateEvent]) {
	dispatcher.Subscribe(
		eventHandlers,
		domain.OrderCreatedEvent,
		domain.OrderReadiedEvent,
		domain.OrderCanceledEvent,
		domain.OrderCompletedEvent,
	)
}

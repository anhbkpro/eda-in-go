package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

func RegisterIntegrationEventHandlers(handler ddd.EventHandler[ddd.AggregateEvent], dispatcher *ddd.EventDispatcher[ddd.AggregateEvent]) {
	dispatcher.Subscribe(handler,
		domain.StoreCreatedEvent,
		domain.StoreParticipationEnabledEvent,
		domain.StoreParticipationDisabledEvent,
		domain.StoreRebrandedEvent)
}

package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

func RegisterMallHandlers(mallHandler ddd.EventHandler[ddd.AggregateEvent], eventSubsricber ddd.EventSubscriber[ddd.AggregateEvent]) {
	eventSubsricber.Subscribe(mallHandler,
		domain.StoreCreatedEvent,
		domain.StoreParticipationEnabledEvent,
		domain.StoreParticipationDisabledEvent,
		domain.StoreRebrandedEvent)
}

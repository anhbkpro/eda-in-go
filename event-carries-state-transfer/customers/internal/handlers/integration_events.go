package handlers

import (
	"eda-in-golang/customers/internal/domain"
	"eda-in-golang/internal/ddd"
)

func RegisterIntegrationEventHandlers(handler ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(handler,
		domain.CustomerRegisteredEvent,
		domain.CustomerSmsChangedEvent,
		domain.CustomerEnabledEvent,
		domain.CustomerDisabledEvent,
	)
}

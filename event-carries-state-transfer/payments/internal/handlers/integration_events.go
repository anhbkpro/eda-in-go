package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/payments/internal/domain"
)

func RegisterIntegrationEventHandlers(eventHandlers ddd.EventHandler[ddd.Event], domainSubscriber ddd.EventSubscriber[ddd.Event]) {
	domainSubscriber.Subscribe(eventHandlers,
		domain.InvoicePaidEvent,
	)
}

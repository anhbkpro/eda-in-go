package handlers

import (
	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

func RegisterCatalogHandlers(catalogHandler ddd.EventHandler[ddd.AggregateEvent], domainSubscriber ddd.EventSubscriber[ddd.AggregateEvent]) {
	domainSubscriber.Subscribe(catalogHandler,
		domain.ProductAddedEvent,
		domain.ProductRebrandedEvent,
		domain.ProductPriceIncreasedEvent,
		domain.ProductPriceDecreasedEvent,
		domain.ProductRemovedEvent)
}

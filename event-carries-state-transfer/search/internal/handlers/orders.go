package handlers

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
	"eda-in-golang/ordering/orderingpb"
)

func RegisterOrderHandlers(orderHandlers ddd.EventHandler[ddd.Event], stream am.EventSubscriber) error {
	evtMsgHandler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, msg am.EventMessage) error {
		return orderHandlers.HandleEvent(ctx, msg)
	})

	// [integration-event-flow.md] 4. Other services consume integration events
	return stream.Subscribe(orderingpb.OrderAggregateChannel, evtMsgHandler, am.MessageFilter{
		orderingpb.OrderCreatedEvent,
		orderingpb.OrderCanceledEvent,
		orderingpb.OrderReadiedEvent,
		orderingpb.OrderCompletedEvent,
	}, am.GroupName("search-orders"))
}

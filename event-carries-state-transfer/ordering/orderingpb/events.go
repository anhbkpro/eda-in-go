package orderingpb

import (
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/serdes"
)

const (
	OrderAggregateChannel = "mallbots.ordering.events.Order"

	// Define event names here, event payloads are defined in the events.proto file
	OrderCreatedEvent   = "ordering.api.OrderCreated"
	OrderReadiedEvent   = "ordering.api.OrderReadied"
	OrderCompletedEvent = "ordering.api.OrderCompleted"
	OrderCanceledEvent  = "ordering.api.OrderCanceled"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerdes(reg)

	// Order events
	if err := serde.Register(&OrderCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderReadied{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderCompleted{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderCanceled{}); err != nil {
		return err
	}
	return nil
}

// Event payloads have Key() methods that return the event name
// ==> this will be used for registering events with sedes ???
func (*OrderCreated) Key() string   { return OrderCreatedEvent }
func (*OrderReadied) Key() string   { return OrderReadiedEvent }
func (*OrderCanceled) Key() string  { return OrderCanceledEvent }
func (*OrderCompleted) Key() string { return OrderCompletedEvent }

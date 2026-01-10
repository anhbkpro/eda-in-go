# Integration Event Flow

**Integration Event Handlers** convert internal domain events into external integration events that other services can consume.

## Event Flow Pattern

```
Domain Event → Event Dispatcher → Integration Handlers → Message Publisher → NATS Stream
```

## Key Components

- **Domain Events**: Internal business logic events (e.g., `OrderCreated`)
- **Event Dispatcher**: Routes domain events to subscribers within the service
- **Integration Handlers**: Transform domain events to protobuf integration events
- **Message Publisher**: Publishes integration events to NATS streams
- **Integration Events**: External API events (e.g., `OrderCreatedEvent`) consumed by other services

## Example: Order Service Integration

```go
// 1. Domain event generated in aggregate
order.CreateOrder(...) // → OrderCreated domain event

// 2. Event dispatcher routes to integration handlers
domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()
handlers.RegisterIntegrationEventHandlers(integrationEventHandlers, domainDispatcher)

// 3. Integration handler transforms and publishes
func (h *IntegrationEventHandlers) onOrderCreated(ctx, event) error {
    // Convert domain.OrderCreated → orderingpb.OrderCreated
    return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel,
        ddd.NewEvent(orderingpb.OrderCreatedEvent, &orderingpb.OrderCreated{...}))
}

// 4. Other services consume integration events
eventStream.Subscribe(orderingpb.OrderAggregateChannel, handler,
    orderingpb.OrderCreatedEvent)
```

## Benefits

- **Clean Architecture**: Internal domain events stay internal
- **Loose Coupling**: Services react to business events, not implementation details
- **Event-Driven Flow**: Maintains asynchronous, reactive architecture
- **API Stability**: Integration events use stable protobuf contracts

## Implementation Details

### Event Handler Registration

```go
handlers.RegisterIntegrationEventHandlers(integrationEventHandlers, domainDispatcher)
```

This registers the integration event handlers as subscribers to domain events within the service's event dispatcher.

### Domain to Integration Event Mapping

| Domain Event | Integration Event | Purpose |
|-------------|------------------|---------|
| `OrderCreated` | `OrderCreatedEvent` | New order notification |
| `OrderReadied` | `OrderReadiedEvent` | Order ready for pickup |
| `OrderCanceled` | `OrderCanceledEvent` | Order cancellation |
| `OrderCompleted` | `OrderCompletedEvent` | Order fulfillment |

### Publisher Setup

Integration handlers receive a message publisher (typically NATS JetStream) that publishes events to named channels that other services can subscribe to.

## Usage Patterns

### Publishing Integration Events

1. Domain logic generates domain events
2. Event dispatcher routes to integration handlers
3. Integration handlers transform events
4. Events published to external channels

### Consuming Integration Events

Other services subscribe to integration event channels to react to business activities:

```go
eventStream.Subscribe(orderingpb.OrderAggregateChannel, orderHandler,
    orderingpb.OrderCreatedEvent, orderingpb.OrderCompletedEvent)
```

This pattern enables the **search service** and other components to index orders, send notifications, and trigger downstream business processes without tight coupling to individual service implementations.

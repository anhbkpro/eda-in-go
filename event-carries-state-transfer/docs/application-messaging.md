# Application Messaging (AM) Module

This document explains the **Application Messaging (AM)** module, which provides the messaging abstraction layer for event-driven communication in this Go application.

---

## 1. Overview

The AM module (`internal/am`) serves as a **generic messaging abstraction** that bridges domain events with message brokers, enabling clean event-driven architecture across microservices.

> **Domain Events → AM Layer → Message Broker**

The AM module provides:
- **Type-safe event publishing and subscribing**
- **Message acknowledgment and lifecycle management**
- **Event serialization/deserialization**
- **Consumer group and filtering capabilities**

---

## 2. Core Concepts

### Message Interface

All messages in the AM system implement the `Message` interface:

```go
type Message interface {
    ddd.IDer
    MessageName() string
    Ack() error   // Acknowledge successful processing
    NAck() error  // Negative acknowledge (retry)
    Extend() error // Extend processing time
    Kill() error   // Stop processing permanently
}
```

### Generic Messaging Interfaces

```go
// Publishing messages of type I
MessagePublisher[I any] interface {
    Publish(ctx context.Context, topicName string, v I) error
}

// Subscribing to messages of type O
MessageSubscriber[O Message] interface {
    Subscribe(topicName string, handler MessageHandler[O], options ...SubscriberOption) error
}

// Combined publisher/subscriber
MessageStream[I any, O Message] interface {
    MessagePublisher[I]
    MessageSubscriber[O]
}
```

---

## 3. Event Stream Layer

### EventMessage Interface

Domain events are wrapped in `EventMessage`:

```go
type EventMessage interface {
    Message
    ddd.Event  // ID, EventName, Payload, Metadata, OccurredAt
}
```

### EventStream Implementation

The `EventStream` provides domain-aware messaging:

```go
type EventStream = MessageStream[ddd.Event, EventMessage]

func NewEventStream(reg registry.Registry, stream MessageStream[RawMessage, RawMessage]) EventStream
```

Key features:
- **Serialization**: Uses registry to serialize/deserialize event payloads
- **Protobuf encoding**: Events are encoded as `EventMessageData` protobuf messages
- **Metadata preservation**: Event metadata is preserved through the messaging layer

---

## 4. Subscriber Configuration

### Configuration Options

```go
type SubscriberConfig struct {
    msgFilter    []string       // Filter by message names
    groupName    string         // Consumer group name
    ackType      AckType        // Auto or Manual acknowledgment
    ackWait      time.Duration  // Ack wait timeout
    maxRedeliver int           // Maximum redelivery attempts
}
```

### Subscriber Options

```go
// Filter messages by event type
am.MessageFilter("OrderCreated", "OrderCancelled")

// Set consumer group for load balancing
am.GroupName("order-processors")

// Configure acknowledgment behavior
am.AckWait(10 * time.Second)
am.MaxRedeliver(3)
```

---

## 5. Usage Patterns

### Publishing Events

```go
// In a service module
eventStream := am.NewEventStream(reg, jetstream.NewStream(streamName, js))

// Publish domain event
event := ddd.NewEvent("OrderCreated", orderCreated{
    OrderID: "order-123",
    CustomerID: "customer-456",
})

err := eventStream.Publish(ctx, "mallbots.events", event)
```

### Subscribing to Events

```go
// Define event handler
handler := am.MessageHandlerFunc[am.EventMessage](func(ctx context.Context, msg am.EventMessage) error {
    switch msg.EventName() {
    case "OrderCreated":
        return handleOrderCreated(ctx, msg)
    case "OrderCancelled":
        return handleOrderCancelled(ctx, msg)
    }
    return nil
})

// Subscribe with filtering
err := eventStream.Subscribe(
    "mallbots.events",
    handler,
    am.MessageFilter("OrderCreated", "OrderCancelled"),
    am.GroupName("order-handlers"),
)
```

### Integration in Services

All services use the AM module for cross-service communication:

```go
func (m Module) Startup(ctx context.Context, mono monolith.Monolith) error {
    // Create event stream
    reg := registry.New()
    eventStream := am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))

    // Use in application layer
    app := application.New(eventStream, /* other deps */)

    // Register integration event handlers
    integrationEventHandlers := application.NewIntegrationEventHandlers(eventStream)

    return nil
}
```

---

## 6. Message Flow

### Publishing Flow

1. **Domain Layer**: Business logic creates domain events
2. **Application Layer**: Events are published via `EventStream.Publish()`
3. **AM Layer**: Events are serialized using registry
4. **Protobuf Layer**: Events wrapped in `EventMessageData`
5. **Transport Layer**: Messages sent to JetStream/NATS

### Subscribing Flow

1. **Transport Layer**: Raw messages received from JetStream/NATS
2. **Protobuf Layer**: `EventMessageData` unmarshaled
3. **AM Layer**: Events deserialized using registry
4. **Application Layer**: Domain events passed to handlers
5. **Acknowledgment**: Success/failure communicated back

---

## 7. Error Handling & Reliability

### Acknowledgment Types

```go
const (
    AckTypeAuto   AckType = iota  // Auto-acknowledge on handler return
    AckTypeManual                 // Manual ack/nack control
)
```

### Message Lifecycle

```go
// Successful processing
msg.Ack()

// Temporary failure (retry)
msg.NAck()

// Extend processing time
msg.Extend()

// Permanent failure (stop retries)
msg.Kill()
```

### Redelivery Configuration

```go
subscriber.Subscribe(topic, handler,
    am.AckWait(30 * time.Second),     // Wait 30s for ack
    am.MaxRedeliver(5),               // Retry up to 5 times
)
```

---

## 8. Architecture Benefits

### Clean Separation

- **Domain Layer**: Pure business logic, no messaging concerns
- **Application Layer**: Orchestrates events and messaging
- **Infrastructure Layer**: AM abstracts message broker details

### Testability

```go
// Mock EventStream for testing
type mockEventStream struct {
    published []ddd.Event
}

func (m *mockEventStream) Publish(ctx context.Context, topic string, event ddd.Event) error {
    m.published = append(m.published, event)
    return nil
}
```

### Broker Agnostic

The AM interfaces allow switching between different message brokers:

```go
// Could use Kafka instead of JetStream
kafkaStream := kafka.NewStream(brokers, config)
eventStream := am.NewEventStream(reg, kafkaStream)
```

---

## 9. Integration Points

### Registry System

AM integrates with the registry for event serialization:

```go
reg := registry.New()
// Register event types
orderingpb.Registrations(reg)

// AM uses registry for serialization
eventStream := am.NewEventStream(reg, stream)
```

### JetStream Transport

Currently uses JetStream (NATS) as the message broker:

```go
stream := jetstream.NewStream(streamName, jsConn)
eventStream := am.NewEventStream(reg, stream)
```

### Monolith Integration

Services register with the monolith container:

```go
func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {
    eventStream := am.NewEventStream(reg, jetstream.NewStream(mono.Config().Nats.Stream, mono.JS()))
    // ... use eventStream
}
```

---

## 10. Common Patterns

### Event-Driven Sagas

```go
// Order service publishes events
eventStream.Publish(ctx, "events", ddd.NewEvent("OrderCreated", orderData))

// Payment service reacts
eventStream.Subscribe("events", paymentHandler,
    am.MessageFilter("OrderCreated"))
```

### CQRS with Events

```go
// Command side publishes events
eventStream.Publish(ctx, "events", ddd.NewEvent("ProductUpdated", productData))

// Query side maintains read models
eventStream.Subscribe("events", searchHandler,
    am.MessageFilter("ProductUpdated"))
```

### Service Integration

```go
// Each service defines its own event handlers
integrationHandlers := application.NewIntegrationEventHandlers(eventStream)

// Register with domain event dispatcher
handlers.RegisterIntegrationEventHandlers(integrationHandlers, domainDispatcher)
```

---

## 11. Configuration Examples

### High-Throughput Consumer

```go
eventStream.Subscribe("events", handler,
    am.GroupName("high-throughput"),
    am.AckType(am.AckTypeAuto),  // Auto-ack for speed
    am.AckWait(5 * time.Second),
)
```

### Reliable Consumer

```go
eventStream.Subscribe("events", handler,
    am.GroupName("reliable-processing"),
    am.AckType(am.AckTypeManual),  // Manual control
    am.AckWait(60 * time.Second),  // Longer timeout
    am.MaxRedeliver(10),           // More retries
)
```

### Filtered Consumer

```go
eventStream.Subscribe("events", handler,
    am.MessageFilter("OrderCreated", "PaymentCompleted", "ShipmentDelivered"),
    am.GroupName("order-lifecycle"),
)
```

---

## 12. Best Practices

### 1. Use Manual Ack for Critical Operations

```go
// Critical business logic should manually control acknowledgment
func handlePayment(ctx context.Context, msg am.EventMessage) error {
    if err := processPayment(msg); err != nil {
        msg.NAck() // Retry
        return err
    }
    msg.Ack() // Success
    return nil
}
```

### 2. Filter Messages Appropriately

```go
// Subscribe only to relevant events
am.MessageFilter("OrderCreated", "OrderCancelled")
// Don't subscribe to all events with "*"
```

### 3. Configure Timeouts Based on Processing Time

```go
// If processing takes ~30 seconds, set ack wait to 60 seconds
am.AckWait(60 * time.Second)
```

### 4. Use Consumer Groups for Load Balancing

```go
// Multiple instances share the load
am.GroupName("order-processors")
```

---

## 13. Troubleshooting

### Common Issues

**Messages not being processed:**
- Check consumer group configuration
- Verify message filtering
- Check JetStream stream configuration

**Duplicate processing:**
- Ensure proper acknowledgment
- Check for multiple consumers with same group name

**Message loss:**
- Verify acknowledgment configuration
- Check JetStream persistence settings

**Performance issues:**
- Review acknowledgment timeouts
- Consider auto-ack for high-throughput scenarios

---

## 14. Related Components

| Component | Responsibility |
|-----------|----------------|
| `internal/am` | Messaging abstraction layer |
| `internal/jetstream` | JetStream/NATS transport |
| `internal/registry` | Event serialization |
| `internal/ddd` | Domain event interfaces |
| `internal/es` | Event sourcing |

---

This AM module provides the foundation for reliable, scalable event-driven communication across all services in the system.

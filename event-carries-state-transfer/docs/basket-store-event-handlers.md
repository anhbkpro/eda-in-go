# Basket Service Store Event Handlers

## Overview

The `store_handlers.go` file implements event-driven synchronization between the **Stores Service** and **Baskets Service** in the event-carries-state-transfer architecture. It maintains a local cache of store information within the baskets service to avoid cross-service calls during basket operations.

## What It Is

This file contains an **event handler** that listens for store-related domain events and updates the baskets service's internal store cache. It's part of the **event-driven data synchronization** pattern used throughout the system.

## Why It's Needed

### Performance Optimization
- **Avoid Cross-Service Calls**: Basket operations need store names for display purposes, but making real-time calls to the stores service would create tight coupling and performance bottlenecks
- **Cache-First Architecture**: By maintaining a local cache, basket operations can retrieve store information instantly without network overhead

### Data Consistency
- **Eventual Consistency**: Store data stays synchronized through events rather than direct queries
- **Decoupled Services**: Baskets service doesn't depend on stores service availability for basic operations

### Business Requirements
- **Store Information in Baskets**: When displaying basket contents, users see store names associated with products
- **Real-time Updates**: Store rebranding should reflect immediately in existing baskets

## When It's Used

### Event Triggers
The handler responds to these specific events from the stores service:

1. **`storespb.StoreCreatedEvent`** - When a new store is created
2. **`storespb.StoreRebrandedEvent`** - When an existing store changes its name

### Execution Context
- **Asynchronous Processing**: Events are processed in the background through the message bus
- **Transaction Boundaries**: Updates happen outside the original store transaction for better decoupling
- **Error Isolation**: Handler failures don't affect the original store operations

## How It Works

### Architecture Pattern

```go
// Event Handler Interface
type StoreHandlers[T ddd.Event] struct {
    cache domain.StoreCacheRepository
}
```

The handler implements the **Event Handler Pattern** with generic constraints for type safety.

### Event Processing Flow

#### 1. Event Reception
```go
func (h StoreHandlers[T]) HandleEvent(ctx context.Context, event T) error {
    switch event.EventName() {
    case storespb.StoreCreatedEvent:
        return h.onStoreCreated(ctx, event)
    case storespb.StoreRebrandedEvent:
        return h.onStoreRebranded(ctx, event)
    }
    return nil
}
```

#### 2. Store Creation Handler
```go
func (h StoreHandlers[T]) onStoreCreated(ctx context.Context, event ddd.Event) error {
    payload := event.Payload().(*storespb.StoreCreated)
    return h.cache.Add(ctx, payload.GetId(), payload.GetName())
}
```
- Extracts store ID and name from the event payload
- Adds the store to the local cache

#### 3. Store Rebranding Handler
```go
func (h StoreHandlers[T]) onStoreRebranded(ctx context.Context, event ddd.Event) error {
    payload := event.Payload().(*storespb.StoreRebranded)
    return h.cache.Rename(ctx, payload.GetId(), payload.GetName())
}
```
- Updates the store name in the cache for the given ID

### Cache Repository Interface

The handler depends on a `StoreCacheRepository` interface:

```go
type StoreCacheRepository interface {
    StoreRepository
    Add(ctx context.Context, storeID, name string) error
    Rename(ctx context.Context, storeID, name string) error
}
```

This interface extends the basic `StoreRepository` with cache-specific operations.

## Integration Points

### Service Registration
The handler is registered in the basket service's module configuration:

```go
// baskets/module.go
storeHandlers := application.NewStoreHandlers(storeCache)
```

### Message Bus Subscription
The baskets service subscribes to store events through the event streaming infrastructure, ensuring the handler receives relevant events.

### Cache Implementation
The actual cache is typically implemented using:
- **In-memory storage** for fast access
- **Database persistence** for durability across service restarts
- **TTL mechanisms** for cache invalidation

## Benefits

### Performance
- **Zero-latency lookups**: Store names retrieved instantly from local cache
- **Reduced network traffic**: No cross-service calls during basket operations
- **Scalability**: Basket service can handle high loads without depending on stores service

### Reliability
- **Fault tolerance**: Basket operations work even if stores service is down
- **Eventual consistency**: Data stays synchronized through reliable event streaming
- **Graceful degradation**: Cached data remains available during network issues

### Maintainability
- **Loose coupling**: Services communicate through events rather than direct API calls
- **Independent deployment**: Services can be updated without coordinating deployments
- **Clear separation of concerns**: Each service owns its domain data

## Error Handling

### Failure Scenarios
- **Cache update failures**: Logged but don't affect the original store operation
- **Event deserialization errors**: Handler returns error, event may be retried
- **Database connection issues**: Cache operations fail gracefully

### Monitoring
- **Event processing metrics**: Track handler performance and success rates
- **Cache consistency checks**: Validate cache data against source of truth
- **Error alerting**: Notify when cache updates consistently fail

## Testing Strategy

### Unit Tests
- **Event handling**: Verify correct cache operations for each event type
- **Error scenarios**: Test handler behavior with invalid events
- **Type safety**: Ensure generic constraints work correctly

### Integration Tests
- **End-to-end synchronization**: Verify store changes propagate to basket cache
- **Event ordering**: Test behavior with out-of-order events
- **Cache consistency**: Validate cached data matches source data

### Performance Tests
- **Cache lookup speed**: Ensure sub-millisecond response times
- **Concurrent updates**: Test behavior under high event load
- **Memory usage**: Monitor cache size and growth patterns

## Future Enhancements

### Potential Improvements
- **Cache warming**: Pre-populate cache on service startup
- **Cache invalidation**: Implement TTL-based expiration
- **Bulk operations**: Handle batch store updates efficiently
- **Cache versioning**: Support multiple cache versions for rolling deployments

### Monitoring Enhancements
- **Cache hit/miss ratios**: Track cache effectiveness
- **Event processing latency**: Measure end-to-end synchronization time
- **Data freshness metrics**: Monitor how current cached data is

This event handler is a critical component of the event-driven architecture, enabling efficient cross-service data synchronization while maintaining loose coupling and high performance.

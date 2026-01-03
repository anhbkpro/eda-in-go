package es

import (
	"context"
	"fmt"

	"eda-in-golang/internal/ddd"
)

type EventPublisher struct {
	AggregateStore
	publisher ddd.EventPublisher[ddd.AggregateEvent]
}

var _ AggregateStore = (*EventPublisher)(nil)

func NewEventPublisher(publisher ddd.EventPublisher[ddd.AggregateEvent]) AggregateStoreMiddleware {
	eventPublisher := EventPublisher{
		publisher: publisher,
	}

	return func(store AggregateStore) AggregateStore {
		eventPublisher.AggregateStore = store
		return eventPublisher
	}
}

func (p EventPublisher) Save(ctx context.Context, aggregate EventSourcedAggregate) error {
	fmt.Println("[Step 15] EventPublisher → EventStore.Save: persisting events to database")
	if err := p.AggregateStore.Save(ctx, aggregate); err != nil {
		return err
	}
	fmt.Println("[Step 16] EventStore: INSERT INTO events completed")

	fmt.Println("[Step 17] EventPublisher → Dispatcher.Publish: publishing events to handlers")
	err := p.publisher.Publish(ctx, aggregate.Events()...)
	fmt.Println("[Step 23] Dispatcher → EventPublisher: all handlers completed")
	return err
}

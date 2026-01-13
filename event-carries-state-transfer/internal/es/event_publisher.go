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
	events := aggregate.Events()
	fmt.Printf("[Step 15.1] EventPublisher → EventStore.Save: saving %d events\n", len(events))
	if err := p.AggregateStore.Save(ctx, aggregate); err != nil {
		fmt.Printf("[Step 15.2] EventPublisher → EventStore.Save: ERROR saving events: %v\n", err)
		return err
	}
	fmt.Println("[Step 15.3] EventStore: INSERT INTO events completed successfully")

	fmt.Println("[Step 16] EventPublisher → Dispatcher.Publish: publishing events to handlers")
	fmt.Printf("[Step 16.1] EventPublisher → Dispatcher.Publish: publishing %d events\n", len(events))
	for i, event := range events {
		fmt.Printf("[Step 16.2] EventPublisher → Dispatcher.Publish: publishing event %d/%d (%s)\n", i+1, len(events), event.EventName())
	}
	err := p.publisher.Publish(ctx, events...)
	if err != nil {
		fmt.Printf("[Step 16.3] EventPublisher → Dispatcher.Publish: ERROR publishing events: %v\n", err)
		return err
	}
	fmt.Println("[Step 16.4] Dispatcher → EventPublisher: events published successfully")
	fmt.Println("[Step 17] Dispatcher → EventPublisher: all handlers completed")
	return err
}

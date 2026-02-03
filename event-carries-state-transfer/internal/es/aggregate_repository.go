package es

import (
	"context"
	"fmt"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
)

// Load() and Save() are the only methods we will use with event-sourced aggregates and their event streams.
type AggregateRepository[T EventSourcedAggregate] struct {
	aggregateName string
	registry      registry.Registry
	store         AggregateStore
}

func NewAggregateRepository[T EventSourcedAggregate](aggregateName string, registry registry.Registry, store AggregateStore) AggregateRepository[T] {
	return AggregateRepository[T]{
		aggregateName: aggregateName,
		registry:      registry,
		store:         store,
	}
}

func (r AggregateRepository[T]) Load(ctx context.Context, aggregateID string) (agg T, err error) {
	fmt.Println("[Step 6] Repo → Store.Load: building aggregate from registry")
	var v any
	// Build the aggregate from the registry
	v, err = r.registry.Build(
		r.aggregateName,
		ddd.WithID(aggregateID),
		ddd.WithName(r.aggregateName),
	)
	if err != nil {
		return agg, err
	}

	var ok bool
	if agg, ok = v.(T); !ok {
		return agg, fmt.Errorf("%T is not the expected type %T", v, agg)
	}

	fmt.Println("[Step 7] Store → EventStore.Load: loading events from database")
	// Pass the new instance of the aggregate into the store.Load() method so it can receive deserialized data
	if err = r.store.Load(ctx, agg); err != nil {
		return agg, err
	}
	fmt.Println("[Step 8] EventStore → Aggregate: replaying events to rebuild state")

	// Return the aggregate if everything was successful
	return agg, nil
}

func (r AggregateRepository[T]) Save(ctx context.Context, aggregate T) error {
	fmt.Println("[Step 9] Repo.Save: checking for pending events")
	if aggregate.Version() == aggregate.PendingVersion() {
		fmt.Println("[Step 10] Repo.Save: no pending events, skipping save")
		return nil
	}

	fmt.Println("[Step 11] Repo.Save: found pending events, starting save process")

	fmt.Println("[Step 12] Repo → Aggregate: applying pending events to update internal state")
	events := aggregate.Events()
	fmt.Printf("[Step 12.1] Repo → Aggregate: processing %d pending events\n", len(events))
	for i, event := range events {
		fmt.Printf("[Step 12.2] Repo → Aggregate: applying event %d/%d (%s)\n", i+1, len(events), event.EventName())
		// Apply any new events the aggregate has created onto itself.
		if err := aggregate.ApplyEvent(event); err != nil {
			fmt.Printf("[Step 12.3] Repo → Aggregate: ERROR applying event %d: %v\n", i+1, err)
			return err
		}
	}
	fmt.Println("[Step 12.4] Repo → Aggregate: all events applied successfully")

	fmt.Println("[Step 13] Repo → Store.Save: delegating to aggregate store (with middleware)")
	// Pass the updated aggregate into the store.Save() method so that it can be serialized into the database.
	err := r.store.Save(ctx, aggregate)
	if err != nil {
		fmt.Printf("[Step 13.1] Repo → Store.Save: ERROR from store: %v\n", err)
		return err
	}
	fmt.Println("[Step 13.2] Repo → Store.Save: store save completed successfully")

	fmt.Println("[Step 14] Repo → Aggregate.CommitEvents: clearing pending events")
	// Update the aggregate version and clear the recently applied events using the aggregate CommitEvents() method.
	aggregate.CommitEvents()
	fmt.Println("[Step 14.1] Repo → Aggregate.CommitEvents: events committed")

	// Return nil if everything was successful
	return nil
}

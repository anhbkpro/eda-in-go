package es

import (
	"context"
	"fmt"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
)

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
	if err = r.store.Load(ctx, agg); err != nil {
		return agg, err
	}
	fmt.Println("[Step 8] EventStore → Aggregate: replaying events to rebuild state")

	return agg, nil
}

func (r AggregateRepository[T]) Save(ctx context.Context, aggregate T) error {
	if aggregate.Version() == aggregate.PendingVersion() {
		fmt.Println("[Step 12] Repo.Save: no pending events, skipping save")
		return nil
	}

	fmt.Println("[Step 13] Repo → Aggregate.ApplyEvent: applying pending events to update state")
	for _, event := range aggregate.Events() {
		if err := aggregate.ApplyEvent(event); err != nil {
			return err
		}
	}

	fmt.Println("[Step 14] Repo → Store.Save: delegating to aggregate store (with middleware)")
	err := r.store.Save(ctx, aggregate)
	if err != nil {
		return err
	}

	fmt.Println("[Step 24] Repo → Aggregate.CommitEvents: clearing pending events")
	aggregate.CommitEvents()

	return nil
}

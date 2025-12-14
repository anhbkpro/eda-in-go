package ddd

import (
	"context"
	"sync"
)

type (
	EventHandler[T Event] interface {
		HandleEvent(ctx context.Context, event T) error
	}

	EventHandlerFunc[T Event] func(ctx context.Context, event T) error

	EventSubscriber[T Event] interface {
		Subscribe(handler EventHandler[T], events ...string)
	}

	EventPublisher[T Event] interface {
		Publish(ctx context.Context, events ...T) error
	}

	EventDispatcher[T Event] struct {
		handlers []eventHandler[T]
		mu       sync.Mutex
	}

	eventHandler[T Event] struct {
		h       EventHandler[T]
		filters map[string]struct{}
	}
)

var _ interface {
	EventSubscriber[Event]
	EventPublisher[Event]
} = (*EventDispatcher[Event])(nil)

func NewEventDispatcher[T Event]() *EventDispatcher[T] {
	return &EventDispatcher[T]{
		handlers: make([]eventHandler[T], 0),
	}
}

func (d *EventDispatcher[T]) Subscribe(handler EventHandler[T], events ...string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var filters map[string]struct{}
	if len(events) > 0 {
		filters = make(map[string]struct{})
		for _, event := range events {
			filters[event] = struct{}{}
		}
	}

	d.handlers = append(d.handlers, eventHandler[T]{
		h:       handler,
		filters: filters,
	})
}

func (d *EventDispatcher[T]) Publish(ctx context.Context, events ...T) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, event := range events {
		for _, handler := range d.handlers {
			if len(handler.filters) > 0 {
				if _, ok := handler.filters[event.EventName()]; !ok {
					continue
				}
			}

			if err := handler.h.HandleEvent(ctx, event); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f EventHandlerFunc[T]) HandleEvent(ctx context.Context, event T) error {
	return f(ctx, event)
}

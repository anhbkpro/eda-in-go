package am

import (
	"context"
	"eda-in-golang/internal/ddd"
)

type (
	// Base interface for all messages with ID, name, and lifecycle methods (Ack/NAck/Extend/Kill)
	Message interface {
		ddd.IDer
		MessageName() string
		Ack() error
		NAck() error
		Extend() error
		Kill() error
	}

	MessageHandler[O Message] interface {
		HandleMessage(ctx context.Context, msg O) error
	}

	MessageHandlerFunc[O Message] func(ctx context.Context, msg O) error

	// Base interface for message publishers that can publish messages of type I
	MessagePublisher[I any] interface {
		Publish(ctx context.Context, topicName string, v I) error
	}

	// Base interface for message subscribers that can subscribe to messages of type O
	MessageSubscriber[O Message] interface {
		Subscribe(topicName string, handler MessageHandler[O], options ...SubscriberOption) error
	}

	// Base interface for message streams that combines publishing and subscribing capabilities
	MessageStream[I any, O Message] interface {
		MessagePublisher[I]
		MessageSubscriber[O]
	}
)

func (f MessageHandlerFunc[O]) HandleMessage(ctx context.Context, msg O) error {
	return f(ctx, msg)
}

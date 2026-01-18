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
	}

	IncomingMessage interface {
		Message
		Ack() error
		NAck() error
		Extend() error
		Kill() error
	}

	MessageHandler[I IncomingMessage] interface {
		HandleMessage(ctx context.Context, msg I) error
	}

	MessageHandlerFunc[I IncomingMessage] func(ctx context.Context, msg I) error

	// Base interface for message publishers that can publish messages of type I
	MessagePublisher[O any] interface {
		Publish(ctx context.Context, topicName string, v O) error
	}

	// Base interface for message subscribers that can subscribe to messages of type O
	MessageSubscriber[I IncomingMessage] interface {
		Subscribe(topicName string, handler MessageHandler[I], options ...SubscriberOption) error
	}

	// Base interface for message streams that combines publishing and subscribing capabilities
	MessageStream[O any, I IncomingMessage] interface {
		MessagePublisher[O]
		MessageSubscriber[I]
	}
)

func (f MessageHandlerFunc[I]) HandleMessage(ctx context.Context, msg I) error {
	return f(ctx, msg)
}

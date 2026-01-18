package sec

import (
	"context"

	"eda-in-golang/internal/am"
	"eda-in-golang/internal/ddd"
)

type (
	StepActionFunc[T any]       func(ctx context.Context, data T) am.Command
	StepReplyHandlerFunc[T any] func(ctx context.Context, data T, reply ddd.Reply) error

	SagaStep[T any] interface {
		Action(fn StepActionFunc[T]) SagaStep[T]
		Compensation(fn StepActionFunc[T]) SagaStep[T]
		OnActionReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T]
		OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T]
		isInvokable(compensating bool) bool
		execute(ctx context.Context, sagaCtx *SagaKontext[T]) stepResult[T]
		handle(ctx context.Context, sagaCtx *SagaKontext[T], reply ddd.Reply) error
	}

	sagaStep[T any] struct {
		actions  map[bool]StepActionFunc[T]
		handlers map[bool]map[string]StepReplyHandlerFunc[T]
	}

	stepResult[T any] struct {
		ctx *SagaKontext[T]
		cmd am.Command
		err error
	}
)

var _ SagaStep[any] = (*sagaStep[any])(nil)

func (s *sagaStep[T]) Action(fn StepActionFunc[T]) SagaStep[T] {
	s.actions[notCompensating] = fn
	return s
}

func (s *sagaStep[T]) Compensation(fn StepActionFunc[T]) SagaStep[T] {
	s.actions[isCompensating] = fn
	return s
}

func (s *sagaStep[T]) OnActionReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T] {
	s.handlers[notCompensating][replyName] = fn
	return s
}

func (s *sagaStep[T]) OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T] {
	s.handlers[isCompensating][replyName] = fn
	return s
}

func (s sagaStep[T]) isInvokable(compensating bool) bool {
	return s.actions[compensating] != nil
}

func (s sagaStep[T]) execute(ctx context.Context, sagaKtx *SagaKontext[T]) stepResult[T] {
	if action := s.actions[sagaKtx.Compensating]; action != nil {
		return stepResult[T]{
			ctx: sagaKtx,
			cmd: action(ctx, sagaKtx.Data),
		}
	}

	return stepResult[T]{ctx: sagaKtx}
}

func (s sagaStep[T]) handle(ctx context.Context, sagaKtx *SagaKontext[T], reply ddd.Reply) error {
	if handler := s.handlers[sagaKtx.Compensating][reply.ReplyName()]; handler != nil {
		return handler(ctx, sagaKtx.Data, reply)
	}
	return nil
}

type StepOption[T any] func(step *sagaStep[T])

// WithAction sets the action function for the saga step.
func WithAction[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.actions[notCompensating] = fn
	}
}

// WithCompensation sets the compensation function for the saga step.
func WithCompensation[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.actions[isCompensating] = fn
	}
}

// OnActionReply registers a reply handler for action replies.
func OnActionReply[T any](replyName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.handlers[notCompensating][replyName] = fn
	}
}

// OnCompensationReply registers a reply handler for compensation replies.
func OnCompensationReply[T any](replyName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.handlers[isCompensating][replyName] = fn
	}
}

package sec

import "eda-in-golang/internal/am"

const (
	SagaCommandIDHandler   = am.CommandHandlerPrefix + "SAGA_ID"
	SagaCommandNameHandler = am.CommandHandlerPrefix + "SAGA_NAME"

	SagaReplyIDHandler   = am.ReplyHandlerPrefix + "SAGA_ID"
	SagaReplyNameHandler = am.ReplyHandlerPrefix + "SAGA_NAME"
)

type (
	SagaKontext[T any] struct {
		ID           string
		Data         T
		Step         int
		Done         bool
		Compensating bool
	}

	Saga[T any] interface {
		AddStep() SagaStep[T]
		Name() string
		ReplyTopic() string
		getSteps() []SagaStep[T]
	}

	saga[T any] struct {
		name       string
		replyTopic string
		steps      []SagaStep[T]
	}
)

const (
	notCompensating = false
	isCompensating  = true
)

func NewSaga[T any](name, replyTopic string) Saga[T] {
	return &saga[T]{
		name:       name,
		replyTopic: replyTopic,
		steps:      make([]SagaStep[T], 0),
	}
}

func (s *saga[T]) AddStep() SagaStep[T] {
	step := &sagaStep[T]{
		actions: map[bool]StepActionFunc[T]{
			notCompensating: nil,
			isCompensating:  nil,
		},
		handlers: map[bool]map[string]StepReplyHandlerFunc[T]{
			notCompensating: make(map[string]StepReplyHandlerFunc[T]),
			isCompensating:  make(map[string]StepReplyHandlerFunc[T]),
		},
	}
	s.steps = append(s.steps, step)
	return step
}

func (s *saga[T]) Name() string {
	return s.name
}

func (s *saga[T]) ReplyTopic() string {
	return s.replyTopic
}

func (s *saga[T]) getSteps() []SagaStep[T] {
	return s.steps
}

func (s *SagaKontext[T]) advance(steps int) {
	var dir = 1
	if s.Compensating {
		dir = -1
	}
	s.Step += dir * steps
}

func (s *SagaKontext[T]) complete() {
	s.Done = true
}

func (s *SagaKontext[T]) compensate() {
	s.Compensating = true
}

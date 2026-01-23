package tm

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"eda-in-golang/internal/am"
)

const messageLimit = 10
const pollingIntervalMilliseconds = 500 * time.Millisecond

type OutboxProcessor interface {
	Start(ctx context.Context) error
}

type outboxProcessor struct {
	publisher am.RawMessageStream
	store     OutboxStore
	logger    zerolog.Logger
}

func NewOutboxProcessor(publisher am.RawMessageStream, store OutboxStore, logger zerolog.Logger) OutboxProcessor {
	return &outboxProcessor{
		publisher: publisher,
		store:     store,
		logger:    logger,
	}
}

func (p outboxProcessor) Start(ctx context.Context) error {
	errC := make(chan error)

	go func() {
		errC <- p.processMessages(ctx)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errC:
		return err
	}
}

func (p outboxProcessor) processMessages(ctx context.Context) error {
	ticker := time.NewTicker(pollingIntervalMilliseconds)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			messages, err := p.store.FindUnpublished(ctx, messageLimit)
			if err != nil {
				return err
			}

			if len(messages) > 0 {
				p.logger.Info().Int("count", len(messages)).Msg("found unpublished messages to process")
			}

			var ids []string
			for _, msg := range messages {
				err := p.publisher.Publish(ctx, msg.Subject(), msg)
				if err != nil {
					return err
				}
				ids = append(ids, msg.ID())
			}

			if len(ids) > 0 {
				err = p.store.MarkAsPublished(ctx, ids...)
				if err != nil {
					return err
				}
			}
		}
	}
}

package logging

import (
	"context"

	"github.com/rs/zerolog"

	"eda-in-golang/notifications/internal/application"
)

type ApplicationLogger struct {
	application.App
	logger zerolog.Logger
}

var _ application.App = (*ApplicationLogger)(nil)

func LogApplicationAccess(app application.App, logger zerolog.Logger) ApplicationLogger {
	return ApplicationLogger{
		App:    app,
		logger: logger,
	}
}

func (a ApplicationLogger) NotifyOrderCreated(ctx context.Context, notify application.OrderCreated) (err error) {
	a.logger.Info().Msg("--> Notifications.NotifyOrderCreated")
	defer func() { a.logger.Info().Err(err).Msg("<-- Notifications.NotifyOrderCreated") }()
	return a.App.NotifyOrderCreated(ctx, notify)
}

func (a ApplicationLogger) NotifyOrderCanceled(ctx context.Context, notify application.OrderCanceled) (err error) {
	a.logger.Info().Msg("--> Notifications.NotifyOrderCanceled")
	defer func() { a.logger.Info().Err(err).Msg("<-- Notifications.NotifyOrderCanceled") }()
	return a.App.NotifyOrderCanceled(ctx, notify)
}

func (a ApplicationLogger) NotifyOrderReady(ctx context.Context, notify application.OrderReady) (err error) {
	a.logger.Info().Msg("--> Notifications.NotifyOrderReady")
	defer func() { a.logger.Info().Err(err).Msg("<-- Notifications.NotifyOrderReady") }()
	return a.App.NotifyOrderReady(ctx, notify)
}

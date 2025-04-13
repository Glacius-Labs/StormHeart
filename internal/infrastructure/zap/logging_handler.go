package handler

import (
	"context"

	"go.uber.org/zap"

	"github.com/glacius-labs/StormHeart/internal/core/event"
)

type LoggingHandler struct {
	logger *zap.Logger
}

func NewLoggingHandler(logger *zap.Logger) *LoggingHandler {
	if logger == nil {
		panic("logger cannot be nil")
	}

	return &LoggingHandler{
		logger: logger,
	}
}

func (h *LoggingHandler) Handle(ctx context.Context, event event.Event) error {
	if event.Error() != nil {
		h.logger.Error("Event received",
			zap.String("type", string(event.Type())),
			zap.String("message", event.Message()),
			zap.Time("timestamp", event.Timestamp()),
			zap.Error(event.Error()),
		)
	} else {
		h.logger.Info("Event received",
			zap.String("type", string(event.Type())),
			zap.String("message", event.Message()),
			zap.Time("timestamp", event.Timestamp()),
		)
	}

	return nil
}

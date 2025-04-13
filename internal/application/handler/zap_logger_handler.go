package handler

import (
	"context"

	"go.uber.org/zap"

	"github.com/glacius-labs/StormHeart/internal/core/event"
)

type ZapLoggerHandler struct {
	logger *zap.Logger
}

func NewZapLoggerHandler(logger *zap.Logger) *ZapLoggerHandler {
	if logger == nil {
		panic("logger cannot be nil")
	}

	return &ZapLoggerHandler{
		logger: logger,
	}
}

func (h *ZapLoggerHandler) Handle(ctx context.Context, event event.Event) {
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
}

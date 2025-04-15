package mqtt

import (
	"context"
	"encoding/json"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"go.uber.org/zap"
)

type EventPublisherHandler struct {
	client Client
	topic  string
	logger *zap.Logger
}

func NewEventPublisherHandler(client Client, topic string, logger *zap.Logger) *EventPublisherHandler {
	if client == nil {
		panic("client cannot be nil")
	}
	if topic == "" {
		panic("topic cannot be empty")
	}
	if logger == nil {
		panic("logger cannot be nil")
	}

	return &EventPublisherHandler{
		client: client,
		topic:  topic,
		logger: logger,
	}
}

func (h *EventPublisherHandler) Handle(ctx context.Context, e event.EventType) {
	h.logger.Info(
		"handling event",
		zap.String("event_type", string(e.Type())),
		zap.String("event_message", e.Message()),
	)

	payload := map[string]string{
		"type":      string(e.Type()),
		"message":   e.Message(),
		"timestamp": e.Timestamp().Format("2006-01-02T15:04:05Z07:00"),
		"error":     "",
	}

	if e.Error() != nil {
		payload["error"] = e.Error().Error()
	}

	data, err := json.Marshal(payload)
	if err != nil {
		h.logger.Error(
			"failed to marshal event",
			zap.String("event_type", string(e.Type())),
			zap.String("event_message", e.Message()),
			zap.Error(err),
		)
		return
	}

	err = h.client.Publish(ctx, h.topic, data)
	if err != nil {
		h.logger.Error(
			"failed to publish event",
			zap.String("event_type", string(e.Type())),
			zap.String("event_message", e.Message()),
			zap.Error(err),
		)
		return
	}

	h.logger.Info(
		"event handled successfully",
		zap.String("event_type", string(e.Type())),
		zap.String("event_message", e.Message()),
		zap.String("topic", h.topic),
	)
}

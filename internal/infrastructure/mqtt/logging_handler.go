package mqtt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/core/event"
)

type LoggingHandler struct {
	client Client
	topic  string
}

func NewLoggingHandler(client Client, topic string) *LoggingHandler {
	if client == nil {
		panic("client cannot be nil")
	}
	if topic == "" {
		panic("topic cannot be empty")
	}

	return &LoggingHandler{
		client: client,
		topic:  topic,
	}
}

func (h *LoggingHandler) Handle(ctx context.Context, e event.Event) error {
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
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = h.client.Publish(ctx, h.topic, data)
	if err != nil {
		return fmt.Errorf("failed to publish event to MQTT: %w", err)
	}

	return nil
}

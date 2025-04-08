package mqtt

import (
	"context"
	"encoding/json"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"go.uber.org/zap"
)

type MQTTWatcher struct {
	client     Client
	topic      string
	sourceName string
	pushFunc   watcher.PushFunc
	logger     *zap.Logger
}

func NewWatcher(client Client, topic, sourceName string, pushFunc watcher.PushFunc, logger *zap.Logger) *MQTTWatcher {
	return &MQTTWatcher{
		client:     client,
		topic:      topic,
		sourceName: sourceName,
		pushFunc:   pushFunc,
		logger:     logger,
	}
}

func (w *MQTTWatcher) Watch(ctx context.Context) error {
	if err := w.client.Connect(); err != nil {
		return err
	}

	w.logger.Info("Connected to MQTT broker")

	if err := w.client.Subscribe(ctx, w.topic, w.handleMessage); err != nil {
		return err
	}

	w.logger.Info("Subscribed to MQTT topic", zap.String("topic", w.topic))

	<-ctx.Done()

	w.logger.Info("Context canceled, disconnecting MQTT client")
	w.client.Disconnect()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	w.pushFunc(shutdownCtx, w.sourceName, []model.Deployment{})

	return nil
}

func (w *MQTTWatcher) handleMessage(ctx context.Context, payload []byte) {
	w.logger.Info("MQTT message received", zap.String("topic", w.topic))

	var deployments []model.Deployment
	if err := json.Unmarshal(payload, &deployments); err != nil {
		w.logger.Error("Invalid deployment message", zap.Error(err))
		return
	}

	w.pushFunc(ctx, w.sourceName, deployments)
}

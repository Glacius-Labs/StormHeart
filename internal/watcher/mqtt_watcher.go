package watcher

import (
	"context"
	"encoding/json"

	"github.com/glacius-labs/StormHeart/internal/model"
	"github.com/glacius-labs/StormHeart/internal/mqtt"
	"go.uber.org/zap"
)

type MQTTWatcher struct {
	client     mqtt.Client
	topic      string
	sourceName string
	pushFunc   PushFunc
	logger     *zap.Logger
}

func NewMQTTWatcher(client mqtt.Client, topic, sourceName string, pushFunc PushFunc, logger *zap.Logger) *MQTTWatcher {
	return &MQTTWatcher{
		client:     client,
		topic:      topic,
		sourceName: sourceName,
		pushFunc:   pushFunc,
		logger:     logger,
	}
}

func (w *MQTTWatcher) Start(ctx context.Context) error {
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

	w.pushFunc(ctx, w.sourceName, []model.Deployment{})

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

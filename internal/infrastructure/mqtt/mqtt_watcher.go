package mqtt

import (
	"context"
	"encoding/json"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"go.uber.org/zap"
)

type MQTTWatcher struct {
	client      Client
	topic       string
	sourceName  string
	handlerFunc watcher.HandlerFunc
	logger      *zap.Logger
}

func NewWatcher(client Client, topic, sourceName string, handlerFunc watcher.HandlerFunc, logger *zap.Logger) *MQTTWatcher {
	if client == nil {
		panic("MQTTWatcher requires a non-nil client")
	}

	if topic == "" {
		panic("MQTTWatcher requires a non-empty topic")
	}

	if sourceName == "" {
		panic("MQTTWatcher requires a non-empty source name")
	}

	if handlerFunc == nil {
		panic("MQTTWatcher requires a non-nil handler func")
	}

	if logger == nil {
		panic("MQTTWatcher requires a non-nil logger")
	}

	return &MQTTWatcher{
		client:      client,
		topic:       topic,
		sourceName:  sourceName,
		handlerFunc: handlerFunc,
		logger:      logger,
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

	w.logger.Info("Initiating shutdown")

	w.client.Disconnect()
	watcher.PushEmptyDeployments(w.handlerFunc, w.sourceName)

	w.logger.Info("Shutdown complete")

	return nil
}

func (w *MQTTWatcher) handleMessage(ctx context.Context, payload []byte) {
	w.logger.Info("MQTT message received", zap.String("topic", w.topic))

	var deployments []model.Deployment
	if err := json.Unmarshal(payload, &deployments); err != nil {
		w.logger.Error("Invalid deployment message", zap.Error(err))
		return
	}

	w.handlerFunc(ctx, w.sourceName, deployments)
}

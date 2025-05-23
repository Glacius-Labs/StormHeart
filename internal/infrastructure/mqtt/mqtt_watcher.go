package mqtt

import (
	"context"
	"encoding/json"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
)

type MQTTWatcher struct {
	client     Client
	topic      string
	sourceName string
	dispatcher *event.Dispatcher
}

func NewWatcher(client Client, topic, sourceName string, dispatcher *event.Dispatcher) *MQTTWatcher {
	if client == nil {
		panic("client cannot be nil")
	}

	if topic == "" {
		panic("topic cannot be empty")
	}

	if sourceName == "" {
		panic("source name cannot be empty")
	}

	if dispatcher == nil {
		panic("dispatcher cannot be nil")
	}

	return &MQTTWatcher{
		client:     client,
		topic:      topic,
		sourceName: sourceName,
		dispatcher: dispatcher,
	}
}

func (w *MQTTWatcher) Watch(ctx context.Context) {
	startedEvent := watcher.NewWatcherStartedEvent(w.sourceName)
	w.dispatcher.Dispatch(ctx, startedEvent)

	if err := w.client.Connect(); err != nil {
		e := watcher.NewWatcherStoppedEvent(w.sourceName, err)
		w.dispatcher.Dispatch(ctx, e)
		return
	}

	if err := w.client.Subscribe(ctx, w.topic, w.handleMessage); err != nil {
		e := watcher.NewWatcherStoppedEvent(w.sourceName, err)
		w.dispatcher.Dispatch(ctx, e)
		return
	}

	<-ctx.Done()

	stoppedEvent := watcher.NewWatcherStoppedEvent(w.sourceName, nil)
	w.dispatcher.Dispatch(ctx, stoppedEvent)
}

func (w *MQTTWatcher) handleMessage(ctx context.Context, payload []byte) {
	var deployments []model.Deployment
	err := json.Unmarshal(payload, &deployments)

	e := watcher.NewDeploymentsReceivedEvent(w.sourceName, deployments, err)
	w.dispatcher.Dispatch(ctx, e)
}

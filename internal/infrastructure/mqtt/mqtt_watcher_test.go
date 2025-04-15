package mqtt_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/mock"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/mqtt"
	"github.com/stretchr/testify/require"
)

func TestNewMQTTWatcher_PanicsOnNilClient(t *testing.T) {
	dispatcher := event.NewDispatcher()

	require.Panics(t, func() {
		_ = mqtt.NewWatcher(nil, "topic", "source", dispatcher)
	}, "expected panic on nil client")
}

func TestNewMQTTWatcher_PanicsOnEmptyTopic(t *testing.T) {
	mockClient := mock.NewMockClient()
	dispatcher := event.NewDispatcher()

	require.Panics(t, func() {
		_ = mqtt.NewWatcher(mockClient, "", "source", dispatcher)
	}, "expected panic on empty topic")
}

func TestNewMQTTWatcher_PanicsOnEmptySourceName(t *testing.T) {
	mockClient := mock.NewMockClient()
	dispatcher := event.NewDispatcher()

	require.Panics(t, func() {
		_ = mqtt.NewWatcher(mockClient, "topic", "", dispatcher)
	}, "expected panic on empty source name")
}

func TestNewMQTTWatcher_PanicsOnNilDispatcher(t *testing.T) {
	mockClient := mock.NewMockClient()

	require.Panics(t, func() {
		_ = mqtt.NewWatcher(mockClient, "topic", "source", nil)
	}, "expected panic on nil dispatcher")
}

func TestMQTTWatcher_StartsAndWaitsOnContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	mockClient := mock.NewMockClient()

	w := mqtt.NewWatcher(mockClient, "test/topic", "mqtt-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	cancel()
	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var started, stopped bool
	for _, e := range events {
		switch e.Type() {
		case watcher.EventTypeWatcherStarted:
			started = true
		case watcher.EventTypeWatcherStopped:
			stopped = true
			require.Nil(t, e.Error(), "expected clean shutdown with no error")
		}
	}

	require.True(t, started, "expected WatcherStartedEvent")
	require.True(t, stopped, "expected WatcherStoppedEvent")
}

func TestMQTTWatcher_FailsOnConnect(t *testing.T) {
	ctx := t.Context()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	mockClient := mock.NewMockClient()
	mockClient.ShouldFailConnect = true

	w := mqtt.NewWatcher(mockClient, "test/topic", "mqtt-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var started, stoppedWithError bool
	for _, e := range events {
		switch e.Type() {
		case watcher.EventTypeWatcherStarted:
			started = true
		case watcher.EventTypeWatcherStopped:
			stoppedWithError = true
			require.NotNil(t, e.Error(), "expected error on WatcherStoppedEvent after connect failure")
		}
	}

	require.True(t, started, "expected WatcherStartedEvent")
	require.True(t, stoppedWithError, "expected WatcherStoppedEvent with error after connect failure")
}

func TestMQTTWatcher_FailsOnSubscribe(t *testing.T) {
	ctx := t.Context()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	mockClient := mock.NewMockClient()
	mockClient.ShouldFailSubscribe = true

	w := mqtt.NewWatcher(mockClient, "test/topic", "mqtt-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var started, stoppedWithError bool
	for _, e := range events {
		switch e.Type() {
		case watcher.EventTypeWatcherStarted:
			started = true
		case watcher.EventTypeWatcherStopped:
			stoppedWithError = true
			require.NotNil(t, e.Error(), "expected error on WatcherStoppedEvent after subscribe failure")
		}
	}

	require.True(t, started, "expected WatcherStartedEvent")
	require.True(t, stoppedWithError, "expected WatcherStoppedEvent with error after subscribe failure")
}

func TestMQTTWatcher_ReceivesAndHandlesMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	mockClient := mock.NewMockClient()

	w := mqtt.NewWatcher(mockClient, "test/topic", "mqtt-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	payload, err := json.Marshal([]model.Deployment{
		{Name: "dummy", Image: "alpine"},
	})
	require.NoError(t, err)

	require.NotNil(t, mockClient.ReceivedHandler, "expected ReceivedHandler to be set")

	mockClient.ReceivedHandler(ctx, payload)

	time.Sleep(100 * time.Millisecond)

	cancel()
	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var deploymentsReceived bool
	for _, e := range events {
		if e.Type() == watcher.EventTypeDeploymentsReceived {
			deploymentsReceived = true

			receivedEvent, ok := e.(watcher.DeploymentsReceivedEvent)
			require.True(t, ok, "expected event to be DeploymentsReceivedEvent type")
			require.Nil(t, receivedEvent.Error(), "expected DeploymentsReceivedEvent to have no error")
			require.Len(t, receivedEvent.Deployments, 1, "expected exactly one deployment in received event")
			require.Equal(t, "dummy", receivedEvent.Deployments[0].Source, "expected deployment name to match")
		}
	}

	require.True(t, deploymentsReceived, "expected a DeploymentsReceivedEvent after message handling")
}

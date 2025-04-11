package static_test

import (
	"context"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/mock"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/static"
	"github.com/stretchr/testify/require"
)

func TestNewWatcher_NotPanicsOnNilDeployments(t *testing.T) {
	dispatcher := event.NewDispatcher()

	require.NotPanics(t, func() {
		_ = static.NewWatcher(nil, dispatcher)
	}, "expected no panic when deployments is nil")
}

func TestNewWatcher_PanicsOnNilDispatcher(t *testing.T) {
	require.Panics(t, func() {
		_ = static.NewWatcher(nil, nil)
	}, "expected panic when dispatcher is nil")
}

func TestStaticWatcher_Watch_EmitsEvents(t *testing.T) {
	// Prepare context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Prepare dispatcher + mock handler
	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	// Prepare deployments
	deployments := []model.Deployment{
		{Name: "test-service", Image: "nginx"},
	}

	// Create watcher
	watcher := static.NewWatcher(deployments, dispatcher)

	// Start watcher in background
	go func() {
		err := watcher.Watch(ctx)
		require.NoError(t, err, "watcher should not return error")
	}()

	// Allow some time for events to be emitted
	time.Sleep(50 * time.Millisecond)

	// Trigger shutdown
	cancel()

	// Allow time for shutdown event
	time.Sleep(50 * time.Millisecond)

	events := handler.Events()

	require.GreaterOrEqual(t, len(events), 3, "expected at least 3 events: started, received, stopped")

	var started, received, stopped bool
	for _, e := range events {
		switch e.Type() {
		case "watcher_started":
			started = true
		case "deployments_received":
			received = true
		case "watcher_stopped":
			stopped = true
		}
	}

	require.True(t, started, "expected WatcherStartedEvent")
	require.True(t, received, "expected DeploymentsReceivedEvent")
	require.True(t, stopped, "expected WatcherStoppedEvent")
}

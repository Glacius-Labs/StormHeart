package file_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/file"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFileWatcher_PanicsOnEmptyPath(t *testing.T) {
	dispatcher := event.NewDispatcher()

	require.Panics(t, func() {
		_ = file.NewWatcher("", "source", dispatcher)
	}, "expected panic on empty path")
}

func TestNewFileWatcher_PanicsOnEmptySource(t *testing.T) {
	dispatcher := event.NewDispatcher()

	require.Panics(t, func() {
		_ = file.NewWatcher("somepath", "", dispatcher)
	}, "expected panic on empty source name")
}

func TestNewFileWatcher_PanicsOnNilDispatcher(t *testing.T) {
	require.Panics(t, func() {
		_ = file.NewWatcher("somepath", "source", nil)
	}, "expected panic on nil dispatcher")
}

func TestFileWatcher_Watch_BadPath_EmitsStoppedWithError(t *testing.T) {
	ctx := t.Context()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	w := file.NewWatcher("/path/that/does/not/exist.json", "bad-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	require.GreaterOrEqual(t, len(events), 2, "expected at least two events: started, stopped with error")

	var started, stoppedWithError bool
	for _, e := range events {
		switch e.Type() {
		case watcher.EventTypeWatcherStarted:
			started = true
		case watcher.EventTypeWatcherStopped:
			stoppedWithError = true
			require.NotNil(t, e.Error(), "expected watcher_stopped event to carry an error")
		}
		require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected recent timestamp for event")
	}

	require.True(t, started, "expected WatcherStartedEvent")
	require.True(t, stoppedWithError, "expected WatcherStoppedEvent with error")
}

func TestFileWatcher_Watch_BadJSON_EmitsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	tmpFile, err := os.CreateTemp("", "badfile-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`{invalid json}`)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	w := file.NewWatcher(tmpFile.Name(), "bad-json-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	cancel()

	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var started, deploymentsReceived, stopped bool
	for _, e := range events {
		switch e.Type() {
		case watcher.EventTypeWatcherStarted:
			started = true
		case watcher.EventTypeDeploymentsReceived:
			deploymentsReceived = true
			require.NotNil(t, e.Error(), "expected deployments_received event to carry unmarshal error")
		case watcher.EventTypeWatcherStopped:
			stopped = true
		}
		require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected recent timestamp for event")
	}

	require.True(t, started, "expected WatcherStartedEvent")
	require.True(t, deploymentsReceived, "expected DeploymentsReceivedEvent even if unmarshal failed")
	require.True(t, stopped, "expected WatcherStoppedEvent")
}

func TestFileWatcher_Watch_FileRemovedDuringWatching(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	tmpFile, err := os.CreateTemp("", "watchedfile-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`[{"name": "dummy", "image": "alpine"}]`)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	w := file.NewWatcher(tmpFile.Name(), "runtime-fail-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	err = os.Remove(tmpFile.Name())
	require.NoError(t, err, "expected temp file to be removed successfully")

	time.Sleep(100 * time.Millisecond)

	cancel()

	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var started, stoppedWithError bool
	for _, e := range events {
		switch e.Type() {
		case watcher.EventTypeWatcherStarted:
			started = true
		case watcher.EventTypeWatcherStopped:
			stoppedWithError = true
			require.NotNil(t, e.Error(), "expected error on watcher_stopped after file removal")
		}
		require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected recent timestamp for event")
	}

	require.True(t, started, "expected WatcherStartedEvent")
	require.True(t, stoppedWithError, "expected WatcherStoppedEvent with error after file disappeared")
}

func TestFileWatcher_Watch_FileRenamedDuringWatching(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	tmpFile, err := os.CreateTemp("", "renamefile-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`[{"name": "dummy", "image": "alpine"}]`)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	w := file.NewWatcher(tmpFile.Name(), "rename-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	newPath := tmpFile.Name() + "-renamed"
	err = os.Rename(tmpFile.Name(), newPath)
	require.NoError(t, err, "expected file to be renamed successfully")
	defer os.Remove(newPath)

	time.Sleep(250 * time.Millisecond)

	cancel()

	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var started, stoppedWithError bool
	for _, e := range events {
		switch e.Type() {
		case watcher.EventTypeWatcherStarted:
			started = true
		case watcher.EventTypeWatcherStopped:
			stoppedWithError = true
			require.NotNil(t, e.Error(), "expected watcher_stopped event with error after rename")
		}
		require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected recent timestamp for event")
	}

	require.True(t, started, "expected WatcherStartedEvent")
	require.True(t, stoppedWithError, "expected WatcherStoppedEvent with error after file rename")
}

func TestFileWatcher_Watch_ShutsDownCleanlyOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	tmpFile, err := os.CreateTemp("", "shutdownfile-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`[{"name": "dummy", "image": "alpine"}]`)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	w := file.NewWatcher(tmpFile.Name(), "shutdown-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	cancel()

	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var stopped bool
	for _, e := range events {
		if e.Type() == watcher.EventTypeWatcherStopped {
			stopped = true
			require.Nil(t, e.Error(), "expected clean shutdown with no error")
		}
		require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected recent timestamp for event")
	}

	require.True(t, stopped, "expected WatcherStoppedEvent on shutdown")
}

func TestFileWatcher_Watch_InitialLoadEmitsDeployments(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	tmpFile, err := os.CreateTemp("", "initialload-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`[{"name": "initial", "image": "alpine:latest"}]`)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	w := file.NewWatcher(tmpFile.Name(), "initial-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	cancel()

	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var started, deploymentsReceived, stopped bool
	for _, e := range events {
		switch e.Type() {
		case watcher.EventTypeWatcherStarted:
			started = true
		case watcher.EventTypeDeploymentsReceived:
			deploymentsReceived = true

			receivedEvent, ok := e.(watcher.DeploymentsReceivedEvent)
			require.True(t, ok, "expected event to be DeploymentsReceivedEvent type")
			require.Nil(t, receivedEvent.Error(), "expected DeploymentsReceivedEvent to have no error on initial load")
			require.Len(t, receivedEvent.Deployments, 1, "expected exactly one deployment on initial load")
			require.Equal(t, "initial", receivedEvent.Deployments[0].Name, "expected deployment name to match")
		case watcher.EventTypeWatcherStopped:
			stopped = true
		}
		require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected recent timestamp for event")
	}

	require.True(t, started, "expected WatcherStartedEvent")
	require.True(t, deploymentsReceived, "expected DeploymentsReceivedEvent from initial file load")
	require.True(t, stopped, "expected WatcherStoppedEvent on clean shutdown")
}

func TestFileWatcher_Watch_HandlesFileChangeAndEmitsDeployments(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	tmpFile, err := os.CreateTemp("", "changefile-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`[{"name": "dummy", "image": "alpine"}]`)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	w := file.NewWatcher(tmpFile.Name(), "change-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(500 * time.Millisecond)

	err = os.WriteFile(tmpFile.Name(), []byte(`[{"name": "updated", "image": "alpine:latest"}]`), 0644)
	require.NoError(t, err)

	time.Sleep(500 * time.Millisecond)

	cancel()

	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var received int
	for _, e := range events {
		if e.Type() == watcher.EventTypeDeploymentsReceived {
			received++

			receivedEvent, ok := e.(watcher.DeploymentsReceivedEvent)
			require.True(t, ok, "expected event to be DeploymentsReceivedEvent type")

			require.Nil(t, receivedEvent.Error(), "expected DeploymentsReceivedEvent to have no error")

			require.NotNil(t, receivedEvent.Deployments, "expected Deployments not to be nil")
			require.NotEmpty(t, receivedEvent.Deployments, "expected at least one deployment received")
		}
		require.WithinDuration(t, time.Now(), e.Timestamp(), 5*time.Second, "expected recent timestamp for event")
	}

	require.Equal(t, 2, received, "expected two deployments_received events (initial + file change)")
}

func TestFileWatcher_Watch_DebouncesRapidFileChanges(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := event.NewDispatcher()
	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	tmpFile, err := os.CreateTemp("", "debouncefile-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`[{"name": "first", "image": "alpine"}]`)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	w := file.NewWatcher(tmpFile.Name(), "debounce-watcher", dispatcher)

	go w.Watch(ctx)

	time.Sleep(100 * time.Millisecond)

	// Rapidly write multiple times
	for range 5 {
		err = os.WriteFile(tmpFile.Name(), []byte(`[{"name": "debounced", "image": "alpine:latest"}]`), 0644)
		require.NoError(t, err)
		time.Sleep(50 * time.Millisecond) // Faster than debounceDelay
	}

	time.Sleep(500 * time.Millisecond) // Wait enough for debounce to fire once

	cancel()
	time.Sleep(100 * time.Millisecond)

	events := handler.Events()

	var deploymentsReceived int
	for _, e := range events {
		if e.Type() == watcher.EventTypeDeploymentsReceived {
			deploymentsReceived++
		}
		require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected recent timestamp for event")
	}

	require.LessOrEqual(t, deploymentsReceived, 3, "expected only 1-3 deployments_received events due to debounce")
}

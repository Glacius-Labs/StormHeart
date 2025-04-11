package watcher_test

import (
	"errors"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"github.com/stretchr/testify/require"
)

func TestWatcherStoppedEvent_Success(t *testing.T) {
	event := watcher.NewWatcherStoppedEvent("file-watcher", nil)

	require.Equal(t, "Watcher stopped cleanly for source file-watcher", event.Message(), "expected clean stop message")
	require.Equal(t, watcher.EventTypeWatcherStopped, event.Type(), "expected correct event type")
	require.Nil(t, event.Error(), "expected no error for clean stop")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "timestamp should be recent")
}

func TestWatcherStoppedEvent_Failure(t *testing.T) {
	err := errors.New("connection lost")
	event := watcher.NewWatcherStoppedEvent("mqtt-watcher", err)

	require.Contains(t, event.Message(), "Watcher stopped with error for source mqtt-watcher", "expected error stop message prefix")
	require.Contains(t, event.Message(), "connection lost", "expected error detail in message")
	require.Equal(t, watcher.EventTypeWatcherStopped, event.Type(), "expected correct event type")
	require.Equal(t, err, event.Error(), "expected correct error reference")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "timestamp should be recent")
}

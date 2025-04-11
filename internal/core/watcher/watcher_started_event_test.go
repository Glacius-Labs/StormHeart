package watcher_test

import (
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"github.com/stretchr/testify/require"
)

func TestWatcherStartedEvent(t *testing.T) {
	event := watcher.NewWatcherStartedEvent("file-watcher")

	require.Equal(t, "Watcher started for source file-watcher", event.Message(), "expected correct start message")
	require.Equal(t, watcher.EventTypeWatcherStarted, event.Type(), "expected correct event type")
	require.Nil(t, event.Error(), "expected no error")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "timestamp should be recent")
}

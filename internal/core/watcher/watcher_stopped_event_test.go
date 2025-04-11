package watcher_test

import (
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"github.com/stretchr/testify/require"
)

func TestWatcherStoppedEvent(t *testing.T) {
	event := watcher.NewWatcherStoppedEvent("mqtt-watcher")

	require.Equal(t, "Watcher stopped for source mqtt-watcher", event.Message(), "expected correct stop message")
	require.Equal(t, watcher.EventTypeWatcherStopped, event.Type(), "expected correct event type")
	require.Nil(t, event.Error(), "expected no error")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "timestamp should be recent")
}

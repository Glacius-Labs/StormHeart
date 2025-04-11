package watcher

import (
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
)

const EventTypeWatcherStopped event.Type = "watcher_stopped"

type WatcherStoppedEvent struct {
	Source    string
	timestamp time.Time
}

func NewWatcherStoppedEvent(source string) WatcherStoppedEvent {
	return WatcherStoppedEvent{
		Source:    source,
		timestamp: time.Now(),
	}
}

func (e WatcherStoppedEvent) Message() string {
	return "Watcher stopped for source " + e.Source
}

func (e WatcherStoppedEvent) Type() event.Type {
	return EventTypeWatcherStopped
}

func (e WatcherStoppedEvent) Error() error {
	return nil
}

func (e WatcherStoppedEvent) Timestamp() time.Time {
	return e.timestamp
}

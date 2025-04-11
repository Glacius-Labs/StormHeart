package watcher

import (
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
)

const EventTypeWatcherStarted event.Type = "watcher_started"

type WatcherStartedEvent struct {
	Source    string
	timestamp time.Time
}

func NewWatcherStartedEvent(source string) WatcherStartedEvent {
	return WatcherStartedEvent{
		Source:    source,
		timestamp: time.Now(),
	}
}

func (e WatcherStartedEvent) Message() string {
	return "Watcher started for source " + e.Source
}

func (e WatcherStartedEvent) Type() event.Type {
	return EventTypeWatcherStarted
}

func (e WatcherStartedEvent) Error() error {
	return nil
}

func (e WatcherStartedEvent) Timestamp() time.Time {
	return e.timestamp
}

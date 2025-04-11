package watcher

import (
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
)

const EventTypeWatcherStopped event.Type = "watcher_stopped"

type WatcherStoppedEvent struct {
	Source    string
	err       error
	timestamp time.Time
}

func NewWatcherStoppedEvent(source string, err error) WatcherStoppedEvent {
	return WatcherStoppedEvent{
		Source:    source,
		err:       err,
		timestamp: time.Now(),
	}
}

func (e WatcherStoppedEvent) Message() string {
	if e.err != nil {
		return "Watcher stopped with error for source " + e.Source + ": " + e.err.Error()
	}
	return "Watcher stopped cleanly for source " + e.Source
}

func (e WatcherStoppedEvent) Type() event.Type {
	return EventTypeWatcherStopped
}

func (e WatcherStoppedEvent) Error() error {
	return e.err
}

func (e WatcherStoppedEvent) Timestamp() time.Time {
	return e.timestamp
}

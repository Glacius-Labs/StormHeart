package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"

	"github.com/fsnotify/fsnotify"
)

const debounceDelay = 250 * time.Millisecond

type FileWatcher struct {
	path          string
	sourceName    string
	dispatcher    *event.Dispatcher
	debounceTimer *time.Timer
	mu            sync.Mutex
}

func NewWatcher(path, sourceName string, dispatcher *event.Dispatcher) *FileWatcher {
	if path == "" {
		panic("FileWatcher requires a non-empty path")
	}

	if sourceName == "" {
		panic("FileWatcher requires a non-empty source name")
	}

	if dispatcher == nil {
		panic("FileWatcher requires a non-nil dispatcher")
	}

	return &FileWatcher{
		path:       path,
		sourceName: sourceName,
		dispatcher: dispatcher,
	}
}

func (w *FileWatcher) Watch(ctx context.Context) {
	startedEvent := watcher.NewWatcherStartedEvent(w.sourceName)
	w.dispatcher.Dispatch(ctx, startedEvent)

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		w.dispatchWatcherStopped(ctx, err)
		return
	}
	defer fsWatcher.Close()

	if err := fsWatcher.Add(w.path); err != nil {
		w.dispatchWatcherStopped(ctx, err)
		return
	}

	w.handleFileChange(ctx)

	for {
		select {
		case event, ok := <-fsWatcher.Events:
			if !ok {
				w.dispatchWatcherStopped(ctx, fmt.Errorf("events channel closed unexpectedly"))
				return
			}
			if event.Has(fsnotify.Remove | fsnotify.Rename) {
				w.dispatchWatcherStopped(ctx, fmt.Errorf("watched file was removed: %s", event.Name))
				return
			}
			w.handleWatchEvent(ctx, event)

		case err, ok := <-fsWatcher.Errors:
			if !ok {
				w.dispatchWatcherStopped(ctx, fmt.Errorf("error channel closed unexpectedly"))
				return
			}
			w.dispatchWatcherStopped(ctx, err)

		case <-ctx.Done():
			w.dispatchWatcherStopped(ctx, nil)

			w.mu.Lock()
			if w.debounceTimer != nil {
				w.debounceTimer.Stop()
				w.debounceTimer = nil
			}
			w.mu.Unlock()

			return
		}
	}
}

func (w *FileWatcher) handleFileChange(ctx context.Context) {
	data, err := os.ReadFile(w.path)
	if err != nil {
		w.dispatchDeploymentsReceived(ctx, nil, fmt.Errorf("failed to read file: %w", err))
		return
	}

	var deployments []model.Deployment
	err = json.Unmarshal(data, &deployments)

	w.dispatchDeploymentsReceived(ctx, deployments, err)
}

func (w *FileWatcher) handleWatchEvent(ctx context.Context, event fsnotify.Event) {
	if !event.Has(fsnotify.Create | fsnotify.Write) {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}

	w.debounceTimer = time.AfterFunc(debounceDelay, func() {
		if ctx.Err() != nil {
			return
		}
		w.handleFileChange(ctx)
	})
}

func (w *FileWatcher) dispatchWatcherStopped(ctx context.Context, err error) {
	event := watcher.NewWatcherStoppedEvent(w.sourceName, err)
	w.dispatcher.Dispatch(ctx, event)
}

func (w *FileWatcher) dispatchDeploymentsReceived(ctx context.Context, deployments []model.Deployment, err error) {
	event := watcher.NewDeploymentsReceivedEvent(w.sourceName, deployments, err)
	w.dispatcher.Dispatch(ctx, event)
}

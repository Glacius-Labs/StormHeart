package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

const debounceDelay = 250 * time.Millisecond

type FileWatcher struct {
	path          string
	sourceName    string
	handlerFunc   watcher.HandlerFunc
	logger        *zap.Logger
	debounceTimer *time.Timer
	mu            sync.Mutex
}

func NewWatcher(path, sourceName string, handlerFunc watcher.HandlerFunc, logger *zap.Logger) *FileWatcher {
	if path == "" {
		panic("FileWatcher requires a non-empty path")
	}

	if sourceName == "" {
		panic("FileWatcher requires a non-empty source name")
	}

	if handlerFunc == nil {
		panic("FileWatcher requires a non-nil handler func")
	}

	if logger == nil {
		panic("FileWatcher requires a non-nil logger")
	}

	return &FileWatcher{
		path:        path,
		sourceName:  sourceName,
		handlerFunc: handlerFunc,
		logger:      logger,
	}
}

func (w *FileWatcher) Watch(ctx context.Context) error {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fsWatcher.Close()

	if err := fsWatcher.Add(w.path); err != nil {
		return err
	}

	w.handleFileChange(ctx)

	for {
		select {
		case event, ok := <-fsWatcher.Events:
			if !ok {
				return nil
			}
			w.handleWatchEvent(ctx, event)

		case err, ok := <-fsWatcher.Errors:
			if !ok {
				w.logger.Error("File system error channel closed unexpectedly")
				return nil
			}
			w.logger.Error("File system error received", zap.Error(err))

		case <-ctx.Done():
			w.logger.Info("Initiating shutdown")
			watcher.PushEmptyDeployments(w.handlerFunc, w.sourceName)
			w.logger.Info("Shutdown complete")

			w.mu.Lock()
			if w.debounceTimer != nil {
				w.debounceTimer.Stop()
				w.debounceTimer = nil
			}
			w.mu.Unlock()

			return nil
		}
	}
}

func (w *FileWatcher) handleFileChange(ctx context.Context) error {
	data, err := os.ReadFile(w.path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var deployments []model.Deployment
	if err := json.Unmarshal(data, &deployments); err != nil {
		return fmt.Errorf("failed to unmarshal deployments: %w", err)
	}

	w.handlerFunc(ctx, w.sourceName, deployments)

	return nil
}

func (w *FileWatcher) handleWatchEvent(ctx context.Context, event fsnotify.Event) {
	if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
		return
	}

	w.logger.Info("Detected file change", zap.String("path", event.Name))

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}

	w.debounceTimer = time.AfterFunc(debounceDelay, func() {
		if ctx.Err() != nil {
			w.logger.Info("Debounced file change canceled due to shutdown")
			return
		}

		if err := w.handleFileChange(ctx); err != nil {
			w.logger.Error("Error while handling file change", zap.Error(err))
		}
	})
}

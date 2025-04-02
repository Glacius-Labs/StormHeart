package watcher

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/glacius-labs/StormHeart/internal/model"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

type FileWatcher struct {
	path       string
	sourceName string
	pushFunc   PushFunc
	logger     *zap.SugaredLogger
}

func NewFileWatcher(path, sourceName string, pushFunc PushFunc, logger *zap.SugaredLogger) *FileWatcher {
	if logger == nil {
		panic("FileWatcher requires a non-nil logger")
	}

	return &FileWatcher{
		path:       path,
		sourceName: sourceName,
		pushFunc:   pushFunc,
		logger:     logger,
	}
}

func (w *FileWatcher) Start(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(w.path); err != nil {
		return err
	}

	if err := w.loadAndPush(); err != nil {
		w.logger.Errorw("Initial load failed", "error", err)
	}

	var (
		mu            sync.Mutex
		debounce      *time.Timer
		debounceDelay = 250 * time.Millisecond
	)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					w.logger.Infow("Detected file change", "path", event.Name)

					mu.Lock()
					if debounce != nil {
						debounce.Stop()
					}
					debounce = time.AfterFunc(debounceDelay, func() {
						if err := w.loadAndPush(); err != nil {
							w.logger.Errorw("Reload failed", "error", err)
						}
					})
					mu.Unlock()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					w.logger.Errorw("Watcher error channel closed unexpectedly")
					return
				}

				w.logger.Errorw("Watcher error received", "error", err)
			case <-ctx.Done():
				w.logger.Infow("Shutting down file watcher")
				return
			}
		}
	}()

	<-done
	return nil
}

func (w *FileWatcher) loadAndPush() error {
	data, err := os.ReadFile(w.path)
	if err != nil {
		return err
	}

	var deployments []model.Deployment
	if err := json.Unmarshal(data, &deployments); err != nil {
		return err
	}

	w.logger.Infow("Loaded deployments", "count", len(deployments), "path", w.path)

	w.pushFunc(w.sourceName, deployments)

	return nil
}

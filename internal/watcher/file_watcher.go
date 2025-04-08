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
	logger     *zap.Logger
}

func NewFileWatcher(path, sourceName string, pushFunc PushFunc, logger *zap.Logger) *FileWatcher {
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

func (w *FileWatcher) Watch(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(w.path); err != nil {
		return err
	}

	if err := w.loadAndPush(ctx); err != nil {
		w.logger.Error("Initial load failed", zap.Error(err))
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
					w.logger.Info("Detected file change", zap.String("path", event.Name))

					mu.Lock()
					if debounce != nil {
						debounce.Stop()
					}
					debounce = time.AfterFunc(debounceDelay, func() {
						if ctx.Err() != nil {
							w.logger.Info("Debounce canceled due to shutdown")
							return
						}

						if err := w.loadAndPush(ctx); err != nil {
							w.logger.Error("Reload failed", zap.Error(err))
						}
					})
					mu.Unlock()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					w.logger.Error("Watcher error channel closed unexpectedly")
					return
				}

				w.logger.Error("Watcher error received", zap.Error(err))
			case <-ctx.Done():
				w.logger.Info("Shutting down file watcher")
				return
			}
		}
	}()

	<-done

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	w.pushFunc(shutdownCtx, w.sourceName, []model.Deployment{})

	return nil
}

func (w *FileWatcher) loadAndPush(ctx context.Context) error {
	data, err := os.ReadFile(w.path)
	if err != nil {
		return err
	}

	var deployments []model.Deployment
	if err := json.Unmarshal(data, &deployments); err != nil {
		return err
	}

	w.logger.Info(
		"Loaded deployments",
		zap.Int("count", len(deployments)),
		zap.String("path", w.path),
	)

	w.pushFunc(ctx, w.sourceName, deployments)

	return nil
}

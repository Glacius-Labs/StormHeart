package app

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/watcher"
	"go.uber.org/zap"
)

type App struct {
	watchers []watcher.Watcher
	logger   *zap.Logger
}

func (a *App) Start(ctx context.Context) error {
	for _, w := range a.watchers {
		go func(w watcher.Watcher) {
			if err := w.Watch(ctx); err != nil {
				if a.logger != nil {
					a.logger.Error("Watcher exited with error", zap.Error(err))
				}
			}
		}(w)
	}

	<-ctx.Done()
	return nil
}

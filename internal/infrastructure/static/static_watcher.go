package static

import (
	"context"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"go.uber.org/zap"
)

const SourceNameStaticWatcher = "static"

type StaticWatcher struct {
	deployments []model.Deployment
	pushFunc    watcher.PushFunc
	logger      *zap.Logger
}

func NewWatcher(deployments []model.Deployment, pushFunc watcher.PushFunc, logger *zap.Logger) *StaticWatcher {
	if logger == nil {
		panic("StaticWatcher requires a non-nil logger")
	}

	return &StaticWatcher{
		deployments: deployments,
		pushFunc:    pushFunc,
		logger:      logger,
	}
}

func (w *StaticWatcher) Watch(ctx context.Context) error {
	w.logger.Info("Pushing static deployments", zap.Int("count", len(w.deployments)))

	w.pushFunc(ctx, SourceNameStaticWatcher, w.deployments)

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	w.pushFunc(shutdownCtx, SourceNameStaticWatcher, []model.Deployment{})

	w.logger.Info("Static watcher shutdown")

	return nil
}

package static

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"go.uber.org/zap"
)

const SourceNameStaticWatcher = "static"

type StaticWatcher struct {
	deployments []model.Deployment
	handlerFunc watcher.HandlerFunc
	logger      *zap.Logger
}

func NewWatcher(deployments []model.Deployment, handlerFunc watcher.HandlerFunc, logger *zap.Logger) *StaticWatcher {
	if deployments == nil {
		deployments = []model.Deployment{}
	}

	if handlerFunc == nil {
		panic("StaticWatcher requires a non-nil handlerFunc")
	}

	if logger == nil {
		panic("StaticWatcher requires a non-nil logger")
	}

	return &StaticWatcher{
		deployments: deployments,
		handlerFunc: handlerFunc,
		logger:      logger,
	}
}

func (w *StaticWatcher) Watch(ctx context.Context) error {
	w.logger.Info("Pushing deployments", zap.Int("count", len(w.deployments)))

	w.handlerFunc(ctx, SourceNameStaticWatcher, w.deployments)

	<-ctx.Done()

	w.logger.Info("Initiating shutdown")
	watcher.PushEmptyDeployments(w.handlerFunc, SourceNameStaticWatcher)
	w.logger.Info("Shutdown complete")

	return nil
}

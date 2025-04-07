package watcher

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/model"
	"go.uber.org/zap"
)

const SourceNameStaticWatcher = "static"

type StaticWatcher struct {
	deployments []model.Deployment
	pushFunc    PushFunc
	logger      *zap.Logger
}

func NewStaticWatcher(deployments []model.Deployment, pushFunc PushFunc, logger *zap.Logger) *StaticWatcher {
	if logger == nil {
		panic("StaticWatcher requires a non-nil logger")
	}

	return &StaticWatcher{
		deployments: deployments,
		pushFunc:    pushFunc,
		logger:      logger,
	}
}

func (w *StaticWatcher) Start(ctx context.Context) error {
	w.logger.Info("Pushing static deployments", zap.Int("count", len(w.deployments)))
	w.pushFunc(ctx, SourceNameStaticWatcher, w.deployments)

	<-ctx.Done()

	w.pushFunc(ctx, SourceNameStaticWatcher, []model.Deployment{})
	w.logger.Info("Static watcher shutdown")

	return nil
}

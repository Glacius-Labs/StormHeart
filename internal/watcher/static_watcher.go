package watcher

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/model"
	"go.uber.org/zap"
)

const SourceNameStaticWatcher = "internal"

type StaticWatcher struct {
	deployments []model.Deployment
	pushFunc    PushFunc
	logger      *zap.SugaredLogger
}

func NewStaticWatcher(deployments []model.Deployment, pushFunc PushFunc, logger *zap.SugaredLogger) *StaticWatcher {
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
	w.logger.Infow("Pushing static deployments", "count", len(w.deployments), "source", SourceNameStaticWatcher)
	w.pushFunc(ctx, SourceNameStaticWatcher, w.deployments)

	<-ctx.Done()
	w.logger.Infow("Static watcher shutdown")
	return nil // << DO NOT return an error on context cancel
}

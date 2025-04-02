package watcher

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/model"
)

const SourceNameInternalWatcher = "internal"

type InternalWatcher struct {
	deplyoments []model.Deployment
	pushFunc    PushFunc
}

func NewInternalWatcher(deployments []model.Deployment, pushFunc PushFunc) *InternalWatcher {
	return &InternalWatcher{
		deplyoments: deployments,
		pushFunc:    pushFunc,
	}
}

func (w *InternalWatcher) Start(ctx context.Context) error {
	w.pushFunc(SourceNameInternalWatcher, w.deplyoments)
	return nil
}

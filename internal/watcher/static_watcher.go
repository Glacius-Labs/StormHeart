package watcher

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/model"
)

const SourceNameStaticWatcher = "internal"

type StaticWatcher struct {
	deplyoments []model.Deployment
	pushFunc    PushFunc
}

func NewStaticWatcher(deployments []model.Deployment, pushFunc PushFunc) *StaticWatcher {
	return &StaticWatcher{
		deplyoments: deployments,
		pushFunc:    pushFunc,
	}
}

func (w *StaticWatcher) Start(ctx context.Context) error {
	w.pushFunc(SourceNameStaticWatcher, w.deplyoments)
	return nil
}

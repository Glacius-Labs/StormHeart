package static

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
)

const SourceNameStaticWatcher = "static"

type StaticWatcher struct {
	deployments []model.Deployment
	dispatcher  *event.Dispatcher
}

func NewWatcher(deployments []model.Deployment, dispatcher *event.Dispatcher) *StaticWatcher {
	if deployments == nil {
		deployments = []model.Deployment{}
	}

	if dispatcher == nil {
		panic("dispatcher cannot be nil")
	}

	return &StaticWatcher{
		deployments: deployments,
		dispatcher:  dispatcher,
	}
}

func (w *StaticWatcher) Watch(ctx context.Context) {
	startedEvent := watcher.NewWatcherStartedEvent(SourceNameStaticWatcher)
	w.dispatcher.Dispatch(ctx, startedEvent)

	receivedDeploymentsEvent := watcher.NewDeploymentsReceivedEvent(SourceNameStaticWatcher, w.deployments, nil)
	w.dispatcher.Dispatch(ctx, receivedDeploymentsEvent)

	<-ctx.Done()

	stoppedEvent := watcher.NewWatcherStoppedEvent(SourceNameStaticWatcher, nil)
	w.dispatcher.Dispatch(ctx, stoppedEvent)
}

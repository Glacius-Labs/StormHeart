package handler

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/application/shared"
	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/reconciler"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
)

type WatcherStoppedHandler struct {
	registry   *shared.DeploymentsRegistry
	reconciler *reconciler.Reconciler
}

func NewWatcherStoppedHandler(registry *shared.DeploymentsRegistry, reconciler *reconciler.Reconciler) *WatcherStoppedHandler {
	if registry == nil {
		panic("registry cannot be nil")
	}
	if reconciler == nil {
		panic("reconciler cannot be nil")
	}

	return &WatcherStoppedHandler{
		registry:   registry,
		reconciler: reconciler,
	}
}

func (h *WatcherStoppedHandler) Handle(ctx context.Context, event event.Event) {
	shutdownEvent, ok := event.(watcher.WatcherStoppedEvent)
	if !ok {
		return
	}

	source := shutdownEvent.Source
	h.registry.Unregister(source)

	allDeployments := h.registry.GetAll()
	h.reconciler.Apply(ctx, allDeployments)
}

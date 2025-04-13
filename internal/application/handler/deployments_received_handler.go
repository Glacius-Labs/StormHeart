package handler

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/application/shared"
	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/reconciler"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
)

type DeploymentsReceivedHandler struct {
	registry   *shared.DeploymentsRegistry
	reconciler *reconciler.Reconciler
}

func NewDeploymentsReceivedHandler(registry *shared.DeploymentsRegistry, reconciler *reconciler.Reconciler) *DeploymentsReceivedHandler {
	if registry == nil {
		panic("registry cannot be nil")
	}
	if reconciler == nil {
		panic("reconciler cannot be nil")
	}

	return &DeploymentsReceivedHandler{
		registry:   registry,
		reconciler: reconciler,
	}
}

func (h *DeploymentsReceivedHandler) Handle(ctx context.Context, event event.Event) {
	deploymentsReceivedEvent, ok := event.(watcher.DeploymentsReceivedEvent)
	if !ok {
		return
	}

	source := deploymentsReceivedEvent.Source
	deployments := deploymentsReceivedEvent.Deployments

	h.registry.Register(source, deployments)

	allDeployments := h.registry.GetAll()

	h.reconciler.Apply(ctx, allDeployments)
}

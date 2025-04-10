package reconciler

import (
	"context"
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/runtime"
)

type Reconciler struct {
	Runtime    runtime.Runtime
	Dispatcher *event.Dispatcher
}

func NewReconciler(runtime runtime.Runtime, dispatcher *event.Dispatcher) *Reconciler {
	if runtime == nil {
		panic("Reconciler requires a non-nil Runtime")
	}

	if dispatcher == nil {
		panic("Reconciler requires a non-nil Dispatcher")
	}

	return &Reconciler{
		Runtime:    runtime,
		Dispatcher: dispatcher,
	}
}

func (r *Reconciler) Apply(ctx context.Context, desired []model.Deployment) {
	actual, err := r.Runtime.List(ctx)
	if err != nil {
		err = fmt.Errorf("failed to list running containers: %w", err)
		e := NewReconciliationPerformedEvent(nil, desired, err)
		r.Dispatcher.Dispatch(ctx, e)
		return
	}

	desiredMap := make(map[string]model.Deployment, len(desired))
	for _, d := range desired {
		desiredMap[d.Name] = d
	}

	actualMap := make(map[string]model.Deployment, len(actual))
	for _, a := range actual {
		actualMap[a.Name] = a
	}

	toStart, toStop := r.diff(desiredMap, actualMap)

	err = r.performDeployments(ctx, toStart, toStop)

	e := NewReconciliationPerformedEvent(actual, desired, err)
	r.Dispatcher.Dispatch(ctx, e)
}

func (r Reconciler) diff(desired, actual map[string]model.Deployment) ([]model.Deployment, []model.Deployment) {
	var toStart, toStop []model.Deployment

	// Determine what to start or restart
	for name, desiredDeployment := range desired {
		actualDeployment, exists := actual[name]
		if !exists {
			toStart = append(toStart, desiredDeployment)
		} else if !desiredDeployment.Equals(actualDeployment) {
			// Deployment exists but changed => must stop old and start new
			toStop = append(toStop, actualDeployment)
			toStart = append(toStart, desiredDeployment)
		}
	}

	// Determine what to stop (deployments no longer desired)
	for name, actualDeployment := range actual {
		if _, exists := desired[name]; !exists {
			toStop = append(toStop, actualDeployment)
		}
	}

	return toStart, toStop
}

func (r Reconciler) performDeployments(ctx context.Context, toStart, toStop []model.Deployment) error {
	hadDeploymentErrors := false

	for _, d := range toStop {
		err := r.Runtime.Remove(ctx, d.Name)

		if err != nil {
			hadDeploymentErrors = true
		}

		e := NewDeploymentRemovedEvent(d, err)
		r.Dispatcher.Dispatch(ctx, e)
	}

	for _, d := range toStart {
		err := r.Runtime.Deploy(ctx, d)

		if err != nil {
			hadDeploymentErrors = true
		}

		e := NewDeploymentCreatedEvent(d, err)
		r.Dispatcher.Dispatch(ctx, e)
	}

	if hadDeploymentErrors {
		return fmt.Errorf("one or more deployment actions failed during reconciliation")
	}

	return nil
}

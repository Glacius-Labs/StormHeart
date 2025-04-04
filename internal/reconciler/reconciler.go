package reconciler

import (
	"context"
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/model"
	"github.com/glacius-labs/StormHeart/internal/runtime"
	"go.uber.org/zap"
)

type Reconciler struct {
	Runtime runtime.Runtime
	Logger  *zap.SugaredLogger
}

func NewReconciler(runtime runtime.Runtime, logger *zap.SugaredLogger) *Reconciler {
	if runtime == nil {
		panic("Reconciler requires a non-nil Runtime")
	}

	if logger == nil {
		panic("Reconciler requires a non-nil Logger")
	}

	return &Reconciler{
		Runtime: runtime,
		Logger:  logger,
	}
}

func (r *Reconciler) Apply(ctx context.Context, desired []model.Deployment) error {
	actual, err := r.Runtime.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list running containers: %w", err)
	}

	desiredMap := make(map[string]model.Deployment, len(desired))
	for _, d := range desired {
		desiredMap[d.Name] = d
	}

	actualMap := make(map[string]model.Deployment, len(actual))
	for _, a := range actual {
		actualMap[a.Name] = a
	}

	var toStart, toStop []model.Deployment

	// Determine what to start or restart
	for name, desiredDeployment := range desiredMap {
		actualDeployment, exists := actualMap[name]
		if !exists || !desiredDeployment.Equals(actualDeployment) {
			toStart = append(toStart, desiredDeployment)
		}
	}

	// Determine what to stop
	for name, actualDeployment := range actualMap {
		if _, exists := desiredMap[name]; !exists {
			toStop = append(toStop, actualDeployment)
		}
	}

	var startErrs, stopErrs int

	for _, d := range toStop {
		if err := r.Runtime.Remove(ctx, d); err != nil {
			stopErrs++
			r.Logger.Errorw("Failed to stop container", "deployment", d, "error", err)
		} else {
			r.Logger.Infow("Stopped container", "deployment", d.Name)
		}
	}

	for _, d := range toStart {
		if err := r.Runtime.Deploy(ctx, d); err != nil {
			startErrs++
			r.Logger.Errorw("Failed to start container", "deployment", d, "error", err)
		} else {
			r.Logger.Infow("Started container", "deployment", d.Name)
		}
	}

	r.Logger.Infow("Reconciliation complete",
		"started", len(toStart),
		"stopped", len(toStop),
		"errors", startErrs+stopErrs,
	)

	if startErrs+stopErrs > 0 {
		return fmt.Errorf("reconciliation failed: %d start errors, %d stop errors", startErrs, stopErrs)
	}

	return nil
}

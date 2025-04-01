package reconciler

import (
	"github.com/glacius-labs/StormHeart/internal/deployment/model"
	"github.com/glacius-labs/StormHeart/internal/runtime"
)

type Reconciler struct {
	runtime runtime.Runtime
}

func NewReconciler(runtime runtime.Runtime) *Reconciler {
	return &Reconciler{
		runtime: runtime,
	}
}

func (r *Reconciler) Reconcile(desired []model.Deployment) error {
	actual, err := r.runtime.List()

	if err != nil {
		return err
	}

	desiredMap := make(map[string]model.Deployment)
	for _, d := range desired {
		desiredMap[d.Name] = d
	}

	actualMap := make(map[string]model.Deployment)
	for _, a := range actual {
		actualMap[a.Name] = a
	}

	for name, desiredDeployment := range desiredMap {
		actualDeployment, exists := actualMap[name]
		if !exists || !desiredDeployment.Equals(actualDeployment) {
			if err := r.runtime.Deploy(desiredDeployment); err != nil {
				return err
			}
		}
	}

	for name := range actualMap {
		if deployment, exists := desiredMap[name]; !exists {
			if err := r.runtime.Remove(deployment); err != nil {
				return err
			}
		}
	}

	return nil
}

package deployment

type Reconciler struct {
	deployer  Deployer
	inspector DeploymentInspector
}

func NewReconciler(deployer Deployer, inspector DeploymentInspector) *Reconciler {
	return &Reconciler{
		deployer:  deployer,
		inspector: inspector,
	}
}

func (r *Reconciler) Reconcile(desired []Deployment) error {
	actual, err := r.inspector.ListActiveDeployments()
	if err != nil {
		return err
	}

	desiredMap := make(map[string]Deployment)
	for _, d := range desired {
		desiredMap[d.Name] = d
	}

	actualMap := make(map[string]Deployment)
	for _, a := range actual {
		actualMap[a.Name] = a
	}

	for name, desiredDeployment := range desiredMap {
		actualDeployment, exists := actualMap[name]
		if !exists || !desiredDeployment.Equals(actualDeployment) {
			if err := r.deployer.Deploy(desiredDeployment); err != nil {
				return err
			}
		}
	}

	for name := range actualMap {
		if _, exists := desiredMap[name]; !exists {
			if err := r.deployer.Remove(name); err != nil {
				return err
			}
		}
	}

	return nil
}

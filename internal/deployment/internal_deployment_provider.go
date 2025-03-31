package deployment

import "context"

const InternalDeploymentProviderSourceName = "internal"

type InternalDeploymentProvider struct {
	deplyoments []Deployment
	push        DeploymentPush
}

func NewInternalDeploymentProvider(deployments []Deployment, push DeploymentPush) *InternalDeploymentProvider {
	return &InternalDeploymentProvider{
		deplyoments: deployments,
		push:        push,
	}
}

func (p *InternalDeploymentProvider) Start(ctx context.Context) error {
	p.push(InternalDeploymentProviderSourceName, p.deplyoments)
	return nil
}

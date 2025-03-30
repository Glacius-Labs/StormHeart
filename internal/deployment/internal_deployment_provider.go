package deployment

import "context"

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
	p.push(p.deplyoments)
	return nil
}

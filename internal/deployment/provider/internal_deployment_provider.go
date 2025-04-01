package provider

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/deployment/model"
)

const InternalDeploymentProviderSourceName = "internal"

type InternalDeploymentProvider struct {
	deplyoments []model.Deployment
	pushFunc    PushFunc
}

func NewInternalDeploymentProvider(deployments []model.Deployment, pushFunc PushFunc) *InternalDeploymentProvider {
	return &InternalDeploymentProvider{
		deplyoments: deployments,
		pushFunc:    pushFunc,
	}
}

func (p *InternalDeploymentProvider) Start(ctx context.Context) error {
	p.pushFunc(InternalDeploymentProviderSourceName, p.deplyoments)
	return nil
}

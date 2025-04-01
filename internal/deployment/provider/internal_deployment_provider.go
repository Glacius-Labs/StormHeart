package provider

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/deployment/model"
)

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
	p.pushFunc(p.deplyoments)
	return nil
}

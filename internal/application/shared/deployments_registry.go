package shared

import (
	"slices"
	"sync"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

type DeploymentsRegistry struct {
	mu      sync.Mutex
	sources map[string][]model.Deployment
}

func NewDeploymentsRegistry() *DeploymentsRegistry {
	return &DeploymentsRegistry{
		sources: make(map[string][]model.Deployment),
	}
}

func (dr *DeploymentsRegistry) Register(source string, deployments []model.Deployment) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	dr.sources[source] = deployments
}

func (dr *DeploymentsRegistry) Unregister(source string) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	delete(dr.sources, source)
}

func (r *DeploymentsRegistry) GetAll() []model.Deployment {
	r.mu.Lock()
	defer r.mu.Unlock()

	var combined []model.Deployment
	for _, deployments := range r.sources {
		combined = append(combined, deployments...)
	}

	var unique []model.Deployment
	for _, candidate := range combined {
		found := slices.ContainsFunc(unique, candidate.Equals)
		if !found {
			unique = append(unique, candidate)
		}
	}

	return unique
}

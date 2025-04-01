package runtime

import (
	"slices"

	"github.com/glacius-labs/StormHeart/internal/deployment/model"
)

type MockDeploymentRuntime struct {
	Active []model.Deployment
}

func NewMockDeploymentRuntime(initial []model.Deployment) *MockDeploymentRuntime {
	return &MockDeploymentRuntime{
		Active: append([]model.Deployment{}, initial...),
	}
}

func (m *MockDeploymentRuntime) Deploy(deployment model.Deployment) error {
	for i, existing := range m.Active {
		if existing.Name == deployment.Name {
			m.Active[i] = deployment
			return nil
		}
	}

	m.Active = append(m.Active, deployment)
	return nil
}

func (m *MockDeploymentRuntime) Remove(deployment model.Deployment) error {
	for i, existing := range m.Active {
		if existing.Name == deployment.Name {
			m.Active = slices.Delete(m.Active, i, i+1)
			return nil
		}
	}

	return nil
}

func (m *MockDeploymentRuntime) List() ([]model.Deployment, error) {
	copied := make([]model.Deployment, len(m.Active))
	copy(copied, m.Active)
	return copied, nil
}

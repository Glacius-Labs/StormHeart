package deployment

import (
	"slices"
)

type MockDeploymentRuntime struct {
	active []Deployment
}

func NewMockDeploymentRuntime(initial []Deployment) *MockDeploymentRuntime {
	return &MockDeploymentRuntime{
		active: append([]Deployment{}, initial...),
	}
}

func (m *MockDeploymentRuntime) Deploy(deployment Deployment) error {
	for i, existing := range m.active {
		if existing.Name == deployment.Name {
			m.active[i] = deployment
			return nil
		}
	}

	m.active = append(m.active, deployment)
	return nil
}

func (m *MockDeploymentRuntime) Remove(name string) error {
	for i, existing := range m.active {
		if existing.Name == name {
			m.active = slices.Delete(m.active, i, i+1)
			return nil
		}
	}

	return nil
}

func (m *MockDeploymentRuntime) ListActiveDeployments() ([]Deployment, error) {
	copied := make([]Deployment, len(m.active))
	copy(copied, m.active)
	return copied, nil
}

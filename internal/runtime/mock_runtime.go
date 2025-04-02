package runtime

import (
	"fmt"
	"slices"

	"github.com/glacius-labs/StormHeart/internal/model"
)

type MockRuntime struct {
	Active     []model.Deployment
	FailDeploy map[string]bool
	FailRemove map[string]bool
	FailList   bool
}

func NewMockRuntime(initial []model.Deployment) *MockRuntime {
	return &MockRuntime{
		Active: append([]model.Deployment{}, initial...),
	}
}

func (m *MockRuntime) Deploy(d model.Deployment) error {
	if m.FailDeploy[d.Name] {
		return fmt.Errorf("simulated deploy failure for %s", d.Name)
	}
	for i, existing := range m.Active {
		if existing.Name == d.Name {
			m.Active[i] = d
			return nil
		}
	}
	m.Active = append(m.Active, d)
	return nil
}

func (m *MockRuntime) Remove(d model.Deployment) error {
	if m.FailRemove[d.Name] {
		return fmt.Errorf("simulated remove failure for %s", d.Name)
	}
	for i, existing := range m.Active {
		if existing.Name == d.Name {
			m.Active = slices.Delete(m.Active, i, i+1)
			return nil
		}
	}
	return nil
}

func (m *MockRuntime) List() ([]model.Deployment, error) {
	if m.FailList {
		return nil, fmt.Errorf("simulated list failure")
	}
	copied := make([]model.Deployment, len(m.Active))
	copy(copied, m.Active)
	return copied, nil
}

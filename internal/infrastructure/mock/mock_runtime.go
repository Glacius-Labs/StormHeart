package mock

import (
	"context"
	"fmt"
	"slices"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

type MockRuntime struct {
	Active     []model.Deployment
	FailDeploy map[string]bool
	FailRemove map[string]bool
	FailList   bool
}

func NewRuntime(initial []model.Deployment) *MockRuntime {
	return &MockRuntime{
		Active: append([]model.Deployment{}, initial...),
	}
}

func (m *MockRuntime) Deploy(ctx context.Context, deployment model.Deployment) error {
	if m.FailDeploy[deployment.Name] {
		return fmt.Errorf("simulated deploy failure for %s", deployment.Name)
	}
	for i, existing := range m.Active {
		if existing.Name == deployment.Name {
			m.Active[i] = deployment
			return nil
		}
	}
	m.Active = append(m.Active, deployment)
	return nil
}

func (m *MockRuntime) Remove(ctx context.Context, deployment model.Deployment) error {
	if m.FailRemove[deployment.Name] {
		return fmt.Errorf("simulated remove failure for %s", deployment.Name)
	}
	for i, existing := range m.Active {
		if existing.Name == deployment.Name {
			m.Active = slices.Delete(m.Active, i, i+1)
			return nil
		}
	}
	return nil
}

func (m *MockRuntime) List(ctx context.Context) ([]model.Deployment, error) {
	if m.FailList {
		return nil, fmt.Errorf("simulated list failure")
	}
	copied := make([]model.Deployment, len(m.Active))
	copy(copied, m.Active)
	return copied, nil
}

package reconciler_test

import (
	"context"
	"testing"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/reconciler"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/mock"
	"github.com/stretchr/testify/require"
)

func TestNewReconciler_PanicsOnNilRuntime(t *testing.T) {
	require.Panics(t, func() {
		_ = reconciler.NewReconciler(nil, event.NewDispatcher())
	}, "expected panic when runtime is nil")
}

func TestNewReconciler_PanicsOnNilDispatcher(t *testing.T) {
	require.Panics(t, func() {
		_ = reconciler.NewReconciler(mock.NewRuntime([]model.Deployment{}), nil)
	}, "expected panic when dispatcher is nil")
}

func TestReconciler_Apply_NoDeployments(t *testing.T) {
	rt := mock.NewRuntime([]model.Deployment{})
	dispatcher := event.NewDispatcher()
	rec := reconciler.NewReconciler(rt, dispatcher)

	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	rec.Apply(context.Background(), []model.Deployment{})

	require.NotEmpty(t, handler.Events(), "expected at least one event (reconciliation)")
}

func TestReconciler_Apply_StartNewDeployment(t *testing.T) {
	rt := mock.NewRuntime([]model.Deployment{})
	dispatcher := event.NewDispatcher()
	rec := reconciler.NewReconciler(rt, dispatcher)

	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	desired := []model.Deployment{
		{Name: "db", Image: "postgres"},
	}

	rec.Apply(context.Background(), desired)

	require.GreaterOrEqual(t, len(handler.Events()), 2, "expected deployment created + reconciliation events")
}

func TestReconciler_Apply_RemoveDeployment(t *testing.T) {
	existing := []model.Deployment{
		{Name: "cache", Image: "redis"},
	}

	rt := mock.NewRuntime(existing)
	dispatcher := event.NewDispatcher()
	rec := reconciler.NewReconciler(rt, dispatcher)

	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	rec.Apply(context.Background(), []model.Deployment{}) // No desired deployments

	require.GreaterOrEqual(t, len(handler.Events()), 2, "expected deployment removed + reconciliation events")
}

func TestReconciler_Apply_ChangedDeployment(t *testing.T) {
	existing := []model.Deployment{
		{Name: "web", Image: "nginx:old"},
	}

	rt := mock.NewRuntime(existing)
	dispatcher := event.NewDispatcher()
	rec := reconciler.NewReconciler(rt, dispatcher)

	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	desired := []model.Deployment{
		{Name: "web", Image: "nginx:new"},
	}

	rec.Apply(context.Background(), desired)

	require.GreaterOrEqual(t, len(handler.Events()), 3, "expected deployment removed + deployment created + reconciliation events")
}

func TestReconciler_Apply_ListFails(t *testing.T) {
	rt := mock.NewRuntime([]model.Deployment{})
	rt.FailList = true // Simulate List() failure

	dispatcher := event.NewDispatcher()
	rec := reconciler.NewReconciler(rt, dispatcher)

	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	rec.Apply(context.Background(), []model.Deployment{})

	require.Len(t, handler.Events(), 1, "expected only reconciliation event with error")
	require.NotNil(t, handler.Events()[0].Error(), "expected reconciliation event to have error")
}

func TestReconciler_Apply_DeploymentActionFails(t *testing.T) {
	// Prepare runtime with a deployment to be removed
	existing := []model.Deployment{
		{Name: "service", Image: "nginx"},
	}

	rt := mock.NewRuntime(existing)
	rt.FailRemove = map[string]bool{
		"service": true, // Simulate remove failure
	}

	dispatcher := event.NewDispatcher()
	rec := reconciler.NewReconciler(rt, dispatcher)

	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	rec.Apply(context.Background(), []model.Deployment{}) // No desired deployments

	events := handler.Events()
	require.GreaterOrEqual(t, len(events), 2, "expected deployment removed + reconciliation events")

	// Find the reconciliation event
	var reconciliationEvent event.Event
	for _, e := range events {
		if e.Type() == "reconciliation_performed" {
			reconciliationEvent = e
			break
		}
	}

	require.NotNil(t, reconciliationEvent, "expected a reconciliation event")
	require.NotNil(t, reconciliationEvent.Error(), "expected reconciliation event to carry an error because deployment action failed")
}

func TestReconciler_Apply_DeploymentCreateFails(t *testing.T) {
	// Prepare runtime with no existing deployments
	rt := mock.NewRuntime([]model.Deployment{})
	rt.FailDeploy = map[string]bool{
		"db": true, // Simulate deploy failure
	}

	dispatcher := event.NewDispatcher()
	rec := reconciler.NewReconciler(rt, dispatcher)

	handler := mock.NewMockHandler()
	dispatcher.Register(handler)

	// Try to start a new deployment
	desired := []model.Deployment{
		{Name: "db", Image: "postgres"},
	}

	rec.Apply(context.Background(), desired)

	events := handler.Events()
	require.GreaterOrEqual(t, len(events), 2, "expected deployment created + reconciliation events")

	// Find the reconciliation event
	var reconciliationEvent event.Event
	for _, e := range events {
		if e.Type() == "reconciliation_performed" {
			reconciliationEvent = e
			break
		}
	}

	require.NotNil(t, reconciliationEvent, "expected a reconciliation event")
	require.NotNil(t, reconciliationEvent.Error(), "expected reconciliation event to carry an error because deployment creation failed")
}

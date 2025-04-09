package reconciler_test

import (
	"context"
	"testing"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/reconciler"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestReconciler_DeploysMissingContainers(t *testing.T) {
	r := mock.NewRuntime(nil)
	logger := zaptest.NewLogger(t)
	rec := reconciler.NewReconciler(r, logger)

	desired := []model.Deployment{
		{Name: "web", Image: "nginx"},
		{Name: "db", Image: "postgres"},
	}

	err := rec.Apply(context.Background(), desired)
	require.NoError(t, err)
	require.Len(t, r.Active, 2)
}

func TestReconciler_RemovesObsoleteContainers(t *testing.T) {
	initial := []model.Deployment{{Name: "stale", Image: "old"}}
	r := mock.NewRuntime(initial)
	logger := zaptest.NewLogger(t)
	rec := reconciler.NewReconciler(r, logger)

	err := rec.Apply(context.Background(), nil)
	require.NoError(t, err)
	require.Empty(t, r.Active)
}

func TestReconciler_RestartsChangedDeployment(t *testing.T) {
	initial := []model.Deployment{{Name: "api", Image: "v1"}}
	r := mock.NewRuntime(initial)
	logger := zaptest.NewLogger(t)
	rec := reconciler.NewReconciler(r, logger)

	desired := []model.Deployment{{Name: "api", Image: "v2"}}
	err := rec.Apply(context.Background(), desired)
	require.NoError(t, err)
	require.Equal(t, 1, len(r.Active))
	require.Equal(t, "v2", r.Active[0].Image)
}

func TestReconciler_DoesNothingWhenInSync(t *testing.T) {
	aligned := []model.Deployment{{Name: "cache", Image: "redis"}}
	r := mock.NewRuntime(aligned)
	logger := zaptest.NewLogger(t)
	rec := reconciler.NewReconciler(r, logger)

	err := rec.Apply(context.Background(), aligned)
	require.NoError(t, err)
	require.Len(t, r.Active, 1)
}

func TestReconciler_DeployFailureIsReported(t *testing.T) {
	r := &mock.MockRuntime{
		FailDeploy: map[string]bool{"web": true},
	}
	logger := zaptest.NewLogger(t)
	rec := reconciler.NewReconciler(r, logger)

	err := rec.Apply(context.Background(), []model.Deployment{
		{Name: "web", Image: "nginx"},
	})

	require.Error(t, err)
}

func TestReconciler_RemoveFailureIsReported(t *testing.T) {
	r := &mock.MockRuntime{
		Active:     []model.Deployment{{Name: "stale", Image: "busybox"}},
		FailRemove: map[string]bool{"stale": true},
	}
	logger := zaptest.NewLogger(t)
	rec := reconciler.NewReconciler(r, logger)

	err := rec.Apply(context.Background(), nil)
	require.Error(t, err)
}

func TestReconciler_ListFailureIsReported(t *testing.T) {
	r := &mock.MockRuntime{
		FailList: true,
	}
	logger := zaptest.NewLogger(t)
	rec := reconciler.NewReconciler(r, logger)

	err := rec.Apply(context.Background(), []model.Deployment{
		{Name: "api", Image: "latest"},
	})

	require.Error(t, err)
}

func TestNewReconciler_PanicsOnNilRuntime(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil Runtime, got none")
		}
	}()
	_ = reconciler.NewReconciler(nil, zaptest.NewLogger(t))
}

func TestNewReconciler_PanicsOnNilLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil Logger, got none")
		}
	}()
	r := mock.NewRuntime(nil)
	_ = reconciler.NewReconciler(r, nil)
}

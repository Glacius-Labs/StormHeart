package reconciler

import (
	"testing"

	"github.com/glacius-labs/StormHeart/internal/model"
	"github.com/glacius-labs/StormHeart/internal/runtime"
	"github.com/stretchr/testify/require"
)

func TestReconciler_ReconcileDeploysMissing(t *testing.T) {
	// Given: an empty runtime and one desired deployment
	mockRuntime := runtime.NewMockDeploymentRuntime([]model.Deployment{})
	rec := NewReconciler(mockRuntime)

	desired := model.Deployment{
		Name:  "flash-a",
		Image: "image:latest",
	}

	// When: we reconcile
	err := rec.Reconcile([]model.Deployment{desired})
	require.NoError(t, err)

	// Then: the deployment should have been added
	require.Len(t, mockRuntime.Active, 1)
	require.True(t, mockRuntime.Active[0].Equals(desired))
}

func TestReconciler_ReconcileRemovesObsolete(t *testing.T) {
	// Given: a runtime with one active deployment, but no desired ones
	active := model.Deployment{
		Name:  "flash-b",
		Image: "old:version",
	}

	mockRuntime := runtime.NewMockDeploymentRuntime([]model.Deployment{active})
	rec := NewReconciler(mockRuntime)

	// When: we reconcile with nothing desired
	err := rec.Reconcile([]model.Deployment{})
	require.NoError(t, err)

	// Then: the active deployment should be removed
	require.Empty(t, mockRuntime.Active)
}

func TestReconciler_ReconcileUpdatesMismatched(t *testing.T) {
	// Given: a runtime with a deployment that has a different image
	active := model.Deployment{
		Name:  "flash-c",
		Image: "image:v1",
	}
	desired := model.Deployment{
		Name:  "flash-c",
		Image: "image:v2",
	}

	mockRuntime := runtime.NewMockDeploymentRuntime([]model.Deployment{active})
	rec := NewReconciler(mockRuntime)

	// When: we reconcile with updated image
	err := rec.Reconcile([]model.Deployment{desired})
	require.NoError(t, err)

	// Then: the image should be updated
	require.Len(t, mockRuntime.Active, 1)
	require.True(t, mockRuntime.Active[0].Equals(desired))
}

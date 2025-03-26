package deployment

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReconciler_ReconcileDeploysMissing(t *testing.T) {
	// Given: an empty runtime and one desired deployment
	mockRuntime := &MockDeploymentRuntime{}
	rec := NewReconciler(mockRuntime, mockRuntime)

	desired := Deployment{
		Name:  "flash-a",
		Image: "image:latest",
	}

	// When: we reconcile
	err := rec.Reconcile([]Deployment{desired})
	require.NoError(t, err)

	// Then: the deployment should have been added
	require.Len(t, mockRuntime.active, 1)
	require.True(t, mockRuntime.active[0].Equals(desired))
}

func TestReconciler_ReconcileRemovesObsolete(t *testing.T) {
	// Given: a runtime with one active deployment, but no desired ones
	active := Deployment{
		Name:  "flash-b",
		Image: "old:version",
	}

	mockRuntime := &MockDeploymentRuntime{
		active: []Deployment{active},
	}
	rec := NewReconciler(mockRuntime, mockRuntime)

	// When: we reconcile with nothing desired
	err := rec.Reconcile([]Deployment{})
	require.NoError(t, err)

	// Then: the active deployment should be removed
	require.Empty(t, mockRuntime.active)
}

func TestReconciler_ReconcileUpdatesMismatched(t *testing.T) {
	// Given: a runtime with a deployment that has a different image
	active := Deployment{
		Name:  "flash-c",
		Image: "image:v1",
	}
	desired := Deployment{
		Name:  "flash-c",
		Image: "image:v2",
	}

	mockRuntime := &MockDeploymentRuntime{
		active: []Deployment{active},
	}
	rec := NewReconciler(mockRuntime, mockRuntime)

	// When: we reconcile with updated image
	err := rec.Reconcile([]Deployment{desired})
	require.NoError(t, err)

	// Then: the image should be updated
	require.Len(t, mockRuntime.active, 1)
	require.True(t, mockRuntime.active[0].Equals(desired))
}

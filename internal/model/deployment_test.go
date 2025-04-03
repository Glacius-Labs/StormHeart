package model_test

import (
	"testing"

	"github.com/glacius-labs/StormHeart/internal/model"
	"github.com/stretchr/testify/require"
)

func TestNewDeployment_ValidInput(t *testing.T) {
	d, err := model.NewDeployment("web", "nginx", model.DeploymentOptions{})
	require.NoError(t, err)
	require.Equal(t, "web", d.Name)
	require.Equal(t, "nginx", d.Image)
	require.NotNil(t, d.Labels)
	require.NotNil(t, d.Environment)
	require.Empty(t, d.Labels)
	require.Empty(t, d.Environment)
}

func TestNewDeployment_EmptyNameFails(t *testing.T) {
	_, err := model.NewDeployment("", "nginx", model.DeploymentOptions{})
	require.Error(t, err)
}

func TestNewDeployment_EmptyImageFails(t *testing.T) {
	_, err := model.NewDeployment("web", "", model.DeploymentOptions{})
	require.Error(t, err)
}

func TestDeployment_Equals_True(t *testing.T) {
	a := model.Deployment{
		Name:  "svc",
		Image: "nginx",
		Labels: map[string]string{
			"tier": "frontend",
		},
		Environment: map[string]string{
			"PORT": "80",
		},
	}

	b := model.Deployment{
		Name:  "svc",
		Image: "nginx",
		Labels: map[string]string{
			"tier": "frontend",
		},
		Environment: map[string]string{
			"PORT": "80",
		},
	}

	require.True(t, a.Equals(b))
}

func TestDeployment_Equals_False(t *testing.T) {
	base := model.Deployment{
		Name:  "svc",
		Image: "nginx",
		Labels: map[string]string{
			"tier": "frontend",
		},
		Environment: map[string]string{
			"PORT": "80",
		},
	}

	// Image mismatch
	diffImage := base
	diffImage.Image = "alpine"
	require.False(t, base.Equals(diffImage))

	// Label mismatch
	diffLabel := base
	diffLabel.Labels = map[string]string{
		"tier": "backend",
	}
	require.False(t, base.Equals(diffLabel))

	// Env mismatch
	diffEnv := base
	diffEnv.Environment = map[string]string{
		"PORT": "8080",
	}
	require.False(t, base.Equals(diffEnv))

	// Extra label
	diffExtraLabel := base
	diffExtraLabel.Labels = map[string]string{
		"tier":  "frontend",
		"extra": "true",
	}
	require.False(t, base.Equals(diffExtraLabel))

	// Extra environment variable
	diffExtraEnv := base
	diffExtraEnv.Environment = map[string]string{
		"PORT": "80",
		"MODE": "debug",
	}
	require.False(t, base.Equals(diffExtraEnv))
}

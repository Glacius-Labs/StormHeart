package event_test

import (
	"errors"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/application/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/stretchr/testify/require"
)

func TestDeploymentsReceivedEvent_Success(t *testing.T) {
	deployment := model.Deployment{
		Name:    "a source",
		Content: []byte("sensor data"),
	}

	e := event.NewDeploymentReceivedEvent(deployment, nil)

	require.Contains(t, e.Message(), "Received deployment from source a source", "expected correct deployment received message")
	require.Equal(t, event.EventTypeDeploymentReceived, e.Type(), "expected correct event type")
	require.Nil(t, e.Error(), "expected no error for successful deployments received")
	require.WithinDuration(t, time.Now(), e.Timestamp(), 5*time.Second, "expected timestamp to be recent")
	require.Equal(t, deployment, e.Deployment, "expected correct deployment reference")
}

func TestDeploymentsReceivedEvent_Failure(t *testing.T) {
	deployment := model.Deployment{
		Name:    "a source",
		Content: []byte(""),
	}
	err := errors.New("failed to read deployments file")

	e := event.NewDeploymentReceivedEvent(deployment, err)

	require.Contains(t, e.Message(), "Received deployment from source a source", "expected correct deployments received message")
	require.Equal(t, event.EventTypeDeploymentReceived, e.Type(), "expected correct event type")
	require.Equal(t, err, e.Error(), "expected correct error attached")
	require.WithinDuration(t, time.Now(), e.Timestamp(), 5*time.Second, "expected timestamp to be recent")
	require.Equal(t, deployment, e.Deployment, "expected correct deployment reference")
}

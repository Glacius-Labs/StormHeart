package watcher_test

import (
	"errors"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"github.com/stretchr/testify/require"
)

func TestDeploymentsReceivedEvent_Success(t *testing.T) {
	deployments := []model.Deployment{
		{Name: "sensor-service", Image: "sensors:v1"},
		{Name: "mqtt-broker", Image: "eclipse-mosquitto"},
	}

	event := watcher.NewDeploymentsReceivedEvent("file-watcher", deployments, nil)

	require.Contains(t, event.Message(), "Received 2 deployments from source file-watcher", "expected correct deployments received message")
	require.Equal(t, watcher.EventTypeDeploymentsReceived, event.Type(), "expected correct event type")
	require.Nil(t, event.Error(), "expected no error for successful deployments received")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "expected timestamp to be recent")

	require.Len(t, event.Deployments, 2, "expected deployments list to have correct length")
	require.Equal(t, "sensor-service", event.Deployments[0].Name, "expected first deployment name to match")
}

func TestDeploymentsReceivedEvent_Failure(t *testing.T) {
	deployments := []model.Deployment{}
	err := errors.New("failed to read deployments file")

	event := watcher.NewDeploymentsReceivedEvent("file-watcher", deployments, err)

	require.Contains(t, event.Message(), "Received 0 deployments from source file-watcher", "expected correct deployments received message")
	require.Equal(t, watcher.EventTypeDeploymentsReceived, event.Type(), "expected correct event type")
	require.Equal(t, err, event.Error(), "expected correct error attached")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "expected timestamp to be recent")

	require.Len(t, event.Deployments, 0, "expected empty deployments list")
}

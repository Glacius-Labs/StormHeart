package watcher_test

import (
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"github.com/stretchr/testify/require"
)

func TestDeploymentsReceivedEvent(t *testing.T) {
	deployments := []model.Deployment{
		{Name: "sensor-service", Image: "sensors:v1"},
		{Name: "mqtt-broker", Image: "eclipse-mosquitto"},
	}

	event := watcher.NewDeploymentsReceivedEvent("file-watcher", deployments)

	require.Contains(t, event.Message(), "Received 2 deployments from source file-watcher", "expected correct deployments received message")
	require.Equal(t, watcher.EventTypeDeploymentsReceived, event.Type(), "expected correct event type")
	require.Nil(t, event.Error(), "expected no error")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "timestamp should be recent")

	require.Len(t, event.Deployments, 2, "expected deployments list to have correct length")
	require.Equal(t, "sensor-service", event.Deployments[0].Name, "expected first deployment name to match")
}

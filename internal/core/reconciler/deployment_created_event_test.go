package reconciler_test

import (
	"errors"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/reconciler"
	"github.com/stretchr/testify/require"
)

func TestDeploymentCreatedEvent_Success(t *testing.T) {
	deployment := model.Deployment{
		Name:  "web",
		Image: "nginx",
	}

	event := reconciler.NewDeploymentCreatedEvent(deployment, nil)

	require.Equal(t, "Successfully created deployment web", event.Message(), "expected success message")
	require.Equal(t, reconciler.EventTypeDeploymentCreated, event.Type(), "expected correct event type")
	require.Nil(t, event.Error(), "expected no error on success")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "expected timestamp to be recent")
}

func TestDeploymentCreatedEvent_Failure(t *testing.T) {
	deployment := model.Deployment{
		Name:  "db",
		Image: "postgres",
	}

	err := errors.New("deployment failure")
	event := reconciler.NewDeploymentCreatedEvent(deployment, err)

	require.Contains(t, event.Message(), "Failed to create deployment db", "expected failure message prefix")
	require.Contains(t, event.Message(), "deployment failure", "expected error detail in message")
	require.Equal(t, reconciler.EventTypeDeploymentCreated, event.Type(), "expected correct event type")
	require.Equal(t, err, event.Error(), "expected correct error reference")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "expected timestamp to be recent")
}

package event_test

import (
	"errors"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/stretchr/testify/require"
)

func TestDeploymentCreatedEvent_Success(t *testing.T) {
	deployment := model.Deployment{
		Name:    "web",
		Content: nil,
	}

	e := event.NewDeploymentCreatedEvent(deployment, nil)

	require.Equal(t, "Successfully created deployment web", e.Message(), "expected success message")
	require.Equal(t, event.EventTypeDeploymentCreated, e.Type(), "expected correct event type")
	require.Nil(t, e.Error(), "expected no error on success")
	require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected timestamp to be recent")
}

func TestDeploymentCreatedEvent_Failure(t *testing.T) {
	deployment := model.Deployment{
		Name:    "db",
		Content: nil,
	}

	err := errors.New("deployment failure")
	e := event.NewDeploymentCreatedEvent(deployment, err)

	require.Contains(t, e.Message(), "Failed to create deployment db", "expected failure message prefix")
	require.Contains(t, e.Message(), "deployment failure", "expected error detail in message")
	require.Equal(t, event.EventTypeDeploymentCreated, e.Type(), "expected correct event type")
	require.Equal(t, err, e.Error(), "expected correct error reference")
	require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected timestamp to be recent")
}

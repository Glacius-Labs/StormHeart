package event_test

import (
	"errors"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/application/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/stretchr/testify/require"
)

func TestDeploymentRemovedEvent_Success(t *testing.T) {
	deployment := model.Deployment{
		Name:    "web",
		Content: make([]byte, 10),
	}

	e := event.NewDeploymentRemovedEvent(deployment, nil)

	require.Equal(t, "Successfully removed deployment web", e.Message(), "expected success message")
	require.Equal(t, event.EventTypeDeploymentRemoved, e.Type(), "expected correct event type")
	require.Nil(t, e.Error(), "expected no error on success")
	require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected timestamp to be recent")
}

func TestDeploymentRemovedEvent_Failure(t *testing.T) {
	deployment := model.Deployment{
		Name:    "db",
		Content: nil,
	}

	err := errors.New("removal failure")
	e := event.NewDeploymentRemovedEvent(deployment, err)

	require.Contains(t, e.Message(), "Failed to remove deployment db", "expected failure message prefix")
	require.Contains(t, e.Message(), "removal failure", "expected error detail in message")
	require.Equal(t, event.EventTypeDeploymentRemoved, e.Type(), "expected correct event type")
	require.Equal(t, err, e.Error(), "expected correct error reference")
	require.WithinDuration(t, time.Now(), e.Timestamp(), time.Second, "expected timestamp to be recent")
}

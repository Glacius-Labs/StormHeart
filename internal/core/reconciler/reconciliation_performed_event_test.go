package reconciler_test

import (
	"errors"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/core/reconciler"
	"github.com/stretchr/testify/require"
)

func TestReconciliationPerformedEvent_Success(t *testing.T) {
	actual := []model.Deployment{
		{Name: "db", Image: "postgres"},
	}
	desired := []model.Deployment{
		{Name: "db", Image: "postgres"},
	}

	event := reconciler.NewReconciliationPerformedEvent(actual, desired, nil)

	require.Contains(t, event.Message(), "Reconciliation completed successfully", "expected success message")
	require.Contains(t, event.Message(), "1 desired deployments", "expected correct desired count")
	require.Contains(t, event.Message(), "1 actual deployments", "expected correct actual count")
	require.Equal(t, reconciler.EventTypeReconciliationPerformed, event.Type(), "expected correct event type")
	require.Nil(t, event.Error(), "expected no error on success")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "expected timestamp to be recent")
}

func TestReconciliationPerformedEvent_Failure(t *testing.T) {
	actual := []model.Deployment{}
	desired := []model.Deployment{
		{Name: "web", Image: "nginx"},
	}

	err := errors.New("reconciliation failure")
	event := reconciler.NewReconciliationPerformedEvent(actual, desired, err)

	require.Contains(t, event.Message(), "Reconciliation failed with error: reconciliation failure", "expected failure message")
	require.Equal(t, reconciler.EventTypeReconciliationPerformed, event.Type(), "expected correct event type")
	require.Equal(t, err, event.Error(), "expected correct error reference")
	require.WithinDuration(t, time.Now(), event.Timestamp(), time.Second, "expected timestamp to be recent")
}

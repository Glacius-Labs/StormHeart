package reconciler

import (
	"fmt"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
)

const EventTypeReconciliationPerformed event.Type = "reconciliation_performed"

type ReconciliationPerformedEvent struct {
	Actual    []model.Deployment
	Desired   []model.Deployment
	err       error
	timestamp time.Time
}

func NewReconciliationPerformedEvent(actual, desired []model.Deployment, err error) ReconciliationPerformedEvent {
	copiedActual := append([]model.Deployment(nil), actual...)
	copiedDesired := append([]model.Deployment(nil), desired...)

	return ReconciliationPerformedEvent{
		Actual:    copiedActual,
		Desired:   copiedDesired,
		err:       err,
		timestamp: time.Now(),
	}
}

func (e ReconciliationPerformedEvent) Message() string {
	if e.err != nil {
		return "Reconciliation failed with error: " + e.err.Error()
	}

	return "Reconciliation completed successfully with " +
		fmt.Sprintf("%d desired deployments and %d actual deployments", len(e.Desired), len(e.Actual))
}

func (e ReconciliationPerformedEvent) Type() event.Type {
	return EventTypeReconciliationPerformed
}

func (e ReconciliationPerformedEvent) Error() error {
	return e.err
}

func (e ReconciliationPerformedEvent) Timestamp() time.Time {
	return e.timestamp
}

package watcher

import (
	"fmt"
	"time"

	"slices"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
)

const EventTypeDeploymentsReceived event.Type = "deployments_received"

type DeploymentsReceivedEvent struct {
	Source      string
	Deployments []model.Deployment
	timestamp   time.Time
}

func NewDeploymentsReceivedEvent(source string, deployments []model.Deployment) DeploymentsReceivedEvent {
	copied := slices.Clone(deployments)

	return DeploymentsReceivedEvent{
		Source:      source,
		Deployments: copied,
		timestamp:   time.Now(),
	}
}

func (e DeploymentsReceivedEvent) Message() string {
	return fmt.Sprintf("Received %d deployments from source %s", len(e.Deployments), e.Source)
}

func (e DeploymentsReceivedEvent) Type() event.Type {
	return EventTypeDeploymentsReceived
}

func (e DeploymentsReceivedEvent) Error() error {
	return nil
}

func (e DeploymentsReceivedEvent) Timestamp() time.Time {
	return e.timestamp
}

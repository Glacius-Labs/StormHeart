package reconciler

import (
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
)

const EventTypeDeploymentRemoved event.Type = "deployment_removed"

type DeploymentRemovedEvent struct {
	Deployment model.Deployment
	err        error
	timestamp  time.Time
}

func NewDeploymentRemovedEvent(deployment model.Deployment, err error) DeploymentRemovedEvent {
	return DeploymentRemovedEvent{
		Deployment: deployment,
		err:        err,
		timestamp:  time.Now(),
	}
}

func (e DeploymentRemovedEvent) Message() string {
	if e.err != nil {
		return "Failed to remove deployment " + e.Deployment.Name + ": " + e.err.Error()
	}
	return "Successfully removed deployment " + e.Deployment.Name
}

func (e DeploymentRemovedEvent) Type() event.Type {
	return EventTypeDeploymentRemoved
}

func (e DeploymentRemovedEvent) Error() error {
	return e.err
}

func (e DeploymentRemovedEvent) Timestamp() time.Time {
	return e.timestamp
}

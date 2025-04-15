package event

import (
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

const EventTypeDeploymentCreated EventType = "deployment_created"

type DeploymentCreatedEvent struct {
	Deployment model.Deployment
	err        error
	timestamp  time.Time
}

func NewDeploymentCreatedEvent(deployment model.Deployment, err error) DeploymentCreatedEvent {
	return DeploymentCreatedEvent{
		Deployment: deployment,
		err:        err,
		timestamp:  time.Now(),
	}
}

func (e DeploymentCreatedEvent) Message() string {
	if e.err != nil {
		return "Failed to create deployment " + e.Deployment.Name + ": " + e.err.Error()
	}
	return "Successfully created deployment " + e.Deployment.Name
}

func (e DeploymentCreatedEvent) Type() EventType {
	return EventTypeDeploymentCreated
}

func (e DeploymentCreatedEvent) Error() error {
	return e.err
}

func (e DeploymentCreatedEvent) Timestamp() time.Time {
	return e.timestamp
}

package event

import (
	"fmt"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

const EventTypeDeploymentReceived EventType = "deployment_received"

type DeploymentReceivedEvent struct {
	Deployment model.Deployment
	err        error
	timestamp  time.Time
}

func NewDeploymentReceivedEvent(deployment model.Deployment, err error) *DeploymentReceivedEvent {
	return &DeploymentReceivedEvent{
		Deployment: deployment,
		err:        err,
		timestamp:  time.Now(),
	}
}

func (e DeploymentReceivedEvent) Message() string {
	return fmt.Sprintf("Received deployment from source %s", e.Deployment.Name)
}

func (e DeploymentReceivedEvent) Type() EventType {
	return EventTypeDeploymentReceived
}

func (e DeploymentReceivedEvent) Error() error {
	return e.err
}

func (e DeploymentReceivedEvent) Timestamp() time.Time {
	return e.timestamp
}

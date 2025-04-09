package event

import (
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

type DeploymentCreatedEvent struct {
	Deployment model.Deployment
	Error      error
	Timestamp  time.Time
}

func (e DeploymentCreatedEvent) ToDispatcherEvent() DispatcherEvent {
	var msg string
	if e.Error != nil {
		msg = "Failed to create deployment " + e.Deployment.Name
	} else {
		msg = "Successfully created deployment " + e.Deployment.Name
	}

	return DispatcherEvent{
		Message:   msg,
		Type:      EventTypeDeployment,
		Error:     e.Error,
		Timestamp: e.Timestamp,
	}
}

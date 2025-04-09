package event

import (
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

type DeploymentRemovedEvent struct {
	Deployment model.Deployment
	Error      error
	Timestamp  time.Time
}

func (e DeploymentRemovedEvent) ToDispatcherEvent() DispatcherEvent {
	var msg string
	if e.Error != nil {
		msg = "Failed to remove deployment " + e.Deployment.Name
	} else {
		msg = "Successfully removed deployment " + e.Deployment.Name
	}

	return DispatcherEvent{
		Message:   msg,
		Type:      EventTypeDeployment,
		Error:     e.Error,
		Timestamp: e.Timestamp,
	}
}

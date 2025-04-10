package event

import (
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/model"
)

const EventTypeDeploymentReceived event.Type = "deployment_received"

type DeploymentReceivedEvent struct {
	Deployment model.Deployment
	Source     string
	Error      error
	Timestamp  time.Time
}

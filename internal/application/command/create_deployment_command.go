package command

import "github.com/glacius-labs/StormHeart/internal/core/model"

type CreateDeploymentCommand struct {
	Deployment model.Deployment
}

func (CreateDeploymentCommand) IsCommand() {}

package command

import "github.com/glacius-labs/StormHeart/internal/core/model"

type UpdateDeploymentCommand struct {
	Deployment model.Deployment
}

func (UpdateDeploymentCommand) IsCommand() {}

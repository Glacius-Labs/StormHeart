package command

import "github.com/glacius-labs/StormHeart/internal/core/model"

const CommandTypeUpdateDeployment CommandType = "update-deployment"

type UpdateDeploymentCommand struct {
	Deployment model.Deployment
}

func (c UpdateDeploymentCommand) CommandType() CommandType {
	return CommandTypeUpdateDeployment
}

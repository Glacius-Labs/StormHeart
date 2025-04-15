package command

import "github.com/glacius-labs/StormHeart/internal/core/model"

const CommandTypeCreateDeployment CommandType = "create-deployment"

type CreateDeploymentCommand struct {
	Deployment model.Deployment
}

func (c CreateDeploymentCommand) CommandType() CommandType {
	return CommandTypeCreateDeployment
}

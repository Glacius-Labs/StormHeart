package command

const CommandTypeRemoveDeployment CommandType = "remove-deployment"

type RemoveDeploymentCommand struct {
	Name string
}

func (c RemoveDeploymentCommand) CommandType() CommandType {
	return CommandTypeRemoveDeployment
}

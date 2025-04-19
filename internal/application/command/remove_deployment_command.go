package command

type RemoveDeploymentCommand struct {
	Name string
}

func (RemoveDeploymentCommand) IsCommand() {}

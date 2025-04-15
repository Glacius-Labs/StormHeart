package handler

import (
	"context"
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/application/command"
)

type RemoveDeploymentHandler struct {
}

func (h *RemoveDeploymentHandler) Handle(ctx context.Context, cmd command.Command) error {
	_, ok := cmd.(command.RemoveDeploymentCommand)
	if !ok {
		return fmt.Errorf("invalid command type: expected RemoveDeploymentCommand, got %T", cmd)
	}

	panic("implement me") // TODO: implement me

	return nil
}

func (h *RemoveDeploymentHandler) CommandType() command.CommandType {
	return command.CommandTypeRemoveDeployment
}

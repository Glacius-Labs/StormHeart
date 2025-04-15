package handler

import (
	"context"
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/application/command"
)

type UpdateDeploymentHandler struct {
}

func (h *UpdateDeploymentHandler) Handle(ctx context.Context, cmd command.Command) error {
	_, ok := cmd.(command.UpdateDeploymentCommand)
	if !ok {
		return fmt.Errorf("invalid command type: expected UpdateDeploymentCommand, got %T", cmd)
	}

	panic("implement me")

	return nil
}

func (h *UpdateDeploymentHandler) CommandType() command.CommandType {
	return command.CommandTypeUpdateDeployment
}

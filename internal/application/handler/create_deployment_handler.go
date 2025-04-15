package handler

import (
	"context"
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/application/command"
)

type CreateDeploymentHandler struct {
}

func (h *CreateDeploymentHandler) Handle(ctx context.Context, cmd command.Command) error {
	_, ok := cmd.(command.CreateDeploymentCommand)
	if !ok {
		return fmt.Errorf("invalid command type: expected CreateDeploymentCommand, got %T", cmd)
	}

	panic("implement me") // TODO: implement me

	return nil
}

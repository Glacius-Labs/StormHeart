package commandhandler

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/application/command"
)

type UpdateDeploymentHandler struct {
}

func (h *UpdateDeploymentHandler) Handle(ctx context.Context, cmd command.UpdateDeploymentCommand) error {
	panic("implement me")
}

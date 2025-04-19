package commandhandler

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/application/command"
)

type RemoveDeploymentHandler struct {
}

func (h *RemoveDeploymentHandler) Handle(ctx context.Context, cmd command.RemoveDeploymentCommand) error {
	panic("implement me")
}

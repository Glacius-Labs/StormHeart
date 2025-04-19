package commandhandler

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/application/command"
)

type Handler[C command.Command] interface {
	Handle(ctx context.Context, cmd C) error
}

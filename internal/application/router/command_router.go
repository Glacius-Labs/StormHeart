package router

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/application/command"
	"github.com/glacius-labs/StormHeart/internal/application/handler"
)

type CommandRouter interface {
	RegisterHandler(h handler.Handler) error
	Publish(ctx context.Context, cmd command.Command) error
	Start() error
	Close() error
}

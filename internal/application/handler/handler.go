package handler

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/application/command"
)

type Handler interface {
	Handle(ctx context.Context, cmd command.Command) error
	CommandType() command.CommandType
}

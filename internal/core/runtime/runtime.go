package runtime

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

type Runtime interface {
	Deploy(ctx context.Context, deployment model.Deployment) error
	Remove(ctx context.Context, name string) error
}

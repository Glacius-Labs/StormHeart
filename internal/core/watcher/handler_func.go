package watcher

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

type HandlerFunc func(ctx context.Context, sourceName string, deployments []model.Deployment)

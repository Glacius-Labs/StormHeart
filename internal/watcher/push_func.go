package watcher

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/model"
)

type PushFunc func(ctx context.Context, sourceName string, deployments []model.Deployment)

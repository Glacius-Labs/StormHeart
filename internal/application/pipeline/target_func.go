package pipeline

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

type TargetFunc func(ctx context.Context, deployments []model.Deployment) error

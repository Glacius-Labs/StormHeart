package pipeline

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/model"
)

type TargetFunc func(ctx context.Context, deployments []model.Deployment) error

package pipeline

import (
	"context"
	"slices"

	"github.com/glacius-labs/StormHeart/internal/core/model"
)

func Deduplicator() Decorator {
	return func(next TargetFunc) TargetFunc {
		return func(ctx context.Context, deployments []model.Deployment) error {
			var out []model.Deployment

			for _, candidate := range deployments {
				found := slices.ContainsFunc(out, candidate.Equals)
				if !found {
					out = append(out, candidate)
				}
			}

			return next(ctx, out)
		}
	}
}

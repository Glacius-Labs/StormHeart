package pipeline

import "github.com/glacius-labs/StormHeart/internal/core/model"

type Filter interface {
	Apply(in []model.Deployment) []model.Deployment
}

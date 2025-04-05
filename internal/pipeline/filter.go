package pipeline

import "github.com/glacius-labs/StormHeart/internal/model"

type Filter interface {
	Apply(in []model.Deployment) []model.Deployment
}

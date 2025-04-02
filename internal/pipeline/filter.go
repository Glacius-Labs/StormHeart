package pipeline

import "github.com/glacius-labs/StormHeart/internal/model"

type Filter interface {
	Filter(in []model.Deployment) []model.Deployment
}

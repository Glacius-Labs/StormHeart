package pipeline

import "github.com/glacius-labs/StormHeart/internal/model"

type Transformer interface {
	Transform(in []model.Deployment) []model.Deployment
}

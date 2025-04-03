package pipeline

import (
	"slices"

	"github.com/glacius-labs/StormHeart/internal/model"
)

type Deduplicator struct{}

func NewDeduplicator() *Deduplicator {
	return &Deduplicator{}
}

func (d Deduplicator) Filter(in []model.Deployment) []model.Deployment {
	var out []model.Deployment
	for _, candidate := range in {
		found := slices.ContainsFunc(out, candidate.Equals)
		if !found {
			out = append(out, candidate)
		}
	}
	return out
}

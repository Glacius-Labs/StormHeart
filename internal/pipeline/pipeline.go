package pipeline

import (
	"slices"
	"sync"

	"github.com/glacius-labs/StormHeart/internal/model"
	"go.uber.org/zap"
)

type Pipeline struct {
	mu           sync.Mutex
	sources      map[string][]model.Deployment
	Target       func([]model.Deployment)
	Transformers []Transformer
	logger       *zap.SugaredLogger
}

func NewPipeline(
	target func([]model.Deployment),
	logger *zap.SugaredLogger,
	transformers ...Transformer,
) *Pipeline {
	if target == nil {
		panic("Pipeline requires a non-nil Target")
	}
	if logger == nil {
		panic("Pipeline requires a non-nil logger")
	}

	return &Pipeline{
		sources:      make(map[string][]model.Deployment),
		Target:       target,
		Transformers: transformers,
		logger:       logger,
	}
}

func (p *Pipeline) Push(source string, deployments []model.Deployment) {
	p.mu.Lock()

	prev, exists := p.sources[source]
	if exists && slices.EqualFunc(prev, deployments, func(a, b model.Deployment) bool {
		return a.Equals(b)
	}) {
		p.mu.Unlock()
		p.logger.Debugw("Push skipped, no changes from source", "source", source)
		return
	}

	p.sources[source] = deployments

	var combined []model.Deployment
	for _, ds := range p.sources {
		combined = append(combined, ds...)
	}
	p.mu.Unlock()

	if exists {
		added, removed := diffDeployments(prev, deployments)
		p.logger.Infow("Source delta",
			"source", source,
			"added", len(added),
			"removed", len(removed),
		)
	}

	originalCount := len(combined)
	transformed := combined
	for _, t := range p.Transformers {
		transformed = t.Transform(transformed)
	}
	finalCount := len(transformed)

	p.logger.Infow("Pipeline push processed",
		"source", source,
		"inputCount", len(deployments),
		"totalBeforeTransforms", originalCount,
		"totalAfterTransforms", finalCount,
	)

	p.Target(transformed)
}

func diffDeployments(old, new []model.Deployment) (added, removed []model.Deployment) {
	for _, d := range new {
		if !slices.ContainsFunc(old, d.Equals) {
			added = append(added, d)
		}
	}
	for _, d := range old {
		if !slices.ContainsFunc(new, d.Equals) {
			removed = append(removed, d)
		}
	}
	return
}

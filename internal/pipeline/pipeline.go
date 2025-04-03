package pipeline

import (
	"sync"

	"github.com/glacius-labs/StormHeart/internal/model"
	"go.uber.org/zap"
)

type Pipeline struct {
	mu      sync.Mutex
	sources map[string][]model.Deployment
	Target  func(deployments []model.Deployment)
	Filters []Filter
	logger  *zap.SugaredLogger
}

func NewPipeline(
	target func(deployments []model.Deployment),
	logger *zap.SugaredLogger,
	filters ...Filter,
) *Pipeline {
	if target == nil {
		panic("Pipeline requires a non-nil Target")
	}
	if logger == nil {
		panic("Pipeline requires a non-nil logger")
	}

	return &Pipeline{
		sources: make(map[string][]model.Deployment),
		Target:  target,
		Filters: filters,
		logger:  logger,
	}
}

func (p *Pipeline) Push(source string, deployments []model.Deployment) {
	p.mu.Lock()

	p.sources[source] = deployments

	var combined []model.Deployment
	for _, ds := range p.sources {
		combined = append(combined, ds...)
	}
	p.mu.Unlock()

	originalCount := len(combined)

	filtered := combined
	for _, t := range p.Filters {
		filtered = t.Filter(filtered)
	}
	finalCount := len(filtered)

	p.logger.Infow("Pipeline push processed",
		"source", source,
		"inputCount", len(deployments),
		"totalBeforeTransforms", originalCount,
		"totalAfterTransforms", finalCount,
	)

	p.Target(filtered)
}

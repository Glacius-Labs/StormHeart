package pipeline

import (
	"context"
	"sync"

	"github.com/glacius-labs/StormHeart/internal/model"
	"go.uber.org/zap"
)

type Pipeline struct {
	mu         sync.Mutex
	sources    map[string][]model.Deployment
	TargetFunc TargetFunc
	Filters    []Filter
	logger     *zap.SugaredLogger
}

func NewPipeline(
	targetFunc TargetFunc,
	logger *zap.SugaredLogger,
	filters ...Filter,
) *Pipeline {
	if targetFunc == nil {
		panic("Pipeline requires a non-nil Target")
	}
	if logger == nil {
		panic("Pipeline requires a non-nil logger")
	}

	return &Pipeline{
		sources:    make(map[string][]model.Deployment),
		TargetFunc: targetFunc,
		Filters:    filters,
		logger:     logger,
	}
}

func (p *Pipeline) Push(ctx context.Context, source string, deployments []model.Deployment) {
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
		filtered = t.Apply(filtered)
	}
	finalCount := len(filtered)

	p.logger.Infow("Pipeline push processed",
		"source", source,
		"inputCount", len(deployments),
		"totalBeforeFilters", originalCount,
		"totalAfterFilters", finalCount,
	)

	if err := p.TargetFunc(ctx, filtered); err != nil {
		p.logger.Errorw("Target function failed", "error", err)
	}
}

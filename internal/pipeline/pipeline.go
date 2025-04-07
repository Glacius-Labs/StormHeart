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
	logger     *zap.Logger
}

func NewPipeline(
	targetFunc TargetFunc,
	logger *zap.Logger,
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

	p.logger.Info("Pipeline push processed",
		zap.String("source", source),
		zap.Int("inputCount", len(deployments)),
		zap.Int("totalBeforeFilters", originalCount),
		zap.Int("totalAfterFilters", finalCount),
	)

	if err := p.TargetFunc(ctx, filtered); err != nil {
		p.logger.Error("Target function failed", zap.Error(err))
	}
}

package pipeline

import (
	"context"
	"sync"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"go.uber.org/zap"
)

type Pipeline struct {
	mu         sync.Mutex
	sources    map[string][]model.Deployment
	targetFunc TargetFunc
	logger     *zap.Logger
}

func NewPipeline(
	targetFunc TargetFunc,
	logger *zap.Logger,
) *Pipeline {
	if targetFunc == nil {
		panic("Pipeline requires a non-nil Target")
	}

	if logger == nil {
		panic("Pipeline requires a non-nil logger")
	}

	return &Pipeline{
		sources:    make(map[string][]model.Deployment),
		targetFunc: targetFunc,
		logger:     logger,
	}
}

func (p *Pipeline) Use(decorator Decorator) {
	p.targetFunc = decorator(p.targetFunc)
}

func (p *Pipeline) Push(ctx context.Context, source string, deployments []model.Deployment) {
	p.mu.Lock()

	p.sources[source] = deployments

	var combined []model.Deployment
	for _, ds := range p.sources {
		combined = append(combined, ds...)
	}

	if err := p.targetFunc(ctx, combined); err != nil {
		p.logger.Error("Target function failed", zap.Error(err))
	}

	p.mu.Unlock()

	p.logger.Info("Pipeline push processed",
		zap.String("source", source),
		zap.Int("inputCount", len(deployments)),
	)
}

package app

import (
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/application/pipeline"
	"github.com/glacius-labs/StormHeart/internal/core/reconciler"
	"github.com/glacius-labs/StormHeart/internal/core/runtime"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
	"go.uber.org/zap"
)

type Builder struct {
	runtime    runtime.Runtime
	reconciler *reconciler.Reconciler
	pipeline   *pipeline.Pipeline
	watchers   []watcher.Watcher
	logger     *zap.Logger
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithLogger(logger *zap.Logger) *Builder {
	b.logger = logger
	return b
}

func (b *Builder) WithRuntime(r runtime.Runtime) *Builder {
	b.runtime = r
	return b
}

func (b *Builder) WithReconciler(r *reconciler.Reconciler) *Builder {
	b.reconciler = r
	return b
}

func (b *Builder) WithPipeline(p *pipeline.Pipeline) *Builder {
	b.pipeline = p
	return b
}

func (b *Builder) WithWatcher(w watcher.Watcher) *Builder {
	b.watchers = append(b.watchers, w)
	return b
}

func (b *Builder) Build() (*App, error) {
	if b.logger == nil {
		return nil, fmt.Errorf("missing logger")
	}

	if b.runtime == nil {
		return nil, fmt.Errorf("missing runtime")
	}

	if b.reconciler == nil {
		return nil, fmt.Errorf("missing reconciler")
	}

	if b.pipeline == nil {
		return nil, fmt.Errorf("missing pipeline")
	}

	if len(b.watchers) == 0 {
		return nil, fmt.Errorf("no watchers configured")
	}

	return &App{
		watchers: b.watchers,
	}, nil
}

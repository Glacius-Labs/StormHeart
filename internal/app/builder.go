package app

import (
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/core/runtime"
	"github.com/glacius-labs/StormHeart/internal/core/watcher"
)

type Builder struct {
	runtime  runtime.Runtime
	watchers []watcher.Watcher
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithRuntime(r runtime.Runtime) *Builder {
	b.runtime = r
	return b
}

func (b *Builder) WithWatcher(w watcher.Watcher) *Builder {
	b.watchers = append(b.watchers, w)
	return b
}

func (b *Builder) Build() (*App, error) {
	if b.runtime == nil {
		return nil, fmt.Errorf("missing runtime")
	}

	if len(b.watchers) == 0 {
		return nil, fmt.Errorf("no watchers configured")
	}

	return &App{
		watchers: b.watchers,
	}, nil
}

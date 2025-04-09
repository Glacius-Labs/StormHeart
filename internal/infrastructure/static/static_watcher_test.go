package static_test

import (
	"context"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/static"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap/zaptest"
)

func TestStaticWatcher_Start_PushesDeployments(t *testing.T) {
	expected := []model.Deployment{
		{Name: "test", Image: "alpine:latest"},
	}

	var called bool
	var gotSource string
	var gotDeployments []model.Deployment

	push := func(ctx context.Context, source string, deployments []model.Deployment) {
		called = true
		gotSource = source
		gotDeployments = deployments
	}

	logger := zaptest.NewLogger(t)
	w := static.NewWatcher(expected, push, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		err := w.Watch(ctx)
		assert.NoError(t, err)
		close(done)
	}()

	// Wait a little to ensure handlerFunc is called
	// (Normally this is instant, but tiny wait is safe)
	<-time.After(50 * time.Millisecond)

	assert.True(t, called, "handlerFunc should have been called")
	assert.Equal(t, static.SourceNameStaticWatcher, gotSource)
	assert.Equal(t, expected, gotDeployments)

	cancel()

	<-done
}

func TestStaticWatcher_NewWatcher_CreatesEmptyCollectionOnNilDeployments(t *testing.T) {
	assert.NotPanics(t, func() {
		static.NewWatcher(nil, func(context.Context, string, []model.Deployment) {}, zaptest.NewLogger(t))
	})
}

func TestStaticWatcher_NewWatcher_PanicsOnNilHandlerFunc(t *testing.T) {
	assert.Panics(t, func() {
		static.NewWatcher(nil, nil, zaptest.NewLogger(t))
	})
}

func TestStaticWatcher_NewWatcher_PanicsOnNilLogger(t *testing.T) {
	assert.Panics(t, func() {
		static.NewWatcher(nil, func(context.Context, string, []model.Deployment) {}, nil)
	})
}

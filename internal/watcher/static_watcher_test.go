package watcher_test

import (
	"context"
	"testing"
	"time"

	"github.com/glacius-labs/StormHeart/internal/model"
	"github.com/glacius-labs/StormHeart/internal/watcher"
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

	logger := zaptest.NewLogger(t).Sugar()
	w := watcher.NewStaticWatcher(expected, push, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		err := w.Start(ctx)
		assert.NoError(t, err)
		close(done)
	}()

	// Wait a little to ensure pushFunc is called
	// (Normally this is instant, but tiny wait is safe)
	<-time.After(50 * time.Millisecond)

	// Validate pushFunc
	assert.True(t, called, "pushFunc should have been called")
	assert.Equal(t, watcher.SourceNameStaticWatcher, gotSource)
	assert.Equal(t, expected, gotDeployments)

	// Now cancel context to shut down watcher
	cancel()

	// Wait for Start(ctx) to return
	<-done
}

func TestStaticWatcher_PanicsOnNilLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil logger, but got none")
		}
	}()

	_ = watcher.NewStaticWatcher(nil, func(context.Context, string, []model.Deployment) {}, nil)
}

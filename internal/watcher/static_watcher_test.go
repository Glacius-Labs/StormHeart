package watcher_test

import (
	"context"
	"testing"

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

	push := func(source string, deployments []model.Deployment) {
		called = true
		gotSource = source
		gotDeployments = deployments
	}

	logger := zaptest.NewLogger(t).Sugar()
	w := watcher.NewStaticWatcher(expected, push, logger)

	err := w.Start(context.Background())
	assert.NoError(t, err)
	assert.True(t, called, "pushFunc should have been called")
	assert.Equal(t, watcher.SourceNameStaticWatcher, gotSource)
	assert.Equal(t, expected, gotDeployments)
}

func TestStaticWatcher_PanicsOnNilLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil logger, but got none")
		}
	}()

	_ = watcher.NewStaticWatcher(nil, func(string, []model.Deployment) {}, nil)
}

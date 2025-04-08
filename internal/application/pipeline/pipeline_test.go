package pipeline_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/glacius-labs/StormHeart/internal/application/pipeline"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestPipeline_SinglePush_CallsTarget(t *testing.T) {
	var called []model.Deployment
	targetFunc := func(ctx context.Context, deployments []model.Deployment) error {
		called = deployments
		return nil
	}

	p := pipeline.NewPipeline(
		targetFunc,
		zaptest.NewLogger(t),
	)

	p.Push(context.Background(), "source1", []model.Deployment{
		{Name: "a", Image: "x"},
	})

	require.Len(t, called, 1)
	require.Equal(t, "a", called[0].Name)
}

func TestNewPipeline_PanicsOnNilTarget(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil target, but got none")
		}
	}()

	_ = pipeline.NewPipeline(nil, zaptest.NewLogger(t))
}

func TestNewPipeline_PanicsOnNilLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil logger, but got none")
		}
	}()

	_ = pipeline.NewPipeline(func(context.Context, []model.Deployment) error { return nil }, nil)
}

func TestPipeline_TargetError_IsHandled(t *testing.T) {
	var called bool

	targetFunc := func(ctx context.Context, deployments []model.Deployment) error {
		called = true
		return fmt.Errorf("simulated target failure")
	}

	p := pipeline.NewPipeline(
		targetFunc,
		zaptest.NewLogger(t),
	)

	p.Push(context.Background(), "source1", []model.Deployment{
		{Name: "fail", Image: "broken"},
	})

	require.True(t, called, "target function should have been called despite error")
}

func TestPipeline_UseDecorator_IsCalled(t *testing.T) {
	var decoratorCalled bool

	targetFunc := func(ctx context.Context, deployments []model.Deployment) error {
		return nil
	}

	p := pipeline.NewPipeline(
		targetFunc,
		zaptest.NewLogger(t),
	)

	p.Use(mockDecorator(&decoratorCalled))

	p.Push(context.Background(), "source1", []model.Deployment{
		{Name: "trigger", Image: "any"},
	})

	require.True(t, decoratorCalled, "decorator should have been called")
}

func mockDecorator(wasCalled *bool) pipeline.Decorator {
	return func(next pipeline.TargetFunc) pipeline.TargetFunc {
		return func(ctx context.Context, deployments []model.Deployment) error {
			*wasCalled = true
			return next(ctx, deployments)
		}
	}
}

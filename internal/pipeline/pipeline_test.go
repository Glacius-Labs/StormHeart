package pipeline_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/glacius-labs/StormHeart/internal/model"
	"github.com/glacius-labs/StormHeart/internal/pipeline"
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
		zaptest.NewLogger(t).Sugar(),
		pipeline.Deduplicator{},
	)

	p.Push(context.Background(), "source1", []model.Deployment{
		{Name: "a", Image: "x"},
	})

	require.Len(t, called, 1)
	require.Equal(t, "a", called[0].Name)
}

func TestPipeline_MultiSource_Deduplication(t *testing.T) {
	var received []model.Deployment
	targetFunc := func(ctx context.Context, deployments []model.Deployment) error {
		received = deployments
		return nil
	}

	p := pipeline.NewPipeline(
		targetFunc,
		zaptest.NewLogger(t).Sugar(),
		pipeline.Deduplicator{},
	)

	ctx := context.Background()

	p.Push(ctx, "file", []model.Deployment{
		{Name: "worker", Image: "alpine"},
	})

	p.Push(ctx, "broker", []model.Deployment{
		{Name: "worker", Image: "alpine"},
		{Name: "db", Image: "pg"},
	})

	require.Len(t, received, 2)
}

func TestNewPipeline_PanicsOnNilTarget(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil target, but got none")
		}
	}()

	_ = pipeline.NewPipeline(nil, zaptest.NewLogger(t).Sugar())
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
		zaptest.NewLogger(t).Sugar(),
		pipeline.Deduplicator{},
	)

	p.Push(context.Background(), "source1", []model.Deployment{
		{Name: "fail", Image: "broken"},
	})

	require.True(t, called, "target function should have been called despite error")
}

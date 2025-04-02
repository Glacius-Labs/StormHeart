package pipeline_test

import (
	"testing"

	"github.com/glacius-labs/StormHeart/internal/model"
	"github.com/glacius-labs/StormHeart/internal/pipeline"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestPipeline_SinglePush_CallsTarget(t *testing.T) {
	var called []model.Deployment
	target := func(in []model.Deployment) {
		called = in
	}

	p := pipeline.NewPipeline(
		target,
		zaptest.NewLogger(t).Sugar(),
		pipeline.Deduplicator{},
	)

	p.Push("source1", []model.Deployment{
		{Name: "a", Image: "x"},
	})

	require.Len(t, called, 1)
	require.Equal(t, "a", called[0].Name)
}

func TestPipeline_MultiSource_Deduplication(t *testing.T) {
	var received []model.Deployment
	target := func(in []model.Deployment) {
		received = in
	}

	p := pipeline.NewPipeline(
		target,
		zaptest.NewLogger(t).Sugar(),
		pipeline.Deduplicator{},
	)

	p.Push("file", []model.Deployment{
		{Name: "worker", Image: "alpine"},
	})

	p.Push("broker", []model.Deployment{
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

	_ = pipeline.NewPipeline(func([]model.Deployment) {}, nil)
}

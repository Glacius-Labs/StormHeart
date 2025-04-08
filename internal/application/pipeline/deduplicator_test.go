package pipeline_test

import (
	"context"
	"testing"

	"github.com/glacius-labs/StormHeart/internal/application/pipeline"
	"github.com/glacius-labs/StormHeart/internal/core/model"
	"github.com/stretchr/testify/require"
)

func TestDeduplicator_RemovesDuplicates(t *testing.T) {
	var received []model.Deployment

	target := func(ctx context.Context, deployments []model.Deployment) error {
		received = deployments
		return nil
	}

	deduplicator := pipeline.Deduplicator()
	decorated := deduplicator(target)

	input := []model.Deployment{
		{Name: "web", Image: "nginx"},
		{Name: "web", Image: "nginx"}, // duplicate
		{Name: "db", Image: "postgres"},
	}

	err := decorated(nil, input)
	require.NoError(t, err)

	require.Len(t, received, 2)
	require.ElementsMatch(t, []model.Deployment{
		{Name: "web", Image: "nginx"},
		{Name: "db", Image: "postgres"},
	}, received)
}

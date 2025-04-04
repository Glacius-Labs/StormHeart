package pipeline_test

import (
	"testing"

	"github.com/glacius-labs/StormHeart/internal/model"
	"github.com/glacius-labs/StormHeart/internal/pipeline"
	"github.com/stretchr/testify/require"
)

func TestDeduplicator_RemovesDuplicates(t *testing.T) {
	d := pipeline.NewDeduplicator()

	input := []model.Deployment{
		{Name: "web", Image: "nginx"},
		{Name: "web", Image: "nginx"}, // duplicate
		{Name: "db", Image: "postgres"},
	}

	result := d.Apply(input)
	require.Len(t, result, 2)
}

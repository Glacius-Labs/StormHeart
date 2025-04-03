package runtime_test

import (
	"strings"
	"testing"

	"github.com/glacius-labs/StormHeart/internal/model"
	"github.com/glacius-labs/StormHeart/internal/runtime"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDockerRuntime_Deploy(t *testing.T) {
	rt, err := runtime.NewDockerRuntime()

	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	d := testDeployment()

	err = rt.Deploy(d)
	require.NoError(t, err)

	t.Cleanup(func() { _ = rt.Remove(d) })
}

func TestDockerRuntime_List(t *testing.T) {
	rt, err := runtime.NewDockerRuntime()

	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	d := testDeployment()

	require.NoError(t, rt.Deploy(d))
	t.Cleanup(func() { _ = rt.Remove(d) })

	all, err := rt.List()
	require.NoError(t, err)

	var found bool
	for _, x := range all {
		if x.Name == d.Name {
			found = true
		}
	}
	require.True(t, found)
}

func TestDockerRuntime_Remove(t *testing.T) {
	rt, err := runtime.NewDockerRuntime()

	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	d := testDeployment()

	require.NoError(t, rt.Deploy(d))
	require.NoError(t, rt.Remove(d))

	all, err := rt.List()
	require.NoError(t, err)
	for _, x := range all {
		require.NotEqual(t, d.Name, x.Name)
	}
}

func testDeployment() model.Deployment {
	return model.Deployment{
		Name:        "stormheart-test-" + strings.ToLower(uuid.NewString()),
		Image:       "alpine",
		Environment: map[string]string{"FOO": "bar"},
		Labels:      map[string]string{"component": "test"},
	}
}

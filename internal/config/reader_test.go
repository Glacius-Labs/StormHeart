package config_test

import (
	"strings"
	"testing"

	"github.com/glacius-labs/StormHeart/internal/config"
	"github.com/stretchr/testify/require"
)

func TestReader_ValidConfig(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"logLevel": "info",
		"runtime": {
			"type": "docker"
		},
		"watchers": {
			"files": [
				{"name": "static", "path": "/some/path/deployments.json"}
			]
		}
	}`

	r := config.NewReader()
	cfg, err := r.Read(strings.NewReader(input))
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, "stormer-alpha", cfg.Identifier)
	require.Equal(t, "info", cfg.LogLevel)
	require.Equal(t, "docker", cfg.Runtime.Type)
	require.Len(t, cfg.Watchers.Files, 1)
	require.Equal(t, "static", cfg.Watchers.Files[0].Name)
	require.Equal(t, "/some/path/deployments.json", cfg.Watchers.Files[0].Path)
}

func TestReader_InvalidConfig_MissingIdentifier(t *testing.T) {
	input := `{
		"logLevel": "info",
		"runtime": {
			"type": "docker"
		},
		"watchers": {
			"files": [
				{"name": "static", "path": "/some/path/deployments.json"}
			]
		}
	}`

	r := config.NewReader()
	_, err := r.Read(strings.NewReader(input))
	require.Error(t, err)
	require.Contains(t, err.Error(), "identifier must not be empty")
}

func TestReader_InvalidConfig_InvalidLogLevel(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"logLevel": "superverbose",
		"runtime": {
			"type": "docker"
		},
		"watchers": {
			"files": [
				{"name": "static", "path": "/some/path/deployments.json"}
			]
		}
	}`

	r := config.NewReader()
	_, err := r.Read(strings.NewReader(input))
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid log level")
}

func TestReader_InvalidConfig_UnsupportedRuntime(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"logLevel": "info",
		"runtime": {
			"type": "unknown-runtime"
		},
		"watchers": {
			"files": [
				{"name": "static", "path": "/some/path/deployments.json"}
			]
		}
	}`

	r := config.NewReader()
	_, err := r.Read(strings.NewReader(input))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported runtime type")
}

func TestReader_InvalidConfig_NoWatchers(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"logLevel": "info",
		"runtime": {
			"type": "docker"
		},
		"watchers": {
			"files": []
		}
	}`

	r := config.NewReader()
	_, err := r.Read(strings.NewReader(input))
	require.Error(t, err)
	require.Contains(t, err.Error(), "at least one file watcher")
}

func TestReader_InvalidConfig_FileWatcherMissingName(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"logLevel": "info",
		"runtime": {
			"type": "docker"
		},
		"watchers": {
			"files": [
				{"name": "", "path": "/some/path/deployments.json"}
			]
		}
	}`

	r := config.NewReader()
	_, err := r.Read(strings.NewReader(input))
	require.Error(t, err)
	require.Contains(t, err.Error(), "file watcher name must not be empty")
}

func TestReader_InvalidConfig_FileWatcherMissingPath(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"logLevel": "info",
		"runtime": {
			"type": "docker"
		},
		"watchers": {
			"files": [
				{"name": "static", "path": ""}
			]
		}
	}`

	r := config.NewReader()
	_, err := r.Read(strings.NewReader(input))
	require.Error(t, err)
	require.Contains(t, err.Error(), "file watcher path must not be empty")
}

func TestReader_InvalidJSON(t *testing.T) {
	input := `{ invalid-json }`

	r := config.NewReader()
	_, err := r.Read(strings.NewReader(input))

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode config")
}

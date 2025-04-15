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
		"stormlink": {
			"host": "localhost",
			"port": 1234
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
	require.Equal(t, "localhost", cfg.StormLink.Host)
	require.Equal(t, 1234, cfg.StormLink.Port)
	require.Len(t, cfg.Watchers.Files, 1)
	require.Equal(t, "static", cfg.Watchers.Files[0].Name)
	require.Equal(t, "/some/path/deployments.json", cfg.Watchers.Files[0].Path)
}

func TestReader_InvalidConfig_MissingIdentifier(t *testing.T) {
	input := `{
		"stormlink": {
			"host": "localhost",
			"port": 1234
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

func TestReader_InvalidConfig_MissingStormLink(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"watchers": {
			"files": [
				{"name": "static", "path": "/some/path/deployments.json"}
			]
		}
	}`

	r := config.NewReader()
	_, err := r.Read(strings.NewReader(input))
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid stormlink")
}

func TestReader_InvalidConfig_MissingStormLinkHost(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"stormlink": {
			"port": 1234
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
	require.Contains(t, err.Error(), "host must not be empty")
}

func TestReader_InvalidConfig_EmptyStormLinkHost(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"stormlink": {
			"host": "",
			"port": 1234
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
	require.Contains(t, err.Error(), "host must not be empty")
}

func TestReader_InvalidConfig_MissingStormLinkPort(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"stormlink": {
			"host": "localhost"
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
	require.Contains(t, err.Error(), "port must be a valid TCP port (1-65535)")
}

func TestReader_InvalidConfig_NegativeStormLinkPort(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"stormlink": {
			"host": "localhost",
			"port": -1
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
	require.Contains(t, err.Error(), "port must be a valid TCP port (1-65535)")
}

func TestReader_InvalidConfig_StormLinkPortTooBig(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"stormlink": {
			"host": "localhost",
			"port": 99999
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
	require.Contains(t, err.Error(), "port must be a valid TCP port (1-65535)")
}

func TestReader_InvalidConfig_NoWatchers(t *testing.T) {
	input := `{
		"identifier": "stormer-alpha",
		"stormlink": {
			"host": "localhost",
			"port": 1234
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
		"stormlink": {
			"host": "localhost",
			"port": 1234
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
		"stormlink": {
			"host": "localhost",
			"port": 1234
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

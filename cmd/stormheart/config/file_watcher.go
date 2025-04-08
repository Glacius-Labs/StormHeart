package config

import "fmt"

type FileWatcher struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (fw FileWatcher) validate() error {
	if fw.Name == "" {
		return fmt.Errorf("file watcher name must not be empty")
	}

	if fw.Path == "" {
		return fmt.Errorf("file watcher path must not be empty")
	}

	// NOTE: We no longer stat-check the file here
	// because the watcher will handle file existence/runtime errors separately

	return nil
}

package config

import "fmt"

type Watchers struct {
	Files []FileWatcher `json:"files"`
}

func (w Watchers) validate() error {
	if len(w.Files) == 0 {
		return fmt.Errorf("at least one file watcher must be configured")
	}

	for i, fw := range w.Files {
		if err := fw.validate(); err != nil {
			return fmt.Errorf("file watcher at index %d: %w", i, err)
		}
	}

	return nil
}

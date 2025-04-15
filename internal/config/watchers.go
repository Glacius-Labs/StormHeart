package config

import "fmt"

type Watchers struct {
	Files []FileWatcher `json:"files"`
}

func (w Watchers) validate() error {
	if len(w.Files) == 0 {
		return nil
	}

	for _, fw := range w.Files {
		if err := fw.validate(); err != nil {
			return fmt.Errorf("file watcher %s: %w", fw.Name, err)
		}
	}

	return nil
}

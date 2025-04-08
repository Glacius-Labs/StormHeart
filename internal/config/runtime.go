package config

import "fmt"

type Runtime struct {
	Type string `json:"type"`
}

func (r Runtime) validate() error {
	if r.Type != "docker" {
		return fmt.Errorf("unsupported runtime type: %s", r.Type)
	}

	return nil
}

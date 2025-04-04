package config

import "fmt"

type Config struct {
	Identifier string   `json:"identifier"`
	LogLevel   string   `json:"logLevel"`
	Runtime    Runtime  `json:"runtime"`
	Watchers   Watchers `json:"watchers"`
}

func (c Config) validate() error {
	if c.Identifier == "" {
		return fmt.Errorf("identifier must not be empty")
	}

	switch c.LogLevel {
	case "debug", "info", "warn", "error", "":
		// OK
	default:
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	if err := c.Runtime.validate(); err != nil {
		return fmt.Errorf("invalid runtime: %w", err)
	}

	if err := c.Watchers.validate(); err != nil {
		return fmt.Errorf("invalid watchers: %w", err)
	}

	return nil
}

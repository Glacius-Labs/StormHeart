package config

import "fmt"

type Config struct {
	Identifier string    `json:"identifier"`
	StormLink  StormLink `json:"stormlink"`
	Watchers   Watchers  `json:"watchers"`
}

func (c Config) validate() error {
	if c.Identifier == "" {
		return fmt.Errorf("identifier must not be empty")
	}

	if err := c.StormLink.validate(); err != nil {
		return fmt.Errorf("invalid stormlink: %w", err)
	}

	if err := c.Watchers.validate(); err != nil {
		return fmt.Errorf("invalid watchers: %w", err)
	}

	return nil
}

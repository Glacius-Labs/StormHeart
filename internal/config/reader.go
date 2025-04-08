package config

import (
	"encoding/json"
	"fmt"
	"io"
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) Read(src io.Reader) (Config, error) {
	var cfg Config

	if err := json.NewDecoder(src).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to decode config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

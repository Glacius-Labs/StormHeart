package config

import "fmt"

type StormLink struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (s StormLink) validate() error {
	if s.Host == "" {
		return fmt.Errorf("host must not be empty")
	}
	if s.Port <= 0 || s.Port > 65535 {
		return fmt.Errorf("port must be a valid TCP port (1-65535)")
	}
	return nil
}

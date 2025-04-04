package main

import (
	"os"

	"github.com/glacius-labs/StormHeart/internal/config"
)

func loadConfig(path string) (config.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return config.Config{}, err
	}
	defer file.Close()

	reader := config.NewReader()
	return reader.Read(file)
}

package main

import (
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/bootstrap"
)

const configPath = "config.json"

func main() {
	cfg, err := loadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	logger, err := setupLogger(cfg.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("failed to setup logger: %v", err))
	}
	defer logger.Sync()

	logger.Infow("Logger initialized", "level", cfg.LogLevel)

	ctx := setupSignalContext(logger)

	if err := bootstrap.Bootstrap(ctx, cfg, logger); err != nil {
		logger.Fatalw("Failed to bootstrap system", "error", err)
	}

	<-ctx.Done()

	logger.Infow("Shutdown complete")
}

package main

import (
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/app"
	"github.com/glacius-labs/StormHeart/internal/application/pipeline"
	"github.com/glacius-labs/StormHeart/internal/application/reconciler"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/docker"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/file"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/mqtt"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/static"
	"go.uber.org/zap"
)

const configPath = "config.json"
const stormLinkSource = "stormlink"

func main() {
	cfg, err := loadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	logger, err := setupLogger()
	if err != nil {
		panic(fmt.Sprintf("failed to setup logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Logger initialized")

	ctx := setupSignalContext(logger)

	runtime, err := docker.NewRuntime()

	if err != nil {
		logger.Fatal("Failed to create runtime", zap.Error(err))
	}

	reconciler := reconciler.NewReconciler(
		runtime,
		logger.With(zap.String("component", "reconciler")),
	)

	pipeline := pipeline.NewPipeline(
		reconciler.Apply,
		logger.With(zap.String("component", "pipeline")),
		pipeline.NewDeduplicator(),
	)

	staticWatcher := static.NewWatcher(
		staticDeployments,
		pipeline.Push,
		logger.With(zap.String("component", "watcher"), zap.String("source", "static")),
	)

	mqttTopic := fmt.Sprintf("stormfleet/%s/deployments", cfg.Identifier)
	mqttUrl := fmt.Sprintf("tcp://%s:%d", cfg.StormLink.Host, cfg.StormLink.Port)
	mqttClient := mqtt.NewPahoClient(mqttUrl)

	mqqtWatcher := mqtt.NewWatcher(
		mqttClient,
		mqttTopic,
		stormLinkSource,
		pipeline.Push,
		logger.With(zap.String("component", "watcher"), zap.String("source", "mqtt")),
	)

	builder := app.
		NewBuilder().
		WithLogger(logger).
		WithRuntime(runtime).
		WithReconciler(reconciler).
		WithPipeline(pipeline).
		WithWatcher(staticWatcher).
		WithWatcher(mqqtWatcher)

	for _, watcherConfig := range cfg.Watchers.Files {
		fileWatcher := file.NewWatcher(
			watcherConfig.Path,
			watcherConfig.Name,
			pipeline.Push,
			logger.With(zap.String("component", "watcher"), zap.String("source", watcherConfig.Name)),
		)

		builder.WithWatcher(fileWatcher)
	}

	app, err := builder.Build()

	if err != nil {
		logger.Fatal("Failed to build app", zap.Error(err))
	}

	if err := app.Start(ctx); err != nil {
		logger.Fatal("Application exited with error", zap.Error(err))
	}

	<-ctx.Done()

	logger.Info("Shutdown complete")
}

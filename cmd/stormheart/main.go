package main

import (
	"fmt"

	"github.com/glacius-labs/StormHeart/internal/app"
	"github.com/glacius-labs/StormHeart/internal/application/handler"
	"github.com/glacius-labs/StormHeart/internal/application/shared"
	"github.com/glacius-labs/StormHeart/internal/core/event"
	"github.com/glacius-labs/StormHeart/internal/core/reconciler"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/docker"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/file"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/logging"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/mqtt"
	"github.com/glacius-labs/StormHeart/internal/infrastructure/zap"
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

	logger.Info("Starting StormHeart")

	ctx := setupSignalContext(logger)

	runtime, err := docker.NewRuntime()

	if err != nil {
		logger.Fatal("failed to create runtime", zap.Error(err))
		return
	}

	registry := shared.NewDeploymentsRegistry()

	dispatcher := event.NewDispatcher()

	reconciler := reconciler.NewReconciler(
		runtime,
		dispatcher,
	)

	mqttTopic := fmt.Sprintf("stormfleet/%s/deployments", cfg.Identifier)
	mqttUrl := fmt.Sprintf("tcp://%s:%d", cfg.StormLink.Host, cfg.StormLink.Port)
	mqttClient := mqtt.NewPahoClient(cfg.Identifier, mqttUrl)

	dispatcher.Register(handler.NewDeploymentHandler(registry, reconciler))
	dispatcher.Register(handler.NewWatcherStoppedHandler(registry, reconciler))
	dispatcher.Register(logging.NewLoggingHandler(logger))
	dispatcher.Register(mqtt.NewEventPublisherHandler(mqttClient, fmt.Sprintf("stormfleet/%s/events", cfg.Identifier)))

	mqqtWatcher := mqtt.NewWatcher(
		mqttClient,
		mqttTopic,
		stormLinkSource,
		dispatcher,
	)

	builder := app.
		NewBuilder().
		WithRuntime(runtime).
		WithWatcher(staticWatcher).
		WithWatcher(mqqtWatcher)

	for _, watcherConfig := range cfg.Watchers.Files {
		fileWatcher := file.NewWatcher(
			watcherConfig.Path,
			watcherConfig.Name,
			dispatcher,
		)

		builder.WithWatcher(fileWatcher)
	}

	app, err := builder.Build()

	if err != nil {
		logger.Fatal("failed to build application", zap.Error(err))
		return
	}

	if err := app.Start(ctx); err != nil {
		logger.Fatal("failed to start application", zap.Error(err))
		return
	}

	<-ctx.Done()

	logger.Info("Shutting down StormHeart")
}

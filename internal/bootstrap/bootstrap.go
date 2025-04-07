package bootstrap

import (
	"context"

	"github.com/glacius-labs/StormHeart/internal/config"
	"github.com/glacius-labs/StormHeart/internal/pipeline"
	"github.com/glacius-labs/StormHeart/internal/reconciler"
	"github.com/glacius-labs/StormHeart/internal/runtime"
	"github.com/glacius-labs/StormHeart/internal/watcher"
	"go.uber.org/zap"
)

func Bootstrap(ctx context.Context, cfg config.Config, logger *zap.SugaredLogger) error {
	// Initialize container runtime
	containerRuntime, err := runtime.NewDockerRuntime()
	if err != nil {
		return err
	}

	// Initialize reconciler with scoped logger
	recLogger := logger.With("component", "reconciler")
	rec := reconciler.NewReconciler(containerRuntime, recLogger)

	// Initialize pipeline with scoped logger
	plLogger := logger.With("component", "pipeline")
	pl := pipeline.NewPipeline(
		rec.Apply,
		plLogger,
		pipeline.NewDeduplicator(),
	)

	// Initialize and start watchers
	staticWatcher := watcher.NewStaticWatcher(
		staticDeployments,
		pl.Push,
		logger.With("component", "watcher", "source", "static"),
	)

	go func() {
		if err := staticWatcher.Start(ctx); err != nil {
			logger.Fatalw("Static watcher exited with error", "source", "static", "error", err)
		}
	}()

	for _, fw := range cfg.Watchers.Files {
		watcherLogger := logger.With(
			"component", "watcher",
			"source", fw.Name,
		)

		w := watcher.NewFileWatcher(
			fw.Path,
			fw.Name,
			pl.Push,
			watcherLogger,
		)

		go func(w *watcher.FileWatcher, name string) {
			if err := w.Start(ctx); err != nil {
				logger.Errorw("File watcher exited with error", "source", name, "error", err)
			}
		}(w, fw.Name)
	}

	logger.Infow("System bootstrapped successfully")
	return nil
}

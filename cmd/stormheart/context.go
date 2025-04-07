package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func setupSignalContext(logger *zap.Logger) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
		cancel()
	}()

	return ctx
}

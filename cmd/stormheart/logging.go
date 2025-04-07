package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger() (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.InfoLevel,
	)

	logFile, err := os.OpenFile("stormheart.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(logFile),
		zap.InfoLevel,
	)

	combinedCore := zapcore.NewTee(consoleCore, fileCore)

	logger := zap.New(combinedCore, zap.AddCaller())
	return logger, nil
}

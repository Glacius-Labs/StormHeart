package main

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger(level string) (*zap.SugaredLogger, error) {
	var lvl zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = zap.DebugLevel
	case "info", "":
		lvl = zap.InfoLevel
	case "warn":
		lvl = zap.WarnLevel
	case "error":
		lvl = zap.ErrorLevel
	default:
		return nil, fmt.Errorf("invalid log level: %s", level)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		lvl,
	)

	logFile, err := os.OpenFile("stormheart.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(logFile),
		lvl,
	)

	combinedCore := zapcore.NewTee(consoleCore, fileCore)

	logger := zap.New(combinedCore, zap.AddCaller())
	return logger.Sugar(), nil
}

package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tclgroup/stock-management/internal/pkg/config"
)

// New creates a new zap logger based on the provided LoggerConfig.
func New(cfg config.LoggerConfig) *zap.Logger {
	level := parseLevel(cfg.Level)

	var log *zap.Logger
	var err error

	if cfg.Env == "development" {
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(level)
		log, err = zapCfg.Build()
	} else {
		zapCfg := zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(level)
		log, err = zapCfg.Build()
	}

	if err != nil {
		panic(fmt.Sprintf("logger: %v", err))
	}

	return log
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

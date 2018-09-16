package provider

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger - return a zap sugar logger which will output to paths
func Logger(paths []string) *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = paths
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := cfg.Build()
	return logger.Sugar()
}

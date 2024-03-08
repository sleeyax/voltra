package bot

import (
	"github.com/sleeyax/voltra/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Creates a new sugared logger form the provided logging options.
func createLogger(options config.LoggingOptions) *zap.SugaredLogger {
	var logger *zap.Logger

	if !options.Enable || options.LogLevel == config.SilentLevel {
		logger = zap.NewNop()
	} else if options.EnableStructuredLogging {
		loggerConfig := zap.NewProductionConfig()
		if logLevel := toZapLogLevel(options.LogLevel); logLevel != zapcore.InvalidLevel {
			loggerConfig.Level.SetLevel(logLevel)
		}
		logger, _ = loggerConfig.Build()
	} else {
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		if logLevel := toZapLogLevel(options.LogLevel); logLevel != zapcore.InvalidLevel {
			loggerConfig.Level.SetLevel(logLevel)
		}
		logger, _ = loggerConfig.Build()
	}

	return logger.Sugar()
}

// Converts a config.LogLevel to a zapcore.Level.
func toZapLogLevel(level config.LogLevel) zapcore.Level {
	switch level {
	case config.DebugLevel:
		return zapcore.DebugLevel
	case config.InfoLevel:
		return zapcore.InfoLevel
	case config.WarnLevel:
		return zapcore.WarnLevel
	case config.ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.InvalidLevel
	}
}

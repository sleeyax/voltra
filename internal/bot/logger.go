package bot

import (
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Creates a new sugared logger form the provided logging options.
func createLogger(options config.LoggingOptions) *zap.SugaredLogger {
	var logger *zap.Logger

	if !options.Enable {
		logger = zap.NewNop()
	} else if options.EnableStructuredLogging {
		logger, _ = zap.NewProduction()
	} else {
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = loggerConfig.Build()
	}

	return logger.Sugar()
}

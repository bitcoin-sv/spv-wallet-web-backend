package logging

import (
	"bux-wallet/config"
	"github.com/spf13/viper"
)

// DefaultLoggerFactory creates an instance of logger factory with default configuration for this app.
func DefaultLoggerFactory() LoggerFactory {
	lvl := viper.GetString(config.EnvLoggingLevel)
	return NewZerologLoggerFactory("bux-wallet", LevelFromString(lvl))
}

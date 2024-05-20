package logging

import (
	"io"
	"os"
	"time"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.elastic.co/ecszerolog"
)

const (
	consoleLogFormat = "console"
	jsonLogFormat    = "json"
)

// GetDefaultLogger generates and returns a default logger instance.
func GetDefaultLogger() *zerolog.Logger {
	logger := ecszerolog.New(os.Stdout, ecszerolog.Level(zerolog.DebugLevel)).
		With().
		Caller().
		Str("application", "spv-wallet-web-backend-default").
		Logger()

	return &logger
}

// CreateLogger create and configure zerolog logger based on viper config.
func CreateLogger() (*zerolog.Logger, error) {
	var writer io.Writer
	if viper.GetString(config.EnvLoggingFormat) == consoleLogFormat {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05.000",
		}
	} else {
		writer = os.Stdout
	}

	parsedLevel, err := zerolog.ParseLevel(viper.GetString(config.EnvLoggingLevel))
	if err != nil {
		return nil, err
	}

	instanceName := viper.GetString(config.EnvLoggingInstanceName)
	logLevel := ecszerolog.Level(parsedLevel)
	origin := ecszerolog.Origin()
	var logger zerolog.Logger

	if viper.GetBool(config.EnvLoggingLogOrigin) {
		logger = ecszerolog.New(writer, logLevel, origin).
			With().
			Str("application", instanceName).
			Logger()
	} else {
		logger = ecszerolog.New(writer, logLevel).
			With().
			Str("application", instanceName).
			Logger()
	}

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local) //nolint:gosmopolitan // We want local time inside logger.
	}

	return &logger, nil
}

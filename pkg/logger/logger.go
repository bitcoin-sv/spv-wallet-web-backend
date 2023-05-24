// Package logger provides a logger for the application.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// NewLogger creates a new logger based on the given module name and output.
func NewLogger(moduleName string, output io.Writer) zerolog.Logger {
	output = zerolog.ConsoleWriter{
		Out:        output,
		TimeFormat: "2006-01-02 15:04:05",
	}

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local)
	}

	logger := zerolog.New(output).With().
		Str("name", moduleName).
		Timestamp().
		Logger()

	return logger
}

// NewAppLogger creates a new logger for the app.
func NewAppLogger() zerolog.Logger {
	return NewLogger("app", os.Stdout)
}

// NewServiceLogger creates a new logger for the service.
func NewServiceLogger() zerolog.Logger {
	return NewLogger("service", os.Stdout)
}

// NewRepositoryLogger creates a new logger for the repository.
func NewRepositoryLogger() zerolog.Logger {
	return NewLogger("repository", os.Stdout)
}

// NewTestLogger creates a new logger for testing purposes.
func NewTestLogger() zerolog.Logger {
	return NewLogger("test", os.Stdout)
}

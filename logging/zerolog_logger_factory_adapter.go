package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type loggerFactory struct {
	application string
	writer      io.Writer
}

// NewZerologLoggerFactory create and configure zerolog logger.
func NewZerologLoggerFactory(appName string, level Level) LoggerFactory {
	setupZerologGlobals(level)

	return &loggerFactory{
		application: appName,
		writer: zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		},
	}
}

func setupZerologGlobals(level Level) {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local) // nolint:gosmopolitan // We want local time inside logger.
	}

	zerolog.SetGlobalLevel(toZerologLevel(level))
}

func (lf *loggerFactory) NewLogger(name string) Logger {
	log := zerolog.New(lf.writer).
		With().
		Str("application", lf.application).
		Str("name", name).
		Timestamp().
		Logger()

	return (*logger)(&log)
}

func (lf *loggerFactory) Level() Level {
	return fromZerologLevel(zerolog.GlobalLevel())
}

func (lf *loggerFactory) SetLevel(level Level) {
	zerolog.SetGlobalLevel(toZerologLevel(level))
}

func fromZerologLevel(level zerolog.Level) (lvl Level) {
	switch level {
	case zerolog.TraceLevel:
		lvl = Trace
	case zerolog.DebugLevel:
		lvl = Debug
	case zerolog.InfoLevel:
		lvl = Info
	case zerolog.WarnLevel:
		lvl = Warn
	case zerolog.ErrorLevel:
		lvl = Error
	case zerolog.FatalLevel:
		lvl = Critical
	case zerolog.Disabled:
		lvl = Off
	case zerolog.NoLevel:
		lvl = Trace
	case zerolog.PanicLevel:
		lvl = Critical
	}
	return lvl
}

func toZerologLevel(level Level) (lvl zerolog.Level) {
	switch level {
	case Trace:
		lvl = zerolog.TraceLevel
	case Debug:
		lvl = zerolog.DebugLevel
	case Info:
		lvl = zerolog.InfoLevel
	case Warn:
		lvl = zerolog.WarnLevel
	case Error:
		lvl = zerolog.ErrorLevel
	case Critical:
		lvl = zerolog.FatalLevel
	case Off:
		lvl = zerolog.Disabled
	}
	return lvl
}

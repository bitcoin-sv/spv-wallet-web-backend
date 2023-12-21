package websocket

import (
	"github.com/centrifugal/centrifuge"
	"github.com/rs/zerolog"
)

type logHandler struct {
	logger *zerolog.Logger
}

func newLogHandler(l *zerolog.Logger) *logHandler {
	centrifugeLogger := l.With().Str("subservice", "centrifuge").Logger()
	return &logHandler{
		logger: &centrifugeLogger,
	}
}

func (l *logHandler) Level() (level centrifuge.LogLevel) {
	// nolint:exhaustive // we don't need to handle all cases
	switch l.logger.GetLevel() {
	case zerolog.TraceLevel:
		level = centrifuge.LogLevelTrace
	case zerolog.DebugLevel:
		level = centrifuge.LogLevelDebug
	case zerolog.InfoLevel:
		level = centrifuge.LogLevelInfo
	case zerolog.WarnLevel:
		level = centrifuge.LogLevelWarn
	case zerolog.ErrorLevel:
		level = centrifuge.LogLevelError
	default:
		level = centrifuge.LogLevelNone
	}
	return level
}

func (l *logHandler) Log(entry centrifuge.LogEntry) {
	var event *zerolog.Event
	switch entry.Level {
	case centrifuge.LogLevelTrace:
		event = l.logger.Trace()
	case centrifuge.LogLevelDebug:
		event = l.logger.Debug()
	case centrifuge.LogLevelInfo:
		event = l.logger.Info()
	case centrifuge.LogLevelWarn:
		event = l.logger.Warn()
	case centrifuge.LogLevelError:
		event = l.logger.Error()
	case centrifuge.LogLevelNone:
		event = nil
	}

	if event == nil {
		return
	}

	event.Msgf("%s Context: %v.", entry.Message, entry.Fields)
}

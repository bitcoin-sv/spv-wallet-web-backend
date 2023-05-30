package logging

import "github.com/rs/zerolog"

type logger zerolog.Logger

func (l *logger) Trace(msg string) {
	log := (*zerolog.Logger)(l)
	log.Trace().Msg(msg)
}

func (l *logger) Debug(msg string) {
	log := (*zerolog.Logger)(l)
	log.Debug().Msg(msg)
}

func (l *logger) Info(msg string) {
	log := (*zerolog.Logger)(l)
	log.Info().Msg(msg)
}

func (l *logger) Warn(msg string) {
	log := (*zerolog.Logger)(l)
	log.Warn().Msg(msg)
}

func (l *logger) Error(msg string) {
	log := (*zerolog.Logger)(l)
	log.Error().Msg(msg)
}

func (l *logger) Critical(msg string) {
	log := (*zerolog.Logger)(l)
	log.Fatal().Msg(msg)
}

func (l *logger) Tracef(msg string, v ...interface{}) {
	log := (*zerolog.Logger)(l)
	log.Trace().Msgf(msg, v...)
}

func (l *logger) Debugf(msg string, v ...interface{}) {
	log := (*zerolog.Logger)(l)
	log.Debug().Msgf(msg, v...)
}

func (l *logger) Infof(msg string, v ...interface{}) {
	log := (*zerolog.Logger)(l)
	log.Info().Msgf(msg, v...)
}

func (l *logger) Warnf(msg string, v ...interface{}) {
	log := (*zerolog.Logger)(l)
	log.Warn().Msgf(msg, v...)
}

func (l *logger) Errorf(msg string, v ...interface{}) {
	log := (*zerolog.Logger)(l)
	log.Error().Msgf(msg, v...)
}

func (l *logger) Criticalf(msg string, v ...interface{}) {
	log := (*zerolog.Logger)(l)
	log.Fatal().Msgf(msg, v...)
}

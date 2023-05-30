package logging

import "strings"

// LevelFromString parses logging level from provided text.
func LevelFromString(level string) (l Level) {
	switch strings.ToLower(level) {
	case "trace":
		l = Trace
	case "debug":
		l = Debug
	case "info":
		l = Info
	case "warn":
		l = Warn
	case "error":
		l = Error
	case "critical":
		l = Critical
	case "off":
		l = Off
	}
	return l
}

func (l Level) String() (s string) {
	switch l {
	case Trace:
		s = "Trace"
	case Debug:
		s = "Debug"
	case Info:
		s = "Info"
	case Warn:
		s = "Warn"
	case Error:
		s = "Error"
	case Critical:
		s = "Critical"
	case Off:
		s = "Off"
	}
	return s
}

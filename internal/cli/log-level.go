package cli

import (
	"errors"
)

// LogLevel wraps a log level CLI setting.
type LogLevel string

const (
	// DebugLogLevel defines the "debug" log level CLI setting.
	DebugLogLevel LogLevel = "debug"

	// InfoLogLevel defines the "info" log level CLI setting.
	InfoLogLevel LogLevel = "info"

	// WarnLogLevel defines the "warn" log level CLI setting.
	WarnLogLevel LogLevel = "warn"

	// ErrorLogLevel defines the "error" log level CLI setting.
	ErrorLogLevel LogLevel = "error"
)

// String returns the wrapped log level.
func (ll *LogLevel) String() string {
	return string(*ll)
}

// Set sets the wrapped log level.
func (ll *LogLevel) Set(s string) error {
	switch LogLevel(s) {
	case DebugLogLevel, InfoLogLevel, WarnLogLevel, ErrorLogLevel:
		// Valid
		*ll = LogLevel(s)
		return nil

	default:
		// Invalid
		return errors.New(`must be one of "debug", "info", "warn", or "error"`)
	}
}

// Type returns the name of of this CLI setting.
func (*LogLevel) Type() string {
	return "LogLevel"
}

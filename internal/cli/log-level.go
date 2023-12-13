package cli

import (
	"errors"
	"log/slog"
	"strings"
)

// LogLevel wraps a slog level..
type LogLevel struct {
	Val slog.Level
}

// String returns the wrapped slog level.
func (ll LogLevel) String() string {
	switch ll.Val {
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return "error"
	}

	return ""
}

// Set sets the wrapped slog level.
func (ll *LogLevel) Set(s string) error {
	switch strings.ToLower(s) {
	case "debug":
		ll.Val = slog.LevelDebug
		return nil
	case "info":
		ll.Val = slog.LevelInfo
		return nil
	case "warn":
		ll.Val = slog.LevelWarn
		return nil
	case "error":
		ll.Val = slog.LevelError
		return nil
	}

	// Invalid
	return errors.New(`must be one of "debug", "info", "warn", or "error"`)
}

// Type returns the name of the slog level type.
func (*LogLevel) Type() string {
	return "LogLevel"
}

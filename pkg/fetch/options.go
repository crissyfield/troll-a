package fetch

import (
	"time"
)

// settings wraps all fetching settings.
type settings struct {
	timeout time.Duration
}

// Option is a function mutating the settings.
type Option func(*settings)

// WithTimeout with set the timeout duration for the fetch operation.
func WithTimeout(timeout time.Duration) func(*settings) {
	return func(s *settings) {
		s.timeout = timeout
	}
}

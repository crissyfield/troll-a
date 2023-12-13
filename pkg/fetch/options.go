package fetch

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// settings wraps all fetching settings.
type settings struct {
	timeout time.Duration
	backOff backoff.BackOff
}

// Option is a function mutating the settings.
type Option func(*settings)

// WithTimeout will set the timeout duration for the fetch operation.
func WithTimeout(timeout time.Duration) func(*settings) {
	return func(s *settings) {
		s.timeout = timeout
	}
}

// WithBackoff will set the backoff strategy for the fetch operation.
func WithBackoff(backOff backoff.BackOff) func(*settings) {
	return func(s *settings) {
		s.backOff = backOff
	}
}

package fetch

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// params wraps all fetching parameters.
type params struct {
	timeout time.Duration
	backOff backoff.BackOff
}

// Option is an option for opening a URL.
type Option func(*params)

// WithTimeout will set the timeout duration for the fetch operation.
func WithTimeout(timeout time.Duration) Option {
	return func(s *params) {
		s.timeout = timeout
	}
}

// WithBackoff will set the backoff strategy for the fetch operation.
func WithBackoff(backOff backoff.BackOff) Option {
	return func(s *params) {
		s.backOff = backOff
	}
}

package fetch

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// state wraps all fetching state.
type state struct {
	timeout time.Duration
	backOff backoff.BackOff
}

// WithTimeout will set the timeout duration for the fetch operation.
func WithTimeout(timeout time.Duration) func(*state) {
	return func(s *state) {
		s.timeout = timeout
	}
}

// WithBackoff will set the backoff strategy for the fetch operation.
func WithBackoff(backOff backoff.BackOff) func(*state) {
	return func(s *state) {
		s.backOff = backOff
	}
}

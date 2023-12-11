package fetch

import (
	"time"
)

// config wraps all settings in a struct.
type config struct {
	timeout time.Duration
}

// Option is a function mutating the config.
type Option func(*config)

// WithTimeout with set the timeout duration for the fetch operation.
func WithTimeout(timeout time.Duration) func(*config) {
	return func(c *config) {
		c.timeout = timeout
	}
}

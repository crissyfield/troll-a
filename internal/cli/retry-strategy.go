package cli

import (
	"errors"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var (
	// BackoffStrategyValNone ...
	RetryStrategyValNever = &backoff.StopBackOff{}

	// RetryStrategyValConstant ...
	RetryStrategyValConstant = backoff.WithMaxRetries(&backoff.ConstantBackOff{Interval: 5 * time.Second}, 5)

	// RetryStrategyValExponential ...
	RetryStrategyValExponential = backoff.NewExponentialBackOff()

	// RetryStrategyValAlways ...
	RetryStrategyValAlways = &backoff.ZeroBackOff{}
)

// RetryStrategy wraps a retry strategy.
type RetryStrategy struct {
	Val backoff.BackOff
}

// String returns the wrapped retry strategy.
func (bo RetryStrategy) String() string {
	switch bo.Val {
	case RetryStrategyValNever:
		return "never"
	case RetryStrategyValConstant:
		return "constant"
	case RetryStrategyValExponential:
		return "exponential"
	case RetryStrategyValAlways:
		return "always"
	}

	return ""
}

// Set sets the wrapped retry strategy.
func (bo *RetryStrategy) Set(s string) error {
	switch strings.ToLower(s) {
	case "never":
		bo.Val = RetryStrategyValNever
		return nil

	case "constant":
		bo.Val = RetryStrategyValConstant
		return nil

	case "exponential":
		bo.Val = RetryStrategyValExponential
		return nil

	case "always":
		bo.Val = RetryStrategyValAlways
		return nil
	}

	// Invalid
	return errors.New(`must be one of "never", "constant", "exponential", or "always"`)
}

// Type returns the name of the backoff retry type.
func (*RetryStrategy) Type() string {
	return "RetryStrategy"
}

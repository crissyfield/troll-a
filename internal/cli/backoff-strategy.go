package cli

import (
	"errors"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var (
	// BackoffStrategyValNone ...
	BackoffStrategyValNone = &backoff.StopBackOff{}

	// BackoffStrategyValConstant ...
	BackoffStrategyValConstant = backoff.WithMaxRetries(&backoff.ConstantBackOff{Interval: 5 * time.Second}, 5)

	// BackoffStrategyValExponential ...
	BackoffStrategyValExponential = backoff.NewExponentialBackOff()

	// BackoffStrategyValZero ...
	BackoffStrategyValZero = &backoff.ZeroBackOff{}
)

// BackoffStrategy wraps a backoff strategy.
type BackoffStrategy struct {
	Val backoff.BackOff
}

// String returns the wrapped backoff strategy.
func (bo BackoffStrategy) String() string {
	switch bo.Val {
	case BackoffStrategyValNone:
		return "none"
	case BackoffStrategyValConstant:
		return "constant"
	case BackoffStrategyValExponential:
		return "exponential"
	case BackoffStrategyValZero:
		return "zero"
	}

	return ""
}

// Set sets the wrapped backoff strategy.
func (bo *BackoffStrategy) Set(s string) error {
	switch strings.ToLower(s) {
	case "none":
		bo.Val = BackoffStrategyValNone
		return nil
	case "constant":
		bo.Val = BackoffStrategyValConstant
		return nil
	case "exponential":
		bo.Val = BackoffStrategyValExponential
		return nil
	case "zero":
		bo.Val = BackoffStrategyValZero
		return nil
	}

	// Invalid
	return errors.New(`must be one of "none", "constant", "exponential", or "zero"`)
}

// Type returns the name of the backoff strategy type.
func (*BackoffStrategy) Type() string {
	return "BackoffStrategy"
}

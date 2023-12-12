package cli

import (
	"errors"
)

// RulesPreset wraps a rules preset CLI setting.
type RulesPreset string

const (
	// AllRulesPreset defines the "all" rules preset CLI setting.
	AllRulesPreset RulesPreset = "all"

	// MostRulesPreset defines the "most" rules preset CLI setting.
	MostRulesPreset RulesPreset = "most"

	// SecretRulesPreset defines the "secret" rules preset CLI setting.
	SecretRulesPreset RulesPreset = "secret"
)

// String returns the wrapped rules preset.
func (rp *RulesPreset) String() string {
	return string(*rp)
}

// Set sets the wrapped rules preset.
func (rp *RulesPreset) Set(s string) error {
	switch RulesPreset(s) {
	case AllRulesPreset, MostRulesPreset, SecretRulesPreset:
		// Valid
		*rp = RulesPreset(s)
		return nil

	default:
		// Invalid
		return errors.New(`must be one of "all", "most", or "secret"`)
	}
}

// Type returns the name of of this CLI setting.
func (*RulesPreset) Type() string {
	return "RulesPreset"
}

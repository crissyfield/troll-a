package cli

import (
	"errors"
	"strings"

	"github.com/crissyfield/troll-a/pkg/detect"
	"github.com/crissyfield/troll-a/pkg/detect/preset"
)

// RulesPreset wraps a rules preset.
type RulesPreset struct {
	Val []detect.GitleaksRuleFunction
}

// String returns the wrapped rules preset.
func (rp RulesPreset) String() string {
	// Hack: we can compare slices, so we compare sizes
	switch len(rp.Val) {
	case len(preset.All):
		return "all"
	case len(preset.Most):
		return "most"
	case len(preset.Secret):
		return "secret"
	}

	return ""
}

// Set sets the wrapped rules preset.
func (rp *RulesPreset) Set(s string) error {
	switch strings.ToLower(s) {
	case "all":
		rp.Val = preset.All
		return nil
	case "most":
		rp.Val = preset.Most
		return nil
	case "secret":
		rp.Val = preset.Secret
		return nil
	}

	// Invalid
	return errors.New(`must be one of "all", "most", or "secret"`)
}

// Type returns the name of the rules preset type.
func (*RulesPreset) Type() string {
	return "rules-preset"
}

package detect

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zricethezav/gitleaks/v8/config"
)

// GitleaksRuleFunction is a function that generates a rule.
type GitleaksRuleFunction func() *config.Rule

// Rule contains information that define details on how to detect secrets.
type Rule struct {
	RuleID      string             // Unique identifier for this rule
	Description string             // Description of the rule
	Entropy     float64            // Minimum Shannon entropy a regex group must have to be considered a secret
	SecretGroup int                // Used to extract secret from regex match
	Regex       AbstractRegexp     // Used to detect secrets
	Allowlists  []config.Allowlist // Allows a rule to be ignored for specific regexes, paths, and/or commits
}

// NewRuleFromGitleaksRule creates a cleaned-up Rule object from a Gitlab's rule.
func NewRuleFromGitleaksRule(r *config.Rule) *Rule {
	return &Rule{
		RuleID:      strings.ToLower(strings.NewReplacer(" ", "-", "_", "-").Replace(r.RuleID)),
		Description: r.Description,
		Entropy:     r.Entropy,
		SecretGroup: r.SecretGroup,
		Regex:       CloneRegexp(r.Regex),
		Allowlists:  r.Allowlists,
	}
}

// NewRuleFromRegExp creates a Rule object from a regular expression.
func NewRuleFromRegExp(re *regexp.Regexp, idx int) *Rule {
	return &Rule{
		RuleID:      fmt.Sprintf("custom-rule-%d", idx),
		Description: fmt.Sprintf("Custom rule #%d (as specified via command line)", idx),
		Regex:       CloneRegexp(re),
	}
}

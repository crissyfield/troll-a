package detect

import (
	"github.com/zricethezav/gitleaks/v8/config"
)

// Rule ...
type Rule struct {
	Description string           // Description of the rule
	RuleID      string           // Unique identifier for this rule
	Entropy     float64          // Minimum Shannon entropy a regex group must have to be considered a secret
	SecretGroup int              // Used to extract secret from regex match
	Regex       AbstractRegexp   // Used to detect secrets
	Allowlist   config.Allowlist // Allows a rule to be ignored for specific regexes, paths, and/or commits
}

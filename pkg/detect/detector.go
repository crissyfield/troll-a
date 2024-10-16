package detect

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/rs/zerolog"
)

// Detector wraps a set of rules to detect secrets.
type Detector struct {
	rules         []*Rule
	combinedRegep AbstractRegexp
	enclosed      bool
}

// NewDetector creates a new Detector object with rules from the given set of Gitleaks rule functions and
// additional custom rules.
func NewDetector(ruleFns []GitleaksRuleFunction, customs []string, enclosed bool) (*Detector, error) {
	// Shut up Gitleaks trace logs
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Create rules and extract raw expressions
	rules := make([]*Rule, 0, len(ruleFns)+len(customs))
	exprs := make([]string, 0, len(ruleFns)+len(customs))

	for _, fn := range ruleFns {
		// Add Gitleaks rules and extract raw expressions
		r := fn()

		rules = append(rules, NewRuleFromGitleaksRule(r))
		exprs = append(exprs, r.Regex.String())
	}

	for i, c := range customs {
		// Add custom rules and extract raw expressions
		re, err := regexp.Compile(c)
		if err != nil {
			return nil, fmt.Errorf("unable to compile custom rule as regular expression [%s]: %w", c, err)
		}

		rules = append(rules, NewRuleFromRegExp(re, i))
		exprs = append(exprs, re.String())
	}

	// Return detector
	return &Detector{
		rules:         rules,
		combinedRegep: MustCompileRegexp(strings.Join(exprs, "|")),
		enclosed:      enclosed,
	}, nil
}

// state wraps some internal state for the detection.
type state struct {
	raw     string   // The string to detect secrets for.
	locator *Locator // A locator to turn indexes into lines and columns.
}

// Detect will detect all secrets in the given reader stream.
func (d *Detector) Detect(r io.Reader) ([]*Finding, error) {
	// Turn the reader into a string
	var s state

	if buf, ok := r.(*bytes.Buffer); ok {
		// Extract underlying data from bytes.Buffer
		s.raw = buf.String()
	} else {
		// Read all data from reader
		d, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("read all from buffer: %w", err)
		}

		s.raw = string(d)
	}

	// Check if any of the rules regexp's matches
	if !d.combinedRegep.MatchString(s.raw) {
		return nil, nil
	}

	// Run through all detection rules and gather findings
	var findings []*Finding

	for _, r := range d.rules {
		findings = append(findings, d.detectRule(&s, r)...)
	}

	return findings, nil
}

// detectRule will detect a single rule.
func (d *Detector) detectRule(s *state, r *Rule) []*Finding {
	// Find all strings matching the rule's regular expression
	var findings []*Finding

	for _, idx := range r.Regex.FindAllStringIndex(s.raw, -1) {
		// Extract match and secret
		start, end := idx[0], idx[1]

		match := strings.Trim(s.raw[start:end], "\n")
		secret := extractSecretFromRegexpSubmatch(r, match)

		// Determine location of the match
		if s.locator == nil {
			s.locator = NewLocator(s.raw)
		}

		loc := s.locator.Find(start, end)

		// Traverse allow lists
		var skip bool

		for _, al := range r.Allowlists {
			// Skip, if the secret is in the list of stopwords
			if al.ContainsStopWord(secret) {
				skip = true
				break
			}

			// Skip, if the secret, match, or line is allowed
			if al.RegexTarget == "match" {
				// Check for match
				if al.RegexAllowed(match) {
					skip = true
					break
				}
			} else if al.RegexTarget == "line" {
				// Check for line
				if al.RegexAllowed(s.raw[loc.StartLineIdx:loc.EndLineIdx]) {
					skip = true
					break
				}
			} else {
				// Check for secret
				if al.RegexAllowed(secret) {
					skip = true
					break
				}
			}
		}

		if skip {
			continue
		}

		// Skip on low entropy (if required)
		if (r.Entropy != 0.0) && checkIfLowEntropy(r, secret) {
			continue
		}

		// Skip if not enclosed (if required)
		if d.enclosed && !checkIfEnclosed(loc.Line(s.raw), secret) {
			continue
		}

		// Append finding
		findings = append(findings, &Finding{
			RuleID:      r.RuleID,
			Description: r.Description,
			Secret:      secret,
			Match:       match,
			Location:    loc,
		})
	}

	return findings
}

// extractSecretFromRegexpSubmatch extracts the secret from the given match for the given rule r.
func extractSecretFromRegexpSubmatch(r *Rule, match string) string {
	if r.SecretGroup > 0 {
		// Pick specific secret group
		groups := r.Regex.FindStringSubmatch(match)
		if len(groups) > r.SecretGroup {
			return groups[r.SecretGroup]
		}
	} else {
		// Otherwise, pick second group (if there are only two)
		groups := r.Regex.FindStringSubmatch(match)
		if len(groups) == 2 {
			return groups[1]
		}
	}

	return match
}

// checkIfLowEntropy checks if the entropy of the given secret is too low for the given rule r.
func checkIfLowEntropy(r *Rule, secret string) bool {
	// Compute entropy and bail if too small
	entropy := shannonEntropy(secret)
	if entropy <= r.Entropy {
		return true
	}

	// Hack borrowed from original Gitleaks code
	if strings.HasPrefix(r.RuleID, "generic") {
		// Skip if there is NO digit in the secret
		var containsDigit bool

		for _, r := range secret {
			if (r >= '1') && (r <= '9') {
				containsDigit = true
				break
			}
		}

		if !containsDigit {
			return true
		}
	}

	return false
}

// Check if secret is enclosed in context.
func checkIfEnclosed(context string, secret string) bool {
	// Trim characters that are likely not part of the secret
	trimmedSecret := strings.TrimFunc(secret, isEnclosureDelimiter)

	for len(context) > 0 {
		// Bail if secret does not exist in context
		idx := strings.Index(context, trimmedSecret)
		if idx == -1 {
			break
		}

		// Check character right before match
		prefixRunes := []rune(context[:idx])

		if (len(prefixRunes) > 0) && !isEnclosureDelimiter(prefixRunes[len(prefixRunes)-1]) {
			context = context[idx+len(trimmedSecret):]
			continue
		}

		// Check character right after match
		suffixRunes := []rune(context[idx+len(trimmedSecret):])

		if (len(suffixRunes) > 0) && !isEnclosureDelimiter(suffixRunes[0]) {
			context = context[idx+len(trimmedSecret):]
			continue
		}

		// Secret is enclosed
		return true
	}

	return false
}

// isEnclosureDelimiter return true if the rune r is an enclosure delimiter.
func isEnclosureDelimiter(r rune) bool {
	return unicode.In(
		r,
		unicode.Cc, // Cc is the set of Unicode characters in category Cc (Other, control).
		unicode.Cf, // Cf is the set of Unicode characters in category Cf (Other, format).
		unicode.Co, // Co is the set of Unicode characters in category Co (Other, private use).
		unicode.Cs, // Cs is the set of Unicode characters in category Cs (Other, surrogate).
		unicode.Mc, // Mc is the set of Unicode characters in category Mc (Mark, spacing combining).
		unicode.Me, // Me is the set of Unicode characters in category Me (Mark, enclosing).
		unicode.Mn, // Mn is the set of Unicode characters in category Mn (Mark, nonspacing).
		unicode.Pc, // Pc is the set of Unicode characters in category Pc (Punctuation, connector).
		unicode.Pe, // Pe is the set of Unicode characters in category Pe (Punctuation, close).
		unicode.Pf, // Pf is the set of Unicode characters in category Pf (Punctuation, final quote).
		unicode.Pi, // Pi is the set of Unicode characters in category Pi (Punctuation, initial quote).
		unicode.Po, // Po is the set of Unicode characters in category Po (Punctuation, other).
		unicode.Ps, // Ps is the set of Unicode characters in category Ps (Punctuation, open).
		unicode.Zl, // Zl is the set of Unicode characters in category Zl (Separator, line).
		unicode.Zp, // Zp is the set of Unicode characters in category Zp (Separator, paragraph).
		unicode.Zs, // Zs is the set of Unicode characters in category Zs (Separator, space).
	)
}

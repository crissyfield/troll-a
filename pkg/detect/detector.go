package detect

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/zricethezav/gitleaks/v8/config"
)

// RuleFunction is a function that generates a rule.
type RuleFunction func() *config.Rule

// Detector wraps a set of rules to detect secrets.
type Detector struct {
	rules         []*Rule
	combinedRegep AbstractRegexp
	enclosed      bool
}

// NewDetector creates a new Detector object with rules from the given set of rule functions.
func NewDetector(ruleFns []RuleFunction, enclosed bool) *Detector {
	// Create rules and extract raw expressions
	rules := make([]*Rule, len(ruleFns))
	exprs := make([]string, len(rules))

	for i, fn := range ruleFns {
		r := fn()

		rules[i] = &Rule{
			Description: r.Description,
			RuleID:      r.RuleID,
			Entropy:     r.Entropy,
			SecretGroup: r.SecretGroup,
			Regex:       DuplicateRegexp(r.Regex),
			Allowlist:   r.Allowlist,
		}

		exprs[i] = r.Regex.String()
	}

	// Return detector
	return &Detector{
		rules:         rules,
		combinedRegep: MustCompileRegexp(strings.Join(exprs, "|")),
		enclosed:      enclosed,
	}
}

// Finding wraps all relevant information for a finding.
type Finding struct {
	ID          string    // ID of the rule responsible for the finding.
	Description string    // Description of the secret found.
	Secret      string    // The actual secret.
	Match       string    // The match containing the secret.
	Location    *Location // The location of the match.
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

	idxs := r.Regex.FindAllStringIndex(s.raw, -1)
	for _, idx := range idxs {
		// Extract match
		start, end := idx[0], idx[1]

		match := strings.Trim(s.raw[start:end], "\n")

		// Extract secret from the regexp submatches
		secret := match

		if r.SecretGroup > 0 {
			// Pick specific secret group
			groups := r.Regex.FindStringSubmatch(secret)
			if len(groups) > r.SecretGroup {
				secret = groups[r.SecretGroup]
			}
		} else {
			// Otherwise, pick second group (if there are only two)
			groups := r.Regex.FindStringSubmatch(secret)
			if len(groups) == 2 {
				secret = groups[1]
			}
		}

		// Skip, if the secret is in the list of stopwords
		if r.Allowlist.ContainsStopWord(secret) {
			continue
		}

		// Determine location of the match
		if s.locator == nil {
			// Deferred, as it might be rather slow
			s.locator = NewLocator(s.raw)
		}

		loc := s.locator.Find(start, end)

		// Skip, if the secret, match, or line is allowed
		switch r.Allowlist.RegexTarget {
		case "match":
			// Check for match
			if r.Allowlist.RegexAllowed(match) {
				continue
			}

		case "line":
			// Check for line
			if r.Allowlist.RegexAllowed(s.raw[loc.StartLineIdx:loc.EndLineIdx]) {
				continue
			}

		default:
			// Check for secret
			if r.Allowlist.RegexAllowed(secret) {
				continue
			}
		}

		// Check for entropy
		if r.Entropy != 0.0 {
			// Compute entropy and bail if too small
			entropy := shannonEntropy(secret)
			if entropy <= r.Entropy {
				continue
			}

			// Hack borrowed from original GitLeaks code
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
					continue
				}
			}
		}

		// Check if the secret is enclosed
		if d.enclosed {
			// Trim characters that are likely not part of the secret
			trimmedSecret := strings.TrimFunc(secret, isLikelyNoSecret)

			// TODO: Factor out, check boundaries, ...
			remainingLine := loc.Line(s.raw)
			isEnclosed := false

			for !isEnclosed && (len(remainingLine) > 0) {
				idx := strings.Index(remainingLine, trimmedSecret)
				if idx == -1 {
					break
				}

				if (idx > 0) && !isLikelyNoSecret([]rune(remainingLine[idx-1:])[0]) {
					remainingLine = remainingLine[idx+len(trimmedSecret):]
					continue
				}

				if (idx+len(trimmedSecret) < len(remainingLine)) && !isLikelyNoSecret([]rune(remainingLine[idx+len(trimmedSecret):])[0]) {
					remainingLine = remainingLine[idx+len(trimmedSecret):]
					continue
				}

				isEnclosed = true
			}

			if !isEnclosed {
				continue
			}
		}

		// Append finding
		findings = append(findings, &Finding{
			ID:          r.RuleID,
			Description: r.Description,
			Secret:      secret,
			Match:       match,
			Location:    loc,
		})
	}

	return findings
}

// isLikelyNoSecret return true if the rune most-likely does not belong to a secret.
func isLikelyNoSecret(r rune) bool {
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

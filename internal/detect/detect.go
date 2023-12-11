package detect

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/zricethezav/gitleaks/v8/config"
)

// state wraps some internal state for the detection.
type state struct {
	raw     string   // The string to detect secrets for.
	locator *Locator // A locator to turn indexes into lines and columns.
}

// Finding wraps all relevant information for a finding.
type Finding struct {
	ID          string // ID of the rule responsible for the finding.
	Description string // Description of the secret found.
	Secret      string // The actual secret.
	Match       string // The match containing the secret.
	Line        string // The line containing the match.
	StartLine   int    // Line of the match start.
	StartColumn int    // Column of the match start.
	EndLine     int    // Line of the match end.
	EndColumn   int    // Line of the column end.
}

// Detect will detect all secrets in the given reader stream.
func Detect(r io.Reader) ([]*Finding, error) {
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

	// Run through all detection rules and gather findings
	var findings []*Finding

	for _, r := range detectionRules {
		findings = append(findings, detectRule(&s, r)...)
	}

	return findings, nil
}

// detectRule will detect a single rule.
func detectRule(s *state, r *config.Rule) []*Finding {
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
			if r.Allowlist.RegexAllowed(s.raw[loc.startLineIndex:loc.endLineIndex]) {
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

		// Append finding
		findings = append(findings, &Finding{
			ID:          r.RuleID,
			Description: r.Description,
			Secret:      secret,
			Match:       match,
			Line:        s.raw[loc.startLineIndex:loc.endLineIndex],
			StartLine:   loc.startLine,
			EndLine:     loc.endLine,
			StartColumn: loc.startColumn,
			EndColumn:   loc.endColumn,
		})
	}

	return findings
}

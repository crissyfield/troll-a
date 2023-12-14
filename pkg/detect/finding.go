package detect

// Finding wraps all relevant information for a finding.
type Finding struct {
	RuleID      string    // ID of the rule responsible for the finding.
	Description string    // Description of the secret found.
	Secret      string    // The actual secret.
	Match       string    // The match containing the secret.
	Location    *Location // The location of the match.
}

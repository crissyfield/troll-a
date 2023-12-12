package detect

// CombinedRegexp is an abstract interface for a combined regular expression. It can only be used for matching.
type CombinedRegexp interface {
	// MatchString reports whether the string s contains any match of the regular expression.
	MatchString(s string) bool
}

package detect

// AbstractRegexp is an abstract interface for a regular expression.
type AbstractRegexp interface {
	// MatchString reports whether the string s contains any match of the regular expression.
	MatchString(s string) bool
}

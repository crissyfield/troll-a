package detect

// AbstractRegexp is an abstract interface for a regular expression.
type AbstractRegexp interface {
	// MatchString reports whether the string s contains any match of the regular expression.
	MatchString(s string) bool

	// FindAllStringIndex returns a slice of all successive matches of the expression. A return value of nil
	// indicates no match.
	FindAllStringIndex(s string, n int) [][]int

	// FindStringSubmatch returns a slice of strings holding the text of the leftmost match of the regular
	// expression in s and the matches, if any, of its subexpressions. A return value of nil indicates no
	// match.
	FindStringSubmatch(s string) []string
}

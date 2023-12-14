//go:build !re2_cgo

package detect

import (
	"regexp"
)

// Ensure interface integrity.
var _ AbstractRegexp = (*regexp.Regexp)(nil)

// CompileRegexp parses the given expression and returns, if successful, an AbstractRegexp.
func CompileRegexp(expr string) (AbstractRegexp, error) {
	return regexp.Compile(expr)
}

// MustCompileRegexp is like CompileRegexp but panics if the expression cannot be parsed. It simplifies safe
// initialization of global variables holding compiled abstract regular expressions.
func MustCompileRegexp(expr string) AbstractRegexp {
	return regexp.MustCompile(expr)
}

// CloneRegexp clones the given regular expression re.
func CloneRegexp(re *regexp.Regexp) AbstractRegexp {
	return re
}

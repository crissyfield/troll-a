//go:build re2_cgo

package detect

import (
	"regexp"

	"github.com/wasilibs/go-re2"
)

// Ensure interface integrity.
var _ AbstractRegexp = (*re2.Regexp)(nil)

// CompileRegexp parses the given expression and returns, if successful, an AbstractRegexp.
func CompileRegexp(expr string) (AbstractRegexp, error) {
	return re2.Compile(expr)
}

// MustCompileRegexp is like CompileRegexp but panics if the expression cannot be parsed. It simplifies safe
// initialization of global variables holding compiled abstract regular expressions.
func MustCompileRegexp(expr string) AbstractRegexp {
	return re2.MustCompile(expr)
}

// DuplicateRegexp "duplicates" the given regular expression by recompiling it as go-re2 regular expression.
func DuplicateRegexp(re *regexp.Regexp) AbstractRegexp {
	return re2.MustCompile(re.String())
}

//go:build re2_cgo

package detect

import (
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

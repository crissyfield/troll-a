//go:build !re2_cgo

package detect

import (
	"regexp"
	"strings"
)

// CompileCombinedRegexp combines the given expressions, parses them and returns, if successful, a
// CombinesRegexp that can be used to match against text.
func CompileCombinedRegexp(exprs []string) (CombinedRegexp, error) {
	return regexp.Compile(strings.Join(exprs, "|"))
}

// MustCompileCombinedRegexp is like CompileCombinedRegexp but panics if the expressions cannot be parsed. It
// simplifies safe initialization of global variables holding compiled regular expressions.
func MustCompileCombinedRegexp(exprs []string) CombinedRegexp {
	return regexp.MustCompile(strings.Join(exprs, "|"))
}

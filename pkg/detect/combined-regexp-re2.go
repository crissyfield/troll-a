//go:build re2_cgo

package detect

import (
	"strings"

	"github.com/wasilibs/go-re2"
)

// CompileCombinedRegexp combines the given expressions, parses them and returns, if successful, a
// CombinesRegexp that can be used to match against text.
func CompileCombinedRegexp(exprs []string) (CombinedRegexp, error) {
	return re2.Compile(strings.Join(exprs, "|"))
}

// MustCompileCombinedRegexp is like CompileCombinedRegexp but panics if the expressions cannot be parsed. It
// simplifies safe initialization of global variables holding compiled regular expressions.
func MustCompileCombinedRegexp(exprs []string) CombinedRegexp {
	return re2.MustCompile(strings.Join(exprs, "|"))
}

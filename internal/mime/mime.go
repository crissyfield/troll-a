package mime

import (
	"github.com/gabriel-vasile/mimetype"
)

const (
	// textPlainMimeType specifies the mime type of "text/plain".
	textPlainMimeType = "text/plain"
)

var (
	// isTextCache caches results of calls to IsText.
	isTextCache = make(map[string]bool)
)

// IsText returns true if the given mime is inherited from "text/plain".
func IsText(mime string) bool {
	// Check cache first
	if isText, ok := isTextCache[mime]; ok {
		return isText
	}

	// Get value bypassing cache
	it := isTextNoCache(mime)

	isTextCache[mime] = it
	return it
}

// isTextNoCache returns true if the given mime is inherited from "text/plain". Any cache is ignored.
func isTextNoCache(mime string) bool {
	// Traverse mime tree
	for mt := mimetype.Lookup(mime); mt != nil; mt = mt.Parent() {
		if mt.Is(textPlainMimeType) {
			return true
		}
	}

	return false
}

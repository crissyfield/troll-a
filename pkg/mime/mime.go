package mime

import (
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

const (
	// textMimeType specifies the mime type family of "text".
	textMimeType = "text"

	// textPlainMimeSubtype specifies the mime subtype of "text/plain".
	textPlainMimeSubtype = "text/plain"
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
	// Remove additional information
	parts := strings.SplitN(mime, ";", 2)
	if len(parts) < 1 {
		return false
	}

	mime = strings.ToLower(strings.TrimSpace(parts[0]))

	// Early exit if already text mime type
	if strings.HasPrefix(mime, textMimeType+"/") {
		return true
	}

	// Traverse mime tree
	for mt := mimetype.Lookup(mime); mt != nil; mt = mt.Parent() {
		if mt.Is(textPlainMimeSubtype) {
			return true
		}
	}

	return false
}

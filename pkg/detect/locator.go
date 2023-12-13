package detect

// Locator allows to locate indexes in a string.
type Locator struct {
	newLineIndexes []int
}

// NewLocator creates a Locator object for the given string.
func NewLocator(s string) *Locator {
	// First line starts before the buffer
	newLineIndexes := []int{0}

	// Add indexes of all newlines
	for i, r := range s {
		if r == '\n' {
			newLineIndexes = append(newLineIndexes, i+1)
		}
	}

	// Last line stops at end of string
	newLineIndexes = append(newLineIndexes, len(s)+1)

	return &Locator{newLineIndexes: newLineIndexes}
}

// Find returns the location for a given index pair.
func (l *Locator) Find(startIdx int, endIdx int) *Location {
	// Find start line
	startLine := 0

	for ; startLine < len(l.newLineIndexes)-1; startLine++ {
		if l.newLineIndexes[startLine+1] >= startIdx {
			break
		}
	}

	// Find end line
	endLine := startLine

	for ; endLine < len(l.newLineIndexes)-1; endLine++ {
		if l.newLineIndexes[endLine+1] > endIdx {
			break
		}
	}

	// Return the location
	return &Location{
		StartIdx:     startIdx,
		EndIdx:       endIdx,
		StartLine:    startLine,
		EndLine:      endLine,
		StartColumn:  startIdx - l.newLineIndexes[startLine],
		EndColumn:    endIdx - l.newLineIndexes[endLine],
		StartLineIdx: l.newLineIndexes[startLine],
		EndLineIdx:   l.newLineIndexes[endLine+1] - 1,
	}
}

package detect

// Locator allows to locate indexes in a string.
type Locator struct {
	newLines []int
}

// NewLocator creates a Locator object for the given string.
func NewLocator(raw string) *Locator {
	// First line starts before the buffer
	newLines := []int{0}

	// Add indexes of all newlines
	for i, r := range raw {
		if r == '\n' {
			newLines = append(newLines, i+1)
		}
	}

	// Last line stops at end of string
	newLines = append(newLines, len(raw)+1)

	return &Locator{newLines: newLines}
}

// Location represents a location in a string.
type Location struct {
	startLine      int
	endLine        int
	startColumn    int
	endColumn      int
	startLineIndex int
	endLineIndex   int
}

// Find returns the location for a given index pair.
func (l *Locator) Find(start int, end int) *Location {
	// Find start line
	startLine := 0

	for ; startLine < len(l.newLines)-1; startLine++ {
		if l.newLines[startLine+1] >= start {
			break
		}
	}

	// Find end line
	endLine := startLine

	for ; endLine < len(l.newLines)-1; endLine++ {
		if l.newLines[endLine+1] > end {
			break
		}
	}

	// Return the location
	return &Location{
		startLine:      startLine,
		endLine:        endLine,
		startColumn:    start - l.newLines[startLine],
		endColumn:      end - l.newLines[endLine],
		startLineIndex: l.newLines[startLine],
		endLineIndex:   l.newLines[endLine+1] - 1,
	}
}

package detect

// Location represents a location in a string.
type Location struct {
	StartIdx     int // Start index of the match
	EndIdx       int // End index of the match
	StartLine    int // Text line the start index falls in
	EndLine      int // Text line the end index falls in
	StartColumn  int // Column in the start line corresponding to the start index
	EndColumn    int // Column in the end line corresponding to the end index
	StartLineIdx int // Index of the beginning of the start line
	EndLineIdx   int // Index of the first character after the end line
}

// Line returns the line(s) of the location within s.
func (l *Location) Line(s string) string {
	return s[l.StartLineIdx:l.EndLineIdx]
}

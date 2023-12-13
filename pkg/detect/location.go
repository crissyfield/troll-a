package detect

// Location represents a location in a string.
type Location struct {
	StartIdx     int
	EndIdx       int
	StartLine    int
	EndLine      int
	StartColumn  int
	EndColumn    int
	StartLineIdx int
	EndLineIdx   int
}

// Line returns the line of the location withing s.
func (l *Location) Line(s string) string {
	return s[l.StartLineIdx:l.EndLineIdx]
}

package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	// InfoStyle defines the style used for informational output.
	InfoStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(11))

	// SuccessStyle defines the style used for success output.
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(10))

	// ErrorStyle defines the style used for error output.
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(9))
)

// Info outputs an informational message to STDOUT.
func Info(format string, a ...any) {
	fmt.Fprintln(os.Stdout, InfoStyle.Render(fmt.Sprintf(format, a...)))
}

// Success outputs a success message to STDERR.
func Success(format string, a ...any) {
	fmt.Fprintln(os.Stderr, SuccessStyle.Render(fmt.Sprintf(format, a...)))
}

// Error outputs an error message to STDERR.
func Error(format string, a ...any) {
	fmt.Fprintln(os.Stderr, ErrorStyle.Render(fmt.Sprintf(format, a...)))
}

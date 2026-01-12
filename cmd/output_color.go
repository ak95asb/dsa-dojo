package cmd

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// Color definitions
var (
	ColorGreen  = color.New(color.FgGreen)
	ColorYellow = color.New(color.FgYellow)
	ColorRed    = color.New(color.FgRed)
	ColorBold   = color.New(color.Bold)
)

// isTerminal checks if file descriptor is a terminal
func isTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

// shouldUseColors checks if colored output should be used
func shouldUseColors() bool {
	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check if stdout is a terminal
	return isTerminal(os.Stdout)
}

// colorize applies color to text if colors enabled
func colorize(text string, c *color.Color) string {
	if !shouldUseColors() {
		return text
	}
	return c.Sprint(text)
}

// colorDifficulty applies difficulty-specific color
func colorDifficulty(difficulty string) string {
	if !shouldUseColors() {
		return difficulty
	}

	switch strings.ToLower(difficulty) {
	case "easy":
		return ColorGreen.Sprint(difficulty)
	case "medium":
		return ColorYellow.Sprint(difficulty)
	case "hard":
		return ColorRed.Sprint(difficulty)
	default:
		return difficulty
	}
}

// colorStatus applies color to status indicator
func colorStatus(solved bool) string {
	if solved {
		return colorize("✓", ColorGreen)
	}
	return colorize("✗", ColorYellow)
}

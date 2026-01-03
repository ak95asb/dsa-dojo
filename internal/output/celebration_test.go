package output

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatCelebration(t *testing.T) {
	t.Run("formats celebration with color", func(t *testing.T) {
		os.Unsetenv("NO_COLOR")

		output := FormatCelebration("Two Sum", 3)

		assert.Contains(t, output, "Congratulations! You solved Two Sum!")
		assert.Contains(t, output, "Solved in 3 attempts")
		assert.Contains(t, output, "dsa status")
		assert.Contains(t, output, "Keep up the great work!")
		assert.Contains(t, output, "\033[") // ANSI color codes present
		assert.Contains(t, output, "\U0001F389") // Party popper emoji
		assert.Contains(t, output, "\U0001F680") // Rocket emoji
	})

	t.Run("formats celebration without color when NO_COLOR set", func(t *testing.T) {
		os.Setenv("NO_COLOR", "1")
		defer os.Unsetenv("NO_COLOR")

		output := FormatCelebration("Add Two Numbers", 1)

		assert.Contains(t, output, "Congratulations! You solved Add Two Numbers!")
		assert.Contains(t, output, "Solved in 1 attempt")
		assert.Contains(t, output, "dsa status")
		assert.Contains(t, output, "Keep up the great work!")
		assert.NotContains(t, output, "\033[") // No ANSI color codes
		assert.NotContains(t, output, "\U0001F389") // No emoji
		assert.NotContains(t, output, "\U0001F680") // No emoji
	})

	t.Run("uses singular 'attempt' for 1 attempt", func(t *testing.T) {
		os.Unsetenv("NO_COLOR")

		output := FormatCelebration("Reverse String", 1)

		assert.Contains(t, output, "Solved in 1 attempt")
		assert.NotContains(t, output, "1 attempts")
	})

	t.Run("uses plural 'attempts' for multiple attempts", func(t *testing.T) {
		os.Unsetenv("NO_COLOR")

		output := FormatCelebration("Valid Palindrome", 5)

		assert.Contains(t, output, "Solved in 5 attempts")
	})

	t.Run("includes problem title in message", func(t *testing.T) {
		os.Unsetenv("NO_COLOR")

		output := FormatCelebration("Longest Substring Without Repeating Characters", 2)

		assert.Contains(t, output, "Longest Substring Without Repeating Characters")
	})

	t.Run("includes status command suggestion", func(t *testing.T) {
		os.Unsetenv("NO_COLOR")

		output := FormatCelebration("Two Sum", 1)

		assert.Contains(t, output, "dsa status")
	})

	t.Run("output is multi-line", func(t *testing.T) {
		os.Unsetenv("NO_COLOR")

		output := FormatCelebration("Two Sum", 1)

		lines := strings.Split(output, "\n")
		assert.GreaterOrEqual(t, len(lines), 4) // At least 4 lines
	})
}

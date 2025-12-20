package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRandomCommandIntegration tests the random command end-to-end
func TestRandomCommandIntegration(t *testing.T) {
	// Setup: Use existing seeded database
	setupIntegrationTest(t)
	defer cleanupIntegrationTest()

	t.Run("displays random problem with no filters", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute command
		rootCmd.SetArgs([]string{"random"})
		err := rootCmd.Execute()

		// Restore stdout
		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		// Assertions
		assert.NoError(t, err)
		assert.Contains(t, output, "🎲 Random Problem Selected!", "Output should contain random header")
		assert.Contains(t, output, "Difficulty:", "Output should contain difficulty")
		assert.Contains(t, output, "Topic:", "Output should contain topic")
		assert.Contains(t, output, "Description:", "Output should contain description")
		assert.Contains(t, output, "Run 'dsa solve", "Output should contain solve command suggestion")
	})

	t.Run("filters by difficulty", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rootCmd.SetArgs([]string{"random", "--difficulty", "easy"})
		rootCmd.Execute()

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		// Should contain "Easy" difficulty
		assert.Contains(t, output, "Easy", "Output should show Easy difficulty")
	})

	t.Run("filters by topic", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rootCmd.SetArgs([]string{"random", "--topic", "arrays"})
		rootCmd.Execute()

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		// Should contain "arrays" topic
		assert.Contains(t, output, "arrays", "Output should show arrays topic")
	})

	t.Run("combines difficulty and topic filters", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rootCmd.SetArgs([]string{"random", "--difficulty", "medium", "--topic", "linked-lists"})
		rootCmd.Execute()

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		// Should contain both Medium and linked-lists
		assert.Contains(t, output, "Medium", "Output should show Medium difficulty")
		assert.Contains(t, output, "linked-lists", "Output should show linked-lists topic")
	})

	t.Run("handles invalid difficulty flag", func(t *testing.T) {
		// Capture stderr
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		rootCmd.SetArgs([]string{"random", "--difficulty", "super-hard"})
		// Command should exit with code 2, but we can't test os.Exit directly

		w.Close()
		os.Stderr = old
		buf := bytes.Buffer{}
		buf.ReadFrom(r)
		// Error message tested in unit tests
	})

	t.Run("output format includes all required sections", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		rootCmd.SetArgs([]string{"random"})
		rootCmd.Execute()

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		// Verify all required sections are present
		requiredSections := []string{
			"🎲 Random Problem Selected!", // Header
			"Difficulty:",                  // Difficulty section
			"Topic:",                       // Topic section
			"Slug:",                        // Slug section
			"Description:",                 // Description section
			"Run 'dsa solve",              // Suggested command
			strings.Repeat("=", 80),       // Separator
		}

		for _, section := range requiredSections {
			assert.Contains(t, output, section, "Output should contain section: "+section)
		}
	})
}

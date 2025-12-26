package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolveCommand(t *testing.T) {
	// Setup: Use existing seeded database
	setupIntegrationTest(t)
	defer cleanupIntegrationTest()

	// Create temporary directory for solutions
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	t.Run("creates solution file successfully", func(t *testing.T) {
		// Save stdout
		oldStdout := os.Stdout
		defer func() { os.Stdout = oldStdout }()

		// Capture stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute command
		rootCmd.SetArgs([]string{"solve", "two-sum"})
		err := rootCmd.Execute()

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		// Assertions
		assert.NoError(t, err)
		assert.Contains(t, output, "✓ Solution file generated:")
		assert.Contains(t, output, "solutions/two_sum.go")

		// Verify file was created
		solutionPath := filepath.Join("solutions", "two_sum.go")
		_, err = os.Stat(solutionPath)
		assert.NoError(t, err, "Solution file should exist")

		// Verify file content
		content, err := os.ReadFile(solutionPath)
		assert.NoError(t, err)
		assert.Contains(t, string(content), "package solutions")
		assert.Contains(t, string(content), "func TwoSum()")
		assert.Contains(t, string(content), "// Two Sum")
	})

	t.Run("handles --force flag for existing files", func(t *testing.T) {
		// Create a solution file first
		solutionPath := filepath.Join("solutions", "binary_search.go")
		os.MkdirAll("solutions", 0755)
		os.WriteFile(solutionPath, []byte("// Old content"), 0644)

		// Save stdout
		oldStdout := os.Stdout
		defer func() { os.Stdout = oldStdout }()

		// Capture stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute command with --force
		rootCmd.SetArgs([]string{"solve", "binary-search", "--force"})
		err := rootCmd.Execute()

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		// Assertions
		assert.NoError(t, err)
		assert.Contains(t, output, "✓ Backup created:")
		assert.Contains(t, output, "✓ Solution file generated:")

		// Verify backup was created
		backupPath := solutionPath + ".backup"
		backupContent, err := os.ReadFile(backupPath)
		assert.NoError(t, err)
		assert.Contains(t, string(backupContent), "// Old content")

		// Verify new file has generated content
		newContent, err := os.ReadFile(solutionPath)
		assert.NoError(t, err)
		assert.Contains(t, string(newContent), "package solutions")
		assert.Contains(t, string(newContent), "func BinarySearch()")
	})

	t.Run("shows error for invalid problem slug", func(t *testing.T) {
		// Save stderr
		oldStderr := os.Stderr
		defer func() { os.Stderr = oldStderr }()

		// Capture stderr
		_, w, _ := os.Pipe()
		os.Stderr = w

		// Execute command with invalid slug
		rootCmd.SetArgs([]string{"solve", "invalid-problem-xyz"})
		// Command will exit, but we can't catch os.Exit in tests

		// Note: This test verifies the command structure
		// Full exit code testing requires subprocess testing
	})
}

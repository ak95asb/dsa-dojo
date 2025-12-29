package solution

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestSlugToFunctionName(t *testing.T) {
	tests := []struct {
		name string
		slug string
		want string
	}{
		{"simple slug", "two-sum", "TwoSum"},
		{"multiple words", "binary-search-tree", "BinarySearchTree"},
		{"with numbers", "3sum", "3sum"},
		{"single word", "arrays", "Arrays"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slugToFunctionName(tt.slug)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateSolution(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	generator := NewGenerator()
	problem := &database.Problem{
		ID:          1,
		Slug:        "two-sum",
		Title:       "Two Sum",
		Difficulty:  "easy",
		Topic:       "arrays",
		Description: "Find two numbers that add up to target",
	}

	t.Run("creates solution file successfully", func(t *testing.T) {
		filePath, err := generator.GenerateSolution(problem, false)

		assert.NoError(t, err)
		assert.Equal(t, filepath.Join("solutions", "two_sum.go"), filePath)

		// Verify file exists
		_, err = os.Stat(filePath)
		assert.NoError(t, err)

		// Verify file contents
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "package solutions")
		assert.Contains(t, contentStr, "// Two Sum")
		assert.Contains(t, contentStr, "// Difficulty: easy")
		assert.Contains(t, contentStr, "// Topic: arrays")
		assert.Contains(t, contentStr, "func TwoSum()")
		assert.Contains(t, contentStr, "// TODO: Implement your solution here")
	})

	t.Run("creates backup when overwriting existing file", func(t *testing.T) {
		// First create a solution file
		filePath, err := generator.GenerateSolution(problem, false)
		assert.NoError(t, err)

		// Modify the file to ensure backup is different
		originalContent := []byte("// Modified content")
		os.WriteFile(filePath, originalContent, 0644)

		// Overwrite with force flag
		newPath, err := generator.GenerateSolution(problem, true)
		assert.NoError(t, err)
		assert.Equal(t, filePath, newPath)

		// Verify backup file exists
		backupPath := filePath + ".backup"
		_, err = os.Stat(backupPath)
		assert.NoError(t, err)

		// Verify backup contains original content
		backupContent, err := os.ReadFile(backupPath)
		assert.NoError(t, err)
		assert.Contains(t, string(backupContent), "// Modified content")

		// Verify new file has generated content
		newContent, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Contains(t, string(newContent), "package solutions")
	})
}

func TestCopyFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "dest.txt")

	// Create source file
	content := []byte("test content")
	err := os.WriteFile(srcPath, content, 0644)
	assert.NoError(t, err)

	// Copy file
	err = copyFile(srcPath, dstPath)
	assert.NoError(t, err)

	// Verify destination file exists and has same content
	dstContent, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.Equal(t, content, dstContent)
}

func TestGenerateSolutionCreatesDirectory(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Ensure solutions directory doesn't exist
	os.RemoveAll("solutions")

	generator := NewGenerator()
	problem := &database.Problem{
		ID:    1,
		Slug:  "test-problem",
		Title: "Test Problem",
	}

	filePath, err := generator.GenerateSolution(problem, false)

	assert.NoError(t, err)
	assert.Contains(t, filePath, "solutions")

	// Verify directory was created
	info, err := os.Stat("solutions")
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

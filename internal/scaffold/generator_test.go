package scaffold

import (
	"os"
	"path/filepath"
	"strings"
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

func TestGenerateBoilerplate(t *testing.T) {
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

	filePath, err := generator.GenerateBoilerplate(problem)

	assert.NoError(t, err)
	assert.Equal(t, filepath.Join("problems", "two_sum.go"), filePath)

	// Verify file exists
	_, err = os.Stat(filePath)
	assert.NoError(t, err)

	// Verify file contents
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "package problems")
	assert.Contains(t, contentStr, "// Two Sum")
	assert.Contains(t, contentStr, "// Difficulty: easy")
	assert.Contains(t, contentStr, "// Topic: arrays")
	assert.Contains(t, contentStr, "func TwoSum()")
	assert.Contains(t, contentStr, "// TODO: Implement your solution here")
}

func TestGenerateTestFile(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	generator := NewGenerator()
	problem := &database.Problem{
		ID:    1,
		Slug:  "two-sum",
		Title: "Two Sum",
	}

	filePath, err := generator.GenerateTestFile(problem)

	assert.NoError(t, err)
	assert.Equal(t, filepath.Join("problems", "two_sum_test.go"), filePath)

	// Verify file exists
	_, err = os.Stat(filePath)
	assert.NoError(t, err)

	// Verify file contents
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "package problems")
	assert.Contains(t, contentStr, "import")
	assert.Contains(t, contentStr, "github.com/stretchr/testify/assert")
	assert.Contains(t, contentStr, "func TestTwoSum(t *testing.T)")
	assert.Contains(t, contentStr, "t.Run(tt.name, func(t *testing.T)")
}

func TestGenerateBoilerplateCreatesDirectory(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Ensure problems directory doesn't exist
	os.RemoveAll("problems")

	generator := NewGenerator()
	problem := &database.Problem{
		ID:    1,
		Slug:  "test-problem",
		Title: "Test Problem",
	}

	filePath, err := generator.GenerateBoilerplate(problem)

	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(filePath, "problems"+string(filepath.Separator)))

	// Verify directory was created
	info, err := os.Stat("problems")
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

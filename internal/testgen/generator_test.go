package testgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_DeriveFunctionName(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		expected string
	}{
		{"simple slug", "two-sum", "TwoSum"},
		{"three words", "add-two-numbers", "AddTwoNumbers"},
		{"single word", "palindrome", "Palindrome"},
		{"hyphenated", "reverse-linked-list", "ReverseLinkedList"},
	}

	gen := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.deriveFunctionName(tt.slug)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerator_GenerateNew(t *testing.T) {
	tmpDir := t.TempDir()

	// Create problems directory within temp dir
	problemsDir := filepath.Join(tmpDir, "problems")
	err := os.Mkdir(problemsDir, 0755)
	assert.NoError(t, err)

	// Change to temp dir for test
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	prob := &problem.ProblemDetails{
		Problem: database.Problem{
			Slug:  "test-problem",
			Title: "Test Problem",
		},
	}

	testCases := []*TestCase{
		{
			Name:     "test case 1",
			Inputs:   []interface{}{1, 2, 3},
			Expected: 6,
		},
		{
			Name:     "test case 2",
			Inputs:   []interface{}{4, 5},
			Expected: 9,
		},
	}

	gen := NewGenerator()
	err = gen.generateNew("problems/test-problem_test.go", prob, testCases)
	assert.NoError(t, err)

	// Verify file was created
	testFilePath := filepath.Join(problemsDir, "test-problem_test.go")
	assert.FileExists(t, testFilePath)

	// Verify file contains expected content
	content, err := os.ReadFile(testFilePath)
	assert.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "package problems")
	assert.Contains(t, contentStr, "func TestTestProblem(t *testing.T)")
	assert.Contains(t, contentStr, `"test case 1"`)
	assert.Contains(t, contentStr, `"test case 2"`)
	assert.Contains(t, contentStr, "assert.Equal")
}

func TestGenerator_AppendToNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create problems directory
	problemsDir := filepath.Join(tmpDir, "problems")
	err := os.Mkdir(problemsDir, 0755)
	assert.NoError(t, err)

	// Change to temp dir
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	prob := &problem.ProblemDetails{
		Problem: database.Problem{
			Slug:  "new-problem",
			Title: "New Problem",
		},
	}

	testCases := []*TestCase{
		{
			Name:     "first test",
			Inputs:   []interface{}{1, 2},
			Expected: 3,
		},
	}

	gen := NewGenerator()
	// Append mode on non-existent file should create new file
	err = gen.Generate(prob, testCases, true)
	assert.NoError(t, err)

	testFilePath := filepath.Join(problemsDir, "new-problem_test.go")
	assert.FileExists(t, testFilePath)
}

func TestFormatValue_DifferentTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"integer", 42, "42"},
		{"string", "hello", `"hello"`},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"float as int", float64(10), "10"},
		{"float with decimal", 3.14, "3.14"},
		{"slice of ints", []interface{}{1, 2, 3}, "[]int{1, 2, 3}"},
		{"empty slice", []interface{}{}, "[]int{}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerator_Generate_InvalidPath(t *testing.T) {
	gen := NewGenerator()
	prob := &problem.ProblemDetails{
		Problem: database.Problem{
			Slug:  "test",
			Title: "Test",
		},
	}
	testCases := []*TestCase{
		{Name: "test", Inputs: []interface{}{}, Expected: 0},
	}

	// Try to write to invalid path
	err := gen.generateNew("/invalid/path/that/does/not/exist/test.go", prob, testCases)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write test file")
}

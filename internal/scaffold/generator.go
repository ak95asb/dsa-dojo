package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/problem"
)

// Generator handles boilerplate and test file generation for custom problems
type Generator struct {
	problemsDir string
}

// NewGenerator creates a new code generator instance
func NewGenerator() *Generator {
	return &Generator{
		problemsDir: "problems",
	}
}

// GenerateBoilerplate creates a boilerplate Go file for the problem
func (g *Generator) GenerateBoilerplate(p *database.Problem) (string, error) {
	// Ensure problems directory exists
	if err := os.MkdirAll(g.problemsDir, 0755); err != nil {
		return "", fmt.Errorf("create problems directory: %w", err)
	}

	// Generate file path
	fileName := problem.SlugToSnakeCase(p.Slug) + ".go"
	filePath := filepath.Join(g.problemsDir, fileName)

	// Generate function name (PascalCase from slug)
	funcName := slugToFunctionName(p.Slug)

	// Prepare template data
	data := struct {
		FunctionName string
		ProblemTitle string
		Description  string
		Difficulty   string
		Topic        string
	}{
		FunctionName: funcName,
		ProblemTitle: p.Title,
		Description:  p.Description,
		Difficulty:   p.Difficulty,
		Topic:        p.Topic,
	}

	// Load and execute template
	tmpl := boilerplateTemplate
	t, err := template.New("boilerplate").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if err := t.Execute(file, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return filePath, nil
}

// GenerateTestFile creates a test file for the problem
func (g *Generator) GenerateTestFile(p *database.Problem) (string, error) {
	// Ensure problems directory exists
	if err := os.MkdirAll(g.problemsDir, 0755); err != nil {
		return "", fmt.Errorf("create problems directory: %w", err)
	}

	// Generate file path
	fileName := problem.SlugToSnakeCase(p.Slug) + "_test.go"
	filePath := filepath.Join(g.problemsDir, fileName)

	// Generate function name (PascalCase from slug)
	funcName := slugToFunctionName(p.Slug)

	// Prepare template data
	data := struct {
		FunctionName string
		ProblemTitle string
	}{
		FunctionName: funcName,
		ProblemTitle: p.Title,
	}

	// Load and execute template
	tmpl := testTemplate
	t, err := template.New("test").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if err := t.Execute(file, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return filePath, nil
}

// slugToFunctionName converts slug to PascalCase function name
// Examples:
//   "two-sum" -> "TwoSum"
//   "binary-search-tree" -> "BinarySearchTree"
func slugToFunctionName(slug string) string {
	parts := strings.Split(slug, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

const boilerplateTemplate = `package problems

// {{.ProblemTitle}}
// Difficulty: {{.Difficulty}}
// Topic: {{.Topic}}
//
// Description:
// {{.Description}}

// {{.FunctionName}} solves the {{.ProblemTitle}} problem
func {{.FunctionName}}() {
	// TODO: Implement your solution here
}
`

const testTemplate = `package problems

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test{{.FunctionName}} tests the {{.ProblemTitle}} solution
func Test{{.FunctionName}}(t *testing.T) {
	tests := []struct {
		name     string
		// Add your test case fields here
		expected interface{}
	}{
		// Add your test cases here
		{
			name:     "example test case",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Call your function and assert results
			// result := {{.FunctionName}}(...)
			// assert.Equal(t, tt.expected, result)
			assert.True(t, true, "Replace with actual test")
		})
	}
}
`

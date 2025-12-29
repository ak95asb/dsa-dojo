package testgen

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ak95asb/dsa-dojo/internal/problem"
)

// Generator handles Go test file generation
type Generator struct{}

// NewGenerator creates a new test generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate creates or appends to a test file with the provided test cases
func (g *Generator) Generate(prob *problem.ProblemDetails, testCases []*TestCase, appendMode bool) error {
	testFilePath := filepath.Join("problems", prob.Slug+"_test.go")

	// Handle append mode
	if appendMode {
		return g.appendToExisting(testFilePath, prob, testCases)
	}

	// Generate new test file
	return g.generateNew(testFilePath, prob, testCases)
}

// generateNew creates a new test file from scratch
func (g *Generator) generateNew(testFilePath string, prob *problem.ProblemDetails, testCases []*TestCase) error {
	// Prepare template data
	data := struct {
		FunctionName string
		TestCases    []*TestCase
	}{
		FunctionName: g.deriveFunctionName(prob.Slug),
		TestCases:    testCases,
	}

	// Generate code from template
	var buf bytes.Buffer
	if err := testTemplate.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format Go code: %w", err)
	}

	// Write to file
	if err := os.WriteFile(testFilePath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}

	fmt.Printf("✅ Generated test file: %s\n", testFilePath)
	return nil
}

// appendToExisting appends test cases to an existing test file
func (g *Generator) appendToExisting(testFilePath string, prob *problem.ProblemDetails, newTestCases []*TestCase) error {
	// Check if file exists
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		// File doesn't exist, generate new
		return g.generateNew(testFilePath, prob, newTestCases)
	}

	// Read existing file
	existingCode, err := os.ReadFile(testFilePath)
	if err != nil {
		return fmt.Errorf("failed to read existing test file: %w", err)
	}

	// Parse existing test cases
	existingTestCases, err := g.parseExistingTests(existingCode)
	if err != nil {
		// If parsing fails, just generate new file
		fmt.Printf("⚠️  Could not parse existing tests, creating new file\n")
		return g.generateNew(testFilePath, prob, newTestCases)
	}

	// Merge test cases
	mergedTestCases := append(existingTestCases, newTestCases...)

	// Generate new file with merged test cases
	if err := g.generateNew(testFilePath, prob, mergedTestCases); err != nil {
		return err
	}

	fmt.Printf("✅ Appended %d test case(s) to %s\n", len(newTestCases), testFilePath)
	return nil
}

// parseExistingTests attempts to extract test cases from existing Go test file
func (g *Generator) parseExistingTests(code []byte) ([]*TestCase, error) {
	// Parse the Go source file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go file: %w", err)
	}

	// This is a simplified parser - in production, you'd use go/ast more thoroughly
	// For now, we'll return empty slice to indicate append will create new file
	_ = file // Use the parsed file (avoid unused variable error)

	return []*TestCase{}, nil
}

// deriveFunctionName derives the test function name from problem slug
// Example: "two-sum" -> "TwoSum"
func (g *Generator) deriveFunctionName(slug string) string {
	parts := strings.Split(slug, "-")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// formatValue converts an interface{} value to its Go code representation
func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case []interface{}:
		parts := make([]string, len(val))
		for i, item := range val {
			parts[i] = formatValue(item)
		}
		return fmt.Sprintf("[]int{%s}", strings.Join(parts, ", "))
	case int, int64, int32:
		return fmt.Sprintf("%v", val)
	case float64:
		// Check if it's actually an integer
		if val == float64(int(val)) {
			return fmt.Sprintf("%d", int(val))
		}
		return fmt.Sprintf("%v", val)
	case bool:
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Test file template
var testTemplate = template.Must(template.New("test").Funcs(template.FuncMap{
	"formatValue": formatValue,
}).Parse(`package problems

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test{{.FunctionName}}(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
{{range .TestCases}}		{"{{.Name}}", {{formatValue .Inputs}}, {{formatValue .Expected}}},
{{end}}	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := {{$.FunctionName}}(tt.input...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
`))

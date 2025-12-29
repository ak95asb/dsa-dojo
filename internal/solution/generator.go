package solution

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/problem"
)

// Generator handles solution file generation
type Generator struct {
	solutionsDir string
}

// NewGenerator creates a new solution file generator
func NewGenerator() *Generator {
	return &Generator{
		solutionsDir: "solutions",
	}
}

// GenerateSolution creates a solution file for the problem
func (g *Generator) GenerateSolution(p *database.Problem, force bool) (string, error) {
	// Ensure solutions directory exists
	if err := os.MkdirAll(g.solutionsDir, 0755); err != nil {
		return "", fmt.Errorf("create solutions directory: %w", err)
	}

	// Generate file path
	fileName := problem.SlugToSnakeCase(p.Slug) + ".go"
	filePath := filepath.Join(g.solutionsDir, fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists
		if !force {
			// Prompt for confirmation
			if !promptOverwrite(filePath) {
				return filePath, nil // User chose not to overwrite
			}
		}

		// Create backup
		backupPath := filePath + ".backup"
		if err := copyFile(filePath, backupPath); err != nil {
			return "", fmt.Errorf("create backup: %w", err)
		}
		fmt.Printf("✓ Backup created: %s\n", backupPath)
	}

	// Generate function name
	funcName := slugToFunctionName(p.Slug)

	// Prepare template data
	data := struct {
		FunctionName string
		ProblemTitle string
		Description  string
		Difficulty   string
		Topic        string
		Slug         string
	}{
		FunctionName: funcName,
		ProblemTitle: p.Title,
		Description:  p.Description,
		Difficulty:   p.Difficulty,
		Topic:        p.Topic,
		Slug:         p.Slug,
	}

	// Execute template
	t, err := template.New("solution").Parse(solutionTemplate)
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

// promptOverwrite asks user for confirmation
func promptOverwrite(filePath string) bool {
	fmt.Printf("Solution file '%s' already exists. Overwrite? [y/N]: ", filePath)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))

	return response == "y" || response == "yes"
}

// copyFile creates a backup copy
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// slugToFunctionName converts "two-sum" -> "TwoSum"
func slugToFunctionName(slug string) string {
	parts := strings.Split(slug, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

const solutionTemplate = `package solutions

// {{.ProblemTitle}}
// Difficulty: {{.Difficulty}}
// Topic: {{.Topic}}
//
// Description:
// {{.Description}}
//
// Run tests: dsa test {{.Slug}}

// {{.FunctionName}} solves the {{.ProblemTitle}} problem
func {{.FunctionName}}() {
	// TODO: Implement your solution here
	//
	// Hints:
	// - Read the problem description above carefully
	// - Consider edge cases (empty inputs, single elements, etc.)
	// - Test your solution with: dsa test {{.Slug}}
	// - Run benchmarks with: dsa bench {{.Slug}}
}
`

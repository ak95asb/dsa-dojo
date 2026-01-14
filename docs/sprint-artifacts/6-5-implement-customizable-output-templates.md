# Story 6.5: Implement Customizable Output Templates

Status: ready-for-dev

## Story

As a **user**,
I want **to define custom output templates**,
So that **I can format output exactly how I need it** (FR34, FR35).

## Acceptance Criteria

### AC1: Custom Template for Problem List

**Given** I want a custom problem list format
**When** I create a template file at `~/.dsa/templates/list.tmpl`:
```
{{range .Problems}}
- [{{if .Solved}}x{{else}} {{end}}] {{.Title}} ({{.Difficulty}})
{{end}}
```
**And** I run `dsa list --template list`
**Then** Output uses my custom template:
```
- [x] Two Sum (Easy)
- [ ] Add Two Numbers (Medium)
```
**And** Template uses Go text/template syntax (Architecture pattern)

### AC2: Template Variables and Data Access

**Given** I want to use template variables
**When** My template includes: `{{.Total}} problems ({{.Solved}} solved)`
**Then** Variables are substituted with actual values
**And** Template has access to all data fields

### AC3: Template Management Commands

**Given** I want to share my template
**When** I run `dsa template export list --output mylist.tmpl`
**Then** The template is exported to the specified file
**And** Another user can import it with `dsa template import mylist.tmpl`

### AC4: List Available Templates

**Given** I want to list available templates
**When** I run `dsa template list`
**Then** I see all built-in and custom templates:
  - table (built-in)
  - json (built-in)
  - csv (built-in)
  - markdown (built-in)
  - list (custom)
**And** Description shows what each template does

## Tasks / Subtasks

- [ ] **Task 1: Create Template Management Command** (AC: 3, 4)
  - [ ] Create `dsa template` parent command
  - [ ] Add subcommands: list, export, import
  - [ ] Define template storage location: ~/.dsa/templates/
  - [ ] Create templates directory on first use

- [ ] **Task 2: Implement Template Discovery** (AC: 1, 4)
  - [ ] Scan ~/.dsa/templates/ for .tmpl files
  - [ ] Build list of available custom templates
  - [ ] Define built-in templates (table, json, csv, markdown)
  - [ ] Create template registry/map

- [ ] **Task 3: Implement Template Loading** (AC: 1, 2)
  - [ ] Use text/template package for parsing
  - [ ] Load template from file
  - [ ] Parse template with error handling
  - [ ] Cache parsed templates for performance

- [ ] **Task 4: Define Template Data Structures** (AC: 2)
  - [ ] Create ListTemplateData struct with fields:
    - Problems (array of ProblemData)
    - Total (int)
    - Solved (int)
    - FilteredBy (string, e.g., "difficulty: easy")
  - [ ] Create ProblemData struct for template use
  - [ ] Create StatusTemplateData for status command

- [ ] **Task 5: Integrate Template Rendering with List Command** (AC: 1, 2)
  - [ ] Add --template flag to list command
  - [ ] Load and parse specified template
  - [ ] Execute template with data
  - [ ] Write rendered output to stdout
  - [ ] Handle template execution errors gracefully

- [ ] **Task 6: Implement Template List Command** (AC: 4)
  - [ ] Scan for custom templates
  - [ ] Display built-in templates
  - [ ] Display custom templates
  - [ ] Show template names and descriptions
  - [ ] Format output as table

- [ ] **Task 7: Implement Template Export Command** (AC: 3)
  - [ ] Read template from ~/.dsa/templates/
  - [ ] Write to specified output file
  - [ ] Handle missing template error
  - [ ] Validate output path

- [ ] **Task 8: Implement Template Import Command** (AC: 3)
  - [ ] Read template from specified file
  - [ ] Validate template syntax
  - [ ] Copy to ~/.dsa/templates/
  - [ ] Handle name conflicts (prompt or --force)

- [ ] **Task 9: Add Unit Tests** (AC: All)
  - [ ] Test template loading and parsing
  - [ ] Test template execution with sample data
  - [ ] Test template discovery (list templates)
  - [ ] Test template export
  - [ ] Test template import
  - [ ] Test invalid template syntax handling

- [ ] **Task 10: Add Integration Tests** (AC: All)
  - [ ] Test `dsa list --template custom` with user template
  - [ ] Test template variables substitution
  - [ ] Test `dsa template list` shows templates
  - [ ] Test `dsa template export` creates file
  - [ ] Test `dsa template import` adds template
  - [ ] Test custom template produces expected output

## Dev Notes

### Architecture Patterns and Constraints

**Template System (text/template):**
- Use Go's `text/template` package (standard library)
- Templates stored in `~/.dsa/templates/`
- File naming: `<template-name>.tmpl`
- Template syntax: `{{.FieldName}}`, `{{range .Items}}`, `{{if .Condition}}`

**Template Location:**
```
~/.dsa/
├── config.yaml
├── dsa.db
└── templates/
    ├── list.tmpl       # Custom template
    ├── checklist.tmpl  # Custom template
    └── ...
```

**Built-in vs Custom Templates:**
- **Built-in:** table, json, csv, markdown (hardcoded formats)
- **Custom:** User-created templates in ~/.dsa/templates/
- Built-in templates can't be overridden (reserved names)

### Source Tree Components

**Files to Create:**
- `cmd/template.go` - Template management commands
- `cmd/output_template.go` - Template rendering helpers

**Files to Modify:**
- `cmd/list.go` - Add --template flag support
- `cmd/status.go` - Add --template flag support (optional)
- `cmd/template_test.go` - Unit tests for template commands

**Standard Library Usage:**
- `text/template` - Template parsing and execution
- `os` - File operations for template management
- `path/filepath` - Path handling

### Implementation Guidance

**1. Template Rendering Helper (cmd/output_template.go):**
```go
package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "text/template"
)

// TemplateData represents data available to templates
type TemplateData struct {
    Problems   []ProblemTemplateData
    Total      int
    Solved     int
    FilteredBy string
}

// ProblemTemplateData represents a problem for templates
type ProblemTemplateData struct {
    ID         string
    Title      string
    Difficulty string
    Topic      string
    Solved     bool
}

// getTemplatesDir returns the user's templates directory
func getTemplatesDir() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(home, ".dsa", "templates"), nil
}

// loadTemplate loads and parses a template file
func loadTemplate(name string) (*template.Template, error) {
    templatesDir, err := getTemplatesDir()
    if err != nil {
        return nil, err
    }

    templatePath := filepath.Join(templatesDir, name+".tmpl")

    // Check if file exists
    if _, err := os.Stat(templatePath); os.IsNotExist(err) {
        return nil, fmt.Errorf("template '%s' not found", name)
    }

    // Parse template file
    tmpl, err := template.ParseFiles(templatePath)
    if err != nil {
        return nil, fmt.Errorf("failed to parse template: %w", err)
    }

    return tmpl, nil
}

// executeTemplate renders template with data to stdout
func executeTemplate(tmpl *template.Template, data interface{}) error {
    if err := tmpl.Execute(os.Stdout, data); err != nil {
        return fmt.Errorf("failed to execute template: %w", err)
    }
    return nil
}

// listCustomTemplates returns list of custom templates
func listCustomTemplates() ([]string, error) {
    templatesDir, err := getTemplatesDir()
    if err != nil {
        return nil, err
    }

    // Check if directory exists
    if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
        return []string{}, nil // No templates yet
    }

    // Read directory
    entries, err := os.ReadDir(templatesDir)
    if err != nil {
        return nil, err
    }

    templates := []string{}
    for _, entry := range entries {
        if !entry.IsDir() && filepath.Ext(entry.Name()) == ".tmpl" {
            // Remove .tmpl extension
            name := entry.Name()[:len(entry.Name())-5]
            templates = append(templates, name)
        }
    }

    return templates, nil
}
```

**2. Integration with List Command (cmd/list.go):**
```go
// Add --template flag
var listTemplate string

func init() {
    listCmd.Flags().StringVar(&listFormat, "format", "table", "Output format (table, json, csv, markdown)")
    listCmd.Flags().StringVar(&listTemplate, "template", "", "Use custom template")
    listCmd.Flags().BoolVar(&listCompact, "compact", false, "Compact output")
    rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
    // Check if custom template specified
    if listTemplate != "" {
        // Load template
        tmpl, err := loadTemplate(listTemplate)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }

        // Query problems
        problems, err := queryProblems()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(3)
        }

        // Build template data
        templateData := buildTemplateData(problems)

        // Execute template
        if err := executeTemplate(tmpl, templateData); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
        return
    }

    // ... existing format-based output logic ...
}

func buildTemplateData(problems []Problem) TemplateData {
    templateProblems := make([]ProblemTemplateData, len(problems))
    solvedCount := 0

    for i, p := range problems {
        templateProblems[i] = ProblemTemplateData{
            ID:         p.Slug,
            Title:      p.Title,
            Difficulty: p.Difficulty,
            Topic:      p.Topic,
            Solved:     p.Solved,
        }
        if p.Solved {
            solvedCount++
        }
    }

    return TemplateData{
        Problems:   templateProblems,
        Total:      len(problems),
        Solved:     solvedCount,
        FilteredBy: "", // Set based on filters if any
    }
}
```

**3. Template Management Commands (cmd/template.go):**
```go
package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/olekukonko/tablewriter"
    "github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
    Use:   "template",
    Short: "Manage output templates",
    Long:  `Manage custom output templates for list and status commands.`,
}

var templateListCmd = &cobra.Command{
    Use:   "list",
    Short: "List available templates",
    Args:  cobra.NoArgs,
    Run:   runTemplateList,
}

var templateExportCmd = &cobra.Command{
    Use:   "export <template-name>",
    Short: "Export a template to file",
    Args:  cobra.ExactArgs(1),
    Run:   runTemplateExport,
}

var templateImportCmd = &cobra.Command{
    Use:   "import <file>",
    Short: "Import a template from file",
    Args:  cobra.ExactArgs(1),
    Run:   runTemplateImport,
}

var exportOutput string
var importForce bool

func init() {
    templateCmd.AddCommand(templateListCmd)
    templateCmd.AddCommand(templateExportCmd)
    templateCmd.AddCommand(templateImportCmd)

    templateExportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path")
    templateImportCmd.Flags().BoolVar(&importForce, "force", false, "Overwrite existing template")

    rootCmd.AddCommand(templateCmd)
}

func runTemplateList(cmd *cobra.Command, args []string) {
    // Built-in templates
    builtInTemplates := []struct {
        Name        string
        Description string
    }{
        {"table", "ASCII table format (default)"},
        {"json", "JSON format for scripting"},
        {"csv", "CSV format for spreadsheets"},
        {"markdown", "Markdown table format"},
    }

    // Custom templates
    customTemplates, err := listCustomTemplates()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Display table
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Name", "Type", "Description"})

    for _, tmpl := range builtInTemplates {
        table.Append([]string{tmpl.Name, "Built-in", tmpl.Description})
    }

    for _, name := range customTemplates {
        table.Append([]string{name, "Custom", "User-defined template"})
    }

    table.Render()
}

func runTemplateExport(cmd *cobra.Command, args []string) {
    templateName := args[0]

    // Check output path
    if exportOutput == "" {
        exportOutput = templateName + ".tmpl"
    }

    // Get templates directory
    templatesDir, err := getTemplatesDir()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    sourcePath := filepath.Join(templatesDir, templateName+".tmpl")

    // Read template
    data, err := os.ReadFile(sourcePath)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Fprintf(os.Stderr, "Error: Template '%s' not found\n", templateName)
            os.Exit(2)
        }
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Write to output
    if err := os.WriteFile(exportOutput, data, 0644); err != nil {
        fmt.Fprintf(os.Stderr, "Error: Failed to write file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("✓ Template '%s' exported to %s\n", templateName, exportOutput)
}

func runTemplateImport(cmd *cobra.Command, args []string) {
    sourceFile := args[0]

    // Extract template name from filename
    templateName := filepath.Base(sourceFile)
    if filepath.Ext(templateName) == ".tmpl" {
        templateName = templateName[:len(templateName)-5]
    }

    // Get templates directory
    templatesDir, err := getTemplatesDir()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Create templates directory if it doesn't exist
    if err := os.MkdirAll(templatesDir, 0755); err != nil {
        fmt.Fprintf(os.Stderr, "Error: Failed to create templates directory: %v\n", err)
        os.Exit(1)
    }

    destPath := filepath.Join(templatesDir, templateName+".tmpl")

    // Check if template already exists
    if _, err := os.Stat(destPath); err == nil && !importForce {
        fmt.Fprintf(os.Stderr, "Error: Template '%s' already exists. Use --force to overwrite.\n", templateName)
        os.Exit(2)
    }

    // Read source file
    data, err := os.ReadFile(sourceFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Validate template syntax
    _, err = template.New(templateName).Parse(string(data))
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: Invalid template syntax: %v\n", err)
        os.Exit(2)
    }

    // Write to templates directory
    if err := os.WriteFile(destPath, data, 0644); err != nil {
        fmt.Fprintf(os.Stderr, "Error: Failed to import template: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("✓ Template '%s' imported successfully\n", templateName)
}
```

### Testing Standards

**Unit Test Coverage:**
- Test template loading from file
- Test template parsing with valid syntax
- Test template execution with sample data
- Test template discovery (listCustomTemplates)
- Test invalid template syntax handling
- Test template export function
- Test template import function
- Test template name conflicts

**Integration Test Coverage:**
- Test `dsa list --template custom` with user template
- Test template variables are substituted correctly
- Test `dsa template list` shows all templates
- Test `dsa template export` creates valid file
- Test `dsa template import` adds template
- Test imported template can be used immediately
- Test invalid template file shows error

**Test Pattern:**
```go
func TestTemplateSystem(t *testing.T) {
    t.Run("load and execute template", func(t *testing.T) {
        // Create temp templates directory
        tmpDir := t.TempDir()
        t.Setenv("HOME", tmpDir)

        templatesDir := filepath.Join(tmpDir, ".dsa", "templates")
        os.MkdirAll(templatesDir, 0755)

        // Create test template
        templateContent := `{{range .Problems}}{{.Title}}: {{.Difficulty}}
{{end}}`
        templatePath := filepath.Join(templatesDir, "test.tmpl")
        os.WriteFile(templatePath, []byte(templateContent), 0644)

        // Load template
        tmpl, err := loadTemplate("test")
        assert.NoError(t, err)
        assert.NotNil(t, tmpl)

        // Execute template
        data := TemplateData{
            Problems: []ProblemTemplateData{
                {Title: "Two Sum", Difficulty: "Easy"},
            },
        }

        var buf bytes.Buffer
        err = tmpl.Execute(&buf, data)
        assert.NoError(t, err)
        assert.Contains(t, buf.String(), "Two Sum: Easy")
    })

    t.Run("list custom templates", func(t *testing.T) {
        tmpDir := t.TempDir()
        t.Setenv("HOME", tmpDir)

        templatesDir := filepath.Join(tmpDir, ".dsa", "templates")
        os.MkdirAll(templatesDir, 0755)

        // Create templates
        os.WriteFile(filepath.Join(templatesDir, "list.tmpl"), []byte("test"), 0644)
        os.WriteFile(filepath.Join(templatesDir, "status.tmpl"), []byte("test"), 0644)

        templates, err := listCustomTemplates()
        assert.NoError(t, err)
        assert.Len(t, templates, 2)
        assert.Contains(t, templates, "list")
        assert.Contains(t, templates, "status")
    })
}
```

### Technical Requirements

**Template System Requirements:**

**1. Template Syntax (text/template):**
- Variables: `{{.FieldName}}`
- Range: `{{range .Items}}...{{end}}`
- Conditionals: `{{if .Condition}}...{{end}}`
- Functions: `{{printf "%s" .Value}}`

**2. Template Data Structure:**
```go
type TemplateData struct {
    Problems   []ProblemTemplateData
    Total      int
    Solved     int
    FilteredBy string
}

type ProblemTemplateData struct {
    ID         string
    Title      string
    Difficulty string
    Topic      string
    Solved     bool
}
```

**3. Template Storage:**
- Location: `~/.dsa/templates/`
- Naming: `<template-name>.tmpl`
- Encoding: UTF-8

**4. Built-in Templates:**
- Reserved names: table, json, csv, markdown
- Cannot be overridden by custom templates
- Always available

**Example Custom Templates:**

**Checklist Template (checklist.tmpl):**
```
{{.Total}} Problems ({{.Solved}} solved)

{{range .Problems}}
- [{{if .Solved}}x{{else}} {{end}}] {{.Title}} ({{.Difficulty}})
{{end}}
```

**Simple List Template (simple.tmpl):**
```
{{range .Problems}}{{.Title}}{{if .Solved}} ✓{{end}}
{{end}}
```

**Detailed Template (detailed.tmpl):**
```
Problem Count: {{.Total}}
Solved: {{.Solved}}
Remaining: {{.Total - .Solved}}

{{range .Problems}}
ID: {{.ID}}
Title: {{.Title}}
Difficulty: {{.Difficulty}}
Topic: {{.Topic}}
Status: {{if .Solved}}Solved{{else}}Not Solved{{end}}
---
{{end}}
```

### Previous Story Intelligence (Stories 6.1-6.4)

**Key Learnings:**
- **Format Flag Pattern:** --format for built-in formats
- **Template Flag:** --template for custom templates
- **Output Helpers:** Separate files for each output type
- **Validation:** Validate input early, fail with helpful errors
- **Exit Codes:** 0 (success), 1 (error), 2 (usage), 3 (database)

**Output System Maturity:**
- Four built-in formats: table, json, csv, markdown
- Now adding: Custom template support
- Consistent pattern across all commands
- Helper functions isolated in separate files

**Test Coverage Standards:**
- 10+ unit tests for helper functions
- 8+ integration tests for commands
- testify/assert for assertions
- Table-driven tests for scenarios

### Project Context Reference

**From Architecture Document:**
- **Template System:** Use text/template (standard library)
- **Extensibility:** Architecture must support customization
- **User Control:** Developers can define their own output formats

**From PRD:**
- **Customization:** Users can format output exactly as needed
- **Sharing:** Templates can be exported and shared
- **Flexibility:** Support wide variety of output styles

**Code Patterns to Follow:**
- Go conventions: gofmt, goimports, golangci-lint
- Use standard library (text/template)
- Co-located tests: *_test.go
- Helper functions in separate files

### Definition of Done

- [ ] Template management commands created (template list, export, import)
- [ ] Templates directory created: ~/.dsa/templates/
- [ ] Template loading helper implemented
- [ ] Template execution helper implemented
- [ ] Template data structures defined
- [ ] List command supports --template flag
- [ ] Custom templates can access all data fields
- [ ] Template list command shows built-in and custom templates
- [ ] Template export command saves template to file
- [ ] Template import command adds template
- [ ] Import validates template syntax before adding
- [ ] Invalid templates show helpful error messages
- [ ] Unit tests: 10+ test scenarios
- [ ] Integration tests: 8+ test scenarios
- [ ] All tests pass: `go test ./...`
- [ ] Build succeeds: `go build`
- [ ] Manual test: Create custom template and use it
- [ ] Manual test: Template list shows templates
- [ ] Manual test: Template export/import works
- [ ] All acceptance criteria satisfied

## Dev Agent Record

### Context Reference

<!-- Path(s) to story context XML will be added here by context workflow -->

### Agent Model Used

<!-- Will be filled by dev agent -->

### Debug Log References

### Completion Notes List

### File List

<!-- Will be filled by dev agent with all modified/created files -->

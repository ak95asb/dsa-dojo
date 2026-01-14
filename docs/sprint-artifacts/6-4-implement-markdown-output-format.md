# Story 6.4: Implement Markdown Output Format

Status: ready-for-dev

## Story

As a **user**,
I want **to generate Markdown tables from CLI output**,
So that **I can easily include results in documentation** (FR33).

## Acceptance Criteria

### AC1: Problem List in Markdown Table

**Given** I want problem list in Markdown
**When** I run `dsa list --format markdown`
**Then** Output is a valid Markdown table:
```markdown
| ID | Title | Difficulty | Topic | Status |
|----|-------|------------|-------|--------|
| two-sum | Two Sum | Easy | Arrays | ✓ |
| add-two-numbers | Add Two Numbers | Medium | Linked Lists | ✗ |
```
**And** Table can be pasted directly into README.md
**And** Table renders correctly in GitHub, GitLab, and other Markdown viewers

### AC2: Status Report in Markdown

**Given** I want status report in Markdown
**When** I run `dsa status --format markdown`
**Then** Output includes:
  - Heading: "# DSA Progress Report"
  - Summary statistics
  - Breakdown tables by difficulty and topic
  - Recent activity list
**And** Output is ready to commit to a progress tracking document

### AC3: Generate Progress Report File

**Given** I want to generate a progress report
**When** I run `dsa report --format markdown --output progress.md`
**Then** A complete Markdown report is generated with:
  - Overall stats
  - Progress charts (ASCII or emoji-based)
  - Problem breakdown tables
  - Recent activity
**And** File can be committed to Git for tracking over time

## Tasks / Subtasks

- [ ] **Task 1: Add Markdown Format Support to List Command** (AC: 1)
  - [ ] Add "markdown" as valid format option to --format flag
  - [ ] Update isValidFormat() to include "markdown"
  - [ ] Add Markdown format branch in runList()

- [ ] **Task 2: Implement Markdown Table Helper** (AC: 1, 2, 3)
  - [ ] Create cmd/output_markdown.go with Markdown helpers
  - [ ] Implement writeMarkdownTable() function
  - [ ] Format table with pipe separators (|)
  - [ ] Generate header separator row (|---|---|)
  - [ ] Handle column alignment

- [ ] **Task 3: Implement Markdown Output for List Command** (AC: 1)
  - [ ] Generate Markdown table for problem list
  - [ ] Columns: ID, Title, Difficulty, Topic, Status
  - [ ] Use ✓ and ✗ for status (UTF-8 compatible)
  - [ ] Ensure table renders in GitHub/GitLab
  - [ ] Write output to stdout (pasteable)

- [ ] **Task 4: Add Markdown Format Support to Status Command** (AC: 2)
  - [ ] Add "markdown" format option to status command
  - [ ] Generate complete progress report with sections
  - [ ] Add heading: "# DSA Progress Report"
  - [ ] Add summary statistics paragraph
  - [ ] Add difficulty breakdown table
  - [ ] Add topic breakdown table
  - [ ] Add recent activity list

- [ ] **Task 5: Implement Markdown Report Sections** (AC: 2)
  - [ ] Create formatMarkdownHeading() helper
  - [ ] Create formatMarkdownList() helper
  - [ ] Create formatMarkdownParagraph() helper
  - [ ] Combine sections into complete report
  - [ ] Ensure proper spacing between sections

- [ ] **Task 6: Add Report Command (Future/Optional)** (AC: 3)
  - [ ] Note: Full `dsa report` command may be future enhancement
  - [ ] For MVP: `dsa status --format markdown` provides report
  - [ ] Add --output flag to save to file (optional)
  - [ ] Document that users can redirect: `dsa status --format markdown > progress.md`

- [ ] **Task 7: Add Unit Tests** (AC: All)
  - [ ] Test Markdown table generation with sample data
  - [ ] Test table header formatting
  - [ ] Test table separator row
  - [ ] Test table column alignment
  - [ ] Test Markdown heading formatting
  - [ ] Test Markdown list formatting
  - [ ] Test complete report structure

- [ ] **Task 8: Add Integration Tests** (AC: All)
  - [ ] Test `dsa list --format markdown` produces valid Markdown
  - [ ] Test Markdown table renders correctly
  - [ ] Test `dsa status --format markdown` generates complete report
  - [ ] Test Markdown output can be saved to file
  - [ ] Test saved file is valid Markdown
  - [ ] Test special characters don't break Markdown syntax

## Dev Notes

### Architecture Patterns and Constraints

**Markdown Table Format (GitHub Flavored Markdown):**
```markdown
| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Value 1  | Value 2  | Value 3  |
| Value 4  | Value 5  | Value 6  |
```

**Markdown Alignment:**
- Left-aligned: `|----------|` (default)
- Right-aligned: `|----------:|`
- Center-aligned: `|:----------:|`

**Output Requirements:**
- Pure Markdown text to stdout
- No ANSI colors or terminal formatting
- UTF-8 characters allowed (✓, ✗ symbols)
- Compatible with GitHub, GitLab, CommonMark

### Source Tree Components

**Files to Create:**
- `cmd/output_markdown.go` - Markdown formatting helpers

**Files to Modify:**
- `cmd/list.go` - Add Markdown format support
- `cmd/status.go` - Add Markdown format support
- `cmd/list_test.go` - Add Markdown unit tests
- `cmd/status_test.go` - Add Markdown unit tests

**Standard Library Usage:**
- `strings` - String building and formatting
- `fmt` - String formatting

### Implementation Guidance

**1. Markdown Helper (cmd/output_markdown.go):**
```go
package cmd

import (
    "fmt"
    "os"
    "strings"
)

// MarkdownTable represents a Markdown table
type MarkdownTable struct {
    Headers []string
    Rows    [][]string
}

// writeMarkdownTable outputs a Markdown table to stdout
func writeMarkdownTable(table MarkdownTable) error {
    if len(table.Headers) == 0 {
        return fmt.Errorf("table headers cannot be empty")
    }

    // Calculate column widths
    widths := calculateColumnWidths(table)

    // Write header row
    headerRow := formatMarkdownRow(table.Headers, widths)
    fmt.Fprintln(os.Stdout, headerRow)

    // Write separator row
    separatorRow := formatMarkdownSeparator(len(table.Headers), widths)
    fmt.Fprintln(os.Stdout, separatorRow)

    // Write data rows
    for _, row := range table.Rows {
        dataRow := formatMarkdownRow(row, widths)
        fmt.Fprintln(os.Stdout, dataRow)
    }

    return nil
}

// calculateColumnWidths determines the width of each column
func calculateColumnWidths(table MarkdownTable) []int {
    widths := make([]int, len(table.Headers))

    // Start with header widths
    for i, header := range table.Headers {
        widths[i] = len(header)
    }

    // Check data row widths
    for _, row := range table.Rows {
        for i, cell := range row {
            if i < len(widths) && len(cell) > widths[i] {
                widths[i] = len(cell)
            }
        }
    }

    return widths
}

// formatMarkdownRow formats a table row with padding
func formatMarkdownRow(cells []string, widths []int) string {
    paddedCells := make([]string, len(cells))
    for i, cell := range cells {
        if i < len(widths) {
            paddedCells[i] = padRight(cell, widths[i])
        } else {
            paddedCells[i] = cell
        }
    }
    return "| " + strings.Join(paddedCells, " | ") + " |"
}

// formatMarkdownSeparator creates the header separator row
func formatMarkdownSeparator(numCols int, widths []int) string {
    separators := make([]string, numCols)
    for i := 0; i < numCols; i++ {
        width := 3 // Minimum separator width
        if i < len(widths) {
            width = widths[i]
        }
        separators[i] = strings.Repeat("-", width)
    }
    return "| " + strings.Join(separators, " | ") + " |"
}

// padRight pads string to specified width
func padRight(s string, width int) string {
    if len(s) >= width {
        return s
    }
    return s + strings.Repeat(" ", width-len(s))
}

// writeMarkdownHeading writes a Markdown heading
func writeMarkdownHeading(text string, level int) {
    prefix := strings.Repeat("#", level)
    fmt.Fprintf(os.Stdout, "%s %s\n\n", prefix, text)
}

// writeMarkdownParagraph writes a paragraph with blank line after
func writeMarkdownParagraph(text string) {
    fmt.Fprintf(os.Stdout, "%s\n\n", text)
}

// writeMarkdownList writes a bulleted list
func writeMarkdownList(items []string) {
    for _, item := range items {
        fmt.Fprintf(os.Stdout, "- %s\n", item)
    }
    fmt.Fprintln(os.Stdout) // Blank line after list
}
```

**2. Integration with List Command (cmd/list.go):**
```go
// In runList function, add Markdown format path

func runList(cmd *cobra.Command, args []string) {
    // Validate format
    validFormats := []string{"table", "json", "csv", "markdown"}
    if !isValidFormat(listFormat, validFormats) {
        fmt.Fprintf(os.Stderr, "Error: Invalid format '%s'. Valid formats: %s\n",
            listFormat, strings.Join(validFormats, ", "))
        os.Exit(2)
    }

    // Query problems
    problems, err := queryProblems()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(3)
    }

    if listFormat == "markdown" {
        table := MarkdownTable{
            Headers: []string{"ID", "Title", "Difficulty", "Topic", "Status"},
            Rows:    make([][]string, len(problems)),
        }

        for i, p := range problems {
            status := "✗"
            if p.Solved {
                status = "✓"
            }

            table.Rows[i] = []string{
                p.Slug,
                p.Title,
                p.Difficulty,
                p.Topic,
                status,
            }
        }

        if err := writeMarkdownTable(table); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
        return
    }

    // ... existing table, json, csv format logic ...
}
```

**3. Integration with Status Command (cmd/status.go):**
```go
// In runStatus function, add Markdown format path

func runStatus(cmd *cobra.Command, args []string) {
    // Validate format
    validFormats := []string{"table", "json", "csv", "markdown"}
    if !isValidFormat(statusFormat, validFormats) {
        fmt.Fprintf(os.Stderr, "Error: Invalid format '%s'. Valid formats: %s\n",
            statusFormat, strings.Join(validFormats, ", "))
        os.Exit(2)
    }

    // Query status data
    statusData, err := queryStatus()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(3)
    }

    if statusFormat == "markdown" {
        // Generate Markdown progress report
        writeMarkdownHeading("DSA Progress Report", 1)

        // Summary
        summary := fmt.Sprintf("**Total Problems:** %d | **Solved:** %d | **Remaining:** %d",
            statusData.Total, statusData.Solved, statusData.Total-statusData.Solved)
        writeMarkdownParagraph(summary)

        // Difficulty breakdown
        writeMarkdownHeading("Progress by Difficulty", 2)
        diffTable := MarkdownTable{
            Headers: []string{"Difficulty", "Solved", "Total", "Percentage"},
            Rows: [][]string{
                {
                    "Easy",
                    fmt.Sprintf("%d", statusData.EasySolved),
                    fmt.Sprintf("%d", statusData.EasyTotal),
                    fmt.Sprintf("%.0f%%", calcPercentage(statusData.EasySolved, statusData.EasyTotal)),
                },
                {
                    "Medium",
                    fmt.Sprintf("%d", statusData.MediumSolved),
                    fmt.Sprintf("%d", statusData.MediumTotal),
                    fmt.Sprintf("%.0f%%", calcPercentage(statusData.MediumSolved, statusData.MediumTotal)),
                },
                {
                    "Hard",
                    fmt.Sprintf("%d", statusData.HardSolved),
                    fmt.Sprintf("%d", statusData.HardTotal),
                    fmt.Sprintf("%.0f%%", calcPercentage(statusData.HardSolved, statusData.HardTotal)),
                },
            },
        }
        writeMarkdownTable(diffTable)
        fmt.Fprintln(os.Stdout) // Blank line

        // Topic breakdown
        if len(statusData.TopicCounts) > 0 {
            writeMarkdownHeading("Progress by Topic", 2)
            topicTable := MarkdownTable{
                Headers: []string{"Topic", "Problems Solved"},
                Rows:    make([][]string, 0, len(statusData.TopicCounts)),
            }
            for topic, count := range statusData.TopicCounts {
                topicTable.Rows = append(topicTable.Rows, []string{topic, fmt.Sprintf("%d", count)})
            }
            writeMarkdownTable(topicTable)
        }

        return
    }

    // ... existing table, json, csv format logic ...
}

func calcPercentage(solved, total int) float64 {
    if total == 0 {
        return 0
    }
    return float64(solved) / float64(total) * 100
}
```

**4. Update Format Validation:**
```go
// Update isValidFormat to accept valid formats list
func isValidFormat(format string, validFormats []string) bool {
    for _, valid := range validFormats {
        if format == valid {
            return true
        }
    }
    return false
}
```

### Testing Standards

**Unit Test Coverage:**
- Test writeMarkdownTable() with sample data
- Test calculateColumnWidths() returns correct widths
- Test formatMarkdownRow() pads columns correctly
- Test formatMarkdownSeparator() creates proper separator
- Test writeMarkdownHeading() formats heading correctly
- Test writeMarkdownList() formats list items
- Test special characters don't break Markdown syntax

**Integration Test Coverage:**
- Test `dsa list --format markdown` produces valid Markdown
- Test Markdown table has proper structure (headers, separator, rows)
- Test `dsa status --format markdown` generates complete report
- Test Markdown sections properly separated
- Test Markdown can be saved to file
- Test saved Markdown file is valid
- Test Markdown renders correctly (manual validation in GitHub/GitLab)

**Test Pattern:**
```go
func TestMarkdownOutput(t *testing.T) {
    t.Run("markdown table structure", func(t *testing.T) {
        table := MarkdownTable{
            Headers: []string{"Name", "Age"},
            Rows: [][]string{
                {"Alice", "30"},
                {"Bob", "25"},
            },
        }

        // Capture stdout
        oldStdout := os.Stdout
        r, w, _ := os.Pipe()
        os.Stdout = w

        err := writeMarkdownTable(table)
        assert.NoError(t, err)

        w.Close()
        os.Stdout = oldStdout

        var buf bytes.Buffer
        io.Copy(&buf, r)
        output := buf.String()

        // Check structure
        lines := strings.Split(strings.TrimSpace(output), "\n")
        assert.Len(t, lines, 4) // Header + separator + 2 rows

        assert.Contains(t, lines[0], "| Name")
        assert.Contains(t, lines[1], "|---")
        assert.Contains(t, lines[2], "| Alice")
    })

    t.Run("column width calculation", func(t *testing.T) {
        table := MarkdownTable{
            Headers: []string{"Short", "Very Long Header"},
            Rows: [][]string{
                {"A", "B"},
                {"C", "D"},
            },
        }

        widths := calculateColumnWidths(table)
        assert.Equal(t, 5, widths[0]) // "Short"
        assert.Equal(t, 16, widths[1]) // "Very Long Header"
    })

    t.Run("markdown heading", func(t *testing.T) {
        // Capture stdout
        oldStdout := os.Stdout
        r, w, _ := os.Pipe()
        os.Stdout = w

        writeMarkdownHeading("Test Heading", 2)

        w.Close()
        os.Stdout = oldStdout

        var buf bytes.Buffer
        io.Copy(&buf, r)
        output := buf.String()

        assert.Equal(t, "## Test Heading\n\n", output)
    })
}
```

### Technical Requirements

**Markdown Table Format:**

**1. Table Structure:**
```markdown
| Column 1 | Column 2 |
|----------|----------|
| Value 1  | Value 2  |
```

**2. Column Padding:**
- Calculate max width for each column
- Pad all cells to match column width
- Consistent spacing improves readability

**3. Separator Row:**
- Second row after header
- Dashes matching column widths
- Format: `|----------|----------|`

**4. UTF-8 Characters:**
- ✓ (U+2713) - Check mark for solved
- ✗ (U+2717) - Cross mark for unsolved
- Compatible with GitHub/GitLab Markdown

**5. Heading Levels:**
- H1 (`#`) - Main report title
- H2 (`##`) - Section headings
- H3 (`###`) - Subsection headings (if needed)

**6. Spacing:**
- Blank line after headings
- Blank line after paragraphs
- Blank line after tables
- Blank line after lists

**Example List Command Markdown Output:**
```markdown
| ID              | Title           | Difficulty | Topic         | Status |
|-----------------|-----------------|------------|---------------|--------|
| two-sum         | Two Sum         | Easy       | Arrays        | ✓      |
| add-two-numbers | Add Two Numbers | Medium     | Linked Lists  | ✗      |
```

**Example Status Command Markdown Output:**
```markdown
# DSA Progress Report

**Total Problems:** 20 | **Solved:** 5 | **Remaining:** 15

## Progress by Difficulty

| Difficulty | Solved | Total | Percentage |
|------------|--------|-------|------------|
| Easy       | 2      | 8     | 25%        |
| Medium     | 2      | 8     | 25%        |
| Hard       | 1      | 4     | 25%        |

## Progress by Topic

| Topic         | Problems Solved |
|---------------|-----------------|
| Arrays        | 3               |
| Linked Lists  | 2               |
```

### Previous Story Intelligence (Stories 6.1, 6.2, 6.3)

**Key Learnings:**
- **Format Flag Pattern:** --format flag with values (table, json, csv, markdown)
- **Format Validation:** isValidFormat() validates against list of formats
- **Output Helpers:** Separate files for each format (output_*.go)
- **Stdout for Data:** All output goes to stdout (redirectable)
- **No Extra Formatting:** Pure Markdown (no colors, no terminal codes)

**Output System Maturity:**
- Four output formats: table, json, csv, markdown
- Consistent pattern across all commands
- Helper functions isolated in separate files
- Format validation with helpful error messages
- Exit codes: 0 (success), 1 (error), 2 (usage), 3 (database)

**Test Coverage Standards:**
- 10+ unit tests for helper functions
- 8+ integration tests for commands
- testify/assert for assertions
- Table-driven tests for scenarios

### Project Context Reference

**From Architecture Document:**
- **Standard Library Preferred:** Use strings, fmt packages
- **No External Dependencies:** Pure Markdown generation
- **Scriptability:** Markdown output for documentation workflows
- **UTF-8 Support:** ✓ and ✗ symbols for status

**From PRD:**
- **Documentation Integration:** Markdown for README and docs
- **Progress Tracking:** Markdown reports for Git commits
- **Community Sharing:** Share progress in Markdown format

**Code Patterns to Follow:**
- Go conventions: gofmt, goimports, golangci-lint passing
- Use standard library (strings, fmt)
- Co-located tests: *_test.go
- Helper functions in separate files
- Table-driven tests

### Definition of Done

- [ ] Markdown format added to valid formats list
- [ ] output_markdown.go created with Markdown helpers
- [ ] writeMarkdownTable() generates valid Markdown tables
- [ ] calculateColumnWidths() handles variable-width columns
- [ ] formatMarkdownRow() pads columns correctly
- [ ] formatMarkdownSeparator() creates proper separator
- [ ] writeMarkdownHeading() formats headings
- [ ] List command supports --format markdown
- [ ] Status command supports --format markdown
- [ ] Markdown tables use pipe separators correctly
- [ ] Column widths calculated and applied
- [ ] UTF-8 symbols (✓, ✗) used correctly
- [ ] Status report has sections (heading, summary, tables)
- [ ] Markdown output to stdout (redirectable)
- [ ] Unit tests: 10+ test scenarios
- [ ] Integration tests: 8+ test scenarios
- [ ] All tests pass: `go test ./...`
- [ ] Build succeeds: `go build`
- [ ] Manual test: `dsa list --format markdown` produces valid table
- [ ] Manual test: `dsa status --format markdown` produces report
- [ ] Manual test: Output renders correctly in GitHub
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

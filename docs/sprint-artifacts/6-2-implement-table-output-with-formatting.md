# Story 6.2: Implement Table Output with Formatting

Status: ready-for-dev

## Story

As a **user**,
I want **well-formatted table output for CLI commands**,
So that **data is easy to read in the terminal** (FR30).

## Acceptance Criteria

### AC1: Basic Table Formatting

**Given** I want to list problems in table format
**When** I run `dsa list --format table` (or just `dsa list`)
**Then** Output is a formatted ASCII table:
```
+------------+-------------------+------------+-----------+--------+
| ID         | Title             | Difficulty | Topic     | Status |
+------------+-------------------+------------+-----------+--------+
| two-sum    | Two Sum           | Easy       | Arrays    | ✓      |
| add-two-no | Add Two Numbers   | Medium     | Linked... | ✗      |
+------------+-------------------+------------+-----------+--------+
```
**And** Columns are aligned and padded properly
**And** Long values are truncated with ellipsis (...)
**And** Table fits within terminal width (respects $COLUMNS or detects automatically)

### AC2: Colored Table Output

**Given** I want colored table output
**When** I run `dsa list` with color enabled
**Then** Table uses colors:
  - Headers: Bold
  - Solved problems (✓): Green
  - Unsolved problems (✗): Yellow
  - Difficulty levels: Easy (Green), Medium (Yellow), Hard (Red)
**And** Colors use ANSI escape codes
**And** Colors respect config setting for color_enabled

### AC3: Narrow Terminal Adaptation

**Given** Terminal width is small
**When** I run `dsa list` in a narrow terminal (< 80 columns)
**Then** Table adapts by:
  - Hiding less important columns (e.g., ID)
  - Shortening column widths
  - Truncating long text with ellipsis
**And** Table remains readable and doesn't break layout

## Tasks / Subtasks

- [ ] **Task 1: Choose and Integrate Table Library** (AC: 1, 2, 3)
  - [ ] Research Go table libraries (olekukonko/tablewriter recommended)
  - [ ] Add dependency to go.mod: `go get github.com/olekukonko/tablewriter`
  - [ ] Create table helper in cmd/output_table.go
  - [ ] Test basic table creation with sample data

- [ ] **Task 2: Implement Terminal Width Detection** (AC: 1, 3)
  - [ ] Add dependency for terminal size detection (golang.org/x/term)
  - [ ] Create getTerminalWidth() helper function
  - [ ] Check $COLUMNS environment variable first
  - [ ] Fallback to terminal size detection
  - [ ] Default to 80 columns if detection fails

- [ ] **Task 3: Implement Color Helper Functions** (AC: 2)
  - [ ] Create color helper in cmd/output_color.go
  - [ ] Add colorize() function with ANSI escape codes
  - [ ] Define color constants (Green, Yellow, Red, Bold)
  - [ ] Create shouldUseColors() function checking:
    - NO_COLOR environment variable
    - color_enabled config setting
    - TTY detection (os.Stdout is terminal)
  - [ ] Return raw text when colors disabled

- [ ] **Task 4: Implement Table Formatting for List Command** (AC: 1, 2, 3)
  - [ ] Detect terminal width
  - [ ] Create table with appropriate columns based on width
  - [ ] Apply column widths and truncation
  - [ ] Format headers with bold (if colors enabled)
  - [ ] Format difficulty with colors (Easy=Green, Medium=Yellow, Hard=Red)
  - [ ] Format status with colors (✓=Green, ✗=Yellow)
  - [ ] Align columns (left for text, center for status)

- [ ] **Task 5: Implement Adaptive Column Hiding** (AC: 3)
  - [ ] Define column priority: Status > Title > Difficulty > Topic > ID
  - [ ] Create column selection logic based on terminal width:
    - Width >= 120: Show all columns
    - Width >= 80: Hide ID column
    - Width < 80: Show only Title, Difficulty, Status
  - [ ] Test table output at different widths (60, 80, 120 columns)

- [ ] **Task 6: Implement Text Truncation** (AC: 1, 3)
  - [ ] Create truncateText() helper function
  - [ ] Add ellipsis (...) when text exceeds max width
  - [ ] Calculate max column widths based on terminal width
  - [ ] Apply truncation to Title and Topic columns

- [ ] **Task 7: Integrate with Status Command** (AC: 1, 2)
  - [ ] Apply table formatting to status output
  - [ ] Create tables for:
    - Problems by difficulty (Easy/Medium/Hard counts)
    - Problems by topic (Topic → Count)
    - Recent activity (Problem, Date, Result)
  - [ ] Use colors for visual hierarchy
  - [ ] Ensure tables fit in terminal width

- [ ] **Task 8: Add Unit Tests** (AC: All)
  - [ ] Test terminal width detection
  - [ ] Test color helper returns ANSI codes when enabled
  - [ ] Test color helper returns plain text when disabled
  - [ ] Test shouldUseColors() respects NO_COLOR
  - [ ] Test truncateText() adds ellipsis correctly
  - [ ] Test column selection based on width
  - [ ] Test table generation with sample data

- [ ] **Task 9: Add Integration Tests** (AC: All)
  - [ ] Test `dsa list` produces formatted table
  - [ ] Test table has proper borders and alignment
  - [ ] Test colors applied when enabled
  - [ ] Test NO_COLOR disables colors
  - [ ] Test narrow terminal hides ID column
  - [ ] Test very narrow terminal shows minimal columns
  - [ ] Test long titles truncated with ellipsis
  - [ ] Test status command tables formatted correctly

## Dev Notes

### Architecture Patterns and Constraints

**Terminal Output Libraries (from Architecture):**
- **Colored output:** `fatih/color` or `gookit/color`
- **Tables:** `olekukonko/tablewriter`
- **Terminal size:** `golang.org/x/term`

**Recommended Stack:**
```bash
go get github.com/olekukonko/tablewriter
go get github.com/fatih/color
go get golang.org/x/term
```

**Color Standards (from Architecture):**
- **NO_COLOR Detection:** Respect NO_COLOR environment variable
- **TTY Detection:** Check if stdout is terminal using isatty
- **ANSI Codes:** Use standard ANSI escape sequences
- **Config Integration:** Respect color_scheme and no_color config settings

**Output Requirements (from Architecture):**
- **Stdout for results:** Table output goes to stdout
- **Stderr for errors:** Error messages to stderr
- **Visual Hierarchy:** Clear use of color and spacing
- **Cross-Platform:** ANSI colors work on Windows 10+, macOS, Linux

### Source Tree Components

**New Files to Create:**
- `cmd/output_table.go` - Table formatting helpers
- `cmd/output_color.go` - Color helper functions

**Files to Modify:**
- `cmd/list.go` - Add table output formatting
- `cmd/status.go` - Add table output for status displays
- `cmd/list_test.go` - Add unit tests for table formatting
- `cmd/status_test.go` - Add unit tests for status tables

**Dependencies to Add:**
```bash
go get github.com/olekukonko/tablewriter
go get github.com/fatih/color
go get golang.org/x/term
```

### Implementation Guidance

**1. Color Helper (cmd/output_color.go):**
```go
package cmd

import (
    "os"
    "github.com/fatih/color"
)

// Color definitions
var (
    ColorGreen  = color.New(color.FgGreen)
    ColorYellow = color.New(color.FgYellow)
    ColorRed    = color.New(color.FgRed)
    ColorBold   = color.New(color.Bold)
)

// shouldUseColors checks if colored output should be used
func shouldUseColors() bool {
    // Check NO_COLOR environment variable
    if os.Getenv("NO_COLOR") != "" {
        return false
    }

    // Check config setting (no_color)
    if viper.GetBool("no_color") {
        return false
    }

    // Check if stdout is a terminal
    return isTerminal(os.Stdout)
}

// isTerminal checks if file descriptor is a terminal
func isTerminal(f *os.File) bool {
    return term.IsTerminal(int(f.Fd()))
}

// colorize applies color to text if colors enabled
func colorize(text string, c *color.Color) string {
    if !shouldUseColors() {
        return text
    }
    return c.Sprint(text)
}

// colorDifficulty applies difficulty-specific color
func colorDifficulty(difficulty string) string {
    if !shouldUseColors() {
        return difficulty
    }

    switch strings.ToLower(difficulty) {
    case "easy":
        return ColorGreen.Sprint(difficulty)
    case "medium":
        return ColorYellow.Sprint(difficulty)
    case "hard":
        return ColorRed.Sprint(difficulty)
    default:
        return difficulty
    }
}

// colorStatus applies color to status indicator
func colorStatus(solved bool) string {
    if solved {
        return colorize("✓", ColorGreen)
    }
    return colorize("✗", ColorYellow)
}
```

**2. Terminal Width Detection (cmd/output_table.go):**
```go
package cmd

import (
    "os"
    "strconv"
    "golang.org/x/term"
)

// getTerminalWidth returns the current terminal width
func getTerminalWidth() int {
    // First, check $COLUMNS environment variable
    if cols := os.Getenv("COLUMNS"); cols != "" {
        if width, err := strconv.Atoi(cols); err == nil && width > 0 {
            return width
        }
    }

    // Try to detect terminal size
    if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
        return width
    }

    // Default to 80 columns if detection fails
    return 80
}

// selectColumns returns column names based on terminal width
func selectColumns(width int) []string {
    if width >= 120 {
        // Full display: all columns
        return []string{"ID", "Title", "Difficulty", "Topic", "Status"}
    } else if width >= 80 {
        // Hide ID column
        return []string{"Title", "Difficulty", "Topic", "Status"}
    } else {
        // Narrow terminal: minimal columns
        return []string{"Title", "Difficulty", "Status"}
    }
}

// truncateText truncates text to maxLen with ellipsis
func truncateText(text string, maxLen int) string {
    if len(text) <= maxLen {
        return text
    }
    if maxLen <= 3 {
        return "..."
    }
    return text[:maxLen-3] + "..."
}
```

**3. Table Formatting (cmd/output_table.go):**
```go
package cmd

import (
    "os"
    "github.com/olekukonko/tablewriter"
)

// Problem represents a problem for table display
type ProblemRow struct {
    ID         string
    Title      string
    Difficulty string
    Topic      string
    Solved     bool
}

// printProblemsTable outputs problems as formatted table
func printProblemsTable(problems []ProblemRow) {
    termWidth := getTerminalWidth()
    columns := selectColumns(termWidth)

    // Create table
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader(columns)
    table.SetBorder(true)
    table.SetAlignment(tablewriter.ALIGN_LEFT)

    // Configure table based on terminal width
    if termWidth < 80 {
        table.SetAutoWrapText(false)
    }

    // Calculate max widths for truncation
    titleMaxWidth := calculateTitleWidth(termWidth, len(columns))
    topicMaxWidth := calculateTopicWidth(termWidth, len(columns))

    // Add rows
    for _, p := range problems {
        row := buildTableRow(p, columns, titleMaxWidth, topicMaxWidth)
        table.Append(row)
    }

    table.Render()
}

// buildTableRow creates a table row based on selected columns
func buildTableRow(p ProblemRow, columns []string, titleMaxWidth, topicMaxWidth int) []string {
    row := []string{}
    for _, col := range columns {
        switch col {
        case "ID":
            row = append(row, p.ID)
        case "Title":
            title := truncateText(p.Title, titleMaxWidth)
            row = append(row, title)
        case "Difficulty":
            diff := colorDifficulty(p.Difficulty)
            row = append(row, diff)
        case "Topic":
            topic := truncateText(p.Topic, topicMaxWidth)
            row = append(row, topic)
        case "Status":
            status := colorStatus(p.Solved)
            row = append(row, status)
        }
    }
    return row
}

// calculateTitleWidth determines max width for title column
func calculateTitleWidth(termWidth, numCols int) int {
    // Reserve space for borders, padding, other columns
    reserved := 30 + (numCols * 3) // 3 chars per column for padding/borders
    available := termWidth - reserved

    // Title gets 50% of remaining space
    titleWidth := available / 2
    if titleWidth < 15 {
        titleWidth = 15 // Minimum readable width
    }
    if titleWidth > 40 {
        titleWidth = 40 // Maximum for better layout
    }
    return titleWidth
}

// calculateTopicWidth determines max width for topic column
func calculateTopicWidth(termWidth, numCols int) int {
    if termWidth >= 120 {
        return 20
    } else if termWidth >= 80 {
        return 15
    } else {
        return 0 // Topic column hidden
    }
}
```

**4. Integration with List Command (cmd/list.go):**
```go
// In runList function, add table output path

func runList(cmd *cobra.Command, args []string) {
    // ... existing validation and query logic ...

    // Get format from flag (default: "table")
    format := listFormat

    if format == "table" {
        // Build problem rows
        rows := make([]ProblemRow, len(problems))
        for i, p := range problems {
            rows[i] = ProblemRow{
                ID:         p.Slug,
                Title:      p.Title,
                Difficulty: p.Difficulty,
                Topic:      p.Topic,
                Solved:     p.Solved,
            }
        }

        // Print formatted table
        printProblemsTable(rows)
        return
    }

    if format == "json" {
        // ... existing JSON output logic ...
    }

    // Invalid format
    fmt.Fprintf(os.Stderr, "Error: Invalid format '%s'\n", format)
    os.Exit(2)
}
```

**5. Integration with Status Command (cmd/status.go):**
```go
// Add table output for status displays

func runStatus(cmd *cobra.Command, args []string) {
    // ... existing validation and query logic ...

    format := statusFormat

    if format == "table" {
        printStatusTable(statusData)
        return
    }

    if format == "json" {
        // ... existing JSON output logic ...
    }
}

func printStatusTable(data StatusData) {
    fmt.Println(colorize("Progress Summary", ColorBold))
    fmt.Printf("Total Problems: %d | Solved: %d\n\n", data.Total, data.Solved)

    // Difficulty breakdown table
    diffTable := tablewriter.NewWriter(os.Stdout)
    diffTable.SetHeader([]string{"Difficulty", "Solved", "Total"})
    diffTable.Append([]string{
        colorDifficulty("Easy"),
        fmt.Sprintf("%d", data.EasySolved),
        fmt.Sprintf("%d", data.EasyTotal),
    })
    diffTable.Append([]string{
        colorDifficulty("Medium"),
        fmt.Sprintf("%d", data.MediumSolved),
        fmt.Sprintf("%d", data.MediumTotal),
    })
    diffTable.Append([]string{
        colorDifficulty("Hard"),
        fmt.Sprintf("%d", data.HardSolved),
        fmt.Sprintf("%d", data.HardTotal),
    })
    diffTable.Render()

    // Topic breakdown table
    if len(data.TopicCounts) > 0 {
        fmt.Println("\nProgress by Topic:")
        topicTable := tablewriter.NewWriter(os.Stdout)
        topicTable.SetHeader([]string{"Topic", "Solved"})
        for topic, count := range data.TopicCounts {
            topicTable.Append([]string{topic, fmt.Sprintf("%d", count)})
        }
        topicTable.Render()
    }
}
```

### Testing Standards

**Unit Test Coverage:**
- Test terminal width detection (COLUMNS env, term.GetSize, default)
- Test shouldUseColors() respects NO_COLOR
- Test shouldUseColors() respects no_color config
- Test shouldUseColors() detects TTY
- Test colorDifficulty() applies correct colors
- Test colorStatus() applies correct colors
- Test truncateText() with various lengths
- Test selectColumns() based on width thresholds
- Test calculateTitleWidth() returns reasonable values

**Integration Test Coverage:**
- Test `dsa list` produces ASCII table with borders
- Test table columns aligned properly
- Test colors applied when enabled (check ANSI codes)
- Test NO_COLOR=1 disables colors
- Test narrow terminal (COLUMNS=60) hides columns
- Test wide terminal (COLUMNS=120) shows all columns
- Test long titles truncated with ellipsis
- Test status command produces formatted tables
- Test table output goes to stdout (can redirect)

**Test Pattern:**
```go
func TestTableFormatting(t *testing.T) {
    t.Run("terminal width detection", func(t *testing.T) {
        // Test COLUMNS env var
        t.Setenv("COLUMNS", "100")
        width := getTerminalWidth()
        assert.Equal(t, 100, width)
    })

    t.Run("column selection for wide terminal", func(t *testing.T) {
        columns := selectColumns(120)
        assert.Equal(t, []string{"ID", "Title", "Difficulty", "Topic", "Status"}, columns)
    })

    t.Run("column selection for narrow terminal", func(t *testing.T) {
        columns := selectColumns(60)
        assert.Equal(t, []string{"Title", "Difficulty", "Status"}, columns)
    })

    t.Run("text truncation with ellipsis", func(t *testing.T) {
        text := "This is a very long title that should be truncated"
        truncated := truncateText(text, 20)
        assert.Equal(t, "This is a very lo...", truncated)
        assert.Equal(t, 20, len(truncated))
    })
}

func TestColorHelpers(t *testing.T) {
    t.Run("NO_COLOR disables colors", func(t *testing.T) {
        t.Setenv("NO_COLOR", "1")
        assert.False(t, shouldUseColors())
    })

    t.Run("difficulty colors", func(t *testing.T) {
        // When colors disabled
        t.Setenv("NO_COLOR", "1")
        assert.Equal(t, "Easy", colorDifficulty("Easy"))

        // When colors enabled (would contain ANSI codes)
        // Test by checking output contains color codes
    })

    t.Run("status colors", func(t *testing.T) {
        t.Setenv("NO_COLOR", "1")
        assert.Equal(t, "✓", colorStatus(true))
        assert.Equal(t, "✗", colorStatus(false))
    })
}
```

### Technical Requirements

**Table Formatting Requirements:**

**1. ASCII Table Structure:**
- Use `+`, `-`, `|` for borders
- Proper column alignment (left for text, center for status)
- Consistent padding (1 space on each side of content)
- Example:
  ```
  +----------+-------+
  | Column 1 | Col 2 |
  +----------+-------+
  | Value    | 123   |
  +----------+-------+
  ```

**2. Color Codes (ANSI Escape Sequences):**
- Green: `\033[32m` (Easy difficulty, Solved status)
- Yellow: `\033[33m` (Medium difficulty, Unsolved status)
- Red: `\033[31m` (Hard difficulty)
- Bold: `\033[1m` (Headers)
- Reset: `\033[0m` (After each colored segment)

**3. Terminal Width Detection Priority:**
1. Check `$COLUMNS` environment variable
2. Use `golang.org/x/term.GetSize()` to detect actual width
3. Default to 80 columns if both fail

**4. Column Priority (for adaptive hiding):**
- Priority 1 (Always show): Title, Status
- Priority 2 (Show if width >= 80): Difficulty, Topic
- Priority 3 (Show if width >= 120): ID

**5. Text Truncation Rules:**
- Title column: Max 40 chars (wide), 25 chars (medium), 20 chars (narrow)
- Topic column: Max 20 chars (wide), 15 chars (medium), hide (narrow)
- Truncation: `"Very Long Title..."` (always 3 dots)

**6. Color Configuration Integration:**
- Check NO_COLOR environment variable first
- Check no_color config setting second
- Check TTY detection third
- If any disables colors, output plain text

### Previous Story Intelligence (Story 6.1)

**Key Learnings from Story 6.1:**
- **Output Format Pattern:** --format flag with values (table, json)
- **Format Validation:** isValidFormat() helper to validate format string
- **Output Helpers:** Separate helper files (output_json.go)
- **Stdout/Stderr Separation:** Results to stdout, errors to stderr
- **Error Messages:** Helpful messages with valid options
- **Testing Pattern:** Unit tests for helpers, integration tests for commands

**Output System Established:**
- Dual output modes: human-friendly (table) + machine-parseable (json)
- Format selection via --format flag
- Compact mode for minified output (--compact)
- Helper functions in separate files (output_json.go)

**Test Coverage from Story 6.1:**
- 10+ unit tests covering output formatting
- 8+ integration tests covering command execution
- All tests use testify/assert
- Table-driven tests for multiple scenarios

### Project Context Reference

**From Architecture Document:**
- **Terminal Libraries:** tablewriter for tables, fatih/color for ANSI colors
- **NO_COLOR Support:** Respect NO_COLOR environment variable
- **TTY Detection:** Check if stdout is terminal for color decisions
- **Visual Hierarchy:** Clear use of color and spacing
- **Cross-Platform:** ANSI colors work on Windows 10+, macOS, Linux

**From PRD:**
- **Human-Friendly Output:** Colored terminal output by default
- **Clear Visual Hierarchy:** Appropriate use of color and spacing
- **Progress Visualizations:** ASCII charts for status

**Color Scheme (from Configuration):**
- Configurable via color_scheme setting (default, solarized, monokai, nord)
- Respects no_color config setting
- MVP: Use default scheme, future stories can add themes

**Code Patterns to Follow:**
- Go conventions: gofmt, goimports, golangci-lint passing
- Co-located tests: *_test.go next to implementation
- Helper functions in separate files
- Table-driven tests for multiple scenarios
- Error messages with actionable guidance

### Definition of Done

- [ ] Table library integrated (olekukonko/tablewriter)
- [ ] Color library integrated (fatih/color)
- [ ] Terminal width detection implemented
- [ ] shouldUseColors() respects NO_COLOR and config
- [ ] colorDifficulty() applies Easy/Medium/Hard colors
- [ ] colorStatus() applies Solved/Unsolved colors
- [ ] truncateText() adds ellipsis correctly
- [ ] selectColumns() adapts to terminal width
- [ ] printProblemsTable() outputs formatted table
- [ ] List command uses table output
- [ ] Status command uses table output
- [ ] Headers displayed in bold
- [ ] Columns aligned and padded
- [ ] Long values truncated with ellipsis
- [ ] Narrow terminal hides less important columns
- [ ] Unit tests: 10+ test scenarios covering helpers
- [ ] Integration tests: 8+ test scenarios covering commands
- [ ] All tests pass: `go test ./...`
- [ ] Build succeeds: `go build`
- [ ] Manual test: `dsa list` shows formatted table
- [ ] Manual test: Colors applied when enabled
- [ ] Manual test: NO_COLOR=1 disables colors
- [ ] Manual test: Narrow terminal (COLUMNS=60) adapts layout
- [ ] Manual test: Wide terminal (COLUMNS=120) shows all columns
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

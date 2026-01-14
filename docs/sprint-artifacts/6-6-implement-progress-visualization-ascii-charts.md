# Story 6.6: Implement Progress Visualization (ASCII Charts)

Status: ready-for-dev

## Story

As a **user**,
I want **to see visual progress charts in the terminal**,
So that **I can quickly understand my progress visually** (FR15).

## Acceptance Criteria

### AC1: ASCII Bar Charts for Progress

**Given** I want to visualize progress
**When** I run `dsa status --visual`
**Then** I see ASCII bar charts showing:
```
Progress by Difficulty:
Easy   [████████████████░░░░] 80% (16/20)
Medium [██████░░░░░░░░░░░░░░] 30% (12/40)
Hard   [███░░░░░░░░░░░░░░░░░] 15% (3/20)

Progress by Topic:
Arrays        [████████████░░░░░░░░] 60% (15/25)
Linked Lists  [██████░░░░░░░░░░░░░░] 30% (6/20)
Trees         [████░░░░░░░░░░░░░░░░] 20% (4/20)
```
**And** Progress bars use Unicode block characters (█ ░)
**And** Bars scale to terminal width automatically

### AC2: Sparkline Charts for Activity

**Given** I want a compact visualization
**When** I run `dsa status --sparkline`
**Then** I see sparkline charts using Unicode:
```
Last 30 days: ▁▂▃▅▇█▇▅▃▂▁▂▃▅▇█
```
**And** Sparklines show activity over time
**And** Useful for at-a-glance progress tracking

### AC3: Emoji-Based Visualization

**Given** I want emoji-based visualization
**When** I run `dsa status --emoji`
**Then** Progress is shown with emoji indicators:
```
Progress: 🟩🟩🟩🟩🟩🟩🟩🟩⬜⬜ 80%
Easy:     ⭐⭐⭐⭐⭐ (5/5)
Medium:   ⭐⭐⭐⬜⬜ (3/5)
Hard:     ⭐⬜⬜⬜⬜ (1/5)
```
**And** Emoji makes output friendly and engaging

### AC4: Save Visualization to File

**Given** I want to save visualization to a file
**When** I run `dsa status --visual --output progress.txt`
**Then** The ASCII visualization is saved to the file
**And** File can be shared or committed to Git

## Tasks / Subtasks

- [ ] **Task 1: Implement ASCII Progress Bar Helper** (AC: 1)
  - [ ] Create cmd/output_visual.go with visualization helpers
  - [ ] Implement renderProgressBar() function
  - [ ] Use Unicode block characters: █ (filled), ░ (empty)
  - [ ] Calculate bar width based on terminal width
  - [ ] Calculate filled blocks based on percentage

- [ ] **Task 2: Implement Terminal Width-Aware Bars** (AC: 1)
  - [ ] Detect terminal width (reuse from Story 6.2)
  - [ ] Calculate max bar width (terminal width - label - percentage - padding)
  - [ ] Scale bar to fit available space
  - [ ] Ensure minimum bar width (10 characters)

- [ ] **Task 3: Add Visual Flag to Status Command** (AC: 1)
  - [ ] Add --visual flag to status command
  - [ ] When enabled, display progress bars instead of table
  - [ ] Show difficulty breakdown with bars
  - [ ] Show topic breakdown with bars
  - [ ] Format output with labels and percentages

- [ ] **Task 4: Implement Sparkline Visualization** (AC: 2)
  - [ ] Add --sparkline flag to status command
  - [ ] Use Unicode sparkline characters: ▁▂▃▅▇█
  - [ ] Map activity values to sparkline heights
  - [ ] Show last 30 days of activity
  - [ ] Query database for daily problem counts

- [ ] **Task 5: Implement Emoji Visualization** (AC: 3)
  - [ ] Add --emoji flag to status command
  - [ ] Use emojis: 🟩 (solved), ⬜ (unsolved), ⭐ (star rating)
  - [ ] Display overall progress with emoji blocks
  - [ ] Display difficulty ratings with stars
  - [ ] Ensure UTF-8 emoji support

- [ ] **Task 6: Add Output File Support** (AC: 4)
  - [ ] Add --output flag to status command
  - [ ] Save visualization to specified file
  - [ ] Write plain text (no ANSI colors)
  - [ ] Preserve Unicode characters in file
  - [ ] Show success message with file path

- [ ] **Task 7: Add Unit Tests** (AC: All)
  - [ ] Test renderProgressBar() with various percentages
  - [ ] Test progress bar width calculation
  - [ ] Test sparkline character mapping
  - [ ] Test emoji visualization rendering
  - [ ] Test output file creation
  - [ ] Test Unicode character handling

- [ ] **Task 8: Add Integration Tests** (AC: All)
  - [ ] Test `dsa status --visual` shows progress bars
  - [ ] Test progress bars scale to terminal width
  - [ ] Test `dsa status --sparkline` shows activity
  - [ ] Test `dsa status --emoji` shows emoji indicators
  - [ ] Test `dsa status --visual --output file.txt` creates file
  - [ ] Test file contains valid UTF-8 content

## Dev Notes

### Architecture Patterns and Constraints

**Unicode Characters for Visualization:**
- **Progress Bars:** █ (U+2588 Full Block), ░ (U+2591 Light Shade)
- **Sparklines:** ▁ ▂ ▃ ▄ ▅ ▆ ▇ █ (U+2581 to U+2588)
- **Emojis:** 🟩 (U+1F7E9), ⬜ (U+2B1C), ⭐ (U+2B50)

**UTF-8 Support:**
- Terminal must support UTF-8 encoding
- Modern terminals (macOS Terminal, iTerm, Windows Terminal) support UTF-8
- Graceful degradation if terminal doesn't support Unicode

**Terminal Width Considerations:**
- Reuse terminal width detection from Story 6.2
- Scale visualizations to fit terminal
- Minimum width for readability

### Source Tree Components

**Files to Create:**
- `cmd/output_visual.go` - Visualization helpers

**Files to Modify:**
- `cmd/status.go` - Add --visual, --sparkline, --emoji, --output flags
- `cmd/status_test.go` - Add visualization unit tests

**Standard Library Usage:**
- `unicode/utf8` - UTF-8 string handling
- `golang.org/x/term` - Terminal width detection (already added in Story 6.2)

### Implementation Guidance

**1. Visualization Helpers (cmd/output_visual.go):**
```go
package cmd

import (
    "fmt"
    "strings"
)

const (
    BlockFull  = "█"
    BlockEmpty = "░"
)

// renderProgressBar creates an ASCII progress bar
func renderProgressBar(label string, current, total int, barWidth int) string {
    if total == 0 {
        return fmt.Sprintf("%-15s [%s] 0%% (0/0)", label, strings.Repeat(BlockEmpty, barWidth))
    }

    percentage := float64(current) / float64(total) * 100
    filledWidth := int(float64(barWidth) * float64(current) / float64(total))
    emptyWidth := barWidth - filledWidth

    bar := strings.Repeat(BlockFull, filledWidth) + strings.Repeat(BlockEmpty, emptyWidth)

    return fmt.Sprintf("%-15s [%s] %.0f%% (%d/%d)",
        label, bar, percentage, current, total)
}

// calculateBarWidth determines the width for progress bars
func calculateBarWidth(terminalWidth int) int {
    // Reserve space for label (15), brackets (2), percentage (7), counts (10), padding (6)
    reserved := 40
    barWidth := terminalWidth - reserved

    // Minimum bar width
    if barWidth < 10 {
        barWidth = 10
    }

    // Maximum bar width for readability
    if barWidth > 40 {
        barWidth = 40
    }

    return barWidth
}

// Sparkline characters (8 levels)
var sparklineChars = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

// renderSparkline creates a sparkline from values
func renderSparkline(values []int) string {
    if len(values) == 0 {
        return ""
    }

    // Find max value for scaling
    maxValue := 0
    for _, v := range values {
        if v > maxValue {
            maxValue = v
        }
    }

    if maxValue == 0 {
        return strings.Repeat("▁", len(values))
    }

    // Map each value to sparkline character
    sparkline := ""
    for _, v := range values {
        index := int(float64(v) / float64(maxValue) * 7)
        if index >= len(sparklineChars) {
            index = len(sparklineChars) - 1
        }
        sparkline += sparklineChars[index]
    }

    return sparkline
}

// renderEmojiProgress creates emoji-based progress visualization
func renderEmojiProgress(label string, current, total int) string {
    if total == 0 {
        return fmt.Sprintf("%s: (0/0)", label)
    }

    percentage := float64(current) / float64(total) * 100

    // Use 10 emoji blocks for visualization
    blocks := 10
    filledBlocks := int(float64(blocks) * float64(current) / float64(total))
    emptyBlocks := blocks - filledBlocks

    visual := strings.Repeat("🟩", filledBlocks) + strings.Repeat("⬜", emptyBlocks)

    return fmt.Sprintf("%s: %s %.0f%%", label, visual, percentage)
}

// renderEmojiStars creates star rating visualization
func renderEmojiStars(label string, current, total int) string {
    if total == 0 {
        return fmt.Sprintf("%s: (0/0)", label)
    }

    // Use 5 stars max
    maxStars := 5
    stars := int(float64(maxStars) * float64(current) / float64(total))
    emptyStars := maxStars - stars

    visual := strings.Repeat("⭐", stars) + strings.Repeat("⬜", emptyStars)

    return fmt.Sprintf("%-10s %s (%d/%d)", label, visual, current, total)
}
```

**2. Integration with Status Command (cmd/status.go):**
```go
// Add flags
var statusVisual bool
var statusSparkline bool
var statusEmoji bool
var statusOutput string

func init() {
    statusCmd.Flags().StringVar(&statusFormat, "format", "table", "Output format (table, json, csv, markdown)")
    statusCmd.Flags().BoolVar(&statusVisual, "visual", false, "Display visual progress bars")
    statusCmd.Flags().BoolVar(&statusSparkline, "sparkline", false, "Display sparkline activity chart")
    statusCmd.Flags().BoolVar(&statusEmoji, "emoji", false, "Display emoji-based visualization")
    statusCmd.Flags().StringVarP(&statusOutput, "output", "o", "", "Save output to file")
    rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
    // Query status data
    statusData, err := queryStatus()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(3)
    }

    // Determine output destination
    output := os.Stdout
    if statusOutput != "" {
        file, err := os.Create(statusOutput)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: Failed to create output file: %v\n", err)
            os.Exit(1)
        }
        defer file.Close()
        output = file
    }

    // Visual mode
    if statusVisual {
        termWidth := getTerminalWidth()
        barWidth := calculateBarWidth(termWidth)

        fmt.Fprintln(output, "Progress by Difficulty:")
        fmt.Fprintln(output, renderProgressBar("Easy", statusData.EasySolved, statusData.EasyTotal, barWidth))
        fmt.Fprintln(output, renderProgressBar("Medium", statusData.MediumSolved, statusData.MediumTotal, barWidth))
        fmt.Fprintln(output, renderProgressBar("Hard", statusData.HardSolved, statusData.HardTotal, barWidth))
        fmt.Fprintln(output)

        if len(statusData.TopicCounts) > 0 {
            fmt.Fprintln(output, "Progress by Topic:")
            for topic, count := range statusData.TopicCounts {
                topicTotal := statusData.TopicTotals[topic] // Need to query this
                fmt.Fprintln(output, renderProgressBar(topic, count, topicTotal, barWidth))
            }
        }

        if statusOutput != "" {
            fmt.Printf("✓ Progress visualization saved to %s\n", statusOutput)
        }
        return
    }

    // Sparkline mode
    if statusSparkline {
        // Query last 30 days activity
        dailyActivity, err := queryDailyActivity(30)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(3)
        }

        sparkline := renderSparkline(dailyActivity)
        fmt.Fprintf(output, "Last 30 days: %s\n", sparkline)

        if statusOutput != "" {
            fmt.Printf("✓ Activity sparkline saved to %s\n", statusOutput)
        }
        return
    }

    // Emoji mode
    if statusEmoji {
        total := statusData.Total
        solved := statusData.Solved

        fmt.Fprintln(output, renderEmojiProgress("Progress", solved, total))
        fmt.Fprintln(output)
        fmt.Fprintln(output, renderEmojiStars("Easy", statusData.EasySolved, statusData.EasyTotal))
        fmt.Fprintln(output, renderEmojiStars("Medium", statusData.MediumSolved, statusData.MediumTotal))
        fmt.Fprintln(output, renderEmojiStars("Hard", statusData.HardSolved, statusData.HardTotal))

        if statusOutput != "" {
            fmt.Printf("✓ Emoji visualization saved to %s\n", statusOutput)
        }
        return
    }

    // ... existing format-based output logic (table, json, csv, markdown) ...
}

// queryDailyActivity returns problem counts for last N days
func queryDailyActivity(days int) ([]int, error) {
    // Query database for daily problem counts
    // Return slice of counts (one per day)
    // This will need to be implemented based on database schema
    return []int{}, nil // Placeholder
}
```

### Testing Standards

**Unit Test Coverage:**
- Test renderProgressBar() with 0%, 50%, 100%
- Test renderProgressBar() with edge cases (0/0, 1/1)
- Test calculateBarWidth() with various terminal widths
- Test renderSparkline() with sample activity data
- Test renderSparkline() with all zeros
- Test renderEmojiProgress() with various percentages
- Test renderEmojiStars() with various ratings

**Integration Test Coverage:**
- Test `dsa status --visual` shows progress bars
- Test progress bars have correct percentages
- Test `dsa status --sparkline` shows activity
- Test `dsa status --emoji` shows emoji blocks
- Test `dsa status --visual --output file.txt` creates file
- Test output file contains UTF-8 characters
- Test visualization adapts to terminal width

**Test Pattern:**
```go
func TestVisualization(t *testing.T) {
    t.Run("progress bar 50%", func(t *testing.T) {
        bar := renderProgressBar("Easy", 10, 20, 20)
        assert.Contains(t, bar, "50%")
        assert.Contains(t, bar, "(10/20)")
        // Check bar has roughly equal filled and empty blocks
        assert.Contains(t, bar, BlockFull)
        assert.Contains(t, bar, BlockEmpty)
    })

    t.Run("progress bar 100%", func(t *testing.T) {
        bar := renderProgressBar("Easy", 20, 20, 20)
        assert.Contains(t, bar, "100%")
        assert.Contains(t, bar, strings.Repeat(BlockFull, 20))
        assert.NotContains(t, bar, BlockEmpty)
    })

    t.Run("progress bar 0%", func(t *testing.T) {
        bar := renderProgressBar("Easy", 0, 20, 20)
        assert.Contains(t, bar, "0%")
        assert.Contains(t, bar, strings.Repeat(BlockEmpty, 20))
        assert.NotContains(t, bar, BlockFull)
    })

    t.Run("sparkline with activity", func(t *testing.T) {
        activity := []int{1, 2, 5, 8, 10, 8, 5, 2, 1}
        sparkline := renderSparkline(activity)
        assert.Len(t, sparkline, len(activity))
        assert.Contains(t, sparkline, "█") // Should have max character
        assert.Contains(t, sparkline, "▁") // Should have min character
    })

    t.Run("emoji progress", func(t *testing.T) {
        visual := renderEmojiProgress("Progress", 8, 10)
        assert.Contains(t, visual, "🟩")
        assert.Contains(t, visual, "⬜")
        assert.Contains(t, visual, "80%")
    })
}
```

### Technical Requirements

**Unicode Character Requirements:**

**1. Progress Bars:**
- Filled block: █ (U+2588)
- Empty block: ░ (U+2591)
- Alternative: Use = and - for ASCII-only terminals

**2. Sparklines:**
- 8 levels: ▁ ▂ ▃ ▄ ▅ ▆ ▇ █ (U+2581 to U+2588)
- Map values to appropriate height
- Show last 30 days of activity

**3. Emojis:**
- Green square: 🟩 (U+1F7E9)
- White square: ⬜ (U+2B1C)
- Star: ⭐ (U+2B50)
- Ensure terminal supports color emojis

**4. Bar Width Calculation:**
```
Terminal Width = 80
Reserved Space = 40 (label + brackets + percentage + counts)
Bar Width = 80 - 40 = 40
```

**5. Percentage Calculation:**
```
Percentage = (current / total) * 100
Filled Blocks = (barWidth * current) / total
Empty Blocks = barWidth - filledBlocks
```

**Example Outputs:**

**Visual Mode:**
```
Progress by Difficulty:
Easy   [████████████████░░░░] 80% (16/20)
Medium [██████░░░░░░░░░░░░░░] 30% (12/40)
Hard   [███░░░░░░░░░░░░░░░░░] 15% (3/20)
```

**Sparkline Mode:**
```
Last 30 days: ▁▂▃▅▇█▇▅▃▂▁▂▃▅▇█▇▅▃▂▁▂▃▅▇█▇▅▃▂▁
```

**Emoji Mode:**
```
Progress: 🟩🟩🟩🟩🟩🟩🟩🟩⬜⬜ 80%

Easy:     ⭐⭐⭐⭐⭐ (5/5)
Medium:   ⭐⭐⭐⬜⬜ (3/5)
Hard:     ⭐⬜⬜⬜⬜ (1/5)
```

### Previous Story Intelligence (Stories 6.1-6.5)

**Key Learnings:**
- **Flag Pattern:** Use boolean flags for visualization modes
- **Terminal Width:** Reuse getTerminalWidth() from Story 6.2
- **Output File:** Use --output flag to save visualizations
- **UTF-8 Support:** Modern terminals support Unicode
- **Helper Functions:** Separate visualization logic in output_visual.go

**Output System Maturity:**
- Five built-in formats: table, json, csv, markdown, visual
- Template support for customization
- Now adding: Progress visualization modes
- Consistent flag patterns across all commands

**Test Coverage Standards:**
- 10+ unit tests for helper functions
- 8+ integration tests for commands
- testify/assert for assertions
- Table-driven tests for scenarios

### Project Context Reference

**From Architecture Document:**
- **Progress Visualizations:** ASCII charts for status (NFR36)
- **Terminal Formatting:** Use appropriate Unicode characters
- **Cross-Platform:** UTF-8 support on modern terminals

**From PRD:**
- **Progress Visualization:** Show progress visually (FR32)
- **Motivation:** Visual feedback drives engagement
- **At-a-Glance:** Quick understanding of progress

**Code Patterns to Follow:**
- Go conventions: gofmt, goimports, golangci-lint
- Use standard library (strings, unicode/utf8)
- Co-located tests: *_test.go
- Helper functions in separate files

### Definition of Done

- [ ] Visualization helpers created in output_visual.go
- [ ] renderProgressBar() generates ASCII progress bars
- [ ] calculateBarWidth() handles terminal width
- [ ] renderSparkline() generates sparkline charts
- [ ] renderEmojiProgress() generates emoji visualizations
- [ ] renderEmojiStars() generates star ratings
- [ ] Status command supports --visual flag
- [ ] Status command supports --sparkline flag
- [ ] Status command supports --emoji flag
- [ ] Status command supports --output flag
- [ ] Progress bars use Unicode block characters
- [ ] Sparklines use Unicode sparkline characters
- [ ] Emoji visualizations use UTF-8 emojis
- [ ] Visualizations scale to terminal width
- [ ] Output file support saves visualizations
- [ ] Unit tests: 10+ test scenarios
- [ ] Integration tests: 8+ test scenarios
- [ ] All tests pass: `go test ./...`
- [ ] Build succeeds: `go build`
- [ ] Manual test: `dsa status --visual` shows bars
- [ ] Manual test: `dsa status --sparkline` shows activity
- [ ] Manual test: `dsa status --emoji` shows emojis
- [ ] Manual test: Visualization saved to file works
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

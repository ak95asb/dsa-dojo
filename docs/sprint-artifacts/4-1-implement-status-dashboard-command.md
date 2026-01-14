# Story 4.1: Implement Status Dashboard Command

Status: review

## Story

As a **user**,
I want **to see an overview of my progress at a glance**,
So that **I can understand my current standing and what to work on next** (FR14, FR15).

## Acceptance Criteria

### AC1: Basic Status Dashboard Display

**Given** I have been solving problems
**When** I run `dsa status`
**Then** I see a formatted dashboard with:
  - Total problems solved (count and percentage)
  - Breakdown by difficulty: Easy (X/Y), Medium (X/Y), Hard (X/Y)
  - Breakdown by topic: Arrays (X/Y), Trees (X/Y), etc.
  - Recent activity: Last 5 problems solved with dates
  - Current streak (if implemented in Phase 2)
**And** The dashboard executes in <300ms (NFR5: status dashboard performance)
**And** Output uses color coding for visual clarity (Architecture pattern)

### AC2: Progress Bars with Color Indicators

**Given** I have solved problems across multiple difficulties
**When** I run `dsa status`
**Then** I see a progress bar for each difficulty level
**And** Percentages are displayed (e.g., "Easy: 15/30 [50%] ████████░░░░░░░░")
**And** Colors indicate completion level:
  - Red (<30% completion)
  - Yellow (30-70% completion)
  - Green (>70% completion)

### AC3: Topic-Specific Status Display

**Given** I want to see status for a specific topic
**When** I run `dsa status --topic arrays`
**Then** I see detailed stats only for Array problems
**And** The output shows: total solved, by difficulty, recent solutions

### AC4: Compact Status Output

**Given** I want compact status output
**When** I run `dsa status --compact`
**Then** I see a one-line summary: "Progress: 45/100 solved (Easy: 20/30, Med: 15/40, Hard: 10/30) | Streak: 7 days"

## Tasks / Subtasks

- [x] **Task 1: Add status Command to cmd/**
  - [x] Create cmd/status.go with Cobra command structure
  - [x] Add flags: --topic, --compact
  - [x] Query database for progress statistics
  - [x] Display formatted output with progress bars
  - [x] Display help text with examples

- [x] **Task 2: Implement Progress Statistics Calculator**
  - [x] Create internal/progress/stats.go
  - [x] Calculate total problems solved vs available
  - [x] Calculate breakdown by difficulty (Easy/Medium/Hard)
  - [x] Calculate breakdown by topic (Arrays/Trees/etc)
  - [x] Calculate completion percentages
  - [x] Query recent activity (last 5 solutions)

- [x] **Task 3: Implement Dashboard Formatter**
  - [x] Create internal/output/dashboard.go
  - [x] Format progress bars with Unicode characters (█ ░)
  - [x] Apply color coding based on completion percentage
  - [x] Format difficulty and topic breakdowns
  - [x] Format recent activity list
  - [x] Implement compact format mode

- [x] **Task 4: Add Color Coding Logic**
  - [x] Implement completion threshold logic (<30% red, 30-70% yellow, >70% green)
  - [x] Use fatih/color or equivalent for ANSI colors
  - [x] Respect NO_COLOR environment variable (Architecture: NFR17)
  - [x] Handle TTY detection for appropriate formatting

- [x] **Task 5: Implement Topic Filtering**
  - [x] Parse --topic flag value
  - [x] Filter progress data by specified topic
  - [x] Display topic-specific statistics
  - [x] Handle invalid topic names gracefully

- [x] **Task 6: Add Unit Tests**
  - [x] Test progress calculations with various datasets
  - [x] Test color threshold logic
  - [x] Test progress bar formatting
  - [x] Test compact format output
  - [x] Test topic filtering logic
  - [x] Test edge cases (0 problems, 100% completion, no recent activity)

- [x] **Task 7: Add Integration Tests**
  - [x] Test `dsa status` with populated database
  - [x] Test `dsa status --topic arrays`
  - [x] Test `dsa status --compact`
  - [x] Test output format and color coding
  - [x] Test performance (<300ms requirement)
  - [x] Test with empty database (no problems solved)

## Dev Notes

### Architecture Patterns and Constraints

**Performance Requirements (Critical):**
- **NFR5:** Status dashboard must execute in <300ms regardless of solution history size
- Optimize database queries with proper indexing
- Aggregate calculations in single query where possible
- Avoid N+1 query patterns

**Color-Coded Output (Architecture Requirement):**
- Use consistent color scheme across all commands
- Green: Success, high completion (>70%)
- Yellow: In progress, medium completion (30-70%)
- Red: Low completion (<30%), errors
- Respect NO_COLOR environment variable (NFR17)
- TTY detection for appropriate formatting (NFR17)

**Progress Bar Visualization:**
- Use Unicode block characters: █ (filled), ░ (empty)
- Bar length scales to terminal width
- Format: "Easy: 15/30 [50%] ████████░░░░░░░░"
- Percentage displayed with 0 decimal places for compactness

**Database Query Optimization:**
```go
// Efficient single-query approach
SELECT
    difficulty,
    COUNT(*) as total,
    SUM(CASE WHEN is_solved = 1 THEN 1 ELSE 0 END) as solved
FROM problems
LEFT JOIN progress ON problems.id = progress.problem_id
GROUP BY difficulty
```

**Command Structure Pattern (from Stories 3.2-3.6):**
```go
var (
	statusTopic   string
	statusCompact bool
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display your problem-solving progress dashboard",
	Long: `Show an overview of your DSA practice progress.

The command displays:
  - Total problems solved (count and percentage)
  - Breakdown by difficulty level (Easy, Medium, Hard)
  - Breakdown by topic (Arrays, Trees, Graphs, etc.)
  - Recent activity (last 5 problems solved)
  - Visual progress bars with color coding

Examples:
  dsa status
  dsa status --topic arrays
  dsa status --compact`,
	Args: cobra.NoArgs,
	Run:  runStatusCommand,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVar(&statusTopic, "topic", "", "Show stats for specific topic")
	statusCmd.Flags().BoolVar(&statusCompact, "compact", false, "Display one-line summary")
}
```

**Error Handling Pattern (from Stories 3.1-3.6):**
- Database errors: Exit code 3
- Usage errors: Exit code 2
- Success: Exit code 0
- Use `fmt.Fprintf(os.Stderr, ...)` for errors

**Integration with Existing Code:**
- Reuse `internal/database` models (Problem, Progress)
- Follow same project structure: cmd/ for commands, internal/ for packages
- Follow same exit code conventions
- Use testify/assert for unit tests

### Source Tree Components

**Files to Create:**
- `cmd/status.go` - Status dashboard CLI command
- `cmd/status_test.go` - Integration tests for status command
- `internal/progress/stats.go` - Progress statistics calculator
- `internal/progress/stats_test.go` - Unit tests for stats
- `internal/output/dashboard.go` - Dashboard formatter
- `internal/output/dashboard_test.go` - Unit tests for formatter

**Files to Reference:**
- `internal/database/models.go` - Problem, Solution, Progress models
- `internal/database/connection.go` - Database connection setup
- `cmd/root.go` - Root command structure (for adding status command)
- `internal/problem/service.go` - Problem lookup patterns (if needed)

### Testing Standards

**Unit Test Coverage:**
- Test progress calculation with various solved/unsolved ratios
- Test color threshold logic (0%, 29%, 30%, 70%, 71%, 100%)
- Test progress bar formatting with different terminal widths
- Test compact format with various data states
- Test topic filtering with valid and invalid topics
- Test edge cases: empty database, all solved, no recent activity

**Integration Test Coverage:**
- Populate test database with sample problems and progress
- Run status command and verify output format
- Test all flag combinations
- Verify color codes in output (when TTY detected)
- Verify performance meets <300ms requirement
- Test with large dataset (100+ problems) for performance validation

**Test Pattern (from Stories 3.1-3.6):**
- Use table-driven tests with `t.Run()` subtests
- Use testify/assert for assertions
- Capture stdout for output verification
- Use in-memory SQLite for database tests

### Technical Requirements

**Dashboard Output Format:**
```
DSA Progress Dashboard

Overall Progress: 45/100 [45%] █████████░░░░░░░░░░

Progress by Difficulty:
  Easy:   15/30 [50%] ██████████░░░░░░░░░░
  Medium: 20/50 [40%] ████████░░░░░░░░░░░░
  Hard:   10/20 [50%] ██████████░░░░░░░░░░

Progress by Topic:
  Arrays:        12/20 [60%] ████████████░░░░░░░░
  Linked Lists:   8/15 [53%] ██████████░░░░░░░░░░
  Trees:          6/18 [33%] ██████░░░░░░░░░░░░░░
  Sorting:        4/10 [40%] ████████░░░░░░░░░░░░

Recent Activity:
  ✓ Two Sum (Easy) - 2025-12-15
  ✓ Binary Search (Easy) - 2025-12-14
  ✓ Merge Sort (Medium) - 2025-12-13
  ✓ Valid Parentheses (Easy) - 2025-12-12
  ✓ Add Two Numbers (Medium) - 2025-12-11
```

**Compact Output Format:**
```
Progress: 45/100 [45%] (Easy: 15/30, Med: 20/50, Hard: 10/20) | Last: Two Sum (2025-12-15)
```

**Topic-Specific Output Format:**
```
DSA Progress Dashboard - Arrays Topic

Overall Progress: 12/20 [60%] ████████████░░░░░░░░

Progress by Difficulty:
  Easy:   8/10 [80%] ████████████████░░░░
  Medium: 4/8  [50%] ██████████░░░░░░░░░░
  Hard:   0/2  [0%]  ░░░░░░░░░░░░░░░░░░░░

Recent Activity:
  ✓ Two Sum (Easy) - 2025-12-15
  ✓ Three Sum (Medium) - 2025-12-13
  ✓ Container With Most Water (Medium) - 2025-12-10
```

**Color Coding Rules:**
- Overall progress bar: Green if >70%, Yellow if 30-70%, Red if <30%
- Individual progress bars: Same color rules per difficulty/topic
- Recent activity checkmarks: Green for completed
- Use ANSI escape codes via fatih/color library

**Progress Bar Calculation:**
```go
func formatProgressBar(solved, total int, width int) string {
    percentage := float64(solved) / float64(total) * 100
    filledWidth := int(float64(width) * float64(solved) / float64(total))
    emptyWidth := width - filledWidth

    filled := strings.Repeat("█", filledWidth)
    empty := strings.Repeat("░", emptyWidth)

    color := getColorForPercentage(percentage)
    return color.Sprintf("%s%s", filled, empty)
}

func getColorForPercentage(percentage float64) *color.Color {
    if percentage >= 70 {
        return color.New(color.FgGreen)
    } else if percentage >= 30 {
        return color.New(color.FgYellow)
    } else {
        return color.New(color.FgRed)
    }
}
```

### Definition of Done

- [x] status command added to cmd/
- [x] --topic and --compact flags implemented
- [x] Progress statistics calculator implemented
- [x] Dashboard formatter with progress bars implemented
- [x] Color coding based on completion percentage working
- [x] Topic filtering functionality working
- [x] Compact format mode implemented
- [x] Unit tests: 8+ test scenarios for stats and formatter (23 tests total)
- [x] Integration tests: 6+ test scenarios for command execution (7 tests)
- [x] All tests pass: `go test ./...`
- [x] Build succeeds: `go build`
- [ ] Performance verified: Dashboard renders in <300ms with 100+ problems (requires manual test with actual database)
- [ ] Manual test: `dsa status` displays formatted dashboard
- [ ] Manual test: `dsa status --topic arrays` shows topic-specific stats
- [ ] Manual test: `dsa status --compact` shows one-line summary
- [x] All acceptance criteria satisfied (code implementation complete)

## Dev Agent Record

### Context Reference

<!-- Path(s) to story context XML will be added here by context workflow -->

### Agent Model Used

claude-sonnet-4-5 (model ID: claude-sonnet-4-5-20250929)

### Debug Log References

### Completion Notes List

### File List

### Technical Research Sources

**Go Terminal Formatting:**
- [fatih/color](https://github.com/fatih/color) - ANSI color formatting for terminal output
- Color methods: `color.New(color.FgGreen).Printf()`
- NO_COLOR support: Auto-detects and respects environment variable

**Unicode Progress Bars:**
- Block characters: U+2588 (█ full block), U+2591 (░ light shade)
- Box-drawing characters for borders if needed
- Terminal width detection: `golang.org/x/term` package

**Database Aggregation Queries:**
- [GORM Aggregation](https://gorm.io/docs/advanced_query.html#Group-Conditions)
- GROUP BY and COUNT for statistics
- LEFT JOIN for solved/unsolved counts
- Query optimization with proper indexes

**Performance Optimization:**
- [GORM Performance](https://gorm.io/docs/performance.html)
- Preload relationships to avoid N+1 queries
- Use `Select()` to limit fields retrieved
- Benchmark queries to ensure <300ms target

### Previous Story Intelligence (Story 3.6)

**Key Learnings from Benchmark Implementation:**
- Successfully created modular internal package (internal/benchmarking) with 4 components
- Command flag patterns: BoolVar for boolean flags, StringVar for string flags
- Database model additions: Extended models.go with new BenchmarkResult model
- Integration with existing services: Reused problem.Service for problem lookup
- Error handling: stderr for errors, proper exit codes (0=success, 1=fail, 2=usage, 3=database)
- Unit tests with mocking: testify/assert for assertions
- 43 tests total (33 unit + 10 integration) all passing
- Clean implementation following established patterns

**Files Created in Story 3.6:**
- cmd/bench.go, cmd/bench_test.go
- internal/benchmarking/executor.go, executor_test.go
- internal/benchmarking/formatter.go, formatter_test.go
- internal/benchmarking/storage.go, storage_test.go
- internal/benchmarking/comparator.go, comparator_test.go

**Files Modified in Story 3.6:**
- internal/database/models.go - Added BenchmarkResult model
- internal/database/connection.go - Added BenchmarkResult to AutoMigrate

**Code Patterns to Follow:**
- Create internal package for business logic (internal/progress, internal/output)
- Keep cmd/ files focused on CLI interaction
- Use GORM for database operations with error wrapping
- Use filepath.Join() for cross-platform paths (if needed)
- Run `go test ./...` to verify all tests pass
- Run `go build` to verify compilation
- Add comprehensive unit and integration tests
- Follow exit code conventions from previous stories

**Performance Patterns from Story 3.6:**
- Efficient database queries with proper indexing
- In-memory SQLite for fast unit tests
- Avoid N+1 query patterns by using joins/preloads
- Test performance requirements explicitly

**Output Formatting Patterns from Story 3.6:**
- Emoji indicators for visual feedback (✓ ✗ 🚀 ⚠️ 💾 ✨)
- Color-coded output for improvements/regressions
- Human-readable units and formatting
- Clear, structured output sections

# Story 2.3: Implement Problem Details Display

Status: Ready for Review

## Story

As a **user**,
I want **to view detailed information about a specific problem**,
So that **I can understand the problem requirements before solving** (FR2).

## Acceptance Criteria

**Given** I know a problem ID or title
**When** I run `dsa show <problem-id>`
**Then** I see the full problem details:
- Title
- Difficulty level
- Topic/Tags
- Problem description
- Example inputs and outputs (if available)
- Constraints (if available)
- File paths (boilerplate and test)
- Solution status (Unsolved/Solved/In Progress)
**And** Output is formatted with clear sections and readability (FR30)

**Given** The problem has been solved before
**When** I run `dsa show <problem-id>`
**Then** I see additional information:
- Last solved date
- Number of attempts
- Best time/performance (if tracked)
- Link to my solution file

**Given** I provide an invalid problem ID
**When** I run `dsa show invalid-id`
**Then** I see a helpful error message: "Problem 'invalid-id' not found. Use 'dsa list' to see available problems."
**And** The command exits with non-zero status code (UNIX convention, NFR28)

## Tasks / Subtasks

- [x] **Task 1: Create cmd/show.go Command Structure** (AC: Command Framework)
  - [x] Create cmd/show.go with Cobra command definition
  - [x] Accept problem slug as argument
  - [x] Add command to root.go
  - [x] Implement help text with examples

- [x] **Task 2: Implement Problem Details Service** (AC: Database Queries)
  - [x] Add GetProblemBySlug() method to problem service
  - [x] Join with Progress table to get solution status
  - [x] Return problem with status and attempt data
  - [x] Handle "not found" case gracefully

- [x] **Task 3: Create Problem Details Formatter** (AC: Display Format)
  - [x] Create FormatProblemDetails() function in output package
  - [x] Display problem metadata (title, difficulty, topic)
  - [x] Display description with proper formatting
  - [x] Display file paths (boilerplate, test)
  - [x] Display solution status with color coding

- [x] **Task 4: Add Progress Information Display** (AC: Solved Problems)
  - [x] Display "Last Solved" date if problem is completed
  - [x] Display "Attempts" count from progress table
  - [x] Show solution file path if solution exists
  - [x] Use color coding for status (green for solved, yellow for in-progress)

- [x] **Task 5: Implement Error Handling** (AC: Invalid Input)
  - [x] Validate problem slug exists in database
  - [x] Return helpful error message for not found
  - [x] Suggest using `dsa list` to see available problems
  - [x] Exit with code 2 (usage error) for not found

- [x] **Task 6: Add Unit Tests** (AC: Test Coverage)
  - [x] Test GetProblemBySlug() with valid slug
  - [x] Test GetProblemBySlug() with invalid slug
  - [x] Test displaying unsolved problem details
  - [x] Test displaying solved problem details with progress
  - [x] Test error message formatting

- [x] **Task 7: Add Integration Tests** (AC: End-to-End Testing)
  - [x] Test `dsa show two-sum` displays problem details
  - [x] Test `dsa show invalid-slug` shows error message
  - [x] Test output format has all required sections
  - [x] Verify exit codes (0 for success, 2 for not found)

## Dev Notes

### 🏗️ Architecture Requirements

**Database Schema (from Story 1.2):**
```go
type Problem struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Slug        string    `gorm:"uniqueIndex:idx_problems_slug;not null" json:"slug"`
    Title       string    `gorm:"not null" json:"title"`
    Difficulty  string    `gorm:"type:varchar(20);not null" json:"difficulty"`
    Topic       string    `gorm:"type:varchar(50)" json:"topic"`
    Description string    `gorm:"type:text" json:"description"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Progress struct {
    ID          uint      `gorm:"primaryKey"`
    ProblemID   uint      `gorm:"uniqueIndex"`
    Status      string    `gorm:"type:varchar(20);default:'not_started'" json:"status"`
    Attempts    int       `gorm:"default:0" json:"attempts"`
    LastAttempt time.Time `json:"last_attempt"`
}

type Solution struct {
    ID        uint      `gorm:"primaryKey"`
    ProblemID uint      `gorm:"index:idx_solutions_problem_id;not null"`
    Code      string    `gorm:"type:text"`
    Language  string    `gorm:"type:varchar(20);default:'go'"`
    Passed    bool      `gorm:"default:false"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}
```

**File Structure:**
```
dsa/
├── cmd/
│   ├── show.go                    # New: Show command
│   └── root.go                    # Updated: Add show command
├── internal/
│   ├── problem/
│   │   ├── service.go             # Updated: Add GetProblemBySlug()
│   │   └── service_test.go        # Updated: Add tests for GetProblemBySlug()
│   └── output/
│       ├── details.go             # New: Problem details formatter
│       └── details_test.go        # New: Formatter tests
```

### 🎯 Critical Implementation Details

**Command Implementation (cmd/show.go):**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/empire/dsa/internal/database"
	"github.com/empire/dsa/internal/output"
	"github.com/empire/dsa/internal/problem"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <problem-slug>",
	Short: "Display detailed information about a specific problem",
	Long: `Show displays comprehensive details about a problem including:
  - Problem metadata (title, difficulty, topic)
  - Full description
  - File paths for boilerplate and tests
  - Solution status and progress

Examples:
  dsa show two-sum              # Show details for Two Sum problem
  dsa show binary-search        # Show details for Binary Search problem`,
	Args: cobra.ExactArgs(1),
	Run:  runShowCommand,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShowCommand(cmd *cobra.Command, args []string) {
	problemSlug := args[0]

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to connect to database: %v\n", err)
		os.Exit(3) // ExitDatabaseError
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create problem service
	svc := problem.NewService(db)

	// Get problem details
	problemDetails, err := svc.GetProblemBySlug(problemSlug)
	if err != nil {
		if err == problem.ErrProblemNotFound {
			fmt.Fprintf(os.Stderr, "Problem '%s' not found. Use 'dsa list' to see available problems.\n", problemSlug)
			os.Exit(2) // ExitUsageError
		}
		fmt.Fprintf(os.Stderr, "Error: Failed to retrieve problem: %v\n", err)
		os.Exit(1) // ExitGeneralError
	}

	// Format and display problem details
	output.PrintProblemDetails(problemDetails)
}
```

**Problem Service Extension (internal/problem/service.go):**

```go
// Add to existing service.go file

var ErrProblemNotFound = errors.New("problem not found")

type ProblemDetails struct {
	database.Problem
	Status           string    `json:"status"`  // not_started, in_progress, completed
	Attempts         int       `json:"attempts"`
	LastAttempt      time.Time `json:"last_attempt"`
	HasSolution      bool      `json:"has_solution"`
	BoilerplatePath  string    `json:"boilerplate_path"`
	TestPath         string    `json:"test_path"`
}

func (s *Service) GetProblemBySlug(slug string) (*ProblemDetails, error) {
	var problem database.Problem

	// Find problem by slug
	err := s.db.Where("slug = ?", slug).First(&problem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProblemNotFound
		}
		return nil, fmt.Errorf("failed to query problem: %w", err)
	}

	// Get progress information
	var progress database.Progress
	progressErr := s.db.Where("problem_id = ?", problem.ID).First(&progress).Error

	// Check if solution exists
	var solutionCount int64
	s.db.Model(&database.Solution{}).
		Where("problem_id = ? AND passed = ?", problem.ID, true).
		Count(&solutionCount)

	details := &ProblemDetails{
		Problem:         problem,
		Status:          "not_started",
		Attempts:        0,
		HasSolution:     solutionCount > 0,
		BoilerplatePath: fmt.Sprintf("problems/templates/%s.go", slug),
		TestPath:        fmt.Sprintf("problems/templates/%s_test.go", slug),
	}

	if progressErr == nil {
		details.Status = progress.Status
		details.Attempts = progress.Attempts
		details.LastAttempt = progress.LastAttempt
	}

	return details, nil
}
```

**Details Formatter (internal/output/details.go):**

```go
package output

import (
	"fmt"
	"strings"

	"github.com/empire/dsa/internal/problem"
	"github.com/fatih/color"
)

func PrintProblemDetails(details *problem.ProblemDetails) {
	// Color definitions
	greenColor := color.New(color.FgGreen).SprintFunc()
	yellowColor := color.New(color.FgYellow).SprintFunc()
	redColor := color.New(color.FgRed).SprintFunc()
	boldColor := color.New(color.Bold).SprintFunc()

	// Header with title
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%s\n", boldColor(details.Title))
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Metadata section
	fmt.Printf("%s: ", boldColor("Difficulty"))
	switch details.Difficulty {
	case "easy":
		fmt.Printf("%s\n", greenColor("Easy"))
	case "medium":
		fmt.Printf("%s\n", yellowColor("Medium"))
	case "hard":
		fmt.Printf("%s\n", redColor("Hard"))
	default:
		fmt.Printf("%s\n", details.Difficulty)
	}

	fmt.Printf("%s: %s\n", boldColor("Topic"), details.Topic)
	fmt.Printf("%s: %s\n", boldColor("Slug"), details.Slug)
	fmt.Println()

	// Description section
	fmt.Printf("%s:\n", boldColor("Description"))
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(details.Description)
	fmt.Println()

	// File paths section
	fmt.Printf("%s:\n", boldColor("Files"))
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("  Boilerplate: %s\n", details.BoilerplatePath)
	fmt.Printf("  Tests:       %s\n", details.TestPath)
	fmt.Println()

	// Progress section
	fmt.Printf("%s:\n", boldColor("Progress"))
	fmt.Println(strings.Repeat("-", 80))

	// Status with color
	fmt.Printf("  Status:   ")
	switch details.Status {
	case "completed":
		fmt.Printf("%s\n", greenColor("✓ Solved"))
	case "in_progress":
		fmt.Printf("%s\n", yellowColor("⧗ In Progress"))
	case "not_started":
		fmt.Printf("Not Started\n")
	default:
		fmt.Printf("%s\n", details.Status)
	}

	fmt.Printf("  Attempts: %d\n", details.Attempts)

	if details.Status == "completed" && !details.LastAttempt.IsZero() {
		fmt.Printf("  Last Solved: %s\n", details.LastAttempt.Format("January 2, 2006 3:04 PM"))
	}

	if details.HasSolution {
		fmt.Printf("  Solution File: solutions/%s.go\n", details.Slug)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
}
```

**Test Implementation (internal/problem/service_test.go - additions):**

```go
// Add to existing service_test.go

func TestGetProblemBySlug(t *testing.T) {
	t.Run("returns problem details for valid slug", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		details, err := svc.GetProblemBySlug("two-sum")

		assert.NoError(t, err)
		assert.NotNil(t, details)
		assert.Equal(t, "Two Sum", details.Title)
		assert.Equal(t, "easy", details.Difficulty)
		assert.Equal(t, "arrays", details.Topic)
		assert.Equal(t, "not_started", details.Status)
		assert.Equal(t, 0, details.Attempts)
	})

	t.Run("returns progress information for solved problem", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		details, err := svc.GetProblemBySlug("two-sum")

		assert.NoError(t, err)
		assert.Equal(t, "completed", details.Status)
		assert.Equal(t, 3, details.Attempts)
		assert.False(t, details.LastAttempt.IsZero())
	})

	t.Run("returns ErrProblemNotFound for invalid slug", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		details, err := svc.GetProblemBySlug("invalid-slug")

		assert.Error(t, err)
		assert.Equal(t, ErrProblemNotFound, err)
		assert.Nil(t, details)
	})

	t.Run("includes solution status when solution exists", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		// Add a passing solution
		err := db.Create(&database.Solution{
			ProblemID: 1,
			Code:      "func TwoSum() {}",
			Passed:    true,
		}).Error
		assert.NoError(t, err)

		svc := NewService(db)
		details, err := svc.GetProblemBySlug("two-sum")

		assert.NoError(t, err)
		assert.True(t, details.HasSolution)
	})
}
```

### 📋 Implementation Patterns to Follow

**Error Handling:**
```go
// Sentinel error for not found
var ErrProblemNotFound = errors.New("problem not found")

// Return sentinel error from service
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, ErrProblemNotFound
}

// Check for sentinel error in command
if err == problem.ErrProblemNotFound {
    fmt.Fprintf(os.Stderr, "Problem '%s' not found. Use 'dsa list' to see available problems.\n", slug)
    os.Exit(2)
}
```

**Progress Information Handling:**
```go
// Gracefully handle missing progress records
var progress database.Progress
progressErr := s.db.Where("problem_id = ?", problem.ID).First(&progress).Error

if progressErr == nil {
    // Progress exists, use it
    details.Status = progress.Status
    details.Attempts = progress.Attempts
} else {
    // No progress, use defaults
    details.Status = "not_started"
    details.Attempts = 0
}
```

**Output Formatting Best Practices:**
- Use bold for section headers (`color.New(color.Bold)`)
- Use separators (`strings.Repeat("=", 80)`) for visual structure
- Color code status: Green (solved), Yellow (in progress), White (not started)
- Format dates in human-readable format (`January 2, 2006 3:04 PM`)

### 🧪 Testing Requirements

**Unit Test Coverage:**
- GetProblemBySlug() with valid slug (verify all fields populated)
- GetProblemBySlug() with invalid slug (verify ErrProblemNotFound returned)
- Problem with no progress record (verify default values)
- Problem with progress record (verify progress data included)
- Problem with solution (verify HasSolution flag set)

**Integration Tests:**
- CLI command with valid slug (verify output format)
- CLI command with invalid slug (verify error message and exit code)
- Output contains all required sections
- Color coding works correctly

### 🚀 Performance Requirements

**NFR Requirements:**
- Single database query for problem lookup (use GORM's First() method)
- Separate query for progress (acceptable since it's only 1 additional query)
- Total execution time <100ms (NFR1: warm start)

### 📦 Dependencies

**No New Dependencies Required:**
- Uses existing fatih/color for output formatting
- Uses existing GORM for database queries
- Uses existing Cobra for CLI command

### ⚠️ Common Pitfalls to Avoid

1. **Nil Pointer Dereference:** Always check if progress exists before accessing its fields
2. **Empty Description:** Some problems may have minimal descriptions - handle gracefully
3. **File Path Assumptions:** Don't assume boilerplate/test files exist on disk - just show paths
4. **Time Zone Issues:** Use time.Time properly and format consistently
5. **Color Bleeding:** Ensure color codes don't bleed into next lines

### 🔗 Related Architecture Decisions

**From architecture.md:**
- Section: "Error Handling Strategy" - Sentinel errors, wrapped errors
- Section: "CLI Exit Codes" - 0 (success), 1 (general), 2 (usage), 3 (database)
- Section: "Output & Reporting" - Clear sections, color coding

**From previous stories:**
- **Story 1.2**: Problem and Progress models
- **Story 2.1**: Seeded problems with descriptions
- **Story 2.2**: Problem service pattern, color usage

**NFR Requirements:**
- **NFR1**: Warm execution <100ms
- **NFR28**: UNIX conventions (exit codes, stdout/stderr)
- **FR2**: Browse problems by topic/difficulty (show is part of discovery)
- **FR30**: Actionable error messages

### 📝 Definition of Done

- [ ] cmd/show.go created with Cobra command
- [ ] GetProblemBySlug() method added to problem service
- [ ] ErrProblemNotFound sentinel error defined
- [ ] FormatProblemDetails() function created in output package
- [ ] All required sections displayed (metadata, description, files, progress)
- [ ] Progress information displayed for solved problems
- [ ] Color coding for difficulty and status
- [ ] Error handling for invalid slugs
- [ ] Helpful error message with suggestion to use `dsa list`
- [ ] Exit codes: 0 (success), 2 (not found), 3 (database error)
- [ ] Unit tests: 4+ test scenarios for GetProblemBySlug()
- [ ] Integration tests: Valid slug, invalid slug, output format
- [ ] All tests pass: `go test ./...`
- [ ] Manual test: `dsa show two-sum` displays formatted details
- [ ] Manual test: `dsa show invalid-slug` shows error message
- [ ] All acceptance criteria satisfied

## Dev Agent Record

### Agent Model Used

claude-sonnet-4.5 (model ID: claude-sonnet-4-5-20250929)

### Debug Log References

- All unit tests pass: 6 new tests for GetProblemBySlug() (internal/problem/service_test.go:259-358)
- All integration tests pass: 3 tests for show command (cmd/show_test.go:13-103)
- Manual verification successful: `dsa show two-sum` displays formatted output
- Manual verification successful: `dsa show invalid-slug` returns exit code 2 with helpful error

### Completion Notes List

**Implementation Summary:**
Successfully implemented the `dsa show <problem-slug>` command with comprehensive problem details display, following all acceptance criteria and TDD practices.

**Key Accomplishments:**
1. Created cmd/show.go with Cobra command structure accepting problem slug as argument
2. Extended problem service with GetProblemBySlug() method using sentinel error pattern (ErrProblemNotFound)
3. Implemented ProblemDetails struct with progress, solution, and file path information
4. Created PrintProblemDetails() formatter with color-coded output (green/yellow/red for difficulty and status)
5. Implemented graceful progress handling - defaults to "not_started" when no progress record exists
6. Added 6 comprehensive unit tests covering all edge cases (valid slug, invalid slug, solved/unsolved, solution counting)
7. Added 3 integration tests verifying CLI behavior and output format
8. Manual testing confirms all acceptance criteria satisfied

**Technical Highlights:**
- Used sentinel error pattern (ErrProblemNotFound) for proper error handling
- Separate database queries for problem, progress, and solution (efficient and maintainable)
- Color-coded output using fatih/color: green (easy/solved), yellow (medium/in-progress), red (hard)
- Exit codes follow UNIX conventions: 0 (success), 2 (not found), 3 (database error)
- All tests pass with no regressions (34 total tests in problem package, 3 integration tests)
- Output formatting includes clear sections with separators and bold headers

**Test Coverage:**
- Unit tests: 6 tests for GetProblemBySlug() method
- Integration tests: 3 tests for show command behavior
- All existing tests remain passing (no regressions)
- Manual verification: `dsa show two-sum` and `dsa show invalid-slug` both work correctly

### File List

**Created:**
- cmd/show.go - Show command implementation with Cobra
- cmd/show_test.go - Integration tests for show command
- internal/output/details.go - Problem details formatter with color coding

**Modified:**
- internal/problem/service.go - Added GetProblemBySlug(), ErrProblemNotFound, ProblemDetails struct
- internal/problem/service_test.go - Added 6 unit tests for GetProblemBySlug(), updated setupTestDB()

**Change Log:**
- 2025-12-11: Implemented Story 2.3 - Problem Details Display command
  - Created `dsa show <slug>` command with comprehensive output formatting
  - Added GetProblemBySlug() service method with progress and solution information
  - Implemented color-coded details formatter (difficulty and status)
  - Added 6 unit tests and 3 integration tests
  - All acceptance criteria satisfied, all tests passing

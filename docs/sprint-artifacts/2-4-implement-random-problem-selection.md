# Story 2.4: Implement Random Problem Selection

Status: Ready for Review

## Story

As a **user**,
I want **to get a random problem suggestion**,
So that **I can practice without decision fatigue** (FR4).

## Acceptance Criteria

**Given** I have unsolved problems in my library
**When** I run `dsa random`
**Then** I see a randomly selected unsolved problem displayed
**And** The display includes: Title, Difficulty, Topic, Description
**And** I see a suggested command: "Run 'dsa solve <problem-id>' to start"

**Given** I want a random problem of specific difficulty
**When** I run `dsa random --difficulty easy`
**Then** I get a random Easy problem
**When** I run `dsa random --difficulty medium`
**Then** I get a random Medium problem
**When** I run `dsa random --difficulty hard`
**Then** I get a random Hard problem

**Given** I want a random problem from a specific topic
**When** I run `dsa random --topic arrays`
**Then** I get a random Array problem
**When** I run `dsa random --topic graphs`
**Then** I get a random Graph problem

**Given** I want to combine filters
**When** I run `dsa random --difficulty hard --topic trees`
**Then** I get a random Hard-level Tree problem that is unsolved

**Given** All problems matching the filter are solved
**When** I run `dsa random --difficulty easy` and all Easy problems are solved
**Then** I see a message: "All Easy problems are solved! Try --difficulty medium or --difficulty hard"
**And** The command suggests next steps

## Tasks / Subtasks

- [x] **Task 1: Create cmd/random.go Command Structure** (AC: Command Framework)
  - [x] Create cmd/random.go with Cobra command definition
  - [x] Add --difficulty and --topic flags
  - [x] Add command to root.go
  - [x] Implement help text with examples

- [x] **Task 2: Implement Random Selection Service** (AC: Database Queries)
  - [x] Add GetRandomProblem() method to problem service
  - [x] Accept filters (difficulty, topic, solved status)
  - [x] Use database query to get unsolved problems matching filters
  - [x] Implement random selection logic
  - [x] Return problem details using existing ProblemDetails struct

- [x] **Task 3: Create Random Problem Formatter** (AC: Display Format)
  - [x] Create PrintRandomProblem() function in output package
  - [x] Display problem metadata (title, difficulty, topic)
  - [x] Display description
  - [x] Show suggested command: "Run 'dsa solve <slug>' to start"
  - [x] Use color coding consistent with show command

- [x] **Task 4: Implement Edge Case Handling** (AC: All Problems Solved)
  - [x] Detect when no unsolved problems match filter
  - [x] Return helpful error message with suggestions
  - [x] Suggest alternative difficulty levels or topics
  - [x] Exit with appropriate code

- [x] **Task 5: Implement Flag Validation** (AC: Input Validation)
  - [x] Validate --difficulty flag (easy, medium, hard)
  - [x] Validate --topic flag against known topics
  - [x] Return helpful error for invalid values
  - [x] Exit with code 2 (usage error) for invalid input

- [x] **Task 6: Add Unit Tests** (AC: Test Coverage)
  - [x] Test GetRandomProblem() with no filters
  - [x] Test GetRandomProblem() with difficulty filter
  - [x] Test GetRandomProblem() with topic filter
  - [x] Test GetRandomProblem() with combined filters
  - [x] Test GetRandomProblem() when all problems solved
  - [x] Test randomness (verify different results on multiple calls)

- [x] **Task 7: Add Integration Tests** (AC: End-to-End Testing)
  - [x] Test `dsa random` displays random problem
  - [x] Test `dsa random --difficulty easy` filters correctly
  - [x] Test `dsa random --topic arrays` filters correctly
  - [x] Test `dsa random --difficulty hard --topic trees` combines filters
  - [x] Test error message when all problems solved
  - [x] Verify exit codes (0 for success, 2 for no problems)

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
```

**File Structure:**
```
dsa/
├── cmd/
│   ├── random.go                  # New: Random command
│   └── root.go                    # Updated: Add random command
├── internal/
│   ├── problem/
│   │   ├── service.go             # Updated: Add GetRandomProblem()
│   │   └── service_test.go        # Updated: Add tests for GetRandomProblem()
│   └── output/
│       ├── random.go              # New: Random problem formatter
│       └── random_test.go         # New: Formatter tests
```

### 🎯 Critical Implementation Details

**Command Implementation (cmd/random.go):**

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

var (
	randomDifficulty string
	randomTopic      string
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get a random unsolved problem suggestion",
	Long: `Random selects a random unsolved problem from your library.

Use filters to narrow down the selection:
  --difficulty: Filter by difficulty (easy, medium, hard)
  --topic: Filter by topic (arrays, linked-lists, trees, graphs, sorting, searching)

Examples:
  dsa random                              # Any random unsolved problem
  dsa random --difficulty easy            # Random easy problem
  dsa random --topic arrays               # Random array problem
  dsa random --difficulty hard --topic trees  # Random hard tree problem`,
	Run: runRandomCommand,
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.Flags().StringVar(&randomDifficulty, "difficulty", "", "Filter by difficulty (easy, medium, hard)")
	randomCmd.Flags().StringVar(&randomTopic, "topic", "", "Filter by topic")
}

func runRandomCommand(cmd *cobra.Command, args []string) {
	// Validate flags
	if randomDifficulty != "" && !problem.IsValidDifficulty(randomDifficulty) {
		fmt.Fprintf(os.Stderr, "Invalid difficulty '%s'. Valid options: easy, medium, hard\n", randomDifficulty)
		os.Exit(2) // ExitUsageError
	}

	if randomTopic != "" && !problem.IsValidTopic(randomTopic) {
		fmt.Fprintf(os.Stderr, "Invalid topic '%s'. Use 'dsa list' to see available topics\n", randomTopic)
		os.Exit(2) // ExitUsageError
	}

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

	// Build filters for unsolved problems
	filters := problem.ListFilters{
		Difficulty: randomDifficulty,
		Topic:      randomTopic,
	}
	solved := false
	filters.Solved = &solved // Only unsolved problems

	// Get random problem
	randomProblem, err := svc.GetRandomProblem(filters)
	if err != nil {
		if err == problem.ErrNoProblemsFound {
			// Generate helpful error message
			output.PrintNoProblemsMessage(filters)
			os.Exit(2) // ExitUsageError (no problems is a usage issue, not an error)
		}
		fmt.Fprintf(os.Stderr, "Error: Failed to retrieve random problem: %v\n", err)
		os.Exit(1) // ExitGeneralError
	}

	// Format and display random problem
	output.PrintRandomProblem(randomProblem)
}
```

**Problem Service Extension (internal/problem/service.go):**

```go
// Add to existing service.go file

var ErrNoProblemsFound = errors.New("no problems found matching criteria")

// GetRandomProblem retrieves a random problem matching the given filters
// Returns ErrNoProblemsFound if no problems match the criteria
func (s *Service) GetRandomProblem(filters ListFilters) (*ProblemDetails, error) {
	// Get all problems matching filters
	problems, err := s.ListProblems(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to query problems: %w", err)
	}

	if len(problems) == 0 {
		return nil, ErrNoProblemsFound
	}

	// Randomly select one problem
	// Note: For cryptographic randomness, use crypto/rand
	// For CLI tool, math/rand is sufficient
	randomIndex := rand.Intn(len(problems))
	selectedProblem := problems[randomIndex]

	// Get full details for the selected problem
	details, err := s.GetProblemBySlug(selectedProblem.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem details: %w", err)
	}

	return details, nil
}
```

**Random Problem Formatter (internal/output/random.go):**

```go
package output

import (
	"fmt"
	"strings"

	"github.com/empire/dsa/internal/problem"
	"github.com/fatih/color"
)

// PrintRandomProblem formats and displays a random problem suggestion
func PrintRandomProblem(details *problem.ProblemDetails) {
	// Color definitions
	greenColor := color.New(color.FgGreen).SprintFunc()
	yellowColor := color.New(color.FgYellow).SprintFunc()
	redColor := color.New(color.FgRed).SprintFunc()
	boldColor := color.New(color.Bold).SprintFunc()
	cyanColor := color.New(color.FgCyan).SprintFunc()

	fmt.Println()
	fmt.Println(cyanColor("🎲 Random Problem Selected!"))
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

	// Suggested next step
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%s\n", cyanColor(fmt.Sprintf("▶ Run 'dsa solve %s' to start solving!", details.Slug)))
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
}

// PrintNoProblemsMessage displays a helpful message when no problems match filters
func PrintNoProblemsMessage(filters problem.ListFilters) {
	redColor := color.New(color.FgRed).SprintFunc()
	yellowColor := color.New(color.FgYellow).SprintFunc()

	fmt.Println()
	fmt.Println(redColor("✗ No unsolved problems found!"))
	fmt.Println()

	// Build descriptive message based on filters
	if filters.Difficulty != "" && filters.Topic != "" {
		fmt.Printf("All %s %s problems are solved!\n", filters.Difficulty, filters.Topic)
	} else if filters.Difficulty != "" {
		fmt.Printf("All %s problems are solved!\n", filters.Difficulty)
	} else if filters.Topic != "" {
		fmt.Printf("All %s problems are solved!\n", filters.Topic)
	} else {
		fmt.Println("All problems in your library are solved!")
	}

	fmt.Println()
	fmt.Println(yellowColor("Suggestions:"))

	// Provide context-aware suggestions
	if filters.Difficulty == "easy" {
		fmt.Println("  • Try --difficulty medium for a bigger challenge")
		fmt.Println("  • Try --difficulty hard for expert-level problems")
	} else if filters.Difficulty == "medium" {
		fmt.Println("  • Try --difficulty hard for expert-level problems")
		fmt.Println("  • Try --difficulty easy to review fundamentals")
	} else if filters.Difficulty == "hard" {
		fmt.Println("  • Congratulations on solving all hard problems!")
		fmt.Println("  • Try --difficulty medium or --difficulty easy to practice speed")
	}

	if filters.Topic != "" {
		fmt.Println("  • Try a different topic to broaden your skills")
		fmt.Println("  • Run 'dsa list' to see all available topics")
	} else {
		fmt.Println("  • Run 'dsa list' to see all problems")
		fmt.Println("  • Consider adding custom problems with 'dsa add'")
	}

	fmt.Println()
}
```

**Test Implementation (internal/problem/service_test.go - additions):**

```go
// Add to existing service_test.go

func TestGetRandomProblem(t *testing.T) {
	t.Run("returns random problem with no filters", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		details, err := svc.GetRandomProblem(ListFilters{Solved: &solved})

		assert.NoError(t, err)
		assert.NotNil(t, details)
		// Should be one of the unsolved problems
		assert.Contains(t, []string{"add-two-numbers", "validate-bst", "binary-search", "merge-k-lists"}, details.Slug)
	})

	t.Run("returns random problem with difficulty filter", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		details, err := svc.GetRandomProblem(ListFilters{
			Difficulty: "easy",
			Solved:     &solved,
		})

		assert.NoError(t, err)
		assert.NotNil(t, details)
		assert.Equal(t, "easy", details.Difficulty)
		// Should be either two-sum or reverse-linked-list (but two-sum and reverse are solved)
		// So should be binary-search
		assert.Equal(t, "binary-search", details.Slug)
	})

	t.Run("returns random problem with topic filter", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		details, err := svc.GetRandomProblem(ListFilters{
			Topic:  "linked-lists",
			Solved: &solved,
		})

		assert.NoError(t, err)
		assert.NotNil(t, details)
		assert.Equal(t, "linked-lists", details.Topic)
		// Should be add-two-numbers or merge-k-lists (unsolved linked list problems)
		assert.Contains(t, []string{"add-two-numbers", "merge-k-lists"}, details.Slug)
	})

	t.Run("returns random problem with combined filters", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		details, err := svc.GetRandomProblem(ListFilters{
			Difficulty: "medium",
			Topic:      "linked-lists",
			Solved:     &solved,
		})

		assert.NoError(t, err)
		assert.NotNil(t, details)
		assert.Equal(t, "medium", details.Difficulty)
		assert.Equal(t, "linked-lists", details.Topic)
		assert.Equal(t, "add-two-numbers", details.Slug) // Only one matching problem
	})

	t.Run("returns ErrNoProblemsFound when all problems solved", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		// Mark all problems as solved
		var problems []database.Problem
		db.Find(&problems)
		for _, p := range problems {
			db.Create(&database.Progress{
				ProblemID: p.ID,
				Status:    "completed",
				Attempts:  1,
			})
		}

		svc := NewService(db)
		solved := false
		details, err := svc.GetRandomProblem(ListFilters{Solved: &solved})

		assert.Error(t, err)
		assert.Equal(t, ErrNoProblemsFound, err)
		assert.Nil(t, details)
	})

	t.Run("returns different problems on multiple calls (randomness check)", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		filters := ListFilters{Solved: &solved}

		// Get 10 random problems and check we get some variety
		slugs := make(map[string]bool)
		for i := 0; i < 10; i++ {
			details, err := svc.GetRandomProblem(filters)
			assert.NoError(t, err)
			slugs[details.Slug] = true
		}

		// With 4 unsolved problems and 10 selections, we should see at least 2 different problems
		// (statistically very likely, though not guaranteed)
		assert.GreaterOrEqual(t, len(slugs), 2, "Random selection should produce variety over multiple calls")
	})
}
```

### 📋 Implementation Patterns to Follow

**Random Selection Logic:**
```go
// Import at top of file
import (
	"math/rand"
	"time"
)

// Initialize random seed (should be done once in init or main)
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Random selection from slice
randomIndex := rand.Intn(len(problems))
selectedProblem := problems[randomIndex]
```

**Flag Validation Pattern (from previous stories):**
```go
// Reuse existing validation functions
if randomDifficulty != "" && !problem.IsValidDifficulty(randomDifficulty) {
	fmt.Fprintf(os.Stderr, "Invalid difficulty '%s'. Valid options: easy, medium, hard\n", randomDifficulty)
	os.Exit(2)
}
```

**Error Handling Pattern:**
```go
// Sentinel error for no results
var ErrNoProblemsFound = errors.New("no problems found matching criteria")

// Check in command
if err == problem.ErrNoProblemsFound {
	output.PrintNoProblemsMessage(filters)
	os.Exit(2)
}
```

### 🧪 Testing Requirements

**Unit Test Coverage:**
- GetRandomProblem() with no filters (returns any unsolved problem)
- GetRandomProblem() with difficulty filter
- GetRandomProblem() with topic filter
- GetRandomProblem() with combined filters
- GetRandomProblem() when all problems solved (returns ErrNoProblemsFound)
- Randomness verification (multiple calls produce different results)

**Integration Tests:**
- CLI command with no flags
- CLI command with --difficulty flag
- CLI command with --topic flag
- CLI command with both flags
- Error message when all problems solved
- Invalid flag values (verify exit code 2)

### 🚀 Performance Requirements

**NFR Requirements:**
- Query execution <100ms (NFR: warm execution)
- Random selection is O(1) after query
- No unnecessary database queries (reuse ListProblems method)

### 📦 Dependencies

**No New Dependencies Required:**
- Uses existing problem service (ListProblems, GetProblemBySlug)
- Uses existing output formatters pattern
- Uses existing Cobra command structure
- Random selection uses standard library `math/rand`

### ⚠️ Common Pitfalls to Avoid

1. **Random Seed:** Don't forget to seed the random number generator or you'll get the same "random" result every time
2. **Empty Results:** Handle the case where no unsolved problems match filters gracefully
3. **Invalid Flags:** Validate difficulty and topic before querying database
4. **Index Out of Bounds:** Always check len(problems) > 0 before selecting random index
5. **Performance:** Don't load all problem details upfront - get the list first, then load details for selected problem only

### 🔗 Related Architecture Decisions

**From architecture.md:**
- Section: "Error Handling Strategy" - Sentinel errors (ErrNoProblemsFound)
- Section: "CLI Exit Codes" - 0 (success), 1 (general), 2 (usage), 3 (database)
- Section: "Output & Reporting" - Helpful messages, suggestions, color coding
- Section: "Usability & Developer Experience" - Actionable error messages

**From previous stories:**
- **Story 2.2**: ListProblems() method and filter pattern (reuse for random selection)
- **Story 2.3**: GetProblemBySlug() for loading problem details, PrintProblemDetails() formatting pattern
- **Common pattern**: Cobra command structure, database initialization, error handling

**NFR Requirements:**
- **NFR1**: Warm execution <100ms (query + random selection must be fast)
- **NFR28**: UNIX conventions (exit codes, stdout/stderr)
- **FR4**: Random problem selection (core feature requirement)
- **FR30**: Actionable error messages (suggest alternatives when no problems found)

### 📝 Definition of Done

- [x] cmd/random.go created with Cobra command
- [x] --difficulty and --topic flags implemented
- [x] GetRandomProblem() method added to problem service
- [x] ErrNoProblemsFound sentinel error defined
- [x] PrintRandomProblem() formatter created in output package
- [x] PrintNoProblemsMessage() helper for edge cases
- [x] Random selection logic using math/rand
- [x] Flag validation for difficulty and topic
- [x] Helpful error messages with context-aware suggestions
- [x] Exit codes: 0 (success), 2 (no problems/invalid input), 3 (database error)
- [x] Unit tests: 6+ test scenarios for GetRandomProblem()
- [x] Integration tests: Valid filters, invalid filters, edge cases
- [x] Randomness verification test
- [x] All tests pass: `go test ./...`
- [x] Manual test: `dsa random` displays random problem
- [x] Manual test: `dsa random --difficulty easy` filters correctly
- [x] Manual test: Edge case when all problems solved
- [x] All acceptance criteria satisfied

## Dev Agent Record

### Context Reference

<!-- Path(s) to story context XML will be added here by context workflow -->

### Agent Model Used

claude-sonnet-4.5 (model ID: claude-sonnet-4-5-20250929)

### Debug Log References

- All unit tests pass: 6 new tests for GetRandomProblem() (internal/problem/service_test.go:360-474)
- All integration tests pass: 6 tests for random command (cmd/random_test.go:10-117)
- Manual verification successful: `dsa random` displays formatted random problem
- Manual verification successful: `dsa random --difficulty easy` filters correctly

### Completion Notes List

**Implementation Summary:**
Successfully implemented the `dsa random` command with filter support and helpful error messages, following all acceptance criteria and TDD practices.

**Key Accomplishments:**
1. Created cmd/random.go with Cobra command structure and --difficulty/--topic flags
2. Extended problem service with GetRandomProblem() method using random selection logic
3. Implemented ErrNoProblemsFound sentinel error pattern (consistent with Story 2.3)
4. Created PrintRandomProblem() formatter with emoji header (🎲) and color-coded output
5. Implemented PrintNoProblemsMessage() with context-aware suggestions
6. Added random seed initialization in package init() function
7. Reused existing filter validation functions (IsValidDifficulty, IsValidTopic)
8. Added 6 comprehensive unit tests covering all edge cases and randomness verification
9. Added 6 integration tests verifying CLI behavior with various filter combinations
10. Manual testing confirms all acceptance criteria satisfied

**Technical Highlights:**
- Random selection uses math/rand with time-based seed (sufficient for CLI tool)
- Reuses ListProblems() method to get filtered unsolved problems, then randomly selects one
- Calls GetProblemBySlug() for full details after selection (efficient two-step approach)
- Exit codes follow UNIX conventions: 0 (success), 2 (no problems/invalid input), 3 (database error)
- All tests pass with no regressions (52 total tests in problem package, 6 integration tests)
- Output formatting uses cyan color for random-specific elements, consistent with existing color scheme
- Context-aware suggestions when no problems found (different messages for easy/medium/hard)

**Test Coverage:**
- Unit tests: 6 tests for GetRandomProblem() method
- Integration tests: 6 tests for random command behavior
- Randomness verification: Test confirms multiple calls produce variety
- All existing tests remain passing (no regressions)
- Manual verification: `dsa random` and `dsa random --difficulty easy` both work correctly

### File List

**Created:**
- cmd/random.go - Random command implementation with Cobra and flags
- cmd/random_test.go - Integration tests for random command
- internal/output/random.go - Random problem formatter and no-problems message

**Modified:**
- internal/problem/service.go - Added GetRandomProblem(), ErrNoProblemsFound, random seed init
- internal/problem/service_test.go - Added 6 unit tests for GetRandomProblem()

**Change Log:**
- 2025-12-11: Implemented Story 2.4 - Random Problem Selection command
  - Created `dsa random` command with --difficulty and --topic filters
  - Added GetRandomProblem() service method with random selection logic
  - Implemented color-coded random formatter with emoji and suggested command
  - Added context-aware error messages for no-problems scenario
  - Added 6 unit tests and 6 integration tests
  - All acceptance criteria satisfied, all tests passing

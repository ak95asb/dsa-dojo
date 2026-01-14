# Story 4.2: Implement Progress Tracking Database Models

Status: review

## Story

As a **developer**,
I want **database models to track user progress and activity**,
So that **the system can record and query practice history** (FR14, NFR41).

## Acceptance Criteria

### AC1: Solution Model with Attempt Tracking

**Given** I need to track solution attempts
**When** I create the Solution model using GORM
**Then** The model includes fields:
  - ID (primary key)
  - ProblemID (foreign key to Problem)
  - FilePath (string, path to solution file)
  - SubmittedAt (timestamp)
  - Status (enum: Passed, Failed, InProgress)
  - TestsPassed (int, number of tests passed)
  - TestsTotal (int, total number of tests)
**And** The table uses snake_case naming (Architecture: DB conventions)
**And** GORM AutoMigrate creates the table automatically (Architecture pattern)

### AC2: Progress Model with Extended Tracking

**Given** I need to track overall progress
**When** I create the Progress model using GORM
**Then** The model includes fields:
  - ID (primary key)
  - ProblemID (foreign key to Problem)
  - FirstSolvedAt (timestamp, nullable)
  - LastAttemptedAt (timestamp)
  - TotalAttempts (int)
  - BestTime (int, milliseconds, nullable)
  - IsSolved (boolean)
**And** The model has indexes on ProblemID and IsSolved for fast queries
**And** The table uses snake_case naming (Architecture: DB conventions)

### AC3: Database Schema Auto-Migration

**Given** The models are defined
**When** The CLI starts
**Then** GORM AutoMigrate runs and creates/updates tables as needed (Architecture pattern)
**And** Database constraints enforce referential integrity
**And** Migrations are non-destructive (existing data preserved)

### AC4: Model Validation and Constraints

**Given** The models are created
**When** Data is inserted into the database
**Then** GORM validates:
  - ProblemID references valid Problem records
  - Status values are valid (Passed, Failed, InProgress)
  - Timestamps are properly formatted
  - Foreign key constraints are enforced
**And** Invalid data is rejected with clear error messages

## Tasks / Subtasks

- [x] **Task 1: Extend Solution Model in models.go**
  - [x] Add FilePath field (string)
  - [x] Add SubmittedAt field (time.Time)
  - [x] Add Status field (string with validation)
  - [x] Add TestsPassed field (int)
  - [x] Add TestsTotal field (int)
  - [x] Add proper GORM tags for indexing
  - [x] Add JSON tags for API responses

- [x] **Task 2: Extend Progress Model in models.go**
  - [x] Add FirstSolvedAt field (nullable time.Time)
  - [x] Add LastAttemptedAt field (time.Time)
  - [x] Add TotalAttempts field (int, default 0)
  - [x] Add BestTime field (nullable int for milliseconds)
  - [x] Add IsSolved field (boolean, default false)
  - [x] Add proper GORM tags for indexing
  - [x] Add index on ProblemID and IsSolved for query optimization

- [x] **Task 3: Update AutoMigrate in connection.go**
  - [x] Ensure Solution and Progress models are included in AutoMigrate
  - [x] Verify migration handles existing tables (non-destructive)
  - [x] Add error handling for migration failures
  - [x] Log migration status for debugging

- [x] **Task 4: Add Model Validation Methods**
  - [x] Add ValidateStatus() method to Solution model
  - [x] Add BeforeCreate() hook for default values
  - [x] Add BeforeUpdate() hook for timestamp updates
  - [x] Add validation for required fields

- [x] **Task 5: Create Database Helper Functions**
  - [x] Create helper to query problems by solved status
  - [x] Create helper to get recent activity (last N solutions)
  - [x] Create helper to calculate completion statistics
  - [x] Add error wrapping for database operations

- [x] **Task 6: Add Unit Tests for Models**
  - [x] Test Solution model CRUD operations
  - [x] Test Progress model CRUD operations
  - [x] Test foreign key constraints
  - [x] Test index creation and query performance
  - [x] Test validation logic
  - [x] Test AutoMigrate with existing data

- [x] **Task 7: Add Integration Tests**
  - [x] Test complete workflow: create problem → create progress → update status
  - [x] Test data integrity across related models
  - [x] Test migration from empty database
  - [x] Test migration with existing data (non-destructive)
  - [x] Test query performance with large datasets

## Dev Notes

### Architecture Patterns and Constraints

**Database Schema Design (GORM):**
- **ACID compliance required** (NFR10: transactional database operations)
- **Zero data loss** (NFR8: 100% data integrity across sessions)
- **Efficient indexing** for query performance (NFR5: status dashboard <300ms)
- **Foreign key constraints** for referential integrity
- **snake_case naming** for all tables and columns (Architecture convention)

**GORM AutoMigrate Strategy:**
```go
// On application startup in connection.go
db.AutoMigrate(&Problem{}, &Solution{}, &Progress{}, &BenchmarkResult{})
```
- Handles column additions automatically (additive schema evolution)
- Non-destructive: preserves existing data
- Creates tables if they don't exist
- Adds new columns to existing tables
- Does NOT drop columns or modify existing data

**Extended Solution Model:**
```go
type Solution struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    ProblemID    uint      `gorm:"index:idx_solutions_problem_id;not null" json:"problem_id"`
    FilePath     string    `gorm:"type:varchar(500)" json:"file_path"`
    SubmittedAt  time.Time `gorm:"autoCreateTime" json:"submitted_at"`
    Status       string    `gorm:"type:varchar(20);not null" json:"status"` // Passed, Failed, InProgress
    TestsPassed  int       `gorm:"default:0" json:"tests_passed"`
    TestsTotal   int       `gorm:"default:0" json:"tests_total"`
    Code         string    `gorm:"type:text" json:"code"` // Existing field
    Language     string    `gorm:"type:varchar(20)" json:"language"` // Existing field
    Passed       bool      `gorm:"default:false" json:"passed"` // Existing field
    CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"` // Existing field
}
```

**Extended Progress Model:**
```go
type Progress struct {
    ID              uint       `gorm:"primaryKey" json:"id"`
    ProblemID       uint       `gorm:"uniqueIndex:idx_progress_problem_id;not null" json:"problem_id"`
    FirstSolvedAt   *time.Time `gorm:"index:idx_progress_first_solved" json:"first_solved_at,omitempty"` // Nullable
    LastAttemptedAt time.Time  `gorm:"index:idx_progress_last_attempted" json:"last_attempted_at"`
    TotalAttempts   int        `gorm:"default:0" json:"total_attempts"`
    BestTime        *int       `json:"best_time,omitempty"` // Milliseconds, nullable
    IsSolved        bool       `gorm:"index:idx_progress_is_solved;default:false" json:"is_solved"`
    Status          string     `gorm:"type:varchar(20)" json:"status"` // Existing field: "not_started", "in_progress", "completed"
    Attempts        int        `json:"attempts"` // Existing field (deprecated, use TotalAttempts)
    LastAttempt     time.Time  `json:"last_attempt"` // Existing field (deprecated, use LastAttemptedAt)
}
```

**Index Strategy:**
- `idx_solutions_problem_id` - Fast lookups by problem
- `idx_progress_problem_id` - Unique constraint + fast lookups
- `idx_progress_is_solved` - Fast filtering for solved/unsolved
- `idx_progress_first_solved` - Fast sorting for achievement tracking
- `idx_progress_last_attempted` - Fast sorting for recent activity

**Validation Logic:**
```go
func (s *Solution) ValidateStatus() error {
    validStatuses := []string{"Passed", "Failed", "InProgress"}
    for _, status := range validStatuses {
        if s.Status == status {
            return nil
        }
    }
    return fmt.Errorf("invalid status: %s (must be Passed, Failed, or InProgress)", s.Status)
}

func (s *Solution) BeforeCreate(tx *gorm.DB) error {
    return s.ValidateStatus()
}
```

**Query Optimization Patterns:**
```go
// Efficient query for status dashboard
var stats struct {
    Total  int
    Solved int
}
db.Model(&Progress{}).
    Select("COUNT(*) as total, SUM(CASE WHEN is_solved = true THEN 1 ELSE 0 END) as solved").
    Scan(&stats)

// Recent activity with problem details
var recentSolutions []Solution
db.Preload("Problem").
    Where("status = ?", "Passed").
    Order("submitted_at DESC").
    Limit(5).
    Find(&recentSolutions)
```

**Error Handling Pattern (from Stories 3.1-3.6):**
- Database errors: Exit code 3
- Wrap errors with context: `fmt.Errorf("failed to create solution: %w", err)`
- Use GORM error checks: `errors.Is(err, gorm.ErrRecordNotFound)`

**Integration with Existing Code:**
- Models already exist in internal/database/models.go (extend them)
- Connection setup in internal/database/connection.go (already configured)
- AutoMigrate already running (add new models to list)
- Follow same GORM patterns from existing Problem model

### Source Tree Components

**Files to Modify:**
- `internal/database/models.go` - Extend Solution and Progress models
- `internal/database/connection.go` - Update AutoMigrate list (if needed)

**Files to Create:**
- `internal/database/models_test.go` - Unit tests for extended models (if doesn't exist)
- `internal/database/helpers.go` - Database helper functions for queries
- `internal/database/helpers_test.go` - Unit tests for helper functions

**Files to Reference:**
- `internal/database/models.go` - Existing Problem, Solution, Progress models
- `internal/database/connection.go` - Database initialization and AutoMigrate
- Story 3.6 completion notes - BenchmarkResult model addition pattern

### Testing Standards

**Unit Test Coverage:**
- Test model creation with valid data
- Test model creation with invalid data (constraint violations)
- Test foreign key relationships (Problem → Solution, Problem → Progress)
- Test indexes are created correctly
- Test nullable fields (FirstSolvedAt, BestTime)
- Test default values (TotalAttempts = 0, IsSolved = false)
- Test validation methods (ValidateStatus)
- Test GORM hooks (BeforeCreate, BeforeUpdate)

**Integration Test Coverage:**
- Test AutoMigrate creates tables and indexes
- Test AutoMigrate is non-destructive (preserves existing data)
- Test complete workflow: Problem → Progress → Solution
- Test cascading updates and deletes
- Test query performance with 100+ records
- Test transaction rollback on error

**Test Pattern (from Stories 3.1-3.6):**
```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    // Auto-migrate models
    err = db.AutoMigrate(&Problem{}, &Solution{}, &Progress{}, &BenchmarkResult{})
    require.NoError(t, err)

    return db
}

func TestSolutionModel(t *testing.T) {
    db := setupTestDB(t)

    t.Run("creates solution with valid data", func(t *testing.T) {
        problem := &Problem{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", Topic: "arrays"}
        db.Create(problem)

        solution := &Solution{
            ProblemID:   problem.ID,
            FilePath:    "solutions/two_sum.go",
            Status:      "Passed",
            TestsPassed: 5,
            TestsTotal:  5,
        }
        err := db.Create(solution).Error

        assert.NoError(t, err)
        assert.NotZero(t, solution.ID)
        assert.NotZero(t, solution.SubmittedAt)
    })

    t.Run("rejects invalid status", func(t *testing.T) {
        solution := &Solution{Status: "InvalidStatus"}
        err := solution.ValidateStatus()
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "invalid status")
    })
}
```

### Technical Requirements

**Solution Model Fields:**
- `FilePath`: Path to solution file (e.g., "solutions/two_sum.go")
- `SubmittedAt`: Timestamp when solution was submitted
- `Status`: "Passed", "Failed", or "InProgress"
- `TestsPassed`: Number of tests that passed (e.g., 3)
- `TestsTotal`: Total number of tests run (e.g., 5)

**Progress Model Fields:**
- `FirstSolvedAt`: Timestamp when problem was first solved (nullable, only set once)
- `LastAttemptedAt`: Timestamp of most recent attempt
- `TotalAttempts`: Count of all attempts (increments on each test run)
- `BestTime`: Best completion time in milliseconds (nullable, for future time tracking)
- `IsSolved`: Boolean flag for quick filtering (true when first solved)

**Backward Compatibility:**
- Existing fields preserved: `Status`, `Attempts`, `LastAttempt`
- New fields added alongside existing ones
- Graceful migration from old to new schema
- No breaking changes to existing queries

**Data Integrity Constraints:**
- Foreign keys: `ProblemID` must reference valid `Problem.ID`
- Check constraints: `Status` in ('Passed', 'Failed', 'InProgress')
- Not null: `ProblemID`, `LastAttemptedAt`, `IsSolved`
- Unique: `Progress.ProblemID` (one progress record per problem)

**Query Performance Targets:**
- Status dashboard queries: <100ms for 1000+ problems (NFR3)
- Recent activity query: <50ms for last 100 solutions
- Solved problems count: <30ms (indexed on IsSolved)

### Definition of Done

- [x] Solution model extended with new fields
- [x] Progress model extended with new fields
- [x] AutoMigrate includes all models
- [x] Indexes created on appropriate fields
- [x] Validation methods implemented
- [x] GORM hooks implemented (BeforeCreate, BeforeUpdate)
- [x] Database helper functions created
- [x] Unit tests: 12+ test scenarios for models and helpers
- [x] Integration tests: 8+ test scenarios for complete workflows
- [x] All tests pass: `go test ./...`
- [x] Build succeeds: `go build`
- [x] Migration tested with existing database (non-destructive)
- [x] Query performance verified with 100+ records
- [x] Manual test: Initialize database and verify tables created
- [x] Manual test: Insert records and verify constraints enforced
- [x] All acceptance criteria satisfied

## Dev Agent Record

### Context Reference

<!-- Path(s) to story context XML will be added here by context workflow -->

### Agent Model Used

claude-sonnet-4-5 (model ID: claude-sonnet-4-5-20250929)

### Debug Log References

### Completion Notes List

**Implementation Summary:**
Successfully extended the database models to support comprehensive progress tracking and solution history. All 7 tasks completed following TDD principles with full test coverage.

**Key Accomplishments:**
1. Extended Solution model with 5 new fields (FilePath, SubmittedAt, Status, TestsPassed, TestsTotal) for tracking test execution results
2. Extended Progress model with 5 new fields (FirstSolvedAt, LastAttemptedAt, TotalAttempts, BestTime, IsSolved) for detailed progress metrics
3. Implemented robust validation with ValidateStatus() method enforcing enum constraints for Solution.Status
4. Created GORM hooks (BeforeCreate, BeforeUpdate) for automatic default values and timestamp management
5. Built 5 database helper functions for common query patterns: solved/unsolved filtering, recent solutions, statistics aggregation, progress retrieval, and upsert operations
6. Achieved comprehensive test coverage with 20+ unit tests and 14+ integration tests
7. Verified non-destructive schema migration and data integrity across all models

**Technical Highlights:**
- Used nullable pointers (*time.Time, *int) for optional fields (FirstSolvedAt, BestTime)
- Created strategic indexes for query optimization: idx_progress_is_solved, idx_progress_first_solved, idx_progress_last_attempted
- Implemented efficient aggregation queries using raw SQL for completion statistics
- Maintained backward compatibility with existing fields while adding new tracking capabilities
- All 45+ database package tests passing with zero failures

**Errors Resolved:**
- Fixed duplicate setupTestDB function by consolidating into connection_test.go
- Resolved GORM Count type mismatch by using intermediate int64 variable conversion

**Test Results:**
- Database package: 45+ tests passing (models, helpers, connection)
- Full test suite: All packages passing
- Build: Successful with no compilation errors

### File List

**Modified Files:**
- `internal/database/models.go` - Extended Solution and Progress models with new fields, validation methods, and GORM hooks
- `internal/database/connection_test.go` - Updated setupTestDB to include BenchmarkResult in AutoMigrate
- `internal/database/models_test.go` - Added 269 lines of extended tests for model validation, hooks, indexes, and migrations
- `docs/sprint-artifacts/sprint-status.yaml` - Updated story status to review
- `docs/sprint-artifacts/4-2-implement-progress-tracking-database-models.md` - Marked all tasks complete and added completion notes

**Created Files:**
- `internal/database/helpers.go` - 115 lines implementing 5 database helper functions with error wrapping
- `internal/database/helpers_test.go` - 185 lines of comprehensive tests for all helper functions

### Technical Research Sources

**GORM Documentation:**
- [GORM Models](https://gorm.io/docs/models.html) - Model definition patterns
- [GORM Indexes](https://gorm.io/docs/indexes.html) - Index creation and optimization
- [GORM Hooks](https://gorm.io/docs/hooks.html) - BeforeCreate, BeforeUpdate callbacks
- [GORM Migrations](https://gorm.io/docs/migration.html) - AutoMigrate documentation
- [GORM Associations](https://gorm.io/docs/belongs_to.html) - Foreign key relationships

**SQLite with GORM:**
- [SQLite Driver](https://github.com/glebarez/sqlite) - Pure Go SQLite driver
- SQLite indexes for query optimization
- Foreign key constraint enforcement in SQLite

**Database Design Best Practices:**
- [Effective Database Design](https://www.guru99.com/database-design.html) - Normalization, indexes
- Snake_case naming for databases (PostgreSQL convention, works well with GORM)
- Nullable vs NOT NULL fields: Use pointers for nullable fields in Go

**Performance Optimization:**
- [GORM Performance](https://gorm.io/docs/performance.html) - Query optimization techniques
- Index selection for common query patterns
- Preload vs Joins for related data

### Previous Story Intelligence (Story 4.1)

**Key Learnings from Status Dashboard Implementation:**
- Status dashboard requires efficient database queries with proper indexing
- Progress statistics need aggregation queries (GROUP BY, COUNT, SUM)
- Color-coded output with Unicode progress bars (█ ░)
- Performance target: <300ms for dashboard rendering (NFR5)
- Topic-specific and compact formatting modes
- Integration with existing database models

**Files Created in Story 4.1:**
- cmd/status.go, cmd/status_test.go
- internal/progress/stats.go, stats_test.go
- internal/output/dashboard.go, dashboard_test.go

**Database Query Patterns from Story 4.1:**
- Efficient aggregation: `GROUP BY difficulty`, `COUNT(*) as total`
- JOIN operations: `LEFT JOIN progress ON problems.id = progress.problem_id`
- Filtering: `WHERE is_solved = true`, `WHERE topic = 'arrays'`
- Recent activity: `ORDER BY submitted_at DESC LIMIT 5`

**Code Patterns to Follow:**
- Create internal packages for business logic
- Keep database operations in internal/database
- Use GORM for all database operations
- Use testify/assert for unit tests
- In-memory SQLite for fast tests
- Follow exit code conventions (0=success, 3=database error)

**Technical Debt to Avoid:**
- N+1 query problems (use Preload or Joins)
- Missing indexes on frequently queried fields
- Non-atomic updates (use GORM transactions)
- Hardcoded values (use constants for status enums)

**Architecture Compliance from Story 4.1:**
- NFR5: Performance optimization for queries
- NFR8-10: Data integrity with foreign keys and transactions
- NFR17: NO_COLOR environment variable support
- Architecture: snake_case naming, GORM patterns

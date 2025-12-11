// Package database provides GORM models and database connection management
// for the DSA CLI application. It handles SQLite database operations for
// storing problems, solutions, and progress tracking data.
package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Problem represents a DSA problem in the problem library.
// Each problem has a unique slug identifier and metadata about
// difficulty, topic category, and problem description.
type Problem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Slug        string    `gorm:"uniqueIndex:idx_problems_slug;not null" json:"slug"`
	Title       string    `gorm:"not null" json:"title"`
	Difficulty  string    `gorm:"type:varchar(20);not null" json:"difficulty"` // easy, medium, hard
	Topic       string    `gorm:"type:varchar(50)" json:"topic"`               // arrays, trees, etc.
	Description string    `gorm:"type:text" json:"description"`
	Tags        string    `gorm:"type:varchar(255)" json:"tags"` // Comma-separated tags (e.g., "bfs,dfs,recursion")
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Solution represents a developer's solution attempt for a problem.
// Multiple solutions can exist for the same problem, tracking code,
// language, test results, and submission details.
type Solution struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProblemID   uint      `gorm:"index:idx_solutions_problem_id;not null" json:"problem_id"`
	Code        string    `gorm:"type:text" json:"code"`
	Language    string    `gorm:"type:varchar(20);default:'go'" json:"language"`
	Passed      bool      `gorm:"default:false" json:"passed"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	FilePath    string    `gorm:"type:varchar(500)" json:"file_path"`
	SubmittedAt time.Time `gorm:"autoCreateTime" json:"submitted_at"`
	Status      string    `gorm:"type:varchar(20);not null;default:'InProgress'" json:"status"` // Passed, Failed, InProgress
	TestsPassed int       `gorm:"default:0" json:"tests_passed"`
	TestsTotal  int       `gorm:"default:0" json:"tests_total"`
}

// Progress tracks a developer's progress on each problem.
// Only one progress record exists per problem, maintaining
// current status, attempt count, solved timestamp, and performance metrics.
type Progress struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	ProblemID       uint       `gorm:"uniqueIndex:idx_progress_problem_id;not null" json:"problem_id"`
	Status          string     `gorm:"type:varchar(20);default:'not_started'" json:"status"` // not_started, in_progress, completed
	Attempts        int        `gorm:"default:0" json:"attempts"`
	LastAttempt     time.Time  `gorm:"" json:"last_attempt"`
	FirstSolvedAt   *time.Time `gorm:"index:idx_progress_first_solved" json:"first_solved_at,omitempty"`
	LastAttemptedAt time.Time  `gorm:"index:idx_progress_last_attempted" json:"last_attempted_at"`
	TotalAttempts   int        `gorm:"default:0" json:"total_attempts"`
	BestTime        *int       `json:"best_time,omitempty"` // Milliseconds, nullable
	IsSolved        bool       `gorm:"index:idx_progress_is_solved;default:false" json:"is_solved"`
}

// BenchmarkResult represents a benchmark run result for a problem solution.
// Stores performance metrics including timing and memory allocations.
type BenchmarkResult struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProblemID   uint      `gorm:"index:idx_benchmarks_problem_id;not null" json:"problem_id"`
	NsPerOp     float64   `gorm:"not null" json:"ns_per_op"`     // Nanoseconds per operation
	AllocsPerOp float64   `gorm:"not null" json:"allocs_per_op"` // Allocations per operation
	BytesPerOp  float64   `gorm:"not null" json:"bytes_per_op"`  // Bytes allocated per operation
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// ValidateStatus checks if the Solution status is valid
func (s *Solution) ValidateStatus() error {
	validStatuses := []string{"Passed", "Failed", "InProgress"}
	for _, status := range validStatuses {
		if s.Status == status {
			return nil
		}
	}
	return fmt.Errorf("invalid status: %s (must be Passed, Failed, or InProgress)", s.Status)
}

// BeforeCreate hook validates data before creating a Solution record
func (s *Solution) BeforeCreate(tx *gorm.DB) error {
	if s.Status == "" {
		s.Status = "InProgress"
	}
	return s.ValidateStatus()
}

// BeforeUpdate hook validates data before updating a Solution record
func (s *Solution) BeforeUpdate(tx *gorm.DB) error {
	if s.Status != "" {
		return s.ValidateStatus()
	}
	return nil
}

// BeforeCreate hook sets default values for Progress record
func (p *Progress) BeforeCreate(tx *gorm.DB) error {
	// LastAttemptedAt defaults to now if not set
	if p.LastAttemptedAt.IsZero() {
		p.LastAttemptedAt = time.Now()
	}
	return nil
}

// BeforeUpdate hook updates timestamps for Progress record
func (p *Progress) BeforeUpdate(tx *gorm.DB) error {
	// Update LastAttemptedAt on each update
	p.LastAttemptedAt = time.Now()
	return nil
}

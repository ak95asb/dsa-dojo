package testing

import (
	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"gorm.io/gorm"
)

// TestResult contains the results of test execution
type TestResult struct {
	AllPassed    bool
	PassedCount  int
	TotalCount   int
	FailedTests  []FailedTest
	Output       string
	Verbose      bool
	RaceDetector bool
}

// FailedTest represents a single failed test case
type FailedTest struct {
	Name     string
	Expected string
	Actual   string
	Message  string
}

// Service provides test execution and management operations
type Service struct {
	db        *gorm.DB
	executor  *Executor
	formatter *Formatter
}

// NewService creates a new testing service
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:        db,
		executor:  NewExecutor(),
		formatter: NewFormatter(),
	}
}

// ExecuteTests runs go test for the specified problem
func (s *Service) ExecuteTests(prob *problem.ProblemDetails, verbose, race bool) (*TestResult, error) {
	return s.executor.Execute(prob, verbose, race)
}

// DisplayResults formats and displays test results to stdout
func (s *Service) DisplayResults(result *TestResult) {
	s.formatter.Display(result)
}

// RecordSolution creates or updates a solution record in the database
func (s *Service) RecordSolution(problemID uint, result *TestResult) error {
	// Check if solution already exists
	var existing database.Solution
	err := s.db.Where("problem_id = ?", problemID).Order("created_at DESC").First(&existing).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// Create new solution record
	solution := &database.Solution{
		ProblemID: problemID,
		Code:      "", // We don't store the actual code in test command
		Language:  "go",
		Passed:    result.AllPassed,
	}

	return s.db.Create(solution).Error
}

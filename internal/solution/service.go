package solution

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"gorm.io/gorm"
)

// Service provides solution-related operations
type Service struct {
	db        *gorm.DB
	generator *Generator
}

// NewService creates a new solution service
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:        db,
		generator: NewGenerator(),
	}
}

// GenerateSolution creates a solution file for the given problem
func (s *Service) GenerateSolution(p *database.Problem, force bool) (string, error) {
	return s.generator.GenerateSolution(p, force)
}

// SubmissionRecord represents a solution submission with metadata
type SubmissionRecord struct {
	ID          uint
	ProblemID   uint
	Code        string
	Language    string
	Passed      bool
	CreatedAt   time.Time
	TestsPassed int
	TestsTotal  int
}

// RecordSubmission saves a solution to history directory and database
// Creates directory structure: solutions/history/<problem-slug>/<timestamp>.go
func (s *Service) RecordSubmission(problemSlug string, problemID uint, solutionPath string, passed bool, testsPassed, testsTotal int) (*SubmissionRecord, error) {
	// Read solution code
	code, err := os.ReadFile(solutionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read solution file: %w", err)
	}

	// Create history directory
	historyDir := filepath.Join("solutions", "history", problemSlug)
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create history directory: %w", err)
	}

	// Generate timestamp-based filename
	timestamp := time.Now().Format("20060102-150405")
	historyPath := filepath.Join(historyDir, fmt.Sprintf("%s.go", timestamp))

	// Copy solution to history
	if err := os.WriteFile(historyPath, code, 0644); err != nil {
		return nil, fmt.Errorf("failed to write history file: %w", err)
	}

	// Create database record
	solution := &database.Solution{
		ProblemID: problemID,
		Code:      string(code),
		Language:  "go",
		Passed:    passed,
	}

	if err := s.db.Create(solution).Error; err != nil {
		return nil, fmt.Errorf("failed to create database record: %w", err)
	}

	return &SubmissionRecord{
		ID:          solution.ID,
		ProblemID:   solution.ProblemID,
		Code:        solution.Code,
		Language:    solution.Language,
		Passed:      solution.Passed,
		CreatedAt:   solution.CreatedAt,
		TestsPassed: testsPassed,
		TestsTotal:  testsTotal,
	}, nil
}

// GetHistory retrieves all solution submissions for a problem
// Returns submissions sorted by most recent first
func (s *Service) GetHistory(problemID uint) ([]SubmissionRecord, error) {
	var solutions []database.Solution

	err := s.db.Where("problem_id = ?", problemID).
		Order("created_at DESC").
		Find(&solutions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query solution history: %w", err)
	}

	// Convert to SubmissionRecord format
	records := make([]SubmissionRecord, len(solutions))
	for i, sol := range solutions {
		records[i] = SubmissionRecord{
			ID:        sol.ID,
			ProblemID: sol.ProblemID,
			Code:      sol.Code,
			Language:  sol.Language,
			Passed:    sol.Passed,
			CreatedAt: sol.CreatedAt,
		}
	}

	return records, nil
}

// GetSubmissionByIndex retrieves a submission by 1-based index (1 = most recent)
func (s *Service) GetSubmissionByIndex(problemID uint, index int) (*SubmissionRecord, error) {
	if index < 1 {
		return nil, fmt.Errorf("index must be >= 1")
	}

	var solution database.Solution

	// Query with LIMIT 1 OFFSET (index-1)
	err := s.db.Where("problem_id = ?", problemID).
		Order("created_at DESC").
		Offset(index - 1).
		Limit(1).
		First(&solution).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("submission #%d not found", index)
		}
		return nil, fmt.Errorf("failed to query submission: %w", err)
	}

	return &SubmissionRecord{
		ID:        solution.ID,
		ProblemID: solution.ProblemID,
		Code:      solution.Code,
		Language:  solution.Language,
		Passed:    solution.Passed,
		CreatedAt: solution.CreatedAt,
	}, nil
}

// BackupCurrentSolution creates a backup of the current solution
func (s *Service) BackupCurrentSolution(problemSlug string, problemID uint) (string, error) {
	solutionPath := filepath.Join("solutions", fmt.Sprintf("%s.go", problemSlug))

	// Check if solution file exists
	if _, err := os.Stat(solutionPath); os.IsNotExist(err) {
		return "", nil // No current solution to backup
	}

	// Create backup with timestamp
	timestamp := time.Now().Format("20060102-150405")
	historyDir := filepath.Join("solutions", "history", problemSlug)
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create history directory: %w", err)
	}

	backupPath := filepath.Join(historyDir, fmt.Sprintf("backup-%s.go", timestamp))

	// Read and write backup
	content, err := os.ReadFile(solutionPath)
	if err != nil {
		return "", fmt.Errorf("failed to read current solution: %w", err)
	}

	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup: %w", err)
	}

	return backupPath, nil
}

// RestoreSolution restores a submission as the current solution
func (s *Service) RestoreSolution(problemSlug string, record *SubmissionRecord) error {
	solutionPath := filepath.Join("solutions", fmt.Sprintf("%s.go", problemSlug))

	// Create solutions directory if doesn't exist
	if err := os.MkdirAll("solutions", 0755); err != nil {
		return fmt.Errorf("failed to create solutions directory: %w", err)
	}

	// Write solution code
	if err := os.WriteFile(solutionPath, []byte(record.Code), 0644); err != nil {
		return fmt.Errorf("failed to write solution: %w", err)
	}

	return nil
}

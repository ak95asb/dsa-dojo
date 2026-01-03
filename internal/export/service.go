package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"gorm.io/gorm"
)

// ExportService handles exporting progress data to various formats
type ExportService struct {
	db *gorm.DB
}

// ExportFilter specifies filtering criteria for exports
type ExportFilter struct {
	Topic      string
	Difficulty string
}

// ExportData represents the complete export structure for JSON
type ExportData struct {
	ExportedAt time.Time       `json:"exported_at"`
	Version    string          `json:"version"`
	Summary    ExportSummary   `json:"summary"`
	Problems   []ProblemExport `json:"problems"`
}

// ExportSummary contains aggregate statistics
type ExportSummary struct {
	TotalProblems      int     `json:"total_problems"`
	ProblemsSolved     int     `json:"problems_solved"`
	OverallSuccessRate float64 `json:"overall_success_rate"`
	AvgAttempts        float64 `json:"avg_attempts"`
}

// ProblemExport represents a problem with its progress and solutions
type ProblemExport struct {
	Slug       string           `json:"slug"`
	Title      string           `json:"title"`
	Difficulty string           `json:"difficulty"`
	Topic      string           `json:"topic"`
	Progress   ProgressExport   `json:"progress"`
	Solutions  []SolutionExport `json:"solutions"`
}

// ProgressExport represents progress data for export
type ProgressExport struct {
	IsSolved        bool       `json:"is_solved"`
	TotalAttempts   int        `json:"total_attempts"`
	FirstSolvedAt   *time.Time `json:"first_solved_at,omitempty"`
	LastAttemptedAt time.Time  `json:"last_attempted_at"`
}

// SolutionExport represents solution data for export
type SolutionExport struct {
	SubmittedAt time.Time `json:"submitted_at"`
	Status      string    `json:"status"`
	TestsPassed int       `json:"tests_passed"`
	TestsTotal  int       `json:"tests_total"`
}

// NewService creates a new export service instance
func NewService(db *gorm.DB) *ExportService {
	return &ExportService{db: db}
}

// ExportToJSON exports progress data to JSON format
func (s *ExportService) ExportToJSON(filter ExportFilter, writer io.Writer) error {
	// Gather export data
	data, err := s.gatherExportData(filter)
	if err != nil {
		return fmt.Errorf("failed to gather export data: %w", err)
	}

	// Encode to JSON with indentation
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// ExportToCSV exports progress data to CSV format
func (s *ExportService) ExportToCSV(filter ExportFilter, writer io.Writer) error {
	// Query problems with progress
	problems, err := s.queryProblemsWithProgress(filter)
	if err != nil {
		return fmt.Errorf("failed to query problems: %w", err)
	}

	// Write CSV
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	header := []string{"Slug", "Title", "Difficulty", "Topic", "IsSolved", "TotalAttempts", "FirstSolvedAt", "LastAttemptedAt"}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, problem := range problems {
		row := []string{
			problem.Slug,
			problem.Title,
			problem.Difficulty,
			problem.Topic,
			fmt.Sprintf("%t", problem.Progress.IsSolved),
			fmt.Sprintf("%d", problem.Progress.TotalAttempts),
			formatTimestamp(problem.Progress.FirstSolvedAt),
			formatTimestamp(&problem.Progress.LastAttemptedAt),
		}
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// gatherExportData collects all data for export
func (s *ExportService) gatherExportData(filter ExportFilter) (*ExportData, error) {
	// Query problems with progress and solutions
	problems, err := s.queryProblemsWithProgress(filter)
	if err != nil {
		return nil, err
	}

	// Calculate summary statistics
	summary := s.calculateSummary(problems)

	// Build export data
	exportProblems := make([]ProblemExport, 0, len(problems))
	for _, problem := range problems {
		exportProblem := ProblemExport{
			Slug:       problem.Slug,
			Title:      problem.Title,
			Difficulty: problem.Difficulty,
			Topic:      problem.Topic,
			Progress: ProgressExport{
				IsSolved:        problem.Progress.IsSolved,
				TotalAttempts:   problem.Progress.TotalAttempts,
				FirstSolvedAt:   problem.Progress.FirstSolvedAt,
				LastAttemptedAt: problem.Progress.LastAttemptedAt,
			},
			Solutions: make([]SolutionExport, 0, len(problem.Solutions)),
		}

		// Add solutions
		for _, solution := range problem.Solutions {
			exportProblem.Solutions = append(exportProblem.Solutions, SolutionExport{
				SubmittedAt: solution.SubmittedAt,
				Status:      solution.Status,
				TestsPassed: solution.TestsPassed,
				TestsTotal:  solution.TestsTotal,
			})
		}

		exportProblems = append(exportProblems, exportProblem)
	}

	return &ExportData{
		ExportedAt: time.Now(),
		Version:    "1.0",
		Summary:    summary,
		Problems:   exportProblems,
	}, nil
}

// ProblemWithProgress represents a problem with related data
type ProblemWithProgress struct {
	database.Problem
	Progress  database.Progress
	Solutions []database.Solution
}

// queryProblemsWithProgress queries problems with progress and solutions
func (s *ExportService) queryProblemsWithProgress(filter ExportFilter) ([]ProblemWithProgress, error) {
	var problems []database.Problem

	// Build query with filters
	query := s.db.Model(&database.Problem{})

	if filter.Difficulty != "" {
		query = query.Where("difficulty = ?", filter.Difficulty)
	}

	if filter.Topic != "" {
		query = query.Where("topic = ?", filter.Topic)
	}

	// Get all matching problems
	if err := query.Find(&problems).Error; err != nil {
		return nil, err
	}

	// Build results with progress and solutions for each problem
	results := make([]ProblemWithProgress, 0, len(problems))
	for _, problem := range problems {
		var progress database.Progress
		var solutions []database.Solution

		// Get progress (may not exist)
		s.db.Where("problem_id = ?", problem.ID).First(&progress)

		// Get solutions (may be empty)
		s.db.Where("problem_id = ?", problem.ID).Find(&solutions)

		results = append(results, ProblemWithProgress{
			Problem:   problem,
			Progress:  progress,
			Solutions: solutions,
		})
	}

	return results, nil
}

// calculateSummary computes summary statistics
func (s *ExportService) calculateSummary(problems []ProblemWithProgress) ExportSummary {
	summary := ExportSummary{}

	summary.TotalProblems = len(problems)

	solvedCount := 0
	totalAttempts := 0

	for _, problem := range problems {
		if problem.Progress.IsSolved {
			solvedCount++
			totalAttempts += problem.Progress.TotalAttempts
		}
	}

	summary.ProblemsSolved = solvedCount

	if summary.TotalProblems > 0 {
		summary.OverallSuccessRate = (float64(solvedCount) / float64(summary.TotalProblems)) * 100
	}

	if solvedCount > 0 {
		summary.AvgAttempts = float64(totalAttempts) / float64(solvedCount)
	}

	return summary
}

// formatTimestamp formats a timestamp for CSV output
func formatTimestamp(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

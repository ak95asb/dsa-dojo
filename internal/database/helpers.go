package database

import (
	"fmt"

	"gorm.io/gorm"
)

// GetProblemsBySolvedStatus queries problems filtered by solved status
func GetProblemsBySolvedStatus(db *gorm.DB, solved bool) ([]Problem, error) {
	var problems []Problem

	query := db.Table("problems").
		Select("problems.*").
		Joins("LEFT JOIN progresses ON problems.id = progresses.problem_id").
		Where("progresses.is_solved = ?", solved)

	if err := query.Find(&problems).Error; err != nil {
		return nil, fmt.Errorf("failed to query problems by solved status: %w", err)
	}

	return problems, nil
}

// GetRecentSolutions retrieves the last N solution submissions
func GetRecentSolutions(db *gorm.DB, limit int) ([]Solution, error) {
	var solutions []Solution

	if err := db.Order("submitted_at DESC").Limit(limit).Find(&solutions).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent solutions: %w", err)
	}

	return solutions, nil
}

// CompletionStats represents aggregate statistics about problem completion
type CompletionStats struct {
	TotalProblems   int
	SolvedProblems  int
	TotalAttempts   int
	AverageAttempts float64
}

// GetCompletionStatistics calculates overall completion statistics
func GetCompletionStatistics(db *gorm.DB) (*CompletionStats, error) {
	stats := &CompletionStats{}

	// Count total problems
	var totalCount int64
	if err := db.Model(&Problem{}).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count total problems: %w", err)
	}
	stats.TotalProblems = int(totalCount)

	// Count solved problems and total attempts
	var result struct {
		SolvedCount  int
		TotalAttempts int
	}

	query := `
		SELECT
			COUNT(CASE WHEN is_solved = true THEN 1 END) as solved_count,
			SUM(total_attempts) as total_attempts
		FROM progresses
	`

	if err := db.Raw(query).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate completion stats: %w", err)
	}

	stats.SolvedProblems = result.SolvedCount
	stats.TotalAttempts = result.TotalAttempts

	// Calculate average attempts per problem
	if stats.SolvedProblems > 0 {
		stats.AverageAttempts = float64(stats.TotalAttempts) / float64(stats.SolvedProblems)
	}

	return stats, nil
}

// GetProgressByProblemID retrieves progress record for a specific problem
func GetProgressByProblemID(db *gorm.DB, problemID uint) (*Progress, error) {
	var progress Progress

	if err := db.Where("problem_id = ?", problemID).First(&progress).Error; err != nil {
		return nil, fmt.Errorf("failed to get progress for problem %d: %w", problemID, err)
	}

	return &progress, nil
}

// UpsertProgress creates or updates a progress record for a problem
func UpsertProgress(db *gorm.DB, progress *Progress) error {
	var existing Progress
	err := db.Where("problem_id = ?", progress.ProblemID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new progress record
		if err := db.Create(progress).Error; err != nil {
			return fmt.Errorf("failed to create progress record: %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check existing progress: %w", err)
	}

	// Update existing record
	progress.ID = existing.ID
	if err := db.Save(progress).Error; err != nil {
		return fmt.Errorf("failed to update progress record: %w", err)
	}

	return nil
}

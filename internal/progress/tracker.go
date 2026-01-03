package progress

import (
	"fmt"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"gorm.io/gorm"
)

// Tracker manages progress tracking for problem solutions
type Tracker struct {
	db *gorm.DB
}

// NewTracker creates a new progress tracker instance
func NewTracker(db *gorm.DB) *Tracker {
	return &Tracker{db: db}
}

// TrackTestCompletion records test results and updates progress tracking.
// It creates/updates Progress and Solution records atomically in a transaction.
// Returns true if this is the first time the problem was solved, false otherwise.
func (t *Tracker) TrackTestCompletion(problemID uint, filePath string, passed bool, testsPassed, testsTotal int) (bool, error) {
	var isFirstTimeSolve bool

	err := t.db.Transaction(func(tx *gorm.DB) error {
		// 1. Verify problem exists (ensures foreign key integrity)
		var problem database.Problem
		err := tx.First(&problem, problemID).Error
		if err != nil {
			return fmt.Errorf("problem not found: %w", err)
		}

		// 2. Get or create Progress record
		var progress database.Progress
		err = tx.Where("problem_id = ?", problemID).FirstOrCreate(&progress, database.Progress{
			ProblemID: problemID,
		}).Error
		if err != nil {
			return fmt.Errorf("failed to get progress: %w", err)
		}

		// 3. Check if this is the first time solving
		isFirstTimeSolve = passed && !progress.IsSolved

		// 4. Prepare updates for Progress record
		now := time.Now()
		updates := map[string]interface{}{
			"last_attempted_at": now,
			"total_attempts":    gorm.Expr("total_attempts + ?", 1),
		}

		// If passed and not previously solved, mark as solved
		if passed && !progress.IsSolved {
			updates["is_solved"] = true
			updates["first_solved_at"] = now
		}

		// 5. Apply updates atomically
		err = tx.Model(&progress).Updates(updates).Error
		if err != nil {
			return fmt.Errorf("failed to update progress: %w", err)
		}

		// 6. Create Solution record
		status := "Failed"
		if passed {
			status = "Passed"
		}

		solution := &database.Solution{
			ProblemID:   problemID,
			FilePath:    filePath,
			Status:      status,
			TestsPassed: testsPassed,
			TestsTotal:  testsTotal,
		}

		err = tx.Create(solution).Error
		if err != nil {
			return fmt.Errorf("failed to create solution: %w", err)
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return isFirstTimeSolve, nil
}

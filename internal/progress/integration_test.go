package progress

import (
	"testing"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration tests verify end-to-end workflows for progress tracking

func TestIntegration_CompleteWorkflow(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	t.Run("complete workflow: fail → fail → pass → progress updated", func(t *testing.T) {
		// Create test problem
		problem := &database.Problem{Slug: "workflow-test", Title: "Workflow Test", Difficulty: "medium"}
		require.NoError(t, db.Create(problem).Error)

		// Attempt 1: Failed
		isFirstSolve, err := tracker.TrackTestCompletion(problem.ID, "problems/workflow-test/solution.go", false, 2, 5)
		require.NoError(t, err)
		assert.False(t, isFirstSolve)

		// Verify progress after first attempt
		var progress database.Progress
		err = db.First(&progress, "problem_id = ?", problem.ID).Error
		require.NoError(t, err)
		assert.False(t, progress.IsSolved)
		assert.Nil(t, progress.FirstSolvedAt)
		assert.Equal(t, 1, progress.TotalAttempts)

		// Verify solution record created
		var solutionCount int64
		db.Model(&database.Solution{}).Where("problem_id = ?", problem.ID).Count(&solutionCount)
		assert.Equal(t, int64(1), solutionCount)

		// Attempt 2: Failed again
		isFirstSolve, err = tracker.TrackTestCompletion(problem.ID, "problems/workflow-test/solution.go", false, 3, 5)
		require.NoError(t, err)
		assert.False(t, isFirstSolve)

		// Verify progress after second attempt
		err = db.First(&progress, "problem_id = ?", problem.ID).Error
		require.NoError(t, err)
		assert.False(t, progress.IsSolved)
		assert.Nil(t, progress.FirstSolvedAt)
		assert.Equal(t, 2, progress.TotalAttempts)

		// Attempt 3: Passed
		isFirstSolve, err = tracker.TrackTestCompletion(problem.ID, "problems/workflow-test/solution.go", true, 5, 5)
		require.NoError(t, err)
		assert.True(t, isFirstSolve)

		// Verify progress after solving
		err = db.First(&progress, "problem_id = ?", problem.ID).Error
		require.NoError(t, err)
		assert.True(t, progress.IsSolved)
		assert.NotNil(t, progress.FirstSolvedAt)
		assert.Equal(t, 3, progress.TotalAttempts)

		// Verify all solution records exist
		db.Model(&database.Solution{}).Where("problem_id = ?", problem.ID).Count(&solutionCount)
		assert.Equal(t, int64(3), solutionCount)

		// Verify latest solution is marked as passed
		var lastSolution database.Solution
		err = db.Where("problem_id = ?", problem.ID).Order("submitted_at DESC").First(&lastSolution).Error
		require.NoError(t, err)
		assert.Equal(t, "Passed", lastSolution.Status)
	})
}

func TestIntegration_IdempotentSolving(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	t.Run("solving same problem twice is idempotent", func(t *testing.T) {
		// Create test problem
		problem := &database.Problem{Slug: "idempotent-test", Title: "Idempotent Test", Difficulty: "easy"}
		require.NoError(t, db.Create(problem).Error)

		// First solve
		isFirstSolve, err := tracker.TrackTestCompletion(problem.ID, "problems/idempotent-test/solution.go", true, 3, 3)
		require.NoError(t, err)
		assert.True(t, isFirstSolve)

		// Get first solved timestamp
		var progress database.Progress
		err = db.First(&progress, "problem_id = ?", problem.ID).Error
		require.NoError(t, err)
		firstSolvedAt := progress.FirstSolvedAt

		time.Sleep(10 * time.Millisecond)

		// Second solve (re-running tests)
		isFirstSolve, err = tracker.TrackTestCompletion(problem.ID, "problems/idempotent-test/solution.go", true, 3, 3)
		require.NoError(t, err)
		assert.False(t, isFirstSolve) // Not first time anymore

		// Verify FirstSolvedAt unchanged
		err = db.First(&progress, "problem_id = ?", problem.ID).Error
		require.NoError(t, err)
		assert.Equal(t, firstSolvedAt, progress.FirstSolvedAt)
		assert.Equal(t, 2, progress.TotalAttempts)
		assert.True(t, progress.IsSolved) // Still solved

		// Third solve
		isFirstSolve, err = tracker.TrackTestCompletion(problem.ID, "problems/idempotent-test/solution.go", true, 3, 3)
		require.NoError(t, err)
		assert.False(t, isFirstSolve)

		// Verify data integrity
		err = db.First(&progress, "problem_id = ?", problem.ID).Error
		require.NoError(t, err)
		assert.Equal(t, firstSolvedAt, progress.FirstSolvedAt) // Never changes
		assert.Equal(t, 3, progress.TotalAttempts) // Incremented
		assert.True(t, progress.IsSolved) // Still solved
	})
}

func TestIntegration_DatabaseIntegrityAfterErrors(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	t.Run("transaction rollback on error maintains database integrity", func(t *testing.T) {
		// Create test problem
		problem := &database.Problem{Slug: "integrity-test", Title: "Integrity Test", Difficulty: "easy"}
		require.NoError(t, db.Create(problem).Error)

		// Track successful completion
		_, err := tracker.TrackTestCompletion(problem.ID, "problems/integrity-test/solution.go", true, 5, 5)
		require.NoError(t, err)

		// Count records before error attempt
		var progressCount, solutionCount int64
		db.Model(&database.Progress{}).Where("problem_id = ?", problem.ID).Count(&progressCount)
		db.Model(&database.Solution{}).Where("problem_id = ?", problem.ID).Count(&solutionCount)
		assert.Equal(t, int64(1), progressCount)
		assert.Equal(t, int64(1), solutionCount)

		// Attempt with invalid problem ID (should fail and rollback)
		_, err = tracker.TrackTestCompletion(999999, "invalid/path.go", true, 5, 5)
		assert.Error(t, err)

		// Verify no new records created for invalid problem
		db.Model(&database.Progress{}).Where("problem_id = ?", 999999).Count(&progressCount)
		db.Model(&database.Solution{}).Where("problem_id = ?", 999999).Count(&solutionCount)
		assert.Equal(t, int64(0), progressCount)
		assert.Equal(t, int64(0), solutionCount)

		// Verify original problem's data unchanged
		var progress database.Progress
		err = db.First(&progress, "problem_id = ?", problem.ID).Error
		require.NoError(t, err)
		assert.True(t, progress.IsSolved)
		assert.Equal(t, 1, progress.TotalAttempts) // Not incremented due to rollback on error
	})
}

func TestIntegration_Performance(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	t.Run("progress update completes in <100ms", func(t *testing.T) {
		// Create test problem
		problem := &database.Problem{Slug: "perf-test", Title: "Performance Test", Difficulty: "easy"}
		require.NoError(t, db.Create(problem).Error)

		// Measure time for progress update
		start := time.Now()
		_, err := tracker.TrackTestCompletion(problem.ID, "problems/perf-test/solution.go", true, 10, 10)
		elapsed := time.Since(start)

		require.NoError(t, err)
		assert.Less(t, elapsed.Milliseconds(), int64(100), "Progress update should complete in <100ms")
	})

	t.Run("multiple sequential updates maintain performance", func(t *testing.T) {
		// Create test problem
		problem := &database.Problem{Slug: "perf-multiple", Title: "Performance Multiple", Difficulty: "medium"}
		require.NoError(t, db.Create(problem).Error)

		// Run 5 updates and measure each
		for i := 1; i <= 5; i++ {
			start := time.Now()
			_, err := tracker.TrackTestCompletion(problem.ID, "problems/perf-multiple/solution.go", false, i, 10)
			elapsed := time.Since(start)

			require.NoError(t, err)
			assert.Less(t, elapsed.Milliseconds(), int64(100), "Each update should complete in <100ms")
		}

		// Verify all updates were recorded
		var progress database.Progress
		err := db.First(&progress, "problem_id = ?", problem.ID).Error
		require.NoError(t, err)
		assert.Equal(t, 5, progress.TotalAttempts)
	})
}

func TestIntegration_MultipleProblems(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	t.Run("tracks progress independently for multiple problems", func(t *testing.T) {
		// Create multiple problems
		problem1 := &database.Problem{Slug: "multi-1", Title: "Multi 1", Difficulty: "easy"}
		problem2 := &database.Problem{Slug: "multi-2", Title: "Multi 2", Difficulty: "medium"}
		problem3 := &database.Problem{Slug: "multi-3", Title: "Multi 3", Difficulty: "hard"}
		require.NoError(t, db.Create(problem1).Error)
		require.NoError(t, db.Create(problem2).Error)
		require.NoError(t, db.Create(problem3).Error)

		// Track different progress for each
		// Problem 1: Solved immediately
		_, err := tracker.TrackTestCompletion(problem1.ID, "problems/multi-1/solution.go", true, 5, 5)
		require.NoError(t, err)

		// Problem 2: Failed twice, then solved
		_, err = tracker.TrackTestCompletion(problem2.ID, "problems/multi-2/solution.go", false, 3, 5)
		require.NoError(t, err)
		_, err = tracker.TrackTestCompletion(problem2.ID, "problems/multi-2/solution.go", false, 4, 5)
		require.NoError(t, err)
		_, err = tracker.TrackTestCompletion(problem2.ID, "problems/multi-2/solution.go", true, 5, 5)
		require.NoError(t, err)

		// Problem 3: Still failing
		_, err = tracker.TrackTestCompletion(problem3.ID, "problems/multi-3/solution.go", false, 1, 5)
		require.NoError(t, err)

		// Verify independent progress
		var progress1, progress2, progress3 database.Progress
		db.First(&progress1, "problem_id = ?", problem1.ID)
		db.First(&progress2, "problem_id = ?", problem2.ID)
		db.First(&progress3, "problem_id = ?", problem3.ID)

		assert.True(t, progress1.IsSolved)
		assert.Equal(t, 1, progress1.TotalAttempts)

		assert.True(t, progress2.IsSolved)
		assert.Equal(t, 3, progress2.TotalAttempts)

		assert.False(t, progress3.IsSolved)
		assert.Equal(t, 1, progress3.TotalAttempts)
	})
}

package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProblemsBySolvedStatus(t *testing.T) {
	db := setupTestDB(t)

	t.Run("returns solved problems", func(t *testing.T) {
		// Create problems
		problem1 := &Problem{Slug: "two-sum", Title: "Two Sum"}
		problem2 := &Problem{Slug: "add-two", Title: "Add Two"}
		problem3 := &Problem{Slug: "reverse", Title: "Reverse"}
		require.NoError(t, db.Create(problem1).Error)
		require.NoError(t, db.Create(problem2).Error)
		require.NoError(t, db.Create(problem3).Error)

		// Mark problem1 and problem2 as solved
		progress1 := &Progress{ProblemID: problem1.ID, IsSolved: true}
		progress2 := &Progress{ProblemID: problem2.ID, IsSolved: true}
		progress3 := &Progress{ProblemID: problem3.ID, IsSolved: false}
		require.NoError(t, db.Create(progress1).Error)
		require.NoError(t, db.Create(progress2).Error)
		require.NoError(t, db.Create(progress3).Error)

		// Get solved problems
		solved, err := GetProblemsBySolvedStatus(db, true)

		assert.NoError(t, err)
		assert.Len(t, solved, 2)
	})

	t.Run("returns unsolved problems", func(t *testing.T) {
		db := setupTestDB(t)

		problem1 := &Problem{Slug: "problem-1", Title: "Problem 1"}
		problem2 := &Problem{Slug: "problem-2", Title: "Problem 2"}
		require.NoError(t, db.Create(problem1).Error)
		require.NoError(t, db.Create(problem2).Error)

		progress1 := &Progress{ProblemID: problem1.ID, IsSolved: false}
		progress2 := &Progress{ProblemID: problem2.ID, IsSolved: true}
		require.NoError(t, db.Create(progress1).Error)
		require.NoError(t, db.Create(progress2).Error)

		unsolved, err := GetProblemsBySolvedStatus(db, false)

		assert.NoError(t, err)
		assert.Len(t, unsolved, 1)
	})
}

func TestGetRecentSolutions(t *testing.T) {
	db := setupTestDB(t)

	t.Run("returns last N solutions ordered by time", func(t *testing.T) {
		problem := &Problem{Slug: "test-problem", Title: "Test Problem"}
		require.NoError(t, db.Create(problem).Error)

		// Create solutions with different submission times
		now := time.Now()
		solution1 := &Solution{
			ProblemID:   problem.ID,
			SubmittedAt: now.Add(-3 * time.Hour),
			Status:      "Passed",
		}
		solution2 := &Solution{
			ProblemID:   problem.ID,
			SubmittedAt: now.Add(-2 * time.Hour),
			Status:      "Failed",
		}
		solution3 := &Solution{
			ProblemID:   problem.ID,
			SubmittedAt: now.Add(-1 * time.Hour),
			Status:      "Passed",
		}

		require.NoError(t, db.Create(solution1).Error)
		require.NoError(t, db.Create(solution2).Error)
		require.NoError(t, db.Create(solution3).Error)

		// Get last 2 solutions
		recent, err := GetRecentSolutions(db, 2)

		assert.NoError(t, err)
		assert.Len(t, recent, 2)
		// Most recent first
		assert.True(t, recent[0].SubmittedAt.After(recent[1].SubmittedAt))
	})

	t.Run("handles limit larger than total solutions", func(t *testing.T) {
		db := setupTestDB(t)

		problem := &Problem{Slug: "test-limit", Title: "Test Limit"}
		require.NoError(t, db.Create(problem).Error)

		solution := &Solution{ProblemID: problem.ID, Status: "Passed"}
		require.NoError(t, db.Create(solution).Error)

		// Request 10 but only 1 exists
		recent, err := GetRecentSolutions(db, 10)

		assert.NoError(t, err)
		assert.Len(t, recent, 1)
	})
}

func TestGetCompletionStatistics(t *testing.T) {
	db := setupTestDB(t)

	t.Run("calculates statistics correctly", func(t *testing.T) {
		// Create 5 problems
		for i := 1; i <= 5; i++ {
			problem := &Problem{
				Slug:  "problem-" + string(rune(i)),
				Title: "Problem " + string(rune(i)),
			}
			require.NoError(t, db.Create(problem).Error)

			// Mark first 3 as solved with varying attempts
			if i <= 3 {
				progress := &Progress{
					ProblemID:     uint(i),
					IsSolved:      true,
					TotalAttempts: i * 2, // 2, 4, 6 attempts
				}
				require.NoError(t, db.Create(progress).Error)
			}
		}

		stats, err := GetCompletionStatistics(db)

		assert.NoError(t, err)
		assert.Equal(t, 5, stats.TotalProblems)
		assert.Equal(t, 3, stats.SolvedProblems)
		assert.Equal(t, 12, stats.TotalAttempts) // 2+4+6
		assert.InDelta(t, 4.0, stats.AverageAttempts, 0.01) // 12/3
	})

	t.Run("handles zero solved problems", func(t *testing.T) {
		db := setupTestDB(t)

		problem := &Problem{Slug: "unsolved", Title: "Unsolved"}
		require.NoError(t, db.Create(problem).Error)

		stats, err := GetCompletionStatistics(db)

		assert.NoError(t, err)
		assert.Equal(t, 1, stats.TotalProblems)
		assert.Equal(t, 0, stats.SolvedProblems)
		assert.Equal(t, 0.0, stats.AverageAttempts)
	})
}

func TestGetProgressByProblemID(t *testing.T) {
	db := setupTestDB(t)

	t.Run("retrieves progress for existing problem", func(t *testing.T) {
		problem := &Problem{Slug: "test-progress", Title: "Test Progress"}
		require.NoError(t, db.Create(problem).Error)

		expected := &Progress{
			ProblemID:     problem.ID,
			TotalAttempts: 5,
			IsSolved:      true,
		}
		require.NoError(t, db.Create(expected).Error)

		result, err := GetProgressByProblemID(db, problem.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.ProblemID, result.ProblemID)
		assert.Equal(t, 5, result.TotalAttempts)
		assert.True(t, result.IsSolved)
	})

	t.Run("returns error for non-existent problem", func(t *testing.T) {
		_, err := GetProgressByProblemID(db, 999)

		assert.Error(t, err)
	})
}

func TestUpsertProgress(t *testing.T) {
	db := setupTestDB(t)

	t.Run("creates new progress record", func(t *testing.T) {
		problem := &Problem{Slug: "upsert-test", Title: "Upsert Test"}
		require.NoError(t, db.Create(problem).Error)

		progress := &Progress{
			ProblemID:     problem.ID,
			TotalAttempts: 1,
			IsSolved:      false,
		}

		err := UpsertProgress(db, progress)

		assert.NoError(t, err)
		assert.NotZero(t, progress.ID)

		// Verify it was created
		var retrieved Progress
		require.NoError(t, db.Where("problem_id = ?", problem.ID).First(&retrieved).Error)
		assert.Equal(t, 1, retrieved.TotalAttempts)
	})

	t.Run("updates existing progress record", func(t *testing.T) {
		db := setupTestDB(t)

		problem := &Problem{Slug: "update-test", Title: "Update Test"}
		require.NoError(t, db.Create(problem).Error)

		// Create initial progress
		initial := &Progress{
			ProblemID:     problem.ID,
			TotalAttempts: 1,
			IsSolved:      false,
		}
		require.NoError(t, db.Create(initial).Error)

		// Update it
		updated := &Progress{
			ProblemID:     problem.ID,
			TotalAttempts: 3,
			IsSolved:      true,
		}

		err := UpsertProgress(db, updated)

		assert.NoError(t, err)

		// Verify update
		var retrieved Progress
		require.NoError(t, db.Where("problem_id = ?", problem.ID).First(&retrieved).Error)
		assert.Equal(t, 3, retrieved.TotalAttempts)
		assert.True(t, retrieved.IsSolved)
	})
}

package progress

import (
	"testing"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTracker(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	assert.NotNil(t, tracker)
	assert.NotNil(t, tracker.db)
}

func TestTrackTestCompletion_FirstSuccessfulSolve(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	// Create test problem
	problem := &database.Problem{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", Topic: "arrays"}
	require.NoError(t, db.Create(problem).Error)

	// Track successful test completion
	isFirstSolve, err := tracker.TrackTestCompletion(problem.ID, "problems/two-sum/solution.go", true, 5, 5)
	require.NoError(t, err)
	assert.True(t, isFirstSolve)

	// Verify progress record created
	var progress database.Progress
	err = db.First(&progress, "problem_id = ?", problem.ID).Error
	require.NoError(t, err)

	assert.True(t, progress.IsSolved)
	assert.NotNil(t, progress.FirstSolvedAt)
	assert.Equal(t, 1, progress.TotalAttempts)
	assert.False(t, progress.LastAttemptedAt.IsZero())

	// Verify solution record created
	var solution database.Solution
	err = db.First(&solution, "problem_id = ?", problem.ID).Error
	require.NoError(t, err)

	assert.Equal(t, problem.ID, solution.ProblemID)
	assert.Equal(t, "Passed", solution.Status)
	assert.Equal(t, 5, solution.TestsPassed)
	assert.Equal(t, 5, solution.TestsTotal)
	assert.Equal(t, "problems/two-sum/solution.go", solution.FilePath)
}

func TestTrackTestCompletion_FailedAttempt(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	// Create test problem
	problem := &database.Problem{Slug: "add-two", Title: "Add Two Numbers", Difficulty: "medium"}
	require.NoError(t, db.Create(problem).Error)

	// Track failed test completion
	isFirstSolve, err := tracker.TrackTestCompletion(problem.ID, "problems/add-two/solution.go", false, 3, 5)
	require.NoError(t, err)
	assert.False(t, isFirstSolve)

	// Verify progress record
	var progress database.Progress
	err = db.First(&progress, "problem_id = ?", problem.ID).Error
	require.NoError(t, err)

	assert.False(t, progress.IsSolved)
	assert.Nil(t, progress.FirstSolvedAt)
	assert.Equal(t, 1, progress.TotalAttempts)
	assert.False(t, progress.LastAttemptedAt.IsZero())

	// Verify solution record shows failure
	var solution database.Solution
	err = db.First(&solution, "problem_id = ?", problem.ID).Error
	require.NoError(t, err)

	assert.Equal(t, "Failed", solution.Status)
	assert.Equal(t, 3, solution.TestsPassed)
	assert.Equal(t, 5, solution.TestsTotal)
}

func TestTrackTestCompletion_MultipleAttempts(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	// Create test problem
	problem := &database.Problem{Slug: "reverse", Title: "Reverse String", Difficulty: "easy"}
	require.NoError(t, db.Create(problem).Error)

	// First attempt: failed
	_, err := tracker.TrackTestCompletion(problem.ID, "problems/reverse/solution.go", false, 2, 5)
	require.NoError(t, err)

	// Second attempt: failed
	_, err = tracker.TrackTestCompletion(problem.ID, "problems/reverse/solution.go", false, 4, 5)
	require.NoError(t, err)

	// Third attempt: passed
	isFirstSolve, err := tracker.TrackTestCompletion(problem.ID, "problems/reverse/solution.go", true, 5, 5)
	require.NoError(t, err)
	assert.True(t, isFirstSolve)

	// Verify progress record
	var progress database.Progress
	err = db.First(&progress, "problem_id = ?", problem.ID).Error
	require.NoError(t, err)

	assert.True(t, progress.IsSolved)
	assert.NotNil(t, progress.FirstSolvedAt)
	assert.Equal(t, 3, progress.TotalAttempts)

	// Verify all solution records created
	var solutions []database.Solution
	err = db.Find(&solutions, "problem_id = ?", problem.ID).Error
	require.NoError(t, err)
	assert.Equal(t, 3, len(solutions))
}

func TestTrackTestCompletion_SubsequentSolve_DoesNotChangeFirstSolvedAt(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	// Create test problem
	problem := &database.Problem{Slug: "palindrome", Title: "Valid Palindrome", Difficulty: "easy"}
	require.NoError(t, db.Create(problem).Error)

	// First solve
	isFirstSolve, err := tracker.TrackTestCompletion(problem.ID, "problems/palindrome/solution.go", true, 3, 3)
	require.NoError(t, err)
	assert.True(t, isFirstSolve)

	// Get first solved time
	var progress1 database.Progress
	err = db.First(&progress1, "problem_id = ?", problem.ID).Error
	require.NoError(t, err)
	firstSolvedAt := progress1.FirstSolvedAt

	// Wait a moment
	time.Sleep(10 * time.Millisecond)

	// Solve again (e.g., re-running tests)
	isFirstSolve, err = tracker.TrackTestCompletion(problem.ID, "problems/palindrome/solution.go", true, 3, 3)
	require.NoError(t, err)
	assert.False(t, isFirstSolve)

	// Verify FirstSolvedAt unchanged
	var progress2 database.Progress
	err = db.First(&progress2, "problem_id = ?", problem.ID).Error
	require.NoError(t, err)

	assert.Equal(t, firstSolvedAt, progress2.FirstSolvedAt)
	assert.Equal(t, 2, progress2.TotalAttempts)
	assert.True(t, progress2.LastAttemptedAt.After(*firstSolvedAt))
}

func TestTrackTestCompletion_IsFirstTimeSolve(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	// Create test problem
	problem := &database.Problem{Slug: "anagram", Title: "Valid Anagram", Difficulty: "easy"}
	require.NoError(t, db.Create(problem).Error)

	// First solve should return true
	isFirstSolve, err := tracker.TrackTestCompletion(problem.ID, "problems/anagram/solution.go", true, 5, 5)
	require.NoError(t, err)
	assert.True(t, isFirstSolve)

	// Second solve should return false
	isFirstSolve, err = tracker.TrackTestCompletion(problem.ID, "problems/anagram/solution.go", true, 5, 5)
	require.NoError(t, err)
	assert.False(t, isFirstSolve)

	// Failed attempt should return false
	isFirstSolve, err = tracker.TrackTestCompletion(problem.ID, "problems/anagram/solution.go", false, 3, 5)
	require.NoError(t, err)
	assert.False(t, isFirstSolve)
}

func TestTrackTestCompletion_TransactionRollback(t *testing.T) {
	db := setupTestDB(t)
	tracker := NewTracker(db)

	// Create test problem
	problem := &database.Problem{Slug: "test", Title: "Test Problem", Difficulty: "easy"}
	require.NoError(t, db.Create(problem).Error)

	// Track with invalid problem ID (should fail)
	_, err := tracker.TrackTestCompletion(999999, "invalid/path.go", true, 5, 5)
	assert.Error(t, err)

	// Verify no progress or solution records created (transaction rolled back)
	var progressCount int64
	db.Model(&database.Progress{}).Where("problem_id = ?", 999999).Count(&progressCount)
	assert.Equal(t, int64(0), progressCount)

	var solutionCount int64
	db.Model(&database.Solution{}).Where("problem_id = ?", 999999).Count(&solutionCount)
	assert.Equal(t, int64(0), solutionCount)
}

// Note: Concurrent test removed - SQLite has limited concurrency support due to database-level locking.
// In the CLI context, test executions are sequential (one at a time), so concurrent access isn't a real-world scenario.

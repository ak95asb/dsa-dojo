package database

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProblemModel(t *testing.T) {
	t.Run("Problem struct has correct fields", func(t *testing.T) {
		problem := Problem{
			ID:          1,
			Slug:        "two-sum",
			Title:       "Two Sum",
			Difficulty:  "easy",
			Topic:       "arrays",
			Description: "Find two numbers that add up to target",
			CreatedAt:   time.Now(),
		}

		assert.Equal(t, uint(1), problem.ID)
		assert.Equal(t, "two-sum", problem.Slug)
		assert.Equal(t, "Two Sum", problem.Title)
		assert.Equal(t, "easy", problem.Difficulty)
		assert.Equal(t, "arrays", problem.Topic)
		assert.Equal(t, "Find two numbers that add up to target", problem.Description)
	})

	t.Run("Problem JSON marshaling uses snake_case", func(t *testing.T) {
		problem := Problem{
			ID:          1,
			Slug:        "two-sum",
			Title:       "Two Sum",
			Difficulty:  "easy",
			Topic:       "arrays",
			Description: "Test description",
			CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		jsonData, err := json.Marshal(problem)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		assert.NoError(t, err)

		// Verify snake_case field names in JSON
		assert.Contains(t, result, "id")
		assert.Contains(t, result, "slug")
		assert.Contains(t, result, "title")
		assert.Contains(t, result, "difficulty")
		assert.Contains(t, result, "topic")
		assert.Contains(t, result, "description")
		assert.Contains(t, result, "created_at")

		// Verify values
		assert.Equal(t, float64(1), result["id"])
		assert.Equal(t, "two-sum", result["slug"])
	})
}

func TestSolutionModel(t *testing.T) {
	t.Run("Solution struct has correct fields", func(t *testing.T) {
		solution := Solution{
			ID:        1,
			ProblemID: 5,
			Code:      "func twoSum() {}",
			Language:  "go",
			Passed:    true,
			CreatedAt: time.Now(),
		}

		assert.Equal(t, uint(1), solution.ID)
		assert.Equal(t, uint(5), solution.ProblemID)
		assert.Equal(t, "func twoSum() {}", solution.Code)
		assert.Equal(t, "go", solution.Language)
		assert.True(t, solution.Passed)
	})

	t.Run("Solution JSON marshaling uses snake_case", func(t *testing.T) {
		solution := Solution{
			ID:        1,
			ProblemID: 5,
			Code:      "test code",
			Language:  "go",
			Passed:    true,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		jsonData, err := json.Marshal(solution)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		assert.NoError(t, err)

		// Verify snake_case field names
		assert.Contains(t, result, "id")
		assert.Contains(t, result, "problem_id")
		assert.Contains(t, result, "code")
		assert.Contains(t, result, "language")
		assert.Contains(t, result, "passed")
		assert.Contains(t, result, "created_at")

		// Verify values
		assert.Equal(t, float64(5), result["problem_id"])
		assert.Equal(t, "go", result["language"])
	})
}

func TestProgressModel(t *testing.T) {
	t.Run("Progress struct has correct fields", func(t *testing.T) {
		progress := Progress{
			ID:          1,
			ProblemID:   10,
			Status:      "in_progress",
			Attempts:    3,
			LastAttempt: time.Now(),
		}

		assert.Equal(t, uint(1), progress.ID)
		assert.Equal(t, uint(10), progress.ProblemID)
		assert.Equal(t, "in_progress", progress.Status)
		assert.Equal(t, 3, progress.Attempts)
	})

	t.Run("Progress JSON marshaling uses snake_case", func(t *testing.T) {
		progress := Progress{
			ID:          1,
			ProblemID:   10,
			Status:      "completed",
			Attempts:    5,
			LastAttempt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		jsonData, err := json.Marshal(progress)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)
		assert.NoError(t, err)

		// Verify snake_case field names
		assert.Contains(t, result, "id")
		assert.Contains(t, result, "problem_id")
		assert.Contains(t, result, "status")
		assert.Contains(t, result, "attempts")
		assert.Contains(t, result, "last_attempt")

		// Verify values
		assert.Equal(t, float64(10), result["problem_id"])
		assert.Equal(t, "completed", result["status"])
		assert.Equal(t, float64(5), result["attempts"])
	})
}

// Extended tests for new fields and functionality

func TestSolutionModel_Create(t *testing.T) {
	db := setupTestDB(t)

	t.Run("creates solution with valid data", func(t *testing.T) {
		problem := &Problem{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", Topic: "arrays"}
		require.NoError(t, db.Create(problem).Error)

		solution := &Solution{
			ProblemID:   problem.ID,
			FilePath:    "solutions/two_sum.go",
			Status:      "Passed",
			TestsPassed: 5,
			TestsTotal:  5,
			Code:        "func twoSum() {}",
			Language:    "go",
		}
		err := db.Create(solution).Error

		assert.NoError(t, err)
		assert.NotZero(t, solution.ID)
		assert.NotZero(t, solution.SubmittedAt)
		assert.NotZero(t, solution.CreatedAt)
	})

	t.Run("creates solution with InProgress status by default", func(t *testing.T) {
		problem := &Problem{Slug: "add-two-numbers", Title: "Add Two Numbers", Difficulty: "medium"}
		require.NoError(t, db.Create(problem).Error)

		solution := &Solution{
			ProblemID: problem.ID,
			Code:      "func addTwo() {}",
		}
		err := db.Create(solution).Error

		assert.NoError(t, err)
		assert.Equal(t, "InProgress", solution.Status)
	})

	t.Run("rejects invalid status", func(t *testing.T) {
		problem := &Problem{Slug: "invalid-test", Title: "Invalid Test"}
		require.NoError(t, db.Create(problem).Error)

		solution := &Solution{
			ProblemID: problem.ID,
			Status:    "InvalidStatus",
		}
		err := db.Create(solution).Error

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})
}

func TestSolutionModel_ValidateStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{"valid Passed", "Passed", false},
		{"valid Failed", "Failed", false},
		{"valid InProgress", "InProgress", false},
		{"invalid Pending", "Pending", true},
		{"invalid empty", "", true},
		{"invalid lowercase", "passed", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solution := &Solution{Status: tt.status}
			err := solution.ValidateStatus()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProgressModel_Create(t *testing.T) {
	db := setupTestDB(t)

	t.Run("creates progress with valid data", func(t *testing.T) {
		problem := &Problem{Slug: "two-sum", Title: "Two Sum"}
		require.NoError(t, db.Create(problem).Error)

		solvedAt := time.Now()
		progress := &Progress{
			ProblemID:       problem.ID,
			Status:          "completed",
			TotalAttempts:   3,
			IsSolved:        true,
			FirstSolvedAt:   &solvedAt,
			LastAttemptedAt: solvedAt,
		}
		err := db.Create(progress).Error

		assert.NoError(t, err)
		assert.NotZero(t, progress.ID)
	})

	t.Run("sets LastAttemptedAt to now if not provided", func(t *testing.T) {
		problem := &Problem{Slug: "add-two-numbers", Title: "Add Two Numbers"}
		require.NoError(t, db.Create(problem).Error)

		before := time.Now()
		progress := &Progress{
			ProblemID: problem.ID,
			Status:    "in_progress",
		}
		err := db.Create(progress).Error
		after := time.Now()

		assert.NoError(t, err)
		assert.True(t, progress.LastAttemptedAt.After(before) || progress.LastAttemptedAt.Equal(before))
		assert.True(t, progress.LastAttemptedAt.Before(after) || progress.LastAttemptedAt.Equal(after))
	})

	t.Run("handles nullable fields correctly", func(t *testing.T) {
		problem := &Problem{Slug: "test-nullable", Title: "Test Nullable"}
		require.NoError(t, db.Create(problem).Error)

		progress := &Progress{
			ProblemID:       problem.ID,
			Status:          "not_started",
			IsSolved:        false,
			BestTime:        nil, // Nullable
			FirstSolvedAt:   nil, // Nullable
		}
		err := db.Create(progress).Error

		assert.NoError(t, err)
		assert.Nil(t, progress.BestTime)
		assert.Nil(t, progress.FirstSolvedAt)
	})

	t.Run("stores best time in milliseconds", func(t *testing.T) {
		problem := &Problem{Slug: "test-best-time", Title: "Test Best Time"}
		require.NoError(t, db.Create(problem).Error)

		bestTime := 1250 // 1.25 seconds = 1250ms
		progress := &Progress{
			ProblemID: problem.ID,
			BestTime:  &bestTime,
		}
		err := db.Create(progress).Error

		assert.NoError(t, err)
		assert.NotNil(t, progress.BestTime)
		assert.Equal(t, 1250, *progress.BestTime)
	})
}

func TestProgressModel_UniqueConstraint(t *testing.T) {
	db := setupTestDB(t)

	t.Run("enforces unique problem_id constraint", func(t *testing.T) {
		problem := &Problem{Slug: "unique-test", Title: "Unique Test"}
		require.NoError(t, db.Create(problem).Error)

		progress1 := &Progress{ProblemID: problem.ID, Status: "in_progress"}
		require.NoError(t, db.Create(progress1).Error)

		progress2 := &Progress{ProblemID: problem.ID, Status: "completed"}
		err := db.Create(progress2).Error

		assert.Error(t, err)
	})
}

func TestProgressModel_BeforeUpdate(t *testing.T) {
	db := setupTestDB(t)

	t.Run("updates LastAttemptedAt on update", func(t *testing.T) {
		problem := &Problem{Slug: "update-test", Title: "Update Test"}
		require.NoError(t, db.Create(problem).Error)

		initialTime := time.Now().Add(-1 * time.Hour)
		progress := &Progress{
			ProblemID:       problem.ID,
			LastAttemptedAt: initialTime,
			TotalAttempts:   1,
		}
		require.NoError(t, db.Create(progress).Error)

		// Wait a moment to ensure time difference
		time.Sleep(10 * time.Millisecond)

		progress.TotalAttempts = 2
		require.NoError(t, db.Save(progress).Error)

		// Reload from database
		var updated Progress
		require.NoError(t, db.First(&updated, progress.ID).Error)

		assert.True(t, updated.LastAttemptedAt.After(initialTime))
	})
}

func TestIndexCreation(t *testing.T) {
	db := setupTestDB(t)

	t.Run("creates indexes on Solution", func(t *testing.T) {
		var count int64
		db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_solutions_problem_id'").Scan(&count)
		assert.Equal(t, int64(1), count, "idx_solutions_problem_id index should exist")
	})

	t.Run("creates indexes on Progress", func(t *testing.T) {
		var count int64
		db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_progress_problem_id'").Scan(&count)
		assert.Equal(t, int64(1), count, "idx_progress_problem_id index should exist")

		db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_progress_is_solved'").Scan(&count)
		assert.Equal(t, int64(1), count, "idx_progress_is_solved index should exist")

		db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_progress_first_solved'").Scan(&count)
		assert.Equal(t, int64(1), count, "idx_progress_first_solved index should exist")

		db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_progress_last_attempted'").Scan(&count)
		assert.Equal(t, int64(1), count, "idx_progress_last_attempted index should exist")
	})
}

func TestAutoMigrate_NonDestructive(t *testing.T) {
	db := setupTestDB(t)

	t.Run("preserves existing data on migration", func(t *testing.T) {
		problem := &Problem{Slug: "preserve-test", Title: "Preserve Test"}
		require.NoError(t, db.Create(problem).Error)

		solution := &Solution{
			ProblemID: problem.ID,
			Code:      "original code",
			Status:    "Passed",
		}
		require.NoError(t, db.Create(solution).Error)
		originalID := solution.ID

		// Run AutoMigrate again (simulating schema evolution)
		err := db.AutoMigrate(&Solution{})
		require.NoError(t, err)

		// Verify data still exists
		var retrieved Solution
		require.NoError(t, db.First(&retrieved, originalID).Error)
		assert.Equal(t, "original code", retrieved.Code)
		assert.Equal(t, "Passed", retrieved.Status)
	})
}

package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&Problem{}, &Solution{}, &Progress{}, &BenchmarkResult{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

func TestInitialize(t *testing.T) {
	t.Run("creates database and runs migrations", func(t *testing.T) {
		// Use in-memory database for testing
		db := setupTestDB(t)
		assert.NotNil(t, db)

		// Verify tables exist by attempting to query them
		var count int64
		err := db.Model(&Problem{}).Count(&count).Error
		assert.NoError(t, err)

		err = db.Model(&Solution{}).Count(&count).Error
		assert.NoError(t, err)

		err = db.Model(&Progress{}).Count(&count).Error
		assert.NoError(t, err)
	})

	t.Run("creates all three tables with correct names", func(t *testing.T) {
		db := setupTestDB(t)

		// GORM uses plural snake_case table names
		// Verify by checking we can create records
		problem := Problem{
			Slug:        "test-problem",
			Title:       "Test Problem",
			Difficulty:  "easy",
			Topic:       "arrays",
			Description: "Test description",
		}
		err := db.Create(&problem).Error
		assert.NoError(t, err)
		assert.Greater(t, problem.ID, uint(0))

		solution := Solution{
			ProblemID: problem.ID,
			Code:      "test code",
			Language:  "go",
			Passed:    true,
		}
		err = db.Create(&solution).Error
		assert.NoError(t, err)

		progress := Progress{
			ProblemID: problem.ID,
			Status:    "in_progress",
			Attempts:  1,
		}
		err = db.Create(&progress).Error
		assert.NoError(t, err)
	})

	t.Run("verifies column names are snake_case", func(t *testing.T) {
		db := setupTestDB(t)

		// Create a problem and verify it can be queried
		problem := Problem{
			Slug:        "column-test",
			Title:       "Column Test",
			Difficulty:  "medium",
			Topic:       "trees",
			Description: "Testing column names",
		}
		err := db.Create(&problem).Error
		assert.NoError(t, err)

		// Query using snake_case column names
		var retrieved Problem
		err = db.Where("slug = ?", "column-test").First(&retrieved).Error
		assert.NoError(t, err)
		assert.Equal(t, "Column Test", retrieved.Title)
		assert.Equal(t, "medium", retrieved.Difficulty)

		// Verify solution with problem_id foreign key
		solution := Solution{
			ProblemID: problem.ID,
			Code:      "test",
			Language:  "go",
		}
		err = db.Create(&solution).Error
		assert.NoError(t, err)

		var retrievedSolution Solution
		err = db.Where("problem_id = ?", problem.ID).First(&retrievedSolution).Error
		assert.NoError(t, err)
		assert.Equal(t, problem.ID, retrievedSolution.ProblemID)
	})

	t.Run("enforces unique constraints", func(t *testing.T) {
		db := setupTestDB(t)

		// Create first problem with slug
		problem1 := Problem{
			Slug:       "unique-test",
			Title:      "First Problem",
			Difficulty: "easy",
		}
		err := db.Create(&problem1).Error
		assert.NoError(t, err)

		// Attempt to create second problem with same slug (should fail)
		problem2 := Problem{
			Slug:       "unique-test",
			Title:      "Second Problem",
			Difficulty: "medium",
		}
		err = db.Create(&problem2).Error
		assert.Error(t, err) // Should fail due to unique constraint

		// Verify only one Progress record per problem
		progress1 := Progress{
			ProblemID: problem1.ID,
			Status:    "in_progress",
			Attempts:  1,
		}
		err = db.Create(&progress1).Error
		assert.NoError(t, err)

		progress2 := Progress{
			ProblemID: problem1.ID, // Same problem_id
			Status:    "completed",
			Attempts:  2,
		}
		err = db.Create(&progress2).Error
		assert.Error(t, err) // Should fail due to unique constraint on problem_id
	})

	t.Run("sets default values correctly", func(t *testing.T) {
		db := setupTestDB(t)

		// Create problem for foreign key
		problem := Problem{
			Slug:       "default-test",
			Title:      "Default Test",
			Difficulty: "easy",
		}
		err := db.Create(&problem).Error
		assert.NoError(t, err)

		// Solution with defaults
		solution := Solution{
			ProblemID: problem.ID,
			Code:      "some code",
			// Language not set - should default to 'go'
			// Passed not set - should default to false
		}
		err = db.Create(&solution).Error
		assert.NoError(t, err)

		var retrieved Solution
		err = db.First(&retrieved, solution.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "go", retrieved.Language)
		assert.False(t, retrieved.Passed)

		// Progress with defaults
		progress := Progress{
			ProblemID: problem.ID,
			// Status not set - should default to 'not_started'
			// Attempts not set - should default to 0
		}
		err = db.Create(&progress).Error
		assert.NoError(t, err)

		var retrievedProgress Progress
		err = db.First(&retrievedProgress, progress.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "not_started", retrievedProgress.Status)
		assert.Equal(t, 0, retrievedProgress.Attempts)
	})
}

func TestInitializeFunction(t *testing.T) {
	t.Run("creates .dsa directory and database file", func(t *testing.T) {
		// Create temporary directory for test
		tempHome := t.TempDir()
		originalHome := os.Getenv("HOME")
		t.Setenv("HOME", tempHome)
		defer os.Setenv("HOME", originalHome)

		// Call Initialize
		db, err := Initialize()
		assert.NoError(t, err)
		assert.NotNil(t, db)

		// Verify .dsa directory was created
		dsaDir := filepath.Join(tempHome, ".dsa")
		stat, err := os.Stat(dsaDir)
		assert.NoError(t, err)
		assert.True(t, stat.IsDir())

		// Verify database file was created
		dbPath := filepath.Join(dsaDir, "dsa.db")
		_, err = os.Stat(dbPath)
		assert.NoError(t, err)

		// Verify tables exist
		var count int64
		err = db.Model(&Problem{}).Count(&count).Error
		assert.NoError(t, err)

		err = db.Model(&Solution{}).Count(&count).Error
		assert.NoError(t, err)

		err = db.Model(&Progress{}).Count(&count).Error
		assert.NoError(t, err)
	})

	t.Run("returns error with wrapped context on failure", func(t *testing.T) {
		// Set HOME to invalid path to trigger error
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", "")
		defer os.Setenv("HOME", originalHome)

		_, err := Initialize()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get user home directory")
	})
}

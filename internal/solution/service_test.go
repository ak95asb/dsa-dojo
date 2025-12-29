package solution

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate models
	err = db.AutoMigrate(&database.Problem{}, &database.Solution{}, &database.Progress{})
	require.NoError(t, err)

	return db
}

// createTestSolutionFile creates a temporary solution file
func createTestSolutionFile(t *testing.T, dir, filename, content string) string {
	path := filepath.Join(dir, filename)
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
	return path
}

func TestRecordSubmission(t *testing.T) {
	t.Run("records passing solution successfully", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		// Create temp directory for test
		tempDir := t.TempDir()
		solutionPath := createTestSolutionFile(t, tempDir, "two-sum.go", "package solutions\n\nfunc TwoSum() {}")

		// Change to temp directory
		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)

		// Record submission
		record, err := svc.RecordSubmission("two-sum", 1, solutionPath, true, 5, 5)

		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, uint(1), record.ProblemID)
		assert.True(t, record.Passed)
		assert.Equal(t, 5, record.TestsPassed)
		assert.Equal(t, 5, record.TestsTotal)

		// Verify history file was created
		historyDir := filepath.Join("solutions", "history", "two-sum")
		entries, err := os.ReadDir(historyDir)
		assert.NoError(t, err)
		assert.Len(t, entries, 1)
	})

	t.Run("records failing solution successfully", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		tempDir := t.TempDir()
		solutionPath := createTestSolutionFile(t, tempDir, "binary-search.go", "package solutions\n\nfunc BinarySearch() {}")

		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)

		record, err := svc.RecordSubmission("binary-search", 2, solutionPath, false, 2, 5)

		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.False(t, record.Passed)
		assert.Equal(t, 2, record.TestsPassed)
		assert.Equal(t, 5, record.TestsTotal)
	})

	t.Run("returns error when solution file does not exist", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		tempDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)

		nonExistentPath := filepath.Join(tempDir, "nonexistent.go")
		record, err := svc.RecordSubmission("test", 1, nonExistentPath, true, 5, 5)

		assert.Error(t, err)
		assert.Nil(t, record)
	})
}

func TestGetHistory(t *testing.T) {
	t.Run("returns all submissions sorted by most recent", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		// Create 3 submissions with different timestamps
		solutions := []database.Solution{
			{ProblemID: 1, Code: "code1", Language: "go", Passed: true, CreatedAt: time.Now().Add(-2 * time.Hour)},
			{ProblemID: 1, Code: "code2", Language: "go", Passed: false, CreatedAt: time.Now().Add(-1 * time.Hour)},
			{ProblemID: 1, Code: "code3", Language: "go", Passed: true, CreatedAt: time.Now()},
		}

		for _, sol := range solutions {
			db.Create(&sol)
		}

		records, err := svc.GetHistory(1)

		assert.NoError(t, err)
		assert.Len(t, records, 3)
		// Verify sorted by most recent first
		assert.Equal(t, "code3", records[0].Code)
		assert.Equal(t, "code2", records[1].Code)
		assert.Equal(t, "code1", records[2].Code)
	})

	t.Run("returns empty list when no submissions exist", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		records, err := svc.GetHistory(999)

		assert.NoError(t, err)
		assert.Len(t, records, 0)
	})

	t.Run("only returns submissions for specified problem", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		// Create submissions for different problems
		db.Create(&database.Solution{ProblemID: 1, Code: "code1", Language: "go", Passed: true})
		db.Create(&database.Solution{ProblemID: 2, Code: "code2", Language: "go", Passed: true})
		db.Create(&database.Solution{ProblemID: 1, Code: "code3", Language: "go", Passed: false})

		records, err := svc.GetHistory(1)

		assert.NoError(t, err)
		assert.Len(t, records, 2)
		for _, record := range records {
			assert.Equal(t, uint(1), record.ProblemID)
		}
	})
}

func TestGetSubmissionByIndex(t *testing.T) {
	t.Run("retrieves submission by 1-based index", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		// Create 3 submissions
		solutions := []database.Solution{
			{ProblemID: 1, Code: "code1", Language: "go", Passed: true, CreatedAt: time.Now().Add(-2 * time.Hour)},
			{ProblemID: 1, Code: "code2", Language: "go", Passed: false, CreatedAt: time.Now().Add(-1 * time.Hour)},
			{ProblemID: 1, Code: "code3", Language: "go", Passed: true, CreatedAt: time.Now()},
		}

		for _, sol := range solutions {
			db.Create(&sol)
		}

		// Index 1 = most recent
		record, err := svc.GetSubmissionByIndex(1, 1)
		assert.NoError(t, err)
		assert.Equal(t, "code3", record.Code)

		// Index 2 = second most recent
		record, err = svc.GetSubmissionByIndex(1, 2)
		assert.NoError(t, err)
		assert.Equal(t, "code2", record.Code)

		// Index 3 = oldest
		record, err = svc.GetSubmissionByIndex(1, 3)
		assert.NoError(t, err)
		assert.Equal(t, "code1", record.Code)
	})

	t.Run("returns error for index less than 1", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		record, err := svc.GetSubmissionByIndex(1, 0)

		assert.Error(t, err)
		assert.Nil(t, record)
		assert.Contains(t, err.Error(), "index must be >= 1")
	})

	t.Run("returns error for index out of range", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		// Create only 2 submissions
		db.Create(&database.Solution{ProblemID: 1, Code: "code1", Language: "go", Passed: true})
		db.Create(&database.Solution{ProblemID: 1, Code: "code2", Language: "go", Passed: true})

		record, err := svc.GetSubmissionByIndex(1, 3)

		assert.Error(t, err)
		assert.Nil(t, record)
	})
}

func TestBackupCurrentSolution(t *testing.T) {
	t.Run("creates backup of current solution", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		tempDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)

		// Create current solution
		solutionDir := filepath.Join(tempDir, "solutions")
		os.MkdirAll(solutionDir, 0755)
		solutionPath := filepath.Join(solutionDir, "two-sum.go")
		solutionContent := "package solutions\n\nfunc TwoSum() { return 42 }"
		os.WriteFile(solutionPath, []byte(solutionContent), 0644)

		backupPath, err := svc.BackupCurrentSolution("two-sum", 1)

		assert.NoError(t, err)
		assert.NotEmpty(t, backupPath)

		// Verify backup file exists and has correct content
		content, err := os.ReadFile(backupPath)
		assert.NoError(t, err)
		assert.Equal(t, solutionContent, string(content))
	})

	t.Run("returns empty string when no current solution exists", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		tempDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)

		backupPath, err := svc.BackupCurrentSolution("nonexistent", 1)

		assert.NoError(t, err)
		assert.Empty(t, backupPath)
	})
}

func TestRestoreSolution(t *testing.T) {
	t.Run("restores solution successfully", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		tempDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)

		// Create submission record
		record := &SubmissionRecord{
			ID:        1,
			ProblemID: 1,
			Code:      "package solutions\n\nfunc TwoSum() { return 123 }",
			Language:  "go",
			Passed:    true,
			CreatedAt: time.Now(),
		}

		err := svc.RestoreSolution("two-sum", record)

		assert.NoError(t, err)

		// Verify solution file was created
		solutionPath := filepath.Join("solutions", "two-sum.go")
		content, err := os.ReadFile(solutionPath)
		assert.NoError(t, err)
		assert.Equal(t, record.Code, string(content))
	})

	t.Run("creates solutions directory if it doesn't exist", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		tempDir := t.TempDir()
		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)

		record := &SubmissionRecord{
			Code: "package solutions\n\nfunc Test() {}",
		}

		err := svc.RestoreSolution("test", record)

		assert.NoError(t, err)

		// Verify directory was created
		solutionDir := filepath.Join("solutions")
		_, err = os.Stat(solutionDir)
		assert.NoError(t, err)
	})
}

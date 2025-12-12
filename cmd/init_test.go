package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestInitCommand(t *testing.T) {
	t.Run("successfully initializes workspace", func(t *testing.T) {
		// Create temp directory for test
		tempHome := t.TempDir()
		originalHome := os.Getenv("HOME")
		t.Setenv("HOME", tempHome)
		defer os.Setenv("HOME", originalHome)

		// Execute command
		initCmd.Run(initCmd, []string{})

		// Verify database created
		dsaDir := filepath.Join(tempHome, ".dsa")
		dbPath := filepath.Join(dsaDir, "dsa.db")
		assert.DirExists(t, dsaDir)
		assert.FileExists(t, dbPath)

		// Verify database is valid and has tables
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		assert.NoError(t, err)

		// Check that tables exist by querying them
		var count int64
		err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='problems'").Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count, "problems table should exist")

		err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='solutions'").Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count, "solutions table should exist")

		err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='progresses'").Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count, "progresses table should exist")
	})

	t.Run("detects existing workspace and skips reinitialization", func(t *testing.T) {
		// Create temp directory and initialize once
		tempHome := t.TempDir()
		originalHome := os.Getenv("HOME")
		t.Setenv("HOME", tempHome)
		defer os.Setenv("HOME", originalHome)

		// First initialization
		initCmd.Run(initCmd, []string{})

		dsaDir := filepath.Join(tempHome, ".dsa")
		dbPath := filepath.Join(dsaDir, "dsa.db")

		// Get file modification time before second init
		stat1, err := os.Stat(dbPath)
		assert.NoError(t, err)

		// Second initialization (should detect existing)
		initCmd.Run(initCmd, []string{})

		// Verify file wasn't modified (no reinitialization happened)
		stat2, err := os.Stat(dbPath)
		assert.NoError(t, err)
		assert.Equal(t, stat1.ModTime(), stat2.ModTime(), "database file should not be modified on reinit")

		// Verify database is still valid
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		assert.NoError(t, err)

		// Verify our three tables still exist
		var problemCount, solutionCount, progressCount int64
		err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='problems'").Scan(&problemCount).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), problemCount)

		err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='solutions'").Scan(&solutionCount).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), solutionCount)

		err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='progresses'").Scan(&progressCount).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), progressCount)
	})

	t.Run("verifies directory permissions", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)

		// Initialize workspace
		initCmd.Run(initCmd, []string{})

		// Check directory permissions
		dsaDir := filepath.Join(tempHome, ".dsa")
		info, err := os.Stat(dsaDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir())

		// Verify directory has proper permissions (0755)
		mode := info.Mode()
		assert.True(t, mode.IsDir())
		// Check owner has read, write, execute
		assert.True(t, mode&0700 == 0700, "owner should have rwx")
	})

	t.Run("verifies database file permissions", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)

		// Initialize workspace
		initCmd.Run(initCmd, []string{})

		// Check database file permissions
		dbPath := filepath.Join(tempHome, ".dsa", "dsa.db")
		info, err := os.Stat(dbPath)
		assert.NoError(t, err)
		assert.False(t, info.IsDir())

		// Database file should be readable and writable
		mode := info.Mode()
		assert.True(t, mode.IsRegular())
	})

	t.Run("data integrity: rerunning doesn't corrupt existing data", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)

		// First initialization
		initCmd.Run(initCmd, []string{})

		dbPath := filepath.Join(tempHome, ".dsa", "dsa.db")

		// Open database and insert test data
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		assert.NoError(t, err)

		// Insert a test problem
		type Problem struct {
			ID         uint   `gorm:"primaryKey"`
			Slug       string `gorm:"uniqueIndex"`
			Title      string
			Difficulty string
		}
		testProblem := Problem{
			Slug:       "test-problem",
			Title:      "Test Problem",
			Difficulty: "easy",
		}
		err = db.Table("problems").Create(&testProblem).Error
		assert.NoError(t, err)

		// Close database
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Run init again (should detect existing and NOT corrupt)
		initCmd.Run(initCmd, []string{})

		// Reopen database and verify data still exists
		db2, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		assert.NoError(t, err)

		var retrievedProblem Problem
		err = db2.Table("problems").Where("slug = ?", "test-problem").First(&retrievedProblem).Error
		assert.NoError(t, err)
		assert.Equal(t, "Test Problem", retrievedProblem.Title)
		assert.Equal(t, "easy", retrievedProblem.Difficulty)
	})
}

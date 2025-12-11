package database

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Initialize creates and initializes the SQLite database with all models.
// It creates the ~/.dsa directory if it doesn't exist, opens the database
// at ~/.dsa/dsa.db, and runs AutoMigrate to create/update all tables.
//
// Returns the database connection or an error with wrapped context if any
// step fails.
func Initialize() (*gorm.DB, error) {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create .dsa directory if it doesn't exist
	dsaDir := filepath.Join(homeDir, ".dsa")
	if err := os.MkdirAll(dsaDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .dsa directory: %w", err)
	}

	// Database file path
	dbPath := filepath.Join(dsaDir, "dsa.db")

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database at %s: %w", dbPath, err)
	}

	// Run AutoMigrate for all models
	if err := db.AutoMigrate(&Problem{}, &Solution{}, &Progress{}, &BenchmarkResult{}); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return db, nil
}

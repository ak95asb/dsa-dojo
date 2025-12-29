package benchmarking

import (
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
	err = db.AutoMigrate(&database.Problem{}, &database.BenchmarkResult{})
	require.NoError(t, err)

	return db
}

func TestStorage_SaveBenchmark(t *testing.T) {
	t.Run("saves benchmark successfully", func(t *testing.T) {
		db := setupTestDB(t)
		storage := NewStorage(db)

		result := &BenchmarkResult{
			NsPerOp:     1234.5,
			BytesPerOp:  512,
			AllocsPerOp: 5,
		}

		err := storage.SaveBenchmark(1, result)

		assert.NoError(t, err)

		// Verify saved
		var saved database.BenchmarkResult
		db.First(&saved)
		assert.Equal(t, uint(1), saved.ProblemID)
		assert.Equal(t, 1234.5, saved.NsPerOp)
		assert.Equal(t, 512.0, saved.BytesPerOp)
		assert.Equal(t, 5.0, saved.AllocsPerOp)
	})
}

func TestStorage_GetBestBenchmark(t *testing.T) {
	t.Run("returns fastest benchmark", func(t *testing.T) {
		db := setupTestDB(t)
		storage := NewStorage(db)

		// Create multiple benchmarks with different speeds
		benchmarks := []database.BenchmarkResult{
			{ProblemID: 1, NsPerOp: 2000, BytesPerOp: 512, AllocsPerOp: 5, CreatedAt: time.Now().Add(-2 * time.Hour)},
			{ProblemID: 1, NsPerOp: 1000, BytesPerOp: 512, AllocsPerOp: 5, CreatedAt: time.Now().Add(-1 * time.Hour)}, // Fastest
			{ProblemID: 1, NsPerOp: 1500, BytesPerOp: 512, AllocsPerOp: 5, CreatedAt: time.Now()},
		}

		for _, b := range benchmarks {
			db.Create(&b)
		}

		best, err := storage.GetBestBenchmark(1)

		assert.NoError(t, err)
		assert.NotNil(t, best)
		assert.Equal(t, 1000.0, best.NsPerOp) // Fastest one
	})

	t.Run("returns nil when no benchmarks exist", func(t *testing.T) {
		db := setupTestDB(t)
		storage := NewStorage(db)

		best, err := storage.GetBestBenchmark(999)

		assert.NoError(t, err)
		assert.Nil(t, best)
	})

	t.Run("only returns benchmarks for specified problem", func(t *testing.T) {
		db := setupTestDB(t)
		storage := NewStorage(db)

		// Create benchmarks for different problems
		db.Create(&database.BenchmarkResult{ProblemID: 1, NsPerOp: 1000, BytesPerOp: 512, AllocsPerOp: 5})
		db.Create(&database.BenchmarkResult{ProblemID: 2, NsPerOp: 500, BytesPerOp: 256, AllocsPerOp: 3})

		best, err := storage.GetBestBenchmark(1)

		assert.NoError(t, err)
		assert.NotNil(t, best)
		assert.Equal(t, uint(1), best.ProblemID)
		assert.Equal(t, 1000.0, best.NsPerOp)
	})
}

func TestStorage_GetBenchmarkHistory(t *testing.T) {
	t.Run("returns all benchmarks sorted by most recent", func(t *testing.T) {
		db := setupTestDB(t)
		storage := NewStorage(db)

		// Create benchmarks with different timestamps
		now := time.Now()
		benchmarks := []database.BenchmarkResult{
			{ProblemID: 1, NsPerOp: 1000, BytesPerOp: 512, AllocsPerOp: 5, CreatedAt: now.Add(-2 * time.Hour)},
			{ProblemID: 1, NsPerOp: 1500, BytesPerOp: 600, AllocsPerOp: 6, CreatedAt: now.Add(-1 * time.Hour)},
			{ProblemID: 1, NsPerOp: 1200, BytesPerOp: 550, AllocsPerOp: 5, CreatedAt: now},
		}

		for _, b := range benchmarks {
			db.Create(&b)
		}

		history, err := storage.GetBenchmarkHistory(1)

		assert.NoError(t, err)
		assert.Len(t, history, 3)
		// Should be sorted by most recent first
		assert.Equal(t, 1200.0, history[0].NsPerOp) // Most recent
		assert.Equal(t, 1500.0, history[1].NsPerOp)
		assert.Equal(t, 1000.0, history[2].NsPerOp) // Oldest
	})

	t.Run("returns empty list when no benchmarks exist", func(t *testing.T) {
		db := setupTestDB(t)
		storage := NewStorage(db)

		history, err := storage.GetBenchmarkHistory(999)

		assert.NoError(t, err)
		assert.Len(t, history, 0)
	})
}

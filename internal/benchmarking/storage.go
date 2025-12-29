package benchmarking

import (
	"fmt"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"gorm.io/gorm"
)

// Storage handles persistence of benchmark results
type Storage struct {
	db *gorm.DB
}

// NewStorage creates a new storage instance
func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

// SaveBenchmark saves a benchmark result to the database
func (s *Storage) SaveBenchmark(problemID uint, result *BenchmarkResult) error {
	benchmark := &database.BenchmarkResult{
		ProblemID:   problemID,
		NsPerOp:     result.NsPerOp,
		AllocsPerOp: result.AllocsPerOp,
		BytesPerOp:  result.BytesPerOp,
	}

	if err := s.db.Create(benchmark).Error; err != nil {
		return fmt.Errorf("failed to save benchmark: %w", err)
	}

	return nil
}

// GetBestBenchmark retrieves the best (fastest) benchmark result for a problem
func (s *Storage) GetBestBenchmark(problemID uint) (*database.BenchmarkResult, error) {
	var benchmark database.BenchmarkResult

	// Find benchmark with lowest ns_per_op (fastest)
	err := s.db.Where("problem_id = ?", problemID).
		Order("ns_per_op ASC").
		First(&benchmark).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No previous benchmarks
		}
		return nil, fmt.Errorf("failed to query best benchmark: %w", err)
	}

	return &benchmark, nil
}

// GetBenchmarkHistory retrieves all benchmark results for a problem
func (s *Storage) GetBenchmarkHistory(problemID uint) ([]database.BenchmarkResult, error) {
	var benchmarks []database.BenchmarkResult

	err := s.db.Where("problem_id = ?", problemID).
		Order("created_at DESC").
		Find(&benchmarks).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query benchmark history: %w", err)
	}

	return benchmarks, nil
}

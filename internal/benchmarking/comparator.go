package benchmarking

import (
	"github.com/ak95asb/dsa-dojo/internal/database"
)

// Comparator handles comparing benchmark results
type Comparator struct{}

// NewComparator creates a new comparator
func NewComparator() *Comparator {
	return &Comparator{}
}

// ComparisonResult represents the comparison between current and previous benchmarks
type ComparisonResult struct {
	TimeDeltaPercent   float64
	MemoryDeltaPercent float64
	AllocsDeltaPercent float64
	IsNewBest          bool
}

// Compare compares current result with previous best
func (c *Comparator) Compare(current *BenchmarkResult, previous *database.BenchmarkResult) *ComparisonResult {
	if previous == nil {
		// No previous results to compare
		return &ComparisonResult{
			IsNewBest: true, // First benchmark is always "best"
		}
	}

	comparison := &ComparisonResult{}

	// Calculate time delta percentage
	// Positive = slower (regression), Negative = faster (improvement)
	comparison.TimeDeltaPercent = ((current.NsPerOp - previous.NsPerOp) / previous.NsPerOp) * 100

	// Calculate memory delta percentage
	comparison.MemoryDeltaPercent = ((current.BytesPerOp - previous.BytesPerOp) / previous.BytesPerOp) * 100

	// Calculate allocations delta percentage
	comparison.AllocsDeltaPercent = ((current.AllocsPerOp - previous.AllocsPerOp) / previous.AllocsPerOp) * 100

	// Determine if this is a new best (faster than previous best)
	comparison.IsNewBest = current.NsPerOp < previous.NsPerOp

	return comparison
}

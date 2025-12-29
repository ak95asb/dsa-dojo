package benchmarking

import (
	"testing"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestComparator_Compare(t *testing.T) {
	t.Run("first benchmark is always best", func(t *testing.T) {
		comp := NewComparator()
		current := &BenchmarkResult{
			NsPerOp:     1000,
			BytesPerOp:  500,
			AllocsPerOp: 10,
		}

		result := comp.Compare(current, nil)

		assert.True(t, result.IsNewBest)
		assert.Equal(t, 0.0, result.TimeDeltaPercent)
		assert.Equal(t, 0.0, result.MemoryDeltaPercent)
		assert.Equal(t, 0.0, result.AllocsDeltaPercent)
	})

	t.Run("faster time shows improvement", func(t *testing.T) {
		comp := NewComparator()
		current := &BenchmarkResult{
			NsPerOp:     800,
			BytesPerOp:  500,
			AllocsPerOp: 10,
		}
		previous := &database.BenchmarkResult{
			NsPerOp:     1000,
			BytesPerOp:  500,
			AllocsPerOp: 10,
			CreatedAt:   time.Now(),
		}

		result := comp.Compare(current, previous)

		assert.True(t, result.IsNewBest)
		assert.Equal(t, -20.0, result.TimeDeltaPercent) // 20% faster
	})

	t.Run("slower time shows regression", func(t *testing.T) {
		comp := NewComparator()
		current := &BenchmarkResult{
			NsPerOp:     1200,
			BytesPerOp:  500,
			AllocsPerOp: 10,
		}
		previous := &database.BenchmarkResult{
			NsPerOp:     1000,
			BytesPerOp:  500,
			AllocsPerOp: 10,
			CreatedAt:   time.Now(),
		}

		result := comp.Compare(current, previous)

		assert.False(t, result.IsNewBest)
		assert.Equal(t, 20.0, result.TimeDeltaPercent) // 20% slower
	})

	t.Run("less memory shows improvement", func(t *testing.T) {
		comp := NewComparator()
		current := &BenchmarkResult{
			NsPerOp:     1000,
			BytesPerOp:  400,
			AllocsPerOp: 10,
		}
		previous := &database.BenchmarkResult{
			NsPerOp:     1000,
			BytesPerOp:  500,
			AllocsPerOp: 10,
			CreatedAt:   time.Now(),
		}

		result := comp.Compare(current, previous)

		assert.Equal(t, -20.0, result.MemoryDeltaPercent) // 20% less memory
	})

	t.Run("fewer allocations shows improvement", func(t *testing.T) {
		comp := NewComparator()
		current := &BenchmarkResult{
			NsPerOp:     1000,
			BytesPerOp:  500,
			AllocsPerOp: 8,
		}
		previous := &database.BenchmarkResult{
			NsPerOp:     1000,
			BytesPerOp:  500,
			AllocsPerOp: 10,
			CreatedAt:   time.Now(),
		}

		result := comp.Compare(current, previous)

		assert.Equal(t, -20.0, result.AllocsDeltaPercent) // 20% fewer allocations
	})
}

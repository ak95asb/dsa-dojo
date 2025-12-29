package benchmarking

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatter_FormatResult(t *testing.T) {
	t.Run("formats basic result correctly", func(t *testing.T) {
		formatter := NewFormatter()
		result := &BenchmarkResult{
			Iterations:  1000000,
			NsPerOp:     1234.5,
			BytesPerOp:  512,
			AllocsPerOp: 5,
		}

		output := formatter.FormatResult(result)

		assert.Contains(t, output, "1,000,000") // Formatted iterations
		assert.Contains(t, output, "µs")        // Time unit
		assert.Contains(t, output, "512 B")     // Memory
		assert.Contains(t, output, "5")         // Allocations
	})
}

func TestFormatter_FormatComparison(t *testing.T) {
	t.Run("shows no comparison message when nil", func(t *testing.T) {
		formatter := NewFormatter()

		output := formatter.FormatComparison(nil)

		assert.Contains(t, output, "No previous benchmark results")
	})

	t.Run("shows improvement message for faster time", func(t *testing.T) {
		formatter := NewFormatter()
		comparison := &ComparisonResult{
			TimeDeltaPercent:   -15.3,
			MemoryDeltaPercent: 0,
			AllocsDeltaPercent: 0,
			IsNewBest:          true,
		}

		output := formatter.FormatComparison(comparison)

		assert.Contains(t, output, "🚀")
		assert.Contains(t, output, "15.3% faster")
		assert.Contains(t, output, "New personal best")
	})

	t.Run("shows regression message for slower time", func(t *testing.T) {
		formatter := NewFormatter()
		comparison := &ComparisonResult{
			TimeDeltaPercent:   20.0,
			MemoryDeltaPercent: 0,
			AllocsDeltaPercent: 0,
			IsNewBest:          false,
		}

		output := formatter.FormatComparison(comparison)

		assert.Contains(t, output, "⚠️")
		assert.Contains(t, output, "20.0% slower")
		assert.NotContains(t, output, "New personal best")
	})

	t.Run("shows memory improvement", func(t *testing.T) {
		formatter := NewFormatter()
		comparison := &ComparisonResult{
			TimeDeltaPercent:   0,
			MemoryDeltaPercent: -30.0,
			AllocsDeltaPercent: 0,
			IsNewBest:          false,
		}

		output := formatter.FormatComparison(comparison)

		assert.Contains(t, output, "💾")
		assert.Contains(t, output, "30.0% less memory")
	})

	t.Run("shows allocations improvement", func(t *testing.T) {
		formatter := NewFormatter()
		comparison := &ComparisonResult{
			TimeDeltaPercent:   0,
			MemoryDeltaPercent: 0,
			AllocsDeltaPercent: -25.0,
			IsNewBest:          false,
		}

		output := formatter.FormatComparison(comparison)

		assert.Contains(t, output, "✨")
		assert.Contains(t, output, "25.0% fewer allocations")
	})
}

func TestFormatter_FormatNumber(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		input    int
		expected string
	}{
		{123, "123"},
		{1234, "1,234"},
		{1234567, "1,234,567"},
		{1000000, "1,000,000"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatter.formatNumber(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatter_FormatTime(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		ns       float64
		contains string
	}{
		{500, "ns"},
		{1500, "µs"},
		{1500000, "ms"},
		{1500000000, "s"},
	}

	for _, tt := range tests {
		t.Run(tt.contains, func(t *testing.T) {
			result := formatter.formatTime(tt.ns)
			assert.True(t, strings.Contains(result, tt.contains))
		})
	}
}

func TestFormatter_FormatBytes(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		bytes    float64
		contains string
	}{
		{500, "B"},
		{2048, "KB"},
		{2097152, "MB"},
		{2147483648, "GB"},
	}

	for _, tt := range tests {
		t.Run(tt.contains, func(t *testing.T) {
			result := formatter.formatBytes(tt.bytes)
			assert.True(t, strings.Contains(result, tt.contains))
		})
	}
}

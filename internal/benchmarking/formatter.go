package benchmarking

import (
	"fmt"
	"strings"
)

// Formatter handles formatting and displaying benchmark results
type Formatter struct{}

// NewFormatter creates a new result formatter
func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatResult formats a benchmark result for display
func (f *Formatter) FormatResult(result *BenchmarkResult) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("\nResults:\n"))
	output.WriteString(fmt.Sprintf("  Iterations:       %s\n", f.formatNumber(result.Iterations)))
	output.WriteString(fmt.Sprintf("  Time per op:      %s\n", f.formatTime(result.NsPerOp)))
	output.WriteString(fmt.Sprintf("  Memory per op:    %s\n", f.formatBytes(result.BytesPerOp)))
	output.WriteString(fmt.Sprintf("  Allocs per op:    %.0f\n", result.AllocsPerOp))

	return output.String()
}

// FormatComparison formats a comparison between current and previous results
func (f *Formatter) FormatComparison(comparison *ComparisonResult) string {
	if comparison == nil {
		return "\nNo previous benchmark results for comparison.\n"
	}

	var output strings.Builder
	output.WriteString("\nComparison with previous best:\n")

	// Time comparison
	if comparison.TimeDeltaPercent != 0 {
		if comparison.TimeDeltaPercent < 0 {
			// Faster (improvement)
			output.WriteString(fmt.Sprintf("  🚀 %.1f%% faster than previous best!\n", -comparison.TimeDeltaPercent))
		} else {
			// Slower (regression)
			output.WriteString(fmt.Sprintf("  ⚠️  %.1f%% slower than previous best\n", comparison.TimeDeltaPercent))
		}
	}

	// Memory comparison
	if comparison.MemoryDeltaPercent != 0 {
		if comparison.MemoryDeltaPercent < 0 {
			// Less memory (improvement)
			output.WriteString(fmt.Sprintf("  💾 %.1f%% less memory than previous best!\n", -comparison.MemoryDeltaPercent))
		} else {
			// More memory (regression)
			output.WriteString(fmt.Sprintf("  ⚠️  %.1f%% more memory than previous best\n", comparison.MemoryDeltaPercent))
		}
	}

	// Allocations comparison
	if comparison.AllocsDeltaPercent != 0 {
		if comparison.AllocsDeltaPercent < 0 {
			// Fewer allocations (improvement)
			output.WriteString(fmt.Sprintf("  ✨ %.1f%% fewer allocations than previous best!\n", -comparison.AllocsDeltaPercent))
		} else {
			// More allocations (regression)
			output.WriteString(fmt.Sprintf("  ⚠️  %.1f%% more allocations than previous best\n", comparison.AllocsDeltaPercent))
		}
	}

	// New best indicator
	if comparison.IsNewBest {
		output.WriteString("\n✓ New personal best recorded!\n")
	}

	return output.String()
}

// formatNumber formats an integer with comma separators
func (f *Formatter) formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	// Add commas for thousands
	if len(s) <= 3 {
		return s
	}

	// Insert commas from right to left
	var result string
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}

// formatTime formats nanoseconds to appropriate time unit
func (f *Formatter) formatTime(ns float64) string {
	if ns < 1000 {
		return fmt.Sprintf("%.2f ns", ns)
	} else if ns < 1000000 {
		return fmt.Sprintf("%.3f µs", ns/1000)
	} else if ns < 1000000000 {
		return fmt.Sprintf("%.3f ms", ns/1000000)
	} else {
		return fmt.Sprintf("%.3f s", ns/1000000000)
	}
}

// formatBytes formats bytes to appropriate unit
func (f *Formatter) formatBytes(b float64) string {
	if b < 1024 {
		return fmt.Sprintf("%.0f B", b)
	} else if b < 1024*1024 {
		return fmt.Sprintf("%.2f KB", b/1024)
	} else if b < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", b/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", b/(1024*1024*1024))
	}
}

package benchmarking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutor_ParseBenchmarkOutput(t *testing.T) {
	t.Run("parses valid benchmark output", func(t *testing.T) {
		executor := NewExecutor()
		output := `goos: darwin
goarch: amd64
BenchmarkTwoSum-8    	1000000	      1234 ns/op	     512 B/op	       5 allocs/op
PASS
ok  	github.com/ak95asb/dsa-dojo/problems	1.234s`

		result, err := executor.parseBenchmarkOutput(output)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "TwoSum", result.BenchmarkName)
		assert.Equal(t, 1000000, result.Iterations)
		assert.Equal(t, 1234.0, result.NsPerOp)
		assert.Equal(t, 512.0, result.BytesPerOp)
		assert.Equal(t, 5.0, result.AllocsPerOp)
	})

	t.Run("parses output with decimal values", func(t *testing.T) {
		executor := NewExecutor()
		output := `BenchmarkBinarySearch-8    	5000000	       234.5 ns/op	      64.0 B/op	       2.0 allocs/op`

		result, err := executor.parseBenchmarkOutput(output)

		assert.NoError(t, err)
		assert.Equal(t, "BinarySearch", result.BenchmarkName)
		assert.Equal(t, 5000000, result.Iterations)
		assert.Equal(t, 234.5, result.NsPerOp)
		assert.Equal(t, 64.0, result.BytesPerOp)
		assert.Equal(t, 2.0, result.AllocsPerOp)
	})

	t.Run("returns error for invalid output", func(t *testing.T) {
		executor := NewExecutor()
		output := `No benchmark output here`

		result, err := executor.parseBenchmarkOutput(output)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no benchmark results found")
	})

	t.Run("handles different GOMAXPROCS values", func(t *testing.T) {
		executor := NewExecutor()
		output := `BenchmarkQuickSort-12    	2000000	       500 ns/op	     256 B/op	       3 allocs/op`

		result, err := executor.parseBenchmarkOutput(output)

		assert.NoError(t, err)
		assert.Equal(t, "QuickSort", result.BenchmarkName)
	})
}

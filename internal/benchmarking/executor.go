package benchmarking

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/ak95asb/dsa-dojo/internal/problem"
)

// Executor handles benchmark execution using go test
type Executor struct{}

// NewExecutor creates a new benchmark executor
func NewExecutor() *Executor {
	return &Executor{}
}

// BenchmarkResult represents parsed benchmark results
type BenchmarkResult struct {
	BenchmarkName string
	Iterations    int
	NsPerOp       float64
	BytesPerOp    float64
	AllocsPerOp   float64
	RawOutput     string
}

// ExecuteOptions contains options for benchmark execution
type ExecuteOptions struct {
	MemProfile bool
	CPUProfile string
	MemProfilePath string
}

// Execute runs benchmarks for a problem and returns parsed results
func (e *Executor) Execute(prob *problem.ProblemDetails, opts ExecuteOptions) (*BenchmarkResult, error) {
	// Construct test file path
	testFilePath := filepath.Join("problems", "templates", fmt.Sprintf("%s_test.go", prob.Slug))

	// Build go test command
	args := []string{"test", "-bench=.", "-benchmem", testFilePath}

	// Add profiling flags if requested
	if opts.CPUProfile != "" {
		args = append(args, fmt.Sprintf("-cpuprofile=%s", opts.CPUProfile))
	}
	if opts.MemProfilePath != "" {
		args = append(args, fmt.Sprintf("-memprofile=%s", opts.MemProfilePath))
	}

	// Execute command
	cmd := exec.Command("go", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("benchmark execution failed: %w\nOutput: %s", err, string(output))
	}

	// Parse benchmark output
	result, err := e.parseBenchmarkOutput(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to parse benchmark output: %w", err)
	}

	result.RawOutput = string(output)
	return result, nil
}

// parseBenchmarkOutput extracts benchmark metrics from go test output
func (e *Executor) parseBenchmarkOutput(output string) (*BenchmarkResult, error) {
	// Regex pattern for benchmark output:
	// Benchmark<Name>-<GOMAXPROCS>  <iterations>  <ns/op> ns/op  <B/op> B/op  <allocs/op> allocs/op
	pattern := `Benchmark(\w+)-\d+\s+(\d+)\s+([\d.]+) ns/op\s+([\d.]+) B/op\s+([\d.]+) allocs/op`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(output)
	if len(matches) < 6 {
		return nil, fmt.Errorf("no benchmark results found in output")
	}

	// Parse values
	iterations, _ := strconv.Atoi(matches[2])
	nsPerOp, _ := strconv.ParseFloat(matches[3], 64)
	bytesPerOp, _ := strconv.ParseFloat(matches[4], 64)
	allocsPerOp, _ := strconv.ParseFloat(matches[5], 64)

	return &BenchmarkResult{
		BenchmarkName: matches[1],
		Iterations:    iterations,
		NsPerOp:       nsPerOp,
		BytesPerOp:    bytesPerOp,
		AllocsPerOp:   allocsPerOp,
	}, nil
}

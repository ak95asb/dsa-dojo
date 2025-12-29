package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ak95asb/dsa-dojo/internal/benchmarking"
	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/spf13/cobra"
)

var (
	benchSave       bool
	benchMem        bool
	benchCPUProfile string
	benchMemProfile string
)

var benchCmd = &cobra.Command{
	Use:   "bench [problem-id]",
	Short: "Run performance benchmarks on your solution",
	Long: `Execute Go benchmarks and measure performance metrics.

The command:
  - Runs go test -bench on the problem's test file
  - Shows iterations, time per operation, allocations, memory
  - Optionally saves results and compares with previous best
  - Supports memory and CPU profiling

Examples:
  dsa bench two-sum
  dsa bench two-sum --save
  dsa bench two-sum --mem
  dsa bench two-sum --cpuprofile=two-sum.cpu.prof`,
	Args: cobra.ExactArgs(1),
	Run:  runBenchCommand,
}

func init() {
	rootCmd.AddCommand(benchCmd)
	benchCmd.Flags().BoolVar(&benchSave, "save", false, "Save benchmark results for comparison")
	benchCmd.Flags().BoolVar(&benchMem, "mem", false, "Enable memory profiling")
	benchCmd.Flags().StringVar(&benchCPUProfile, "cpuprofile", "", "Write CPU profile to file")
	benchCmd.Flags().StringVar(&benchMemProfile, "memprofile", "", "Write memory profile to file")
}

func runBenchCommand(cmd *cobra.Command, args []string) {
	slug := args[0]

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to connect to database: %v\n", err)
		os.Exit(3) // ExitDatabaseError
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Get problem by slug
	problemSvc := problem.NewService(db)
	prob, err := problemSvc.GetProblemBySlug(slug)
	if err != nil {
		if errors.Is(err, problem.ErrProblemNotFound) {
			fmt.Fprintf(os.Stderr, "Problem '%s' not found. Run 'dsa list' to see available problems.\n", slug)
			os.Exit(2) // ExitUsageError
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create benchmarking components
	executor := benchmarking.NewExecutor()
	formatter := benchmarking.NewFormatter()
	storage := benchmarking.NewStorage(db)
	comparator := benchmarking.NewComparator()

	// Build execution options
	opts := benchmarking.ExecuteOptions{
		MemProfile:     benchMem || benchMemProfile != "",
		CPUProfile:     benchCPUProfile,
		MemProfilePath: benchMemProfile,
	}

	// Execute benchmarks
	fmt.Printf("Running benchmarks for %s...\n\n", slug)
	result, err := executor.Execute(prob, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running benchmarks: %v\n", err)
		os.Exit(1)
	}

	// Display raw benchmark output
	fmt.Println(result.RawOutput)

	// Display formatted results
	fmt.Print(formatter.FormatResult(result))

	// Handle comparison if save flag is set or if previous benchmarks exist
	if benchSave {
		// Get previous best for comparison
		previousBest, err := storage.GetBestBenchmark(prob.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to retrieve previous benchmarks: %v\n", err)
		}

		// Compare results
		comparison := comparator.Compare(result, previousBest)
		fmt.Print(formatter.FormatComparison(comparison))

		// Save current result
		if err := storage.SaveBenchmark(prob.ID, result); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving benchmark: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✓ Benchmark results saved")
	}

	// Display profiling messages
	if benchCPUProfile != "" {
		fmt.Printf("\nCPU profile saved to %s\n", benchCPUProfile)
		fmt.Printf("View with: go tool pprof %s\n", benchCPUProfile)
		fmt.Printf("Analyze with: go tool pprof -http=:8080 %s\n", benchCPUProfile)
	}

	if benchMemProfile != "" {
		fmt.Printf("\nMemory profile saved to %s\n", benchMemProfile)
		fmt.Printf("View with: go tool pprof %s\n", benchMemProfile)
		fmt.Printf("Analyze with: go tool pprof -http=:8080 %s\n", benchMemProfile)
	}

	os.Exit(0)
}

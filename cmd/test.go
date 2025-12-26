package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/output"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/ak95asb/dsa-dojo/internal/progress"
	testingpkg "github.com/ak95asb/dsa-dojo/internal/testing"
	"github.com/spf13/cobra"
)

var (
	testVerbose bool
	testRace    bool
	testWatch   bool
)

var testCmd = &cobra.Command{
	Use:   "test [problem-id]",
	Short: "Run tests for a problem solution",
	Long: `Execute Go tests for your solution and display results.

The command:
  - Runs tests for the specified problem
  - Shows colored pass/fail status
  - Updates progress when all tests pass
  - Supports verbose and race detection modes

Examples:
  dsa test two-sum
  dsa test binary-search --verbose
  dsa test merge-intervals --race
  dsa test quick-sort --watch
  dsa test quick-sort --watch --verbose`,
	Args: cobra.ExactArgs(1),
	Run:  runTestCommand,
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().BoolVarP(&testVerbose, "verbose", "v", false, "Show detailed test output")
	testCmd.Flags().BoolVar(&testRace, "race", false, "Run tests with race detector")
	testCmd.Flags().BoolVarP(&testWatch, "watch", "w", false, "Watch for file changes and re-run tests")
}

func runTestCommand(cmd *cobra.Command, args []string) {
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

	// Create test service
	testSvc := testingpkg.NewService(db)

	// Route to watch mode if --watch flag is set
	if testWatch {
		if err := testSvc.Watch(prob, problemSvc, testVerbose, testRace); err != nil {
			fmt.Fprintf(os.Stderr, "Error in watch mode: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0) // Clean exit from watch mode
	}

	// Execute tests (normal mode)
	result, err := testSvc.ExecuteTests(prob, testVerbose, testRace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running tests: %v\n", err)
		os.Exit(1)
	}

	// Display results
	testSvc.DisplayResults(result)

	// Track progress (for both passed and failed tests)
	tracker := progress.NewTracker(db)
	filePath := fmt.Sprintf("problems/%s/solution.go", prob.Slug)
	isFirstTimeSolve, err := tracker.TrackTestCompletion(
		prob.ID,
		filePath,
		result.AllPassed,
		result.PassedCount,
		result.TotalCount,
	)
	if err != nil {
		// Log error but don't fail the command - progress tracking is non-critical
		fmt.Fprintf(os.Stderr, "Warning: Failed to update progress: %v\n", err)
	}

	// Display celebration message on first-time solve
	if result.AllPassed && isFirstTimeSolve {
		// Get the progress to determine attempt count
		var progressRecord database.Progress
		err = db.Where("problem_id = ?", prob.ID).First(&progressRecord).Error
		if err == nil {
			celebration := output.FormatCelebration(prob.Title, progressRecord.TotalAttempts)
			fmt.Println("\n" + celebration)
		} else {
			// Fallback to simple message if progress query fails
			fmt.Println("\n🎉 Congratulations! You solved it!")
		}
	} else if result.AllPassed {
		// Subsequent solve - simple message
		fmt.Println("\n✓ All tests passed!")
	}

	// Exit with appropriate code
	if result.AllPassed {
		os.Exit(0)
	} else {
		os.Exit(1) // Tests failed
	}
}

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/ak95asb/dsa-dojo/internal/solution"
	testingpkg "github.com/ak95asb/dsa-dojo/internal/testing"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit [problem-id]",
	Short: "Submit and save your solution to history",
	Long: `Submit your solution for a problem after verifying tests pass.

The command:
  - Runs tests to verify solution passes
  - Saves solution to solutions/history/<problem-id>/<timestamp>.go
  - Records submission in database with pass/fail status
  - Displays confirmation message

Examples:
  dsa submit two-sum
  dsa submit binary-search`,
	Args: cobra.ExactArgs(1),
	Run:  runSubmitCommand,
}

func init() {
	rootCmd.AddCommand(submitCmd)
}

func runSubmitCommand(cmd *cobra.Command, args []string) {
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

	// Check if solution file exists
	solutionPath := filepath.Join("solutions", fmt.Sprintf("%s.go", slug))
	if _, err := os.Stat(solutionPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Solution file not found: %s\n", solutionPath)
		fmt.Fprintf(os.Stderr, "Run 'dsa solve %s' to create a solution file.\n", slug)
		os.Exit(1)
	}

	// Run tests to verify solution
	fmt.Println("Running tests...")
	testSvc := testingpkg.NewService(db)
	result, err := testSvc.ExecuteTests(prob, false, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running tests: %v\n", err)
		os.Exit(1)
	}

	// Display test results
	testSvc.DisplayResults(result)

	// Submit solution (regardless of pass/fail)
	solutionSvc := solution.NewService(db)
	record, err := solutionSvc.RecordSubmission(
		slug,
		prob.ID,
		solutionPath,
		result.AllPassed,
		result.PassedCount,
		result.TotalCount,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error recording submission: %v\n", err)
		os.Exit(1)
	}

	// Display confirmation
	fmt.Println()
	if result.AllPassed {
		fmt.Printf("✓ Solution submitted and saved to history\n")
	} else {
		fmt.Printf("✓ Solution submitted (failed tests) and saved to history\n")
	}
	fmt.Printf("  Submission ID: %d\n", record.ID)
	fmt.Printf("  Status: ")
	if record.Passed {
		fmt.Printf("✓ Passed (%d/%d tests)\n", result.PassedCount, result.TotalCount)
	} else {
		fmt.Printf("✗ Failed (%d/%d tests)\n", result.PassedCount, result.TotalCount)
	}
	fmt.Printf("  Timestamp: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))

	os.Exit(0)
}

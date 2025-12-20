package cmd

import (
	"fmt"
	"os"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/output"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <problem-slug>",
	Short: "Display detailed information about a specific problem",
	Long: `Show displays comprehensive details about a problem including:
  - Problem metadata (title, difficulty, topic)
  - Full description
  - File paths for boilerplate and tests
  - Solution status and progress

Examples:
  dsa show two-sum              # Show details for Two Sum problem
  dsa show binary-search        # Show details for Binary Search problem`,
	Args: cobra.ExactArgs(1),
	Run:  runShowCommand,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShowCommand(cmd *cobra.Command, args []string) {
	problemSlug := args[0]

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

	// Create problem service
	svc := problem.NewService(db)

	// Get problem details
	problemDetails, err := svc.GetProblemBySlug(problemSlug)
	if err != nil {
		if err == problem.ErrProblemNotFound {
			fmt.Fprintf(os.Stderr, "Problem '%s' not found. Use 'dsa list' to see available problems.\n", problemSlug)
			os.Exit(2) // ExitUsageError
		}
		fmt.Fprintf(os.Stderr, "Error: Failed to retrieve problem: %v\n", err)
		os.Exit(1) // ExitGeneralError
	}

	// Format and display problem details
	output.PrintProblemDetails(problemDetails)
}

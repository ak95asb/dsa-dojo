package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ak95asb/dsa-dojo/internal/database"
	editorpkg "github.com/ak95asb/dsa-dojo/internal/editor"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/ak95asb/dsa-dojo/internal/solution"
	"github.com/spf13/cobra"
)

var (
	solveOpen  bool
	solveForce bool
)

var solveCmd = &cobra.Command{
	Use:   "solve [problem-id]",
	Short: "Start solving a problem with generated boilerplate",
	Long: `Generate a solution file with boilerplate code and function signature.

The command creates:
  - A solution file at solutions/<slug>.go
  - Boilerplate with function signature and helpful comments
  - Optional: Opens the file in your configured editor

Examples:
  dsa solve two-sum
  dsa solve binary-search --open
  dsa solve merge-intervals --force`,
	Args: cobra.ExactArgs(1),
	Run:  runSolveCommand,
}

func init() {
	rootCmd.AddCommand(solveCmd)
	solveCmd.Flags().BoolVarP(&solveOpen, "open", "o", false, "Open solution in editor after generation")
	solveCmd.Flags().BoolVarP(&solveForce, "force", "f", false, "Overwrite existing solution without confirmation")
}

func runSolveCommand(cmd *cobra.Command, args []string) {
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

	// Generate solution file
	solutionSvc := solution.NewService(db)
	solutionPath, err := solutionSvc.GenerateSolution(&prob.Problem, solveForce)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating solution: %v\n", err)
		os.Exit(1)
	}

	// Update progress tracking
	if err := problemSvc.UpdateProgress(prob.ID, "in_progress"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to update progress: %v\n", err)
	}

	fmt.Printf("✓ Solution file generated: %s\n", solutionPath)

	// Open in editor if requested
	if solveOpen {
		editorCmd := editorpkg.Detect()
		if err := editorpkg.Launch(editorCmd, solutionPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to open editor: %v\n", err)
		} else {
			fmt.Printf("✓ Opened in %s\n", editorCmd)
		}
	}
}

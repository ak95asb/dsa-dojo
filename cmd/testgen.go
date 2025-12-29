package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/ak95asb/dsa-dojo/internal/testgen"
	"github.com/spf13/cobra"
)

var (
	testGenAppend   bool
	testGenFromFile string
)

var testGenCmd = &cobra.Command{
	Use:   "test-gen [problem-id]",
	Short: "Generate test cases for a custom problem",
	Long: `Interactively generate test cases or import from JSON file.

The command:
  - Prompts for test case inputs interactively (default)
  - Imports test cases from JSON file (--from-file)
  - Appends to existing test file (--append)
  - Generates table-driven tests with testify/assert
  - Follows Go testing conventions

Examples:
  dsa test-gen my-problem
  dsa test-gen my-problem --from-file tests.json
  dsa test-gen my-problem --append
  dsa test-gen my-problem --append --from-file tests.json`,
	Args: cobra.ExactArgs(1),
	Run:  runTestGenCommand,
}

func init() {
	rootCmd.AddCommand(testGenCmd)
	testGenCmd.Flags().BoolVarP(&testGenAppend, "append", "a", false, "Append to existing test file")
	testGenCmd.Flags().StringVarP(&testGenFromFile, "from-file", "f", "", "Import test cases from JSON file")
}

func runTestGenCommand(cmd *cobra.Command, args []string) {
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

	// Create testgen service
	testGenSvc := testgen.NewService()

	// Route to appropriate mode
	if testGenFromFile != "" {
		// JSON file import mode
		if err := testGenSvc.GenerateFromFile(prob, testGenFromFile, testGenAppend); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating tests from file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Interactive mode
		if err := testGenSvc.GenerateInteractive(prob, testGenAppend); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating tests interactively: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("\n✅ Test file generated successfully!")
	os.Exit(0)
}

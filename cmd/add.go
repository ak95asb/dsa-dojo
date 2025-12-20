package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/output"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/ak95asb/dsa-dojo/internal/scaffold"
	"github.com/spf13/cobra"
)

var (
	addDifficulty string
	addTopic      string
	addTags       string
)

var addCmd = &cobra.Command{
	Use:   "add [problem title]",
	Short: "Add a custom problem to your library",
	Long: `Add creates a custom problem with boilerplate code and test file.

The command creates:
  - A new problem entry in the database
  - A boilerplate Go file at problems/<slug>.go
  - A test file at problems/<slug>_test.go

Examples:
  dsa add "Two Sum" --difficulty easy --topic arrays
  dsa add "Custom DFS Problem" --difficulty hard --topic graphs --tags "dfs,backtracking"`,
	Args: cobra.ExactArgs(1), // Require problem title
	Run:  runAddCommand,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&addDifficulty, "difficulty", "", "Difficulty level (easy, medium, hard) [required]")
	addCmd.Flags().StringVar(&addTopic, "topic", "", "Problem topic (arrays, linked-lists, trees, etc.) [required]")
	addCmd.Flags().StringVar(&addTags, "tags", "", "Comma-separated tags (optional)")
	addCmd.MarkFlagRequired("difficulty")
	addCmd.MarkFlagRequired("topic")
}

func runAddCommand(cmd *cobra.Command, args []string) {
	title := args[0]

	// Validate flags
	if !problem.IsValidDifficulty(addDifficulty) {
		fmt.Fprintf(os.Stderr, "Invalid difficulty '%s'. Valid options: easy, medium, hard\n", addDifficulty)
		os.Exit(2) // ExitUsageError
	}

	if !problem.IsValidTopic(addTopic) {
		fmt.Fprintf(os.Stderr, "Invalid topic '%s'. Valid topics: arrays, linked-lists, trees, graphs, sorting, searching\n", addTopic)
		os.Exit(2) // ExitUsageError
	}

	// Prompt for description (interactive)
	fmt.Println("Enter problem description (press Ctrl+D or Ctrl+Z when done):")
	description, err := readMultilineInput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading description: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(description) == "" {
		fmt.Fprintf(os.Stderr, "Description cannot be empty\n")
		os.Exit(2)
	}

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

	// Create problem
	newProblem, err := svc.CreateProblem(problem.CreateProblemInput{
		Title:       title,
		Difficulty:  addDifficulty,
		Topic:       addTopic,
		Description: description,
		Tags:        addTags,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating problem: %v\n", err)
		os.Exit(1)
	}

	// Generate boilerplate and test files
	generator := scaffold.NewGenerator()
	boilerplatePath, err := generator.GenerateBoilerplate(newProblem)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating boilerplate: %v\n", err)
		os.Exit(1)
	}

	testPath, err := generator.GenerateTestFile(newProblem)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating test file: %v\n", err)
		os.Exit(1)
	}

	// Display success message
	output.PrintProblemCreated(newProblem, boilerplatePath, testPath)
}

// readMultilineInput reads multi-line input from stdin until EOF
func readMultilineInput() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

package cmd

import (
	"fmt"
	"os"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/output"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/spf13/cobra"
)

var (
	randomDifficulty string
	randomTopic      string
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Get a random unsolved problem suggestion",
	Long: `Random selects a random unsolved problem from your library.

Use filters to narrow down the selection:
  --difficulty: Filter by difficulty (easy, medium, hard)
  --topic: Filter by topic (arrays, linked-lists, trees, graphs, sorting, searching)

Examples:
  dsa random                              # Any random unsolved problem
  dsa random --difficulty easy            # Random easy problem
  dsa random --topic arrays               # Random array problem
  dsa random --difficulty hard --topic trees  # Random hard tree problem`,
	Run: runRandomCommand,
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.Flags().StringVar(&randomDifficulty, "difficulty", "", "Filter by difficulty (easy, medium, hard)")
	randomCmd.Flags().StringVar(&randomTopic, "topic", "", "Filter by topic")
}

func runRandomCommand(cmd *cobra.Command, args []string) {
	// Validate flags
	if randomDifficulty != "" && !problem.IsValidDifficulty(randomDifficulty) {
		fmt.Fprintf(os.Stderr, "Invalid difficulty '%s'. Valid options: easy, medium, hard\n", randomDifficulty)
		os.Exit(2) // ExitUsageError
	}

	if randomTopic != "" && !problem.IsValidTopic(randomTopic) {
		fmt.Fprintf(os.Stderr, "Invalid topic '%s'. Use 'dsa list' to see available topics\n", randomTopic)
		os.Exit(2) // ExitUsageError
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

	// Build filters for unsolved problems
	filters := problem.ListFilters{
		Difficulty: randomDifficulty,
		Topic:      randomTopic,
	}
	solved := false
	filters.Solved = &solved // Only unsolved problems

	// Get random problem
	randomProblem, err := svc.GetRandomProblem(filters)
	if err != nil {
		if err == problem.ErrNoProblemsFound {
			// Generate helpful error message
			output.PrintNoProblemsMessage(filters)
			os.Exit(2) // ExitUsageError (no problems is a usage issue, not an error)
		}
		fmt.Fprintf(os.Stderr, "Error: Failed to retrieve random problem: %v\n", err)
		os.Exit(1) // ExitGeneralError
	}

	// Format and display random problem
	output.PrintRandomProblem(randomProblem)
}

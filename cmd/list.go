/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/spf13/cobra"
)

var (
	listDifficulty string
	listTopic      string
	listSolved     bool
	listUnsolved   bool
	listFormat     string
	listCompact    bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available problems with optional filters",
	Long: `List all problems in your library with optional filtering by difficulty, topic, or completion status.

Examples:
  dsa list                                  # List all problems
  dsa list --difficulty easy                # List only Easy problems
  dsa list --topic arrays                   # List only Array problems
  dsa list --difficulty medium --topic trees  # Combined filters
  dsa list --unsolved                       # List only unsolved problems
  dsa list --format json                    # Output as formatted JSON
  dsa list --format json --compact          # Output as compact JSON (single line)
  dsa list --format csv                     # Output as CSV
  dsa list --format csv > problems.csv      # Export to CSV file`,
	Run: runListCommand,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().StringVarP(&listDifficulty, "difficulty", "d", "", "Filter by difficulty (easy, medium, hard)")
	listCmd.Flags().StringVarP(&listTopic, "topic", "t", "", "Filter by topic (arrays, linked-lists, trees, graphs, sorting, searching)")
	listCmd.Flags().BoolVar(&listSolved, "solved", false, "Show only solved problems")
	listCmd.Flags().BoolVar(&listUnsolved, "unsolved", false, "Show only unsolved problems")
	listCmd.Flags().StringVar(&listFormat, "format", "table", "Output format (table, json)")
	listCmd.Flags().BoolVar(&listCompact, "compact", false, "Compact JSON output (no indentation)")
}

func runListCommand(cmd *cobra.Command, args []string) {
	// Validate format
	if !isValidFormat(listFormat) {
		fmt.Fprintf(os.Stderr, "Error: Invalid format '%s'. Valid formats: table, json, csv\n", listFormat)
		os.Exit(2)
	}

	// Validate conflicting flags
	if listSolved && listUnsolved {
		fmt.Fprintln(os.Stderr, "Error: Cannot use both --solved and --unsolved flags")
		os.Exit(2) // ExitUsageError
	}

	// Validate difficulty
	if listDifficulty != "" && !problem.IsValidDifficulty(listDifficulty) {
		fmt.Fprintf(os.Stderr, "Error: Invalid difficulty '%s'. Must be one of: easy, medium, hard\n", listDifficulty)
		os.Exit(2)
	}

	// Validate topic
	if listTopic != "" && !problem.IsValidTopic(listTopic) {
		fmt.Fprintf(os.Stderr, "Error: Invalid topic '%s'. Must be one of: arrays, linked-lists, trees, graphs, sorting, searching\n", listTopic)
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

	// Build filters
	filters := problem.ListFilters{
		Difficulty: listDifficulty,
		Topic:      listTopic,
		Solved:     nil, // Will be set based on flags
	}

	if listSolved {
		solved := true
		filters.Solved = &solved
	} else if listUnsolved {
		solved := false
		filters.Solved = &solved
	}

	// Query problems
	problems, err := svc.ListProblems(filters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to list problems: %v\n", err)
		os.Exit(1) // ExitGeneralError
	}

	// Handle empty results
	if len(problems) == 0 {
		if listFormat == "json" {
			// Empty JSON response
			response := ListResponse{
				Problems: []ProblemJSON{},
				Total:    0,
				Solved:   0,
			}
			if err := outputJSON(response, listCompact); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
		fmt.Println("No problems found matching the specified filters.")
		fmt.Println("Run 'dsa list' without filters to see all problems.")
		return
	}

	// CSV output
	if listFormat == "csv" {
		headers := []string{"ID", "Title", "Difficulty", "Topic", "Solved", "FirstSolvedAt"}
		rows := make([][]string, len(problems))

		for i, p := range problems {
			rows[i] = []string{
				p.Slug,                              // ID
				p.Title,                             // Title
				p.Difficulty,                        // Difficulty
				p.Topic,                             // Topic
				formatBoolCSV(p.IsSolved),           // Solved
				formatDateISO(p.FirstSolvedAt),      // FirstSolvedAt (empty if unsolved)
			}
		}

		if err := writeCSV(headers, rows); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// JSON output
	if listFormat == "json" {
		jsonProblems := make([]ProblemJSON, len(problems))
		solvedCount := 0
		for i, p := range problems {
			jsonProblems[i] = ProblemJSON{
				ID:         p.Slug,
				Title:      p.Title,
				Difficulty: p.Difficulty,
				Topic:      p.Topic,
				Solved:     p.IsSolved,
			}
			if p.IsSolved {
				solvedCount++
			}
		}

		response := ListResponse{
			Problems: jsonProblems,
			Total:    len(problems),
			Solved:   solvedCount,
		}

		if err := outputJSON(response, listCompact); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Format and display output (default table)
	// Convert to ProblemRow format
	problemRows := make([]ProblemRow, len(problems))
	for i, p := range problems {
		problemRows[i] = ProblemRow{
			ID:         p.Slug,
			Title:      p.Title,
			Difficulty: p.Difficulty,
			Topic:      p.Topic,
			Solved:     p.IsSolved,
		}
	}
	printProblemsTable(problemRows)
}

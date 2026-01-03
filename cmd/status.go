package cmd

import (
	"fmt"
	"os"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/output"
	"github.com/ak95asb/dsa-dojo/internal/progress"
	"github.com/spf13/cobra"
)

var (
	statusTopic   string
	statusCompact bool
	statusFormat  string
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display your problem-solving progress dashboard",
	Long: `Show an overview of your DSA practice progress.

The command displays:
  - Total problems solved (count and percentage)
  - Breakdown by difficulty level (Easy, Medium, Hard)
  - Breakdown by topic (Arrays, Trees, Graphs, etc.)
  - Recent activity (last 5 problems solved)
  - Visual progress bars with color coding

Examples:
  dsa status
  dsa status --topic arrays
  dsa status --compact
  dsa status --format json                 # Output as formatted JSON
  dsa status --format json --compact       # Output as compact JSON
  dsa status --format csv                  # Output as CSV
  dsa status --format csv > stats.csv      # Export stats to CSV file`,
	Args: cobra.NoArgs,
	Run:  runStatusCommand,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVar(&statusTopic, "topic", "", "Show stats for specific topic")
	statusCmd.Flags().BoolVar(&statusCompact, "compact", false, "Display one-line summary")
	statusCmd.Flags().StringVar(&statusFormat, "format", "table", "Output format (table, json, csv)")
}

func runStatusCommand(cmd *cobra.Command, args []string) {
	// Validate format
	if !isValidFormat(statusFormat) {
		fmt.Fprintf(os.Stderr, "Error: Invalid format '%s'. Valid formats: table, json, csv\n", statusFormat)
		os.Exit(2)
	}

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize database: %v\n", err)
		os.Exit(3)
	}

	// Create progress service
	progressService := progress.NewService(db)

	// Get statistics
	stats, err := progressService.GetStats(statusTopic)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to calculate statistics: %v\n", err)
		os.Exit(3)
	}

	// CSV output
	if statusFormat == "csv" {
		headers := []string{"Category", "Value", "Total", "Solved", "Unsolved"}
		rows := [][]string{
			{"Overall", "Problems", fmt.Sprintf("%d", stats.TotalProblems), fmt.Sprintf("%d", stats.TotalSolved), fmt.Sprintf("%d", stats.TotalProblems-stats.TotalSolved)},
		}

		// Add difficulty breakdown
		for _, diff := range []string{"easy", "medium", "hard"} {
			if diffStats, ok := stats.ByDifficulty[diff]; ok {
				rows = append(rows, []string{
					"Difficulty",
					diff,
					fmt.Sprintf("%d", diffStats.Total),
					fmt.Sprintf("%d", diffStats.Solved),
					fmt.Sprintf("%d", diffStats.Total-diffStats.Solved),
				})
			}
		}

		// Add topic breakdown
		for topic, topicStats := range stats.ByTopic {
			rows = append(rows, []string{
				"Topic",
				topic,
				fmt.Sprintf("%d", topicStats.Total),
				fmt.Sprintf("%d", topicStats.Solved),
				fmt.Sprintf("%d", topicStats.Total-topicStats.Solved),
			})
		}

		if err := writeCSV(headers, rows); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// JSON output
	if statusFormat == "json" {
		// Build by_difficulty map
		byDifficulty := make(map[string]int)
		for diff, diffStats := range stats.ByDifficulty {
			byDifficulty[diff] = diffStats.Solved
		}

		// Build by_topic map
		byTopic := make(map[string]int)
		for topic, topicStats := range stats.ByTopic {
			byTopic[topic] = topicStats.Solved
		}

		// Build recent activity
		recentActivity := make([]RecentActivityJSON, len(stats.RecentActivity))
		for i, ra := range stats.RecentActivity {
			recentActivity[i] = RecentActivityJSON{
				ProblemID: ra.Slug,
				Title:     ra.Title,
				Date:      ra.SolvedAt.Format("2006-01-02T15:04:05Z07:00"), // RFC3339
				Passed:    true, // If it's in RecentActivity, it was solved/passed
			}
		}

		response := StatusResponse{
			TotalProblems:  stats.TotalProblems,
			ProblemsSolved: stats.TotalSolved,
			ByDifficulty:   byDifficulty,
			ByTopic:        byTopic,
			RecentActivity: recentActivity,
			// Streak is Phase 2, omit for now (will be 0 and omitted due to omitempty)
		}

		if err := outputJSON(response, statusCompact); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Table output with statistics tables
	if statusCompact {
		// Use dashboard compact format for one-line summary
		dashboard := output.NewDashboard(stats, statusCompact, statusTopic)
		dashboardOutput := dashboard.Render()
		fmt.Print(dashboardOutput)
	} else {
		// Use table format for full display
		// Print header
		if statusTopic != "" {
			fmt.Printf("DSA Progress Dashboard - %s Topic\n\n", statusTopic)
		} else {
			fmt.Println("DSA Progress Dashboard")
			fmt.Println()
		}

		// Print overall summary
		percentage := 0
		if stats.TotalProblems > 0 {
			percentage = (stats.TotalSolved * 100) / stats.TotalProblems
		}
		fmt.Printf("Overall Progress: %d/%d problems solved (%d%%)\n\n",
			stats.TotalSolved, stats.TotalProblems, percentage)

		// Convert difficulty stats to StatsRow format
		difficultyStats := make(map[string]StatsRow)
		for diff, diffStats := range stats.ByDifficulty {
			difficultyStats[diff] = StatsRow{
				Category: diff,
				Total:    diffStats.Total,
				Solved:   diffStats.Solved,
				Unsolved: diffStats.Total - diffStats.Solved,
			}
		}

		// Convert topic stats to StatsRow format
		topicStats := make(map[string]StatsRow)
		for topic, topicData := range stats.ByTopic {
			topicStats[topic] = StatsRow{
				Category: topic,
				Total:    topicData.Total,
				Solved:   topicData.Solved,
				Unsolved: topicData.Total - topicData.Solved,
			}
		}

		// Print stats tables
		printStatsTable(difficultyStats, topicStats)

		// Print recent activity
		if len(stats.RecentActivity) > 0 {
			fmt.Println("\nRecent Activity:")
			for _, recent := range stats.RecentActivity {
				checkmark := colorize("✓", ColorGreen)
				dateStr := recent.SolvedAt.Format("2006-01-02")
				fmt.Printf("  %s %s (%s) - %s\n",
					checkmark, recent.Title, colorDifficulty(recent.Difficulty), dateStr)
			}
		}
	}
}

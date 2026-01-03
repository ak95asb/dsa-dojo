package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ak95asb/dsa-dojo/internal/analytics"
	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/output"
	"github.com/spf13/cobra"
)

var (
	analyticsTopic      string
	analyticsDifficulty string
	analyticsJSON       bool
)

var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "View detailed analytics and insights about your practice patterns",
	Long: `Display comprehensive analytics about your DSA practice performance.

The command shows:
  - Overall success rate across all attempted problems
  - Success rates broken down by difficulty level (Easy, Medium, Hard)
  - Success rates broken down by topic (Arrays, Trees, Graphs, etc.)
  - Average number of attempts needed to solve problems
  - Practice pattern insights (most/least practiced topics, strengths/weaknesses)

You can filter analytics by topic or difficulty, and export results as JSON.

Examples:
  dsa analytics                      # Show all analytics
  dsa analytics --topic arrays       # Analytics for arrays only
  dsa analytics --difficulty medium  # Analytics for medium problems
  dsa analytics --json              # Output as JSON
  dsa analytics --topic strings --json`,
	Args: cobra.NoArgs,
	Run:  runAnalyticsCommand,
}

func init() {
	rootCmd.AddCommand(analyticsCmd)
	analyticsCmd.Flags().StringVar(&analyticsTopic, "topic", "", "Filter analytics by topic")
	analyticsCmd.Flags().StringVar(&analyticsDifficulty, "difficulty", "", "Filter analytics by difficulty (easy, medium, hard)")
	analyticsCmd.Flags().BoolVar(&analyticsJSON, "json", false, "Output analytics as JSON")
}

func runAnalyticsCommand(cmd *cobra.Command, args []string) {
	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize database: %v\n", err)
		os.Exit(3)
	}

	// Create analytics service
	analyticsService := analytics.NewAnalyticsService(db)

	// Build filter
	filter := analytics.AnalyticsFilter{
		Topic:      analyticsTopic,
		Difficulty: analyticsDifficulty,
	}

	// Calculate statistics
	stats, err := analyticsService.CalculateStats(filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to calculate analytics: %v\n", err)
		os.Exit(3)
	}

	// Output results
	if analyticsJSON {
		// JSON output
		jsonData, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to marshal JSON: %v\n", err)
			os.Exit(3)
		}
		fmt.Println(string(jsonData))
	} else {
		// Dashboard output
		formatter := output.NewAnalyticsFormatter(stats, filter)
		result := formatter.Render()
		fmt.Print(result)
	}
}

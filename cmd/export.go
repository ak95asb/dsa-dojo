package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/export"
	"github.com/spf13/cobra"
)

var (
	exportFormat     string
	exportOutput     string
	exportDifficulty string
	exportTopic      string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export progress data to external formats",
	Long: `Export your practice progress to JSON or CSV format.

The command supports:
  - JSON export with full details (problems, progress, solutions, analytics)
  - CSV export for spreadsheet compatibility
  - Filtering by difficulty and topic
  - Output to file or stdout (for piping)

Examples:
  dsa export --format json --output progress.json
  dsa export --format csv --output progress.csv
  dsa export --format json --difficulty medium
  dsa export --format json | jq .summary`,
	Args: cobra.NoArgs,
	Run:  runExportCommand,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVar(&exportFormat, "format", "json", "Export format (json or csv)")
	exportCmd.Flags().StringVar(&exportOutput, "output", "", "Output file (default: stdout)")
	exportCmd.Flags().StringVar(&exportDifficulty, "difficulty", "", "Filter by difficulty (easy, medium, hard)")
	exportCmd.Flags().StringVar(&exportTopic, "topic", "", "Filter by topic")
}

func runExportCommand(cmd *cobra.Command, args []string) {
	// Validate format
	if exportFormat != "json" && exportFormat != "csv" {
		fmt.Fprintf(os.Stderr, "Error: Invalid format '%s'. Must be 'json' or 'csv'\n", exportFormat)
		os.Exit(2)
	}

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Create output writer
	var writer io.Writer
	if exportOutput == "" {
		writer = os.Stdout
	} else {
		file, err := os.Create(exportOutput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to create output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = file
	}

	// Create export service
	service := export.NewService(db)

	// Build filter
	filter := export.ExportFilter{
		Difficulty: exportDifficulty,
		Topic:      exportTopic,
	}

	// Export data
	var exportErr error
	if exportFormat == "json" {
		exportErr = service.ExportToJSON(filter, writer)
	} else {
		exportErr = service.ExportToCSV(filter, writer)
	}

	if exportErr != nil {
		fmt.Fprintf(os.Stderr, "Error: Export failed: %v\n", exportErr)
		os.Exit(1)
	}

	// Success message to stderr (not stdout) - only if writing to file
	if exportOutput != "" {
		fmt.Fprintf(os.Stderr, "✓ Export completed: %s\n", exportOutput)
	}
}

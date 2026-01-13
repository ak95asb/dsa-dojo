package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

// getTerminalWidth returns the current terminal width
func getTerminalWidth() int {
	// First, check $COLUMNS environment variable
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}

	// Try to detect terminal size
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
		return width
	}

	// Default to 80 columns if detection fails
	return 80
}

// selectColumns returns column names based on terminal width
func selectColumns(width int) []string {
	if width >= 120 {
		// Full display: all columns
		return []string{"ID", "Title", "Difficulty", "Topic", "Status"}
	} else if width >= 80 {
		// Hide ID column
		return []string{"Title", "Difficulty", "Topic", "Status"}
	} else {
		// Narrow terminal: minimal columns
		return []string{"Title", "Difficulty", "Status"}
	}
}

// truncateText truncates text to maxLen with ellipsis
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return "..."
	}
	return text[:maxLen-3] + "..."
}

// calculateTitleWidth determines max width for title column
func calculateTitleWidth(termWidth, numCols int) int {
	// Reserve space for borders, padding, other columns
	reserved := 30 + (numCols * 3) // 3 chars per column for padding/borders
	available := termWidth - reserved

	// Title gets 50% of remaining space
	titleWidth := available / 2
	if titleWidth < 15 {
		titleWidth = 15 // Minimum readable width
	}
	if titleWidth > 40 {
		titleWidth = 40 // Maximum for better layout
	}
	return titleWidth
}

// calculateTopicWidth determines max width for topic column
func calculateTopicWidth(termWidth, numCols int) int {
	if termWidth >= 120 {
		return 20
	} else if termWidth >= 80 {
		return 15
	} else {
		return 0 // Topic column hidden
	}
}

// ProblemRow represents a problem for table display
type ProblemRow struct {
	ID         string
	Title      string
	Difficulty string
	Topic      string
	Solved     bool
}

// buildTableRow creates a table row based on selected columns
func buildTableRow(p ProblemRow, columns []string, titleMaxWidth, topicMaxWidth int) []string {
	row := []string{}
	for _, col := range columns {
		switch col {
		case "ID":
			row = append(row, p.ID)
		case "Title":
			title := truncateText(p.Title, titleMaxWidth)
			row = append(row, title)
		case "Difficulty":
			diff := colorDifficulty(p.Difficulty)
			row = append(row, diff)
		case "Topic":
			topic := truncateText(p.Topic, topicMaxWidth)
			row = append(row, topic)
		case "Status":
			status := colorStatus(p.Solved)
			row = append(row, status)
		}
	}
	return row
}

// printProblemsTable outputs problems as formatted table
func printProblemsTable(problems []ProblemRow) {
	termWidth := getTerminalWidth()
	columns := selectColumns(termWidth)

	// Create table
	table := tablewriter.NewWriter(os.Stdout)

	// Set header
	table.Header(columns)

	// Calculate max widths for truncation
	titleMaxWidth := calculateTitleWidth(termWidth, len(columns))
	topicMaxWidth := calculateTopicWidth(termWidth, len(columns))

	// Count solved problems
	solvedCount := 0
	for _, p := range problems {
		if p.Solved {
			solvedCount++
		}
	}

	// Add rows
	for _, p := range problems {
		row := buildTableRow(p, columns, titleMaxWidth, topicMaxWidth)
		table.Append(row)
	}

	table.Render()

	// Print summary
	fmt.Printf("\nTotal: %d problems (%d solved, %d unsolved)\n",
		len(problems), solvedCount, len(problems)-solvedCount)
}

// StatsRow represents a statistics row for table display
type StatsRow struct {
	Category string
	Total    int
	Solved   int
	Unsolved int
}

// printStatsTable outputs statistics as formatted tables
func printStatsTable(difficultyStats map[string]StatsRow, topicStats map[string]StatsRow) {
	// Print difficulty breakdown
	if len(difficultyStats) > 0 {
		fmt.Println("Progress by Difficulty:")
		diffTable := tablewriter.NewWriter(os.Stdout)
		diffTable.Header([]string{"Difficulty", "Total", "Solved", "Unsolved", "Progress"})

		// Add rows in order: easy, medium, hard
		for _, diff := range []string{"easy", "medium", "hard"} {
			if stats, ok := difficultyStats[diff]; ok {
				percentage := 0
				if stats.Total > 0 {
					percentage = (stats.Solved * 100) / stats.Total
				}
				row := []string{
					colorDifficulty(stats.Category),
					fmt.Sprintf("%d", stats.Total),
					fmt.Sprintf("%d", stats.Solved),
					fmt.Sprintf("%d", stats.Unsolved),
					fmt.Sprintf("%d%%", percentage),
				}
				diffTable.Append(row)
			}
		}
		diffTable.Render()
		fmt.Println()
	}

	// Print topic breakdown
	if len(topicStats) > 0 {
		fmt.Println("Progress by Topic:")
		topicTable := tablewriter.NewWriter(os.Stdout)
		topicTable.Header([]string{"Topic", "Total", "Solved", "Unsolved", "Progress"})

		for _, stats := range topicStats {
			percentage := 0
			if stats.Total > 0 {
				percentage = (stats.Solved * 100) / stats.Total
			}
			row := []string{
				stats.Category,
				fmt.Sprintf("%d", stats.Total),
				fmt.Sprintf("%d", stats.Solved),
				fmt.Sprintf("%d", stats.Unsolved),
				fmt.Sprintf("%d%%", percentage),
			}
			topicTable.Append(row)
		}
		topicTable.Render()
	}
}

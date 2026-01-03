package output

import (
	"fmt"
	"strings"

	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/fatih/color"
)

// PrintProblemsTable formats and displays problems as a color-coded table
// Difficulty colors: Easy (green), Medium (yellow), Hard (red)
// Status colors: Solved (green checkmark), Unsolved (white)
func PrintProblemsTable(problems []problem.ProblemWithStatus) {
	// Color definitions (fatih/color respects NO_COLOR environment variable automatically)
	greenColor := color.New(color.FgGreen).SprintFunc()
	yellowColor := color.New(color.FgYellow).SprintFunc()
	redColor := color.New(color.FgRed).SprintFunc()

	// Calculate column widths
	slugWidth, titleWidth, diffWidth, topicWidth, statusWidth := 10, 30, 10, 15, 10
	for _, p := range problems {
		if len(p.Slug) > slugWidth {
			slugWidth = len(p.Slug)
		}
		if len(p.Title) > titleWidth {
			titleWidth = len(p.Title)
		}
		if len(p.Topic) > topicWidth {
			topicWidth = len(p.Topic)
		}
	}

	// Print header
	fmt.Println(strings.Repeat("-", slugWidth+titleWidth+diffWidth+topicWidth+statusWidth+13))
	fmt.Printf("| %-*s | %-*s | %-*s | %-*s | %-*s |\n",
		slugWidth, "Slug",
		titleWidth, "Title",
		diffWidth, "Difficulty",
		topicWidth, "Topic",
		statusWidth, "Status")
	fmt.Println(strings.Repeat("-", slugWidth+titleWidth+diffWidth+topicWidth+statusWidth+13))

	// Print rows with color-coded difficulty and status
	for _, p := range problems {
		// Format difficulty with color coding
		var difficultyStr string
		switch p.Difficulty {
		case "easy":
			difficultyStr = greenColor("Easy")
		case "medium":
			difficultyStr = yellowColor("Medium")
		case "hard":
			difficultyStr = redColor("Hard")
		default:
			difficultyStr = p.Difficulty
		}

		// Format status with color (green checkmark for solved)
		var statusStr string
		if p.IsSolved {
			statusStr = greenColor("✓ Solved")
		} else {
			statusStr = "Unsolved"
		}

		// Truncate title if too long
		title := p.Title
		if len(title) > titleWidth {
			title = title[:titleWidth-3] + "..."
		}

		fmt.Printf("| %-*s | %-*s | %-*s | %-*s | %-*s |\n",
			slugWidth, p.Slug,
			titleWidth, title,
			diffWidth, difficultyStr,
			topicWidth, p.Topic,
			statusWidth, statusStr)
	}
	fmt.Println(strings.Repeat("-", slugWidth+titleWidth+diffWidth+topicWidth+statusWidth+13))

	// Print summary statistics
	solvedCount := 0
	for _, p := range problems {
		if p.IsSolved {
			solvedCount++
		}
	}
	fmt.Printf("\nTotal: %d problems (%d solved, %d unsolved)\n",
		len(problems), solvedCount, len(problems)-solvedCount)
}

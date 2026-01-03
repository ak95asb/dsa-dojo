package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/ak95asb/dsa-dojo/internal/progress"
	"github.com/fatih/color"
)

// Dashboard formats progress statistics as a visual dashboard
type Dashboard struct {
	stats     *progress.Stats
	compact   bool
	topicName string
}

// NewDashboard creates a new dashboard formatter
func NewDashboard(stats *progress.Stats, compact bool, topicName string) *Dashboard {
	return &Dashboard{
		stats:     stats,
		compact:   compact,
		topicName: topicName,
	}
}

// Render generates the formatted dashboard output
func (d *Dashboard) Render() string {
	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}

	if d.compact {
		return d.renderCompact()
	}
	return d.renderFull()
}

// renderFull generates the full dashboard output
func (d *Dashboard) renderFull() string {
	var output strings.Builder

	// Title
	if d.topicName != "" {
		output.WriteString(fmt.Sprintf("DSA Progress Dashboard - %s Topic\n\n", strings.Title(d.topicName)))
	} else {
		output.WriteString("DSA Progress Dashboard\n\n")
	}

	// Overall progress
	output.WriteString(d.formatOverallProgress())
	output.WriteString("\n")

	// Progress by difficulty
	output.WriteString("Progress by Difficulty:\n")
	for _, difficulty := range []string{"easy", "medium", "hard"} {
		if stats, ok := d.stats.ByDifficulty[difficulty]; ok {
			output.WriteString(d.formatDifficultyLine(difficulty, stats))
		}
	}
	output.WriteString("\n")

	// Progress by topic (only if not filtered)
	if d.topicName == "" && len(d.stats.ByTopic) > 0 {
		output.WriteString("Progress by Topic:\n")
		for topic, stats := range d.stats.ByTopic {
			output.WriteString(d.formatTopicLine(topic, stats))
		}
		output.WriteString("\n")
	}

	// Recent activity
	if len(d.stats.RecentActivity) > 0 {
		output.WriteString("Recent Activity:\n")
		for _, recent := range d.stats.RecentActivity {
			checkmark := color.GreenString("✓")
			dateStr := recent.SolvedAt.Format("2006-01-02")
			output.WriteString(fmt.Sprintf("  %s %s (%s) - %s\n",
				checkmark, recent.Title, strings.Title(recent.Difficulty), dateStr))
		}
	}

	return output.String()
}

// renderCompact generates compact one-line summary
func (d *Dashboard) renderCompact() string {
	percentage := 0
	if d.stats.TotalProblems > 0 {
		percentage = (d.stats.TotalSolved * 100) / d.stats.TotalProblems
	}

	var parts []string

	// Overall progress
	parts = append(parts, fmt.Sprintf("Progress: %d/%d [%d%%]",
		d.stats.TotalSolved, d.stats.TotalProblems, percentage))

	// By difficulty
	var diffParts []string
	for _, diff := range []string{"easy", "medium", "hard"} {
		if stats, ok := d.stats.ByDifficulty[diff]; ok {
			label := "Easy"
			if diff == "medium" {
				label = "Med"
			} else if diff == "hard" {
				label = "Hard"
			}
			diffParts = append(diffParts, fmt.Sprintf("%s: %d/%d", label, stats.Solved, stats.Total))
		}
	}
	if len(diffParts) > 0 {
		parts = append(parts, "("+strings.Join(diffParts, ", ")+")")
	}

	// Last solved
	if len(d.stats.RecentActivity) > 0 {
		recent := d.stats.RecentActivity[0]
		dateStr := recent.SolvedAt.Format("2006-01-02")
		parts = append(parts, fmt.Sprintf("| Last: %s (%s)", recent.Title, dateStr))
	}

	return strings.Join(parts, " ")
}

// formatOverallProgress formats the overall progress line with bar
func (d *Dashboard) formatOverallProgress() string {
	percentage := 0
	if d.stats.TotalProblems > 0 {
		percentage = (d.stats.TotalSolved * 100) / d.stats.TotalProblems
	}

	bar := d.formatProgressBar(d.stats.TotalSolved, d.stats.TotalProblems, 20)
	return fmt.Sprintf("Overall Progress: %d/%d [%d%%] %s\n",
		d.stats.TotalSolved, d.stats.TotalProblems, percentage, bar)
}

// formatDifficultyLine formats a difficulty progress line
func (d *Dashboard) formatDifficultyLine(difficulty string, stats progress.DifficultyStats) string {
	percentage := 0
	if stats.Total > 0 {
		percentage = (stats.Solved * 100) / stats.Total
	}

	bar := d.formatProgressBar(stats.Solved, stats.Total, 20)
	label := strings.Title(difficulty)

	// Pad label for alignment
	paddedLabel := fmt.Sprintf("%-7s", label+":")

	return fmt.Sprintf("  %s %2d/%2d [%3d%%] %s\n",
		paddedLabel, stats.Solved, stats.Total, percentage, bar)
}

// formatTopicLine formats a topic progress line
func (d *Dashboard) formatTopicLine(topic string, stats progress.TopicStats) string {
	percentage := 0
	if stats.Total > 0 {
		percentage = (stats.Solved * 100) / stats.Total
	}

	bar := d.formatProgressBar(stats.Solved, stats.Total, 20)
	label := strings.Title(topic) + ":"

	// Pad label for alignment
	paddedLabel := fmt.Sprintf("%-15s", label)

	return fmt.Sprintf("  %s %2d/%2d [%3d%%] %s\n",
		paddedLabel, stats.Solved, stats.Total, percentage, bar)
}

// formatProgressBar creates a colored progress bar using Unicode blocks
func (d *Dashboard) formatProgressBar(solved, total, width int) string {
	if total == 0 {
		return strings.Repeat("░", width)
	}

	percentage := float64(solved) / float64(total) * 100
	filledWidth := int(float64(width) * float64(solved) / float64(total))
	emptyWidth := width - filledWidth

	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)

	// Apply color based on percentage
	c := d.getColorForPercentage(percentage)
	return c.Sprint(filled + empty)
}

// getColorForPercentage returns the appropriate color for a percentage
func (d *Dashboard) getColorForPercentage(percentage float64) *color.Color {
	if percentage >= 70 {
		return color.New(color.FgGreen)
	} else if percentage >= 30 {
		return color.New(color.FgYellow)
	} else {
		return color.New(color.FgRed)
	}
}

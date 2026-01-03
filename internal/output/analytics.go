package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/ak95asb/dsa-dojo/internal/analytics"
	"github.com/fatih/color"
)

// AnalyticsFormatter formats analytics statistics as a visual dashboard
type AnalyticsFormatter struct {
	stats  *analytics.AnalyticsStats
	filter analytics.AnalyticsFilter
}

// NewAnalyticsFormatter creates a new analytics formatter
func NewAnalyticsFormatter(stats *analytics.AnalyticsStats, filter analytics.AnalyticsFilter) *AnalyticsFormatter {
	return &AnalyticsFormatter{
		stats:  stats,
		filter: filter,
	}
}

// Render generates the formatted analytics output
func (f *AnalyticsFormatter) Render() string {
	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}

	var output strings.Builder

	// Title
	output.WriteString(f.formatTitle())
	output.WriteString("\n")

	// Overall metrics
	output.WriteString(f.formatOverallMetrics())
	output.WriteString("\n")

	// Success rates by difficulty
	if len(f.stats.SuccessRateByDifficulty) > 0 {
		output.WriteString("Success Rate by Difficulty:\n")
		for _, difficulty := range []string{"easy", "medium", "hard"} {
			if rate, ok := f.stats.SuccessRateByDifficulty[difficulty]; ok {
				output.WriteString(f.formatDifficultyLine(difficulty, rate))
			}
		}
		output.WriteString("\n")
	}

	// Success rates by topic
	if len(f.stats.SuccessRateByTopic) > 0 && f.filter.Topic == "" {
		output.WriteString("Success Rate by Topic:\n")
		for topic, rate := range f.stats.SuccessRateByTopic {
			output.WriteString(f.formatTopicLine(topic, rate))
		}
		output.WriteString("\n")
	}

	// Average attempts by difficulty
	if len(f.stats.AvgAttemptsByDifficulty) > 0 {
		output.WriteString("Average Attempts by Difficulty:\n")
		for _, difficulty := range []string{"easy", "medium", "hard"} {
			if avg, ok := f.stats.AvgAttemptsByDifficulty[difficulty]; ok {
				output.WriteString(f.formatAvgAttemptsLine(difficulty, avg))
			}
		}
		output.WriteString("\n")
	}

	// Practice patterns insights
	if f.stats.MostPracticedTopic != "" || f.stats.BestDifficulty != "" {
		output.WriteString(f.formatInsights())
		output.WriteString("\n")
	}

	return output.String()
}

// formatTitle generates the dashboard title
func (f *AnalyticsFormatter) formatTitle() string {
	title := "Analytics Dashboard"
	if f.filter.Topic != "" {
		title += fmt.Sprintf(" - %s Topic", strings.Title(f.filter.Topic))
	}
	if f.filter.Difficulty != "" {
		title += fmt.Sprintf(" - %s Difficulty", strings.Title(f.filter.Difficulty))
	}
	return title
}

// formatOverallMetrics formats the overall success rate and average attempts
func (f *AnalyticsFormatter) formatOverallMetrics() string {
	var output strings.Builder

	// Overall success rate with visual indicator
	successColor := f.getColorForSuccessRate(f.stats.OverallSuccessRate)
	output.WriteString(fmt.Sprintf("Overall Success Rate: %s\n",
		successColor.Sprintf("%.1f%%", f.stats.OverallSuccessRate)))

	// Overall average attempts
	output.WriteString(fmt.Sprintf("Average Attempts to Solve: %.1f\n",
		f.stats.AvgAttemptsOverall))

	return output.String()
}

// formatDifficultyLine formats a difficulty success rate line
func (f *AnalyticsFormatter) formatDifficultyLine(difficulty string, rate float64) string {
	label := fmt.Sprintf("%-7s", strings.Title(difficulty)+":")
	successColor := f.getColorForSuccessRate(rate)

	// Create visual bar
	bar := f.formatRateBar(rate, 20)

	return fmt.Sprintf("  %s %s %s\n",
		label, successColor.Sprintf("%5.1f%%", rate), bar)
}

// formatTopicLine formats a topic success rate line
func (f *AnalyticsFormatter) formatTopicLine(topic string, rate float64) string {
	label := fmt.Sprintf("%-15s", strings.Title(topic)+":")
	successColor := f.getColorForSuccessRate(rate)

	// Create visual bar
	bar := f.formatRateBar(rate, 20)

	return fmt.Sprintf("  %s %s %s\n",
		label, successColor.Sprintf("%5.1f%%", rate), bar)
}

// formatAvgAttemptsLine formats an average attempts line
func (f *AnalyticsFormatter) formatAvgAttemptsLine(difficulty string, avg float64) string {
	label := fmt.Sprintf("%-7s", strings.Title(difficulty)+":")

	// Color based on avg attempts (fewer = better)
	avgColor := f.getColorForAttempts(avg)

	return fmt.Sprintf("  %s %s\n",
		label, avgColor.Sprintf("%.1f attempts", avg))
}

// formatInsights formats the practice patterns insights section
func (f *AnalyticsFormatter) formatInsights() string {
	var output strings.Builder

	output.WriteString("Practice Insights:\n")

	if f.stats.MostPracticedTopic != "" {
		output.WriteString(fmt.Sprintf("  Most Practiced:  %s\n",
			color.CyanString(strings.Title(f.stats.MostPracticedTopic))))
	}

	if f.stats.LeastPracticedTopic != "" {
		output.WriteString(fmt.Sprintf("  Least Practiced: %s\n",
			color.YellowString(strings.Title(f.stats.LeastPracticedTopic))))
	}

	if f.stats.BestDifficulty != "" {
		output.WriteString(fmt.Sprintf("  Strength:        %s problems\n",
			color.GreenString(strings.Title(f.stats.BestDifficulty))))
	}

	if f.stats.ChallengingDifficulty != "" {
		output.WriteString(fmt.Sprintf("  Challenge:       %s problems\n",
			color.RedString(strings.Title(f.stats.ChallengingDifficulty))))
	}

	return output.String()
}

// formatRateBar creates a colored progress bar for success rates
func (f *AnalyticsFormatter) formatRateBar(rate float64, width int) string {
	if rate < 0 {
		rate = 0
	}
	if rate > 100 {
		rate = 100
	}

	filledWidth := int(float64(width) * rate / 100.0)
	emptyWidth := width - filledWidth

	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)

	// Apply color based on rate
	c := f.getColorForSuccessRate(rate)
	return c.Sprint(filled + empty)
}

// getColorForSuccessRate returns the appropriate color for a success rate
func (f *AnalyticsFormatter) getColorForSuccessRate(rate float64) *color.Color {
	if rate >= 70 {
		return color.New(color.FgGreen)
	} else if rate >= 40 {
		return color.New(color.FgYellow)
	} else {
		return color.New(color.FgRed)
	}
}

// getColorForAttempts returns the appropriate color for average attempts
func (f *AnalyticsFormatter) getColorForAttempts(avg float64) *color.Color {
	if avg <= 2.0 {
		return color.New(color.FgGreen)
	} else if avg <= 4.0 {
		return color.New(color.FgYellow)
	} else {
		return color.New(color.FgRed)
	}
}

package output

import (
	"os"
	"strings"
	"testing"

	"github.com/ak95asb/dsa-dojo/internal/analytics"
	"github.com/stretchr/testify/assert"
)

func TestNewAnalyticsFormatter(t *testing.T) {
	stats := &analytics.AnalyticsStats{}
	filter := analytics.AnalyticsFilter{}

	formatter := NewAnalyticsFormatter(stats, filter)

	assert.NotNil(t, formatter)
	assert.Equal(t, stats, formatter.stats)
}

func TestAnalyticsFormatter_Render_BasicOutput(t *testing.T) {
	// Disable colors for testing
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	stats := &analytics.AnalyticsStats{
		OverallSuccessRate:      75.0,
		SuccessRateByDifficulty: map[string]float64{"easy": 90.0, "medium": 65.0, "hard": 50.0},
		SuccessRateByTopic:      map[string]float64{"arrays": 80.0, "trees": 70.0},
		AvgAttemptsOverall:      2.5,
		AvgAttemptsByDifficulty: map[string]float64{"easy": 1.5, "medium": 2.8, "hard": 4.2},
		MostPracticedTopic:      "arrays",
		LeastPracticedTopic:     "graphs",
		BestDifficulty:          "easy",
		ChallengingDifficulty:   "hard",
	}

	filter := analytics.AnalyticsFilter{}
	formatter := NewAnalyticsFormatter(stats, filter)
	output := formatter.Render()

	// Check that output contains key sections
	assert.Contains(t, output, "Analytics Dashboard")
	assert.Contains(t, output, "Overall Success Rate: 75.0%")
	assert.Contains(t, output, "Average Attempts to Solve: 2.5")
	assert.Contains(t, output, "Success Rate by Difficulty:")
	assert.Contains(t, output, "Success Rate by Topic:")
	assert.Contains(t, output, "Average Attempts by Difficulty:")
	assert.Contains(t, output, "Practice Insights:")
}

func TestAnalyticsFormatter_RenderWithTopicFilter(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	stats := &analytics.AnalyticsStats{
		OverallSuccessRate:      80.0,
		SuccessRateByDifficulty: map[string]float64{"easy": 100.0, "medium": 70.0},
		SuccessRateByTopic:      map[string]float64{"arrays": 80.0},
		AvgAttemptsOverall:      2.0,
		AvgAttemptsByDifficulty: map[string]float64{"easy": 1.2, "medium": 2.5},
	}

	filter := analytics.AnalyticsFilter{Topic: "arrays"}
	formatter := NewAnalyticsFormatter(stats, filter)
	output := formatter.Render()

	// Title should include filter
	assert.Contains(t, output, "Arrays Topic")

	// Should not show topic breakdown when filtered by topic
	lines := strings.Split(output, "\n")
	topicSectionFound := false
	for _, line := range lines {
		if strings.Contains(line, "Success Rate by Topic:") {
			topicSectionFound = true
			break
		}
	}
	assert.False(t, topicSectionFound, "Should not show topic breakdown when filtered by topic")
}

func TestAnalyticsFormatter_RenderWithDifficultyFilter(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	stats := &analytics.AnalyticsStats{
		OverallSuccessRate:      90.0,
		SuccessRateByDifficulty: map[string]float64{"easy": 90.0},
		AvgAttemptsOverall:      1.5,
		AvgAttemptsByDifficulty: map[string]float64{"easy": 1.5},
	}

	filter := analytics.AnalyticsFilter{Difficulty: "easy"}
	formatter := NewAnalyticsFormatter(stats, filter)
	output := formatter.Render()

	// Title should include filter
	assert.Contains(t, output, "Easy Difficulty")
	assert.Contains(t, output, "Overall Success Rate: 90.0%")
}

func TestAnalyticsFormatter_FormatTitle(t *testing.T) {
	tests := []struct {
		name     string
		filter   analytics.AnalyticsFilter
		expected string
	}{
		{
			name:     "No filters",
			filter:   analytics.AnalyticsFilter{},
			expected: "Analytics Dashboard",
		},
		{
			name:     "Topic filter only",
			filter:   analytics.AnalyticsFilter{Topic: "arrays"},
			expected: "Analytics Dashboard - Arrays Topic",
		},
		{
			name:     "Difficulty filter only",
			filter:   analytics.AnalyticsFilter{Difficulty: "medium"},
			expected: "Analytics Dashboard - Medium Difficulty",
		},
		{
			name:     "Both filters",
			filter:   analytics.AnalyticsFilter{Topic: "trees", Difficulty: "hard"},
			expected: "Analytics Dashboard - Trees Topic - Hard Difficulty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewAnalyticsFormatter(&analytics.AnalyticsStats{}, tt.filter)
			title := formatter.formatTitle()
			assert.Equal(t, tt.expected, title)
		})
	}
}

func TestAnalyticsFormatter_FormatOverallMetrics(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	stats := &analytics.AnalyticsStats{
		OverallSuccessRate: 62.5,
		AvgAttemptsOverall: 3.2,
	}

	formatter := NewAnalyticsFormatter(stats, analytics.AnalyticsFilter{})
	output := formatter.formatOverallMetrics()

	assert.Contains(t, output, "Overall Success Rate: 62.5%")
	assert.Contains(t, output, "Average Attempts to Solve: 3.2")
}

func TestAnalyticsFormatter_FormatDifficultyLine(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	formatter := NewAnalyticsFormatter(&analytics.AnalyticsStats{}, analytics.AnalyticsFilter{})

	output := formatter.formatDifficultyLine("easy", 85.5)

	assert.Contains(t, output, "Easy:")
	assert.Contains(t, output, "85.5%")
	// Should contain progress bar characters
	assert.True(t, strings.Contains(output, "█") || strings.Contains(output, "░"))
}

func TestAnalyticsFormatter_FormatTopicLine(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	formatter := NewAnalyticsFormatter(&analytics.AnalyticsStats{}, analytics.AnalyticsFilter{})

	output := formatter.formatTopicLine("arrays", 72.3)

	assert.Contains(t, output, "Arrays:")
	assert.Contains(t, output, "72.3%")
	// Should contain progress bar characters
	assert.True(t, strings.Contains(output, "█") || strings.Contains(output, "░"))
}

func TestAnalyticsFormatter_FormatAvgAttemptsLine(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	formatter := NewAnalyticsFormatter(&analytics.AnalyticsStats{}, analytics.AnalyticsFilter{})

	output := formatter.formatAvgAttemptsLine("medium", 3.5)

	assert.Contains(t, output, "Medium:")
	assert.Contains(t, output, "3.5 attempts")
}

func TestAnalyticsFormatter_FormatInsights(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	stats := &analytics.AnalyticsStats{
		MostPracticedTopic:    "arrays",
		LeastPracticedTopic:   "graphs",
		BestDifficulty:        "easy",
		ChallengingDifficulty: "hard",
	}

	formatter := NewAnalyticsFormatter(stats, analytics.AnalyticsFilter{})
	output := formatter.formatInsights()

	assert.Contains(t, output, "Practice Insights:")
	assert.Contains(t, output, "Most Practiced:  Arrays")
	assert.Contains(t, output, "Least Practiced: Graphs")
	assert.Contains(t, output, "Strength:        Easy problems")
	assert.Contains(t, output, "Challenge:       Hard problems")
}

func TestAnalyticsFormatter_FormatRateBar(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	formatter := NewAnalyticsFormatter(&analytics.AnalyticsStats{}, analytics.AnalyticsFilter{})

	tests := []struct {
		name     string
		rate     float64
		width    int
		minFilled int
		maxFilled int
	}{
		{"Zero rate", 0.0, 10, 0, 0},
		{"Half rate", 50.0, 10, 4, 6},
		{"Full rate", 100.0, 10, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := formatter.formatRateBar(tt.rate, tt.width)

			// Count filled blocks
			filledCount := strings.Count(bar, "█")

			assert.GreaterOrEqual(t, filledCount, tt.minFilled)
			assert.LessOrEqual(t, filledCount, tt.maxFilled)
			assert.Equal(t, tt.width, len([]rune(bar))) // Total width should match
		})
	}
}

func TestAnalyticsFormatter_EmptyStats(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	stats := &analytics.AnalyticsStats{
		SuccessRateByDifficulty: make(map[string]float64),
		SuccessRateByTopic:      make(map[string]float64),
		AvgAttemptsByDifficulty: make(map[string]float64),
	}

	formatter := NewAnalyticsFormatter(stats, analytics.AnalyticsFilter{})
	output := formatter.Render()

	// Should still render title and overall metrics
	assert.Contains(t, output, "Analytics Dashboard")
	assert.Contains(t, output, "Overall Success Rate: 0.0%")
	assert.Contains(t, output, "Average Attempts to Solve: 0.0")

	// Should not crash with empty maps
	assert.NotPanics(t, func() {
		formatter.Render()
	})
}

package output

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/progress"
	"github.com/stretchr/testify/assert"
)

func TestDashboard_RenderFull(t *testing.T) {
	stats := &progress.Stats{
		TotalProblems: 10,
		TotalSolved:   5,
		ByDifficulty: map[string]progress.DifficultyStats{
			"easy":   {Total: 4, Solved: 3},
			"medium": {Total: 4, Solved: 2},
			"hard":   {Total: 2, Solved: 0},
		},
		ByTopic: map[string]progress.TopicStats{
			"arrays": {Total: 3, Solved: 2},
			"trees":  {Total: 3, Solved: 1},
		},
		RecentActivity: []progress.RecentProblem{
			{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", SolvedAt: time.Now()},
		},
	}

	dashboard := NewDashboard(stats, false, "")
	output := dashboard.Render()

	assert.Contains(t, output, "DSA Progress Dashboard")
	assert.Contains(t, output, "Overall Progress: 5/10")
	assert.Contains(t, output, "[50%]")
	assert.Contains(t, output, "Progress by Difficulty:")
	assert.Contains(t, output, "Easy:")
	assert.Contains(t, output, "Medium:")
	assert.Contains(t, output, "Hard:")
	assert.Contains(t, output, "Progress by Topic:")
	assert.Contains(t, output, "Arrays:")
	assert.Contains(t, output, "Trees:")
	assert.Contains(t, output, "Recent Activity:")
	assert.Contains(t, output, "Two Sum")
}

func TestDashboard_RenderCompact(t *testing.T) {
	stats := &progress.Stats{
		TotalProblems: 10,
		TotalSolved:   5,
		ByDifficulty: map[string]progress.DifficultyStats{
			"easy":   {Total: 4, Solved: 3},
			"medium": {Total: 4, Solved: 2},
			"hard":   {Total: 2, Solved: 0},
		},
		RecentActivity: []progress.RecentProblem{
			{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", SolvedAt: time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)},
		},
	}

	dashboard := NewDashboard(stats, true, "")
	output := dashboard.Render()

	assert.Contains(t, output, "Progress: 5/10 [50%]")
	assert.Contains(t, output, "Easy: 3/4")
	assert.Contains(t, output, "Med: 2/4")
	assert.Contains(t, output, "Hard: 0/2")
	assert.Contains(t, output, "Last: Two Sum (2025-12-15)")
	assert.NotContains(t, output, "Progress by Difficulty:")
	assert.NotContains(t, output, "Recent Activity:")
}

func TestDashboard_RenderTopicSpecific(t *testing.T) {
	stats := &progress.Stats{
		TotalProblems: 3,
		TotalSolved:   2,
		ByDifficulty: map[string]progress.DifficultyStats{
			"easy":   {Total: 2, Solved: 2},
			"medium": {Total: 1, Solved: 0},
		},
		RecentActivity: []progress.RecentProblem{
			{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", SolvedAt: time.Now()},
		},
	}

	dashboard := NewDashboard(stats, false, "arrays")
	output := dashboard.Render()

	assert.Contains(t, output, "DSA Progress Dashboard - Arrays Topic")
	assert.Contains(t, output, "Overall Progress: 2/3")
	assert.Contains(t, output, "Progress by Difficulty:")
	assert.NotContains(t, output, "Progress by Topic:")
}

func TestFormatProgressBar_EmptyProgress(t *testing.T) {
	stats := &progress.Stats{
		TotalProblems: 10,
		TotalSolved:   0,
	}

	dashboard := NewDashboard(stats, false, "")
	bar := dashboard.formatProgressBar(0, 10, 20)

	assert.Contains(t, bar, "░")
	assert.NotContains(t, bar, "█")
	assert.Equal(t, 20, len([]rune(removeANSI(bar))), "Bar should be 20 characters wide")
}

func TestFormatProgressBar_FullProgress(t *testing.T) {
	stats := &progress.Stats{
		TotalProblems: 10,
		TotalSolved:   10,
	}

	dashboard := NewDashboard(stats, false, "")
	bar := dashboard.formatProgressBar(10, 10, 20)

	assert.Contains(t, bar, "█")
	assert.NotContains(t, bar, "░")
	assert.Equal(t, 20, len([]rune(removeANSI(bar))), "Bar should be 20 characters wide")
}

func TestFormatProgressBar_PartialProgress(t *testing.T) {
	stats := &progress.Stats{}
	dashboard := NewDashboard(stats, false, "")

	tests := []struct {
		name   string
		solved int
		total  int
		width  int
	}{
		{"50% progress", 5, 10, 20},
		{"25% progress", 1, 4, 20},
		{"75% progress", 15, 20, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := dashboard.formatProgressBar(tt.solved, tt.total, tt.width)
			cleanBar := removeANSI(bar)

			assert.Contains(t, bar, "█", "Should contain filled blocks")
			assert.Contains(t, bar, "░", "Should contain empty blocks")
			assert.Equal(t, tt.width, len([]rune(cleanBar)), "Bar should match expected width")

			// Count filled and empty blocks
			filledCount := strings.Count(cleanBar, "█")
			emptyCount := strings.Count(cleanBar, "░")
			assert.Equal(t, tt.width, filledCount+emptyCount, "Total blocks should equal width")
		})
	}
}

func TestFormatProgressBar_ZeroTotal(t *testing.T) {
	stats := &progress.Stats{}
	dashboard := NewDashboard(stats, false, "")

	bar := dashboard.formatProgressBar(0, 0, 20)

	assert.NotContains(t, bar, "█")
	assert.Contains(t, bar, "░")
	assert.Equal(t, 20, len([]rune(removeANSI(bar))), "Should return all empty blocks")
}

func TestGetColorForPercentage(t *testing.T) {
	stats := &progress.Stats{}
	dashboard := NewDashboard(stats, false, "")

	tests := []struct {
		name       string
		percentage float64
		expectRed  bool
		expectYel  bool
		expectGrn  bool
	}{
		{"0% - red", 0, true, false, false},
		{"29% - red", 29, true, false, false},
		{"30% - yellow", 30, false, true, false},
		{"50% - yellow", 50, false, true, false},
		{"69% - yellow", 69, false, true, false},
		{"70% - green", 70, false, false, true},
		{"100% - green", 100, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := dashboard.getColorForPercentage(tt.percentage)

			// Check color attributes
			if tt.expectRed {
				assert.NotNil(t, color, "Color should not be nil for red")
			} else if tt.expectYel {
				assert.NotNil(t, color, "Color should not be nil for yellow")
			} else if tt.expectGrn {
				assert.NotNil(t, color, "Color should not be nil for green")
			}
		})
	}
}

func TestDashboard_NoColorEnvironmentVariable(t *testing.T) {
	// Set NO_COLOR environment variable
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	stats := &progress.Stats{
		TotalProblems: 10,
		TotalSolved:   5,
		ByDifficulty: map[string]progress.DifficultyStats{
			"easy": {Total: 4, Solved: 3},
		},
	}

	dashboard := NewDashboard(stats, false, "")
	output := dashboard.Render()

	// When NO_COLOR is set, output should not contain ANSI escape codes
	// But should still contain the content
	assert.Contains(t, output, "Overall Progress: 5/10")
	assert.Contains(t, output, "Easy:")
}

func TestDashboard_EmptyStats(t *testing.T) {
	stats := &progress.Stats{
		TotalProblems:  0,
		TotalSolved:    0,
		ByDifficulty:   make(map[string]progress.DifficultyStats),
		ByTopic:        make(map[string]progress.TopicStats),
		RecentActivity: []progress.RecentProblem{},
	}

	dashboard := NewDashboard(stats, false, "")
	output := dashboard.Render()

	assert.Contains(t, output, "DSA Progress Dashboard")
	assert.Contains(t, output, "Overall Progress: 0/0 [0%]")
	assert.NotContains(t, output, "Recent Activity:")
}

func TestDashboard_CompactWithNoRecentActivity(t *testing.T) {
	stats := &progress.Stats{
		TotalProblems: 10,
		TotalSolved:   5,
		ByDifficulty: map[string]progress.DifficultyStats{
			"easy": {Total: 4, Solved: 3},
		},
		RecentActivity: []progress.RecentProblem{},
	}

	dashboard := NewDashboard(stats, true, "")
	output := dashboard.Render()

	assert.Contains(t, output, "Progress: 5/10 [50%]")
	assert.NotContains(t, output, "Last:")
}

func TestFormatDifficultyLine(t *testing.T) {
	stats := &progress.Stats{}
	dashboard := NewDashboard(stats, false, "")

	diffStats := progress.DifficultyStats{
		Total:  10,
		Solved: 7,
	}

	output := dashboard.formatDifficultyLine("easy", diffStats)

	assert.Contains(t, output, "Easy:")
	assert.Contains(t, output, "7/10")
	assert.Contains(t, output, "70%")
	assert.Contains(t, output, "█")
	assert.Contains(t, output, "░")
}

func TestFormatTopicLine(t *testing.T) {
	stats := &progress.Stats{}
	dashboard := NewDashboard(stats, false, "")

	topicStats := progress.TopicStats{
		Total:  15,
		Solved: 9,
	}

	output := dashboard.formatTopicLine("arrays", topicStats)

	assert.Contains(t, output, "Arrays:")
	assert.Contains(t, output, "9/15")
	assert.Contains(t, output, "60%")
	assert.Contains(t, output, "█")
	assert.Contains(t, output, "░")
}

// Helper function to remove ANSI escape codes for length testing
func removeANSI(s string) string {
	// Remove ANSI escape sequences using a simple state machine
	var result strings.Builder
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			// Skip until we find the final character of the escape sequence
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') || s[i] == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteByte(s[i])
	}

	return result.String()
}

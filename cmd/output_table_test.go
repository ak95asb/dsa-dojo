package cmd

import (
	"os"
	"testing"
)

// TestGetTerminalWidth tests terminal width detection
func TestGetTerminalWidth(t *testing.T) {
	tests := []struct {
		name          string
		columnsEnv    string
		expectedWidth int
		description   string
	}{
		{
			name:          "COLUMNS environment variable set",
			columnsEnv:    "100",
			expectedWidth: 100,
			description:   "Should use COLUMNS env var when set",
		},
		{
			name:          "COLUMNS with invalid value",
			columnsEnv:    "invalid",
			expectedWidth: 80,
			description:   "Should fallback to default 80 when COLUMNS is invalid",
		},
		{
			name:          "COLUMNS with zero",
			columnsEnv:    "0",
			expectedWidth: 80,
			description:   "Should fallback to default 80 when COLUMNS is 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.columnsEnv != "" {
				os.Setenv("COLUMNS", tt.columnsEnv)
				defer os.Unsetenv("COLUMNS")
			}

			width := getTerminalWidth()

			// For valid COLUMNS, check exact match
			// For invalid, we expect default (80) or actual terminal width
			if tt.columnsEnv != "" && tt.columnsEnv != "invalid" && tt.columnsEnv != "0" {
				if width != tt.expectedWidth {
					t.Errorf("%s: expected %d, got %d", tt.description, tt.expectedWidth, width)
				}
			} else {
				// For invalid cases, just ensure we get a reasonable width (not 0, not negative)
				if width <= 0 {
					t.Errorf("%s: expected positive width, got %d", tt.description, width)
				}
			}
		})
	}
}

// TestSelectColumns tests column selection based on terminal width
func TestSelectColumns(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		expected []string
	}{
		{
			name:     "Wide terminal (>= 120)",
			width:    130,
			expected: []string{"ID", "Title", "Difficulty", "Topic", "Status"},
		},
		{
			name:     "Exact 120 width",
			width:    120,
			expected: []string{"ID", "Title", "Difficulty", "Topic", "Status"},
		},
		{
			name:     "Medium terminal (>= 80, < 120)",
			width:    100,
			expected: []string{"Title", "Difficulty", "Topic", "Status"},
		},
		{
			name:     "Exact 80 width",
			width:    80,
			expected: []string{"Title", "Difficulty", "Topic", "Status"},
		},
		{
			name:     "Narrow terminal (< 80)",
			width:    60,
			expected: []string{"Title", "Difficulty", "Status"},
		},
		{
			name:     "Very narrow terminal",
			width:    40,
			expected: []string{"Title", "Difficulty", "Status"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selectColumns(tt.width)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d columns, got %d", len(tt.expected), len(result))
				return
			}

			for i, col := range result {
				if col != tt.expected[i] {
					t.Errorf("Column %d: expected %s, got %s", i, tt.expected[i], col)
				}
			}
		})
	}
}

// TestTruncateText tests text truncation with ellipsis
func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "Text shorter than max length",
			text:     "Hello",
			maxLen:   10,
			expected: "Hello",
		},
		{
			name:     "Text equal to max length",
			text:     "Hello",
			maxLen:   5,
			expected: "Hello",
		},
		{
			name:     "Text longer than max length",
			text:     "Hello World",
			maxLen:   8,
			expected: "Hello...",
		},
		{
			name:     "Max length is 3",
			text:     "Hello",
			maxLen:   3,
			expected: "...",
		},
		{
			name:     "Max length is 2",
			text:     "Hello",
			maxLen:   2,
			expected: "...",
		},
		{
			name:     "Max length is 1",
			text:     "Hello",
			maxLen:   1,
			expected: "...",
		},
		{
			name:     "Empty string",
			text:     "",
			maxLen:   5,
			expected: "",
		},
		{
			name:     "Long title truncation",
			text:     "This is a very long problem title that needs truncation",
			maxLen:   20,
			expected: "This is a very lo...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateText(tt.text, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestCalculateTitleWidth tests title width calculation
func TestCalculateTitleWidth(t *testing.T) {
	tests := []struct {
		name        string
		termWidth   int
		numCols     int
		expectedMin int
		expectedMax int
	}{
		{
			name:        "Narrow terminal",
			termWidth:   80,
			numCols:     4,
			expectedMin: 15,
			expectedMax: 40,
		},
		{
			name:        "Wide terminal",
			termWidth:   150,
			numCols:     5,
			expectedMin: 15,
			expectedMax: 40,
		},
		{
			name:        "Very narrow terminal",
			termWidth:   60,
			numCols:     3,
			expectedMin: 15,
			expectedMax: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := calculateTitleWidth(tt.termWidth, tt.numCols)

			if width < tt.expectedMin {
				t.Errorf("Title width %d is below minimum %d", width, tt.expectedMin)
			}
			if width > tt.expectedMax {
				t.Errorf("Title width %d exceeds maximum %d", width, tt.expectedMax)
			}
		})
	}
}

// TestCalculateTopicWidth tests topic width calculation
func TestCalculateTopicWidth(t *testing.T) {
	tests := []struct {
		name      string
		termWidth int
		numCols   int
		expected  int
	}{
		{
			name:      "Wide terminal (>= 120)",
			termWidth: 130,
			numCols:   5,
			expected:  20,
		},
		{
			name:      "Medium terminal (>= 80, < 120)",
			termWidth: 100,
			numCols:   4,
			expected:  15,
		},
		{
			name:      "Narrow terminal (< 80)",
			termWidth: 70,
			numCols:   3,
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := calculateTopicWidth(tt.termWidth, tt.numCols)
			if width != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, width)
			}
		})
	}
}

// TestBuildTableRow tests table row building
func TestBuildTableRow(t *testing.T) {
	problem := ProblemRow{
		ID:         "two-sum",
		Title:      "Two Sum Problem",
		Difficulty: "easy",
		Topic:      "arrays",
		Solved:     true,
	}

	tests := []struct {
		name           string
		columns        []string
		titleMaxWidth  int
		topicMaxWidth  int
		expectedLength int
	}{
		{
			name:           "All columns",
			columns:        []string{"ID", "Title", "Difficulty", "Topic", "Status"},
			titleMaxWidth:  30,
			topicMaxWidth:  15,
			expectedLength: 5,
		},
		{
			name:           "Without ID",
			columns:        []string{"Title", "Difficulty", "Topic", "Status"},
			titleMaxWidth:  30,
			topicMaxWidth:  15,
			expectedLength: 4,
		},
		{
			name:           "Minimal columns",
			columns:        []string{"Title", "Difficulty", "Status"},
			titleMaxWidth:  20,
			topicMaxWidth:  0,
			expectedLength: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := buildTableRow(problem, tt.columns, tt.titleMaxWidth, tt.topicMaxWidth)

			if len(row) != tt.expectedLength {
				t.Errorf("Expected row length %d, got %d", tt.expectedLength, len(row))
			}

			// Verify row contains expected values
			for i, col := range tt.columns {
				switch col {
				case "ID":
					if row[i] != problem.ID {
						t.Errorf("Expected ID %s, got %s", problem.ID, row[i])
					}
				case "Title":
					// Title might be truncated
					if len(row[i]) == 0 {
						t.Error("Title should not be empty")
					}
				case "Difficulty":
					// Difficulty might have color codes
					// Just check it's not empty
					if len(row[i]) == 0 {
						t.Error("Difficulty should not be empty")
					}
				case "Topic":
					// Topic might be truncated
					if len(row[i]) == 0 {
						t.Error("Topic should not be empty")
					}
				case "Status":
					// Status will have color codes and symbols
					if len(row[i]) == 0 {
						t.Error("Status should not be empty")
					}
				}
			}
		})
	}
}

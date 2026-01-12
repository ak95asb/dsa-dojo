package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// TestShouldUseColors tests color detection logic
func TestShouldUseColors(t *testing.T) {
	tests := []struct {
		name        string
		noColorEnv  string
		expected    bool
		description string
	}{
		{
			name:        "NO_COLOR not set",
			noColorEnv:  "",
			expected:    false, // Will be false in test environment (not a terminal)
			description: "Should check terminal when NO_COLOR is not set",
		},
		{
			name:        "NO_COLOR set to 1",
			noColorEnv:  "1",
			expected:    false,
			description: "Should disable colors when NO_COLOR is set",
		},
		{
			name:        "NO_COLOR set to any value",
			noColorEnv:  "true",
			expected:    false,
			description: "Should disable colors when NO_COLOR has any value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set NO_COLOR environment variable
			if tt.noColorEnv != "" {
				os.Setenv("NO_COLOR", tt.noColorEnv)
				defer os.Unsetenv("NO_COLOR")
			} else {
				os.Unsetenv("NO_COLOR")
			}

			result := shouldUseColors()

			// In test environment, stdout is not a terminal
			// So we expect false regardless when NO_COLOR is not set
			if tt.noColorEnv != "" {
				if result != tt.expected {
					t.Errorf("%s: expected %v, got %v", tt.description, tt.expected, result)
				}
			} else {
				// When NO_COLOR is not set, result depends on terminal detection
				// In tests, it should be false (not a terminal)
				if result != false {
					t.Errorf("Expected false when running in test (not a terminal), got %v", result)
				}
			}
		})
	}
}

// TestColorDifficulty tests difficulty color application
func TestColorDifficulty(t *testing.T) {
	// Disable colors for predictable testing
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	tests := []struct {
		name       string
		difficulty string
		expected   string
	}{
		{
			name:       "Easy difficulty",
			difficulty: "easy",
			expected:   "easy",
		},
		{
			name:       "Medium difficulty",
			difficulty: "medium",
			expected:   "medium",
		},
		{
			name:       "Hard difficulty",
			difficulty: "hard",
			expected:   "hard",
		},
		{
			name:       "Easy with capital",
			difficulty: "Easy",
			expected:   "Easy",
		},
		{
			name:       "Unknown difficulty",
			difficulty: "unknown",
			expected:   "unknown",
		},
		{
			name:       "Empty difficulty",
			difficulty: "",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := colorDifficulty(tt.difficulty)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestColorDifficultyWithColors tests difficulty coloring with colors enabled
func TestColorDifficultyWithColors(t *testing.T) {
	// Enable colors for this test
	os.Unsetenv("NO_COLOR")
	color.NoColor = false

	tests := []struct {
		name         string
		difficulty   string
		shouldContain string
	}{
		{
			name:         "Easy should contain easy",
			difficulty:   "easy",
			shouldContain: "easy",
		},
		{
			name:         "Medium should contain medium",
			difficulty:   "medium",
			shouldContain: "medium",
		},
		{
			name:         "Hard should contain hard",
			difficulty:   "hard",
			shouldContain: "hard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := colorDifficulty(tt.difficulty)

			// Result should contain the difficulty text (may have ANSI codes)
			if !strings.Contains(result, tt.shouldContain) {
				t.Errorf("Expected result to contain '%s', got '%s'", tt.shouldContain, result)
			}
		})
	}
}

// TestColorStatus tests status color application
func TestColorStatus(t *testing.T) {
	// Disable colors for predictable testing
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	tests := []struct {
		name     string
		solved   bool
		expected string
	}{
		{
			name:     "Solved problem",
			solved:   true,
			expected: "✓",
		},
		{
			name:     "Unsolved problem",
			solved:   false,
			expected: "✗",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := colorStatus(tt.solved)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestColorize tests the colorize helper function
func TestColorize(t *testing.T) {
	// Test with NO_COLOR set
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	text := "Hello"
	c := ColorGreen

	result := colorize(text, c)
	if result != text {
		t.Errorf("Expected '%s' (no color), got '%s'", text, result)
	}
}

// TestColorizeWithColors tests colorize with colors enabled
func TestColorizeWithColors(t *testing.T) {
	// Enable colors
	os.Unsetenv("NO_COLOR")
	color.NoColor = false

	text := "Hello"
	c := ColorGreen

	result := colorize(text, c)

	// Result should contain the text (may have ANSI codes when colors are supported)
	if !strings.Contains(result, text) {
		t.Errorf("Expected result to contain '%s', got '%s'", text, result)
	}
}

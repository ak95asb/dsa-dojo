package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestIntegration_ListTableOutput tests list command table output
func TestIntegration_ListTableOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary for testing
	buildCmd := exec.Command("go", "build", "-o", "dsa_test")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("../dsa_test")

	tests := []struct {
		name              string
		args              []string
		envVars           map[string]string
		expectedInOutput  []string
		unexpectedInOutput []string
	}{
		{
			name: "List all problems with default table format",
			args: []string{"list"},
			expectedInOutput: []string{
				"TITLE",
				"DIFFICULTY",
				"STATUS",
				"Total:",
			},
		},
		{
			name: "List with wide terminal (all columns)",
			args: []string{"list"},
			envVars: map[string]string{
				"COLUMNS": "130",
			},
			expectedInOutput: []string{
				"ID",
				"TITLE",
				"DIFFICULTY",
				"TOPIC",
				"STATUS",
			},
		},
		{
			name: "List with medium terminal (no ID column)",
			args: []string{"list"},
			envVars: map[string]string{
				"COLUMNS": "100",
			},
			expectedInOutput: []string{
				"TITLE",
				"DIFFICULTY",
				"TOPIC",
				"STATUS",
			},
			unexpectedInOutput: []string{
				"│ ID",
			},
		},
		{
			name: "List with narrow terminal (minimal columns)",
			args: []string{"list"},
			envVars: map[string]string{
				"COLUMNS": "70",
			},
			expectedInOutput: []string{
				"TITLE",
				"DIFFICULTY",
				"STATUS",
			},
			unexpectedInOutput: []string{
				"│ ID",
				"│ TOPIC",
			},
		},
		{
			name: "List with NO_COLOR environment variable",
			args: []string{"list", "--difficulty", "easy"},
			envVars: map[string]string{
				"NO_COLOR": "1",
			},
			expectedInOutput: []string{
				"easy",
				"TITLE",
			},
		},
		{
			name: "List with difficulty filter",
			args: []string{"list", "--difficulty", "hard"},
			expectedInOutput: []string{
				"DIFFICULTY",
				"hard",
			},
		},
		{
			name: "List with format json",
			args: []string{"list", "--format", "json"},
			expectedInOutput: []string{
				`"problems"`,
				`"total"`,
				`"solved"`,
			},
			unexpectedInOutput: []string{
				"TITLE",
				"┌",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../dsa_test", tt.args...)

			// Set environment variables
			cmd.Env = os.Environ()
			for key, value := range tt.envVars {
				cmd.Env = append(cmd.Env, key+"="+value)
			}

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("Command output: %s", string(output))
				t.Fatalf("Command failed: %v", err)
			}

			outputStr := string(output)

			// Check expected strings are in output
			for _, expected := range tt.expectedInOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain '%s', but it didn't.\nOutput:\n%s",
						expected, outputStr)
				}
			}

			// Check unexpected strings are NOT in output
			for _, unexpected := range tt.unexpectedInOutput {
				if strings.Contains(outputStr, unexpected) {
					t.Errorf("Expected output to NOT contain '%s', but it did.\nOutput:\n%s",
						unexpected, outputStr)
				}
			}
		})
	}
}

// TestIntegration_StatusTableOutput tests status command table output
func TestIntegration_StatusTableOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary for testing
	buildCmd := exec.Command("go", "build", "-o", "dsa_test")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("../dsa_test")

	tests := []struct {
		name              string
		args              []string
		envVars           map[string]string
		expectedInOutput  []string
		unexpectedInOutput []string
	}{
		{
			name: "Status with default table format",
			args: []string{"status"},
			expectedInOutput: []string{
				"DSA Progress Dashboard",
				"Overall Progress:",
				"Progress by Difficulty:",
				"Progress by Topic:",
				"DIFFICULTY",
				"TOPIC",
			},
		},
		{
			name: "Status with compact flag",
			args: []string{"status", "--compact"},
			expectedInOutput: []string{
				"Progress:",
				"Easy:",
				"Med:",
				"Hard:",
			},
			unexpectedInOutput: []string{
				"┌",
				"│",
				"└",
			},
		},
		{
			name: "Status with format json",
			args: []string{"status", "--format", "json"},
			expectedInOutput: []string{
				`"total_problems"`,
				`"problems_solved"`,
				`"by_difficulty"`,
				`"by_topic"`,
			},
			unexpectedInOutput: []string{
				"Dashboard",
				"┌",
			},
		},
		{
			name: "Status with NO_COLOR",
			args: []string{"status"},
			envVars: map[string]string{
				"NO_COLOR": "1",
			},
			expectedInOutput: []string{
				"Progress by Difficulty:",
				"Progress by Topic:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../dsa_test", tt.args...)

			// Set environment variables
			cmd.Env = os.Environ()
			for key, value := range tt.envVars {
				cmd.Env = append(cmd.Env, key+"="+value)
			}

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("Command output: %s", string(output))
				t.Fatalf("Command failed: %v", err)
			}

			outputStr := string(output)

			// Check expected strings are in output
			for _, expected := range tt.expectedInOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain '%s', but it didn't.\nOutput:\n%s",
						expected, outputStr)
				}
			}

			// Check unexpected strings are NOT in output
			for _, unexpected := range tt.unexpectedInOutput {
				if strings.Contains(outputStr, unexpected) {
					t.Errorf("Expected output to NOT contain '%s', but it did.\nOutput:\n%s",
						unexpected, outputStr)
				}
			}
		})
	}
}

// TestIntegration_TableBorders tests that table borders are properly rendered
func TestIntegration_TableBorders(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary for testing
	buildCmd := exec.Command("go", "build", "-o", "dsa_test")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("../dsa_test")

	cmd := exec.Command("../dsa_test", "list", "--difficulty", "easy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command output: %s", string(output))
		t.Fatalf("Command failed: %v", err)
	}

	outputStr := string(output)

	// Check for table border characters
	borderChars := []string{"┌", "─", "┐", "│", "├", "┼", "┤", "└", "┴", "┘"}
	for _, char := range borderChars {
		if !strings.Contains(outputStr, char) {
			t.Errorf("Expected table to contain border character '%s'", char)
		}
	}
}

// TestIntegration_TableOutputConsistency tests that table output is consistent
func TestIntegration_TableOutputConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary for testing
	buildCmd := exec.Command("go", "build", "-o", "dsa_test")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("../dsa_test")

	cmd := exec.Command("../dsa_test", "list", "--difficulty", "easy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command output: %s", string(output))
		t.Fatalf("Command failed: %v", err)
	}

	outputStr := string(output)

	// Verify table output has required elements
	requiredElements := []string{
		"TITLE",
		"DIFFICULTY",
		"STATUS",
		"Total:",
		"problems",
	}

	for _, element := range requiredElements {
		if !strings.Contains(outputStr, element) {
			t.Errorf("Expected output to contain '%s'", element)
		}
	}

	// Verify output has table borders
	if !strings.Contains(outputStr, "┌") || !strings.Contains(outputStr, "│") {
		t.Error("Expected output to contain table border characters")
	}

	// Verify output has summary line
	if !strings.Contains(outputStr, "solved") && !strings.Contains(outputStr, "unsolved") {
		t.Error("Expected output to contain summary statistics")
	}
}

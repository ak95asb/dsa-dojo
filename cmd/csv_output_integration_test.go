package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestIntegration_ListCSVOutput tests list command CSV output
func TestIntegration_ListCSVOutput(t *testing.T) {
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
		expectedInOutput  []string
		unexpectedInOutput []string
	}{
		{
			name: "List CSV with headers",
			args: []string{"list", "--format", "csv", "--difficulty", "easy"},
			expectedInOutput: []string{
				"ID,Title,Difficulty,Topic,Solved,FirstSolvedAt",
				"easy",
				"false",
			},
		},
		{
			name: "List CSV with all problems",
			args: []string{"list", "--format", "csv"},
			expectedInOutput: []string{
				"ID,Title,Difficulty,Topic,Solved,FirstSolvedAt",
				"easy",
				"medium",
				"hard",
			},
		},
		{
			name: "List CSV is not table format",
			args: []string{"list", "--format", "csv"},
			unexpectedInOutput: []string{
				"┌",
				"│",
				"└",
				"Total:",
			},
		},
		{
			name: "List CSV is not JSON format",
			args: []string{"list", "--format", "csv"},
			unexpectedInOutput: []string{
				`"problems"`,
				`"total"`,
				"{",
				"}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../dsa_test", tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("Command output: %s", string(output))
				t.Fatalf("Command failed: %v", err)
			}

			outputStr := string(output)

			// Check expected strings
			for _, expected := range tt.expectedInOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain '%s', but it didn't.\nOutput:\n%s",
						expected, outputStr)
				}
			}

			// Check unexpected strings
			for _, unexpected := range tt.unexpectedInOutput {
				if strings.Contains(outputStr, unexpected) {
					t.Errorf("Expected output to NOT contain '%s', but it did.\nOutput:\n%s",
						unexpected, outputStr)
				}
			}
		})
	}
}

// TestIntegration_StatusCSVOutput tests status command CSV output
func TestIntegration_StatusCSVOutput(t *testing.T) {
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
		expectedInOutput  []string
		unexpectedInOutput []string
	}{
		{
			name: "Status CSV with headers",
			args: []string{"status", "--format", "csv"},
			expectedInOutput: []string{
				"Category,Value,Total,Solved,Unsolved",
				"Overall,Problems",
				"Difficulty,easy",
				"Difficulty,medium",
				"Difficulty,hard",
				"Topic,",
			},
		},
		{
			name: "Status CSV is not table format",
			args: []string{"status", "--format", "csv"},
			unexpectedInOutput: []string{
				"DSA Progress Dashboard",
				"┌",
				"│",
				"└",
			},
		},
		{
			name: "Status CSV is not JSON format",
			args: []string{"status", "--format", "csv"},
			unexpectedInOutput: []string{
				`"total_problems"`,
				`"problems_solved"`,
				"{",
				"}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("../dsa_test", tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("Command output: %s", string(output))
				t.Fatalf("Command failed: %v", err)
			}

			outputStr := string(output)

			// Check expected strings
			for _, expected := range tt.expectedInOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Expected output to contain '%s', but it didn't.\nOutput:\n%s",
						expected, outputStr)
				}
			}

			// Check unexpected strings
			for _, unexpected := range tt.unexpectedInOutput {
				if strings.Contains(outputStr, unexpected) {
					t.Errorf("Expected output to NOT contain '%s', but it did.\nOutput:\n%s",
						unexpected, outputStr)
				}
			}
		})
	}
}

// TestIntegration_CSVFileRedirect tests CSV output can be redirected to file
func TestIntegration_CSVFileRedirect(t *testing.T) {
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

	// Create temp file for CSV output
	tmpFile := "/tmp/test_problems.csv"
	defer os.Remove(tmpFile)

	// Run command with redirect
	cmd := exec.Command("sh", "-c", "../dsa_test list --format csv --difficulty easy > "+tmpFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	// Read the file
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	// Verify CSV structure
	if !strings.Contains(contentStr, "ID,Title,Difficulty,Topic,Solved,FirstSolvedAt") {
		t.Error("CSV file missing header row")
	}

	// Verify it's pure CSV (no extra formatting)
	lines := strings.Split(strings.TrimSpace(contentStr), "\n")
	if len(lines) < 2 {
		t.Error("CSV file should have at least header + 1 data row")
	}

	// Each line should be comma-separated values
	for i, line := range lines {
		if !strings.Contains(line, ",") {
			t.Errorf("Line %d is not comma-separated: %s", i, line)
		}
	}
}

// TestIntegration_CSVDateFormat tests date formatting in CSV
func TestIntegration_CSVDateFormat(t *testing.T) {
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

	cmd := exec.Command("../dsa_test", "list", "--format", "csv")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command output: %s", string(output))
		t.Fatalf("Command failed: %v", err)
	}

	outputStr := string(output)

	// Verify header has FirstSolvedAt column
	if !strings.Contains(outputStr, "FirstSolvedAt") {
		t.Error("CSV should have FirstSolvedAt column in header")
	}

	// Split into lines and check each row
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")
	for i, line := range lines[1:] { // Skip header
		fields := strings.Split(line, ",")
		if len(fields) < 6 {
			t.Errorf("Row %d has fewer than 6 fields: %s", i, line)
			continue
		}

		// LastSolved field is at index 5
		dateField := fields[5]

		// Date should either be empty or in ISO 8601 format (YYYY-MM-DD)
		if dateField != "" && !isISO8601Date(dateField) {
			t.Errorf("Row %d has invalid date format '%s': %s", i, dateField, line)
		}
	}
}

// TestIntegration_CSVEmptyFirstSolvedAt tests unsolved problems have empty FirstSolvedAt
func TestIntegration_CSVEmptyFirstSolvedAt(t *testing.T) {
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

	cmd := exec.Command("../dsa_test", "list", "--format", "csv", "--difficulty", "easy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command output: %s", string(output))
		t.Fatalf("Command failed: %v", err)
	}

	outputStr := string(output)
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")

	// Check for rows with false (unsolved) that should have empty FirstSolvedAt
	for i, line := range lines[1:] { // Skip header
		if strings.Contains(line, ",false,") {
			// This is an unsolved problem
			fields := strings.Split(line, ",")
			if len(fields) >= 6 {
				dateField := fields[5]
				if dateField != "" {
					t.Errorf("Row %d: unsolved problem should have empty FirstSolvedAt, got '%s': %s",
						i, dateField, line)
				}
			}
		}
	}
}

// Helper to check if string is valid ISO 8601 date (YYYY-MM-DD)
func isISO8601Date(s string) bool {
	// Simple validation: 10 characters, format YYYY-MM-DD
	if len(s) != 10 {
		return false
	}
	if s[4] != '-' || s[7] != '-' {
		return false
	}
	return true
}

package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// TestFormatDateISO tests ISO 8601 date formatting
func TestFormatDateISO(t *testing.T) {
	tests := []struct {
		name     string
		date     *time.Time
		expected string
	}{
		{
			name:     "valid date",
			date:     timePtr(time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)),
			expected: "2025-01-15",
		},
		{
			name:     "nil date",
			date:     nil,
			expected: "",
		},
		{
			name:     "zero time",
			date:     &time.Time{},
			expected: "",
		},
		{
			name:     "different date",
			date:     timePtr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
			expected: "2024-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDateISO(tt.date)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestFormatBoolCSV tests boolean formatting for CSV
func TestFormatBoolCSV(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected string
	}{
		{
			name:     "true value",
			value:    true,
			expected: "true",
		},
		{
			name:     "false value",
			value:    false,
			expected: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBoolCSV(tt.value)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestWriteCSV_SimpleData tests basic CSV writing
func TestWriteCSV_SimpleData(t *testing.T) {
	headers := []string{"Name", "Age", "City"}
	rows := [][]string{
		{"Alice", "30", "New York"},
		{"Bob", "25", "Los Angeles"},
	}

	output := captureStdout(t, func() {
		err := writeCSV(headers, rows)
		if err != nil {
			t.Fatalf("writeCSV failed: %v", err)
		}
	})

	// Verify headers
	if !strings.Contains(output, "Name,Age,City") {
		t.Errorf("Expected headers 'Name,Age,City', got: %s", output)
	}

	// Verify data rows
	if !strings.Contains(output, "Alice,30,New York") {
		t.Errorf("Expected row 'Alice,30,New York', got: %s", output)
	}
	if !strings.Contains(output, "Bob,25,Los Angeles") {
		t.Errorf("Expected row 'Bob,25,Los Angeles', got: %s", output)
	}
}

// TestWriteCSV_QuotingCommas tests CSV quoting for values with commas
func TestWriteCSV_QuotingCommas(t *testing.T) {
	headers := []string{"Title", "Description"}
	rows := [][]string{
		{"Problem 1, Part A", "Easy problem"},
		{"Problem 2", "Medium problem, with comma"},
	}

	output := captureStdout(t, func() {
		err := writeCSV(headers, rows)
		if err != nil {
			t.Fatalf("writeCSV failed: %v", err)
		}
	})

	// Values with commas should be quoted
	if !strings.Contains(output, `"Problem 1, Part A"`) {
		t.Errorf("Expected quoted value '\"Problem 1, Part A\"', got: %s", output)
	}
	if !strings.Contains(output, `"Medium problem, with comma"`) {
		t.Errorf("Expected quoted value '\"Medium problem, with comma\"', got: %s", output)
	}
}

// TestWriteCSV_QuotingQuotes tests CSV escaping for values with quotes
func TestWriteCSV_QuotingQuotes(t *testing.T) {
	headers := []string{"Title"}
	rows := [][]string{
		{`Title with "quotes"`},
	}

	output := captureStdout(t, func() {
		err := writeCSV(headers, rows)
		if err != nil {
			t.Fatalf("writeCSV failed: %v", err)
		}
	})

	// Quotes should be escaped by doubling them
	if !strings.Contains(output, `"Title with ""quotes"""`) {
		t.Errorf("Expected escaped quotes, got: %s", output)
	}
}

// TestWriteCSV_EmptyValues tests empty value handling
func TestWriteCSV_EmptyValues(t *testing.T) {
	headers := []string{"ID", "Title", "Date"}
	rows := [][]string{
		{"1", "Problem 1", "2025-01-15"},
		{"2", "Problem 2", ""}, // Empty date
	}

	output := captureStdout(t, func() {
		err := writeCSV(headers, rows)
		if err != nil {
			t.Fatalf("writeCSV failed: %v", err)
		}
	})

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 { // Header + 2 rows
		t.Fatalf("Expected 3 lines, got %d", len(lines))
	}

	// Second row should have empty date field
	if !strings.HasSuffix(lines[2], ",") {
		t.Errorf("Expected row to end with empty field (comma), got: %s", lines[2])
	}
}

// TestWriteCSV_NewlineInValue tests quoting for values with newlines
func TestWriteCSV_NewlineInValue(t *testing.T) {
	headers := []string{"Description"}
	rows := [][]string{
		{"Line 1\nLine 2"},
	}

	output := captureStdout(t, func() {
		err := writeCSV(headers, rows)
		if err != nil {
			t.Fatalf("writeCSV failed: %v", err)
		}
	})

	// Values with newlines should be quoted
	if !strings.Contains(output, `"Line 1`) {
		t.Errorf("Expected quoted value with newline, got: %s", output)
	}
}

// TestWriteCSV_MultiRow tests multiple row CSV output
func TestWriteCSV_MultiRow(t *testing.T) {
	headers := []string{"ID", "Value"}
	rows := [][]string{
		{"1", "A"},
		{"2", "B"},
		{"3", "C"},
	}

	output := captureStdout(t, func() {
		err := writeCSV(headers, rows)
		if err != nil {
			t.Fatalf("writeCSV failed: %v", err)
		}
	})

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 4 { // Header + 3 rows
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

// Helper function to capture stdout
func captureStdout(t *testing.T, fn func()) string {
	// Save original stdout
	oldStdout := os.Stdout

	// Create pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Redirect stdout to pipe
	os.Stdout = w

	// Run function
	fn()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

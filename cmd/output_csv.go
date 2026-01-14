package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

// writeCSV writes data as CSV to stdout
func writeCSV(headers []string, rows [][]string) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write header row
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// formatDateISO formats time pointer as ISO 8601 date (YYYY-MM-DD)
func formatDateISO(t *time.Time) string {
	if t == nil || t.IsZero() {
		return "" // Empty string for nil or zero time
	}
	return t.Format("2006-01-02")
}

// formatBoolCSV formats boolean for CSV (true/false)
func formatBoolCSV(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

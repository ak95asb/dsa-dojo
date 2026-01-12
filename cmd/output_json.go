package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// ListResponse represents JSON output for list command
//
// JSON Schema:
//
//	{
//	  "problems": [ /* array of ProblemJSON objects */ ],
//	  "total": 20,      // total number of problems
//	  "solved": 5       // number of solved problems
//	}
type ListResponse struct {
	Problems []ProblemJSON `json:"problems"`
	Total    int           `json:"total"`
	Solved   int           `json:"solved"`
}

// ProblemJSON represents a problem in JSON format
//
// JSON Schema:
//
//	{
//	  "id": "two-sum",
//	  "title": "Two Sum",
//	  "difficulty": "easy",  // easy, medium, or hard
//	  "topic": "arrays",
//	  "solved": false
//	}
type ProblemJSON struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	Topic      string `json:"topic"`
	Solved     bool   `json:"solved"`
}

// StatusResponse represents JSON output for status command
//
// JSON Schema:
//
//	{
//	  "total_problems": 20,
//	  "problems_solved": 5,
//	  "by_difficulty": {
//	    "easy": 2,
//	    "medium": 2,
//	    "hard": 1
//	  },
//	  "by_topic": {
//	    "arrays": 3,
//	    "trees": 2
//	  },
//	  "recent_activity": [ /* optional, array of RecentActivityJSON */ ],
//	  "streak": 5       // optional, Phase 2 feature
//	}
type StatusResponse struct {
	TotalProblems  int                  `json:"total_problems"`
	ProblemsSolved int                  `json:"problems_solved"`
	ByDifficulty   map[string]int       `json:"by_difficulty"`
	ByTopic        map[string]int       `json:"by_topic"`
	RecentActivity []RecentActivityJSON `json:"recent_activity,omitempty"`
	Streak         int                  `json:"streak,omitempty"` // Phase 2 feature
}

// RecentActivityJSON represents recent problem activity
//
// JSON Schema:
//
//	{
//	  "problem_id": "two-sum",
//	  "title": "Two Sum",
//	  "date": "2025-01-15T14:30:00Z",  // RFC3339 format
//	  "passed": true
//	}
type RecentActivityJSON struct {
	ProblemID string `json:"problem_id"`
	Title     string `json:"title"`
	Date      string `json:"date"` // RFC3339 format
	Passed    bool   `json:"passed"`
}

// outputJSON marshals data to JSON and prints to stdout
func outputJSON(data interface{}, compact bool) error {
	var output []byte
	var err error

	if compact {
		// Compact JSON (single line, no indentation)
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(data); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
		// Remove trailing newline added by Encoder
		output = bytes.TrimSpace(buffer.Bytes())
	} else {
		// Formatted JSON (indented)
		output, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	}

	fmt.Fprintln(os.Stdout, string(output))
	return nil
}

// isValidFormat checks if format string is valid
func isValidFormat(format string) bool {
	validFormats := []string{"table", "json", "csv"}
	for _, valid := range validFormats {
		if format == valid {
			return true
		}
	}
	return false
}

package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListResponse_JSONSchema(t *testing.T) {
	t.Run("formatted JSON has indentation", func(t *testing.T) {
		problems := []ProblemJSON{
			{ID: "two-sum", Title: "Two Sum", Difficulty: "easy", Topic: "arrays", Solved: true},
		}
		response := ListResponse{
			Problems: problems,
			Total:    1,
			Solved:   1,
		}

		output, err := json.MarshalIndent(response, "", "  ")
		assert.NoError(t, err)
		assert.Contains(t, string(output), "\n")
		assert.Contains(t, string(output), "  \"problems\"")
	})

	t.Run("JSON uses snake_case field names", func(t *testing.T) {
		response := ListResponse{
			Problems: []ProblemJSON{},
			Total:    0,
			Solved:   0,
		}

		output, err := json.Marshal(response)
		assert.NoError(t, err)

		// Verify snake_case fields exist
		assert.Contains(t, string(output), "\"problems\"")
		assert.Contains(t, string(output), "\"total\"")
		assert.Contains(t, string(output), "\"solved\"")
	})

	t.Run("problem JSON uses snake_case", func(t *testing.T) {
		problem := ProblemJSON{
			ID:         "two-sum",
			Title:      "Two Sum",
			Difficulty: "easy",
			Topic:      "arrays",
			Solved:     true,
		}

		output, err := json.Marshal(problem)
		assert.NoError(t, err)

		// Verify all fields are snake_case
		assert.Contains(t, string(output), "\"id\"")
		assert.Contains(t, string(output), "\"title\"")
		assert.Contains(t, string(output), "\"difficulty\"")
		assert.Contains(t, string(output), "\"topic\"")
		assert.Contains(t, string(output), "\"solved\"")
	})
}

func TestStatusResponse_JSONSchema(t *testing.T) {
	t.Run("status JSON uses snake_case field names", func(t *testing.T) {
		response := StatusResponse{
			TotalProblems:  20,
			ProblemsSolved: 5,
			ByDifficulty: map[string]int{
				"easy":   2,
				"medium": 2,
				"hard":   1,
			},
			ByTopic: map[string]int{
				"arrays": 3,
			},
		}

		output, err := json.Marshal(response)
		assert.NoError(t, err)

		// Verify snake_case fields
		assert.Contains(t, string(output), "\"total_problems\"")
		assert.Contains(t, string(output), "\"problems_solved\"")
		assert.Contains(t, string(output), "\"by_difficulty\"")
		assert.Contains(t, string(output), "\"by_topic\"")
	})

	t.Run("recent activity uses RFC3339 date format", func(t *testing.T) {
		testTime := time.Date(2025, 1, 15, 14, 30, 0, 0, time.UTC)
		activity := RecentActivityJSON{
			ProblemID: "two-sum",
			Title:     "Two Sum",
			Date:      testTime.Format(time.RFC3339),
			Passed:    true,
		}

		output, err := json.Marshal(activity)
		assert.NoError(t, err)

		// Verify RFC3339 format
		assert.Contains(t, string(output), "2025-01-15T14:30:00Z")
		assert.Contains(t, string(output), "\"problem_id\"")
		assert.Contains(t, string(output), "\"passed\"")
	})
}

func TestOutputJSON(t *testing.T) {
	t.Run("formatted JSON output", func(t *testing.T) {
		response := ListResponse{
			Problems: []ProblemJSON{
				{ID: "test", Title: "Test", Difficulty: "easy", Topic: "arrays", Solved: false},
			},
			Total:  1,
			Solved: 0,
		}

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := outputJSON(response, false)
		assert.NoError(t, err)

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// Verify formatted JSON
		assert.Contains(t, output, "\n")
		assert.Contains(t, output, "  ")
		assert.Contains(t, output, "\"problems\"")
	})

	t.Run("compact JSON is single line", func(t *testing.T) {
		response := ListResponse{
			Problems: []ProblemJSON{},
			Total:    0,
			Solved:   0,
		}

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := outputJSON(response, true)
		assert.NoError(t, err)

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := strings.TrimSpace(buf.String())

		// Verify compact JSON (should not contain newlines within JSON itself)
		lines := strings.Split(output, "\n")
		assert.Len(t, lines, 1, "Compact JSON should be on a single line")
	})

	t.Run("empty arrays are not null", func(t *testing.T) {
		response := ListResponse{
			Problems: []ProblemJSON{},
			Total:    0,
			Solved:   0,
		}

		output, err := json.Marshal(response)
		assert.NoError(t, err)

		// Verify empty array is [] not null
		assert.Contains(t, string(output), "\"problems\":[]")
		assert.NotContains(t, string(output), "\"problems\":null")
	})
}

func TestIsValidFormat(t *testing.T) {
	t.Run("table is valid", func(t *testing.T) {
		assert.True(t, isValidFormat("table"))
	})

	t.Run("json is valid", func(t *testing.T) {
		assert.True(t, isValidFormat("json"))
	})

	t.Run("invalid format returns false", func(t *testing.T) {
		assert.False(t, isValidFormat("xml"))
		assert.False(t, isValidFormat("yaml"))
		assert.False(t, isValidFormat(""))
	})
}

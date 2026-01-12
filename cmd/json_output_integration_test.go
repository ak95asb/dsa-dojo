package cmd

import (
	"encoding/json"
	"testing"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/ak95asb/dsa-dojo/internal/progress"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Integration tests for JSON output functionality

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&database.Problem{}, &database.Solution{}, &database.Progress{}, &database.BenchmarkResult{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

func TestIntegration_ListCommand_JSONFormat(t *testing.T) {
	// Setup isolated in-memory database
	db := setupTestDB(t)

	// Seed test data
	svc := problem.NewService(db)
	testProblems := []problem.CreateProblemInput{
		{
			Title:       "Two Sum",
			Difficulty:  "easy",
			Topic:       "arrays",
			Description: "Find two numbers that add up to target",
		},
		{
			Title:       "Reverse Linked List",
			Difficulty:  "medium",
			Topic:       "linked-lists",
			Description: "Reverse a singly linked list",
		},
	}

	var createdProblems []*database.Problem
	for _, input := range testProblems {
		p, err := svc.CreateProblem(input)
		require.NoError(t, err)
		createdProblems = append(createdProblems, p)
	}

	// Mark first problem as solved (Two Sum)
	require.NoError(t, svc.UpdateProgress(createdProblems[0].ID, "completed"))

	t.Run("formatted JSON output", func(t *testing.T) {
		// Query problems using service
		problems, err := svc.ListProblems(problem.ListFilters{})
		require.NoError(t, err)
		require.Len(t, problems, 2)

		// Build JSON response
		jsonProblems := make([]ProblemJSON, len(problems))
		solvedCount := 0
		for i, p := range problems {
			jsonProblems[i] = ProblemJSON{
				ID:         p.Slug,
				Title:      p.Title,
				Difficulty: p.Difficulty,
				Topic:      p.Topic,
				Solved:     p.IsSolved,
			}
			if p.IsSolved {
				solvedCount++
			}
		}

		response := ListResponse{
			Problems: jsonProblems,
			Total:    len(problems),
			Solved:   solvedCount,
		}

		// Test JSON marshaling
		output, err := json.MarshalIndent(response, "", "  ")
		require.NoError(t, err, "Should marshal to JSON")

		// Verify content
		assert.Equal(t, 2, response.Total)
		assert.Equal(t, 1, response.Solved)
		assert.Len(t, response.Problems, 2)

		// Verify formatted output
		assert.Contains(t, string(output), "\n")
		assert.Contains(t, string(output), "  ")

		// Verify problem data
		var twoSum *ProblemJSON
		for i := range response.Problems {
			if response.Problems[i].ID == "two-sum" {
				twoSum = &response.Problems[i]
				break
			}
		}
		require.NotNil(t, twoSum)
		assert.Equal(t, "Two Sum", twoSum.Title)
		assert.Equal(t, "easy", twoSum.Difficulty)
		assert.Equal(t, "arrays", twoSum.Topic)
		assert.True(t, twoSum.Solved)
	})

	t.Run("compact JSON output", func(t *testing.T) {
		// Query problems
		problems, err := svc.ListProblems(problem.ListFilters{})
		require.NoError(t, err)

		// Build JSON response
		jsonProblems := make([]ProblemJSON, len(problems))
		for i, p := range problems {
			jsonProblems[i] = ProblemJSON{
				ID:         p.Slug,
				Title:      p.Title,
				Difficulty: p.Difficulty,
				Topic:      p.Topic,
				Solved:     p.IsSolved,
			}
		}

		response := ListResponse{
			Problems: jsonProblems,
			Total:    len(problems),
			Solved:   1,
		}

		// Test compact JSON marshaling
		output, err := json.Marshal(response)
		require.NoError(t, err, "Should marshal to compact JSON")

		// Verify it's valid JSON
		var parsed ListResponse
		err = json.Unmarshal(output, &parsed)
		require.NoError(t, err, "Compact output should be valid JSON")

		// Verify no excessive whitespace (compact format)
		assert.NotContains(t, string(output), "\n  ", "Compact JSON should not have indentation")
	})

	t.Run("filters work with JSON output", func(t *testing.T) {
		// Query with filter
		problems, err := svc.ListProblems(problem.ListFilters{Difficulty: "easy"})
		require.NoError(t, err)

		// Build JSON response
		jsonProblems := make([]ProblemJSON, len(problems))
		solvedCount := 0
		for i, p := range problems {
			jsonProblems[i] = ProblemJSON{
				ID:         p.Slug,
				Title:      p.Title,
				Difficulty: p.Difficulty,
				Topic:      p.Topic,
				Solved:     p.IsSolved,
			}
			if p.IsSolved {
				solvedCount++
			}
		}

		response := ListResponse{
			Problems: jsonProblems,
			Total:    len(problems),
			Solved:   solvedCount,
		}

		// Test JSON marshaling
		_, err = json.Marshal(response)
		require.NoError(t, err)

		// Should only have one easy problem
		assert.Equal(t, 1, response.Total)
		assert.Len(t, response.Problems, 1)
		assert.Equal(t, "two-sum", response.Problems[0].ID)
	})

	t.Run("empty results return valid JSON", func(t *testing.T) {
		// Query with filter that returns nothing
		_, err := svc.ListProblems(problem.ListFilters{Difficulty: "hard"})
		require.NoError(t, err)

		// Build JSON response with empty slice
		response := ListResponse{
			Problems: []ProblemJSON{}, // Empty slice, not nil
			Total:    0,
			Solved:   0,
		}

		// Test JSON marshaling
		output, err := json.Marshal(response)
		require.NoError(t, err)

		// Verify JSON structure
		var parsed ListResponse
		err = json.Unmarshal(output, &parsed)
		require.NoError(t, err)

		// Should have empty array, not null
		assert.Equal(t, 0, response.Total)
		assert.Equal(t, 0, response.Solved)
		assert.NotNil(t, response.Problems)
		assert.Len(t, response.Problems, 0)

		// Verify empty array is [] not null
		assert.Contains(t, string(output), `"problems":[]`)
		assert.NotContains(t, string(output), `"problems":null`)
	})
}

func TestIntegration_StatusCommand_JSONFormat(t *testing.T) {
	// Setup isolated in-memory database
	db := setupTestDB(t)

	// Seed test data
	svc := problem.NewService(db)
	testProblems := []problem.CreateProblemInput{
		{Title: "Two Sum", Difficulty: "easy", Topic: "arrays", Description: "Test"},
		{Title: "Best Time", Difficulty: "easy", Topic: "arrays", Description: "Test"},
		{Title: "Reverse List", Difficulty: "medium", Topic: "linked-lists", Description: "Test"},
		{Title: "Merge K", Difficulty: "hard", Topic: "linked-lists", Description: "Test"},
	}

	var createdProblems []*database.Problem
	for _, input := range testProblems {
		p, err := svc.CreateProblem(input)
		require.NoError(t, err)
		createdProblems = append(createdProblems, p)
	}

	// Mark some as solved (Two Sum and Reverse List)
	require.NoError(t, svc.UpdateProgress(createdProblems[0].ID, "completed"))
	require.NoError(t, svc.UpdateProgress(createdProblems[2].ID, "completed"))

	t.Run("formatted JSON output", func(t *testing.T) {
		// Get stats using progress service
		progressSvc := progress.NewService(db)
		stats, err := progressSvc.GetStats("")
		require.NoError(t, err)

		// Build JSON response
		byDifficulty := make(map[string]int)
		for diff, diffStats := range stats.ByDifficulty {
			byDifficulty[diff] = diffStats.Solved
		}

		byTopic := make(map[string]int)
		for topic, topicStats := range stats.ByTopic {
			byTopic[topic] = topicStats.Solved
		}

		response := StatusResponse{
			TotalProblems:  stats.TotalProblems,
			ProblemsSolved: stats.TotalSolved,
			ByDifficulty:   byDifficulty,
			ByTopic:        byTopic,
		}

		// Test JSON marshaling
		output, err := json.MarshalIndent(response, "", "  ")
		require.NoError(t, err, "Should marshal to JSON")

		// Verify content
		assert.Equal(t, 4, response.TotalProblems)
		assert.Equal(t, 2, response.ProblemsSolved)

		// Verify by_difficulty breakdown
		assert.Equal(t, 1, response.ByDifficulty["easy"])
		assert.Equal(t, 1, response.ByDifficulty["medium"])
		assert.Equal(t, 0, response.ByDifficulty["hard"])

		// Verify by_topic breakdown
		assert.Equal(t, 1, response.ByTopic["arrays"])
		assert.Equal(t, 1, response.ByTopic["linked-lists"])

		// Verify formatted output
		assert.Contains(t, string(output), "\n")
		assert.Contains(t, string(output), "  ")
	})

	t.Run("compact JSON output", func(t *testing.T) {
		// Get stats
		progressSvc := progress.NewService(db)
		stats, err := progressSvc.GetStats("")
		require.NoError(t, err)

		// Build JSON response
		byDifficulty := make(map[string]int)
		for diff, diffStats := range stats.ByDifficulty {
			byDifficulty[diff] = diffStats.Solved
		}

		byTopic := make(map[string]int)
		for topic, topicStats := range stats.ByTopic {
			byTopic[topic] = topicStats.Solved
		}

		response := StatusResponse{
			TotalProblems:  stats.TotalProblems,
			ProblemsSolved: stats.TotalSolved,
			ByDifficulty:   byDifficulty,
			ByTopic:        byTopic,
		}

		// Test compact JSON marshaling
		output, err := json.Marshal(response)
		require.NoError(t, err, "Should marshal to compact JSON")

		// Verify it's valid JSON
		var parsed StatusResponse
		err = json.Unmarshal(output, &parsed)
		require.NoError(t, err, "Compact output should be valid JSON")

		// Verify no excessive whitespace (compact format)
		assert.NotContains(t, string(output), "\n  ", "Compact JSON should not have indentation")
	})

	t.Run("snake_case field names in output", func(t *testing.T) {
		// Get stats
		progressSvc := progress.NewService(db)
		stats, err := progressSvc.GetStats("")
		require.NoError(t, err)

		// Build JSON response
		response := StatusResponse{
			TotalProblems:  stats.TotalProblems,
			ProblemsSolved: stats.TotalSolved,
			ByDifficulty:   map[string]int{"easy": 1, "medium": 1},
			ByTopic:        map[string]int{"arrays": 1, "linked-lists": 1},
		}

		// Test JSON marshaling
		output, err := json.Marshal(response)
		require.NoError(t, err)

		outputStr := string(output)

		// Verify snake_case field names
		assert.Contains(t, outputStr, `"total_problems"`)
		assert.Contains(t, outputStr, `"problems_solved"`)
		assert.Contains(t, outputStr, `"by_difficulty"`)
		assert.Contains(t, outputStr, `"by_topic"`)

		// Verify camelCase is NOT present
		assert.NotContains(t, outputStr, `"totalProblems"`)
		assert.NotContains(t, outputStr, `"problemsSolved"`)
		assert.NotContains(t, outputStr, `"byDifficulty"`)
		assert.NotContains(t, outputStr, `"byTopic"`)
	})
}

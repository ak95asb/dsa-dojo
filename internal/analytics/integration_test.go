package analytics

import (
	"fmt"
	"testing"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Integration tests verify end-to-end analytics workflows

func setupIntegrationDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&database.Problem{}, &database.Progress{}, &database.Solution{})
	require.NoError(t, err)

	return db
}

func TestIntegration_CompleteAnalyticsWorkflow(t *testing.T) {
	db := setupIntegrationDB(t)
	service := NewAnalyticsService(db)

	t.Run("complete workflow: seed data → calculate analytics → verify results", func(t *testing.T) {
		// Seed realistic problem data
		problems := []database.Problem{
			{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", Topic: "arrays"},
			{Slug: "reverse-linked-list", Title: "Reverse Linked List", Difficulty: "easy", Topic: "linked-lists"},
			{Slug: "valid-anagram", Title: "Valid Anagram", Difficulty: "easy", Topic: "strings"},
			{Slug: "binary-search", Title: "Binary Search", Difficulty: "medium", Topic: "arrays"},
			{Slug: "merge-intervals", Title: "Merge Intervals", Difficulty: "medium", Topic: "arrays"},
			{Slug: "lru-cache", Title: "LRU Cache", Difficulty: "hard", Topic: "design"},
		}

		for _, p := range problems {
			require.NoError(t, db.Create(&p).Error)
		}

		// Seed progress data with realistic attempt patterns
		progressRecords := []database.Progress{
			{ProblemID: 1, IsSolved: true, TotalAttempts: 1},  // easy arrays - solved in 1
			{ProblemID: 2, IsSolved: true, TotalAttempts: 2},  // easy linked-lists - solved in 2
			{ProblemID: 3, IsSolved: true, TotalAttempts: 1},  // easy strings - solved in 1
			{ProblemID: 4, IsSolved: true, TotalAttempts: 3},  // medium arrays - solved in 3
			{ProblemID: 5, IsSolved: false, TotalAttempts: 4}, // medium arrays - not solved
			{ProblemID: 6, IsSolved: false, TotalAttempts: 5}, // hard design - not solved
		}

		for _, p := range progressRecords {
			require.NoError(t, db.Create(&p).Error)
		}

		// Calculate analytics
		stats, err := service.CalculateStats(AnalyticsFilter{})
		require.NoError(t, err)
		require.NotNil(t, stats)

		// Verify overall success rate: 4 solved / 6 attempted = 66.67%
		// Note: Arrays topic has 3 problems (2 solved, 1 not) in the data
		expectedRate := (4.0 / 6.0) * 100 // 4 solved out of 6 total
		assert.InDelta(t, expectedRate, stats.OverallSuccessRate, 0.1)

		// Verify success rate by difficulty
		assert.InDelta(t, 100.0, stats.SuccessRateByDifficulty["easy"], 0.1)   // 3/3 = 100%
		assert.InDelta(t, 50.0, stats.SuccessRateByDifficulty["medium"], 0.1)  // 1/2 = 50%
		assert.InDelta(t, 0.0, stats.SuccessRateByDifficulty["hard"], 0.1)     // 0/1 = 0%

		// Verify success rate by topic
		// Arrays: problem 1 (solved), problem 4 (solved), problem 5 (not solved) = 2/3 = 66.67%
		assert.InDelta(t, 66.67, stats.SuccessRateByTopic["arrays"], 0.1)

		// Verify average attempts (only for solved problems)
		// (1 + 2 + 1 + 3) / 4 = 1.75
		assert.InDelta(t, 1.75, stats.AvgAttemptsOverall, 0.1)

		// Verify average attempts by difficulty
		assert.InDelta(t, 1.33, stats.AvgAttemptsByDifficulty["easy"], 0.1) // (1+2+1)/3
		assert.InDelta(t, 3.0, stats.AvgAttemptsByDifficulty["medium"], 0.1) // 3/1

		// Verify practice patterns
		assert.NotEmpty(t, stats.MostPracticedTopic)
		assert.NotEmpty(t, stats.BestDifficulty)
		assert.Equal(t, "easy", stats.BestDifficulty) // 100% success rate
	})
}

func TestIntegration_FilteringWorkflow(t *testing.T) {
	db := setupIntegrationDB(t)
	service := NewAnalyticsService(db)

	// Seed test data
	problems := []database.Problem{
		{Slug: "array-1", Title: "Array 1", Difficulty: "easy", Topic: "arrays"},
		{Slug: "array-2", Title: "Array 2", Difficulty: "medium", Topic: "arrays"},
		{Slug: "tree-1", Title: "Tree 1", Difficulty: "easy", Topic: "trees"},
		{Slug: "tree-2", Title: "Tree 2", Difficulty: "hard", Topic: "trees"},
	}
	for _, p := range problems {
		require.NoError(t, db.Create(&p).Error)
	}

	progressRecords := []database.Progress{
		{ProblemID: 1, IsSolved: true, TotalAttempts: 1},
		{ProblemID: 2, IsSolved: false, TotalAttempts: 3},
		{ProblemID: 3, IsSolved: true, TotalAttempts: 2},
		{ProblemID: 4, IsSolved: false, TotalAttempts: 5},
	}
	for _, p := range progressRecords {
		require.NoError(t, db.Create(&p).Error)
	}

	t.Run("filter by topic - arrays only", func(t *testing.T) {
		stats, err := service.CalculateStats(AnalyticsFilter{Topic: "arrays"})
		require.NoError(t, err)

		// Should only include arrays: 1 solved / 2 attempted = 50%
		assert.InDelta(t, 50.0, stats.OverallSuccessRate, 0.1)

		// Should only have arrays in topic breakdown
		assert.Len(t, stats.SuccessRateByTopic, 1)
		assert.Contains(t, stats.SuccessRateByTopic, "arrays")
	})

	t.Run("filter by difficulty - easy only", func(t *testing.T) {
		stats, err := service.CalculateStats(AnalyticsFilter{Difficulty: "easy"})
		require.NoError(t, err)

		// Should only include easy: 2 solved / 2 attempted = 100%
		assert.InDelta(t, 100.0, stats.OverallSuccessRate, 0.1)

		// Should only have easy in difficulty breakdown
		assert.Len(t, stats.SuccessRateByDifficulty, 1)
		assert.Contains(t, stats.SuccessRateByDifficulty, "easy")
	})

	t.Run("filter by both topic and difficulty", func(t *testing.T) {
		stats, err := service.CalculateStats(AnalyticsFilter{Topic: "arrays", Difficulty: "easy"})
		require.NoError(t, err)

		// Should only include easy arrays: 1 solved / 1 attempted = 100%
		assert.InDelta(t, 100.0, stats.OverallSuccessRate, 0.1)
		assert.InDelta(t, 1.0, stats.AvgAttemptsOverall, 0.1)
	})
}

func TestIntegration_Performance(t *testing.T) {
	db := setupIntegrationDB(t)
	service := NewAnalyticsService(db)

	t.Run("analytics calculation completes in <300ms", func(t *testing.T) {
		// Seed larger dataset (50 problems)
		for i := 1; i <= 50; i++ {
			problem := &database.Problem{
				Slug:       fmt.Sprintf("perf-problem-%d", i),
				Title:      fmt.Sprintf("Problem %d", i),
				Difficulty: []string{"easy", "medium", "hard"}[i%3],
				Topic:      []string{"arrays", "trees", "graphs", "strings"}[i%4],
			}
			require.NoError(t, db.Create(problem).Error)

			progress := &database.Progress{
				ProblemID:    uint(i),
				IsSolved:     i%2 == 0,
				TotalAttempts: (i % 5) + 1,
			}
			require.NoError(t, db.Create(progress).Error)
		}

		// Measure analytics calculation time
		start := time.Now()
		stats, err := service.CalculateStats(AnalyticsFilter{})
		elapsed := time.Since(start)

		require.NoError(t, err)
		require.NotNil(t, stats)
		assert.Less(t, elapsed.Milliseconds(), int64(300), "Analytics calculation should complete in <300ms")
	})

	t.Run("filtered analytics maintains performance", func(t *testing.T) {
		start := time.Now()
		stats, err := service.CalculateStats(AnalyticsFilter{Topic: "arrays"})
		elapsed := time.Since(start)

		require.NoError(t, err)
		require.NotNil(t, stats)
		assert.Less(t, elapsed.Milliseconds(), int64(300), "Filtered analytics should complete in <300ms")
	})
}

func TestIntegration_EdgeCases(t *testing.T) {
	db := setupIntegrationDB(t)
	service := NewAnalyticsService(db)

	t.Run("empty database returns zero stats", func(t *testing.T) {
		stats, err := service.CalculateStats(AnalyticsFilter{})
		require.NoError(t, err)
		require.NotNil(t, stats)

		assert.Equal(t, 0.0, stats.OverallSuccessRate)
		assert.Equal(t, 0.0, stats.AvgAttemptsOverall)
		assert.Empty(t, stats.MostPracticedTopic)
		assert.Empty(t, stats.LeastPracticedTopic)
	})

	t.Run("all problems solved returns 100%", func(t *testing.T) {
		problems := []database.Problem{
			{Slug: "all-1", Title: "All 1", Difficulty: "easy", Topic: "arrays"},
			{Slug: "all-2", Title: "All 2", Difficulty: "medium", Topic: "trees"},
		}
		for _, p := range problems {
			require.NoError(t, db.Create(&p).Error)
		}

		progressRecords := []database.Progress{
			{ProblemID: 1, IsSolved: true, TotalAttempts: 1},
			{ProblemID: 2, IsSolved: true, TotalAttempts: 2},
		}
		for _, p := range progressRecords {
			require.NoError(t, db.Create(&p).Error)
		}

		stats, err := service.CalculateStats(AnalyticsFilter{})
		require.NoError(t, err)

		assert.Equal(t, 100.0, stats.OverallSuccessRate)
	})

	t.Run("no solved problems returns 0% with valid attempt data", func(t *testing.T) {
		// Clear previous data
		db.Exec("DELETE FROM progresses")
		db.Exec("DELETE FROM problems")

		problems := []database.Problem{
			{Slug: "none-1", Title: "None 1", Difficulty: "hard", Topic: "graphs"},
		}
		for _, p := range problems {
			require.NoError(t, db.Create(&p).Error)
		}

		progressRecords := []database.Progress{
			{ProblemID: 1, IsSolved: false, TotalAttempts: 3},
		}
		for _, p := range progressRecords {
			require.NoError(t, db.Create(&p).Error)
		}

		stats, err := service.CalculateStats(AnalyticsFilter{})
		require.NoError(t, err)

		assert.Equal(t, 0.0, stats.OverallSuccessRate)
		assert.Equal(t, 0.0, stats.AvgAttemptsOverall) // No solved = no average
	})
}

func TestIntegration_DataIntegrity(t *testing.T) {
	db := setupIntegrationDB(t)
	service := NewAnalyticsService(db)

	t.Run("analytics don't modify database", func(t *testing.T) {
		// Seed data
		problem := &database.Problem{Slug: "integrity", Title: "Integrity", Difficulty: "easy", Topic: "arrays"}
		require.NoError(t, db.Create(problem).Error)

		progress := &database.Progress{ProblemID: 1, IsSolved: true, TotalAttempts: 2}
		require.NoError(t, db.Create(progress).Error)

		// Count records before
		var countBefore int64
		db.Model(&database.Progress{}).Count(&countBefore)

		// Calculate analytics
		_, err := service.CalculateStats(AnalyticsFilter{})
		require.NoError(t, err)

		// Count records after
		var countAfter int64
		db.Model(&database.Progress{}).Count(&countAfter)

		// Verify no changes
		assert.Equal(t, countBefore, countAfter)

		// Verify data unchanged
		var progressAfter database.Progress
		db.First(&progressAfter, "problem_id = ?", 1)
		assert.True(t, progressAfter.IsSolved)
		assert.Equal(t, 2, progressAfter.TotalAttempts)
	})
}

func TestIntegration_PracticePatternInsights(t *testing.T) {
	db := setupIntegrationDB(t)
	service := NewAnalyticsService(db)

	t.Run("identifies practice patterns correctly", func(t *testing.T) {
		// Create scenario where arrays are most practiced, graphs least practiced
		problems := []database.Problem{
			{Slug: "arr-1", Title: "Array 1", Difficulty: "easy", Topic: "arrays"},
			{Slug: "arr-2", Title: "Array 2", Difficulty: "easy", Topic: "arrays"},
			{Slug: "arr-3", Title: "Array 3", Difficulty: "medium", Topic: "arrays"},
			{Slug: "graph-1", Title: "Graph 1", Difficulty: "hard", Topic: "graphs"},
		}
		for _, p := range problems {
			require.NoError(t, db.Create(&p).Error)
		}

		progressRecords := []database.Progress{
			{ProblemID: 1, IsSolved: true, TotalAttempts: 5},   // arrays: 5 attempts
			{ProblemID: 2, IsSolved: true, TotalAttempts: 3},   // arrays: 3 attempts
			{ProblemID: 3, IsSolved: false, TotalAttempts: 7},  // arrays: 7 attempts (total: 15)
			{ProblemID: 4, IsSolved: false, TotalAttempts: 2},  // graphs: 2 attempts
		}
		for _, p := range progressRecords {
			require.NoError(t, db.Create(&p).Error)
		}

		stats, err := service.CalculateStats(AnalyticsFilter{})
		require.NoError(t, err)

		// Most practiced should be arrays (15 total attempts)
		assert.Equal(t, "arrays", stats.MostPracticedTopic)

		// Least practiced should be graphs (2 total attempts)
		assert.Equal(t, "graphs", stats.LeastPracticedTopic)

		// Best difficulty should be easy (2/2 = 100%)
		assert.Equal(t, "easy", stats.BestDifficulty)

		// Challenging should be hard (0/1 = 0%)
		assert.Equal(t, "hard", stats.ChallengingDifficulty)
	})
}

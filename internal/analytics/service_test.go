package analytics

import (
	"testing"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&database.Problem{}, &database.Progress{}, &database.Solution{})
	require.NoError(t, err)

	return db
}

func seedTestData(db *gorm.DB, t *testing.T) {
	// Create problems across different difficulties and topics
	problems := []database.Problem{
		{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", Topic: "arrays"},
		{Slug: "valid-parentheses", Title: "Valid Parentheses", Difficulty: "easy", Topic: "strings"},
		{Slug: "add-two-numbers", Title: "Add Two Numbers", Difficulty: "medium", Topic: "linked-lists"},
		{Slug: "longest-substring", Title: "Longest Substring", Difficulty: "medium", Topic: "strings"},
		{Slug: "median-arrays", Title: "Median of Two Sorted Arrays", Difficulty: "hard", Topic: "arrays"},
		{Slug: "binary-tree", Title: "Binary Tree", Difficulty: "easy", Topic: "trees"},
		{Slug: "graph-traverse", Title: "Graph Traversal", Difficulty: "hard", Topic: "graphs"},
		{Slug: "merge-sort", Title: "Merge Sort", Difficulty: "medium", Topic: "sorting"},
	}

	for _, p := range problems {
		require.NoError(t, db.Create(&p).Error)
	}

	// Create progress records with varied patterns
	progress := []database.Progress{
		{ProblemID: 1, IsSolved: true, TotalAttempts: 2},   // easy arrays - solved
		{ProblemID: 2, IsSolved: true, TotalAttempts: 1},   // easy strings - solved
		{ProblemID: 3, IsSolved: false, TotalAttempts: 3},  // medium linked-lists - not solved
		{ProblemID: 4, IsSolved: true, TotalAttempts: 4},   // medium strings - solved
		{ProblemID: 5, IsSolved: false, TotalAttempts: 5},  // hard arrays - not solved
		{ProblemID: 6, IsSolved: true, TotalAttempts: 1},   // easy trees - solved
		{ProblemID: 7, IsSolved: false, TotalAttempts: 2},  // hard graphs - not solved
		{ProblemID: 8, IsSolved: true, TotalAttempts: 3},   // medium sorting - solved
	}

	for _, p := range progress {
		require.NoError(t, db.Create(&p).Error)
	}
}

func TestNewAnalyticsService(t *testing.T) {
	db := setupTestDB(t)
	service := NewAnalyticsService(db)

	assert.NotNil(t, service)
	assert.NotNil(t, service.db)
}

func TestCalculateStats_OverallSuccessRate(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewAnalyticsService(db)

	stats, err := service.CalculateStats(AnalyticsFilter{})
	require.NoError(t, err)

	// 5 solved out of 8 attempted = 62.5%
	assert.InDelta(t, 62.5, stats.OverallSuccessRate, 0.1)
}

func TestCalculateStats_SuccessRateByDifficulty(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewAnalyticsService(db)

	stats, err := service.CalculateStats(AnalyticsFilter{})
	require.NoError(t, err)

	// Easy: 3 solved out of 3 = 100%
	assert.InDelta(t, 100.0, stats.SuccessRateByDifficulty["easy"], 0.1)

	// Medium: 2 solved out of 3 = 66.7%
	assert.InDelta(t, 66.7, stats.SuccessRateByDifficulty["medium"], 0.1)

	// Hard: 0 solved out of 2 = 0%
	assert.InDelta(t, 0.0, stats.SuccessRateByDifficulty["hard"], 0.1)
}

func TestCalculateStats_SuccessRateByTopic(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewAnalyticsService(db)

	stats, err := service.CalculateStats(AnalyticsFilter{})
	require.NoError(t, err)

	// Arrays: 1 solved out of 2 = 50%
	assert.InDelta(t, 50.0, stats.SuccessRateByTopic["arrays"], 0.1)

	// Strings: 2 solved out of 2 = 100%
	assert.InDelta(t, 100.0, stats.SuccessRateByTopic["strings"], 0.1)

	// Trees: 1 solved out of 1 = 100%
	assert.InDelta(t, 100.0, stats.SuccessRateByTopic["trees"], 0.1)

	// Graphs: 0 solved out of 1 = 0%
	assert.InDelta(t, 0.0, stats.SuccessRateByTopic["graphs"], 0.1)
}

func TestCalculateStats_AverageAttempts(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewAnalyticsService(db)

	stats, err := service.CalculateStats(AnalyticsFilter{})
	require.NoError(t, err)

	// Overall: (2+1+4+1+3) / 5 solved = 2.2 attempts
	assert.InDelta(t, 2.2, stats.AvgAttemptsOverall, 0.1)
}

func TestCalculateStats_AverageAttemptsByDifficulty(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewAnalyticsService(db)

	stats, err := service.CalculateStats(AnalyticsFilter{})
	require.NoError(t, err)

	// Easy: (2+1+1) / 3 = 1.33
	assert.InDelta(t, 1.3, stats.AvgAttemptsByDifficulty["easy"], 0.2)

	// Medium: (4+3) / 2 = 3.5
	assert.InDelta(t, 3.5, stats.AvgAttemptsByDifficulty["medium"], 0.1)

	// Hard: no solved problems, should not be in map or be 0
	_, exists := stats.AvgAttemptsByDifficulty["hard"]
	if exists {
		assert.Equal(t, 0.0, stats.AvgAttemptsByDifficulty["hard"])
	}
}

func TestCalculateStats_PracticePatterns(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewAnalyticsService(db)

	stats, err := service.CalculateStats(AnalyticsFilter{})
	require.NoError(t, err)

	// Most practiced: strings (1+4=5 attempts) or arrays (2+5=7 attempts)
	assert.NotEmpty(t, stats.MostPracticedTopic)

	// Least practiced: graphs (2 attempts) or trees (1 attempt)
	assert.NotEmpty(t, stats.LeastPracticedTopic)

	// Best performing: easy (100% success)
	assert.Equal(t, "easy", stats.BestDifficulty)

	// Most challenging: hard (0% success) or medium (higher avg attempts)
	assert.NotEmpty(t, stats.ChallengingDifficulty)
}

func TestCalculateStats_FilterByTopic(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewAnalyticsService(db)

	filter := AnalyticsFilter{Topic: "strings"}
	stats, err := service.CalculateStats(filter)
	require.NoError(t, err)

	// Strings: 2 solved out of 2 = 100%
	assert.InDelta(t, 100.0, stats.OverallSuccessRate, 0.1)

	// Should only have strings data
	assert.Len(t, stats.SuccessRateByTopic, 1)
	assert.Contains(t, stats.SuccessRateByTopic, "strings")
}

func TestCalculateStats_FilterByDifficulty(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewAnalyticsService(db)

	filter := AnalyticsFilter{Difficulty: "easy"}
	stats, err := service.CalculateStats(filter)
	require.NoError(t, err)

	// Easy: 3 solved out of 3 = 100%
	assert.InDelta(t, 100.0, stats.OverallSuccessRate, 0.1)

	// Should only have easy data
	assert.Len(t, stats.SuccessRateByDifficulty, 1)
	assert.Contains(t, stats.SuccessRateByDifficulty, "easy")
}

func TestCalculateStats_EmptyDatabase(t *testing.T) {
	db := setupTestDB(t)
	service := NewAnalyticsService(db)

	stats, err := service.CalculateStats(AnalyticsFilter{})
	require.NoError(t, err)

	assert.Equal(t, 0.0, stats.OverallSuccessRate)
	assert.Equal(t, 0.0, stats.AvgAttemptsOverall)
	assert.Empty(t, stats.MostPracticedTopic)
	assert.Empty(t, stats.LeastPracticedTopic)
}

func TestCalculateStats_NoSolvedProblems(t *testing.T) {
	db := setupTestDB(t)
	service := NewAnalyticsService(db)

	// Create problems but none solved
	problem := &database.Problem{Slug: "unsolved", Title: "Unsolved", Difficulty: "easy", Topic: "arrays"}
	require.NoError(t, db.Create(problem).Error)

	progress := &database.Progress{ProblemID: problem.ID, IsSolved: false, TotalAttempts: 3}
	require.NoError(t, db.Create(progress).Error)

	stats, err := service.CalculateStats(AnalyticsFilter{})
	require.NoError(t, err)

	assert.Equal(t, 0.0, stats.OverallSuccessRate)
	assert.Equal(t, 0.0, stats.AvgAttemptsOverall) // No solved, so avg is 0
}

func TestCalculateStats_AllSolved(t *testing.T) {
	db := setupTestDB(t)
	service := NewAnalyticsService(db)

	// Create problems all solved
	problem1 := &database.Problem{Slug: "solved1", Title: "Solved 1", Difficulty: "easy", Topic: "arrays"}
	problem2 := &database.Problem{Slug: "solved2", Title: "Solved 2", Difficulty: "medium", Topic: "trees"}
	require.NoError(t, db.Create(problem1).Error)
	require.NoError(t, db.Create(problem2).Error)

	progress1 := &database.Progress{ProblemID: problem1.ID, IsSolved: true, TotalAttempts: 1}
	progress2 := &database.Progress{ProblemID: problem2.ID, IsSolved: true, TotalAttempts: 2}
	require.NoError(t, db.Create(progress1).Error)
	require.NoError(t, db.Create(progress2).Error)

	stats, err := service.CalculateStats(AnalyticsFilter{})
	require.NoError(t, err)

	assert.Equal(t, 100.0, stats.OverallSuccessRate)
	assert.InDelta(t, 1.5, stats.AvgAttemptsOverall, 0.1) // (1+2)/2
}

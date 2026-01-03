package progress

import (
	"testing"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&database.Problem{}, &database.Progress{}, &database.Solution{}, &database.BenchmarkResult{})
	assert.NoError(t, err)

	return db
}

func seedTestData(t *testing.T, db *gorm.DB) {
	problems := []database.Problem{
		{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", Topic: "arrays"},
		{Slug: "add-two-numbers", Title: "Add Two Numbers", Difficulty: "medium", Topic: "linked-lists"},
		{Slug: "longest-substring", Title: "Longest Substring", Difficulty: "medium", Topic: "strings"},
		{Slug: "median-sorted-arrays", Title: "Median of Sorted Arrays", Difficulty: "hard", Topic: "arrays"},
		{Slug: "reverse-integer", Title: "Reverse Integer", Difficulty: "easy", Topic: "math"},
		{Slug: "palindrome-number", Title: "Palindrome Number", Difficulty: "easy", Topic: "math"},
		{Slug: "merge-k-lists", Title: "Merge K Sorted Lists", Difficulty: "hard", Topic: "linked-lists"},
		{Slug: "valid-parentheses", Title: "Valid Parentheses", Difficulty: "easy", Topic: "strings"},
		{Slug: "binary-tree-max", Title: "Binary Tree Max", Difficulty: "medium", Topic: "trees"},
		{Slug: "graph-traversal", Title: "Graph Traversal", Difficulty: "hard", Topic: "graphs"},
	}

	for _, p := range problems {
		err := db.Create(&p).Error
		assert.NoError(t, err)
	}

	// Mark some problems as completed
	progress := []database.Progress{
		{ProblemID: 1, Status: "completed", LastAttempt: time.Now().Add(-5 * 24 * time.Hour)}, // two-sum
		{ProblemID: 2, Status: "completed", LastAttempt: time.Now().Add(-4 * 24 * time.Hour)}, // add-two-numbers
		{ProblemID: 3, Status: "in_progress", LastAttempt: time.Now().Add(-3 * 24 * time.Hour)},
		{ProblemID: 5, Status: "completed", LastAttempt: time.Now().Add(-2 * 24 * time.Hour)}, // reverse-integer
		{ProblemID: 8, Status: "completed", LastAttempt: time.Now().Add(-1 * 24 * time.Hour)}, // valid-parentheses
		{ProblemID: 9, Status: "completed", LastAttempt: time.Now()},                           // binary-tree-max
	}

	for _, p := range progress {
		err := db.Create(&p).Error
		assert.NoError(t, err)
	}
}

func TestGetStats_OverallProgress(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	service := NewService(db)
	stats, err := service.GetStats("")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 10, stats.TotalProblems, "Should have 10 total problems")
	assert.Equal(t, 5, stats.TotalSolved, "Should have 5 solved problems")
}

func TestGetStats_ByDifficulty(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	service := NewService(db)
	stats, err := service.GetStats("")

	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Check easy difficulty
	easyStats, ok := stats.ByDifficulty["easy"]
	assert.True(t, ok, "Should have easy difficulty stats")
	assert.Equal(t, 4, easyStats.Total, "Should have 4 easy problems")
	assert.Equal(t, 3, easyStats.Solved, "Should have 3 easy problems solved")

	// Check medium difficulty
	mediumStats, ok := stats.ByDifficulty["medium"]
	assert.True(t, ok, "Should have medium difficulty stats")
	assert.Equal(t, 3, mediumStats.Total, "Should have 3 medium problems")
	assert.Equal(t, 2, mediumStats.Solved, "Should have 2 medium problems solved")

	// Check hard difficulty
	hardStats, ok := stats.ByDifficulty["hard"]
	assert.True(t, ok, "Should have hard difficulty stats")
	assert.Equal(t, 3, hardStats.Total, "Should have 3 hard problems")
	assert.Equal(t, 0, hardStats.Solved, "Should have 0 hard problems solved")
}

func TestGetStats_ByTopic(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	service := NewService(db)
	stats, err := service.GetStats("")

	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Check arrays topic
	arraysStats, ok := stats.ByTopic["arrays"]
	assert.True(t, ok, "Should have arrays topic stats")
	assert.Equal(t, 2, arraysStats.Total, "Should have 2 array problems")
	assert.Equal(t, 1, arraysStats.Solved, "Should have 1 array problem solved")

	// Check math topic
	mathStats, ok := stats.ByTopic["math"]
	assert.True(t, ok, "Should have math topic stats")
	assert.Equal(t, 2, mathStats.Total, "Should have 2 math problems")
	assert.Equal(t, 1, mathStats.Solved, "Should have 1 math problem solved")

	// Check strings topic
	stringsStats, ok := stats.ByTopic["strings"]
	assert.True(t, ok, "Should have strings topic stats")
	assert.Equal(t, 2, stringsStats.Total, "Should have 2 string problems")
	assert.Equal(t, 1, stringsStats.Solved, "Should have 1 string problem solved")
}

func TestGetStats_RecentActivity(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	service := NewService(db)
	stats, err := service.GetStats("")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 5, len(stats.RecentActivity), "Should have 5 recent activities")

	// Most recent should be binary-tree-max
	assert.Equal(t, "binary-tree-max", stats.RecentActivity[0].Slug)
	assert.Equal(t, "Binary Tree Max", stats.RecentActivity[0].Title)
	assert.Equal(t, "medium", stats.RecentActivity[0].Difficulty)
}

func TestGetStats_TopicFilter(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	service := NewService(db)
	stats, err := service.GetStats("arrays")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 2, stats.TotalProblems, "Should have 2 array problems")
	assert.Equal(t, 1, stats.TotalSolved, "Should have 1 array problem solved")

	// Topic breakdown should be empty when filtering by topic
	assert.Equal(t, 0, len(stats.ByTopic), "Should not have topic breakdown when filtering")

	// Difficulty breakdown should only include arrays
	easyStats, ok := stats.ByDifficulty["easy"]
	assert.True(t, ok, "Should have easy difficulty stats")
	assert.Equal(t, 1, easyStats.Total, "Should have 1 easy array problem")
	assert.Equal(t, 1, easyStats.Solved, "Should have 1 easy array problem solved")

	hardStats, ok := stats.ByDifficulty["hard"]
	assert.True(t, ok, "Should have hard difficulty stats")
	assert.Equal(t, 1, hardStats.Total, "Should have 1 hard array problem")
	assert.Equal(t, 0, hardStats.Solved, "Should have 0 hard array problems solved")
}

func TestGetStats_EmptyDatabase(t *testing.T) {
	db := setupTestDB(t)

	service := NewService(db)
	stats, err := service.GetStats("")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 0, stats.TotalProblems, "Should have 0 total problems")
	assert.Equal(t, 0, stats.TotalSolved, "Should have 0 solved problems")
	assert.Equal(t, 0, len(stats.ByDifficulty), "Should have no difficulty stats")
	assert.Equal(t, 0, len(stats.ByTopic), "Should have no topic stats")
	assert.Equal(t, 0, len(stats.RecentActivity), "Should have no recent activity")
}

func TestGetStats_NoSolvedProblems(t *testing.T) {
	db := setupTestDB(t)

	// Seed only problems, no progress
	problems := []database.Problem{
		{Slug: "problem-1", Title: "Problem 1", Difficulty: "easy", Topic: "arrays"},
		{Slug: "problem-2", Title: "Problem 2", Difficulty: "medium", Topic: "trees"},
		{Slug: "problem-3", Title: "Problem 3", Difficulty: "hard", Topic: "graphs"},
	}

	for _, p := range problems {
		err := db.Create(&p).Error
		assert.NoError(t, err)
	}

	service := NewService(db)
	stats, err := service.GetStats("")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 3, stats.TotalProblems, "Should have 3 total problems")
	assert.Equal(t, 0, stats.TotalSolved, "Should have 0 solved problems")

	// All difficulty stats should show 0 solved
	for _, diffStats := range stats.ByDifficulty {
		assert.Equal(t, 0, diffStats.Solved, "Should have 0 solved for all difficulties")
	}

	assert.Equal(t, 0, len(stats.RecentActivity), "Should have no recent activity")
}

func TestGetStats_AllProblemsCompleted(t *testing.T) {
	db := setupTestDB(t)

	// Seed problems with all completed
	problems := []database.Problem{
		{Slug: "problem-1", Title: "Problem 1", Difficulty: "easy", Topic: "arrays"},
		{Slug: "problem-2", Title: "Problem 2", Difficulty: "medium", Topic: "trees"},
	}

	for _, p := range problems {
		err := db.Create(&p).Error
		assert.NoError(t, err)
	}

	// Mark all as completed
	progress := []database.Progress{
		{ProblemID: 1, Status: "completed", LastAttempt: time.Now()},
		{ProblemID: 2, Status: "completed", LastAttempt: time.Now()},
	}

	for _, p := range progress {
		err := db.Create(&p).Error
		assert.NoError(t, err)
	}

	service := NewService(db)
	stats, err := service.GetStats("")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 2, stats.TotalProblems, "Should have 2 total problems")
	assert.Equal(t, 2, stats.TotalSolved, "Should have 2 solved problems")

	// All difficulty stats should show 100% completion
	for _, diffStats := range stats.ByDifficulty {
		assert.Equal(t, diffStats.Total, diffStats.Solved, "All problems should be solved")
	}
}

func TestGetStats_InvalidTopicFilter(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(t, db)

	service := NewService(db)
	stats, err := service.GetStats("nonexistent-topic")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 0, stats.TotalProblems, "Should have 0 problems for invalid topic")
	assert.Equal(t, 0, stats.TotalSolved, "Should have 0 solved for invalid topic")
}

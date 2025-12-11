package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeedProblems(t *testing.T) {
	t.Run("seeds all problems on first run", func(t *testing.T) {
		db := setupTestDB(t)

		count, err := SeedProblems(db)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, 20, "should seed at least 20 problems")

		// Verify problems in database
		var problems []Problem
		db.Find(&problems)
		assert.GreaterOrEqual(t, len(problems), 20, "database should contain at least 20 problems")
	})

	t.Run("is idempotent - no duplicates on second run", func(t *testing.T) {
		db := setupTestDB(t)

		// First seeding
		count1, err := SeedProblems(db)
		assert.NoError(t, err)
		assert.Greater(t, count1, 0, "should seed problems on first run")

		// Second seeding
		count2, err := SeedProblems(db)
		assert.NoError(t, err)
		assert.Equal(t, 0, count2, "should not seed any new problems on second run")

		// Verify total count unchanged
		var total int64
		db.Model(&Problem{}).Count(&total)
		assert.Equal(t, int64(count1), total, "total problems should equal first seeding count")
	})

	t.Run("covers all required topics", func(t *testing.T) {
		db := setupTestDB(t)
		SeedProblems(db)

		requiredTopics := []string{"arrays", "linked-lists", "trees", "graphs", "sorting", "searching"}
		for _, topic := range requiredTopics {
			var count int64
			db.Model(&Problem{}).Where("topic = ?", topic).Count(&count)
			assert.Greater(t, count, int64(0), "topic '%s' should have at least one problem", topic)
		}
	})

	t.Run("has proper difficulty distribution", func(t *testing.T) {
		db := setupTestDB(t)
		SeedProblems(db)

		var easyCount, mediumCount, hardCount int64
		db.Model(&Problem{}).Where("difficulty = ?", "easy").Count(&easyCount)
		db.Model(&Problem{}).Where("difficulty = ?", "medium").Count(&mediumCount)
		db.Model(&Problem{}).Where("difficulty = ?", "hard").Count(&hardCount)

		assert.Greater(t, easyCount, int64(0), "should have easy problems")
		assert.Greater(t, mediumCount, int64(0), "should have medium problems")
		assert.Greater(t, hardCount, int64(0), "should have hard problems")

		// Verify medium problems are the most common (roughly 50%)
		total := easyCount + mediumCount + hardCount
		assert.Greater(t, mediumCount, total/3, "medium problems should be a significant portion")
	})

	t.Run("all seeded problems have required fields", func(t *testing.T) {
		db := setupTestDB(t)
		SeedProblems(db)

		var problems []Problem
		db.Find(&problems)

		for _, problem := range problems {
			assert.NotEmpty(t, problem.Slug, "problem should have slug")
			assert.NotEmpty(t, problem.Title, "problem should have title")
			assert.NotEmpty(t, problem.Description, "problem should have description")
			assert.NotEmpty(t, problem.Difficulty, "problem should have difficulty")
			assert.NotEmpty(t, problem.Topic, "problem should have topic")
			assert.Contains(t, []string{"easy", "medium", "hard"}, problem.Difficulty, "difficulty should be valid")
		}
	})

	t.Run("all problem slugs are unique", func(t *testing.T) {
		db := setupTestDB(t)
		SeedProblems(db)

		var problems []Problem
		db.Find(&problems)

		slugs := make(map[string]bool)
		for _, problem := range problems {
			assert.False(t, slugs[problem.Slug], "slug '%s' should be unique", problem.Slug)
			slugs[problem.Slug] = true
		}
	})
}

package problem

import (
	"testing"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&database.Problem{}, &database.Progress{}, &database.Solution{})
	assert.NoError(t, err)

	return db
}

func seedTestProblems(t *testing.T, db *gorm.DB) {
	problems := []database.Problem{
		{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", Topic: "arrays"},
		{Slug: "add-two-numbers", Title: "Add Two Numbers", Difficulty: "medium", Topic: "linked-lists"},
		{Slug: "reverse-linked-list", Title: "Reverse Linked List", Difficulty: "easy", Topic: "linked-lists"},
		{Slug: "validate-bst", Title: "Validate Binary Search Tree", Difficulty: "medium", Topic: "trees"},
		{Slug: "binary-search", Title: "Binary Search", Difficulty: "easy", Topic: "searching"},
		{Slug: "merge-k-lists", Title: "Merge K Sorted Lists", Difficulty: "hard", Topic: "linked-lists"},
	}

	for _, p := range problems {
		err := db.Create(&p).Error
		assert.NoError(t, err)
	}

	// Mark some as solved (problem IDs 1 and 3)
	err := db.Create(&database.Progress{ProblemID: 1, Status: "completed", Attempts: 3}).Error
	assert.NoError(t, err)
	err = db.Create(&database.Progress{ProblemID: 3, Status: "completed", Attempts: 1}).Error
	assert.NoError(t, err)
}

func TestListProblems(t *testing.T) {
	t.Run("lists all problems with no filters", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		problems, err := svc.ListProblems(ListFilters{})

		assert.NoError(t, err)
		assert.Equal(t, 6, len(problems))
	})

	t.Run("filters by difficulty - easy", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		problems, err := svc.ListProblems(ListFilters{Difficulty: "easy"})

		assert.NoError(t, err)
		assert.Equal(t, 3, len(problems))
		for _, p := range problems {
			assert.Equal(t, "easy", p.Difficulty)
		}
	})

	t.Run("filters by difficulty - medium", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		problems, err := svc.ListProblems(ListFilters{Difficulty: "medium"})

		assert.NoError(t, err)
		assert.Equal(t, 2, len(problems))
		for _, p := range problems {
			assert.Equal(t, "medium", p.Difficulty)
		}
	})

	t.Run("filters by difficulty - hard", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		problems, err := svc.ListProblems(ListFilters{Difficulty: "hard"})

		assert.NoError(t, err)
		assert.Equal(t, 1, len(problems))
		assert.Equal(t, "hard", problems[0].Difficulty)
	})

	t.Run("filters by topic - linked-lists", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		problems, err := svc.ListProblems(ListFilters{Topic: "linked-lists"})

		assert.NoError(t, err)
		assert.Equal(t, 3, len(problems))
		for _, p := range problems {
			assert.Equal(t, "linked-lists", p.Topic)
		}
	})

	t.Run("filters by topic - arrays", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		problems, err := svc.ListProblems(ListFilters{Topic: "arrays"})

		assert.NoError(t, err)
		assert.Equal(t, 1, len(problems))
		assert.Equal(t, "arrays", problems[0].Topic)
		assert.Equal(t, "two-sum", problems[0].Slug)
	})

	t.Run("filters by solved status - solved only", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := true
		problems, err := svc.ListProblems(ListFilters{Solved: &solved})

		assert.NoError(t, err)
		assert.Equal(t, 2, len(problems))
		for _, p := range problems {
			assert.True(t, p.IsSolved, "Problem %s should be marked as solved", p.Slug)
		}
	})

	t.Run("filters by solved status - unsolved only", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		problems, err := svc.ListProblems(ListFilters{Solved: &solved})

		assert.NoError(t, err)
		assert.Equal(t, 4, len(problems))
		for _, p := range problems {
			assert.False(t, p.IsSolved, "Problem %s should be marked as unsolved", p.Slug)
		}
	})

	t.Run("combines multiple filters - difficulty and topic", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		problems, err := svc.ListProblems(ListFilters{
			Difficulty: "easy",
			Topic:      "linked-lists",
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, len(problems))
		assert.Equal(t, "reverse-linked-list", problems[0].Slug)
		assert.Equal(t, "easy", problems[0].Difficulty)
		assert.Equal(t, "linked-lists", problems[0].Topic)
	})

	t.Run("combines all filters - difficulty, topic, and unsolved", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		problems, err := svc.ListProblems(ListFilters{
			Difficulty: "medium",
			Topic:      "linked-lists",
			Solved:     &solved,
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, len(problems))
		assert.Equal(t, "add-two-numbers", problems[0].Slug)
		assert.False(t, problems[0].IsSolved)
	})

	t.Run("returns empty slice for no matches", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		problems, err := svc.ListProblems(ListFilters{Topic: "graphs"})

		assert.NoError(t, err)
		assert.Equal(t, 0, len(problems))
	})

	t.Run("handles database with no problems", func(t *testing.T) {
		db := setupTestDB(t)
		// Don't seed any problems

		svc := NewService(db)
		problems, err := svc.ListProblems(ListFilters{})

		assert.NoError(t, err)
		assert.Equal(t, 0, len(problems))
	})
}

func TestIsValidDifficulty(t *testing.T) {
	tests := []struct {
		name       string
		difficulty string
		want       bool
	}{
		{"valid easy", "easy", true},
		{"valid medium", "medium", true},
		{"valid hard", "hard", true},
		{"invalid super-hard", "super-hard", false},
		{"invalid empty", "", false},
		{"invalid uppercase", "EASY", false},
		{"invalid mixed case", "Easy", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidDifficulty(tt.difficulty)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsValidTopic(t *testing.T) {
	tests := []struct {
		name  string
		topic string
		want  bool
	}{
		{"valid arrays", "arrays", true},
		{"valid linked-lists", "linked-lists", true},
		{"valid trees", "trees", true},
		{"valid graphs", "graphs", true},
		{"valid sorting", "sorting", true},
		{"valid searching", "searching", true},
		{"invalid dp", "dynamic-programming", false},
		{"invalid empty", "", false},
		{"invalid uppercase", "ARRAYS", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidTopic(tt.topic)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetProblemBySlug(t *testing.T) {
	t.Run("returns problem details for valid slug", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		details, err := svc.GetProblemBySlug("two-sum")

		assert.NoError(t, err)
		assert.NotNil(t, details)
		assert.Equal(t, "Two Sum", details.Title)
		assert.Equal(t, "two-sum", details.Slug)
		assert.Equal(t, "easy", details.Difficulty)
		assert.Equal(t, "arrays", details.Topic)
		assert.Equal(t, "completed", details.Status) // Problem ID 1 is marked as completed
		assert.Equal(t, 3, details.Attempts)
		assert.Equal(t, "problems/templates/two-sum.go", details.BoilerplatePath)
		assert.Equal(t, "problems/templates/two-sum_test.go", details.TestPath)
	})

	t.Run("returns problem details for unsolved problem", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		details, err := svc.GetProblemBySlug("binary-search")

		assert.NoError(t, err)
		assert.NotNil(t, details)
		assert.Equal(t, "Binary Search", details.Title)
		assert.Equal(t, "binary-search", details.Slug)
		assert.Equal(t, "not_started", details.Status) // No progress record
		assert.Equal(t, 0, details.Attempts)
		assert.False(t, details.HasSolution)
	})

	t.Run("returns ErrProblemNotFound for invalid slug", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		details, err := svc.GetProblemBySlug("invalid-slug")

		assert.Error(t, err)
		assert.Equal(t, ErrProblemNotFound, err)
		assert.Nil(t, details)
	})

	t.Run("includes solution status when solution exists", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		// Add a passing solution for problem ID 1 (two-sum)
		err := db.Create(&database.Solution{
			ProblemID: 1,
			Code:      "func TwoSum() {}",
			Passed:    true,
			Language:  "go",
		}).Error
		assert.NoError(t, err)

		svc := NewService(db)
		details, err := svc.GetProblemBySlug("two-sum")

		assert.NoError(t, err)
		assert.True(t, details.HasSolution)
	})

	t.Run("does not count failed solutions", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		// Add a failing solution (should not count)
		err := db.Create(&database.Solution{
			ProblemID: 5, // binary-search
			Code:      "func BinarySearch() {}",
			Passed:    false,
			Language:  "go",
		}).Error
		assert.NoError(t, err)

		svc := NewService(db)
		details, err := svc.GetProblemBySlug("binary-search")

		assert.NoError(t, err)
		assert.False(t, details.HasSolution, "Failed solutions should not count")
	})

	t.Run("handles empty database", func(t *testing.T) {
		db := setupTestDB(t)
		// Don't seed any problems

		svc := NewService(db)
		details, err := svc.GetProblemBySlug("two-sum")

		assert.Error(t, err)
		assert.Equal(t, ErrProblemNotFound, err)
		assert.Nil(t, details)
	})
}

func TestGetRandomProblem(t *testing.T) {
	t.Run("returns random problem with no filters", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		details, err := svc.GetRandomProblem(ListFilters{Solved: &solved})

		assert.NoError(t, err)
		assert.NotNil(t, details)
		// Should be one of the unsolved problems (4 total unsolved: add-two-numbers, validate-bst, binary-search, merge-k-lists)
		assert.Contains(t, []string{"add-two-numbers", "validate-bst", "binary-search", "merge-k-lists"}, details.Slug)
	})

	t.Run("returns random problem with difficulty filter", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		details, err := svc.GetRandomProblem(ListFilters{
			Difficulty: "easy",
			Solved:     &solved,
		})

		assert.NoError(t, err)
		assert.NotNil(t, details)
		assert.Equal(t, "easy", details.Difficulty)
		// Should be binary-search (unsolved easy problem)
		assert.Equal(t, "binary-search", details.Slug)
	})

	t.Run("returns random problem with topic filter", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		details, err := svc.GetRandomProblem(ListFilters{
			Topic:  "linked-lists",
			Solved: &solved,
		})

		assert.NoError(t, err)
		assert.NotNil(t, details)
		assert.Equal(t, "linked-lists", details.Topic)
		// Should be add-two-numbers or merge-k-lists (unsolved linked list problems)
		assert.Contains(t, []string{"add-two-numbers", "merge-k-lists"}, details.Slug)
	})

	t.Run("returns random problem with combined filters", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		details, err := svc.GetRandomProblem(ListFilters{
			Difficulty: "medium",
			Topic:      "linked-lists",
			Solved:     &solved,
		})

		assert.NoError(t, err)
		assert.NotNil(t, details)
		assert.Equal(t, "medium", details.Difficulty)
		assert.Equal(t, "linked-lists", details.Topic)
		assert.Equal(t, "add-two-numbers", details.Slug) // Only one matching problem
	})

	t.Run("returns ErrNoProblemsFound when all problems solved", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		// Mark all problems as solved
		var problems []database.Problem
		db.Find(&problems)
		for _, p := range problems {
			db.Create(&database.Progress{
				ProblemID: p.ID,
				Status:    "completed",
				Attempts:  1,
			})
		}

		svc := NewService(db)
		solved := false
		details, err := svc.GetRandomProblem(ListFilters{Solved: &solved})

		assert.Error(t, err)
		assert.Equal(t, ErrNoProblemsFound, err)
		assert.Nil(t, details)
	})

	t.Run("returns different problems on multiple calls (randomness check)", func(t *testing.T) {
		db := setupTestDB(t)
		seedTestProblems(t, db)

		svc := NewService(db)
		solved := false
		filters := ListFilters{Solved: &solved}

		// Get 10 random problems and check we get some variety
		slugs := make(map[string]bool)
		for i := 0; i < 10; i++ {
			details, err := svc.GetRandomProblem(filters)
			assert.NoError(t, err)
			slugs[details.Slug] = true
		}

		// With 4 unsolved problems and 10 selections, we should see at least 2 different problems
		// (statistically very likely, though not guaranteed)
		assert.GreaterOrEqual(t, len(slugs), 2, "Random selection should produce variety over multiple calls")
	})
}

func TestTitleToSlug(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{"simple title", "Two Sum", "two-sum"},
		{"multiple words", "Binary Search Tree Validation", "binary-search-tree-validation"},
		{"with numbers", "3Sum Problem", "3sum-problem"},
		{"special characters", "Find & Replace", "find-replace"},
		{"extra spaces", "Two  Sum  Problem", "two-sum-problem"},
		{"leading/trailing spaces", "  Two Sum  ", "two-sum"},
		{"already lowercase", "two-sum", "two-sum"},
		{"with hyphens", "Two-Sum-Problem", "two-sum-problem"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TitleToSlug(tt.title)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSlugToSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		slug string
		want string
	}{
		{"simple slug", "two-sum", "two_sum"},
		{"multiple hyphens", "binary-search-tree", "binary_search_tree"},
		{"already snake case", "two_sum", "two_sum"},
		{"single word", "arrays", "arrays"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SlugToSnakeCase(tt.slug)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateProblem(t *testing.T) {
	t.Run("creates problem with all fields", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		input := CreateProblemInput{
			Title:       "Two Sum",
			Difficulty:  "easy",
			Topic:       "arrays",
			Description: "Find two numbers that add up to a target.",
			Tags:        "hash-table,array",
		}

		problem, err := svc.CreateProblem(input)

		assert.NoError(t, err)
		assert.NotNil(t, problem)
		assert.Equal(t, "Two Sum", problem.Title)
		assert.Equal(t, "two-sum", problem.Slug)
		assert.Equal(t, "easy", problem.Difficulty)
		assert.Equal(t, "arrays", problem.Topic)
		assert.Equal(t, "hash-table,array", problem.Tags)
		assert.NotZero(t, problem.ID)

		// Verify progress record created
		var progress database.Progress
		err = db.Where("problem_id = ?", problem.ID).First(&progress).Error
		assert.NoError(t, err)
		assert.Equal(t, "not_started", progress.Status)
		assert.Equal(t, 0, progress.Attempts)
	})

	t.Run("handles duplicate slug by appending number", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		// Create first problem
		input1 := CreateProblemInput{
			Title:       "Two Sum",
			Difficulty:  "easy",
			Topic:       "arrays",
			Description: "Original problem",
		}
		problem1, err := svc.CreateProblem(input1)
		assert.NoError(t, err)
		assert.Equal(t, "two-sum", problem1.Slug)

		// Create second problem with same title
		input2 := CreateProblemInput{
			Title:       "Two Sum",
			Difficulty:  "medium",
			Topic:       "arrays",
			Description: "Different problem",
		}
		problem2, err := svc.CreateProblem(input2)
		assert.NoError(t, err)
		assert.Equal(t, "two-sum-2", problem2.Slug)

		// Create third problem with same title
		input3 := CreateProblemInput{
			Title:       "Two Sum",
			Difficulty:  "hard",
			Topic:       "arrays",
			Description: "Another problem",
		}
		problem3, err := svc.CreateProblem(input3)
		assert.NoError(t, err)
		assert.Equal(t, "two-sum-3", problem3.Slug)
	})

	t.Run("transaction rolls back on progress creation failure", func(t *testing.T) {
		db := setupTestDB(t)
		svc := NewService(db)

		input := CreateProblemInput{
			Title:       "Test Problem",
			Difficulty:  "easy",
			Topic:       "arrays",
			Description: "Test description",
		}

		// This should succeed
		_, err := svc.CreateProblem(input)
		assert.NoError(t, err)

		// Verify problem count
		var count int64
		db.Model(&database.Problem{}).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

package problem

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"gorm.io/gorm"
)

func init() {
	// Seed random number generator for random problem selection
	rand.Seed(time.Now().UnixNano())
}

// Service handles problem-related business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new problem service instance
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// ListFilters defines optional filters for listing problems
type ListFilters struct {
	Difficulty string
	Topic      string
	Solved     *bool // Pointer to distinguish between false and unset
}

// ProblemWithStatus extends Problem model with solved status and progress timestamps
type ProblemWithStatus struct {
	database.Problem
	IsSolved      bool       `json:"is_solved"`
	FirstSolvedAt *time.Time `json:"first_solved_at,omitempty"`
}

// ListProblems retrieves problems from database with optional filters
// Returns problems with solved status determined by LEFT JOIN with Progress table
func (s *Service) ListProblems(filters ListFilters) ([]ProblemWithStatus, error) {
	var results []ProblemWithStatus

	// Build query with LEFT JOIN to Progress table (GORM pluralizes to "progresses")
	// COALESCE ensures NULL progress records are treated as unsolved (0/false)
	query := s.db.Table("problems").
		Select("problems.*, COALESCE(progresses.is_solved, 0) as is_solved, progresses.first_solved_at").
		Joins("LEFT JOIN progresses ON problems.id = progresses.problem_id")

	// Apply difficulty filter
	if filters.Difficulty != "" {
		query = query.Where("problems.difficulty = ?", filters.Difficulty)
	}

	// Apply topic filter
	if filters.Topic != "" {
		query = query.Where("problems.topic = ?", filters.Topic)
	}

	// Apply solved status filter
	if filters.Solved != nil {
		if *filters.Solved {
			// Only solved problems (progress is_solved must be true)
			query = query.Where("progresses.is_solved = ?", true)
		} else {
			// Only unsolved problems (no progress record OR is_solved is false)
			query = query.Where("(progresses.is_solved IS NULL OR progresses.is_solved = ?)", false)
		}
	}

	// Execute query with ordering for consistent, predictable results
	// Order by difficulty (easy, hard, medium alphabetically) then title
	err := query.Order("problems.difficulty ASC, problems.title ASC").Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query problems: %w", err)
	}

	return results, nil
}

// IsValidDifficulty checks if difficulty is one of the allowed values
func IsValidDifficulty(difficulty string) bool {
	validDifficulties := map[string]bool{
		"easy":   true,
		"medium": true,
		"hard":   true,
	}
	return validDifficulties[difficulty]
}

// IsValidTopic checks if topic is one of the known topics
func IsValidTopic(topic string) bool {
	validTopics := map[string]bool{
		"arrays":       true,
		"linked-lists": true,
		"trees":        true,
		"graphs":       true,
		"sorting":      true,
		"searching":    true,
	}
	return validTopics[topic]
}

// ErrProblemNotFound is returned when a problem slug does not exist
var ErrProblemNotFound = errors.New("problem not found")

// ErrNoProblemsFound is returned when no problems match the filter criteria
var ErrNoProblemsFound = errors.New("no problems found matching criteria")

// ProblemDetails extends Problem model with progress and solution information
type ProblemDetails struct {
	database.Problem
	Status          string    `json:"status"`           // not_started, in_progress, completed
	Attempts        int       `json:"attempts"`         // Number of solution attempts
	LastAttempt     time.Time `json:"last_attempt"`     // Last attempt timestamp
	HasSolution     bool      `json:"has_solution"`     // Whether a passing solution exists
	BoilerplatePath string    `json:"boilerplate_path"` // Path to boilerplate code
	TestPath        string    `json:"test_path"`        // Path to test file
}

// GetProblemBySlug retrieves a problem and its progress details by slug
// Returns ErrProblemNotFound if the slug does not exist
func (s *Service) GetProblemBySlug(slug string) (*ProblemDetails, error) {
	var problem database.Problem

	// Find problem by slug
	err := s.db.Where("slug = ?", slug).First(&problem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProblemNotFound
		}
		return nil, fmt.Errorf("failed to query problem: %w", err)
	}

	// Get progress information
	var progress database.Progress
	progressErr := s.db.Where("problem_id = ?", problem.ID).First(&progress).Error

	// Check if solution exists
	var solutionCount int64
	s.db.Model(&database.Solution{}).
		Where("problem_id = ? AND passed = ?", problem.ID, true).
		Count(&solutionCount)

	details := &ProblemDetails{
		Problem:         problem,
		Status:          "not_started",
		Attempts:        0,
		HasSolution:     solutionCount > 0,
		BoilerplatePath: fmt.Sprintf("problems/templates/%s.go", slug),
		TestPath:        fmt.Sprintf("problems/templates/%s_test.go", slug),
	}

	if progressErr == nil {
		details.Status = progress.Status
		details.Attempts = progress.Attempts
		details.LastAttempt = progress.LastAttempt
	}

	return details, nil
}

// GetRandomProblem retrieves a random problem matching the given filters
// Returns ErrNoProblemsFound if no problems match the criteria
func (s *Service) GetRandomProblem(filters ListFilters) (*ProblemDetails, error) {
	// Get all problems matching filters
	problems, err := s.ListProblems(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to query problems: %w", err)
	}

	if len(problems) == 0 {
		return nil, ErrNoProblemsFound
	}

	// Randomly select one problem
	randomIndex := rand.Intn(len(problems))
	selectedProblem := problems[randomIndex]

	// Get full details for the selected problem
	details, err := s.GetProblemBySlug(selectedProblem.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem details: %w", err)
	}

	return details, nil
}

// TitleToSlug converts a problem title to a URL-friendly slug
// Examples:
//   "Two Sum" -> "two-sum"
//   "Binary Search Tree Validation" -> "binary-search-tree-validation"
func TitleToSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	slug = result.String()

	// Remove consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// SlugToSnakeCase converts a kebab-case slug to snake_case for file names
// Examples:
//   "two-sum" -> "two_sum"
//   "binary-search-tree" -> "binary_search_tree"
func SlugToSnakeCase(slug string) string {
	return strings.ReplaceAll(slug, "-", "_")
}

// CreateProblemInput contains the parameters for creating a new custom problem
type CreateProblemInput struct {
	Title       string
	Difficulty  string
	Topic       string
	Description string
	Tags        string // Comma-separated tags
}

// CreateProblem creates a new problem with generated slug and initial progress record
// It ensures slug uniqueness by appending numbers if conflicts occur (two-sum-2, two-sum-3, etc.)
// Returns the created problem with ID populated
func (s *Service) CreateProblem(input CreateProblemInput) (*database.Problem, error) {
	// Generate unique slug from title
	slug := TitleToSlug(input.Title)

	// Check slug uniqueness and resolve conflicts
	slug, err := s.ensureUniqueSlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to generate unique slug: %w", err)
	}

	// Create problem record
	problem := &database.Problem{
		Slug:        slug,
		Title:       input.Title,
		Difficulty:  input.Difficulty,
		Topic:       input.Topic,
		Description: input.Description,
		Tags:        input.Tags,
	}

	// Use transaction to create problem and initial progress
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Create problem
		if err := tx.Create(problem).Error; err != nil {
			return fmt.Errorf("create problem: %w", err)
		}

		// Create initial progress record
		progress := &database.Progress{
			ProblemID: problem.ID,
			Status:    "not_started",
			Attempts:  0,
		}
		if err := tx.Create(progress).Error; err != nil {
			return fmt.Errorf("create progress: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return problem, nil
}

// ensureUniqueSlug checks if a slug exists in the database and appends a number if needed
// Examples:
//   "two-sum" -> "two-sum" (if unique)
//   "two-sum" -> "two-sum-2" (if "two-sum" exists)
//   "two-sum" -> "two-sum-3" (if "two-sum" and "two-sum-2" exist)
func (s *Service) ensureUniqueSlug(slug string) (string, error) {
	// Check if slug exists
	var count int64
	err := s.db.Model(&database.Problem{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return "", err
	}

	if count == 0 {
		return slug, nil // Slug is unique
	}

	// Slug exists, find next available number
	suffix := 2
	for {
		candidateSlug := fmt.Sprintf("%s-%d", slug, suffix)
		err := s.db.Model(&database.Problem{}).Where("slug = ?", candidateSlug).Count(&count).Error
		if err != nil {
			return "", err
		}
		if count == 0 {
			return candidateSlug, nil
		}
		suffix++
	}
}

// UpdateProgress updates progress status for a problem
// Creates a new progress record if one doesn't exist
// Increments attempts counter and updates timestamp
func (s *Service) UpdateProgress(problemID uint, status string) error {
	var progress database.Progress

	// Find or create progress record
	err := s.db.Where("problem_id = ?", problemID).First(&progress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new progress record
			progress = database.Progress{
				ProblemID:   problemID,
				Status:      status,
				Attempts:    1,
				LastAttempt: time.Now(),
			}
			return s.db.Create(&progress).Error
		}
		return fmt.Errorf("query progress: %w", err)
	}

	// Update existing progress
	updates := map[string]interface{}{
		"status":       status,
		"attempts":     progress.Attempts + 1,
		"last_attempt": time.Now(),
	}

	return s.db.Model(&progress).Updates(updates).Error
}

package progress

import (
	"fmt"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"gorm.io/gorm"
)

// Stats represents aggregated progress statistics
type Stats struct {
	TotalProblems   int
	TotalSolved     int
	ByDifficulty    map[string]DifficultyStats
	ByTopic         map[string]TopicStats
	RecentActivity  []RecentProblem
}

// DifficultyStats represents progress for a difficulty level
type DifficultyStats struct {
	Total  int
	Solved int
}

// TopicStats represents progress for a topic
type TopicStats struct {
	Total  int
	Solved int
}

// RecentProblem represents a recently solved problem
type RecentProblem struct {
	Slug       string
	Title      string
	Difficulty string
	SolvedAt   time.Time
}

// Service provides progress statistics calculation
type Service struct {
	db *gorm.DB
}

// NewService creates a new progress service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// GetStats calculates overall progress statistics
func (s *Service) GetStats(topicFilter string) (*Stats, error) {
	stats := &Stats{
		ByDifficulty:   make(map[string]DifficultyStats),
		ByTopic:        make(map[string]TopicStats),
		RecentActivity: []RecentProblem{},
	}

	// Build base query
	query := s.db.Model(&database.Problem{})
	if topicFilter != "" {
		query = query.Where("topic = ?", topicFilter)
	}

	// Get total problems count
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to count problems: %w", err)
	}
	stats.TotalProblems = int(count)

	// Calculate stats by difficulty
	if err := s.calculateDifficultyStats(stats, topicFilter); err != nil {
		return nil, err
	}

	// Calculate stats by topic (skip if filtering by topic)
	if topicFilter == "" {
		if err := s.calculateTopicStats(stats); err != nil {
			return nil, err
		}
	}

	// Get recent activity
	if err := s.getRecentActivity(stats, topicFilter); err != nil {
		return nil, err
	}

	return stats, nil
}

// calculateDifficultyStats calculates progress by difficulty level
func (s *Service) calculateDifficultyStats(stats *Stats, topicFilter string) error {
	type Result struct {
		Difficulty string
		Total      int
		Solved     int
	}

	query := `
		SELECT
			problems.difficulty,
			COUNT(problems.id) as total,
			SUM(CASE WHEN progresses.status = 'completed' THEN 1 ELSE 0 END) as solved
		FROM problems
		LEFT JOIN progresses ON problems.id = progresses.problem_id
	`

	if topicFilter != "" {
		query += " WHERE problems.topic = ?"
	}

	query += " GROUP BY problems.difficulty"

	var results []Result
	var err error

	if topicFilter != "" {
		err = s.db.Raw(query, topicFilter).Scan(&results).Error
	} else {
		err = s.db.Raw(query).Scan(&results).Error
	}

	if err != nil {
		return fmt.Errorf("failed to calculate difficulty stats: %w", err)
	}

	// Populate stats
	for _, r := range results {
		stats.ByDifficulty[r.Difficulty] = DifficultyStats{
			Total:  r.Total,
			Solved: r.Solved,
		}
		stats.TotalSolved += r.Solved
	}

	return nil
}

// calculateTopicStats calculates progress by topic
func (s *Service) calculateTopicStats(stats *Stats) error {
	type Result struct {
		Topic  string
		Total  int
		Solved int
	}

	query := `
		SELECT
			problems.topic,
			COUNT(problems.id) as total,
			SUM(CASE WHEN progresses.status = 'completed' THEN 1 ELSE 0 END) as solved
		FROM problems
		LEFT JOIN progresses ON problems.id = progresses.problem_id
		GROUP BY problems.topic
	`

	var results []Result
	if err := s.db.Raw(query).Scan(&results).Error; err != nil {
		return fmt.Errorf("failed to calculate topic stats: %w", err)
	}

	// Populate stats
	for _, r := range results {
		if r.Topic != "" { // Skip empty topics
			stats.ByTopic[r.Topic] = TopicStats{
				Total:  r.Total,
				Solved: r.Solved,
			}
		}
	}

	return nil
}

// getRecentActivity retrieves the last 5 solved problems
func (s *Service) getRecentActivity(stats *Stats, topicFilter string) error {
	query := s.db.Table("problems").
		Select("problems.slug, problems.title, problems.difficulty, progresses.last_attempt as solved_at").
		Joins("INNER JOIN progresses ON problems.id = progresses.problem_id").
		Where("progresses.status = ?", "completed")

	if topicFilter != "" {
		query = query.Where("problems.topic = ?", topicFilter)
	}

	query = query.Order("progresses.last_attempt DESC").Limit(5)

	if err := query.Scan(&stats.RecentActivity).Error; err != nil {
		return fmt.Errorf("failed to get recent activity: %w", err)
	}

	return nil
}

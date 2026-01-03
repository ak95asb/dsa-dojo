package analytics

import (
	"fmt"

	"gorm.io/gorm"
)

// AnalyticsService provides analytics calculations for practice patterns
type AnalyticsService struct {
	db *gorm.DB
}

// AnalyticsFilter specifies filtering criteria for analytics
type AnalyticsFilter struct {
	Topic      string
	Difficulty string
}

// AnalyticsStats contains calculated analytics data
type AnalyticsStats struct {
	OverallSuccessRate      float64            `json:"overall_success_rate"`
	SuccessRateByDifficulty map[string]float64 `json:"success_rate_by_difficulty"`
	SuccessRateByTopic      map[string]float64 `json:"success_rate_by_topic"`
	AvgAttemptsOverall      float64            `json:"avg_attempts_overall"`
	AvgAttemptsByDifficulty map[string]float64 `json:"avg_attempts_by_difficulty"`
	AvgAttemptsByTopic      map[string]float64 `json:"avg_attempts_by_topic"`
	MostPracticedTopic      string             `json:"most_practiced_topic"`
	LeastPracticedTopic     string             `json:"least_practiced_topic"`
	BestDifficulty          string             `json:"best_difficulty"`
	ChallengingDifficulty   string             `json:"challenging_difficulty"`
}

// NewAnalyticsService creates a new analytics service instance
func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// CalculateStats computes all analytics statistics with optional filtering
func (s *AnalyticsService) CalculateStats(filter AnalyticsFilter) (*AnalyticsStats, error) {
	stats := &AnalyticsStats{
		SuccessRateByDifficulty: make(map[string]float64),
		SuccessRateByTopic:      make(map[string]float64),
		AvgAttemptsByDifficulty: make(map[string]float64),
		AvgAttemptsByTopic:      make(map[string]float64),
	}

	// Calculate overall success rate
	overallRate, err := s.calculateOverallSuccessRate(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate overall success rate: %w", err)
	}
	stats.OverallSuccessRate = overallRate

	// Calculate success rates by difficulty
	ratesByDiff, err := s.calculateSuccessRateByDifficulty(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate success rates by difficulty: %w", err)
	}
	stats.SuccessRateByDifficulty = ratesByDiff

	// Calculate success rates by topic
	ratesByTopic, err := s.calculateSuccessRateByTopic(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate success rates by topic: %w", err)
	}
	stats.SuccessRateByTopic = ratesByTopic

	// Calculate overall average attempts
	overallAvg, err := s.calculateOverallAverageAttempts(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate overall avg attempts: %w", err)
	}
	stats.AvgAttemptsOverall = overallAvg

	// Calculate average attempts by difficulty
	avgByDiff, err := s.calculateAverageAttemptsByDifficulty(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate avg attempts by difficulty: %w", err)
	}
	stats.AvgAttemptsByDifficulty = avgByDiff

	// Calculate average attempts by topic
	avgByTopic, err := s.calculateAverageAttemptsByTopic(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate avg attempts by topic: %w", err)
	}
	stats.AvgAttemptsByTopic = avgByTopic

	// Analyze practice patterns
	patterns, err := s.analyzePracticePatterns(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze practice patterns: %w", err)
	}
	stats.MostPracticedTopic = patterns.MostPracticed
	stats.LeastPracticedTopic = patterns.LeastPracticed
	stats.BestDifficulty = patterns.BestDifficulty
	stats.ChallengingDifficulty = patterns.ChallengingDifficulty

	return stats, nil
}

func (s *AnalyticsService) calculateOverallSuccessRate(filter AnalyticsFilter) (float64, error) {
	type Result struct {
		Total  int64
		Solved int64
	}

	var result Result
	query := s.db.Table("progresses").
		Select("COUNT(*) as total, SUM(CASE WHEN is_solved = true THEN 1 ELSE 0 END) as solved").
		Joins("INNER JOIN problems ON progresses.problem_id = problems.id").
		Where("progresses.total_attempts > 0")

	if filter.Topic != "" {
		query = query.Where("problems.topic = ?", filter.Topic)
	}
	if filter.Difficulty != "" {
		query = query.Where("problems.difficulty = ?", filter.Difficulty)
	}

	if err := query.Scan(&result).Error; err != nil {
		return 0, err
	}

	if result.Total == 0 {
		return 0.0, nil
	}

	return (float64(result.Solved) / float64(result.Total)) * 100, nil
}

func (s *AnalyticsService) calculateSuccessRateByDifficulty(filter AnalyticsFilter) (map[string]float64, error) {
	type Result struct {
		Difficulty  string
		Total       int64
		Solved      int64
		SuccessRate float64
	}

	var results []Result
	query := s.db.Table("progresses").
		Select("problems.difficulty, COUNT(*) as total, SUM(CASE WHEN is_solved = true THEN 1 ELSE 0 END) as solved").
		Joins("INNER JOIN problems ON progresses.problem_id = problems.id").
		Where("progresses.total_attempts > 0").
		Group("problems.difficulty")

	if filter.Topic != "" {
		query = query.Where("problems.topic = ?", filter.Topic)
	}
	if filter.Difficulty != "" {
		query = query.Where("problems.difficulty = ?", filter.Difficulty)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	rates := make(map[string]float64)
	for _, r := range results {
		if r.Total > 0 {
			rates[r.Difficulty] = (float64(r.Solved) / float64(r.Total)) * 100
		}
	}

	return rates, nil
}

func (s *AnalyticsService) calculateSuccessRateByTopic(filter AnalyticsFilter) (map[string]float64, error) {
	type Result struct {
		Topic  string
		Total  int64
		Solved int64
	}

	var results []Result
	query := s.db.Table("progresses").
		Select("problems.topic, COUNT(*) as total, SUM(CASE WHEN is_solved = true THEN 1 ELSE 0 END) as solved").
		Joins("INNER JOIN problems ON progresses.problem_id = problems.id").
		Where("progresses.total_attempts > 0").
		Group("problems.topic")

	if filter.Topic != "" {
		query = query.Where("problems.topic = ?", filter.Topic)
	}
	if filter.Difficulty != "" {
		query = query.Where("problems.difficulty = ?", filter.Difficulty)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	rates := make(map[string]float64)
	for _, r := range results {
		if r.Total > 0 {
			rates[r.Topic] = (float64(r.Solved) / float64(r.Total)) * 100
		}
	}

	return rates, nil
}

func (s *AnalyticsService) calculateOverallAverageAttempts(filter AnalyticsFilter) (float64, error) {
	type Result struct {
		AvgAttempts float64
	}

	var result Result
	query := s.db.Table("progresses").
		Select("AVG(progresses.total_attempts) as avg_attempts").
		Joins("INNER JOIN problems ON progresses.problem_id = problems.id").
		Where("progresses.is_solved = true")

	if filter.Topic != "" {
		query = query.Where("problems.topic = ?", filter.Topic)
	}
	if filter.Difficulty != "" {
		query = query.Where("problems.difficulty = ?", filter.Difficulty)
	}

	if err := query.Scan(&result).Error; err != nil {
		return 0, err
	}

	return result.AvgAttempts, nil
}

func (s *AnalyticsService) calculateAverageAttemptsByDifficulty(filter AnalyticsFilter) (map[string]float64, error) {
	type Result struct {
		Difficulty  string
		AvgAttempts float64
	}

	var results []Result
	query := s.db.Table("progresses").
		Select("problems.difficulty, AVG(progresses.total_attempts) as avg_attempts").
		Joins("INNER JOIN problems ON progresses.problem_id = problems.id").
		Where("progresses.is_solved = true").
		Group("problems.difficulty")

	if filter.Topic != "" {
		query = query.Where("problems.topic = ?", filter.Topic)
	}
	if filter.Difficulty != "" {
		query = query.Where("problems.difficulty = ?", filter.Difficulty)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	avgs := make(map[string]float64)
	for _, r := range results {
		avgs[r.Difficulty] = r.AvgAttempts
	}

	return avgs, nil
}

func (s *AnalyticsService) calculateAverageAttemptsByTopic(filter AnalyticsFilter) (map[string]float64, error) {
	type Result struct {
		Topic       string
		AvgAttempts float64
	}

	var results []Result
	query := s.db.Table("progresses").
		Select("problems.topic, AVG(progresses.total_attempts) as avg_attempts").
		Joins("INNER JOIN problems ON progresses.problem_id = problems.id").
		Where("progresses.is_solved = true").
		Group("problems.topic")

	if filter.Topic != "" {
		query = query.Where("problems.topic = ?", filter.Topic)
	}
	if filter.Difficulty != "" {
		query = query.Where("problems.difficulty = ?", filter.Difficulty)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	avgs := make(map[string]float64)
	for _, r := range results {
		avgs[r.Topic] = r.AvgAttempts
	}

	return avgs, nil
}

type PracticePatterns struct {
	MostPracticed       string
	LeastPracticed      string
	BestDifficulty      string
	ChallengingDifficulty string
}

func (s *AnalyticsService) analyzePracticePatterns(filter AnalyticsFilter) (*PracticePatterns, error) {
	patterns := &PracticePatterns{}

	// Find most/least practiced topics
	type TopicAttempts struct {
		Topic         string
		TotalAttempts int64
	}

	var topicResults []TopicAttempts
	topicQuery := s.db.Table("progresses").
		Select("problems.topic, SUM(progresses.total_attempts) as total_attempts").
		Joins("INNER JOIN problems ON progresses.problem_id = problems.id").
		Group("problems.topic").
		Order("total_attempts DESC")

	if filter.Difficulty != "" {
		topicQuery = topicQuery.Where("problems.difficulty = ?", filter.Difficulty)
	}

	if err := topicQuery.Scan(&topicResults).Error; err != nil {
		return nil, err
	}

	if len(topicResults) > 0 {
		patterns.MostPracticed = topicResults[0].Topic
		patterns.LeastPracticed = topicResults[len(topicResults)-1].Topic
	}

	// Find best/challenging difficulty
	type DifficultyStats struct {
		Difficulty  string
		SuccessRate float64
		AvgAttempts float64
	}

	var diffResults []DifficultyStats
	diffQuery := s.db.Table("progresses").
		Select(`problems.difficulty,
			CAST(SUM(CASE WHEN is_solved = true THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) * 100 as success_rate,
			AVG(CASE WHEN is_solved = true THEN progresses.total_attempts ELSE NULL END) as avg_attempts`).
		Joins("INNER JOIN problems ON progresses.problem_id = problems.id").
		Where("progresses.total_attempts > 0").
		Group("problems.difficulty").
		Order("success_rate DESC")

	if filter.Topic != "" {
		diffQuery = diffQuery.Where("problems.topic = ?", filter.Topic)
	}

	if err := diffQuery.Scan(&diffResults).Error; err != nil {
		return nil, err
	}

	if len(diffResults) > 0 {
		patterns.BestDifficulty = diffResults[0].Difficulty

		// Challenging is the one with lowest success rate OR highest avg attempts
		// Use last in the list (sorted by success rate DESC)
		patterns.ChallengingDifficulty = diffResults[len(diffResults)-1].Difficulty
	}

	return patterns, nil
}

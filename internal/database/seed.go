package database

import (
	"fmt"

	"github.com/ak95asb/dsa-dojo/problems"
	"gorm.io/gorm"
)

// SeedProblems populates the database with the initial problem library
// Returns the number of problems seeded and any error encountered
// This function is idempotent - it will not create duplicates if run multiple times
func SeedProblems(db *gorm.DB) (int, error) {
	seedData := problems.SeedData()
	seededCount := 0

	for _, seed := range seedData {
		// Check if problem already exists by slug
		var existing Problem
		result := db.Where("slug = ?", seed.Slug).First(&existing)

		// Skip if already exists
		if result.Error == nil {
			continue
		}

		// Only proceed if error is "record not found"
		if result.Error != gorm.ErrRecordNotFound {
			return seededCount, fmt.Errorf("failed to check existing problem '%s': %w", seed.Slug, result.Error)
		}

		// Create new problem
		problem := Problem{
			Slug:        seed.Slug,
			Title:       seed.Title,
			Description: seed.Description,
			Difficulty:  seed.Difficulty,
			Topic:       seed.Topic,
		}

		if err := db.Create(&problem).Error; err != nil {
			return seededCount, fmt.Errorf("failed to seed problem '%s': %w", seed.Slug, err)
		}

		seededCount++
	}

	return seededCount, nil
}

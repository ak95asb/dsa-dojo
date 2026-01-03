package export

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Integration tests verify end-to-end export workflows

func setupIntegrationDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&database.Problem{}, &database.Progress{}, &database.Solution{})
	require.NoError(t, err)

	return db
}

func seedLargeDataset(db *gorm.DB, t *testing.T, count int) {
	topics := []string{"arrays", "strings", "trees", "graphs", "dynamic-programming"}
	difficulties := []string{"easy", "medium", "hard"}

	for i := 1; i <= count; i++ {
		problem := &database.Problem{
			Slug:       string(rune('a'+(i%26))) + "-problem-" + string(rune(i)),
			Title:      "Problem " + string(rune(i)),
			Difficulty: difficulties[i%3],
			Topic:      topics[i%5],
		}
		require.NoError(t, db.Create(problem).Error)

		progress := &database.Progress{
			ProblemID:     uint(i),
			IsSolved:      i%2 == 0,
			TotalAttempts: (i % 5) + 1,
		}
		require.NoError(t, db.Create(progress).Error)

		if i%2 == 0 {
			solution := &database.Solution{
				ProblemID:   uint(i),
				Status:      "Passed",
				TestsPassed: 5,
				TestsTotal:  5,
				Passed:      true,
			}
			require.NoError(t, db.Create(solution).Error)
		}
	}
}

func TestIntegration_CompleteExportWorkflow(t *testing.T) {
	db := setupIntegrationDB(t)
	service := NewService(db)

	t.Run("complete JSON export workflow", func(t *testing.T) {
		// Seed test data
		seedTestData(db, t)

		// Export to buffer
		var buf bytes.Buffer
		err := service.ExportToJSON(ExportFilter{}, &buf)
		require.NoError(t, err)

		// Verify JSON is valid
		var data ExportData
		err = json.Unmarshal(buf.Bytes(), &data)
		require.NoError(t, err)

		// Verify all components
		assert.Equal(t, "1.0", data.Version)
		assert.NotEmpty(t, data.ExportedAt)
		assert.Greater(t, data.Summary.TotalProblems, 0)
		assert.NotEmpty(t, data.Problems)

		// Verify data integrity
		for _, problem := range data.Problems {
			assert.NotEmpty(t, problem.Slug)
			assert.NotEmpty(t, problem.Title)
			assert.NotEmpty(t, problem.Difficulty)
			assert.NotEmpty(t, problem.Topic)
		}
	})

	t.Run("complete CSV export workflow", func(t *testing.T) {
		// Export to buffer
		var buf bytes.Buffer
		err := service.ExportToCSV(ExportFilter{}, &buf)
		require.NoError(t, err)

		// Parse CSV
		reader := csv.NewReader(&buf)
		records, err := reader.ReadAll()
		require.NoError(t, err)

		// Verify structure
		assert.Greater(t, len(records), 1) // Header + data
		assert.Len(t, records[0], 8)       // 8 columns

		// Verify data consistency
		for i, record := range records[1:] {
			assert.Len(t, record, 8, "Row %d should have 8 columns", i)
		}
	})
}

func TestIntegration_FilteredExports(t *testing.T) {
	db := setupIntegrationDB(t)
	seedTestData(db, t)
	service := NewService(db)

	t.Run("filter by difficulty - JSON", func(t *testing.T) {
		var buf bytes.Buffer
		filter := ExportFilter{Difficulty: "easy"}
		err := service.ExportToJSON(filter, &buf)
		require.NoError(t, err)

		var data ExportData
		json.Unmarshal(buf.Bytes(), &data)

		// All problems should be easy
		for _, problem := range data.Problems {
			assert.Equal(t, "easy", problem.Difficulty)
		}
	})

	t.Run("filter by topic - JSON", func(t *testing.T) {
		var buf bytes.Buffer
		filter := ExportFilter{Topic: "arrays"}
		err := service.ExportToJSON(filter, &buf)
		require.NoError(t, err)

		var data ExportData
		json.Unmarshal(buf.Bytes(), &data)

		// All problems should be arrays
		for _, problem := range data.Problems {
			assert.Equal(t, "arrays", problem.Topic)
		}
	})

	t.Run("combined filters - CSV", func(t *testing.T) {
		var buf bytes.Buffer
		filter := ExportFilter{Difficulty: "easy", Topic: "strings"}
		err := service.ExportToCSV(filter, &buf)
		require.NoError(t, err)

		reader := csv.NewReader(&buf)
		records, err := reader.ReadAll()
		require.NoError(t, err)

		// Should have specific filtered data
		for _, record := range records[1:] {
			assert.Equal(t, "easy", record[2])    // Difficulty column
			assert.Equal(t, "strings", record[3]) // Topic column
		}
	})
}

func TestIntegration_Performance(t *testing.T) {
	db := setupIntegrationDB(t)
	service := NewService(db)

	t.Run("JSON export completes in <5s with 100+ problems", func(t *testing.T) {
		// Seed large dataset
		seedLargeDataset(db, t, 100)

		// Measure export time
		start := time.Now()
		var buf bytes.Buffer
		err := service.ExportToJSON(ExportFilter{}, &buf)
		elapsed := time.Since(start)

		require.NoError(t, err)
		assert.Less(t, elapsed.Milliseconds(), int64(5000), "JSON export should complete in <5s")

		// Verify data was exported
		var data ExportData
		json.Unmarshal(buf.Bytes(), &data)
		assert.Equal(t, 100, data.Summary.TotalProblems)
	})

	t.Run("CSV export completes in <5s with 100+ problems", func(t *testing.T) {
		start := time.Now()
		var buf bytes.Buffer
		err := service.ExportToCSV(ExportFilter{}, &buf)
		elapsed := time.Since(start)

		require.NoError(t, err)
		assert.Less(t, elapsed.Milliseconds(), int64(5000), "CSV export should complete in <5s")

		// Verify data was exported
		reader := csv.NewReader(&buf)
		records, err := reader.ReadAll()
		require.NoError(t, err)
		assert.Equal(t, 101, len(records)) // 100 problems + header
	})
}

func TestIntegration_FileOutput(t *testing.T) {
	db := setupIntegrationDB(t)
	seedTestData(db, t)
	service := NewService(db)

	t.Run("export to JSON file", func(t *testing.T) {
		filename := "test_export.json"
		defer os.Remove(filename)

		file, err := os.Create(filename)
		require.NoError(t, err)
		defer file.Close()

		err = service.ExportToJSON(ExportFilter{}, file)
		require.NoError(t, err)

		// Close file to flush
		file.Close()

		// Read file back
		content, err := os.ReadFile(filename)
		require.NoError(t, err)

		var data ExportData
		err = json.Unmarshal(content, &data)
		require.NoError(t, err)
		assert.Greater(t, data.Summary.TotalProblems, 0)
	})

	t.Run("export to CSV file", func(t *testing.T) {
		filename := "test_export.csv"
		defer os.Remove(filename)

		file, err := os.Create(filename)
		require.NoError(t, err)
		defer file.Close()

		err = service.ExportToCSV(ExportFilter{}, file)
		require.NoError(t, err)

		// Close file to flush
		file.Close()

		// Read file back
		content, err := os.ReadFile(filename)
		require.NoError(t, err)

		reader := csv.NewReader(bytes.NewReader(content))
		records, err := reader.ReadAll()
		require.NoError(t, err)
		assert.Greater(t, len(records), 1)
	})
}

func TestIntegration_DataIntegrity(t *testing.T) {
	db := setupIntegrationDB(t)
	service := NewService(db)

	t.Run("export doesn't modify database", func(t *testing.T) {
		// Seed data
		seedTestData(db, t)

		// Count records before
		var countBefore int64
		db.Model(&database.Problem{}).Count(&countBefore)

		// Export
		var buf bytes.Buffer
		err := service.ExportToJSON(ExportFilter{}, &buf)
		require.NoError(t, err)

		// Count records after
		var countAfter int64
		db.Model(&database.Problem{}).Count(&countAfter)

		// Should be unchanged
		assert.Equal(t, countBefore, countAfter)
	})

	t.Run("multiple exports produce consistent results", func(t *testing.T) {
		var buf1, buf2 bytes.Buffer

		// First export
		err := service.ExportToJSON(ExportFilter{}, &buf1)
		require.NoError(t, err)

		// Second export
		err = service.ExportToJSON(ExportFilter{}, &buf2)
		require.NoError(t, err)

		// Parse both
		var data1, data2 ExportData
		json.Unmarshal(buf1.Bytes(), &data1)
		json.Unmarshal(buf2.Bytes(), &data2)

		// Summary should be identical (excluding timestamp)
		assert.Equal(t, data1.Summary.TotalProblems, data2.Summary.TotalProblems)
		assert.Equal(t, data1.Summary.ProblemsSolved, data2.Summary.ProblemsSolved)
		assert.Equal(t, data1.Summary.OverallSuccessRate, data2.Summary.OverallSuccessRate)
		assert.Len(t, data1.Problems, len(data2.Problems))
	})
}

func TestIntegration_EdgeCases(t *testing.T) {
	db := setupIntegrationDB(t)
	service := NewService(db)

	t.Run("export with no data", func(t *testing.T) {
		var buf bytes.Buffer
		err := service.ExportToJSON(ExportFilter{}, &buf)
		require.NoError(t, err)

		var data ExportData
		json.Unmarshal(buf.Bytes(), &data)

		assert.Equal(t, 0, data.Summary.TotalProblems)
		assert.Empty(t, data.Problems)
	})

	t.Run("export with invalid filter returns empty", func(t *testing.T) {
		seedTestData(db, t)

		var buf bytes.Buffer
		filter := ExportFilter{Topic: "nonexistent-topic"}
		err := service.ExportToJSON(filter, &buf)
		require.NoError(t, err)

		var data ExportData
		json.Unmarshal(buf.Bytes(), &data)

		assert.Equal(t, 0, len(data.Problems))
	})
}

func TestIntegration_SolutionHistory(t *testing.T) {
	db := setupIntegrationDB(t)
	seedTestData(db, t)
	service := NewService(db)

	t.Run("JSON export includes full solution history", func(t *testing.T) {
		var buf bytes.Buffer
		err := service.ExportToJSON(ExportFilter{}, &buf)
		require.NoError(t, err)

		var data ExportData
		json.Unmarshal(buf.Bytes(), &data)

		// Find two-sum problem
		var twoSum *ProblemExport
		for i := range data.Problems {
			if data.Problems[i].Slug == "two-sum" {
				twoSum = &data.Problems[i]
				break
			}
		}

		require.NotNil(t, twoSum)
		assert.Len(t, twoSum.Solutions, 2) // Should have 2 solutions

		// Verify solution details
		assert.Equal(t, "Failed", twoSum.Solutions[0].Status)
		assert.Equal(t, 3, twoSum.Solutions[0].TestsPassed)
		assert.Equal(t, 5, twoSum.Solutions[0].TestsTotal)

		assert.Equal(t, "Passed", twoSum.Solutions[1].Status)
		assert.Equal(t, 5, twoSum.Solutions[1].TestsPassed)
		assert.Equal(t, 5, twoSum.Solutions[1].TestsTotal)
	})
}

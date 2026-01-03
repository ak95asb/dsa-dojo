package export

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
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
	// Create test problems
	problems := []database.Problem{
		{Slug: "two-sum", Title: "Two Sum", Difficulty: "easy", Topic: "arrays"},
		{Slug: "valid-parentheses", Title: "Valid Parentheses", Difficulty: "easy", Topic: "strings"},
		{Slug: "merge-sort", Title: "Merge Sort", Difficulty: "medium", Topic: "sorting"},
		{Slug: "binary-tree", Title: "Binary Tree", Difficulty: "hard", Topic: "trees"},
	}

	for _, p := range problems {
		require.NoError(t, db.Create(&p).Error)
	}

	// Create progress records
	progressRecords := []database.Progress{
		{ProblemID: 1, IsSolved: true, TotalAttempts: 2},
		{ProblemID: 2, IsSolved: true, TotalAttempts: 1},
		{ProblemID: 3, IsSolved: false, TotalAttempts: 3},
		{ProblemID: 4, IsSolved: false, TotalAttempts: 5},
	}

	for _, p := range progressRecords {
		require.NoError(t, db.Create(&p).Error)
	}

	// Create solution records
	solutions := []database.Solution{
		{ProblemID: 1, Status: "Failed", TestsPassed: 3, TestsTotal: 5},
		{ProblemID: 1, Status: "Passed", TestsPassed: 5, TestsTotal: 5, Passed: true},
		{ProblemID: 2, Status: "Passed", TestsPassed: 3, TestsTotal: 3, Passed: true},
	}

	for _, s := range solutions {
		require.NoError(t, db.Create(&s).Error)
	}
}

func TestNewService(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	assert.NotNil(t, service)
	assert.NotNil(t, service.db)
}

func TestExportToJSON_ValidSchema(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewService(db)

	var buf bytes.Buffer
	err := service.ExportToJSON(ExportFilter{}, &buf)

	require.NoError(t, err)

	// Validate JSON structure
	var data ExportData
	err = json.Unmarshal(buf.Bytes(), &data)
	require.NoError(t, err)

	// Verify metadata
	assert.Equal(t, "1.0", data.Version)
	assert.NotEmpty(t, data.ExportedAt)

	// Verify summary
	assert.Equal(t, 4, data.Summary.TotalProblems)
	assert.Equal(t, 2, data.Summary.ProblemsSolved)
	assert.InDelta(t, 50.0, data.Summary.OverallSuccessRate, 0.1)

	// Verify problems data
	assert.Len(t, data.Problems, 4)
}

func TestExportToJSON_IncludesSolutions(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewService(db)

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
	assert.Equal(t, "Two Sum", twoSum.Title)
	assert.Len(t, twoSum.Solutions, 2) // Two solutions submitted
	assert.True(t, twoSum.Progress.IsSolved)
	assert.Equal(t, 2, twoSum.Progress.TotalAttempts)
}

func TestExportToJSON_FilterByDifficulty(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewService(db)

	var buf bytes.Buffer
	filter := ExportFilter{Difficulty: "easy"}
	err := service.ExportToJSON(filter, &buf)

	require.NoError(t, err)

	var data ExportData
	json.Unmarshal(buf.Bytes(), &data)

	// Should only have easy problems
	assert.Len(t, data.Problems, 2)
	for _, problem := range data.Problems {
		assert.Equal(t, "easy", problem.Difficulty)
	}
}

func TestExportToJSON_FilterByTopic(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewService(db)

	var buf bytes.Buffer
	filter := ExportFilter{Topic: "arrays"}
	err := service.ExportToJSON(filter, &buf)

	require.NoError(t, err)

	var data ExportData
	json.Unmarshal(buf.Bytes(), &data)

	// Should only have arrays problems
	assert.Len(t, data.Problems, 1)
	assert.Equal(t, "arrays", data.Problems[0].Topic)
}

func TestExportToJSON_CombinedFilters(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewService(db)

	var buf bytes.Buffer
	filter := ExportFilter{Difficulty: "easy", Topic: "strings"}
	err := service.ExportToJSON(filter, &buf)

	require.NoError(t, err)

	var data ExportData
	json.Unmarshal(buf.Bytes(), &data)

	// Should only have easy strings problems
	assert.Len(t, data.Problems, 1)
	assert.Equal(t, "valid-parentheses", data.Problems[0].Slug)
}

func TestExportToJSON_EmptyDatabase(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	var buf bytes.Buffer
	err := service.ExportToJSON(ExportFilter{}, &buf)

	require.NoError(t, err)

	var data ExportData
	json.Unmarshal(buf.Bytes(), &data)

	assert.Equal(t, 0, data.Summary.TotalProblems)
	assert.Empty(t, data.Problems)
}

func TestExportToCSV_ValidFormat(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewService(db)

	var buf bytes.Buffer
	err := service.ExportToCSV(ExportFilter{}, &buf)

	require.NoError(t, err)

	// Parse CSV
	reader := csv.NewReader(&buf)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Verify header
	assert.Equal(t, "Slug", records[0][0])
	assert.Equal(t, "Title", records[0][1])
	assert.Equal(t, "Difficulty", records[0][2])
	assert.Equal(t, "Topic", records[0][3])
	assert.Equal(t, "IsSolved", records[0][4])
	assert.Equal(t, "TotalAttempts", records[0][5])
	assert.Equal(t, "FirstSolvedAt", records[0][6])
	assert.Equal(t, "LastAttemptedAt", records[0][7])

	// Verify data rows (4 problems + 1 header)
	assert.Len(t, records, 5)
}

func TestExportToCSV_DataAccuracy(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewService(db)

	var buf bytes.Buffer
	err := service.ExportToCSV(ExportFilter{}, &buf)

	require.NoError(t, err)

	reader := csv.NewReader(&buf)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Find two-sum row
	var twoSumRow []string
	for _, record := range records[1:] {
		if record[0] == "two-sum" {
			twoSumRow = record
			break
		}
	}

	require.NotNil(t, twoSumRow)
	assert.Equal(t, "Two Sum", twoSumRow[1])
	assert.Equal(t, "easy", twoSumRow[2])
	assert.Equal(t, "arrays", twoSumRow[3])
	assert.Equal(t, "true", twoSumRow[4])
	assert.Equal(t, "2", twoSumRow[5])
}

func TestExportToCSV_FilterByDifficulty(t *testing.T) {
	db := setupTestDB(t)
	seedTestData(db, t)
	service := NewService(db)

	var buf bytes.Buffer
	filter := ExportFilter{Difficulty: "medium"}
	err := service.ExportToCSV(filter, &buf)

	require.NoError(t, err)

	reader := csv.NewReader(&buf)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Should have header + 1 medium problem
	assert.Len(t, records, 2)
	assert.Equal(t, "merge-sort", records[1][0])
}

func TestExportToCSV_SpecialCharacters(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	// Create problem with special characters
	problem := &database.Problem{
		Slug:       "test-special",
		Title:      `Problem with "quotes" and, commas`,
		Difficulty: "easy",
		Topic:      "arrays",
	}
	require.NoError(t, db.Create(problem).Error)

	progress := &database.Progress{
		ProblemID:     1,
		IsSolved:      false,
		TotalAttempts: 1,
	}
	require.NoError(t, db.Create(progress).Error)

	var buf bytes.Buffer
	err := service.ExportToCSV(ExportFilter{}, &buf)

	require.NoError(t, err)

	// Verify CSV is properly escaped
	reader := csv.NewReader(&buf)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Find the test problem
	for _, record := range records[1:] {
		if record[0] == "test-special" {
			assert.Equal(t, `Problem with "quotes" and, commas`, record[1])
		}
	}
}

func TestExportToCSV_EmptyDatabase(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	var buf bytes.Buffer
	err := service.ExportToCSV(ExportFilter{}, &buf)

	require.NoError(t, err)

	reader := csv.NewReader(&buf)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Should only have header
	assert.Len(t, records, 1)
}

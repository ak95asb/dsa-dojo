package testing

import (
	"testing"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestParseTestResults(t *testing.T) {
	executor := NewExecutor()

	tests := []struct {
		name            string
		output          string
		expectedPassed  int
		expectedTotal   int
		expectedAllPass bool
		expectedFails   int
	}{
		{
			name: "all tests pass",
			output: `=== RUN   TestTwoSum
--- PASS: TestTwoSum (0.00s)
=== RUN   TestBinarySearch
--- PASS: TestBinarySearch (0.00s)
PASS
ok  	github.com/ak95asb/dsa-dojo/problems	0.123s`,
			expectedPassed:  2,
			expectedTotal:   2,
			expectedAllPass: true,
			expectedFails:   0,
		},
		{
			name: "some tests fail",
			output: `=== RUN   TestTwoSum
--- PASS: TestTwoSum (0.00s)
=== RUN   TestBinarySearch
--- FAIL: TestBinarySearch (0.00s)
    Error: Not equal:
           expected: 0
           actual  : -1
FAIL
FAIL	github.com/ak95asb/dsa-dojo/problems	0.456s`,
			expectedPassed:  1,
			expectedTotal:   2,
			expectedAllPass: false,
			expectedFails:   1,
		},
		{
			name: "all tests fail",
			output: `=== RUN   TestTwoSum
--- FAIL: TestTwoSum (0.00s)
    Error: Not equal
=== RUN   TestBinarySearch
--- FAIL: TestBinarySearch (0.00s)
    Error: Not equal
FAIL
FAIL	github.com/ak95asb/dsa-dojo/problems	0.789s`,
			expectedPassed:  0,
			expectedTotal:   2,
			expectedAllPass: false,
			expectedFails:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &TestResult{
				AllPassed: tt.expectedAllPass,
			}
			executor.parseTestResults(result, tt.output)

			assert.Equal(t, tt.expectedPassed, result.PassedCount, "Passed count mismatch")
			assert.Equal(t, tt.expectedTotal, result.TotalCount, "Total count mismatch")
			assert.Equal(t, tt.expectedFails, len(result.FailedTests), "Failed tests count mismatch")
		})
	}
}

func TestShouldUseColor(t *testing.T) {
	tests := []struct {
		name     string
		noColor  string
		expected bool
	}{
		{
			name:     "NO_COLOR not set",
			noColor:  "",
			expected: false, // Will be false in test environment (no TTY)
		},
		{
			name:     "NO_COLOR set to 1",
			noColor:  "1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.noColor != "" {
				t.Setenv("NO_COLOR", tt.noColor)
			}

			result := shouldUseColor()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatterDisplay(t *testing.T) {
	formatter := NewFormatter()

	t.Run("displays success for all passing tests", func(t *testing.T) {
		result := &TestResult{
			AllPassed:   true,
			PassedCount: 5,
			TotalCount:  5,
			FailedTests: []FailedTest{},
		}

		// This should not panic
		formatter.Display(result)
	})

	t.Run("displays failure for failing tests", func(t *testing.T) {
		result := &TestResult{
			AllPassed:   false,
			PassedCount: 2,
			TotalCount:  5,
			FailedTests: []FailedTest{
				{
					Name:     "TestFoo",
					Expected: "5",
					Actual:   "3",
					Message:  "Not equal",
				},
			},
		}

		// This should not panic
		formatter.Display(result)
	})
}

func TestServiceRecordSolution(t *testing.T) {
	// Skip database tests in unit tests - covered by integration tests
	// This test verifies the service structure compiles correctly

	t.Run("service creation", func(t *testing.T) {
		// Initialize test database
		db, err := database.Initialize()
		if err != nil {
			t.Skip("Database not available for unit test")
			return
		}
		defer func() {
			sqlDB, _ := db.DB()
			sqlDB.Close()
		}()

		// Seed test data
		database.SeedProblems(db)

		service := NewService(db)
		assert.NotNil(t, service)

		// Test with a valid problem ID
		result := &TestResult{
			AllPassed:   true,
			PassedCount: 5,
			TotalCount:  5,
		}

		err = service.RecordSolution(1, result)
		assert.NoError(t, err)
	})
}

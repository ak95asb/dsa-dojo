package testgen

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONImporter handles importing test cases from JSON files
type JSONImporter struct{}

// NewJSONImporter creates a new JSON importer
func NewJSONImporter() *JSONImporter {
	return &JSONImporter{}
}

// JSONTestFile represents the structure of the JSON test file
type JSONTestFile struct {
	Tests []JSONTestCase `json:"tests"`
}

// JSONTestCase represents a single test case in JSON format
type JSONTestCase struct {
	Name     string        `json:"name"`
	Inputs   []interface{} `json:"inputs"`
	Expected interface{}   `json:"expected"`
}

// ImportFromFile reads and parses test cases from a JSON file
func (j *JSONImporter) ImportFromFile(filePath string) ([]*TestCase, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var jsonFile JSONTestFile
	if err := json.Unmarshal(data, &jsonFile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate schema
	if err := j.validate(&jsonFile); err != nil {
		return nil, fmt.Errorf("invalid JSON schema: %w", err)
	}

	// Convert to internal representation
	testCases := make([]*TestCase, len(jsonFile.Tests))
	for i, jsonTest := range jsonFile.Tests {
		testCases[i] = &TestCase{
			Name:     jsonTest.Name,
			Inputs:   jsonTest.Inputs,
			Expected: jsonTest.Expected,
		}
	}

	fmt.Printf("📦 Imported %d test case(s) from %s\n", len(testCases), filePath)
	return testCases, nil
}

// validate checks that the JSON structure is valid
func (j *JSONImporter) validate(jsonFile *JSONTestFile) error {
	if len(jsonFile.Tests) == 0 {
		return fmt.Errorf("no test cases found in JSON file")
	}

	for i, test := range jsonFile.Tests {
		if test.Name == "" {
			return fmt.Errorf("test case %d is missing 'name' field", i+1)
		}
		if test.Inputs == nil {
			return fmt.Errorf("test case '%s' is missing 'inputs' field", test.Name)
		}
		if test.Expected == nil {
			return fmt.Errorf("test case '%s' is missing 'expected' field", test.Name)
		}
	}

	return nil
}

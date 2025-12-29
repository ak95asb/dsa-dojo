package testgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONImporter_ImportFromFile_ValidFile(t *testing.T) {
	// Create temporary JSON file
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "tests.json")

	jsonContent := `{
  "tests": [
    {
      "name": "test case 1",
      "inputs": [1, 2, 3],
      "expected": 6
    },
    {
      "name": "test case 2",
      "inputs": [4, 5],
      "expected": 9
    }
  ]
}`

	err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
	assert.NoError(t, err)

	// Import test cases
	importer := NewJSONImporter()
	testCases, err := importer.ImportFromFile(jsonFile)

	assert.NoError(t, err)
	assert.Len(t, testCases, 2)
	assert.Equal(t, "test case 1", testCases[0].Name)
	assert.Equal(t, []interface{}{float64(1), float64(2), float64(3)}, testCases[0].Inputs)
	assert.Equal(t, float64(6), testCases[0].Expected)
}

func TestJSONImporter_ImportFromFile_FileNotFound(t *testing.T) {
	importer := NewJSONImporter()
	_, err := importer.ImportFromFile("/nonexistent/file.json")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

func TestJSONImporter_ImportFromFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "invalid.json")

	// Write invalid JSON
	err := os.WriteFile(jsonFile, []byte("{ invalid json }"), 0644)
	assert.NoError(t, err)

	importer := NewJSONImporter()
	_, err = importer.ImportFromFile(jsonFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

func TestJSONImporter_ImportFromFile_MissingTestsField(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "empty.json")

	// Write JSON without tests
	err := os.WriteFile(jsonFile, []byte("{}"), 0644)
	assert.NoError(t, err)

	importer := NewJSONImporter()
	_, err = importer.ImportFromFile(jsonFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no test cases found")
}

func TestJSONImporter_ImportFromFile_MissingNameField(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "missing_name.json")

	jsonContent := `{
  "tests": [
    {
      "inputs": [1, 2],
      "expected": 3
    }
  ]
}`

	err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
	assert.NoError(t, err)

	importer := NewJSONImporter()
	_, err = importer.ImportFromFile(jsonFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing 'name' field")
}

func TestJSONImporter_ImportFromFile_MissingInputsField(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "missing_inputs.json")

	jsonContent := `{
  "tests": [
    {
      "name": "test",
      "expected": 3
    }
  ]
}`

	err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
	assert.NoError(t, err)

	importer := NewJSONImporter()
	_, err = importer.ImportFromFile(jsonFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing 'inputs' field")
}

func TestJSONImporter_ImportFromFile_MissingExpectedField(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "missing_expected.json")

	jsonContent := `{
  "tests": [
    {
      "name": "test",
      "inputs": [1, 2]
    }
  ]
}`

	err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
	assert.NoError(t, err)

	importer := NewJSONImporter()
	_, err = importer.ImportFromFile(jsonFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing 'expected' field")
}

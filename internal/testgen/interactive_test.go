package testgen

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInteractiveInput_ParseValue_Integer(t *testing.T) {
	ii := NewInteractiveInput()

	value, err := ii.parseValue("42")
	assert.NoError(t, err)
	assert.Equal(t, 42, value)
}

func TestInteractiveInput_ParseValue_Float(t *testing.T) {
	ii := NewInteractiveInput()

	value, err := ii.parseValue("3.14")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, value)
}

func TestInteractiveInput_ParseValue_Boolean(t *testing.T) {
	ii := NewInteractiveInput()

	value, err := ii.parseValue("true")
	assert.NoError(t, err)
	assert.Equal(t, true, value)

	value, err = ii.parseValue("false")
	assert.NoError(t, err)
	assert.Equal(t, false, value)
}

func TestInteractiveInput_ParseValue_String(t *testing.T) {
	ii := NewInteractiveInput()

	value, err := ii.parseValue(`"hello"`)
	assert.NoError(t, err)
	assert.Equal(t, "hello", value)

	// Plain string without quotes
	value, err = ii.parseValue("world")
	assert.NoError(t, err)
	assert.Equal(t, "world", value)
}

func TestInteractiveInput_ParseValue_Array(t *testing.T) {
	ii := NewInteractiveInput()

	value, err := ii.parseValue("[1,2,3]")
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{1, 2, 3}, value)

	// Empty array
	value, err = ii.parseValue("[]")
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{}, value)
}

func TestInteractiveInput_ParseInputs(t *testing.T) {
	ii := NewInteractiveInput()

	tests := []struct {
		name     string
		input    string
		expected []interface{}
	}{
		{"single value", "5", []interface{}{5}},
		{"multiple values", "1,2,3", []interface{}{1, 2, 3}},
		{"with spaces", "1, 2, 3", []interface{}{1, 2, 3}},
		{"empty string", "", []interface{}{}},
		{"mixed types", "1,true,hello", []interface{}{1, true, "hello"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ii.parseInputs(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInteractiveInput_Collect_WithMockedInput(t *testing.T) {
	// Mock stdin with test data
	input := strings.NewReader("test 1\n1,2,3\n6\ndone\n")

	ii := &InteractiveInput{
		scanner: bufio.NewScanner(input),
	}

	testCases, err := ii.Collect()
	assert.NoError(t, err)
	assert.Len(t, testCases, 1)
	assert.Equal(t, "test 1", testCases[0].Name)
	assert.Equal(t, []interface{}{1, 2, 3}, testCases[0].Inputs)
	assert.Equal(t, 6, testCases[0].Expected)
}

func TestInteractiveInput_Collect_MultipleTestCases(t *testing.T) {
	input := strings.NewReader("first test\n1,2\n3\nsecond test\n4,5\n9\ndone\n")

	ii := &InteractiveInput{
		scanner: bufio.NewScanner(input),
	}

	testCases, err := ii.Collect()
	assert.NoError(t, err)
	assert.Len(t, testCases, 2)

	assert.Equal(t, "first test", testCases[0].Name)
	assert.Equal(t, []interface{}{1, 2}, testCases[0].Inputs)
	assert.Equal(t, 3, testCases[0].Expected)

	assert.Equal(t, "second test", testCases[1].Name)
	assert.Equal(t, []interface{}{4, 5}, testCases[1].Inputs)
	assert.Equal(t, 9, testCases[1].Expected)
}

func TestInteractiveInput_Collect_ImmediateDone(t *testing.T) {
	input := strings.NewReader("done\n")

	ii := &InteractiveInput{
		scanner: bufio.NewScanner(input),
	}

	_, err := ii.Collect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no test cases provided")
}

func TestInteractiveInput_Collect_EmptyTestCaseName(t *testing.T) {
	// Empty name should be skipped and ask again
	input := strings.NewReader("\nvalid name\n1,2\n3\ndone\n")

	ii := &InteractiveInput{
		scanner: bufio.NewScanner(input),
	}

	testCases, err := ii.Collect()
	assert.NoError(t, err)
	assert.Len(t, testCases, 1)
	assert.Equal(t, "valid name", testCases[0].Name)
}

package testgen

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// InteractiveInput handles interactive test case collection from user
type InteractiveInput struct {
	scanner *bufio.Scanner
}

// NewInteractiveInput creates a new interactive input handler
func NewInteractiveInput() *InteractiveInput {
	return &InteractiveInput{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// Collect prompts user for test cases and returns them
func (i *InteractiveInput) Collect() ([]*TestCase, error) {
	var testCases []*TestCase

	fmt.Println("📝 Interactive Test Case Generator")
	fmt.Println("Enter test cases one by one. Type 'done' when finished.")

	for {
		// Prompt for test case name
		fmt.Print("Enter test case name (or 'done' to finish): ")
		if !i.scanner.Scan() {
			return nil, fmt.Errorf("failed to read test case name")
		}
		name := strings.TrimSpace(i.scanner.Text())

		if name == "done" {
			break
		}

		if name == "" {
			fmt.Println("❌ Test case name cannot be empty. Try again.")
			continue
		}

		// Prompt for inputs
		fmt.Print("Enter inputs (comma-separated, e.g., 1,2,3): ")
		if !i.scanner.Scan() {
			return nil, fmt.Errorf("failed to read inputs")
		}
		inputsStr := strings.TrimSpace(i.scanner.Text())

		inputs, err := i.parseInputs(inputsStr)
		if err != nil {
			fmt.Printf("❌ Invalid input format: %v. Try again.\n", err)
			continue
		}

		// Prompt for expected output
		fmt.Print("Enter expected output: ")
		if !i.scanner.Scan() {
			return nil, fmt.Errorf("failed to read expected output")
		}
		expectedStr := strings.TrimSpace(i.scanner.Text())

		expected, err := i.parseValue(expectedStr)
		if err != nil {
			fmt.Printf("❌ Invalid expected output: %v. Try again.\n", err)
			continue
		}

		testCase := &TestCase{
			Name:     name,
			Inputs:   inputs,
			Expected: expected,
		}

		testCases = append(testCases, testCase)
		fmt.Printf("✅ Added test case: %s\n\n", name)
	}

	if len(testCases) == 0 {
		return nil, fmt.Errorf("no test cases provided")
	}

	fmt.Printf("\n📊 Collected %d test case(s)\n", len(testCases))
	return testCases, nil
}

// parseInputs parses comma-separated input values
func (i *InteractiveInput) parseInputs(inputsStr string) ([]interface{}, error) {
	if inputsStr == "" {
		return []interface{}{}, nil
	}

	parts := strings.Split(inputsStr, ",")
	inputs := make([]interface{}, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		value, err := i.parseValue(part)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, value)
	}

	return inputs, nil
}

// parseValue attempts to parse a string value into appropriate Go type
func (i *InteractiveInput) parseValue(valueStr string) (interface{}, error) {
	valueStr = strings.TrimSpace(valueStr)

	// Try to parse as integer
	if intVal, err := strconv.Atoi(valueStr); err == nil {
		return intVal, nil
	}

	// Try to parse as float
	if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return floatVal, nil
	}

	// Try to parse as boolean
	if boolVal, err := strconv.ParseBool(valueStr); err == nil {
		return boolVal, nil
	}

	// Check for array/slice notation: [1,2,3]
	if strings.HasPrefix(valueStr, "[") && strings.HasSuffix(valueStr, "]") {
		innerStr := strings.TrimSpace(valueStr[1 : len(valueStr)-1])
		if innerStr == "" {
			return []interface{}{}, nil
		}
		return i.parseInputs(innerStr)
	}

	// Default to string (remove quotes if present)
	if strings.HasPrefix(valueStr, "\"") && strings.HasSuffix(valueStr, "\"") {
		return valueStr[1 : len(valueStr)-1], nil
	}

	// Return as-is if it's a plain string
	return valueStr, nil
}

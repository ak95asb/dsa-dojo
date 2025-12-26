package testing

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ak95asb/dsa-dojo/internal/problem"
)

// Executor handles test execution
type Executor struct{}

// NewExecutor creates a new test executor
func NewExecutor() *Executor {
	return &Executor{}
}

// Execute runs go test for the specified problem
func (e *Executor) Execute(prob *problem.ProblemDetails, verbose, race bool) (*TestResult, error) {
	// Construct test file path: problems/<slug>_test.go
	testFile := filepath.Join("problems", problem.SlugToSnakeCase(prob.Slug)+"_test.go")

	// Build go test command arguments
	args := []string{"test"}

	if verbose {
		args = append(args, "-v")
	}

	if race {
		args = append(args, "-race")
	}

	args = append(args, testFile)

	// Execute go test
	cmd := exec.Command("go", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			// Command failed to execute (not just tests failing)
			return nil, fmt.Errorf("failed to execute go test: %w", err)
		}
	}

	// Combine stdout and stderr
	output := stdout.String() + stderr.String()

	// Parse test results
	result := &TestResult{
		AllPassed:    exitCode == 0,
		Output:       output,
		Verbose:      verbose,
		RaceDetector: race,
	}

	e.parseTestResults(result, output)

	return result, nil
}

// parseTestResults extracts pass/fail counts and failed test details from go test output
func (e *Executor) parseTestResults(result *TestResult, output string) {
	lines := strings.Split(output, "\n")

	// Parse test counts and failed tests
	passCount := 0
	failCount := 0
	var currentFailedTest *FailedTest

	for _, line := range lines {
		// Match individual test results
		// Example: "--- PASS: TestTwoSum (0.00s)"
		if strings.Contains(line, "--- PASS:") {
			passCount++
		} else if strings.Contains(line, "--- FAIL:") {
			failCount++
			// Extract test name
			testNameRegex := regexp.MustCompile(`--- FAIL: (\w+)`)
			if matches := testNameRegex.FindStringSubmatch(line); len(matches) > 1 {
				currentFailedTest = &FailedTest{
					Name: matches[1],
				}
				result.FailedTests = append(result.FailedTests, *currentFailedTest)
			}
		}

		// Parse testify assertion failures
		// Example: "Error Trace: ..." followed by "Error: ..." and "Test: ..."
		if currentFailedTest != nil {
			if strings.Contains(line, "expected:") || strings.Contains(line, "Expected:") {
				// Extract expected value
				currentFailedTest.Expected = strings.TrimSpace(strings.Split(line, ":")[1])
			} else if strings.Contains(line, "actual:") || strings.Contains(line, "Actual:") {
				// Extract actual value
				currentFailedTest.Actual = strings.TrimSpace(strings.Split(line, ":")[1])
			} else if strings.Contains(line, "Error:") && !strings.Contains(line, "Error Trace") {
				// Extract error message
				parts := strings.SplitN(line, "Error:", 2)
				if len(parts) > 1 {
					currentFailedTest.Message = strings.TrimSpace(parts[1])
				}
			}
		}

		// Reset current failed test when we hit a new test
		if strings.Contains(line, "=== RUN") {
			currentFailedTest = nil
		}
	}

	// Check for summary line which is more reliable
	// Example: "FAIL	github.com/ak95asb/dsa-dojo/problems	0.123s"
	// Or: "ok  	github.com/ak95asb/dsa-dojo/problems	0.123s"
	for _, line := range lines {
		if strings.HasPrefix(line, "FAIL\t") || strings.HasPrefix(line, "ok  \t") {
			// If we have FAIL, parse the output more carefully
			// For now, use the counts we've gathered
			break
		}
	}

	result.PassedCount = passCount
	result.TotalCount = passCount + failCount

	// Handle case where no individual test results were found
	// This can happen with compilation errors or when tests don't output details
	if result.TotalCount == 0 {
		// Try to detect compilation errors or other issues
		if strings.Contains(output, "FAIL") && (strings.Contains(output, "syntax error") || strings.Contains(output, "undefined:")) {
			result.TotalCount = 1
			result.PassedCount = 0
			result.FailedTests = []FailedTest{
				{
					Name:    "Compilation",
					Message: "Test file failed to compile",
				},
			}
		} else if !result.AllPassed {
			// Some other error occurred
			result.TotalCount = 1
			result.PassedCount = 0
			result.FailedTests = []FailedTest{
				{
					Name:    "Unknown",
					Message: "Test execution failed",
				},
			}
		}
	}
}

// extractNumberFromLine extracts a number from a line containing a label and number
// Example: "PASS: 5" returns 5
func extractNumberFromLine(line, label string) int {
	parts := strings.Split(line, label)
	if len(parts) < 2 {
		return 0
	}
	numStr := strings.TrimSpace(parts[1])
	// Extract just the number part
	numRegex := regexp.MustCompile(`\d+`)
	if matches := numRegex.FindString(numStr); matches != "" {
		num, _ := strconv.Atoi(matches)
		return num
	}
	return 0
}

package testing

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Formatter handles formatting and displaying test results
type Formatter struct {
	useColor bool
}

// NewFormatter creates a new formatter with color detection
func NewFormatter() *Formatter {
	return &Formatter{
		useColor: shouldUseColor(),
	}
}

// shouldUseColor determines if colored output should be used
// Returns false if NO_COLOR is set or stdout is not a TTY
func shouldUseColor() bool {
	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check if stdout is a terminal
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	// Check if output is a character device (TTY)
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// Display formats and prints test results
func (f *Formatter) Display(result *TestResult) {
	if result.AllPassed {
		f.displaySuccess(result)
	} else {
		f.displayFailure(result)
	}

	// Display verbose output if requested
	if result.Verbose && result.Output != "" {
		fmt.Println("\nDetailed Output:")
		fmt.Println("================")
		fmt.Println(result.Output)
	}
}

// displaySuccess shows results when all tests pass
func (f *Formatter) displaySuccess(result *TestResult) {
	if f.useColor {
		green := color.New(color.FgGreen, color.Bold)
		green.Printf("✓ All tests passed! (%d/%d)\n", result.PassedCount, result.TotalCount)
	} else {
		fmt.Printf("✓ All tests passed! (%d/%d)\n", result.PassedCount, result.TotalCount)
	}
}

// displayFailure shows results when tests fail
func (f *Formatter) displayFailure(result *TestResult) {
	if f.useColor {
		red := color.New(color.FgRed, color.Bold)
		red.Printf("✗ Tests failed (%d/%d passed)\n", result.PassedCount, result.TotalCount)
	} else {
		fmt.Printf("✗ Tests failed (%d/%d passed)\n", result.PassedCount, result.TotalCount)
	}

	// Display failed test details
	if len(result.FailedTests) > 0 {
		fmt.Println("\nFailed Tests:")
		fmt.Println("=============")

		for i, test := range result.FailedTests {
			fmt.Printf("\n%d. %s\n", i+1, test.Name)

			if test.Message != "" {
				fmt.Printf("   Error: %s\n", test.Message)
			}

			if test.Expected != "" {
				if f.useColor {
					fmt.Print("   Expected: ")
					color.New(color.FgGreen).Println(test.Expected)
				} else {
					fmt.Printf("   Expected: %s\n", test.Expected)
				}
			}

			if test.Actual != "" {
				if f.useColor {
					fmt.Print("   Actual:   ")
					color.New(color.FgRed).Println(test.Actual)
				} else {
					fmt.Printf("   Actual:   %s\n", test.Actual)
				}
			}
		}
	}
}

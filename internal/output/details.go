package output

import (
	"fmt"
	"strings"

	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/fatih/color"
)

// PrintProblemDetails formats and displays comprehensive problem details
// Includes metadata, description, file paths, and progress information
func PrintProblemDetails(details *problem.ProblemDetails) {
	// Color definitions
	greenColor := color.New(color.FgGreen).SprintFunc()
	yellowColor := color.New(color.FgYellow).SprintFunc()
	redColor := color.New(color.FgRed).SprintFunc()
	boldColor := color.New(color.Bold).SprintFunc()

	// Header with title
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%s\n", boldColor(details.Title))
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Metadata section
	fmt.Printf("%s: ", boldColor("Difficulty"))
	switch details.Difficulty {
	case "easy":
		fmt.Printf("%s\n", greenColor("Easy"))
	case "medium":
		fmt.Printf("%s\n", yellowColor("Medium"))
	case "hard":
		fmt.Printf("%s\n", redColor("Hard"))
	default:
		fmt.Printf("%s\n", details.Difficulty)
	}

	fmt.Printf("%s: %s\n", boldColor("Topic"), details.Topic)
	fmt.Printf("%s: %s\n", boldColor("Slug"), details.Slug)
	fmt.Println()

	// Description section
	fmt.Printf("%s:\n", boldColor("Description"))
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(details.Description)
	fmt.Println()

	// File paths section
	fmt.Printf("%s:\n", boldColor("Files"))
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("  Boilerplate: %s\n", details.BoilerplatePath)
	fmt.Printf("  Tests:       %s\n", details.TestPath)
	fmt.Println()

	// Progress section
	fmt.Printf("%s:\n", boldColor("Progress"))
	fmt.Println(strings.Repeat("-", 80))

	// Status with color
	fmt.Printf("  Status:   ")
	switch details.Status {
	case "completed":
		fmt.Printf("%s\n", greenColor("✓ Solved"))
	case "in_progress":
		fmt.Printf("%s\n", yellowColor("⧗ In Progress"))
	case "not_started":
		fmt.Printf("Not Started\n")
	default:
		fmt.Printf("%s\n", details.Status)
	}

	fmt.Printf("  Attempts: %d\n", details.Attempts)

	if details.Status == "completed" && !details.LastAttempt.IsZero() {
		fmt.Printf("  Last Solved: %s\n", details.LastAttempt.Format("January 2, 2006 3:04 PM"))
	}

	if details.HasSolution {
		fmt.Printf("  Solution File: solutions/%s.go\n", details.Slug)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
}

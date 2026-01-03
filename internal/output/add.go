package output

import (
	"fmt"
	"strings"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/fatih/color"
)

// PrintProblemCreated displays success message after creating a custom problem
func PrintProblemCreated(problem *database.Problem, boilerplatePath, testPath string) {
	greenColor := color.New(color.FgGreen).SprintFunc()
	yellowColor := color.New(color.FgYellow).SprintFunc()
	redColor := color.New(color.FgRed).SprintFunc()
	boldColor := color.New(color.Bold).SprintFunc()
	cyanColor := color.New(color.FgCyan).SprintFunc()

	fmt.Println()
	fmt.Println(greenColor("✓ Custom Problem Created Successfully!"))
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%s\n", boldColor(problem.Title))
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Metadata
	fmt.Printf("%s: %s\n", boldColor("Slug"), problem.Slug)
	fmt.Printf("%s: ", boldColor("Difficulty"))
	switch problem.Difficulty {
	case "easy":
		fmt.Printf("%s\n", greenColor("Easy"))
	case "medium":
		fmt.Printf("%s\n", yellowColor("Medium"))
	case "hard":
		fmt.Printf("%s\n", redColor("Hard"))
	default:
		fmt.Printf("%s\n", problem.Difficulty)
	}

	fmt.Printf("%s: %s\n", boldColor("Topic"), problem.Topic)

	if problem.Tags != "" {
		fmt.Printf("%s: %s\n", boldColor("Tags"), problem.Tags)
	}
	fmt.Println()

	// Files created
	fmt.Println(boldColor("Files Created:"))
	fmt.Printf("  📝 Boilerplate: %s\n", boilerplatePath)
	fmt.Printf("  🧪 Test file:   %s\n", testPath)
	fmt.Println()

	// Next steps
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%s\n", cyanColor(fmt.Sprintf("▶ Run 'dsa solve %s' to start solving!", problem.Slug)))
	fmt.Printf("%s\n", cyanColor(fmt.Sprintf("▶ Run 'dsa show %s' to view problem details", problem.Slug)))
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
}

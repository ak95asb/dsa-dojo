package output

import (
	"fmt"
	"strings"

	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/fatih/color"
)

// PrintRandomProblem formats and displays a random problem suggestion
func PrintRandomProblem(details *problem.ProblemDetails) {
	// Color definitions
	greenColor := color.New(color.FgGreen).SprintFunc()
	yellowColor := color.New(color.FgYellow).SprintFunc()
	redColor := color.New(color.FgRed).SprintFunc()
	boldColor := color.New(color.Bold).SprintFunc()
	cyanColor := color.New(color.FgCyan).SprintFunc()

	fmt.Println()
	fmt.Println(cyanColor("🎲 Random Problem Selected!"))
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

	// Suggested next step
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%s\n", cyanColor(fmt.Sprintf("▶ Run 'dsa solve %s' to start solving!", details.Slug)))
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
}

// PrintNoProblemsMessage displays a helpful message when no problems match filters
func PrintNoProblemsMessage(filters problem.ListFilters) {
	redColor := color.New(color.FgRed).SprintFunc()
	yellowColor := color.New(color.FgYellow).SprintFunc()

	fmt.Println()
	fmt.Println(redColor("✗ No unsolved problems found!"))
	fmt.Println()

	// Build descriptive message based on filters
	if filters.Difficulty != "" && filters.Topic != "" {
		fmt.Printf("All %s %s problems are solved!\n", filters.Difficulty, filters.Topic)
	} else if filters.Difficulty != "" {
		fmt.Printf("All %s problems are solved!\n", filters.Difficulty)
	} else if filters.Topic != "" {
		fmt.Printf("All %s problems are solved!\n", filters.Topic)
	} else {
		fmt.Println("All problems in your library are solved!")
	}

	fmt.Println()
	fmt.Println(yellowColor("Suggestions:"))

	// Provide context-aware suggestions
	if filters.Difficulty == "easy" {
		fmt.Println("  • Try --difficulty medium for a bigger challenge")
		fmt.Println("  • Try --difficulty hard for expert-level problems")
	} else if filters.Difficulty == "medium" {
		fmt.Println("  • Try --difficulty hard for expert-level problems")
		fmt.Println("  • Try --difficulty easy to review fundamentals")
	} else if filters.Difficulty == "hard" {
		fmt.Println("  • Congratulations on solving all hard problems!")
		fmt.Println("  • Try --difficulty medium or --difficulty easy to practice speed")
	}

	if filters.Topic != "" {
		fmt.Println("  • Try a different topic to broaden your skills")
		fmt.Println("  • Run 'dsa list' to see all available topics")
	} else {
		fmt.Println("  • Run 'dsa list' to see all problems")
		fmt.Println("  • Consider adding custom problems with 'dsa add'")
	}

	fmt.Println()
}

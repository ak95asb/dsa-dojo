package output

import (
	"fmt"
	"os"
)

// FormatCelebration creates a congratulations message for first-time problem solves
func FormatCelebration(problemTitle string, attempts int) string {
	noColor := os.Getenv("NO_COLOR") != ""

	var output string

	// Congratulations header
	if noColor {
		output += fmt.Sprintf("Congratulations! You solved %s!\n\n", problemTitle)
	} else {
		output += fmt.Sprintf("\U0001F389 \033[1;32mCongratulations! You solved %s!\033[0m\n\n", problemTitle)
	}

	// Attempt count
	attemptText := "attempt"
	if attempts != 1 {
		attemptText = "attempts"
	}

	if noColor {
		output += fmt.Sprintf("Solved in %d %s\n", attempts, attemptText)
	} else {
		output += fmt.Sprintf("\u23F1\uFE0F  Solved in %d %s\n", attempts, attemptText)
	}

	// Next steps
	if noColor {
		output += "View your progress with: dsa status\n\n"
	} else {
		output += "\U0001F4CA View your progress with: \033[1mdsa status\033[0m\n\n"
	}

	// Encouragement
	if noColor {
		output += "Keep up the great work!"
	} else {
		output += "Keep up the great work! \U0001F680"
	}

	return output
}

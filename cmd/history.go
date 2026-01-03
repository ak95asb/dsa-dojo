package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/ak95asb/dsa-dojo/internal/solution"
	"github.com/spf13/cobra"
)

var (
	historyShow    int
	historyRestore int
)

var historyCmd = &cobra.Command{
	Use:   "history [problem-id]",
	Short: "View solution submission history",
	Long: `Display all solution attempts for a problem.

Options:
  --show N     Display solution code from Nth attempt (1 = most recent)
  --restore N  Restore Nth attempt as current solution (1 = most recent)

Examples:
  dsa history two-sum
  dsa history two-sum --show 2
  dsa history two-sum --restore 3`,
	Args: cobra.ExactArgs(1),
	Run:  runHistoryCommand,
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().IntVar(&historyShow, "show", 0, "Display solution code from Nth attempt")
	historyCmd.Flags().IntVar(&historyRestore, "restore", 0, "Restore Nth attempt as current solution")
}

func runHistoryCommand(cmd *cobra.Command, args []string) {
	slug := args[0]

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to connect to database: %v\n", err)
		os.Exit(3) // ExitDatabaseError
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Get problem by slug
	problemSvc := problem.NewService(db)
	prob, err := problemSvc.GetProblemBySlug(slug)
	if err != nil {
		if errors.Is(err, problem.ErrProblemNotFound) {
			fmt.Fprintf(os.Stderr, "Problem '%s' not found. Run 'dsa list' to see available problems.\n", slug)
			os.Exit(2) // ExitUsageError
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create solution service
	solutionSvc := solution.NewService(db)

	// Route based on flags
	if historyShow > 0 {
		showSolution(solutionSvc, prob.ID, slug, historyShow)
	} else if historyRestore > 0 {
		restoreSolution(solutionSvc, prob.ID, slug, historyRestore)
	} else {
		listHistory(solutionSvc, prob.ID, slug)
	}

	os.Exit(0)
}

// listHistory displays all submissions for a problem
func listHistory(svc *solution.Service, problemID uint, slug string) {
	records, err := svc.GetHistory(problemID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving history: %v\n", err)
		os.Exit(1)
	}

	if len(records) == 0 {
		fmt.Printf("No submission history for %s\n", slug)
		fmt.Printf("\nSubmit your first solution with: dsa submit %s\n", slug)
		return
	}

	// Display header
	fmt.Printf("Solution History for %s:\n\n", slug)
	fmt.Printf(" #  Date & Time           Status    Tests\n")
	fmt.Printf("==  ====================  ========  =====\n")

	// Display submissions
	for i, record := range records {
		index := i + 1
		dateTime := record.CreatedAt.Format("2006-01-02 15:04:05")

		var status string
		if record.Passed {
			status = "✓ Passed"
		} else {
			status = "✗ Failed"
		}

		// Note: We don't store test counts in database.Solution model
		// For this display, we'll show passed/failed status only
		fmt.Printf("%2d  %s  %-8s  -\n", index, dateTime, status)
	}

	// Display usage hints
	fmt.Printf("\nUse 'dsa history %s --show N' to view solution #N\n", slug)
	fmt.Printf("Use 'dsa history %s --restore N' to restore solution #N\n", slug)
}

// showSolution displays code for a specific submission
func showSolution(svc *solution.Service, problemID uint, slug string, index int) {
	record, err := svc.GetSubmissionByIndex(problemID, index)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Display header
	fmt.Printf("Solution #%d for %s\n", index, slug)
	fmt.Printf("Date: %s\n", record.CreatedAt.Format("2006-01-02 15:04:05"))

	if record.Passed {
		fmt.Printf("Status: ✓ Passed\n")
	} else {
		fmt.Printf("Status: ✗ Failed\n")
	}

	// Display code
	fmt.Println("\n--- Code ---")
	fmt.Println(record.Code)
}

// restoreSolution restores a previous submission as current solution
func restoreSolution(svc *solution.Service, problemID uint, slug string, index int) {
	// Get submission record
	record, err := svc.GetSubmissionByIndex(problemID, index)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Prompt for confirmation
	timestamp := record.CreatedAt.Format("2006-01-02 15:04:05")
	fmt.Printf("Restore solution from %s? Current solution will be backed up. [y/N]: ", timestamp)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))

	if response != "y" && response != "yes" {
		fmt.Println("Restoration cancelled.")
		return
	}

	// Backup current solution
	backupPath, err := svc.BackupCurrentSolution(slug, problemID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating backup: %v\n", err)
		os.Exit(1)
	}

	if backupPath != "" {
		fmt.Printf("✓ Current solution backed up to %s\n", backupPath)
	}

	// Restore solution
	if err := svc.RestoreSolution(slug, record); err != nil {
		fmt.Fprintf(os.Stderr, "Error restoring solution: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Solution restored from history")
	fmt.Printf("  File: solutions/%s.go\n", slug)
	fmt.Printf("  Timestamp: %s\n", timestamp)
}

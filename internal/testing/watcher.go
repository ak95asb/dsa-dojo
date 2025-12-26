package testing

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ak95asb/dsa-dojo/internal/problem"
	"github.com/fsnotify/fsnotify"
)

// TestState tracks the previous test result for transition detection
type TestState struct {
	LastResult *TestResult
}

// DetectTransition compares current result with previous and returns transition type
func (s *TestState) DetectTransition(current *TestResult) string {
	if s.LastResult == nil {
		s.LastResult = current
		return "" // First run
	}

	if !s.LastResult.AllPassed && current.AllPassed {
		s.LastResult = current
		return "fail_to_pass" // 🎉 Tests now passing!
	}

	if s.LastResult.AllPassed && !current.AllPassed {
		s.LastResult = current
		return "pass_to_fail" // ⚠️  Tests broken
	}

	s.LastResult = current
	return "no_change"
}

// Watch monitors the solution file for changes and re-runs tests automatically
func (s *Service) Watch(prob *problem.ProblemDetails, problemSvc *problem.Service, verbose, race bool) error {
	// Construct solution file path
	solutionPath := filepath.Join("solutions", prob.Slug+".go")
	solutionDir := "solutions"

	// Verify solution file exists
	if _, err := os.Stat(solutionPath); os.IsNotExist(err) {
		return fmt.Errorf("solution file not found: %s", solutionPath)
	}

	// Get absolute path for display
	absPath, _ := filepath.Abs(solutionPath)

	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}
	defer watcher.Close()

	// Watch the parent directory (solutions/) not individual file
	// This handles atomic saves where editors write to temp file then rename
	if err := watcher.Add(solutionDir); err != nil {
		return fmt.Errorf("failed to watch directory %s: %w", solutionDir, err)
	}

	// Set up signal handling for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initialize test state for transition detection
	testState := &TestState{}

	// Display initial message
	fmt.Printf("👀 Watching %s for changes... (Press Ctrl+C to stop)\n\n", absPath)

	// Run tests immediately before starting watch
	s.runTestsInWatchMode(prob, problemSvc, verbose, race, testState, true)

	// Debouncing variables
	var debounceTimer *time.Timer
	const debounceDuration = 100 * time.Millisecond

	// Watch loop
	for {
		select {
		case <-sigChan:
			// Ctrl+C pressed - clean shutdown
			fmt.Println("\nWatch mode stopped")
			return nil

		case event, ok := <-watcher.Events:
			if !ok {
				return fmt.Errorf("watcher events channel closed")
			}

			// Filter: only process Write events for our target file
			// Also handle Create/Rename for atomic saves (vim, emacs)
			if event.Name == solutionPath && (event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Rename == fsnotify.Rename) {

				// Debounce: ignore rapid successive writes within 100ms
				if debounceTimer != nil {
					debounceTimer.Stop()
				}

				debounceTimer = time.AfterFunc(debounceDuration, func() {
					s.runTestsInWatchMode(prob, problemSvc, verbose, race, testState, false)
				})
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return fmt.Errorf("watcher errors channel closed")
			}
			fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
		}
	}
}

// runTestsInWatchMode executes tests and displays results with transition detection
func (s *Service) runTestsInWatchMode(prob *problem.ProblemDetails, problemSvc *problem.Service, verbose, race bool, testState *TestState, isInitial bool) {
	// Clear terminal and show re-running message (skip on initial run)
	if !isInitial {
		clearTerminal()
		fmt.Println("🔄 Re-running tests...")
		fmt.Println()
	}

	// Show timestamp
	fmt.Printf("⏰ %s\n\n", time.Now().Format("15:04:05"))

	// Execute tests
	result, err := s.ExecuteTests(prob, verbose, race)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running tests: %v\n", err)
		return
	}

	// Display results
	s.DisplayResults(result)

	// Detect state transition
	transition := testState.DetectTransition(result)

	// Handle transitions
	switch transition {
	case "fail_to_pass":
		fmt.Println("\n🎉 Tests now passing!")

		// Update progress to completed
		if err := problemSvc.UpdateProgress(prob.ID, "completed"); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to update progress: %v\n", err)
		}

		// Create/update solution record
		if err := s.RecordSolution(prob.ID, result); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to record solution: %v\n", err)
		}

	case "pass_to_fail":
		fmt.Println("\n⚠️  Tests broken - check your changes")

	case "no_change":
		// No special message for no change
		if result.AllPassed && !isInitial {
			// Still passing
			fmt.Println("\n✅ Tests still passing")
		}
	}

	fmt.Println() // Add spacing before next watch message
}

// clearTerminal clears the terminal screen using ANSI escape codes
// Works on Windows 10+, macOS, and Linux
func clearTerminal() {
	fmt.Print("\033[2J\033[H") // ANSI: clear screen + move cursor to home
}

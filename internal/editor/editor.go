package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/viper"
)

// Detect returns the editor command to use
// Priority: Viper config > $EDITOR env var > platform default
func Detect() string {
	// 1. Check Viper config
	if editor := viper.GetString("editor"); editor != "" {
		return editor
	}

	// 2. Check $EDITOR environment variable
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	// 3. Platform defaults
	switch runtime.GOOS {
	case "windows":
		return "notepad"
	default: // Unix-like (macOS, Linux)
		return "vi"
	}
}

// Launch opens the file in the detected editor
// Uses cmd.Start() to avoid blocking the CLI
func Launch(editorCmd, filePath string) error {
	cmd := exec.Command(editorCmd, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch %s: %w", editorCmd, err)
	}

	// Don't wait for editor to close (allow background editing)
	return nil
}

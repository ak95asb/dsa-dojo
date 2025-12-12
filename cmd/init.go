/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ak95asb/dsa-dojo/internal/database"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your DSA practice workspace",
	Long: `Initialize creates the ~/.dsa directory and database for storing
your practice progress, solutions, and problem data.

This command is safe to run multiple times - it will detect an existing
workspace and skip initialization.

Exit Codes:
  0 - Success (workspace initialized or already exists)
  3 - Database error (check directory permissions)`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get home directory for displaying in messages
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "✗ Error: Failed to get home directory")
			fmt.Fprintf(os.Stderr, "  %v\n", err)
			os.Exit(3)
		}

		dsaDir := filepath.Join(homeDir, ".dsa")
		dbPath := filepath.Join(dsaDir, "dsa.db")

		// Check if workspace already exists
		if _, err := os.Stat(dbPath); err == nil {
			fmt.Printf("Workspace already initialized at %s\n", dsaDir)
			return // Exit with code 0 (success)
		}

		// Initialize database
		db, err := database.Initialize()
		if err != nil {
			fmt.Fprintln(os.Stderr, "✗ Error: Failed to initialize workspace")
			fmt.Fprintf(os.Stderr, "  %v\n", err)
			fmt.Fprintln(os.Stderr, "\nTroubleshooting:")
			fmt.Fprintln(os.Stderr, "  • Check directory permissions for ~/.dsa")
			fmt.Fprintln(os.Stderr, "  • Ensure sufficient disk space")
			fmt.Fprintln(os.Stderr, "  • Verify SQLite dependencies are available")
			os.Exit(3)
		}

		// Seed initial problem library
		count, err := database.SeedProblems(db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Warning: Failed to seed problem library: %v\n", err)
			// Don't exit - workspace is still usable without seeded problems
		} else if count > 0 {
			fmt.Printf("✓ Seeded %d problems to library\n", count)
		}

		// Verify database connection (optional sanity check)
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}

		fmt.Printf("✓ Workspace initialized at %s\n", dsaDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

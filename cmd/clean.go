package cmd

import (
	"aipad/internal/state"
	syncpkg "aipad/internal/sync"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove context from rules and config files",
	Long: `Clean up all synced context from provider rules directories and config files.

This command will:
- Remove scratchpad copies from .claude/rules/ and .agent/rules/
- Remove the managed context block from CLAUDE.md and AGENTS.md
- Keep the original scratchpad in .aipad/ intact`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if session exists
		s, err := state.Load()
		if err != nil {
			fmt.Println("No active session found.")
			os.Exit(1)
		}

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Cleaning synced context...")

		// Clean all provider rules directories and config files
		for name, config := range s.Providers {
			// Skip duplicates (ag is same as antigravity)
			if name == "ag" {
				continue
			}

			// Remove scratchpad from rules directory
			rulesScrtachpad := filepath.Join(cwd, config.RulesDir, "scratchpad.md")
			if err := os.Remove(rulesScrtachpad); err != nil {
				if !os.IsNotExist(err) {
					fmt.Printf("Warning: Could not remove %s: %v\n", rulesScrtachpad, err)
				}
			} else {
				fmt.Printf("Removed %s\n", rulesScrtachpad)
			}

			// Clear managed block from config file
			configPath := filepath.Join(cwd, config.ConfigFile)
			if err := syncpkg.UpdateConfigWithManagedBlock(configPath, ""); err != nil {
				fmt.Printf("Warning: Could not clear %s: %v\n", configPath, err)
			} else {
				fmt.Printf("Cleared managed block in %s\n", config.ConfigFile)
			}
		}

		fmt.Println("\nClean complete! Original scratchpad in .aipad/ is preserved.")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

package cmd

import (
	"aipad/internal/state"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show conversation history",
	Long: `Display the conversation history from the scratchpad.
Shows all context entries with their timestamps.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if session exists
		_, err := state.Load()
		if err != nil {
			fmt.Println("No active session found. Run 'aipad new <provider>' to start.")
			os.Exit(1)
		}

		// Read scratchpad
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		scratchpadPath := filepath.Join(cwd, state.AIPadDir, state.ScratchpadFile)

		content, err := os.ReadFile(scratchpadPath)
		if err != nil {
			fmt.Println("No scratchpad found.")
			os.Exit(1)
		}

		if len(strings.TrimSpace(string(content))) == 0 {
			fmt.Println("Scratchpad is empty. Use 'aipad convo \"<text>\"' to add context.")
			return
		}

		fmt.Println("╔══════════════════════════════════════════╗")
		fmt.Println("║         Conversation History             ║")
		fmt.Println("╚══════════════════════════════════════════╝")
		fmt.Println()

		// Parse entries
		entries := parseEntries(string(content))
		for i, entry := range entries {
			fmt.Printf("  [%d] %s\n", i+1, entry.Timestamp)
			// Truncate content to 80 chars
			preview := strings.ReplaceAll(entry.Content, "\n", " ")
			if len(preview) > 80 {
				preview = preview[:80] + "..."
			}
			fmt.Printf("      %s\n\n", preview)
		}
	},
}

type Entry struct {
	Timestamp string
	Content   string
}

func parseEntries(content string) []Entry {
	var entries []Entry
	// Match pattern: ## [timestamp] Context Update
	pattern := regexp.MustCompile(`## \[([^\]]+)\] Context Update\n([\s\S]*?)(?:---|$)`)
	matches := pattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			entries = append(entries, Entry{
				Timestamp: match[1],
				Content:   strings.TrimSpace(match[2]),
			})
		}
	}
	return entries
}

func init() {
	rootCmd.AddCommand(listCmd)
}

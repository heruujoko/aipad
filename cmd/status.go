package cmd

import (
	"aipad/internal/state"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current provider and session info",
	Long: `Display the current AIPad session status including:
- Current AI provider
- Session ID
- Created at timestamp
- Last sync timestamp
- Number of context entries`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := state.Load()
		if err != nil {
			fmt.Println("No active session found. Run 'aipad new <provider>' to start.")
			os.Exit(1)
		}

		fmt.Println("╔══════════════════════════════════════════╗")
		fmt.Println("║           AIPad Session Status           ║")
		fmt.Println("╚══════════════════════════════════════════╝")
		fmt.Printf("  Provider:    %s\n", s.CurrentProvider)
		fmt.Printf("  Session ID:  %s\n", s.SessionID)
		fmt.Printf("  Created:     %s\n", s.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Last Sync:   %s\n", s.LastSync.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Entries:     %d context(s)\n", len(s.ContextHashes))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

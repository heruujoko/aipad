package cmd

import (
	"aipad/internal/crypto"
	"aipad/internal/state"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// convoCmd represents the convo command
var convoCmd = &cobra.Command{
	Use:   "convo \"<text>\"",
	Short: "Add conversation context to the scratchpad",
	Long: `Append conversation context to the scratchpad with a timestamp.
The content is hashed and checked for duplicates before being added.

Example:
  aipad convo "Discussed the new API design with focus on REST principles"`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one argument: the conversation text")
		}
		if len(args[0]) == 0 {
			return fmt.Errorf("conversation text cannot be empty")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		text := args[0]

		// 1. Load existing state
		s, err := state.Load()
		if err != nil {
			fmt.Printf("Error: No active session found. Run 'aipad new <provider>' first.\n")
			os.Exit(1)
		}

		// 2. Generate hash for deduplication
		hash := crypto.GenerateHash(text)

		// 3. Check for duplicates
		if crypto.IsDuplicate(hash, s.ContextHashes) {
			fmt.Println("Duplicate content detected. Skipping addition.")
			return
		}

		// 4. Append to scratchpad.md
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		scratchpadPath := filepath.Join(cwd, state.AIPadDir, state.ScratchpadFile)

		// Format the entry
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		entry := fmt.Sprintf("\n## [%s] Context Update\n%s\n---\n", timestamp, text)

		// Append to file
		f, err := os.OpenFile(scratchpadPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("Error opening scratchpad: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		if _, err := f.WriteString(entry); err != nil {
			fmt.Printf("Error writing to scratchpad: %v\n", err)
			os.Exit(1)
		}

		// 5. Update state with new hash
		s.ContextHashes = append(s.ContextHashes, hash)
		s.LastSync = time.Now()
		if err := s.Save(); err != nil {
			fmt.Printf("Error saving state: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Context added to scratchpad.")
	},
}

func init() {
	rootCmd.AddCommand(convoCmd)
}

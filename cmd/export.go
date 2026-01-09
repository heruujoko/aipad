package cmd

import (
	"aipad/internal/state"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export [filename]",
	Short: "Export conversation history to a file",
	Long: `Export the conversation history from the scratchpad to a file.

This command will:
- Read the current scratchpad content
- Export it to the specified file (default: export-<timestamp>.md)
- Include metadata about the session in the export

Supported formats: .md (markdown), .txt (text), .json

Example:
  aipad export
  aipad export conversation.md
  aipad export conversation.json`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("accepts at most one argument: [filename]")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Load existing state
		s, err := state.Load()
		if err != nil {
			fmt.Printf("Error: No active session found. Run 'aipad new <provider>' first.\n")
			os.Exit(1)
		}

		// 2. Determine output filename
		outputFile := "export-" + time.Now().Format("20060102-150405") + ".md"
		if len(args) == 1 {
			outputFile = args[0]
		}

		// 3. Read scratchpad content
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		scratchpadPath := filepath.Join(cwd, state.AIPadDir, state.ScratchpadFile)

		content, err := os.ReadFile(scratchpadPath)
		if err != nil {
			fmt.Printf("Error reading scratchpad: %v\n", err)
			os.Exit(1)
		}

		// 4. Determine file extension and format
		ext := filepath.Ext(outputFile)
		var exportContent string

		switch ext {
		case ".json":
			exportContent = fmt.Sprintf(`{
  "session_id": "%s",
  "provider": "%s",
  "created_at": "%s",
  "exported_at": "%s",
  "entries": %d
}`,
				s.SessionID,
				s.CurrentProvider,
				s.CreatedAt.Format(time.RFC3339),
				time.Now().Format(time.RFC3339),
				len(s.ContextHashes))
		case ".txt":
			exportContent = fmt.Sprintf("AIPad Conversation Export\n"+
				"==========================\n"+
				"Session ID: %s\n"+
				"Provider: %s\n"+
				"Created: %s\n"+
				"Exported: %s\n"+
				"Total Entries: %d\n"+
				"\n%s",
				s.SessionID,
				s.CurrentProvider,
				s.CreatedAt.Format("2006-01-02 15:04:05"),
				time.Now().Format("2006-01-02 15:04:05"),
				len(s.ContextHashes),
				string(content))
		default: // .md or any other format
			exportContent = fmt.Sprintf("# AIPad Conversation Export\n\n"+
				"**Session ID:** %s\n\n"+
				"**Provider:** %s\n\n"+
				"**Created:** %s\n\n"+
				"**Exported:** %s\n\n"+
				"**Total Entries:** %d\n\n"+
				"---\n\n"+
				"%s",
				s.SessionID,
				s.CurrentProvider,
				s.CreatedAt.Format("2006-01-02 15:04:05"),
				time.Now().Format("2006-01-02 15:04:05"),
				len(s.ContextHashes),
				string(content))
		}

		// 5. Write to file
		outputPath := filepath.Join(cwd, outputFile)
		if err := os.WriteFile(outputPath, []byte(exportContent), 0644); err != nil {
			fmt.Printf("Error writing export file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Exported conversation history to: %s\n", outputFile)
		fmt.Printf("  Session: %s\n", s.SessionID)
		fmt.Printf("  Entries: %d\n", len(s.ContextHashes))
		fmt.Printf("  Format: %s\n", ext)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

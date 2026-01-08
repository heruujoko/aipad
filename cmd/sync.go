package cmd

import (
	"aipad/internal/state"
	syncpkg "aipad/internal/sync"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync [provider]",
	Short: "Manually trigger context sync to provider",
	Long: `Manually sync the scratchpad content to the provider's configuration.

This command will:
- Read the current scratchpad content
- Create the provider's rules directory if needed
- Copy the scratchpad to the rules directory
- Update the provider's config file with the current context
- Update the last_sync timestamp in state.json

If no provider is specified, it syncs to the current provider.

Valid providers are: claude, antigravity, ag

Example:
  aipad sync
  aipad sync antigravity`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("accepts at most one argument: [provider]")
		}
		if len(args) == 1 {
			provider := args[0]
			// Check if provider exists in the builtin or custom providers
			s := state.NewState("claude")
			if _, ok := s.Providers[provider]; !ok {
				return fmt.Errorf("unsupported provider: '%s'. Use 'aipad providers list' to see available providers", provider)
			}
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

		// 2. Determine provider to sync to
		provider := s.CurrentProvider
		if len(args) == 1 {
			provider = args[0]
		}

		fmt.Printf("Syncing context to provider: %s\n", provider)

		// 3. Get provider config
		providerConfig, ok := s.Providers[provider]
		if !ok {
			fmt.Printf("Error: Provider configuration not found for '%s'\n", provider)
			os.Exit(1)
		}

		// 4. Create rules directory
		if err := syncpkg.EnsureRulesDir(providerConfig.RulesDir); err != nil {
			fmt.Printf("Error creating rules directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Ensured rules directory: %s\n", providerConfig.RulesDir)

		// 5. Copy scratchpad to rules directory
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		scratchpadPath := filepath.Join(cwd, state.AIPadDir, state.ScratchpadFile)

		if err := syncpkg.CopyScratchpadToRules(scratchpadPath, providerConfig.RulesDir); err != nil {
			fmt.Printf("Error copying scratchpad: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Copied scratchpad to %s\n", providerConfig.RulesDir)

		// 6. Update config file with managed block
		configPath := filepath.Join(cwd, providerConfig.ConfigFile)
		if err := syncpkg.SyncProviderConfig(configPath, scratchpadPath); err != nil {
			fmt.Printf("Error updating config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Updated %s with current context\n", providerConfig.ConfigFile)

		// 7. Update last sync timestamp
		s.LastSync = time.Now()
		if err := s.Save(); err != nil {
			fmt.Printf("Error saving state: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nSync complete! Last sync: %s\n", s.LastSync.Format("2006-01-02 15:04:05"))
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

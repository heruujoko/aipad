package cmd

import (
	"aipad/internal/state"
	syncpkg "aipad/internal/sync"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use <provider>",
	Short: "Switch to a different AI provider",
	Long: `Switch to a different AI provider and sync the context.

This command will:
- Update the current provider in state.json
- Create the provider's rules directory if needed
- Copy the scratchpad to the rules directory
- Update the provider's config file with the current context

Valid providers are: claude, antigravity, ag

Example:
  aipad use antigravity`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one argument: <provider>")
		}
		provider := args[0]
		if !validProviders[provider] {
			return fmt.Errorf("unsupported provider: '%s'. Valid providers are: claude, antigravity, ag", provider)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		provider := args[0]

		// 1. Load existing state
		s, err := state.Load()
		if err != nil {
			fmt.Printf("Error: No active session found. Run 'aipad new <provider>' first.\n")
			os.Exit(1)
		}

		// 2. Update current provider
		s.CurrentProvider = provider
		if err := s.Save(); err != nil {
			fmt.Printf("Error saving state: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Switched to provider: %s\n", provider)

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

		fmt.Println("\nProvider switch complete! The AI assistant should now have access to your context.")
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}

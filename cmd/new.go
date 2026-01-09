package cmd

import (
	"aipad/internal/state"
	syncpkg "aipad/internal/sync"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var validProviders = map[string]bool{
	"claude":      true,
	"antigravity": true,
	"ag":          true,
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <provider>",
	Short: "Initialize a new AIPad session",
	Long: `Initialize a new AIPad session with a specific AI provider.
Valid providers are: claude, antigravity, ag.

This command will:
- Create the .aipad/ directory
- Initialize state.json
- Initialize scratchpad.md
- Set up provider configuration`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one argument: <provider>")
		}
		provider := args[0]
		// Check if provider exists in the builtin or custom providers
		s := state.NewState("claude")
		if _, ok := s.Providers[provider]; !ok {
			return fmt.Errorf("unsupported provider: '%s'. Use 'aipad providers list' to see available providers", provider)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		provider := args[0]
		fmt.Printf("Initializing AIPad session for provider: %s\n", provider)

		// 1. Initialize .aipad directory
		if err := state.InitAIPadDir(); err != nil {
			fmt.Printf("Error creating .aipad directory: %v\n", err)
			os.Exit(1)
		}

		// 2. Create/Update state.json
		s := state.NewState(provider)
		if err := s.Save(); err != nil {
			fmt.Printf("Error saving state: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Initialized .aipad/state.json")

		// 3. Initialize scratchpad.md
		if err := state.EnsureScratchpad(); err != nil {
			fmt.Printf("Error creating scratchpad: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Initialized .aipad/scratchpad.md")

		// 4. Create provider-specific config if missing
		// Note from specs: "Set up provider-specific configuration file if it doesn't exist"
		providerConfig := s.Providers[provider]
		configPath := providerConfig.ConfigFile // Config file is relative to cwd

		// Create file if it doesn't exist, populated with Agent Awareness instructions
		if err := ensureConfigFileWithInstructions(configPath); err != nil {
			fmt.Printf("Error ensuring provider config %s: %v\n", configPath, err)
			os.Exit(1)
		}
		fmt.Printf("Ensured provider config exists: %s\n", configPath)

		// 5. Initial Sync: Ensure the config file has the current scratchpad content (likely empty or just init)
		// This also ensures the managed block structure is correct even if the file existed but was empty/malformed
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		scratchpadPath := filepath.Join(cwd, state.AIPadDir, state.ScratchpadFile)

		if err := syncpkg.SyncProviderConfig(configPath, scratchpadPath); err != nil {
			fmt.Printf("Error syncing initial context to config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Synced initial context to %s\n", configPath)

		fmt.Printf("Successfully started session! You are now using: %s\n", provider)
	},
}

func ensureConfigFileWithInstructions(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create with basic header, SyncProviderConfig will fill in the managed block
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		// Write a minimal header
		_, err = file.WriteString(fmt.Sprintf("# %s Configuration\n\n", filename))
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)
}

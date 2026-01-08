package cmd

import (
	"aipad/internal/state"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var validProviders = map[string]bool{
	"claude":      true,
	"antigravity": true,
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <provider>",
	Short: "Initialize a new AIPad session",
	Long: `Initialize a new AIPad session with a specific AI provider.
Valid providers are: claude, antigravity.

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
		if !validProviders[provider] {
			return fmt.Errorf("unsupported provider: '%s'. Valid providers are: claude, antigravity", provider)
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
		err := ensureFileExists(providerConfig.ConfigFile)
		if err != nil {
			fmt.Printf("Error creating provider config %s: %v\n", providerConfig.ConfigFile, err)
			os.Exit(1)
		}
		fmt.Printf("Ensured provider config exists: %s\n", providerConfig.ConfigFile)

		fmt.Printf("Successfully started session! You are now using: %s\n", provider)
	},
}

func ensureFileExists(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		// Optionally write a header
		_, err = file.WriteString(fmt.Sprintf("# %s Configuration\n", filename))
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)
}

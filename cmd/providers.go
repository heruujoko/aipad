package cmd

import (
	"aipad/internal/config"
	"aipad/internal/state"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// providersCmd represents the providers command
var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "Manage custom provider configurations",
	Long: `Manage custom provider configurations.

This command allows you to add, remove, and list custom AI providers
beyond the built-in ones (claude, antigravity, ag).

Custom providers are stored in .aipad/providers.json or ~/.aipad/providers.json

Example:
  aipad providers add myai MYAI.md .myai/rules/
  aipad providers remove myai
  aipad providers list`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'aipad providers --help' to see available subcommands")
	},
}

// addProviderCmd represents the providers add command
var addProviderCmd = &cobra.Command{
	Use:   "add <name> <config-file> <rules-dir>",
	Short: "Add a custom provider",
	Long: `Add a custom provider configuration.

Arguments:
  name        - Unique identifier for the provider
  config-file - Path to the provider's config file (e.g., MYAI.md)
  rules-dir   - Path to the provider's rules directory (e.g., .myai/rules/)

Example:
  aipad providers add myai MYAI.md .myai/rules/`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return fmt.Errorf("requires exactly three arguments: <name> <config-file> <rules-dir>")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		configFile := args[1]
		rulesDir := args[2]

		// Check if it's a builtin provider
		if validProviders[name] {
			fmt.Printf("Error: Cannot override builtin provider '%s'\n", name)
			os.Exit(1)
		}

		if err := config.AddProvider(name, configFile, rulesDir); err != nil {
			fmt.Printf("Error adding provider: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully added custom provider '%s'\n", name)
		fmt.Printf("  Config file: %s\n", configFile)
		fmt.Printf("  Rules dir:   %s\n", rulesDir)
		fmt.Println("\nYou can now use this provider with:")
		fmt.Printf("  aipad new %s\n", name)
		fmt.Printf("  aipad use %s\n", name)
	},
}

// removeProviderCmd represents the providers remove command
var removeProviderCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a custom provider",
	Long: `Remove a custom provider configuration.

Example:
  aipad providers remove myai`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one argument: <name>")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// Check if it's a builtin provider
		if validProviders[name] {
			fmt.Printf("Error: Cannot remove builtin provider '%s'\n", name)
			os.Exit(1)
		}

		if err := config.RemoveProvider(name); err != nil {
			fmt.Printf("Error removing provider: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully removed custom provider '%s'\n", name)
	},
}

// listProvidersCmd represents the providers list command
var listProvidersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all providers",
	Long: `List all available providers (builtin and custom).

Shows the provider name, config file, and rules directory for each provider.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("╔══════════════════════════════════════════╗")
		fmt.Println("║          Available Providers            ║")
		fmt.Println("╚══════════════════════════════════════════╝")
		fmt.Println("\nBuiltin Providers:")

		// Load state to get all providers
		s, err := state.Load()
		if err != nil {
			// If no session, just show builtin providers
			s = state.NewState("claude")
		}

		// List builtin providers
		builtinProviders := []string{"claude", "antigravity", "ag"}
		for _, name := range builtinProviders {
			if providerConfig, ok := s.Providers[name]; ok {
				fmt.Printf("  %-15s -> %s (rules: %s)\n", name, providerConfig.ConfigFile, providerConfig.RulesDir)
			}
		}

		// List custom providers
		customProviders, err := config.ListProviders()
		if err == nil && len(customProviders) > 0 {
			fmt.Println("\nCustom Providers:")
			for _, p := range customProviders {
				if p.Enabled {
					fmt.Printf("  %-15s -> %s (rules: %s)\n", p.Name, p.ConfigFile, p.RulesDir)
				}
			}
		} else {
			fmt.Println("\nCustom Providers: (none)")
		}
	},
}

func init() {
	rootCmd.AddCommand(providersCmd)
	providersCmd.AddCommand(addProviderCmd)
	providersCmd.AddCommand(removeProviderCmd)
	providersCmd.AddCommand(listProvidersCmd)
}

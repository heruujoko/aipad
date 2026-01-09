package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	MarkerStart = "<!-- AIPAD_CONTEXT_START -->"
	MarkerEnd   = "<!-- AIPAD_CONTEXT_END -->"
)

// AgentAwarenessInstructions contains the text to inject into config files
const AgentAwarenessInstructions = `## AIPad Context Management

This project uses **AIPad** for context switching between AI assistants.

### How to Save Context
When you complete a significant task or conversation milestone, save the context using:
` + "```bash" + `
aipad convo "Summary of what was accomplished"
` + "```" + `

### When to Save
- After completing a feature or bug fix
- Before switching to a different topic
- When the user requests a context save
- At natural conversation breakpoints

### Reading Context
The shared scratchpad is located at ` + "`.aipad/scratchpad.md`" + `. Review it to understand prior context.
`

// EnsureRulesDir creates the provider's rules directory if it doesn't exist
func EnsureRulesDir(rulesDir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(cwd, rulesDir)
	return os.MkdirAll(fullPath, 0755)
}

// CopyScratchpadToRules copies the scratchpad to the provider's rules directory
func CopyScratchpadToRules(scratchpadPath, rulesDir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Read scratchpad content
	content, err := os.ReadFile(scratchpadPath)
	if err != nil {
		return fmt.Errorf("failed to read scratchpad: %w", err)
	}

	// Write to rules directory
	destPath := filepath.Join(cwd, rulesDir, "scratchpad.md")
	return os.WriteFile(destPath, content, 0644)
}

// UpdateConfigWithManagedBlock updates the config file with managed block content
func UpdateConfigWithManagedBlock(configPath string, newContent string) error {
	// Read existing content
	existingContent, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new file with markers
			content := fmt.Sprintf("%s\n%s\n%s\n", MarkerStart, newContent, MarkerEnd)
			return os.WriteFile(configPath, []byte(content), 0644)
		}
		return err
	}

	contentStr := string(existingContent)

	// Check if markers exist
	if strings.Contains(contentStr, MarkerStart) && strings.Contains(contentStr, MarkerEnd) {
		// Replace content between markers
		pattern := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(MarkerStart) + `.*?` + regexp.QuoteMeta(MarkerEnd))
		replacement := fmt.Sprintf("%s\n%s\n%s", MarkerStart, newContent, MarkerEnd)
		contentStr = pattern.ReplaceAllString(contentStr, replacement)
	} else {
		// Append markers and content at the end
		contentStr = contentStr + "\n" + MarkerStart + "\n" + newContent + "\n" + MarkerEnd + "\n"
	}

	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

// SyncProviderConfig syncs the scratchpad content into the provider's config file managed block
func SyncProviderConfig(configPath, scratchpadPath string) error {
	// Read scratchpad content
	scratchpadContent, err := os.ReadFile(scratchpadPath)
	if err != nil {
		return fmt.Errorf("failed to read scratchpad: %w", err)
	}

	// Build the managed block content
	managedContent := AgentAwarenessInstructions + "\n## Current Session Context\n\n" + string(scratchpadContent)

	return UpdateConfigWithManagedBlock(configPath, managedContent)
}

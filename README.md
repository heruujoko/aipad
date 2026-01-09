# AIPad CLI

AIPad is a powerful context management tool designed for developers working with multiple AI assistants (like Claude, Antigravity, etc.). It ensures your conversation context is preserved, deduplicated, and synchronized across different AI platforms seamlessly.

## 🚀 Key Features

- **Multi-Provider Support**: Built-in support for Claude and Antigravity, with easy custom provider configuration.
- **Context Synchronization**: Automatically syncs your conversation "scratchpad" into provider-specific rule files (e.g., `CLAUDE.md`, `AGENTS.md`).
- **Smart Deduplication**: Uses SHA256 hashing and fuzzy matching (>80% similarity) to prevent redundant context from cluttering your files.
- **Provider Switching**: Seamlessly switch between AI assistants while carrying over your relevant context.
- **Managed Blocks**: Injects context into your existing documentation files using non-destructive markers (`<!-- AIPAD_CONTEXT_START -->`).
- **History & Status**: Track your session details and review your conversation history at any time.

## 🛠 Installation

### Via Homebrew (Preferred)

```bash
brew tap heruujoko/tap
brew install aipad
```

### From Source

Ensuring you have Go installed on your system:

```bash
# Clone the repository
git clone https://github.com/heruujoko/aipad.git
cd aipad

# Build the binary
go build -o aipad

# Move to your path (optional)
mv aipad /usr/local/bin/
```

## 📖 Usage

### 1. Initialize a New Session
Start a session with your preferred AI assistant:
```bash
aipad new claude
# or
aipad new ag  # Alias for antigravity
```

### 2. Add Conversation Context
Save important milestones or task summaries to your project's context:
```bash
aipad convo "Implemented the user authentication layer using JWT."
```
*Note: AIPad will automatically reject duplicates or near-duplicate entries.*

### 3. Switch Providers
Switching from Claude to another assistant? AIPad will sync the context to the new provider's rules:
```bash
aipad use ag
```

### 4. Manage Custom Providers
Add your own AI provider configurations:
```bash
aipad providers add my-bot MY_BOT.md .mybot/rules/
aipad providers list
```

### 5. Utility Commands
- **Status**: View current session details.
  ```bash
  aipad status
  ```
- **List**: Review conversation history.
  ```bash
  aipad list
  ```
- **Sync**: Manually force a synchronization.
  ```bash
  aipad sync
  ```
- **Clean**: Remove synced context blocks from your project files.
  ```bash
  aipad clean
  ```
- **Export**: Export context history to Markdown, JSON, or Text.
  ```bash
  aipad export history.json
  ```

## 🏗 Project Structure

```text
.
├── .aipad/
│   ├── state.json          # Session metadata and history
│   └── scratchpad.md       # The master context file
├── CLAUDE.md               # Claude managed block
├── AGENTS.md               # Antigravity managed block
└── ...
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

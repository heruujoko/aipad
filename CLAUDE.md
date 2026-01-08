# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AIPad is a Go-based CLI tool (using Cobra framework) that manages context switching between different AI assistants (Claude, Antigravity, etc.) by preserving conversation context and syncing configuration files across platforms.

## Development Commands

### Build
```bash
go build -o aipad main.go
```

### Run
```bash
go run main.go
```

### Clean build artifacts
```bash
rm -f aipad bin/
```

### Test
```bash
go test ./...
```

### Run single test
```bash
go test -v ./internal/package_name -run TestFunctionName
```

## Architecture

### Directory Structure
```
aipad/
├── cmd/
│   └── root.go             # CLI command definitions (Cobra)
├── internal/
│   ├── state/              # state.json management (not yet implemented)
│   ├── sync/               # scratchpad and config sync logic (not yet implemented)
│   └── crypto/             # Hashing/deduplication (not yet implemented)
├── main.go                 # Entry point
├── .aipad/                 # Runtime state directory (gitignored)
│   ├── state.json          # Session state and metadata
│   └── scratchpad.md       # Shared context scratchpad
├── CLAUDE.md               # Claude-specific config (when aipad-managed)
└── AGENTS.md               # Antigravity-specific config (when aipad-managed)
```

### Key Design Patterns

#### Managed Block Strategy
Config files (CLAUDE.md, AGENTS.md) use marker-based updates to preserve user content:
```markdown
<!-- AIPAD_CONTEXT_START -->
[Dynamic aipad-managed content goes here]
<!-- AIPAD_CONTEXT_END -->
```
This prevents destructive edits when syncing context.

#### State File Schema (.aipad/state.json)
```json
{
  "version": "1.0",
  "current_provider": "claude",
  "session_id": "uuid-here",
  "created_at": "2026-01-08T10:00:00Z",
  "last_sync": "2026-01-08T10:30:00Z",
  "context_hashes": ["hash1", "hash2"],
  "providers": {
    "claude": {
      "config_file": "CLAUDE.md",
      "rules_dir": ".claude/rules/"
    },
    "antigravity": {
      "config_file": "AGENTS.md",
      "rules_dir": ".agent/rules/"
    }
  }
}
```

#### Deduplication
- Content is hashed (MD5/SHA256) before appending to scratchpad
- Hashes stored in state.json prevent duplicate additions
- Fuzzy matching (>80% similarity) for near-duplicate detection

### Planned CLI Commands (from specs.md)

- `aipad new <provider>` - Initialize session with provider
- `aipad convo "<text>"` - Append conversation context to scratchpad
- `aipad use <provider>` - Switch active provider and sync configs
- `aipad status` - Show current provider and session info
- `aipad list` - Show conversation history
- `aipad sync` - Manually trigger context sync
- `aipad clean` - Remove old/duplicate context

### Provider Configuration
- **Claude**: `.claude/rules/` directory, `CLAUDE.md` config file
- **Antigravity**: `.agent/rules/` directory, `AGENTS.md` config file

## Implementation Status

The project is iterating through the implementation phases defined in `specs.md`.

### Completed
- **Phase 1 (Core Infrastructure)**: 
  - Go module and Cobra CLI init
  - `.aipad/` directory structure
  - `state.json` management (creation, loading, saving)
- **Phase 2 (Session Management)**:
  - `aipad new <provider>` command implemented

### In Progress
- **Phase 2**: `aipad convo` command
- **Phase 3**: Deduplication System
- **Phase 4**: Provider Switching

## License

MIT License - See LICENSE file for details.

<!-- AIPAD_CONTEXT_START -->
<!-- AIPAD_CONTEXT_END -->

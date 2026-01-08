# AGENTS.md

This file provides guidance to AI agent assistants when working with code in this repository.

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
## AIPad Context Management

This project uses **AIPad** for context switching between AI assistants.

### How to Save Context
When you complete a significant task or conversation milestone, save the context using:
```bash
./aipad convo "Summary of what was accomplished"
```

### When to Save
- After completing a feature or bug fix
- Before switching to a different topic
- When the user requests a context save
- At natural conversation breakpoints

### Reading Context
The shared scratchpad is located at `.aipad/scratchpad.md`. Review it to understand prior context.

## Current Session Context


## [2026-01-08 23:25:43] Context Update
The project uses Go and Cobra for CLI development.
---

## [2026-01-08 23:26:29] Context Update
We are also using SHA256 for hashing.
---

## [2026-01-08 23:27:26] Context Update
Implemented fuzzy matching for deduplication:
- Added Levenshtein distance and similarity ratio logic in internal/crypto.
- Updated State struct to store ContextHistory (raw text) for similarity checks.
- Updated aipad convo to check for exact hash match AND >80% similarity.
- Added comprehensive unit tests for crypto package.
- Verified functionality with local tests (exact match rejection + fuzzy match rejection).
---

<!-- AIPAD_CONTEXT_END -->
```
This prevents destructive edits when syncing context between AI platforms.

#### State File Schema (.aipad/state.json)
```json
{
  "version": "1.0",
  "current_provider": "antigravity",
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

#### Deduplication Strategy
- Content is hashed (MD5/SHA256) before appending to scratchpad
- Hashes stored in state.json prevent duplicate additions
- Fuzzy matching (>80% similarity) for near-duplicate detection

### Planned CLI Commands

- `aipad new <provider>` - Initialize session with provider
- `aipad convo "<text>"` - Append conversation context to scratchpad
- `aipad use <provider>` - Switch active provider and sync configs
- `aipad status` - Show current provider and session info
- `aipad list` - Show conversation history
- `aipad sync` - Manually trigger context sync
- `aipad clean` - Remove old/duplicate context

### Provider Configuration
- **Antigravity**: `.agent/rules/` directory, `AGENTS.md` config file
- **Claude**: `.claude/rules/` directory, `CLAUDE.md` config file

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

## Context Preservation

When switching between AI assistants, AIPad maintains:
- Shared scratchpad with conversation history
- Provider-specific rules directories
- Synced configuration via managed blocks in config files
- Duplicate detection via content hashing

## License

MIT License - See LICENSE file for details.

<!-- AIPAD_CONTEXT_START -->
## AIPad Context Management

This project uses **AIPad** for context switching between AI assistants.

### How to Save Context
When you complete a significant task or conversation milestone, save the context using:
```bash
./aipad convo "Summary of what was accomplished"
```

### When to Save
- After completing a feature or bug fix
- Before switching to a different topic
- When the user requests a context save
- At natural conversation breakpoints

### Reading Context
The shared scratchpad is located at `.aipad/scratchpad.md`. Review it to understand prior context.

## Current Session Context


## [2026-01-08 23:25:43] Context Update
The project uses Go and Cobra for CLI development.
---

## [2026-01-08 23:26:29] Context Update
We are also using SHA256 for hashing.
---

## [2026-01-08 23:27:26] Context Update
Implemented fuzzy matching for deduplication:
- Added Levenshtein distance and similarity ratio logic in internal/crypto.
- Updated State struct to store ContextHistory (raw text) for similarity checks.
- Updated aipad convo to check for exact hash match AND >80% similarity.
- Added comprehensive unit tests for crypto package.
- Verified functionality with local tests (exact match rejection + fuzzy match rejection).
---

<!-- AIPAD_CONTEXT_END -->

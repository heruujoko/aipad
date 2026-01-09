# AIPad CLI - Product Requirements Document

## Overview
AIPad is a CLI tool that manages context switching between different AI assistants (Claude, Antigravity, etc.) by preserving conversation context and syncing configuration files across different AI platforms.

## Problem Statement
When users run out of AI credits on one platform and switch to another, they lose all conversation context, forcing them to manually recreate the context or start from scratch. This is inefficient and frustrating.

## Solution
A CLI tool that:
- Maintains a shared scratchpad for conversation context
- Automatically syncs platform-specific configuration files
- Detects and prevents duplicate context additions
- Provides seamless switching between AI providers

---

## Core Features

### 1. Session Initialization
**Command:** `aipad new <provider>`

**Behavior:**
- Creates `.aipad/` directory in current working directory if it doesn't exist
- Initializes `scratchpad.md` for the session
- Creates/updates `.aipad/state.json` to track current provider and session metadata
- Sets up provider-specific configuration file if it doesn't exist:
  - Claude: `CLAUDE.md`
  - Antigravity: `AGENTS.md`

**State to Track:**
- Current provider
- Session start time
- Last sync timestamp
- Context checksum/hash for deduplication

---

### 2. Context Saving
**Command:** `aipad convo "<conversation text>"`

**Behavior:**
- Appends conversation context to `scratchpad.md` with timestamp
- Generates a content hash/signature of the added context
- Stores hash in state file to prevent duplicate additions
- Formats context in a structured way:
  ```markdown
  ## [Timestamp] Context Update
  <conversation text>
  ---
  ```

**Deduplication Logic:**
- Before appending, compute hash of new content
- Compare against existing hashes in state file
- If similar content exists (e.g., >80% similarity), skip or merge
- Use simple hash comparison or fuzzy matching

---

### 3. Provider Switching
**Command:** `aipad use <provider>`

**Behavior:**
- Updates current provider in state file
- Links/copies `scratchpad.md` to the target provider's rules directory:
  - Claude: Creates/updates `.claude/rules/` and links scratchpad
  - Antigravity: Creates/updates `.agent/rules/` and links scratchpad
- Maintains a **Managed Block** in the provider-specific config file (CLAUDE.md or AGENTS.md):
  - Instead of rewriting the file, the tool identifies aipad-managed markers.
  - If markers don't exist, it appends them to the end of the file.
  - If they do exist, it only updates the content within those markers.

**Managed Block Strategy:**
- Markers:
  ```markdown
  <!-- AIPAD_CONTEXT_START -->
  [Dynamic content goes here]
  <!-- AIPAD_CONTEXT_END -->
  ```
- This "Safe Append" prevents destructive edits and allows the user to keep their own custom rules outside the block.

---

## Key Action Items for Implementation

### Phase 1: Core Infrastructure (Go)
- [x] **Action 1.1:** Initialize Go module (`go mod init aipad`) and set up Cobra CLI framework.
- [x] **Action 1.2:** Implement `.aipad/` directory initialization and state management.
- [x] **Action 1.3:** Create state.json schema with fields: `current_provider`, `session_id`, `last_sync`, `context_hashes[]`.
- [ ] **Action 1.4:** Set up GitHub Actions workflow for automated cross-compilation (build for Linux, macOS, Windows).

### Phase 2: Session & Context Management
- [x] **Action 2.1:** Implement `aipad new <provider>` command
  - Initialize scratchpad.md
  - Create provider config file if missing
  - Set initial state
- [x] **Action 2.2:** Implement `aipad convo "<text>"` command
  - Append to scratchpad with timestamp
  - Generate content hash (MD5 or SHA256)
  - Store hash in state file
  - Format output consistently

### Phase 3: Deduplication System
- [x] **Action 3.1:** Create hash generation function for conversation text
- [x] **Action 3.2:** Implement similarity detection algorithm (exact match or fuzzy)
- [x] **Action 3.3:** Build deduplication check for both scratchpad.md and config files
- [x] **Action 3.4:** Add distinctive markers to config files for easy detection/replacement

### Phase 4: Provider Switching
- [x] **Action 4.1:** Implement `aipad use <provider>` command
- [x] **Action 4.2:** Create symlink or copy logic for scratchpad → provider rules directory
- [x] **Action 4.3:** Implement config file (CLAUDE.md/AGENTS.md) update with marker-based replacement
- [x] **Action 4.4:** Ensure directory structure creation for `.claude/rules/` and `.agent/rules/`
- [x] **Action 4.5:** Inject "Agent Awareness" instructions into config files (telling the agent to use `aipad convo` for saving context).

### Phase 5: User Experience
- [x] **Action 5.1:** Implement `aipad --version` to display current version from build metadata.
- [x] **Action 5.2:** Add colored output for better readability.
- [x] **Action 5.3:** Implement `aipad status` command to show current provider and session info.
- [x] **Action 5.4:** Add `aipad list` command to show conversation history.
- [x] **Action 5.5:** Create help documentation and usage examples.
- [x] **Action 5.6:** Add error handling and user-friendly error messages.
- [x] **Action 5.7:** add aipad clean to remove conetxt from rules and agets.md and claude.md

### Phase 6: Advanced Features (Optional)
- [x] **Action 6.1:** Implement `aipad sync` to manually trigger context sync
- [x] **Action 6.2:** Support custom provider configurations via config file
- [x] **Action 6.3:** Add `aipad export` to export conversation history

---

## Technical Specifications

### File Structure
```
project/
├── .aipad/
│   ├── state.json          # Session state and metadata
│   └── scratchpad.md       # Shared context scratchpad
├── CLAUDE.md               # Claude-specific config
├── AGENTS.md               # Antigravity-specific config
├── .claude/
│   └── rules/
│       └── scratchpad.md   # Symlink or copy
└── .agent/
    └── rules/
        └── scratchpad.md   # Symlink or copy
```

### State File Schema (state.json)
```json
{
  "version": "1.0",
  "current_provider": "claude",
  "session_id": "uuid-here",
  "created_at": "2026-01-08T10:00:00Z",
  "last_sync": "2026-01-08T10:30:00Z",
  "context_hashes": [
    "hash1",
    "hash2"
  ],
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

### Deduplication Strategy
1. **Hash-based:** Generate MD5/SHA256 hash of normalized content (trim whitespace, lowercase)
2. **Marker-based:** Use HTML-style comments to mark aipad-managed sections in config files
3. **Similarity threshold:** Consider content duplicate if hash matches or >80% similar

---

## Success Criteria
- Users can switch between AI providers without losing context
- No duplicate context additions in config files
- Simple, intuitive CLI commands
- Reliable state management across sessions
- Clear feedback and error messages

## Non-Goals (v1)
- GUI interface
- Cloud sync between machines
- Support for more than 2 providers initially
- Automatic credit monitoring
- Conversation encryption

---

## Distribution & Publishing

### GitHub Actions & GoReleaser
The primary distribution method will be via GitHub Releases.

**Action Items:**
- [x] **Action 7.1:** Create `.goreleaser.yaml` to handle binary compilation and packaging.
- [x] **Action 7.2:** Set up GitHub Actions workflow to trigger GoReleaser on git tags.
- [x] **Action 7.3:** Automate generation of SHA256 hashes and build artifacts for Darwin/Amd64, Darwin/Arm64, Linux, and Windows.
- [x] **Action 7.4:** Create CI workflow (`ci.yml`) to compile and test on every PR update.

### Homebrew Distribution
**Action Items:**
- [ ] **Action 8.1:** Create/Update Homebrew formula to fetch the latest binary from GitHub Releases.
- [ ] **Action 8.2:** Test installation: `brew install yourusername/tap/aipad`.

### NPM Distribution (Optional Wrapper)
**Action Items:**
- [ ] **Action 9.1:** (If needed) Create a thin Node.js wrapper that downloads the appropriate Go binary for the user's platform.

---

## Technical Specifications

### Implementation Approach: Go (Golang)
- **Framework:** Cobra for CLI commands and flags.
- **Persistence:** JSON for `state.json`.
- **Cross-Compilation:** Native Go support + GoReleaser.
- **Safety:** Managed blocks in markdown files to prevent accidental data loss.
- **License:** MIT (Permissive, No Liability).

---

## Project Structure

```
aipad/
├── cmd/
│   └── root.go             # Entry point and command definitions
├── internal/
│   ├── state/              # state.json management
│   ├── sync/               # scratchpad and config logic
│   └── crypto/             # Hashing/deduplication
├── main.go                 # Go main entry
├── .goreleaser.yaml        # Build configuration
├── go.mod                  # Go dependencies
├── LICENSE                 # MIT License
├── CLAUDE.md
└── AGENTS.md
```

### Release Checklist
- [ ] Update version in `main.go` or version file
- [ ] Update `CHANGELOG.md`
- [ ] Git commit and tag: `git tag -a v0.1.0 -m "Release v0.1.0"`
- [ ] Push tag: `git push origin v0.1.0`
- [ ] GitHub Action triggers GoReleaser
- [ ] Verify binaries in GitHub Release
- [ ] Update Homebrew formula (if not automated)
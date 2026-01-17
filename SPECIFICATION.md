# fizz - Fizzy.do CLI Tool

## Overview

`fizz` is a command-line interface for Fizzy.do (https://fizzy.do), a Kanban-style project management tool. It provides a simple, consistent way to interact with all Fizzy API features through the terminal, using the libfizz-go library (https://github.com/visionik/libfizz-go).

**Key Principles:**
- Simple, predictable syntax: `fizz noun verb --flags "data"`
- Unix philosophy: composable, scriptable, works with pipes
- Multiple output formats: human-readable tables and machine-readable JSON/YAML
- Environment-based authentication (12-factor app)
- Comprehensive: supports all 11 Fizzy API services

## Requirements

### Functional Requirements

**MUST have:**
- Support all 11 libfizz-go services (Identity, Boards, Cards, Comments, Reactions, Steps, Tags, Columns, Users, Notifications, Uploads)
- Environment variable authentication (`FIZZY_TOKEN`, `FIZZY_ACCOUNT`)
- Multiple output formats (table, JSON, YAML)
- Shell completion generation (bash, zsh, fish)
- Help system with `--ai-help` flag mentioned
- Auto-pagination for list operations with optional `--limit`
- Accept both card IDs (UUID) and numbers seamlessly
- Flexible input: flags, stdin, and file input
- Proper error handling with helpful messages and hints
- Exit codes: 0=success, 1=user error, 2=system error

**SHOULD have:**
- Colored output for human-readable formats
- `--debug` flag for verbose output
- Example commands in help text
- Real API integration tests with ephemeral test board

**MUST NOT:**
- Require config files (env vars only)
- Make assumptions about user's workflow
- Show stack traces without `--debug` flag

### Non-Functional Requirements

- **Performance:** Commands MUST complete within 5 seconds for typical operations
- **Reliability:** 80%+ test coverage, integration tests against real API
- **Usability:** Zero-config operation with just env vars set
- **Compatibility:** Support Go 1.21+, work on macOS/Linux/Windows
- **Maintainability:** Standard Go project layout, clear separation of concerns

## Architecture

### Technology Stack

- **Language:** Go 1.21+
- **CLI Framework:** cobra (command structure, completion)
- **HTTP Client:** libfizz-go (https://github.com/visionik/libfizz-go)
- **Output Formatting:** 
  - Tables: olekukonko/tablewriter or charmbracelet/lipgloss
  - JSON: encoding/json (stdlib)
  - YAML: gopkg.in/yaml.v3
- **Build System:** Task (https://taskfile.dev)
- **Release:** goreleaser

### Project Structure

```
fizz/
├── cmd/
│   ├── root.go           # Root command, global flags
│   ├── identity.go       # Identity commands
│   ├── boards.go         # Board commands
│   ├── cards.go          # Card commands
│   ├── comments.go       # Comment commands
│   ├── reactions.go      # Reaction commands
│   ├── steps.go          # Step commands
│   ├── tags.go           # Tag commands
│   ├── columns.go        # Column commands
│   ├── users.go          # User commands
│   ├── notifications.go  # Notification commands
│   ├── uploads.go        # Upload commands
│   └── completion.go     # Shell completion commands
├── internal/
│   ├── client/
│   │   ├── client.go     # Fizzy client wrapper
│   │   └── resolver.go   # Card ID/number resolver
│   ├── format/
│   │   ├── formatter.go  # Output formatter interface
│   │   ├── table.go      # Table formatter
│   │   ├── json.go       # JSON formatter
│   │   └── yaml.go       # YAML formatter
│   ├── input/
│   │   ├── parser.go     # Input parser (flags/stdin/file)
│   │   └── stdin.go      # Stdin reader
│   └── config/
│       └── env.go        # Environment variable handling
├── tests/
│   ├── integration/      # Real API tests
│   └── mocks/            # Mock HTTP server tests
├── main.go               # Entry point
├── go.mod
├── go.sum
├── Taskfile.yml          # Build tasks
├── .goreleaser.yml       # Release configuration
├── README.md
└── LICENSE
```

### Command Structure

All commands follow pattern: `fizz <noun> <verb> [flags] [args]`

**Nouns (Services):**
- `identity` - Identity operations
- `boards` - Board management
- `cards` - Card operations
- `comments` - Comment management
- `reactions` - Reactions on comments
- `steps` - Checklist steps
- `tags` - Tag management
- `columns` - Column operations
- `users` - User listing
- `notifications` - Notification management
- `uploads` - File upload operations

**Common Verbs:**
- `list` - List resources
- `get` - Get single resource
- `create` - Create new resource
- `update` - Update existing resource
- `delete` - Delete resource
- Service-specific: `close`, `reopen`, `assign`, etc.

**Global Flags:**
- `--format` - Output format: table (default), json, yaml
- `--debug` - Enable debug output
- `--help, -h` - Show help (mentions --ai-help)
- `--ai-help` - Show AI-powered help

**Common Flags:**
- `--limit` - Limit results for list operations
- `--input, -i` - Input from file
- `-` - Read from stdin

### Authentication Flow

1. Check `FIZZY_TOKEN` environment variable (MUST be set)
2. Check `FIZZY_ACCOUNT` environment variable (MUST be set)
3. If either missing, print helpful error with setup instructions
4. Create libfizz-go client with credentials
5. All commands use shared client instance

### Card ID/Number Resolution

Since Fizzy API uses card numbers (integers) in URLs but cards also have UUIDs:

1. Parse user input (could be number or UUID)
2. If looks like number (digits only), use directly
3. If looks like UUID, call `cards list` to find matching card
4. Extract number from matched card
5. Use number for API operation
6. Cache ID→Number mappings in memory for session

**Trade-off:** Slight performance cost for better UX

## Implementation Plan

### Phase 1: Foundation

#### Subphase 1.1: Project Setup
**Dependencies:** None

**Tasks:**

**Task 1.1.1: Initialize Go module and dependencies**
- Create `go.mod` with Go 1.21+
- Add dependencies:
  - `github.com/visionik/libfizz-go/fizzy`
  - `github.com/spf13/cobra`
  - `gopkg.in/yaml.v3`
  - `github.com/olekukonko/tablewriter` OR `github.com/charmbracelet/lipgloss`
- Add dev dependencies:
  - `github.com/stretchr/testify`
- Acceptance: `go mod tidy` succeeds, no errors

**Task 1.1.2: Create project structure**
- Create directory structure as specified in Architecture
- Create placeholder files with package declarations
- Create `main.go` with basic cobra setup
- Acceptance: `go build` succeeds

**Task 1.1.3: Setup Taskfile**
- Create `Taskfile.yml` with tasks:
  - `build` - Build binary
  - `test` - Run unit tests
  - `test:coverage` - Generate coverage report
  - `test:integration` - Run integration tests
  - `lint` - Run linters (go vet, go fmt)
  - `check` - Pre-commit checks (fmt, lint, test)
- Acceptance: `task --list` shows all tasks

#### Subphase 1.2: Core Infrastructure
**Dependencies:** Subphase 1.1

**Task 1.2.1: Implement environment config**
- File: `internal/config/env.go`
- Read `FIZZY_TOKEN` and `FIZZY_ACCOUNT` from env
- Return clear error if either missing
- Include setup instructions in error message
- Write unit tests
- Acceptance: Tests pass, helpful error shown when vars missing

**Task 1.2.2: Implement Fizzy client wrapper**
- File: `internal/client/client.go`
- Wrap libfizz-go client creation
- Handle authentication
- Support debug mode (verbose logging)
- Write unit tests with mock env vars
- Acceptance: Client creates successfully with valid creds

**Task 1.2.3: Implement output formatters**
- Files: `internal/format/*.go`
- Interface: `Formatter` with `Format(data interface{}) error`
- Implementations:
  - `TableFormatter` - Pretty tables with tablewriter
  - `JSONFormatter` - Standard JSON output
  - `YAMLFormatter` - YAML output
- Factory function: `NewFormatter(format string) Formatter`
- Write unit tests for each formatter
- Acceptance: All formatters produce correct output

**Task 1.2.4: Implement input parser**
- Files: `internal/input/*.go`
- Parse flags into structs
- Detect stdin input (check if `-` provided or piped)
- Read from files if `--input` provided
- Merge flag values with stdin/file JSON
- Write unit tests
- Acceptance: All input methods work, precedence correct

**Task 1.2.5: Setup root command**
- File: `cmd/root.go`
- Create root cobra command
- Add global flags: `--format`, `--debug`
- Show help by default if no command specified
- Mention `--ai-help` in help text
- Implement `--ai-help` flag (print extended help message)
- Acceptance: `fizz` shows help, mentions --ai-help

### Phase 2: Identity Service
**Dependencies:** Phase 1

#### Subphase 2.1: Identity Commands
**Dependencies:** None

**Task 2.1.1: Implement identity get command**
- File: `cmd/identity.go`
- Command: `fizz identity get`
- Call `client.Identity.Get()`
- Format output (show accounts, user info)
- Handle errors gracefully
- Write unit tests
- Acceptance: Command shows user identity

**Task 2.1.2: Add completion for identity**
- Generate shell completion for identity commands
- Test: `fizz identity <TAB>` suggests `get`
- Acceptance: Completion works in bash/zsh

**Task 2.1.3: Integration test for identity**
- File: `tests/integration/identity_test.go`
- Test against real Fizzy API
- Verify identity data returned
- Acceptance: Test passes with real credentials

### Phase 3: Boards Service
**Dependencies:** Phase 1

#### Subphase 3.1: Board Commands
**Dependencies:** None

**Task 3.1.1: Implement boards list**
- File: `cmd/boards.go`
- Command: `fizz boards list`
- Auto-fetch all pages (use `ListAll`)
- Support `--limit` flag
- Format as table/JSON/YAML
- Write unit tests
- Acceptance: Lists all boards

**Task 3.1.2: Implement boards get**
- Command: `fizz boards get <board-id>`
- Show single board details
- Write unit tests
- Acceptance: Shows board info

**Task 3.1.3: Implement boards create**
- Command: `fizz boards create --name="Name" [--description="Desc"]`
- Support stdin JSON input
- Support `--input=file.json`
- Show created board
- Write unit tests
- Acceptance: Creates board successfully

**Task 3.1.4: Implement boards update**
- Command: `fizz boards update <board-id> [--name=""] [--description=""]`
- Support partial updates
- Write unit tests
- Acceptance: Updates board

**Task 3.1.5: Implement boards delete**
- Command: `fizz boards delete <board-id>`
- Confirm deletion (or add `--force` to skip)
- Write unit tests
- Acceptance: Deletes board

**Task 3.1.6: Integration tests for boards**
- File: `tests/integration/boards_test.go`
- Create test board
- List, get, update operations
- Delete test board
- Acceptance: All operations work against real API

### Phase 4: Cards Service
**Dependencies:** Phase 1

#### Subphase 4.1: Card ID Resolution
**Dependencies:** None

**Task 4.1.1: Implement card resolver**
- File: `internal/client/resolver.go`
- Function: `ResolveCard(input string) (cardNumber string, err error)`
- Detect if input is number or UUID
- If UUID, search cards to find number
- Cache results in memory map
- Write unit tests
- Acceptance: Both IDs and numbers work

#### Subphase 4.2: Core Card Commands
**Dependencies:** Subphase 4.1

**Task 4.2.1: Implement cards list**
- File: `cmd/cards.go`
- Command: `fizz cards list [--board=ID] [--status=open] [--limit=N]`
- Auto-fetch all pages
- Support filtering by board, status, tags, column
- Write unit tests
- Acceptance: Lists cards with filters

**Task 4.2.2: Implement cards get**
- Command: `fizz cards get <card-id-or-number>`
- Use resolver to handle both ID and number
- Show full card details
- Write unit tests
- Acceptance: Works with both ID and number

**Task 4.2.3: Implement cards create**
- Command: `fizz cards create --board=ID --title="Title" [--body="Body"]`
- Support stdin/file input for complex data
- Show created card with both ID and number
- Write unit tests
- Acceptance: Creates card, returns both identifiers

**Task 4.2.4: Implement cards update**
- Command: `fizz cards update <card-id-or-number> [--title=""] [--body=""]`
- Support partial updates
- Write unit tests
- Acceptance: Updates card

**Task 4.2.5: Implement cards delete**
- Command: `fizz cards delete <card-id-or-number>`
- Write unit tests
- Acceptance: Deletes card

#### Subphase 4.3: Card Actions
**Dependencies:** Subphase 4.2

**Task 4.3.1: Implement card state actions**
- Commands:
  - `fizz cards close <card-id-or-number>`
  - `fizz cards reopen <card-id-or-number>`
  - `fizz cards postpone <card-id-or-number>`
  - `fizz cards triage <card-id-or-number>`
- Write unit tests for each
- Acceptance: All state transitions work

**Task 4.3.2: Implement card assignment/tagging**
- Commands:
  - `fizz cards assign <card-id-or-number> <user-id>`
  - `fizz cards tag <card-id-or-number> <tag-name>`
  - `fizz cards move <card-id-or-number> --column=<column-id>`
- Write unit tests
- Acceptance: All operations work

**Task 4.3.3: Implement card watching**
- Commands:
  - `fizz cards watch <card-id-or-number>`
  - `fizz cards unwatch <card-id-or-number>`
- Write unit tests
- Acceptance: Watch/unwatch work

**Task 4.3.4: Implement golden card operations**
- Commands:
  - `fizz cards golden <card-id-or-number>` - Mark as golden
  - `fizz cards ungolden <card-id-or-number>` - Remove golden
- Write unit tests
- Acceptance: Golden status toggles

**Task 4.3.5: Integration tests for cards**
- File: `tests/integration/cards_test.go`
- Create test board, create cards
- Test all operations (CRUD, actions, etc.)
- Clean up test data
- Acceptance: All card operations work with real API

### Phase 5: Comments Service
**Dependencies:** Phase 1, Phase 4 (for card operations)

#### Subphase 5.1: Comment Commands
**Dependencies:** None

**Task 5.1.1: Implement comments list**
- File: `cmd/comments.go`
- Command: `fizz comments list <card-id-or-number>`
- Show all comments for card
- Write unit tests
- Acceptance: Lists comments

**Task 5.1.2: Implement comments create**
- Command: `fizz comments create <card-id-or-number> --body="Comment text"`
- Support stdin input for long comments
- Write unit tests
- Acceptance: Creates comment

**Task 5.1.3: Implement comments update**
- Command: `fizz comments update <card-id-or-number> <comment-id> --body="Updated"`
- Write unit tests
- Acceptance: Updates comment

**Task 5.1.4: Implement comments delete**
- Command: `fizz comments delete <card-id-or-number> <comment-id>`
- Write unit tests
- Acceptance: Deletes comment

**Task 5.1.5: Integration tests**
- File: `tests/integration/comments_test.go`
- Test all comment operations on test card
- Acceptance: All operations work

### Phase 6: Reactions Service
**Dependencies:** Phase 1, Phase 5

#### Subphase 6.1: Reaction Commands
**Dependencies:** None

**Task 6.1.1: Implement reactions list**
- File: `cmd/reactions.go`
- Command: `fizz reactions list <card-id-or-number> <comment-id>`
- Show all reactions on comment
- Write unit tests
- Acceptance: Lists reactions

**Task 6.1.2: Implement reactions create**
- Command: `fizz reactions create <card-id-or-number> <comment-id> --emoji="👍"`
- Write unit tests
- Acceptance: Adds reaction

**Task 6.1.3: Implement reactions delete**
- Command: `fizz reactions delete <card-id-or-number> <comment-id> <reaction-id>`
- Write unit tests
- Acceptance: Removes reaction

**Task 6.1.4: Integration tests**
- File: `tests/integration/reactions_test.go`
- Test reaction operations
- Acceptance: All operations work

### Phase 7: Steps Service
**Dependencies:** Phase 1, Phase 4

#### Subphase 7.1: Step Commands
**Dependencies:** None

**Task 7.1.1: Implement steps list**
- File: `cmd/steps.go`
- Command: `fizz steps list <card-id-or-number>`
- Show all checklist steps
- Write unit tests
- Acceptance: Lists steps

**Task 7.1.2: Implement steps get**
- Command: `fizz steps get <card-id-or-number> <step-id>`
- Show single step
- Write unit tests
- Acceptance: Shows step details

**Task 7.1.3: Implement steps create**
- Command: `fizz steps create <card-id-or-number> --content="Step text" [--completed=true]`
- Write unit tests
- Acceptance: Creates step

**Task 7.1.4: Implement steps update**
- Command: `fizz steps update <card-id-or-number> <step-id> [--content=""] [--completed=true]`
- Write unit tests
- Acceptance: Updates step

**Task 7.1.5: Implement steps delete**
- Command: `fizz steps delete <card-id-or-number> <step-id>`
- Write unit tests
- Acceptance: Deletes step

**Task 7.1.6: Integration tests**
- File: `tests/integration/steps_test.go`
- Test all step operations
- Acceptance: All operations work

### Phase 8: Tags Service
**Dependencies:** Phase 1

#### Subphase 8.1: Tag Commands
**Dependencies:** None

**Task 8.1.1: Implement tags list**
- File: `cmd/tags.go`
- Command: `fizz tags list`
- Show all account tags
- Write unit tests
- Acceptance: Lists tags

**Task 8.1.2: Implement tags create**
- Command: `fizz tags create --name="bug" [--color="#ff0000"]`
- Write unit tests
- Acceptance: Creates tag

**Task 8.1.3: Integration tests**
- File: `tests/integration/tags_test.go`
- Test tag operations
- Acceptance: All operations work

### Phase 9: Columns Service
**Dependencies:** Phase 1, Phase 3

#### Subphase 9.1: Column Commands
**Dependencies:** None

**Task 9.1.1: Implement columns list**
- File: `cmd/columns.go`
- Command: `fizz columns list <board-id>`
- Show all columns for board
- Write unit tests
- Acceptance: Lists columns

**Task 9.1.2: Implement columns get**
- Command: `fizz columns get <board-id> <column-id>`
- Write unit tests
- Acceptance: Shows column

**Task 9.1.3: Implement columns create**
- Command: `fizz columns create <board-id> --name="Column Name"`
- Write unit tests
- Acceptance: Creates column

**Task 9.1.4: Implement columns update**
- Command: `fizz columns update <board-id> <column-id> [--name=""] [--position=N]`
- Write unit tests
- Acceptance: Updates column

**Task 9.1.5: Implement columns delete**
- Command: `fizz columns delete <board-id> <column-id>`
- Write unit tests
- Acceptance: Deletes column

**Task 9.1.6: Integration tests**
- File: `tests/integration/columns_test.go`
- Test column operations on test board
- Acceptance: All operations work

### Phase 10: Users Service
**Dependencies:** Phase 1

#### Subphase 10.1: User Commands
**Dependencies:** None

**Task 10.1.1: Implement users list**
- File: `cmd/users.go`
- Command: `fizz users list`
- Show all account users
- Write unit tests
- Acceptance: Lists users

**Task 10.1.2: Integration tests**
- File: `tests/integration/users_test.go`
- Test user listing
- Acceptance: Lists users successfully

### Phase 11: Notifications Service
**Dependencies:** Phase 1

#### Subphase 11.1: Notification Commands
**Dependencies:** None

**Task 11.1.1: Implement notifications list**
- File: `cmd/notifications.go`
- Command: `fizz notifications list`
- Show all notifications
- Write unit tests
- Acceptance: Lists notifications

**Task 11.1.2: Implement notification actions**
- Commands:
  - `fizz notifications read <notification-id>`
  - `fizz notifications unread <notification-id>`
  - `fizz notifications read-all`
- Write unit tests
- Acceptance: All actions work

**Task 11.1.3: Integration tests**
- File: `tests/integration/notifications_test.go`
- Test notification operations
- Acceptance: All operations work

### Phase 12: Uploads Service
**Dependencies:** Phase 1

#### Subphase 12.1: Upload Commands
**Dependencies:** None

**Task 12.1.1: Implement upload command**
- File: `cmd/uploads.go`
- Command: `fizz uploads create <file-path>`
- Use `client.Uploads.UploadFile()`
- Return upload URL/metadata
- Write unit tests (with test file)
- Acceptance: Uploads file successfully

**Task 12.1.2: Integration tests**
- File: `tests/integration/uploads_test.go`
- Upload test file, verify response
- Acceptance: Upload works with real API

### Phase 13: Shell Completion
**Dependencies:** All service commands (Phase 2-12)

#### Subphase 13.1: Completion Implementation
**Dependencies:** None

**Task 13.1.1: Implement completion command**
- File: `cmd/completion.go`
- Commands:
  - `fizz completion bash`
  - `fizz completion zsh`
  - `fizz completion fish`
- Generate using cobra's built-in completion
- Include installation instructions in output
- Write tests
- Acceptance: Completion scripts generate

**Task 13.1.2: Test completion**
- Manually test completion in bash/zsh
- Verify tab-completion works for:
  - Command names
  - Subcommands
  - Flags
- Acceptance: Completion works as expected

### Phase 14: Polish & Documentation
**Dependencies:** Phase 13

#### Subphase 14.1: User Experience
**Dependencies:** None

**Task 14.1.1: Add colors to table output**
- Use color library (e.g., fatih/color)
- Colorize status fields, headers
- Respect `NO_COLOR` env var
- Acceptance: Tables are colorized appropriately

**Task 14.1.2: Improve error messages**
- Review all error messages
- Add context and hints where helpful
- Example: "Card not found. Did you mean card number instead of ID?"
- Acceptance: Errors are clear and actionable

**Task 14.1.3: Add examples to help text**
- Add real examples to each command's help
- Include common use cases
- Acceptance: Help text includes examples

#### Subphase 14.2: Documentation
**Dependencies:** None

**Task 14.2.1: Write comprehensive README**
- Installation instructions (all methods)
- Quick start guide
- Authentication setup
- Command reference or link to generated docs
- Examples for common workflows
- Acceptance: README is complete and clear

**Task 14.2.2: Create CONTRIBUTING.md**
- Development setup instructions
- Testing guidelines
- Code style expectations
- PR process
- Acceptance: Contributors can get started easily

**Task 14.2.3: Add LICENSE**
- Choose license (MIT recommended)
- Add LICENSE file
- Acceptance: Project is properly licensed

### Phase 15: Release Infrastructure
**Dependencies:** Phase 14

#### Subphase 15.1: Release Automation
**Dependencies:** None

**Task 15.1.1: Setup goreleaser**
- Create `.goreleaser.yml`
- Configure builds for:
  - macOS (amd64, arm64)
  - Linux (amd64, arm64)
  - Windows (amd64)
- Create archives (tar.gz, zip)
- Generate checksums
- Write release notes template
- Acceptance: `goreleaser check` passes

**Task 15.1.2: Create GitHub Actions workflow**
- File: `.github/workflows/release.yml`
- Trigger on tag push (v*)
- Run tests first
- Use goreleaser action
- Upload artifacts to GitHub Releases
- Acceptance: Workflow defined and valid

**Task 15.1.3: Test release process**
- Create test tag
- Run goreleaser locally
- Verify artifacts created
- Acceptance: Release artifacts look correct

**Task 15.1.4: Document release process**
- Add to README or CONTRIBUTING
- Steps to create a release
- Versioning scheme (semver)
- Acceptance: Release process documented

## Testing Strategy

### Unit Tests

**Coverage Target:** ≥80%

**Scope:**
- All internal packages (client, format, input, config)
- Command flag parsing logic
- Error handling paths
- Input parsing (flags, stdin, files)
- Output formatting (table, JSON, YAML)

**Tools:**
- `testing` package (stdlib)
- `github.com/stretchr/testify/assert` for assertions
- `github.com/stretchr/testify/mock` for mocks

**Run:** `task test` or `go test ./...`

### Integration Tests

**Scope:**
- All service commands against real Fizzy API
- Create ephemeral test board: `fizz-test-{timestamp}`
- Perform operations (create, update, delete)
- Clean up test board after tests
- Verify operations via libfizz-go directly

**Requirements:**
- `FIZZY_TOKEN` and `FIZZY_ACCOUNT` env vars must be set
- Use build tag: `-tags=integration`

**Files:** `tests/integration/*_test.go`

**Run:** `task test:integration` or `go test -tags=integration ./tests/integration/...`

**Cleanup Strategy:**
- Use `t.Cleanup()` to ensure board deletion
- If tests fail, board still gets deleted
- Add timestamp to board name for easy identification

### Manual Testing

**Pre-release Checklist:**
1. Test all commands with real data
2. Verify output formats (table, JSON, YAML)
3. Test error scenarios (invalid token, missing board, etc.)
4. Test shell completion in bash and zsh
5. Test on macOS and Linux
6. Verify `--ai-help` and `--debug` flags work

## Deployment

### Installation Methods

#### Go Install
```bash
go install github.com/visionik/fizz@latest
```

#### Pre-built Binaries
Download from GitHub Releases:
```bash
# macOS (Homebrew in future)
curl -L https://github.com/visionik/fizz/releases/latest/download/fizz_darwin_amd64.tar.gz | tar xz
sudo mv fizz /usr/local/bin/

# Linux
curl -L https://github.com/visionik/fizz/releases/latest/download/fizz_linux_amd64.tar.gz | tar xz
sudo mv fizz /usr/local/bin/

# Windows (via PowerShell)
# Download from releases page and add to PATH
```

#### Future: Homebrew
```bash
brew install visionik/tap/fizz
```

### Release Process

1. Ensure all tests pass: `task check`
2. Update CHANGELOG.md
3. Create git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
4. Push tag: `git push origin v1.0.0`
5. GitHub Actions runs goreleaser
6. Binaries published to GitHub Releases
7. Announce release

### Versioning

Follow Semantic Versioning (semver):
- `v1.0.0` - Major release (breaking changes)
- `v1.1.0` - Minor release (new features, backwards compatible)
- `v1.0.1` - Patch release (bug fixes)

## Out of Scope

The following are explicitly NOT included in this specification:

- Config file support (only env vars)
- Interactive TUI/browser
- Local caching of API data
- Offline mode
- Plugin system
- Desktop notifications
- Webhooks/automation
- Brew tap maintenance (until project matures)

## Open Questions & Risks

### Risks

1. **Card ID/Number confusion** - Mitigated by flexible resolver
2. **API rate limits** - Mitigated by libfizz-go's retry logic
3. **Large datasets** - Auto-pagination might be slow; `--limit` provides escape hatch
4. **Breaking API changes** - Mitigated by using stable libfizz-go library

### Future Enhancements

Consider after v1.0:
- Config file support (if users request)
- Homebrew tap
- Shell aliases/shortcuts
- Watch mode for notifications
- Bulk operations
- Templates for common workflows

---

## Interview Summary

### Questions & Answers

1. **CLI Command Structure** → Option 1: Direct mapping (fizz boards list, fizz cards create)
2. **Authentication** → Option 1: Environment variables only (FIZZY_TOKEN, FIZZY_ACCOUNT)
3. **Output Format** → Option 3: Multiple formats (--format=table/json/yaml)
4. **Input Handling** → Option 3: Flags with stdin/file support
5. **Error Handling** → Option 2: Helpful errors with hints, --debug flag
6. **Pagination** → Option 1: Auto-fetch all with --limit option
7. **Shell Completion** → Option 2: Help + completion scripts + --ai-help flag (show help by default)
8. **Card ID vs Number** → Option 2: Flexible - accept both, resolve internally
9. **Testing** → Option 2: Unit + Integration with real API (ephemeral test board)
10. **Installation** → Option 3: go install + GitHub releases + future Homebrew
11. **Dashdash Integration** → Just --ai-help flag mentioned in help
12. **Implementation Timeline** → Option 2: Feature complete (all 11 services)
13. **Project Structure** → Option 2: Standard Go CLI layout

---

## Next Steps

Type `implement SPECIFICATION.md` to begin implementation following this specification.

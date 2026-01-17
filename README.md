# fizz - Fizzy.do CLI

A simple, powerful command-line interface for [Fizzy.do](https://fizzy.do), the Kanban project management tool.

Built with [libfizz-go](https://github.com/visionik/libfizz-go), the official Go library for the Fizzy.do API.

## Features

- **Complete API Coverage**: All 11 Fizzy services supported (Identity, Boards, Cards, Comments, Reactions, Steps, Tags, Columns, Users, Notifications, Uploads)
- **Multiple Output Formats**: Table (default), JSON, and YAML
- **Flexible Authentication**: Environment variables only (12-factor app)
- **Shell Completion**: Bash, Zsh, and Fish
- **AI Help**: Built-in `--ai-help` flag for enhanced assistance
- **Simple Syntax**: `fizz noun verb --flags`

## Installation

### Using go install

```bash
go install github.com/visionik/fizz@latest
```

### Using pre-built binaries

Download from [GitHub Releases](https://github.com/visionik/fizz/releases):

```bash
# macOS
curl -L https://github.com/visionik/fizz/releases/latest/download/fizz_darwin_amd64.tar.gz | tar xz
sudo mv fizz /usr/local/bin/

# Linux
curl -L https://github.com/visionik/fizz/releases/latest/download/fizz_linux_amd64.tar.gz | tar xz
sudo mv fizz /usr/local/bin/
```

## Quick Start

### 1. Set up authentication

```bash
export FIZZY_TOKEN="your-api-token"
export FIZZY_ACCOUNT="your-account-id"
```

Add these to your `~/.zshrc` or `~/.bashrc` to make them permanent.

### 2. Try some commands

```bash
# Get your identity
fizz identity get

# List all boards
fizz boards list

# Create a card
fizz cards create --board=123 --title="Fix bug" --body="Description here"

# List cards as JSON
fizz cards list --format=json
```

## Usage

### Common Commands

#### Boards
```bash
fizz boards list
fizz boards get <board-id>
fizz boards create --name="My Board"
fizz boards update <board-id> --name="Updated"
fizz boards delete <board-id>
```

#### Cards
```bash
fizz cards list
fizz cards list --limit=10
fizz cards get <card-id>
fizz cards create --board=<id> --title="Title"
fizz cards update <card-id> --title="New title"
fizz cards delete <card-id>

# Card actions
fizz cards close <card-id>
fizz cards reopen <card-id>
fizz cards assign <card-id> <user-id>
fizz cards tag <card-id> <tag-name>
fizz cards move <card-id> --column=<column-id>
fizz cards watch <card-id>
fizz cards golden <card-id>
```

#### Comments
```bash
fizz comments list <card-id>
fizz comments create <card-id> --body="Great work!"
fizz comments update <card-id> <comment-id> --body="Updated"
fizz comments delete <card-id> <comment-id>
```

#### Other Services
```bash
# Reactions
fizz reactions list <card-id> <comment-id>
fizz reactions create <card-id> <comment-id> --emoji="👍"

# Steps (checklist items)
fizz steps list <card-id>
fizz steps create <card-id> --content="Review code"

# Tags
fizz tags list
fizz tags create --name="bug" --color="#ff0000"

# Columns
fizz columns list <board-id>
fizz columns create <board-id> --name="In Progress"

# Users
fizz users list

# Notifications
fizz notifications list
fizz notifications read <notification-id>
fizz notifications read-all

# Uploads
fizz uploads create ./image.png
```

### Output Formats

```bash
# Human-readable table (default)
fizz boards list

# JSON for scripting
fizz boards list --format=json

# YAML
fizz boards list --format=yaml
```

### Shell Completion

```bash
# Bash
fizz completion bash > /etc/bash_completion.d/fizz

# Zsh
fizz completion zsh > "${fpath[1]}/_fizz"

# Fish
fizz completion fish > ~/.config/fish/completions/fizz.fish
```

### Debug Mode

```bash
fizz --debug boards list
```

### AI Help

```bash
fizz boards --ai-help
```

## Development

### Prerequisites

- Go 1.21+
- Task ([taskfile.dev](https://taskfile.dev))

### Building

```bash
task build
```

### Testing

```bash
# Unit tests
task test

# With coverage
task test:coverage

# Integration tests (requires FIZZY_TOKEN and FIZZY_ACCOUNT)
task test:integration

# All checks
task check
```

### Project Structure

```
fizz/
├── cmd/               # Command implementations
├── internal/
│   ├── client/       # Fizzy client wrapper
│   ├── format/       # Output formatters
│   ├── input/        # Input parsers
│   └── config/       # Configuration
├── tests/
│   └── integration/  # Integration tests
├── main.go
├── Taskfile.yml
└── README.md
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `task check` to ensure tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Links

- [Fizzy.do](https://fizzy.do)
- [libfizz-go](https://github.com/visionik/libfizz-go)
- [Issues](https://github.com/visionik/fizz/issues)

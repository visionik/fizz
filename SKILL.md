---
name: fizz
description: Fizzy.do Kanban CLI for managing boards, cards, comments, tags, and more.
---

# fizz - Fizzy.do CLI

Manage Fizzy.do Kanban boards via the command line. Supports all API operations with multiple output formats.

## Quick Reference

### Authentication
```bash
export FIZZY_TOKEN="your-api-token"
export FIZZY_ACCOUNT="your-account-id"
```

### Boards
```bash
fizz boards list
fizz boards get <id>
fizz boards create --name="Board Name"
fizz boards update <id> --name="New Name"
fizz boards delete <id>
```

### Cards
```bash
fizz cards list [--board=<id>] [--limit=N]
fizz cards get <id>
fizz cards create --board=<id> --title="Title" [--body="Description"]
fizz cards update <id> [--title="New"] [--body="Updated"]
fizz cards delete <id>

# Actions
fizz cards close <id>
fizz cards reopen <id>
fizz cards assign <card-id> <user-id>
fizz cards tag <card-id> <tag-name>
fizz cards move <card-id> --column=<column-id>
fizz cards watch <card-id>
fizz cards golden <card-id>
```

### Comments
```bash
fizz comments list <card-id>
fizz comments create <card-id> --body="Comment text"
fizz comments update <card-id> <comment-id> --body="Updated"
fizz comments delete <card-id> <comment-id>
```

### Reactions
```bash
fizz reactions list <card-id> <comment-id>
fizz reactions create <card-id> <comment-id> --emoji="👍"
fizz reactions delete <card-id> <comment-id> <reaction-id>
```

### Steps (Checklists)
```bash
fizz steps list <card-id>
fizz steps create <card-id> --content="Task item"
fizz steps update <card-id> <step-id> --content="Updated"
fizz steps check <card-id> <step-id>
fizz steps uncheck <card-id> <step-id>
fizz steps delete <card-id> <step-id>
```

### Tags
```bash
fizz tags list
fizz tags get <id>
fizz tags create --name="bug" --color="#ff0000"
fizz tags update <id> --name="feature"
fizz tags delete <id>
```

### Columns
```bash
fizz columns list <board-id>
fizz columns create <board-id> --name="In Progress"
fizz columns update <board-id> <column-id> --name="Done"
fizz columns delete <board-id> <column-id>
```

### Users
```bash
fizz users list
fizz users get <id>
```

### Notifications
```bash
fizz notifications list
fizz notifications read <id>
fizz notifications read-all
```

### Uploads
```bash
fizz uploads create ./image.png
```

### Identity
```bash
fizz identity get
```

## Output Formats

```bash
fizz boards list              # Table (default)
fizz boards list --format=json
fizz boards list --format=yaml
```

## Flags

| Flag | Description |
|------|-------------|
| `--format` | Output format: table, json, yaml |
| `--debug` | Enable debug logging |
| `--ai-help` | Show AI/LLM usage guidance |

## Tips

- Use `--format=json` for scripting and piping to `jq`
- Card IDs are returned when creating; save them for subsequent operations
- Use `fizz <command> --ai-help` for detailed guidance on any command

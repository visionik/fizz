package aihelp

const Content = `---
abt: https://github.com/visionik/dashdash
sub: false
ver: https://github.com/visionik/dashdash/blob/main/README.md#v1.0.0
acl: interact
web: https://fizzy.do
cli: https://github.com/visionik/fizz
mcp: none
api: https://fizzy.do/docs/api
---

# fizz - Fizzy.do CLI AI Help Guide

## Overview

fizz is a command-line interface for Fizzy.do, a Kanban project management tool. It supports all 11 Fizzy API services: Identity, Boards, Cards, Comments, Reactions, Steps, Tags, Columns, Users, Notifications, and Uploads.

## Setup/Prerequisites

### Installation

` + "`" + `` + "`" + `bash
go install github.com/visionik/fizz@latest
` + "`" + `` + "`" + `

### Authentication

fizz requires two environment variables:

1. Get your API token from https://fizzy.do/settings/tokens
2. Get your account ID from https://fizzy.do/settings/account
3. Set environment variables:

` + "`" + `` + "`" + `bash
export FIZZY_TOKEN="your-api-token"
export FIZZY_ACCOUNT="your-account-id"
` + "`" + `` + "`" + `

⚠️ **Never commit these credentials to version control**

## Command Reference

All commands follow the pattern: ` + "`" + `fizz <noun> <verb> [flags]` + "`" + `

### Global Flags

- ` + "`" + `--format FORMAT` + "`" + ` - Output format: table (default), json, yaml
- ` + "`" + `--debug` + "`" + ` - Enable debug output
- ` + "`" + `--ai-help` + "`" + ` - Show this AI-oriented help

### Services (Nouns)

1. **identity** - Get user identity and account information
2. **boards** - Manage boards (list, get, create, update, delete)
3. **cards** - Manage cards with full CRUD + 16 actions
4. **comments** - Manage comments on cards
5. **reactions** - Manage emoji reactions on comments
6. **steps** - Manage checklist steps on cards
7. **tags** - Manage tags
8. **columns** - Manage board columns
9. **users** - List users in account
10. **notifications** - Manage notifications
11. **uploads** - Upload files

## Input Specification

### ID Formats

- **Board IDs**: String format (e.g., "03fbhgtekgu3r5adlafa4qd22")
- **Card IDs**: String format, accepts both UUID and number
- **Column IDs**: String format
- **User IDs**: String format

### Common Parameters

- ` + "`" + `--limit N` + "`" + ` - Limit results (0 = all)
- ` + "`" + `--name "Name"` + "`" + ` - Resource name (use quotes for spaces)
- ` + "`" + `--title "Title"` + "`" + ` - Card title
- ` + "`" + `--body "Text"` + "`" + ` - Body/description text
- ` + "`" + `--board ID` + "`" + ` - Board identifier
- ` + "`" + `--column ID` + "`" + ` - Column identifier

## Output Formats

⚠️ **LLM Best Practice: Always Use --format=json**

✅ RECOMMENDED for programmatic parsing:
` + "`" + `` + "`" + `bash
fizz boards list --format=json
fizz cards get CARD_ID --format=json
` + "`" + `` + "`" + `

❌ NOT RECOMMENDED for parsing:
` + "`" + `` + "`" + `bash
fizz boards list  # Table format, hard to parse reliably
` + "`" + `` + "`" + `

### JSON Output Structure

All list commands return arrays of objects:
` + "`" + `` + "`" + `json
[
  {
    "id": "string",
    "name": "string",
    "created_at": "ISO 8601 datetime",
    ...
  }
]
` + "`" + `` + "`" + `

Single resource commands return objects:
` + "`" + `` + "`" + `json
{
  "id": "string",
  "name": "string",
  ...
}
` + "`" + `` + "`" + `

## Examples

### Identity
` + "`" + `` + "`" + `bash
fizz identity get --format=json
` + "`" + `` + "`" + `

### Boards
` + "`" + `` + "`" + `bash
# List all boards
fizz boards list --format=json

# Get specific board
fizz boards get BOARD_ID --format=json

# Create board
fizz boards create --name="My Board" --format=json

# Create with description
fizz boards create --name="Project" --description="Description here" --format=json

# Update board
fizz boards update BOARD_ID --name="New Name" --format=json

# Delete board
fizz boards delete BOARD_ID
` + "`" + `` + "`" + `

### Cards
` + "`" + `` + "`" + `bash
# List all cards
fizz cards list --format=json

# List limited cards
fizz cards list --limit=10 --format=json

# Get specific card (by ID or number)
fizz cards get CARD_ID --format=json
fizz cards get 123 --format=json

# Create card
fizz cards create --board=BOARD_ID --title="Bug fix" --format=json
fizz cards create --board=BOARD_ID --title="Feature" --body="Description" --format=json

# Update card
fizz cards update CARD_ID --title="New title" --format=json

# Delete card
fizz cards delete CARD_ID

# Card actions
fizz cards close CARD_ID
fizz cards reopen CARD_ID
fizz cards postpone CARD_ID
fizz cards triage CARD_ID
fizz cards assign CARD_ID USER_ID
fizz cards tag CARD_ID "bug"
fizz cards move CARD_ID --column=COLUMN_ID
fizz cards watch CARD_ID
fizz cards unwatch CARD_ID
fizz cards golden CARD_ID
fizz cards ungolden CARD_ID
` + "`" + `` + "`" + `

### Comments
` + "`" + `` + "`" + `bash
# List comments on card
fizz comments list CARD_ID --format=json

# Create comment
fizz comments create CARD_ID --body="Great work!" --format=json

# Update comment
fizz comments update CARD_ID COMMENT_ID --body="Updated text" --format=json

# Delete comment
fizz comments delete CARD_ID COMMENT_ID
` + "`" + `` + "`" + `

### Other Services
` + "`" + `` + "`" + `bash
# Reactions
fizz reactions list CARD_ID COMMENT_ID --format=json
fizz reactions create CARD_ID COMMENT_ID --emoji="👍" --format=json
fizz reactions delete CARD_ID COMMENT_ID REACTION_ID

# Steps (checklist items)
fizz steps list CARD_ID --format=json
fizz steps create CARD_ID --content="Review code" --format=json
fizz steps update CARD_ID STEP_ID --completed=true --format=json
fizz steps delete CARD_ID STEP_ID

# Tags
fizz tags list --format=json
fizz tags create --name="bug" --color="#ff0000" --format=json

# Columns
fizz columns list BOARD_ID --format=json
fizz columns create BOARD_ID --name="In Progress" --format=json
fizz columns update BOARD_ID COLUMN_ID --name="Done" --format=json
fizz columns delete BOARD_ID COLUMN_ID

# Users
fizz users list --format=json

# Notifications
fizz notifications list --format=json
fizz notifications read NOTIFICATION_ID
fizz notifications unread NOTIFICATION_ID
fizz notifications read-all

# Uploads
fizz uploads create ./file.png --format=json
` + "`" + `` + "`" + `

## Troubleshooting

### Error: "FIZZY_TOKEN environment variable is not set"

**Cause:** Missing authentication credentials

**Solution:**
` + "`" + `` + "`" + `bash
export FIZZY_TOKEN="your-token"
export FIZZY_ACCOUNT="your-account-id"
` + "`" + `` + "`" + `

### Error: "failed to ... : 401 Unauthorized"

**Cause:** Invalid or expired API token

**Solution:**
1. Verify token at https://fizzy.do/settings/tokens
2. Generate new token if needed
3. Update environment variable

### Error: "invalid board ID" or "invalid card ID"

**Cause:** Incorrect ID format or non-existent resource

**Solution:**
1. List resources first to get valid IDs:
   ` + "`" + `fizz boards list --format=json | jq '.[].id'` + "`" + `
2. Verify ID format (should be long alphanumeric strings)

### Error: "--name is required" or "--title is required"

**Cause:** Missing required flag

❌ WRONG:
` + "`" + `` + "`" + `bash
fizz boards create
fizz cards create --board=123
` + "`" + `` + "`" + `

✅ CORRECT:
` + "`" + `` + "`" + `bash
fizz boards create --name="Board Name"
fizz cards create --board=123 --title="Card Title"
` + "`" + `` + "`" + `

## Best Practices

💡 **Always use --format=json** for programmatic parsing

💡 **Store credentials in environment variables**, not in commands

💡 **Quote string values with spaces**:
- ✅ ` + "`" + `--title="My Card Title"` + "`" + `
- ❌ ` + "`" + `--title=My Card Title` + "`" + ` (will fail)

💡 **Use card numbers for simplicity**: Both work, but numbers are easier:
- ` + "`" + `fizz cards get 123` + "`" + ` (simple)
- ` + "`" + `fizz cards get 03ff063zx6myr3tdt98q3bjfj` + "`" + ` (UUID)

💡 **Check exit codes**:
- 0 = Success
- 1 = Error (check stderr for details)

💡 **Enable debug for troubleshooting**:
` + "`" + `` + "`" + `bash
fizz --debug boards list
` + "`" + `` + "`" + `

## Authentication and Prerequisites

**Required Credentials:**
- ` + "`" + `FIZZY_TOKEN` + "`" + ` - API authentication token
- ` + "`" + `FIZZY_ACCOUNT` + "`" + ` - Account ID

**Obtain Credentials:**
1. Sign up at https://fizzy.do
2. Navigate to Settings > API Tokens
3. Generate new token
4. Find your account ID in Settings > Account

**Setup for Shell Session:**
` + "`" + `` + "`" + `bash
export FIZZY_TOKEN="your_token_here"
export FIZZY_ACCOUNT="your_account_id"
` + "`" + `` + "`" + `

**Permanent Setup (add to ~/.zshrc or ~/.bashrc):**
` + "`" + `` + "`" + `bash
echo 'export FIZZY_TOKEN="your_token"' >> ~/.zshrc
echo 'export FIZZY_ACCOUNT="your_account_id"' >> ~/.zshrc
source ~/.zshrc
` + "`" + `` + "`" + `

## Rate Limits and Performance

**Rate Limits:**
- Managed by Fizzy.do API
- No client-side rate limiting enforced
- Tool will return API errors if limits exceeded

**Performance Tips:**
- Use ` + "`" + `--limit` + "`" + ` flag to reduce response size
- Batch operations when possible
- Use JSON format for faster parsing

**Timeouts:**
- Default timeout: 30 seconds per request
- Controlled by underlying HTTP client

## Alternative Access Methods

### Web Interface
- URL: https://fizzy.do
- Full feature parity with CLI
- Supports all CRUD operations through GUI

### API Documentation
- URL: https://fizzy.do/docs/api
- Direct REST API access
- OpenAPI/Swagger documentation available

### MCP Server
- Status: Not available
- Check https://github.com/visionik for future releases

## Quick Reference

| Task | Command |
|------|---------|
| List boards | ` + "`" + `fizz boards list --format=json` + "`" + ` |
| Create board | ` + "`" + `fizz boards create --name="Name" --format=json` + "`" + ` |
| List cards | ` + "`" + `fizz cards list --format=json` + "`" + ` |
| Create card | ` + "`" + `fizz cards create --board=ID --title="Title" --format=json` + "`" + ` |
| Close card | ` + "`" + `fizz cards close CARD_ID` + "`" + ` |
| Add comment | ` + "`" + `fizz comments create CARD_ID --body="Text" --format=json` + "`" + ` |
| List tags | ` + "`" + `fizz tags list --format=json` + "`" + ` |
| List users | ` + "`" + `fizz users list --format=json` + "`" + ` |
| Get identity | ` + "`" + `fizz identity get --format=json` + "`" + ` |

## Operational Semantics

### Idempotency

**Idempotent Operations:**
- GET operations (get, list) - Safe to repeat
- DELETE operations - Returns success even if already deleted
- UPDATE operations - Overwrites state

**Non-Idempotent Operations:**
- CREATE operations - Each call creates new resource

### Side Effects

| Command | Side Effects | Reversible? |
|---------|-------------|-------------|
| list, get | None (read-only) | N/A |
| create | Creates resource | Manual delete required |
| update | Modifies resource | Manual revert required |
| delete | Removes resource | Not reversible |
| close/reopen | Changes card state | Use reopen/close to reverse |

### Exit Codes

- ` + "`" + `0` + "`" + ` - Command succeeded
- ` + "`" + `1` + "`" + ` - Command failed (check error message)
`

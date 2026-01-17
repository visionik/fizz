# Product Requirements Document: fizz

**Generated**: 2026-01-17
**Status**: Ready for AI Interview

## Initial Input

**Project Description**: a sane, simple CLI for fizzy.do written in go and using libfizz-go

**I want to build fizz that has the following features:**
1. talk to fizzy.do via the libfizz-go library https://github.com/visionik/libfizz-go
2. support all of the feature of libfizz-go
3. make the cli easy to use with a syntax of fizz noun verb --flags "data"
4. support https://github.com/visionik/dashdash
---

# Specification Generation

Agent workflow for creating project specifications via structured interview.

Legend (from RFC2119): !=MUST, ~=SHOULD, ≉=SHOULD NOT, ⊗=MUST NOT, ?=MAY.

## Input Template

```
I want to build fizz that has the following features:
1. [feature]
2. [feature]
...
N. [feature]
```

## Interview Process

- ~ Use Claude AskInterviewQuestion when available (emulate it if not available)
- ! If Input Template fields are empty: ask overview, then features, then details
- ! Ask one focused, non-trivial question per step
- ~ Provide numbered answer options when appropriate
- ! Include "other" option for custom/unknown responses
- ! make it clear which option you feel is RECOMMENDED
- ! when you are done, append to the end of this file all questions asked and answers given.

**Question Areas:**

- ! Missing decisions (language, framework, deployment)
- ! Edge cases (errors, boundaries, failure modes)
- ! Implementation details (architecture, patterns, libraries)
- ! Requirements (performance, security, scalability)
- ! UX/constraints (users, timeline, compatibility)
- ! Tradeoffs (simplicity vs features, speed vs safety)

**Completion:**

- ! Continue until little ambiguity remains
- ! Ensure spec is comprehensive enough to implement

## Output Generation

- ! Generate as SPECIFICATION.md
- ! follow all relevant warping guidelines
- ! use RFC2119 MUST, SHOULD, MAY, SHOULD NOT, MUST NOT wording
- ! Break into phases, subphases, tasks
- ! end of each phase/subphase must implement and run testing until it passes
- ! Mark all dependencies explicitly: "Phase 2 (depends on: Phase 1)"
- ! Design for parallel work (multiple agents)
- ⊗ Write code (specification only)

## Afterwards

- ! let user know to type "implement SPECIFICATION.md" to start implementation

**Structure:**

```markdown
# Project Name

## Overview

## Requirements

## Architecture

## Implementation Plan

### Phase 1: Foundation

#### Subphase 1.1: Setup

- Task 1.1.1: (description, dependencies, acceptance criteria)

#### Subphase 1.2: Core (depends on: 1.1)

### Phase 2: Features (depends on: Phase 1)

## Testing Strategy

## Deployment
```

## Best Practices

- ! Detailed enough to implement without guesswork
- ! Clear scope boundaries (in vs out)
- ! Include rationale for major decisions
- ~ Size tasks for 1-4 hours
- ! Minimize inter-task dependencies
- ! Define clear component interfaces

## Anti-Patterns

- ⊗ Multiple questions at once
- ⊗ Assumptions without clarifying
- ⊗ Vague requirements
- ⊗ Missing dependencies
- ⊗ Sequential tasks that could be parallel

---

## Interview Results

**Date**: 2026-01-17
**Status**: Complete

### Q1: CLI Command Structure
**Question**: How should the CLI commands be structured?

Options:
1. Direct service mapping: `fizz boards list`, `fizz cards create` (RECOMMENDED - simplest)
2. Grouped commands: `fizz board --list`, `fizz card --create`
3. Single command with subcommands: `fizz --boards list`, `fizz --cards create`

**Answer**: Option 1

### Q2: Authentication
**Question**: How should authentication be handled?

Options:
1. Environment variables only: `FIZZY_TOKEN`, `FIZZY_ACCOUNT` (RECOMMENDED - 12-factor app)
2. Config file: `~/.fizz.yaml`
3. Both env vars and config file (config file overrides)

**Answer**: Option 1

### Q3: Output Format
**Question**: What output formats should be supported?

Options:
1. Human-readable tables only
2. JSON only (machine-readable)
3. Multiple formats with `--format` flag: table (default), json, yaml (RECOMMENDED)

**Answer**: Option 3

### Q4: Input Handling
**Question**: How should complex input data be provided?

Options:
1. Flags only: `fizz cards create --board=123 --title="Bug" --body="Description"`
2. JSON stdin only: `echo '{"title":"Bug"}' | fizz cards create --board=123`
3. Both flags and stdin/file support (RECOMMENDED - most flexible)

**Answer**: Option 3

### Q5: Error Handling
**Question**: How should errors be displayed?

Options:
1. Minimal: just error message
2. Standard: error message + helpful hints (RECOMMENDED)
3. Verbose: full stack traces

**Answer**: Option 2

### Q6: Pagination
**Question**: How should paginated list results be handled?

Options:
1. Auto-fetch all pages (RECOMMENDED - simpler UX)
2. Page-by-page with `--page` flag
3. Cursor-based with `--cursor` and `--limit`

**Answer**: Option 1 (with optional `--limit` flag)

### Q7: Discoverability
**Question**: How should users discover commands and options?

Options:
1. Built-in help only: `fizz help`, `fizz boards --help`
2. Help + shell completion scripts (RECOMMENDED)
3. Help + completion + interactive TUI

**Answer**: Option 2 (with `--ai-help` flag and default help display)

### Q8: Card ID vs Number
**Question**: Fizzy uses both card IDs (UUID) and numbers. How should the CLI handle this?

Options:
1. ID only (UUID)
2. Accept both, resolve internally (RECOMMENDED - best UX)
3. Number only

**Answer**: Option 2

### Q9: Testing Strategy
**Question**: What testing approach should be used?

Options:
1. Unit tests only
2. Unit + integration tests with real API (RECOMMENDED)
3. Unit + integration + E2E with mock server

**Answer**: Option 2 (with ephemeral test board)

### Q10: Installation
**Question**: How should users install the CLI?

Options:
1. `go install` only
2. Pre-built binaries via GitHub Releases
3. Multiple methods: go install + releases + Homebrew (RECOMMENDED)

**Answer**: Option 3

### Q11: Dashdash Integration
**Question**: How should dashdash (https://github.com/visionik/dashdash) be integrated?

**Answer**: Just include `--ai-help` flag in all help output

### Q12: Implementation Timeline
**Question**: Which services should be implemented?

Options:
1. Core services only (boards, cards)
2. All 11 services from libfizz-go (RECOMMENDED - feature complete)
3. MVP first, then iterate

**Answer**: Option 2

### Q13: Project Structure
**Question**: How should the Go project be structured?

Options:
1. Flat structure (all in root)
2. Standard Go CLI layout: cmd/, internal/ (RECOMMENDED)
3. Domain-driven: services/, models/, handlers/

**Answer**: Option 2

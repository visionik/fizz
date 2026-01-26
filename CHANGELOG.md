# Changelog

All notable changes to fizz will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-01-26

### Added
- Complete Fizzy.do API coverage (11 services)
  - Identity, Boards, Cards, Comments, Reactions
  - Steps, Tags, Columns, Users, Notifications, Uploads
- Multiple output formats: table, JSON, YAML
- Card actions: close, reopen, assign, tag, move, watch, golden
- Compact table display for cards and boards
- `--ai-help` flag for LLM/agent guidance
- Shell completion for Bash, Zsh, Fish
- Debug mode (`--debug`)
- SKILL.md for Clawdbot integration
- Built on [libfizz-go](https://github.com/visionik/libfizz-go)

[Unreleased]: https://github.com/visionik/fizz/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/visionik/fizz/releases/tag/v0.1.0

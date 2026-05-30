# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build

```bash
make build
```

Or:

```bash
go build -o open-notebook ./main.go
```

Version is set at build time via ldflags from git tags: `-X github.com/darimuri/open-notebook-cli/cmd.version=$(git describe --tags)`

### Install locally

```bash
make install-local
```

This builds and copies the binary to `~/.local/bin/open-notebook`.

## Test

```bash
go test ./... -v
```

## Architecture

### Command Structure (cmd/)

Commands use Cobra. Each domain has its own file:
- `root.go` - Root command, global flags, config loading
- `notebooks.go` - Notebook CRUD operations
- `notes.go` - Note operations
- `sources.go` - Source management including recursive crawling, embed-batch
- `search.go` - Search and ask operations
- `commands.go` - Command job listing and cancel
- `output.go` - JSON/table output formatting
- `config_backup.go` - Config backup/restore
- `update.go` - CLI self-update

### Internal Packages (internal/)

- `api/` - REST API client with Get/Post/Delete methods
- `auth/` - Authentication middleware
- `config/` - Config file loading (viper-based)
- `crawler/` - HTML link extraction for recursive source adding
- `formatter/` - Output formatting

### API Client

Located at `internal/api/client.go`. Uses `viper` for config with keys like `api_url`, `api_key`. API responses are typed in `internal/api/types.go`.

## Configuration

Config precedence: CLI flags → Environment variables → Config file (`~/.config/open-notebook/config.yaml`) → Defaults.

Environment variables: `OPEN_NOTEBOOK_API_URL`, `OPEN_NOTEBOOK_API_KEY`, `OPEN_NOTEBOOK_OUTPUT`.

## Skills

CLI skill definitions are in the separate repo. When modifying CLI commands, update the skill:
- Skill path: `../open-notebook-cli-skills/skills/open-notebook/SKILL.md`

## Server Source Reference

When adding or modifying commands, always verify the API endpoint in the server source first:
- Server source: `../../lfnovo/open-notebook/api/routers/`
- Check the corresponding router (e.g., `sources.py`, `commands.py`) for API signature and parameters

## Command Development Workflow

1. **Before implementing**: Check server source to understand API behavior
2. **Implement CLI command** in `cmd/`
3. **Update skill** at `../open-notebook-cli-skills/skills/open-notebook/SKILL.md`
4. **Do not bump skill version** unless explicitly required

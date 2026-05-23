---
name: open-notebook
description: CLI tool for Open Notebook API - Manage notebooks, sources, search, and chat with your research materials
version: "0.0.4"
trigger: /open-notebook
---

# Open Notebook

CLI tool for Open Notebook API - Research Assistant. This skill provides commands to manage notebooks, sources, search, and chat with your research materials.

## Setup

Make sure `open-notebook` CLI is installed and accessible in PATH. See [CLI Installation](#cli-installation) below.

## Commands

### Notebooks

```bash
# List all notebooks
/open-notebook notebooks list

# Get notebook details
/open-notebook notebooks get <notebook_id>

# Create a new notebook
/open-notebook notebooks create "Research Notes"

# Update notebook name
/open-notebook notebooks update <notebook_id> "New Name"

# Delete notebook (sources become unattached)
/open-notebook notebooks delete <notebook_id>

# Delete notebook and its exclusive sources
/open-notebook notebooks delete --delete-sources <notebook_id>

# Preview what will be deleted
/open-notebook notebooks delete-preview <notebook_id>

# Add existing source to notebook
/open-notebook notebooks add-source <notebook_id> <source_id>

# Remove source from notebook
/open-notebook notebooks remove-source <notebook_id> <source_id>
```

### Sources

```bash
# Add a single URL
/open-notebook sources add https://example.com/article

# Add multiple URLs
/open-notebook sources add https://site.com/page1 https://site.com/page2

# Recursively add all internal links from a page
/open-notebook sources add -r https://docs.site.com/guide

# Recursively with depth limit
/open-notebook sources add -r --depth 3 https://docs.site.com/guide

# Add from file (one URL per line)
/open-notebook sources add --file urls.txt

# Add text content
/open-notebook sources add --text "Important notes"

# Add to specific notebook
/open-notebook sources add -n <notebook_id> https://example.com

# Add URL but skip embedding (for bulk import)
/open-notebook sources add --skip-embed https://example.com

# Upload a file
/open-notebook sources upload /path/to/file.pdf

# Download a source
/open-notebook sources download <source_id>

# Retry failed source
/open-notebook sources retry <source_id>

# Get source insights
/open-notebook sources insights <source_id>

# Embed a source for vector search
/open-notebook sources embed <source_id>
/open-notebook sources embed <source_id> --wait
/open-notebook sources embed <source_id> --wait --polling-period 5

# Embed all non-embedded sources (batch)
/open-notebook sources embed-batch
/open-notebook sources embed-batch --max 10
/open-notebook sources embed-batch --notebook <notebook_id>
/open-notebook sources embed-batch --notebook <notebook_id> --max 5
/open-notebook sources embed-batch --polling-period 5

# Check source status
/open-notebook sources status <source_id>
```

### Sources (list command)

```bash
# List all sources
/open-notebook sources list

# Filter by notebook
/open-notebook sources list --notebook <notebook_id>

# Limit number of sources
/open-notebook sources list --max 20

# Combine filters
/open-notebook sources list --notebook <notebook_id> --max 50
```

### Search

```bash
# Search notebooks
/open-notebook search search "machine learning"

# Ask a question (uses default model from API)
open-notebook search ask "What is the main contribution of this paper?"

# Ask with specific notebook
open-notebook search --notebook <notebook_id> ask "What is covered?"

# Ask with specific models
open-notebook search --strategy-model model:xxx --answer-model model:xxx --final-answer-model model:xxx ask "question"

# Simple ask (non-streaming response)
open-notebook search simple "Summarize the key points"
```

### Notes

```bash
# List all notes
/open-notebook notes list

# Get note details
/open-notebook notes get <note_id>

# Create a note in a notebook
/open-notebook notes create <notebook_id> "Note content"

# Update note
/open-notebook notes update <note_id> "Updated content"

# Delete note
/open-notebook notes delete <note_id>
```

### Config Backup & Restore

```bash
# Backup models, credentials, and defaults to a directory
/open-notebook config backup <output_dir>

# Restore from backup directory
/open-notebook config restore <input_dir>
```

### Commands

```bash
# List recent commands (server v1.8.5+ only)
/open-notebook commands list

# Note: List may return empty if not implemented on server
# Check command status
/open-notebook commands status <command_id>
```

## CLI Installation

### Quick install (Linux/macOS)
```bash
curl -sL https://raw.githubusercontent.com/darimuri/open-notebook-cli/main/install.sh | bash
```

### Quick install (Windows)
```powershell
irm https://raw.githubusercontent.com/darimuri/open-notebook-cli/main/install.ps1 | iex
```

### From source
```bash
git clone https://github.com/darimuri/open-notebook-cli.git
cd open-notebook-cli
go build -o open-notebook ./main.go
# Add to PATH or use ./open-notebook
```

### Manual download
Download from https://github.com/darimuri/open-notebook-cli/releases/latest

### Global Flags
```
--api-url string     API server URL
--api-key string    API key for authentication
--config string     Config file path
--notebook string   Default notebook ID
--debug             Enable debug output
-o, --output string Output format (table, json)
--version           Show version
```

### Update

```bash
# Update to the latest version
/open-notebook update
```

### Configuration

Environment variables:
- `OPEN_NOTEBOOK_API_URL` - API server URL (default: http://localhost:8080)
- `OPEN_NOTEBOOK_API_KEY` - API key for authentication
- `OPEN_NOTEBOOK_OUTPUT` - Default output format (table or json)

Config file: `~/.config/open-notebook/config.yaml`

```yaml
api-url: "https://open-notebook.darimuri.me"
notebook: "notebook:trzgg4xyju0awzcfltrp"  # default notebook ID
```

## Examples

### Research Workflow

```bash
# 1. Create a research notebook
/open-notebook notebooks create "ML Papers Review"

# 2. Add sources recursively from documentation
/open-notebook sources add -r --depth 2 -n <notebook_id> https://docs.site.com

# 3. Ask questions about the sources
/open-notebook search ask "What are the main topics covered?"

# 4. Add notes
/open-notebook notes create <notebook_id> "Key insight: The model uses..."

# 5. Check notebook status
/open-notebook notebooks get <notebook_id>
```

### Source Management

```bash
# Add multiple URLs from file
/open-notebook sources add --file paper-urls.txt

# Crawl entire documentation site
/open-notebook sources add -r --depth 5 https://docs.site.com

# Link sources to notebooks
/open-notebook notebooks add-source <notebook_id> <source_id>
```
# Open Notebook CLI

CLI tool for Open Notebook API - Research Assistant.

## Installation

### CLI Only

```bash
go install github.com/darimuri/open-notebook-cli@latest
```

Or build from source:

```bash
git clone https://github.com/darimuri/open-notebook-cli.git
cd open-notebook-cli
go build -o open-notebook-cli ./main.go
```

### Claude Code Plugin Installation

Install the Open Notebook plugin for use in Claude Code:

```bash
claude plugin marketplace add darimuri/open-notebook-cli
claude plugin install open-notebook@open-notebook-cli
```

After installation, you can use `/open-notebook` commands directly in Claude Code:

```bash
/open-notebook notebooks list
/open-notebook sources add -r https://docs.site.com
/open-notebook search ask "What is machine learning?"
```

## Configuration

Configuration is loaded in the following order (highest priority first):

1. **Command-line flags**
2. **Environment variables**
3. **Config file** (`~/.config/open-notebook/config.yaml`)
4. **Defaults**

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `OPEN_NOTEBOOK_API_URL` | API server URL | `http://localhost:8080` |
| `OPEN_NOTEBOOK_API_KEY` | API key for authentication | (empty) |
| `OPEN_NOTEBOOK_OUTPUT` | Default output format | `table` |

### Config File

Create `~/.config/open-notebook/config.yaml`:

```yaml
api_url: "http://localhost:8080"
api_key: "your-api-key"
output: "table"  # or "json"
```

## Usage

### Global Flags

```
--api-url string     API server URL
--api-key string      API key for authentication
--config string      Config file path
-o, --output string   Output format (table, json) (default "table")
```

### Notebooks

```bash
# List all notebooks
open-notebook-cli notebooks list

# Get a specific notebook
open-notebook-cli notebooks get <notebook_id>

# Create a new notebook
open-notebook-cli notebooks create "My Notebook"

# Update a notebook
open-notebook-cli notebooks update <notebook_id> "New Name"

# Delete a notebook (sources become unattached)
open-notebook-cli notebooks delete <notebook_id>

# Delete a notebook and its exclusive sources
open-notebook-cli notebooks delete --delete-sources <notebook_id>

# Preview what will be deleted
open-notebook-cli notebooks delete-preview <notebook_id>

# Add source to notebook
open-notebook-cli notebooks add-source <notebook_id> <source_id>

# Remove source from notebook
open-notebook-cli notebooks remove-source <notebook_id> <source_id>
```

### Notes

```bash
# List all notes
open-notebook-cli notes list

# Get a note
open-notebook-cli notes get <note_id>

# Create a note
open-notebook-cli notes create <notebook_id> "Note content here"

# Update a note
open-notebook-cli notes update <note_id> "Updated content"

# Delete a note
open-notebook-cli notes delete <note_id>
```

### Sources

```bash
# List all sources
open-notebook-cli sources list

# Add a single URL
open-notebook-cli sources add https://example.com/article

# Add multiple URLs
open-notebook-cli sources add https://site.com/page1 https://site.com/page2

# Add with recursive crawling (all internal links)
open-notebook-cli sources add -r https://docs.site.com/guide

# Add with depth limit
open-notebook-cli sources add -r --depth 3 https://docs.site.com/guide

# Add from file (one URL per line)
open-notebook-cli sources add --file urls.txt

# Add text content
open-notebook-cli sources add --text "Important notes"

# Add to specific notebook
open-notebook-cli sources add -n <notebook_id> https://example.com

# Upload a file
open-notebook-cli sources upload /path/to/file.pdf

# Download a source
open-notebook-cli sources download <source_id>

# Retry a failed source
open-notebook-cli sources retry <source_id>

# Get source insights
open-notebook-cli sources insights <source_id>

# Check source status
open-notebook-cli sources status <source_id>
```

### Search

```bash
# Search notebooks
open-notebook-cli search search "your query"

# Ask a question (detailed response)
open-notebook-cli search ask "What is machine learning?"

# Simple ask (quick answer)
open-notebook-cli search simple "What is AI?"
```

## Output Formats

### Table (default)

```
ID    NAME           DESCRIPTION
----  -------------  -------------
1     My Notebook    A test notebook
2     Research       AI research notes
```

### JSON

```bash
open-notebook-cli notebooks list --output json
```

```json
[
  {
    "id": "1",
    "name": "My Notebook",
    "description": "A test notebook",
    "archived": false,
    "created": "2024-01-01T00:00:00Z",
    "updated": "2024-01-02T00:00:00Z",
    "source_count": 5,
    "note_count": 10
  }
]
```

## Examples

### Research Workflow

```bash
# 1. Create a research notebook
open-notebook-cli notebooks create "ML Papers Review"

# 2. Add sources recursively from documentation
open-notebook-cli sources add -r --depth 2 https://docs.site.com

# 3. Ask questions about the sources
open-notebook-cli search ask "What are the main topics covered?"

# 4. Add notes
open-notebook-cli notes create <notebook_id> "Key insight: The model uses..."

# 5. Check notebook status
open-notebook-cli notebooks get <notebook_id>
```

### Source Management

```bash
# Add multiple URLs from file
open-notebook-cli sources add --file paper-urls.txt

# Crawl entire documentation site
open-notebook-cli sources add -r --depth 5 https://docs.site.com

# Link sources to notebooks
open-notebook-cli notebooks add-source <notebook_id> <source_id>
```

## Development

### Build

```bash
go build -o open-notebook-cli ./main.go
```

### Test

```bash
go test ./... -v
```

### Run

```bash
./open-notebook-cli --help
```

## Project Structure

```
open-notebook-cli/
├── cmd/                 # Cobra commands
│   ├── root.go         # Root command
│   ├── notebooks.go    # Notebook commands
│   ├── notes.go        # Note commands
│   ├── sources.go      # Source commands
│   ├── search.go       # Search commands
│   └── skills.go       # Skills commands
├── internal/
│   ├── api/           # API client
│   ├── auth/          # Auth middleware
│   ├── config/        # Config loading
│   ├── crawler/       # HTML link extraction
│   └── formatter/      # Output formatter
├── skills/
│   └── open-notebook/ # Claude Code skill
│       └── SKILL.md    # Skill definition
├── tests/
│   ├── unit/          # Unit tests
│   └── integration/   # Integration tests
├── main.go
└── go.mod
```

## License

MIT
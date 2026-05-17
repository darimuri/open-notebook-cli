# Open Notebook CLI

CLI tool for Open Notebook API - Research Assistant.

## Installation

```bash
go install github.com/darimuri/open-notebook-cli@latest
```

Or build from source:

```bash
git clone https://github.com/darimuri/open-notebook-cli.git
cd open-notebook-cli
go build -o open-notebook-cli ./main.go
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

# Delete a notebook
open-notebook-cli notebooks delete <notebook_id>

# Preview what will be deleted
open-notebook-cli notebooks delete-preview <notebook_id>
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

# Upload a source
open-notebook-cli sources upload /path/to/file

# Download a source
open-notebook-cli sources download <source_id>

# Retry a failed source
open-notebook-cli sources retry <source_id>
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

### Skills

```bash
# List available skills
open-notebook-cli skills list

# Invoke a skill
open-notebook-cli skills invoke brainstorming
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

```bash
# Use a different API server
open-notebook-cli --api-url http://192.168.1.100:8080 notebooks list

# Use API key authentication
open-notebook-cli --api-key my-secret-key notebooks list

# Use environment variable
export OPEN_NOTEBOOK_API_KEY=my-secret-key
open-notebook-cli notebooks list

# Get JSON output
open-notebook-cli notebooks list --output json
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
│   └── formatter/      # Output formatter
├── tests/
│   ├── unit/          # Unit tests
│   └── integration/   # Integration tests
├── main.go
└── go.mod
```

## License

MIT
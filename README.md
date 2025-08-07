# AI Git Auto

A Go library and CLI tool that automatically stages changes, generates intelligent commit messages using local Ollama models, commits, and pushes - all in one command.

## ✨ Features

- � **Complete Workflow**: `git add . → AI commit message → git commit → git push`
- �🔍 **Automatic Change Detection**: Scans staged files and analyzes diffs
- 🤖 **AI-Powered Messages**: Uses local Ollama models for meaningful commit messages
- 📝 **Conventional Commits**: Follows conventional commit format
- ⚙️ **Highly Configurable**: Customizable AI model, temperature, and behavior
- 🛡️ **Safe & Interactive**: Confirmation prompts with dry-run mode
- 🌍 **Global CLI**: Install once, use anywhere
- � **Library + CLI**: Use as Go library or standalone CLI tool

## 🚀 Quick Start (One Command)

```bash
# Install globally
curl -fsSL https://raw.githubusercontent.com/TheRealMasterK/Ai-Git-Comments-Auto/main/install.sh | bash

# Or build from source
make global-install

# Navigate to your Git repository and run
ai-git-auto
```

This single command will:
1. 📝 Stage all changes (`git add .`)
2. 🔍 Analyze the changes
3. 🤖 Generate an AI commit message
4. 💾 Commit with the message
5. 📤 Push to remote (with confirmation)

## Prerequisites
   ```bash
   # Install Ollama (macOS)
   brew install ollama

   # Start Ollama service
   ollama serve

   # Pull a model (e.g., llama2, codellama, mistral)
   ollama pull llama2
   ```

2. **Go**: Version 1.21 or higher

## Installation

### As a Go Module

```bash
go get github.com/TheRealMasterK/Ai-Git-Comments-Auto
```

### Build the CLI Tool

```bash
git clone https://github.com/TheRealMasterK/Ai-Git-Comments-Auto.git
cd Ai-Git-Comments-Auto
go build -o git-ai-commit ./cmd/git-ai-commit
```

## Quick Start

### Using the Library

```go
package main

import (
    "fmt"
    "log"

    "github.com/TheRealMasterK/Ai-Git-Comments-Auto"
)

func main() {
    // Create configuration
    config := gitcommenter.DefaultConfig()
    config.Model = "codellama"  // or "llama2", "mistral", etc.

    // Create commenter
    commenter := gitcommenter.New(config)

    // Scan staged changes
    changes, err := commenter.ScanStagedChanges()
    if err != nil {
        log.Fatal(err)
    }

    // Generate commit message
    suggestion, err := commenter.GenerateCommitMessage(changes)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Suggested commit: %s\\n", suggestion.Subject)
}
```

### Using the CLI Tool

```bash
# Stage your changes
git add .

# Generate commit message
./git-ai-commit

# Interactive mode (will prompt to commit)
./git-ai-commit -interactive

# Use specific model
./git-ai-commit -model codellama

# List available models
./git-ai-commit -list-models
```

## Configuration

The `Config` struct allows you to customize the behavior:

```go
type Config struct {
    OllamaEndpoint string        // Default: "http://localhost:11434"
    Model         string         // Default: "llama2"
    MaxTokens     int           // Default: 150
    Temperature   float64       // Default: 0.7 (0.0-1.0)
    RepositoryPath string       // Default: "."
    Timeout       time.Duration // Default: 30s
}
```

## CLI Options

```bash
Usage of git-ai-commit:
  -endpoint string
        Ollama endpoint (default "http://localhost:11434")
  -interactive
        Interactive mode to approve commit message
  -list-models
        List available Ollama models
  -max-tokens int
        Maximum tokens for response (default 150)
  -model string
        Ollama model to use (default "llama2")
  -repo string
        Path to git repository (default ".")
  -temperature float
        Temperature for AI model (0.0-1.0) (default 0.7)
```

## API Reference

### Main Types

#### `GitCommenter`
The main struct that handles scanning and message generation.

```go
func New(config *Config) *GitCommenter
func (gc *GitCommenter) ScanStagedChanges() ([]FileChange, error)
func (gc *GitCommenter) GenerateCommitMessage(changes []FileChange) (*CommitSuggestion, error)
func (gc *GitCommenter) ListAvailableModels() ([]string, error)
```

#### `FileChange`
Represents a changed file with its metadata.

```go
type FileChange struct {
    FilePath     string // Path to the changed file
    ChangeType   string // "added", "modified", "deleted", "renamed"
    Diff         string // Git diff output
    LinesAdded   int    // Number of lines added
    LinesRemoved int    // Number of lines removed
}
```

#### `CommitSuggestion`
Contains the AI-generated commit message suggestion.

```go
type CommitSuggestion struct {
    Subject       string   // Commit subject line
    Body         string   // Commit body (optional)
    Confidence   float64  // Confidence score (0.0-1.0)
    FilesAffected []string // List of affected file paths
}
```

## Examples

### Basic Usage

```go
// examples/basic/main.go
commenter := gitcommenter.New(gitcommenter.DefaultConfig())
changes, _ := commenter.ScanStagedChanges()
suggestion, _ := commenter.GenerateCommitMessage(changes)
fmt.Println(suggestion.Subject)
```

### Custom Configuration

```go
config := &gitcommenter.Config{
    OllamaEndpoint: "http://localhost:11434",
    Model:         "codellama",
    Temperature:   0.3,  // More focused responses
    MaxTokens:     100,  // Shorter messages
}

commenter := gitcommenter.New(config)
```

### Integration with Git Hooks

Create a pre-commit hook that suggests commit messages:

```bash
#!/bin/bash
# .git/hooks/prepare-commit-msg

if [ -z "$2" ]; then
    ./git-ai-commit > /tmp/ai-commit-msg
    echo "# AI-suggested commit message:" > "$1"
    cat /tmp/ai-commit-msg >> "$1"
    echo "" >> "$1"
    echo "# Please review and edit as needed" >> "$1"
fi
```

## Recommended Models

- **codellama**: Best for code-related commits
- **llama2**: General purpose, good balance
- **mistral**: Fast and efficient
- **codellama:13b**: More accurate but slower

## Workflow

1. Make your code changes
2. Stage changes with `git add .`
3. Run the AI commit tool
4. Review and approve the suggested message
5. Commit with the generated message

## Error Handling

The library handles common scenarios:
- No staged changes
- Invalid Git repository
- Ollama service unavailable
- Model not found
- Network timeouts

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Troubleshooting

### Ollama Not Running
```
Error: failed to call Ollama API: connection refused
```
Solution: Start Ollama with `ollama serve`

### Model Not Found
```
Error: model 'modelname' not found
```
Solution: Pull the model with `ollama pull modelname`

### No Staged Changes
```
No staged changes found
```
Solution: Stage your changes with `git add .`

### Permission Denied
```
Error: not in a git repository
```
Solution: Run from within a Git repository directory

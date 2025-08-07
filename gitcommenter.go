// Package gitcommenter provides functionality to scan Git changes and generate
// commit messages using a local Ollama model.
package gitcommenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// Config holds the configuration for the Git commenter
type Config struct {
	// OllamaEndpoint is the Ollama API endpoint (default: http://localhost:11434)
	OllamaEndpoint string
	// Model is the Ollama model to use (default: llama2)
	Model string
	// MaxTokens is the maximum number of tokens for the response
	MaxTokens int
	// Temperature controls randomness in the response (0.0 to 1.0)
	Temperature float64
	// RepositoryPath is the path to the Git repository
	RepositoryPath string
	// Timeout is the HTTP request timeout
	Timeout time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		OllamaEndpoint: "http://localhost:11434",
		Model:         "llama2",
		MaxTokens:     150,
		Temperature:   0.7,
		RepositoryPath: ".",
		Timeout:       30 * time.Second,
	}
}

// GitCommenter handles scanning Git changes and generating commit messages
type GitCommenter struct {
	config *Config
	client *http.Client
}

// New creates a new GitCommenter with the given configuration
func New(config *Config) *GitCommenter {
	if config == nil {
		config = DefaultConfig()
	}

	return &GitCommenter{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// FileChange represents a changed file with its diff
type FileChange struct {
	FilePath   string
	ChangeType string // "added", "modified", "deleted", "renamed"
	Diff       string
	LinesAdded int
	LinesRemoved int
}

// CommitSuggestion represents a suggested commit message
type CommitSuggestion struct {
	Subject     string
	Body        string
	Confidence  float64
	FilesAffected []string
}

// ScanStagedChanges scans the staged changes in the Git repository
func (gc *GitCommenter) ScanStagedChanges() ([]FileChange, error) {
	// Check if we're in a git repository
	if err := gc.ensureGitRepository(); err != nil {
		return nil, fmt.Errorf("not in a git repository: %w", err)
	}

	// Get list of staged files
	cmd := exec.Command("git", "diff", "--cached", "--name-status")
	cmd.Dir = gc.config.RepositoryPath
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []FileChange{}, nil // No staged changes
	}

	var changes []FileChange
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		status := parts[0]
		filepath := parts[1]

		change := FileChange{
			FilePath:   filepath,
			ChangeType: gc.parseChangeType(status),
		}

		// Get the diff for this file
		diff, linesAdded, linesRemoved, err := gc.getFileDiff(filepath)
		if err != nil {
			// Log error but continue with other files
			fmt.Printf("Warning: failed to get diff for %s: %v\n", filepath, err)
			continue
		}

		change.Diff = diff
		change.LinesAdded = linesAdded
		change.LinesRemoved = linesRemoved

		changes = append(changes, change)
	}

	return changes, nil
}

// GenerateCommitMessage generates a commit message based on the changes
func (gc *GitCommenter) GenerateCommitMessage(changes []FileChange) (*CommitSuggestion, error) {
	if len(changes) == 0 {
		return nil, fmt.Errorf("no changes to analyze")
	}

	// Build context for the AI model
	context := gc.buildChangeContext(changes)

	// Create prompt for the AI model
	prompt := gc.buildPrompt(context, changes)

	// Call Ollama API
	response, err := gc.callOllama(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate commit message: %w", err)
	}

	// Parse and return the suggestion
	suggestion := gc.parseCommitSuggestion(response, changes)
	return suggestion, nil
}

// ensureGitRepository checks if the current directory is a Git repository
func (gc *GitCommenter) ensureGitRepository() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = gc.config.RepositoryPath
	_, err := cmd.Output()
	return err
}

// parseChangeType converts Git status to readable change type
func (gc *GitCommenter) parseChangeType(status string) string {
	switch status[0] {
	case 'A':
		return "added"
	case 'M':
		return "modified"
	case 'D':
		return "deleted"
	case 'R':
		return "renamed"
	case 'C':
		return "copied"
	default:
		return "modified"
	}
}

// getFileDiff gets the diff for a specific file
func (gc *GitCommenter) getFileDiff(filepath string) (string, int, int, error) {
	cmd := exec.Command("git", "diff", "--cached", "--", filepath)
	cmd.Dir = gc.config.RepositoryPath
	output, err := cmd.Output()
	if err != nil {
		return "", 0, 0, err
	}

	diff := string(output)
	linesAdded, linesRemoved := gc.countDiffLines(diff)

	return diff, linesAdded, linesRemoved, nil
}

// countDiffLines counts added and removed lines in a diff
func (gc *GitCommenter) countDiffLines(diff string) (added, removed int) {
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			added++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			removed++
		}
	}
	return added, removed
}

// buildChangeContext creates a summary of changes for the AI model
func (gc *GitCommenter) buildChangeContext(changes []FileChange) string {
	var context strings.Builder

	context.WriteString("Files changed:\n")
	for _, change := range changes {
		context.WriteString(fmt.Sprintf("- %s (%s): +%d -%d lines\n",
			change.FilePath, change.ChangeType, change.LinesAdded, change.LinesRemoved))
	}

	return context.String()
}

// buildPrompt creates the prompt for the AI model
func (gc *GitCommenter) buildPrompt(context string, changes []FileChange) string {
	var prompt strings.Builder

	prompt.WriteString("You are a helpful assistant that generates concise, meaningful Git commit messages based on code changes.\n\n")
	prompt.WriteString("Here are the changes made to the repository:\n\n")
	prompt.WriteString(context)
	prompt.WriteString("\n")

	// Add some diff context for key changes
	for i, change := range changes {
		if i >= 3 { // Limit to first 3 files to avoid token limits
			break
		}
		if change.Diff != "" {
			prompt.WriteString(fmt.Sprintf("Diff for %s:\n", change.FilePath))
			// Truncate diff if too long
			diff := change.Diff
			if len(diff) > 1000 {
				diff = diff[:1000] + "\n... (truncated)"
			}
			prompt.WriteString(diff)
			prompt.WriteString("\n---\n\n")
		}
	}

	prompt.WriteString("Please generate a commit message following conventional commit format:\n")
	prompt.WriteString("- Use a clear, concise subject line (50 characters or less)\n")
	prompt.WriteString("- Start with a type (feat, fix, docs, style, refactor, test, chore)\n")
	prompt.WriteString("- Use present tense, imperative mood\n")
	prompt.WriteString("- Include a body if needed to explain what and why\n\n")
	prompt.WriteString("Respond with only the commit message, no additional text or formatting.")

	return prompt.String()
}

// OllamaRequest represents a request to the Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Options struct {
		Temperature float64 `json:"temperature"`
		NumPredict  int     `json:"num_predict"`
	} `json:"options"`
}

// OllamaResponse represents a response from the Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// callOllama makes a request to the Ollama API
func (gc *GitCommenter) callOllama(prompt string) (string, error) {
	req := OllamaRequest{
		Model:  gc.config.Model,
		Prompt: prompt,
		Stream: false,
	}
	req.Options.Temperature = gc.config.Temperature
	req.Options.NumPredict = gc.config.MaxTokens

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := gc.client.Post(gc.config.OllamaEndpoint+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return strings.TrimSpace(ollamaResp.Response), nil
}

// parseCommitSuggestion parses the AI response into a CommitSuggestion
func (gc *GitCommenter) parseCommitSuggestion(response string, changes []FileChange) *CommitSuggestion {
	lines := strings.Split(response, "\n")

	var subject, body string
	var filesAffected []string

	if len(lines) > 0 {
		subject = strings.TrimSpace(lines[0])
	}

	if len(lines) > 1 {
		bodyLines := lines[1:]
		// Remove empty lines at the beginning
		for i, line := range bodyLines {
			if strings.TrimSpace(line) != "" {
				bodyLines = bodyLines[i:]
				break
			}
		}
		body = strings.Join(bodyLines, "\n")
	}

	for _, change := range changes {
		filesAffected = append(filesAffected, change.FilePath)
	}

	return &CommitSuggestion{
		Subject:       subject,
		Body:         strings.TrimSpace(body),
		Confidence:   0.8, // Default confidence
		FilesAffected: filesAffected,
	}
}

// GetRepository returns the current repository path
func (gc *GitCommenter) GetRepository() string {
	return gc.config.RepositoryPath
}

// SetModel changes the Ollama model
func (gc *GitCommenter) SetModel(model string) {
	gc.config.Model = model
}

// ListAvailableModels lists available Ollama models
func (gc *GitCommenter) ListAvailableModels() ([]string, error) {
	resp, err := gc.client.Get(gc.config.OllamaEndpoint + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var models []string
	for _, model := range response.Models {
		models = append(models, model.Name)
	}

	return models, nil
}

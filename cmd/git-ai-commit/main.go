package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gitcommenter "github.com/TheRealMasterK/Ai-Git-Comments-Auto"
)

func main() {
	var (
		repoPath    = flag.String("repo", ".", "Path to git repository")
		model       = flag.String("model", "llama2", "Ollama model to use")
		endpoint    = flag.String("endpoint", "http://localhost:11434", "Ollama endpoint")
		temperature = flag.Float64("temperature", 0.7, "Temperature for AI model (0.0-1.0)")
		maxTokens   = flag.Int("max-tokens", 150, "Maximum tokens for response")
		listModels  = flag.Bool("list-models", false, "List available Ollama models")
		interactive = flag.Bool("interactive", false, "Interactive mode to approve commit message")
	)
	flag.Parse()

	// Create configuration
	config := &gitcommenter.Config{
		OllamaEndpoint: *endpoint,
		Model:         *model,
		MaxTokens:     *maxTokens,
		Temperature:   *temperature,
		RepositoryPath: *repoPath,
	}

	// Create commenter
	commenter := gitcommenter.New(config)

	// List models if requested
	if *listModels {
		models, err := commenter.ListAvailableModels()
		if err != nil {
			log.Fatalf("Failed to list models: %v", err)
		}

		fmt.Println("Available Ollama models:")
		for _, model := range models {
			fmt.Printf("  - %s\n", model)
		}
		return
	}

	// Get absolute path for better error messages
	absPath, err := filepath.Abs(*repoPath)
	if err != nil {
		log.Fatalf("Invalid repository path: %v", err)
	}

	fmt.Printf("Scanning staged changes in: %s\n", absPath)

	// Scan staged changes
	changes, err := commenter.ScanStagedChanges()
	if err != nil {
		log.Fatalf("Failed to scan changes: %v", err)
	}

	if len(changes) == 0 {
		fmt.Println("No staged changes found. Run 'git add .' first to stage your changes.")
		return
	}

	// Display found changes
	fmt.Printf("\nFound %d staged file(s):\n", len(changes))
	for _, change := range changes {
		fmt.Printf("  %s: %s (+%d -%d lines)\n",
			change.ChangeType, change.FilePath, change.LinesAdded, change.LinesRemoved)
	}

	fmt.Println("\nGenerating commit message...")

	// Generate commit message
	suggestion, err := commenter.GenerateCommitMessage(changes)
	if err != nil {
		log.Fatalf("Failed to generate commit message: %v", err)
	}

	// Display the suggestion
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("SUGGESTED COMMIT MESSAGE")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Subject: %s\n", suggestion.Subject)

	if suggestion.Body != "" {
		fmt.Printf("\nBody:\n%s\n", suggestion.Body)
	}

	fmt.Printf("\nFiles affected: %s\n", strings.Join(suggestion.FilesAffected, ", "))
	fmt.Printf("Confidence: %.1f%%\n", suggestion.Confidence*100)
	fmt.Println(strings.Repeat("=", 60))

	// Interactive mode
	if *interactive {
		fmt.Print("\nDo you want to use this commit message? (y/n): ")
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
			if err := commitChanges(suggestion.Subject, suggestion.Body, *repoPath); err != nil {
				log.Fatalf("Failed to commit changes: %v", err)
			}
			fmt.Println("✅ Changes committed successfully!")
		} else {
			fmt.Println("Commit cancelled. You can manually commit with:")
			fmt.Printf("git commit -m \"%s\"\n", suggestion.Subject)
		}
	} else {
		fmt.Println("\nTo commit with this message, run:")
		if suggestion.Body != "" {
			fmt.Printf("git commit -m \"%s\" -m \"%s\"\n", suggestion.Subject, suggestion.Body)
		} else {
			fmt.Printf("git commit -m \"%s\"\n", suggestion.Subject)
		}
	}
}

// commitChanges commits the staged changes with the generated message
func commitChanges(subject, body, repoPath string) error {
	args := []string{"commit", "-m", subject}
	if body != "" {
		args = append(args, "-m", body)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

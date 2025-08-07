package main

import (
	"fmt"
	"log"

	gitcommenter "github.com/TheRealMasterK/Ai-Git-Comments-Auto"
)

func main() {
	// Create a default configuration
	config := gitcommenter.DefaultConfig()

	// Customize configuration if needed
	config.Model = "codellama"              // Use CodeLlama model
	config.Temperature = 0.5                // Lower temperature for more focused responses
	config.RepositoryPath = "."            // Current directory

	// Create the commenter
	commenter := gitcommenter.New(config)

	// List available models
	fmt.Println("Available Ollama models:")
	models, err := commenter.ListAvailableModels()
	if err != nil {
		log.Printf("Warning: Could not fetch models: %v", err)
	} else {
		for _, model := range models {
			fmt.Printf("  - %s\n", model)
		}
	}

	// Scan for staged changes
	fmt.Println("\nScanning staged changes...")
	changes, err := commenter.ScanStagedChanges()
	if err != nil {
		log.Fatalf("Error scanning changes: %v", err)
	}

	if len(changes) == 0 {
		fmt.Println("No staged changes found. Please run 'git add .' first.")
		return
	}

	// Display changes
	fmt.Printf("Found %d staged files:\n", len(changes))
	for _, change := range changes {
		fmt.Printf("  %s: %s (+%d -%d lines)\n",
			change.ChangeType, change.FilePath, change.LinesAdded, change.LinesRemoved)
	}

	// Generate commit message
	fmt.Println("\nGenerating commit message...")
	suggestion, err := commenter.GenerateCommitMessage(changes)
	if err != nil {
		log.Fatalf("Error generating commit message: %v", err)
	}

	// Display suggestion
	fmt.Println("\n--- Suggested Commit Message ---")
	fmt.Printf("Subject: %s\n", suggestion.Subject)
	if suggestion.Body != "" {
		fmt.Printf("Body: %s\n", suggestion.Body)
	}
	fmt.Printf("Confidence: %.1f%%\n", suggestion.Confidence*100)
	fmt.Printf("Files affected: %v\n", suggestion.FilesAffected)
}

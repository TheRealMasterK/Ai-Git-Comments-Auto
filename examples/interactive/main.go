package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	gitcommenter "github.com/TheRealMasterK/Ai-Git-Comments-Auto"
)

func main() {
	fmt.Println("🤖 AI Git Commit Message Generator")
	fmt.Println("==================================")

	// Check if Ollama is running
	if !isOllamaRunning() {
		fmt.Println("❌ Ollama is not running. Please start it with: ollama serve")
		return
	}

	// Check if we're in a git repository
	if !isGitRepository() {
		fmt.Println("❌ Not in a Git repository. Please run this from a Git repository.")
		return
	}

	// Create commenter with interactive model selection
	commenter := createCommenterWithModelSelection()
	if commenter == nil {
		return
	}

	// Check for staged changes
	changes, err := commenter.ScanStagedChanges()
	if err != nil {
		log.Fatalf("❌ Error scanning changes: %v", err)
	}

	if len(changes) == 0 {
		fmt.Println("📝 No staged changes found.")
		fmt.Println("💡 Tip: Stage your changes first with 'git add .' or 'git add <files>'")
		return
	}

	// Display changes summary
	displayChangesSummary(changes)

	// Generate commit message
	fmt.Println("🎯 Generating AI commit message...")
	suggestion, err := commenter.GenerateCommitMessage(changes)
	if err != nil {
		log.Fatalf("❌ Error generating commit message: %v", err)
	}

	// Display suggestion
	displayCommitSuggestion(suggestion)

	// Interactive approval
	if askForApproval() {
		if commitWithMessage(suggestion) {
			fmt.Println("✅ Successfully committed changes!")
		} else {
			fmt.Println("❌ Failed to commit changes.")
		}
	} else {
		fmt.Println("📋 You can manually commit with:")
		if suggestion.Body != "" {
			fmt.Printf("   git commit -m \"%s\" -m \"%s\"\n", suggestion.Subject, suggestion.Body)
		} else {
			fmt.Printf("   git commit -m \"%s\"\n", suggestion.Subject)
		}
	}
}

func isOllamaRunning() bool {
	commenter := gitcommenter.New(gitcommenter.DefaultConfig())
	_, err := commenter.ListAvailableModels()
	return err == nil
}

func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	_, err := cmd.Output()
	return err == nil
}

func createCommenterWithModelSelection() *gitcommenter.GitCommenter {
	config := gitcommenter.DefaultConfig()

	// Create temporary commenter to list models
	tempCommenter := gitcommenter.New(config)
	models, err := tempCommenter.ListAvailableModels()
	if err != nil {
		fmt.Printf("❌ Error listing models: %v\n", err)
		return nil
	}

	if len(models) == 0 {
		fmt.Println("❌ No Ollama models found. Please pull a model first:")
		fmt.Println("   ollama pull llama2")
		return nil
	}

	// Display available models
	fmt.Println("\n📚 Available models:")
	for i, model := range models {
		fmt.Printf("   %d. %s\n", i+1, model)
	}

	// Ask user to select model
	fmt.Print("\n🔧 Select model (1-", len(models), ") or press Enter for default (llama2): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		var selection int
		n, err := fmt.Sscanf(input, "%d", &selection)
		if err == nil && n == 1 && selection >= 1 && selection <= len(models) {
			config.Model = models[selection-1]
		}
	}

	fmt.Printf("🤖 Using model: %s\n", config.Model)
	return gitcommenter.New(config)
}

func displayChangesSummary(changes []gitcommenter.FileChange) {
	fmt.Printf("\n📊 Found %d staged file(s):\n", len(changes))

	totalAdded, totalRemoved := 0, 0
	for _, change := range changes {
		icon := getChangeIcon(change.ChangeType)
		fmt.Printf("   %s %s (+%d -%d)\n",
			icon, change.FilePath, change.LinesAdded, change.LinesRemoved)
		totalAdded += change.LinesAdded
		totalRemoved += change.LinesRemoved
	}

	fmt.Printf("\n📈 Total: +%d -%d lines\n", totalAdded, totalRemoved)
}

func getChangeIcon(changeType string) string {
	switch changeType {
	case "added":
		return "➕"
	case "modified":
		return "📝"
	case "deleted":
		return "🗑️"
	case "renamed":
		return "📛"
	default:
		return "📄"
	}
}

func displayCommitSuggestion(suggestion *gitcommenter.CommitSuggestion) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎯 AI-GENERATED COMMIT MESSAGE")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📝 Subject: %s\n", suggestion.Subject)

	if suggestion.Body != "" {
		fmt.Printf("\n📄 Body:\n%s\n", suggestion.Body)
	}

	fmt.Printf("\n📊 Confidence: %.0f%%\n", suggestion.Confidence*100)
	fmt.Printf("📁 Files: %s\n", strings.Join(suggestion.FilesAffected, ", "))
	fmt.Println(strings.Repeat("=", 60))
}

func askForApproval() bool {
	fmt.Print("\n❓ Do you want to commit with this message? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y" || response == "yes"
}

func commitWithMessage(suggestion *gitcommenter.CommitSuggestion) bool {
	args := []string{"commit", "-m", suggestion.Subject}
	if suggestion.Body != "" {
		args = append(args, "-m", suggestion.Body)
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run() == nil
}

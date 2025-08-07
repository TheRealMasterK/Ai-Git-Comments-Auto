package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	gitcommenter "github.com/TheRealMasterK/Ai-Git-Comments-Auto"
)

const (
	version = "1.0.0"
)

func main() {
	var (
		model       = flag.String("model", "llama2", "Ollama model to use")
		endpoint    = flag.String("endpoint", "http://localhost:11434", "Ollama endpoint")
		temperature = flag.Float64("temperature", 0.7, "Temperature for AI model (0.0-1.0)")
		maxTokens   = flag.Int("max-tokens", 150, "Maximum tokens for response")
		listModels  = flag.Bool("list-models", false, "List available Ollama models")
		interactive = flag.Bool("interactive", true, "Interactive mode to approve commit message (default: true)")
		skipAdd     = flag.Bool("skip-add", false, "Skip 'git add .' and only commit staged files")
		skipPush    = flag.Bool("skip-push", false, "Skip 'git push' after committing")
		dryRun      = flag.Bool("dry-run", false, "Show what would be done without executing")
		showVersion = flag.Bool("version", false, "Show version information")
		force       = flag.Bool("force", false, "Skip confirmation prompts")
	)
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("AI Git Auto v%s\n", version)
		fmt.Println("Automated Git workflow with AI-generated commit messages")
		return
	}

	// Print header
	fmt.Println("🚀 AI Git Auto - Automated Git Workflow")
	fmt.Println("======================================")

	// Create configuration
	config := &gitcommenter.Config{
		OllamaEndpoint: *endpoint,
		Model:         *model,
		MaxTokens:     *maxTokens,
		Temperature:   *temperature,
		RepositoryPath: ".",
	}

	// Create commenter
	commenter := gitcommenter.New(config)

	// List models if requested
	if *listModels {
		models, err := commenter.ListAvailableModels()
		if err != nil {
			log.Fatalf("❌ Failed to list models: %v", err)
		}

		fmt.Println("📚 Available Ollama models:")
		for _, model := range models {
			fmt.Printf("  - %s\n", model)
		}
		return
	}

	// Verify prerequisites
	if err := verifyPrerequisites(); err != nil {
		log.Fatalf("❌ %v", err)
	}

	// Get current directory for display
	pwd, _ := os.Getwd()
	fmt.Printf("📂 Working directory: %s\n", pwd)

	// Step 1: Git add (unless skipped)
	if !*skipAdd {
		fmt.Println("\n📝 Step 1: Staging changes...")
		if *dryRun {
			fmt.Println("   [DRY RUN] Would run: git add .")
		} else {
			if err := runGitAdd(); err != nil {
				log.Fatalf("❌ Failed to stage changes: %v", err)
			}
			fmt.Println("✅ Changes staged successfully")
		}
	} else {
		fmt.Println("\n📝 Step 1: Using already staged changes...")
	}

	// Step 2: Scan changes and generate commit message
	fmt.Println("\n🔍 Step 2: Scanning staged changes...")
	changes, err := commenter.ScanStagedChanges()
	if err != nil {
		log.Fatalf("❌ Failed to scan changes: %v", err)
	}

	if len(changes) == 0 {
		fmt.Println("📄 No staged changes found.")
		if !*skipAdd {
			fmt.Println("💡 Tip: Make sure you have changes to commit")
		} else {
			fmt.Println("💡 Tip: Stage your changes first with 'git add <files>'")
		}
		return
	}

	// Display changes summary
	displayChangesSummary(changes)

	fmt.Printf("\n🤖 Step 3: Generating AI commit message (using %s)...\n", *model)
	suggestion, err := commenter.GenerateCommitMessage(changes)
	if err != nil {
		log.Fatalf("❌ Failed to generate commit message: %v", err)
	}

	// Display the suggestion
	displayCommitSuggestion(suggestion)

	// Step 4: Commit
	fmt.Println("\n💾 Step 4: Committing changes...")
	commitApproved := !*interactive || *force || askForApproval("commit with this message")

	if *dryRun {
		fmt.Printf("   [DRY RUN] Would run: git commit -m \"%s\"", suggestion.Subject)
		if suggestion.Body != "" {
			fmt.Printf(" -m \"%s\"", suggestion.Body)
		}
		fmt.Println()
	} else if commitApproved {
		if err := runGitCommit(suggestion); err != nil {
			log.Fatalf("❌ Failed to commit: %v", err)
		}
		fmt.Println("✅ Changes committed successfully")
	} else {
		fmt.Println("❌ Commit cancelled by user")
		return
	}

	// Step 5: Push (unless skipped)
	if !*skipPush {
		fmt.Println("\n📤 Step 5: Pushing to remote...")

		// Check if there's a remote configured
		if !hasRemoteConfigured() {
			fmt.Println("⚠️  No remote repository configured, skipping push")
		} else {
			pushApproved := !*interactive || *force || askForApproval("push to remote")

			if *dryRun {
				fmt.Println("   [DRY RUN] Would run: git push")
			} else if pushApproved {
				if err := runGitPush(); err != nil {
					log.Printf("⚠️  Failed to push: %v", err)
					fmt.Println("💡 You can push manually later with: git push")
				} else {
					fmt.Println("✅ Changes pushed successfully")
				}
			} else {
				fmt.Println("📝 Push skipped. You can push manually with: git push")
			}
		}
	}

	fmt.Println("\n🎉 Workflow completed!")
}

func verifyPrerequisites() error {
	// Check if in git repository
	if !isGitRepository() {
		return fmt.Errorf("not in a Git repository")
	}

	// Check if Ollama is running
	config := gitcommenter.DefaultConfig()
	commenter := gitcommenter.New(config)
	if _, err := commenter.ListAvailableModels(); err != nil {
		return fmt.Errorf("Ollama is not running or not accessible at %s. Please start it with: ollama serve", config.OllamaEndpoint)
	}

	return nil
}

func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	_, err := cmd.Output()
	return err == nil
}

func hasRemoteConfigured() bool {
	cmd := exec.Command("git", "remote")
	output, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(output)) != ""
}

func runGitAdd() error {
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runGitCommit(suggestion *gitcommenter.CommitSuggestion) error {
	args := []string{"commit", "-m", suggestion.Subject}
	if suggestion.Body != "" {
		args = append(args, "-m", suggestion.Body)
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runGitPush() error {
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func displayChangesSummary(changes []gitcommenter.FileChange) {
	fmt.Printf("📊 Found %d staged file(s):\n", len(changes))

	totalAdded, totalRemoved := 0, 0
	for _, change := range changes {
		icon := getChangeIcon(change.ChangeType)
		fmt.Printf("   %s %s (+%d -%d)\n",
			icon, change.FilePath, change.LinesAdded, change.LinesRemoved)
		totalAdded += change.LinesAdded
		totalRemoved += change.LinesRemoved
	}

	fmt.Printf("📈 Total: +%d -%d lines\n", totalAdded, totalRemoved)
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
	fmt.Println(strings.Repeat("=", 60))
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

func askForApproval(action string) bool {
	fmt.Printf("❓ Do you want to %s? (Y/n): ", action)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	// Default to yes if empty response
	return response == "" || response == "y" || response == "yes"
}

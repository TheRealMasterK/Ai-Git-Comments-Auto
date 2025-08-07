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
	fmt.Println("🔍 Verifying prerequisites...")
	fmt.Println("   ➤ Checking Git repository...")
	if err := verifyPrerequisites(); err != nil {
		log.Fatalf("❌ %v", err)
	}
	fmt.Printf("   ✅ Git repository confirmed\n")

	// Check Ollama connection and model
	fmt.Printf("   ➤ Testing connection to Ollama at %s...\n", *endpoint)
	availableModels, err := commenter.ListAvailableModels()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Ollama: %v", err)
	}
	fmt.Printf("   ✅ Connected successfully (%d models available)\n", len(availableModels))

	// Verify selected model exists or let user choose
	modelExists := false
	for _, availableModel := range availableModels {
		if availableModel == *model {
			modelExists = true
			break
		}
	}

	if !modelExists {
		fmt.Printf("   ⚠️  Model '%s' not found.\n", *model)

		if len(availableModels) == 0 {
			log.Fatalf("❌ No Ollama models available. Please pull a model first:\n   ollama pull llama3.2")
		}

		// Interactive model selection
		fmt.Println("   📚 Available models:")
		for i, availableModel := range availableModels {
			recommendation := getModelRecommendation(availableModel)
			fmt.Printf("      %d. %s%s\n", i+1, availableModel, recommendation)
		}

		selectedModel, err := promptUserForModel(availableModels)
		if err != nil {
			log.Fatalf("❌ Model selection cancelled")
		}
		*model = selectedModel
	}

	fmt.Printf("   ✅ Using AI model: %s\n", *model)

	// Update config with selected model
	config.Model = *model

	// Get current directory for display
	pwd, _ := os.Getwd()
	fmt.Printf("   📂 Working directory: %s\n", pwd)

	// Step 1: Git add (unless skipped)
	if !*skipAdd {
		fmt.Println("\n📝 Step 1: Staging changes (git add .)...")

		// Show what files will be staged
		fmt.Println("   ➤ Checking for unstaged changes...")
		unstagedFiles, err := getUnstagedFiles()
		if err != nil {
			fmt.Printf("   ⚠️  Warning: Could not list unstaged files: %v\n", err)
		} else if len(unstagedFiles) > 0 {
			fmt.Printf("   ➤ Found %d unstaged file(s):\n", len(unstagedFiles))
			for i, file := range unstagedFiles {
				if i >= 5 { // Limit display to first 5 files
					fmt.Printf("      ... and %d more files\n", len(unstagedFiles)-5)
					break
				}
				fmt.Printf("      • %s\n", file)
			}
		} else {
			fmt.Println("   ➤ No unstaged files found")
		}

		if *dryRun {
			fmt.Println("   [DRY RUN] Would run: git add .")
		} else {
			fmt.Println("   ➤ Running: git add .")
			if err := runGitAdd(); err != nil {
				log.Fatalf("❌ Failed to stage changes: %v", err)
			}
			fmt.Println("   ✅ Changes staged successfully")
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
	fmt.Println("   ➤ Analyzing file changes and diffs...")
	fmt.Printf("   ➤ Sending context to Ollama model '%s'...\n", *model)

	suggestion, err := commenter.GenerateCommitMessage(changes)
	if err != nil {
		log.Fatalf("❌ Failed to generate commit message: %v", err)
	}

	fmt.Printf("   ✅ AI commit message generated (confidence: %.0f%%)\n", suggestion.Confidence*100)

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
		fmt.Println("   ➤ Running git commit...")
		if err := runGitCommit(suggestion); err != nil {
			log.Fatalf("❌ Failed to commit: %v", err)
		}
		fmt.Println("   ✅ Changes committed successfully")

		// Show commit hash
		if hash, err := getLastCommitHash(); err == nil {
			fmt.Printf("   📝 Commit hash: %s\n", hash)
		}
	} else {
		fmt.Println("   ❌ Commit cancelled by user")
		return
	}

	// Step 5: Push (unless skipped)
	if !*skipPush {
		fmt.Println("\n📤 Step 5: Pushing to remote...")

		// Check if there's a remote configured
		fmt.Println("   ➤ Checking for remote repositories...")
		remotes, err := getConfiguredRemotes()
		if err != nil || len(remotes) == 0 {
			fmt.Println("   ⚠️  No remote repository configured, skipping push")
			fmt.Println("   💡 Add a remote with: git remote add origin <url>")
		} else {
			fmt.Printf("   ➤ Found remote(s): %s\n", strings.Join(remotes, ", "))

			// Check current branch
			branch, err := getCurrentBranch()
			if err == nil {
				fmt.Printf("   ➤ Current branch: %s\n", branch)
			}

			pushApproved := !*interactive || *force || askForApproval("push to remote")

			if *dryRun {
				fmt.Println("   [DRY RUN] Would run: git push")
			} else if pushApproved {
				fmt.Println("   ➤ Running: git push")
				if err := runGitPush(); err != nil {
					log.Printf("   ⚠️  Failed to push: %v", err)
					fmt.Println("   💡 You can push manually later with: git push")
				} else {
					fmt.Println("   ✅ Changes pushed successfully")
				}
			} else {
				fmt.Println("   📝 Push skipped. You can push manually with: git push")
			}
		}
	} else {
		fmt.Println("\n📤 Step 5: Skipping push (--skip-push flag used)")
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
	fmt.Printf("   📊 Found %d staged file(s):\n", len(changes))

	totalAdded, totalRemoved := 0, 0
	filesByType := make(map[string]int)

	for _, change := range changes {
		icon := getChangeIcon(change.ChangeType)
		fmt.Printf("      %s %s (+%d -%d lines)\n",
			icon, change.FilePath, change.LinesAdded, change.LinesRemoved)
		totalAdded += change.LinesAdded
		totalRemoved += change.LinesRemoved
		filesByType[change.ChangeType]++
	}

	fmt.Printf("   📈 Total changes: +%d -%d lines\n", totalAdded, totalRemoved)

	// Show summary by change type
	var summary []string
	for changeType, count := range filesByType {
		if count > 0 {
			summary = append(summary, fmt.Sprintf("%d %s", count, changeType))
		}
	}
	if len(summary) > 0 {
		fmt.Printf("   📋 Summary: %s\n", strings.Join(summary, ", "))
	}
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

func getUnstagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}

	// Also get untracked files
	cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard")
	output, err = cmd.Output()
	if err == nil {
		untrackedLines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range untrackedLines {
			if line != "" {
				files = append(files, line+" (untracked)")
			}
		}
	}

	return files, nil
}

func getLastCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output))[:7], nil // Return short hash
}

func getConfiguredRemotes() ([]string, error) {
	cmd := exec.Command("git", "remote")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var remotes []string
	for _, line := range lines {
		if line != "" {
			remotes = append(remotes, line)
		}
	}
	return remotes, nil
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func promptUserForModel(availableModels []string) (string, error) {
	if len(availableModels) == 0 {
		return "", fmt.Errorf("no models available")
	}

	fmt.Print("\n   🤖 Please select a model (1-", len(availableModels), ") or press Enter for default: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)

	// If empty input, use first available model
	if input == "" {
		fmt.Printf("   ➤ Using default model: %s\n", availableModels[0])
		return availableModels[0], nil
	}

	// Parse selection
	var selection int
	n, err := fmt.Sscanf(input, "%d", &selection)
	if n != 1 || err != nil || selection < 1 || selection > len(availableModels) {
		fmt.Printf("   ❌ Invalid selection. Using default model: %s\n", availableModels[0])
		return availableModels[0], nil
	}

	selectedModel := availableModels[selection-1]
	fmt.Printf("   ➤ Selected model: %s\n", selectedModel)
	return selectedModel, nil
}

func getModelRecommendation(modelName string) string {
	modelLower := strings.ToLower(modelName)

	switch {
	case strings.Contains(modelLower, "llama3"):
		return " 🌟 (Recommended - Great for code)"
	case strings.Contains(modelLower, "codellama"):
		return " 💻 (Best for coding)"
	case strings.Contains(modelLower, "qwen"):
		if strings.Contains(modelLower, "32b") {
			return " 🚀 (Powerful but slow)"
		} else if strings.Contains(modelLower, "7b") {
			return " ⚡ (Good balance)"
		}
		return " 🧠 (Smart choice)"
	case strings.Contains(modelLower, "mistral"):
		return " ⚡ (Fast and efficient)"
	case strings.Contains(modelLower, "llama2"):
		return " 🏛️ (Reliable classic)"
	case strings.Contains(modelLower, "3b"):
		return " ⚡ (Fast and light)"
	case strings.Contains(modelLower, "7b"):
		return " ⚖️ (Balanced)"
	case strings.Contains(modelLower, "13b") || strings.Contains(modelLower, "32b"):
		return " 🐢 (Slow but accurate)"
	default:
		return ""
	}
}

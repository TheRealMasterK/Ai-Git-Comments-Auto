class AiGitAuto < Formula
  desc "Automate your Git workflow with AI-powered commit messages using local Ollama models"
  homepage "https://github.com/TheRealMasterK/Ai-Git-Comments-Auto"
  url "https://github.com/TheRealMasterK/Ai-Git-Comments-Auto/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "placeholder_sha256_will_be_updated_on_release"
  license "MIT"
  version "1.0.0"

  depends_on "go" => :build
  depends_on "git"

  def install
    # Build the binary
    system "go", "build", "-o", "ai-git-auto", "./cmd/ai-git-auto"
    
    # Install binary
    bin.install "ai-git-auto"
    
    # Install completion scripts
    generate_completions_from_executable(bin/"ai-git-auto", "completion")
  end

  def caveats
    <<~EOS
      🚀 AI Git Auto installed successfully!

      Prerequisites:
        1. Install Ollama: brew install ollama
        2. Start Ollama: brew services start ollama
        3. Install an AI model: ollama pull llama3.2:3b

      Usage:
        cd your-git-repo
        ai-git-auto

      The tool will automatically:
        • Stage your changes (git add .)
        • Generate AI commit messages
        • Commit and push your changes
    EOS
  end

  test do
    system "#{bin}/ai-git-auto", "--version"
  end
end

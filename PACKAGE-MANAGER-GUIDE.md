# 🚀 AI Git Auto - Package Manager Installation Guide

## ✨ New Installation Methods

We've added support for the most popular package managers! Choose your preferred method:

### 📦 Package Managers (Easiest!)

#### 🍺 Homebrew (macOS/Linux)
```bash
# One command installation
brew install TheRealMasterK/tap/ai-git-auto

# Then just use it
cd your-git-repo
ai-git-auto
```

**What Homebrew handles:**
- ✅ Downloads and builds the binary automatically
- ✅ Installs to your PATH
- ✅ Handles Go build dependency
- ✅ Shows helpful setup instructions

#### 📦 npm (Cross-platform)
```bash
# Global installation
npm install -g ai-git-auto

# Then just use it
cd your-git-repo
ai-git-auto
```

**What npm handles:**
- ✅ Checks system requirements (Node.js, Git, Go)
- ✅ Builds the Go binary automatically
- ✅ Installs globally to your PATH
- ✅ Shows setup instructions for Ollama

### 🔧 Script Installation (Previous method still available)
```bash
# One-click installation with all prerequisites
curl -fsSL https://raw.githubusercontent.com/TheRealMasterK/Ai-Git-Comments-Auto/main/install.sh | bash
```

## 🎯 Installation Comparison

| Method | Prerequisites Handled | Build Required | Global Install | Best For |
|--------|---------------------|----------------|---------------|----------|
| **Homebrew** | Go (build only) | ✅ Automatic | ✅ Yes | macOS/Linux users |
| **npm** | Node.js, checks Go | ✅ Automatic | ✅ Yes | Cross-platform devs |
| **Script** | Everything! | ✅ Automatic | ✅ Yes | First-time users |
| **Manual** | User handles | ❌ Manual | ❌ Manual | Developers |

## 📋 Prerequisites (All Methods)

**Required for AI functionality:**
1. **Ollama**: AI model hosting
   ```bash
   # macOS/Linux
   curl -fsSL https://ollama.ai/install.sh | sh

   # Or via Homebrew
   brew install ollama
   ```

2. **AI Model**: Choose one
   ```bash
   ollama pull llama3.2:3b     # Recommended: Fast, efficient
   ollama pull codellama:7b    # Code-specialized
   ollama pull mistral:7b      # General purpose
   ```

## 🚀 Usage (Same for all installation methods)

After installation with any method:

```bash
cd your-git-repository
ai-git-auto
```

**What happens:**
1. 📝 Stages your changes (`git add .`)
2. 🤖 Lets you choose AI model interactively
3. 🔍 Analyzes your code changes
4. ✨ Generates intelligent commit message
5. 💾 Commits and pushes automatically

## 🎉 For Developers: Publishing Updates

### Homebrew Release
```bash
# 1. Create GitHub release
# 2. Generate tarball
make brew-prepare

# 3. Update Formula/ai-git-auto.rb with new SHA256
# 4. Push to homebrew-tap repository
```

### npm Release
```bash
# 1. Update version in package.json
# 2. Prepare package
make npm-prepare

# 3. Test locally
npm pack && npm install -g ai-git-auto-*.tgz

# 4. Publish
npm login
npm publish
```

## 🔍 Testing Your Installation

```bash
# Check version
ai-git-auto --version

# Test in a repo (dry run)
cd /path/to/git/repo
ai-git-auto --dry-run

# Full test
ai-git-auto
```

## 🎯 Success!

**Mission accomplished!** Developers can now install AI Git Auto using their preferred package manager:

- **Mac developers**: `brew install TheRealMasterK/tap/ai-git-auto`
- **Node.js developers**: `npm install -g ai-git-auto`
- **Everyone else**: `curl -fsSL install.sh | bash`

All methods result in the same powerful AI-driven Git workflow automation! 🚀

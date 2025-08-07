# 🍺 Publishing to Homebrew - Step by Step Guide

## Prerequisites

Before publishing to Homebrew, you'll need:

1. **GitHub repository for the tap**: `TheRealMasterK/homebrew-tap`
2. **Release on main repo**: Create a v1.0.0 release on `Ai-Git-Comments-Auto`
3. **Homebrew installed**: `brew --version`

## Step 1: Create Homebrew Tap Repository

```bash
# Create a new repository on GitHub named: homebrew-tap
# Repository URL should be: https://github.com/TheRealMasterK/homebrew-tap

# Clone it locally
git clone https://github.com/TheRealMasterK/homebrew-tap.git
cd homebrew-tap

# Create Formula directory
mkdir -p Formula
```

## Step 2: Copy the Formula

```bash
# Copy our formula to the tap repository
cp /Users/kylelloyd/Documents/GitHub/Ai-Git-Comments-Auto/Formula/ai-git-auto.rb Formula/

# Commit and push
git add Formula/ai-git-auto.rb
git commit -m "Add ai-git-auto formula"
git push origin main
```

## Step 3: Create GitHub Release

```bash
# In your main repository (Ai-Git-Comments-Auto)
cd /Users/kylelloyd/Documents/GitHub/Ai-Git-Comments-Auto

# Create a git tag
git tag -a v1.0.0 -m "Release v1.0.0 with Homebrew and npm support"
git push origin v1.0.0

# Create a GitHub release from the tag v1.0.0
# GitHub will automatically create the source tarball at:
# https://github.com/TheRealMasterK/Ai-Git-Comments-Auto/archive/refs/tags/v1.0.0.tar.gz
```

## Step 4: Test the Homebrew Formula

```bash
# Test the formula locally
brew install --build-from-source TheRealMasterK/tap/ai-git-auto

# Or test without installing
brew install --build-from-source ./Formula/ai-git-auto.rb --verbose

# Test that it works
ai-git-auto --version
```

## Step 5: Users Can Now Install

After completing steps 1-3, users can install with:

```bash
# Add the tap (one time)
brew tap TheRealMasterK/tap

# Install the package
brew install ai-git-auto
```

## Commands to Run

Here are the exact commands you need to run:

```bash
# 1. Create homebrew-tap repo on GitHub, then:
git clone https://github.com/TheRealMasterK/homebrew-tap.git
cd homebrew-tap
mkdir -p Formula

# 2. Copy formula
cp /Users/kylelloyd/Documents/GitHub/Ai-Git-Comments-Auto/Formula/ai-git-auto.rb Formula/
git add Formula/ai-git-auto.rb
git commit -m "feat: add ai-git-auto formula v1.0.0"
git push origin main

# 3. Create release tag
cd /Users/kylelloyd/Documents/GitHub/Ai-Git-Comments-Auto
git tag -a v1.0.0 -m "Release v1.0.0: Homebrew and npm support"
git push origin v1.0.0

# 4. Create GitHub release from web interface using tag v1.0.0
```

## Formula Details

Our formula at `Formula/ai-git-auto.rb` has:
- ✅ Correct URL pointing to GitHub release tarball
- ✅ SHA256 hash: `b856de7b67244eaf3bec8aecfe262110538dd5e1a1e52d6424f24ff09264277d`
- ✅ Go build dependency
- ✅ Installation instructions in caveats
- ✅ Test command

## Next Steps

After Homebrew is published, users will be able to:

```bash
brew tap TheRealMasterK/tap
brew install ai-git-auto
ai-git-auto --version
```

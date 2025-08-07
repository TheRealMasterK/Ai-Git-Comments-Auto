# 🍺 Homebrew Tap for AI Git Auto

This directory contains the Homebrew formula for AI Git Auto.

## For Users

### Install AI Git Auto via Homebrew

```bash
# Add the tap (only needed once)
brew tap TheRealMasterK/tap

# Install AI Git Auto
brew install ai-git-auto
```

### Prerequisites

The formula will automatically handle Go as a build dependency, but you'll still need:

1. **Ollama**: Install with `brew install ollama`
2. **AI Model**: Download with `ollama pull llama3.2:3b`

### Usage

After installation:
```bash
cd your-git-repository
ai-git-auto
```

## For Maintainers

### Updating the Formula

1. Update the version and URL in `Formula/ai-git-auto.rb`
2. Calculate new SHA256: `shasum -a 256 ai-git-auto-VERSION.tar.gz`
3. Update the `sha256` field in the formula
4. Test the formula: `brew install --build-from-source ./Formula/ai-git-auto.rb`

### Release Process

1. Create a new release on GitHub
2. Generate the tarball: `make brew-prepare`
3. Update the formula with new version and SHA256
4. Push to the tap repository

## Formula Location

This formula should be placed in a separate Homebrew tap repository:
- Repository: `TheRealMasterK/homebrew-tap`
- Formula path: `Formula/ai-git-auto.rb`

## Testing

```bash
# Test the formula locally
brew install --build-from-source ./Formula/ai-git-auto.rb

# Test that it works
ai-git-auto --version
```

# Contributing to AI Git Auto

We welcome contributions! Here's how to get started:

## 🚀 Quick Setup for Contributors

1. **Fork and Clone**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/Ai-Git-Comments-Auto.git
   cd Ai-Git-Comments-Auto
   ```

2. **Install Dependencies** (our script makes this easy):
   ```bash
   # One command to set up everything
   ./install.sh
   ```

3. **Build and Test**:
   ```bash
   go mod tidy
   go build -o ai-git-auto ./cmd/ai-git-auto
   ./ai-git-auto --help
   ```

## 🛠️ Development Guidelines

### Code Style
- Follow standard Go conventions
- Use `gofmt` and `golint`
- Add comments for exported functions
- Keep functions focused and small

### Testing
```bash
go test ./...
```

### Building
```bash
# Build for current platform
go build -o ai-git-auto ./cmd/ai-git-auto

# Build for multiple platforms
make build-all
```

## 📝 Making Changes

1. Create a new branch: `git checkout -b feature/your-feature`
2. Make your changes
3. Test thoroughly
4. Commit using conventional commits:
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation
   - `refactor:` for code refactoring

## 🎯 What We're Looking For

- **Bug fixes**: Always welcome
- **New AI model integrations**: Support for more Ollama models
- **Platform support**: Windows support would be great
- **Performance improvements**: Faster git operations, better error handling
- **Documentation**: Examples, tutorials, better README

## 📚 Project Structure

```
├── cmd/
│   └── ai-git-auto/     # Main CLI application
├── examples/            # Usage examples
├── gitcommenter.go      # Core library
├── go.mod              # Go module definition
├── install.sh          # One-click installation script
├── Makefile           # Build automation
└── README.md          # Project documentation
```

## 🧪 Testing Your Changes

Before submitting a PR, make sure to:

1. **Test the core functionality**:
   ```bash
   # In a test git repo
   echo "test change" > test.txt
   git add .
   ./ai-git-auto --dry-run
   ```

2. **Test the installation script**:
   ```bash
   # Test in a clean environment or Docker container
   ./install.sh
   ```

3. **Test with different Ollama models**:
   ```bash
   ollama pull llama3.2:3b
   ollama pull codellama:7b
   # Test model selection works
   ```

## 🚀 Submitting Changes

1. Push your changes to your fork
2. Create a Pull Request with:
   - Clear description of what you changed
   - Why the change was needed
   - Any breaking changes
   - Screenshots if UI-related

## ❓ Need Help?

- Open an issue for bugs or questions
- Check existing issues before creating new ones
- Be detailed in issue descriptions

## 🎉 Recognition

Contributors will be added to the README and releases. Thank you for helping make AI Git Auto better!

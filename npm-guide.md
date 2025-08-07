# 📦 npm Package Publication Guide

## For Users

### Install AI Git Auto via npm

```bash
npm install -g ai-git-auto
```

### Prerequisites

The npm package will handle building the Go binary automatically, but you'll need:

1. **Node.js**: Version 14 or higher
2. **Go**: For building the binary (will be checked during installation)
3. **Ollama**: Install from https://ollama.ai/
4. **AI Model**: Download with `ollama pull llama3.2:3b`

### Usage

After installation:
```bash
cd your-git-repository
ai-git-auto
```

## For Maintainers

### Publishing to npm

1. **Prepare the package**:
   ```bash
   make npm-prepare
   ```

2. **Test locally**:
   ```bash
   npm pack
   npm install -g ai-git-auto-1.0.0.tgz
   ai-git-auto --version
   ```

3. **Publish to npm**:
   ```bash
   npm login
   npm publish
   ```

### Release Process

1. Update version in `package.json`
2. Update version in `Makefile`
3. Update version in `cmd/ai-git-auto/main.go`
4. Run `make npm-prepare` to test build
5. Commit changes and create git tag
6. Publish: `npm publish`

### Package Structure

```
├── package.json          # npm package configuration
├── bin/
│   └── ai-git-auto       # Built binary (created by postinstall)
├── scripts/
│   ├── preinstall.js     # System requirements check
│   └── postinstall.js    # Build and install binary
├── install.sh            # Fallback installation script
└── README.md             # Documentation
```

## Testing

### Local Testing

```bash
# Test the preinstall script
node scripts/preinstall.js

# Test building
make npm-prepare

# Test the postinstall script
node scripts/postinstall.js
```

### Integration Testing

```bash
# Pack the npm package
npm pack

# Install globally from the tarball
npm install -g ai-git-auto-*.tgz

# Test it works
ai-git-auto --version

# Test in a git repository
cd /path/to/git/repo
ai-git-auto --dry-run
```

## Troubleshooting

### Common Issues

1. **Go not found**: The postinstall script will show an error if Go isn't installed
2. **Binary build fails**: Check Go version and dependencies
3. **Permission errors**: May need `sudo npm install -g` on some systems
4. **Platform not supported**: Only macOS and Linux are supported

### Debug Installation

```bash
# Install with verbose logging
npm install -g ai-git-auto --verbose

# Check installed files
npm list -g --depth=0 | grep ai-git-auto
which ai-git-auto
```

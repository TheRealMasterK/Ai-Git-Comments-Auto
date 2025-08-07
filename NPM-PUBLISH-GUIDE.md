# 📦 Publishing to npm - Step by Step Guide

## Prerequisites

Before publishing to npm, you'll need:

1. **npm account**: Create at https://www.npmjs.com/signup
2. **npm CLI**: Already installed with Node.js
3. **Unique package name**: `ai-git-auto` (we'll check availability)

## Step 1: Set up npm Account

```bash
# Create account on npmjs.com first, then login
npm login

# Enter your credentials:
# Username: your-npm-username
# Password: your-npm-password
# Email: your-email@example.com
# OTP (if 2FA enabled): 123456

# Verify you're logged in
npm whoami
```

## Step 2: Check Package Name Availability

```bash
# Check if 'ai-git-auto' is available
npm info ai-git-auto

# If it returns "npm ERR! 404 'ai-git-auto@latest' is not in the npm registry"
# then the name is available!

# If it exists, you'll need to choose a different name like:
# - ai-git-commit-auto
# - ai-git-committer
# - git-ai-auto
# - ollama-git-auto
```

## Step 3: Prepare for Publication

```bash
# Make sure everything is built and ready
cd /Users/kylelloyd/Documents/GitHub/Ai-Git-Comments-Auto

# Clean and rebuild
make clean
make npm-prepare

# Test the package locally
npm pack
npm install -g ./ai-git-auto-1.0.0.tgz

# Test it works
ai-git-auto --version
ai-git-auto --help

# If testing passed, uninstall the test version
npm uninstall -g ai-git-auto
```

## Step 4: Publish to npm

```bash
# Publish the package
npm publish

# If successful, you'll see output like:
# + ai-git-auto@1.0.0
```

## Step 5: Verify Publication

```bash
# Check your package is available
npm info ai-git-auto

# Test installation from npm
npm install -g ai-git-auto

# Test it works
ai-git-auto --version
```

## Commands to Run

Here are the exact commands you need to run:

```bash
# 1. Login to npm (do this first)
npm login
npm whoami  # verify login

# 2. Check name availability
npm info ai-git-auto

# 3. Prepare package
cd /Users/kylelloyd/Documents/GitHub/Ai-Git-Comments-Auto
make npm-prepare

# 4. Test locally
npm pack
npm install -g ./ai-git-auto-1.0.0.tgz
ai-git-auto --version
npm uninstall -g ai-git-auto

# 5. Publish
npm publish

# 6. Verify
npm info ai-git-auto
npm install -g ai-git-auto
ai-git-auto --version
```

## Troubleshooting

### Common Issues:

1. **Name conflict**: If `ai-git-auto` is taken, try:
   ```bash
   # Update package.json with new name
   # Then publish with new name
   ```

2. **Permission errors**: Make sure you're logged in:
   ```bash
   npm logout
   npm login
   ```

3. **2FA required**: Enable and use OTP:
   ```bash
   npm publish --otp=123456
   ```

4. **Binary missing**: Make sure `bin/ai-git-auto` exists:
   ```bash
   ls -la bin/
   make npm-prepare  # rebuilds binary
   ```

## After Publication

Once published, users worldwide can install with:

```bash
npm install -g ai-git-auto
```

## Package Details

Our npm package includes:
- ✅ Binary in `bin/ai-git-auto`
- ✅ Preinstall system check
- ✅ Postinstall binary building
- ✅ Global installation support
- ✅ Cross-platform (macOS/Linux)
- ✅ Size: ~4.8MB

## Next Steps

After npm publication:

1. Update README.md with npm installation instructions
2. Test installation on different systems
3. Monitor npm stats at https://www.npmjs.com/package/ai-git-auto
4. Consider setting up automated releases with GitHub Actions

#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

const colors = {
    red: '\x1b[31m',
    green: '\x1b[32m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    purple: '\x1b[35m',
    reset: '\x1b[0m'
};

function log(color, emoji, message) {
    console.log(`${colors[color]}${emoji} ${message}${colors.reset}`);
}

function runCommand(command, description) {
    try {
        log('blue', '🔄', description);
        execSync(command, { stdio: 'inherit' });
        log('green', '✅', `${description} completed`);
        return true;
    } catch (error) {
        log('red', '❌', `Failed: ${description}`);
        return false;
    }
}

console.log(`${colors.blue}
┌─────────────────────────────────────────────────────────┐
│                                                         │
│      🚀 AI Git Auto - npm Installation Complete        │
│                                                         │
└─────────────────────────────────────────────────────────┘
${colors.reset}`);

// Check if the binary exists (should be included in npm package)
const binDir = path.join(__dirname, '..', 'bin');
const platform = os.platform();
const binaryName = platform === 'win32' ? 'ai-git-auto.exe' : 'ai-git-auto';
const binaryPath = path.join(binDir, binaryName);

if (fs.existsSync(binaryPath)) {
    log('green', '✅', 'AI Git Auto binary is ready!');
    
    // Make binary executable on Unix systems
    if (platform !== 'win32') {
        try {
            execSync(`chmod +x "${binaryPath}"`);
        } catch (error) {
            log('yellow', '⚠️', 'Could not make binary executable');
        }
    }
} else {
    log('red', '❌', 'Binary not found. Attempting to build from source...');
    
    // Try to build from source as fallback
    try {
        if (!fs.existsSync(binDir)) {
            fs.mkdirSync(binDir, { recursive: true });
        }
        
        log('blue', '🔄', 'Building AI Git Auto binary...');
        const buildCommand = `go build -o "${binaryPath}" ./cmd/ai-git-auto`;
        execSync(buildCommand, { 
            cwd: path.join(__dirname, '..'),
            stdio: 'inherit' 
        });
        log('green', '✅', 'Binary built successfully');
        
        // Make binary executable on Unix systems
        if (platform !== 'win32') {
            try {
                execSync(`chmod +x "${binaryPath}"`);
            } catch (error) {
                log('yellow', '⚠️', 'Could not make binary executable');
            }
        }
    } catch (error) {
        log('red', '❌', 'Failed to build binary. Make sure Go is installed.');
        process.exit(1);
    }
}// Check if binary works
try {
    execSync(`"${binaryPath}" --version`, { stdio: 'ignore' });
    log('green', '✅', 'AI Git Auto installed successfully!');
} catch (error) {
    log('red', '❌', 'Binary installation verification failed');
}

console.log(`
${colors.yellow}📋 Next Steps:${colors.reset}

1. Install Ollama (if not already installed):
   ${colors.blue}curl -fsSL https://ollama.ai/install.sh | sh${colors.reset}

2. Install an AI model:
   ${colors.blue}ollama pull llama3.2:3b${colors.reset}

3. Navigate to any Git repository:
   ${colors.blue}cd /path/to/your/project${colors.reset}

4. Run AI Git Auto:
   ${colors.blue}ai-git-auto${colors.reset}

${colors.green}🚀 Ready to automate your Git workflow with AI! 🚀${colors.reset}
`);

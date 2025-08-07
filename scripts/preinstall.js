#!/usr/bin/env node

const os = require('os');
const { execSync } = require('child_process');

console.log('🔍 Checking system requirements for AI Git Auto...');

// Check operating system
const platform = os.platform();
if (platform !== 'darwin' && platform !== 'linux') {
    console.error('❌ AI Git Auto only supports macOS and Linux');
    process.exit(1);
}

console.log(`✅ Platform supported: ${platform}`);

// Check if git is installed
try {
    execSync('git --version', { stdio: 'ignore' });
    console.log('✅ Git is installed');
} catch (error) {
    console.error('❌ Git is not installed. Please install Git first.');
    process.exit(1);
}

// Check if Go is installed (needed for building)
try {
    execSync('go version', { stdio: 'ignore' });
    console.log('✅ Go is installed');
} catch (error) {
    console.log('⚠️  Go not found. Will attempt to install during build process.');
}

console.log('✅ System requirements check passed');

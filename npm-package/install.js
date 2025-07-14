#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');

// Determine platform
const platform = os.platform();
const arch = os.arch();

// Map to GoReleaser format for choosing binaries
const goPlatform = platform === 'win32' ? 'windows' : platform;
const goArch = arch === 'x64' ? 'amd64' : arch;

// Find the correct binary directory (GoReleaser adds arch version suffixes)
const binDirRoot = path.join(__dirname, 'bin');
let binPath;
let binDir;

if (fs.existsSync(binDirRoot)) {
  const dirs = fs.readdirSync(binDirRoot);
  const matchingDir = dirs.find(dir => {
    // Remove the 'lamacli_' prefix and split by '_'
    const withoutPrefix = dir.replace(/^lamacli_/, '');
    const parts = withoutPrefix.split('_');
    return parts.length >= 2 && parts[0] === goPlatform && parts[1].startsWith(goArch);
  });
  
  if (matchingDir) {
    binDir = path.join(binDirRoot, matchingDir);
    binPath = path.join(binDir, 'lamacli' + (platform === 'win32' ? '.exe' : ''));
  }
}

// Check if binary exists
if (!binPath || !fs.existsSync(binPath)) {
  console.error(`Binary not found for platform ${goPlatform}_${goArch}`);
  console.error('Available binaries:');
  
  if (fs.existsSync(binDirRoot)) {
    const dirs = fs.readdirSync(binDirRoot);
    dirs.forEach(dir => {
      const files = fs.readdirSync(path.join(binDirRoot, dir));
      console.error(`  ${dir}: ${files.join(', ')}`);
    });
  } else {
    console.error('  No bin directory found');
  }
  
  process.exit(1);
}

// Make binary executable (Unix-like systems)
if (platform !== 'win32') {
  try {
    fs.chmodSync(binPath, '755');
    console.log(`Made ${binPath} executable`);
  } catch (error) {
    console.error('Failed to make binary executable:', error);
    process.exit(1);
  }
}

console.log(`LamaCLI binary installed successfully for ${goPlatform}_${goArch}`);
console.log(`Binary location: ${binPath}`);
console.log('You can now use "lamacli" or "npx lamacli" to run the tool.');

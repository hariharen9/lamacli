#!/usr/bin/env node

const { execSync } = require('child_process');
const path = require('path');
const os = require('os');

// Determine platform
const platform = os.platform();
const arch = os.arch();

// Map to GoReleaser format for choosing binaries
const goPlatform = platform === 'win32' ? 'windows' : platform;
const goArch = arch === 'x64' ? 'amd64' : arch;

// Find the correct binary directory (GoReleaser adds arch version suffixes)
const fs = require('fs');
const binDir = path.join(__dirname, 'bin');
let binPath;

if (fs.existsSync(binDir)) {
  const dirs = fs.readdirSync(binDir);
  const matchingDir = dirs.find(dir => {
    // Remove the 'lamacli_' prefix and split by '_'
    const withoutPrefix = dir.replace(/^lamacli_/, '');
    const parts = withoutPrefix.split('_');
    return parts.length >= 2 && parts[0] === goPlatform && parts[1].startsWith(goArch);
  });
  
  if (matchingDir) {
    binPath = path.join(binDir, matchingDir, 'lamacli' + (platform === 'win32' ? '.exe' : ''));
  }
}

if (!binPath || !fs.existsSync(binPath)) {
  console.error(`Binary not found for platform ${goPlatform}_${goArch}`);
  console.error('Available binaries:');
  if (fs.existsSync(binDir)) {
    const dirs = fs.readdirSync(binDir);
    dirs.forEach(dir => {
      const files = fs.readdirSync(path.join(binDir, dir));
      console.error(`  ${dir}: ${files.join(', ')}`);
    });
  }
  process.exit(1);
}

try {
  // Pass all command line arguments to the binary
  const args = process.argv.slice(2);
  const command = `"${binPath}" ${args.map(arg => `"${arg}"`).join(' ')}`;
  execSync(command, { stdio: 'inherit' });
} catch (error) {
  console.error('Failed to execute LamaCLI binary:', error);
  process.exit(1);
}

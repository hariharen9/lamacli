# LamaCLI npm Distribution Setup Guide

This guide explains how to set up LamaCLI for distribution via npm, allowing users to install and use your Go CLI tool with `npm install -g lamacli` or `npx lamacli`.

## Files Created

### 1. GoReleaser Configuration (`.goreleaser.yml`)
- Configures multi-platform builds for Linux, macOS, and Windows
- Supports amd64 and arm64 architectures
- Creates archives and checksums

### 2. npm Package Structure (`npm-package/`)
- `package.json` - npm package configuration
- `index.js` - JavaScript shim that executes the appropriate binary
- `install.js` - Post-install script that sets up the binary for the current platform
- `README.md` - npm package documentation
- `LICENSE` - License file (copied from your main LICENSE.md)

### 3. Build Scripts
- `build-release.sh` - Bash script for Unix-like systems
- `build-release.ps1` - PowerShell script for Windows
- `test-npm-package.ps1` - Test script to verify setup works locally

### 4. GitHub Actions (`.github/workflows/release.yml`)
- Automated workflow for building and publishing on tag creation
- Handles both GitHub releases and npm publishing

## Prerequisites

1. **GoReleaser**: Install GoReleaser
   ```bash
   # On macOS/Linux
   curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
   
   # On Windows with Scoop
   scoop install goreleaser
   ```

2. **npm Account**: Create an npm account and get an authentication token
   ```bash
   npm login
   npm token create --access public
   ```

3. **GitHub Secrets**: Add the following secrets to your GitHub repository:
   - `NPM_TOKEN` - Your npm authentication token

## How It Works

1. **GoReleaser** builds binaries for multiple platforms
2. **Build script** copies binaries to the npm package structure
3. **npm publish** publishes the package to npm registry
4. **Users install** with `npm install -g lamacli` or use `npx lamacli`
5. **install.js** runs post-install to set up the correct binary for the user's platform
6. **index.js** acts as a shim to execute the appropriate binary

## Usage

### Local Development & Testing

1. Test the setup locally:
   ```powershell
   # On Windows
   .\test-npm-package.ps1
   
   # On Unix-like systems
   ./test-npm-package.sh
   ```

### Manual Release

1. Build and publish manually:
   ```powershell
   # On Windows
   .\build-release.ps1
   
   # On Unix-like systems
   ./build-release.sh
   ```

### Automated Release via GitHub Actions

1. Create a new tag and push it:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The GitHub Action will automatically:
   - Build binaries with GoReleaser
   - Create a GitHub release
   - Publish to npm

## Binary Structure

After building, the npm package will have this structure:
```
npm-package/
├── bin/
│   ├── linux_amd64/
│   │   └── lamacli
│   ├── linux_arm64/
│   │   └── lamacli
│   ├── darwin_amd64/
│   │   └── lamacli
│   ├── darwin_arm64/
│   │   └── lamacli
│   └── windows_amd64/
│       └── lamacli.exe
├── index.js
├── install.js
├── package.json
├── README.md
└── LICENSE
```

## User Experience

Once published, users can:

1. **Install globally:**
   ```bash
   npm install -g lamacli
   lamacli
   ```

2. **Use with npx:**
   ```bash
   npx lamacli
   ```

3. **Use in package.json scripts:**
   ```json
   {
     "scripts": {
       "chat": "lamacli"
     }
   }
   ```

## Troubleshooting

### Common Issues

1. **Binary not found**: Ensure the GoReleaser build completed successfully
2. **Permission denied**: The install.js script should make binaries executable on Unix systems
3. **Architecture mismatch**: Check that your GoReleaser config includes the target platform/architecture

### Debugging

1. Check the post-install output:
   ```bash
   npm install lamacli --verbose
   ```

2. Verify binary exists:
   ```bash
   ls -la $(npm root -g)/lamacli/bin/
   ```

3. Test the shim directly:
   ```bash
   node $(npm root -g)/lamacli/index.js --help
   ```

## Version Management

To update the version:

1. Update `npm-package/package.json`
2. Create a new git tag
3. Push the tag to trigger the release workflow

The build scripts will automatically sync the version from package.json.

## Package Size Considerations

The npm package will be larger than typical Node.js packages because it includes binaries for all platforms. Consider:

- Total size will be ~50-100MB depending on your binary size
- Users only execute one binary but download all
- This is normal for Go CLI tools distributed via npm

## Alternative Approaches

If package size is a concern, consider:
1. Platform-specific packages (e.g., `lamacli-darwin`, `lamacli-linux`)
2. Download binaries on-demand during post-install
3. Use a package manager detection system to install only the needed binary

This setup provides a complete solution for distributing your Go CLI tool via npm while maintaining the simplicity of `npm install -g lamacli` for end users.

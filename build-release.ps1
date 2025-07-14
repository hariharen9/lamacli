# This script orchestrates the Go build, and packages it for npm on Windows.

Write-Host "Building LamaCLI for multi-platform distribution..."

# Run GoReleaser to build the binaries
Write-Host "Running GoReleaser..."
goreleaser release --clean --skip-publish --debug

if ($LASTEXITCODE -ne 0) {
    Write-Error "GoReleaser failed with exit code $LASTEXITCODE"
    exit 1
}

# Create necessary directories
Write-Host "Creating npm package directories..."
New-Item -ItemType Directory -Force -Path "npm-package/bin"

# Get version from package.json
$packageJsonContent = Get-Content "npm-package/package.json" | ConvertFrom-Json
$version = $packageJsonContent.version

Write-Host "Current version: $version"

# Copy the binaries to the npm package bin directory
Write-Host "Copying binaries to npm package..."
Get-ChildItem -Path "dist" -Directory | ForEach-Object {
    $dirName = $_.Name
    if ($dirName -match "^lamacli_.*") {
        # Use the original directory name from GoReleaser
        $binDir = "npm-package/bin/$dirName"
        New-Item -ItemType Directory -Force -Path $binDir
        
        # Copy the binary (handle .exe for Windows)
        $sourceBinary = Join-Path $_.FullName "lamacli"
        $sourceExe = Join-Path $_.FullName "lamacli.exe"
        
        if (Test-Path $sourceExe) {
            Copy-Item $sourceExe (Join-Path $binDir "lamacli.exe")
            Write-Host "Copied Windows binary: $sourceExe -> $binDir/lamacli.exe"
        } elseif (Test-Path $sourceBinary) {
            Copy-Item $sourceBinary (Join-Path $binDir "lamacli")
            Write-Host "Copied binary: $sourceBinary -> $binDir/lamacli"
        } else {
            Write-Warning "No binary found in $($_.FullName)"
        }
    }
}

# Update npm package version
Write-Host "Updating npm package version..."
$packageJsonContent.version = $version
$packageJsonContent | ConvertTo-Json -Depth 10 | Set-Content "npm-package/package.json"

# Publish to npm
Write-Host "Publishing to npm..."
Set-Location "npm-package"
npm publish --access public

if ($LASTEXITCODE -eq 0) {
    Write-Host "Successfully published LamaCLI v$version to npm!" -ForegroundColor Green
} else {
    Write-Error "Failed to publish to npm with exit code $LASTEXITCODE"
    exit 1
}

Set-Location ".."
Write-Host "Build and publish complete!"

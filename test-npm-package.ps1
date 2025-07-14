# Test script to verify the npm package setup works locally

Write-Host "Testing LamaCLI npm package setup..." -ForegroundColor Green

# Build a test version with GoReleaser
Write-Host "Building test binaries..." -ForegroundColor Yellow
goreleaser build --clean --snapshot

if ($LASTEXITCODE -ne 0) {
    Write-Error "GoReleaser build failed"
    exit 1
}

# Create test npm package structure
Write-Host "Setting up test npm package..." -ForegroundColor Yellow
$testDir = "test-npm-package"
Remove-Item -Recurse -Force $testDir -ErrorAction SilentlyContinue
Copy-Item -Recurse "npm-package" $testDir

# Copy test binaries
New-Item -ItemType Directory -Force -Path "$testDir/bin"
Get-ChildItem -Path "dist" -Directory | ForEach-Object {
    $dirName = $_.Name
    if ($dirName -match "^lamacli_.*") {
        # Use the original directory name from GoReleaser
        $binDir = "$testDir/bin/$dirName"
        New-Item -ItemType Directory -Force -Path $binDir
        
        $sourceBinary = Join-Path $_.FullName "lamacli"
        $sourceExe = Join-Path $_.FullName "lamacli.exe"
        
        if (Test-Path $sourceExe) {
            Copy-Item $sourceExe (Join-Path $binDir "lamacli.exe")
            Write-Host "Copied test binary: $sourceExe -> $binDir/lamacli.exe"
        } elseif (Test-Path $sourceBinary) {
            Copy-Item $sourceBinary (Join-Path $binDir "lamacli")
            Write-Host "Copied test binary: $sourceBinary -> $binDir/lamacli"
        }
    }
}

# Test the package
Write-Host "Testing package installation..." -ForegroundColor Yellow
Set-Location $testDir
node install.js

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Package installation test passed!" -ForegroundColor Green
    
    Write-Host "Testing binary execution..." -ForegroundColor Yellow
    node index.js --help
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Binary execution test passed!" -ForegroundColor Green
    } else {
        Write-Host "❌ Binary execution test failed!" -ForegroundColor Red
    }
} else {
    Write-Host "❌ Package installation test failed!" -ForegroundColor Red
}

Set-Location ".."
Write-Host "Test completed. Check $testDir directory for results." -ForegroundColor Blue

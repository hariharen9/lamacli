#!/bin/bash

# This script orchestrates the Go build, and packages it for npm.

# Run GoReleaser to build the binaries.
goreleaser release --clean --skip-publish --debug

# Create necessary directories
mkdir -p npm-package/bin

# Get version from package.json
VERSION=$(jq -r '.version' npm-package/package.json)

# Copy the binaries to the npm package bin directory
for dir in dist/lamacli_*; do
  if [ -d "$dir" ]; then
    dirname=$(basename "$dir")
    mkdir -p "npm-package/bin/$dirname"
    if [ -f "$dir/lamacli.exe" ]; then
      cp "$dir/lamacli.exe" "npm-package/bin/$dirname/lamacli.exe"
    else
      cp "$dir/lamacli" "npm-package/bin/$dirname/lamacli"
      chmod +x "npm-package/bin/$dirname/lamacli"
    fi
  fi
done

# Update npm package version
jq --arg VERSION "$VERSION" '.version = $VERSION' npm-package/package.json > npm-package/temp.json && mv npm-package/temp.json npm-package/package.json

# Publish to npm
cd npm-package
npm publish --access public


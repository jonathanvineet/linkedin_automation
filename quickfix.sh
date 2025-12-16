#!/bin/bash

# One-Command Fix - Run this if setup.sh failed
# This handles: build tools installation + compilation + verification

echo "⚡ One-Command Fix"
echo "=================="
echo ""
echo "This will:"
echo "  1. Install build-essential (gcc, make, etc.)"
echo "  2. Download Go dependencies"
echo "  3. Build the Go backend with CGO enabled"
echo "  4. Verify the build"
echo ""

read -p "Continue? [Y/n] " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]] && [[ ! -z $REPLY ]]; then
    exit 1
fi

set -e  # Exit on error

# Install build tools
echo "📦 Installing build tools..."
sudo apt-get update -qq
sudo apt-get install -y build-essential
echo "✅ Build tools installed"
echo ""

# Make scripts executable
chmod +x *.sh 2>/dev/null || true

# Create directories
mkdir -p bin data logs
echo "✅ Directories created"
echo ""

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod download
go mod tidy
echo "✅ Dependencies ready"
echo ""

# Build with CGO
echo "🔨 Building Go backend..."
export CGO_ENABLED=1
go build -v -o bin/automation ./cmd/app
echo ""
echo "✅ Build successful!"
echo ""

# Verify
if [ -f "bin/automation" ] && [ -x "bin/automation" ]; then
    SIZE=$(du -h bin/automation | cut -f1)
    echo "✅ Binary verified: bin/automation ($SIZE)"
else
    echo "❌ Binary verification failed"
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ ALL DONE!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Next steps:"
echo "  1. nano .env       # Add LinkedIn test account credentials"
echo "  2. ./start.sh      # Start backend + frontend"
echo ""
echo "Or just the backend:"
echo "  ./bin/automation"
echo ""

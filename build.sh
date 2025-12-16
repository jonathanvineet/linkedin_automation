#!/bin/bash

# Complete Build Script - Handles all dependencies and edge cases
# This script will compile the Go backend with proper CGO settings for SQLite

set -e  # Exit on error

echo "🔨 LinkedIn Automation - Complete Build"
echo "========================================"
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Step 1: Check prerequisites
echo "📋 Checking prerequisites..."

if ! command -v go &> /dev/null; then
    echo -e "${RED}✗${NC} Go not found. Please install Go 1.21+"
    exit 1
fi
echo -e "${GREEN}✓${NC} Go $(go version | awk '{print $3}')"

if ! command -v gcc &> /dev/null && ! command -v clang &> /dev/null; then
    echo -e "${YELLOW}⚠${NC} Warning: No C compiler found (gcc/clang)"
    echo "   SQLite requires CGO which needs a C compiler"
    echo "   Install: sudo apt-get install build-essential"
fi

# Step 2: Create directories
echo ""
echo "📁 Creating directories..."
mkdir -p bin data logs config
echo -e "${GREEN}✓${NC} Directories created"

# Step 3: Download Go dependencies
echo ""
echo "📦 Downloading Go dependencies..."
if go mod download; then
    echo -e "${GREEN}✓${NC} Dependencies downloaded"
else
    echo -e "${RED}✗${NC} Failed to download dependencies"
    exit 1
fi

# Step 4: Tidy go.mod
echo ""
echo "🧹 Tidying go.mod..."
if go mod tidy; then
    echo -e "${GREEN}✓${NC} go.mod tidied"
else
    echo -e "${RED}✗${NC} Failed to tidy go.mod"
    exit 1
fi

# Step 5: Build with CGO enabled (required for SQLite)
echo ""
echo "🔨 Building Go backend (with CGO for SQLite)..."
echo "   Command: CGO_ENABLED=1 go build -v -o bin/automation ./cmd/app"
echo ""

export CGO_ENABLED=1

if go build -v -o bin/automation ./cmd/app; then
    echo ""
    echo -e "${GREEN}✓${NC} Backend built successfully!"
    echo "   Binary: bin/automation"
    echo "   Size: $(du -h bin/automation | cut -f1)"
else
    echo ""
    echo -e "${RED}✗${NC} Build failed!"
    echo ""
    echo "Common issues:"
    echo "  1. Missing C compiler (install: apt-get install build-essential)"
    echo "  2. Import path mismatch in go files"
    echo "  3. Missing dependencies (run: go mod tidy)"
    echo ""
    exit 1
fi

# Step 6: Verify binary
echo ""
echo "✅ Verifying binary..."
if [ -f "bin/automation" ] && [ -x "bin/automation" ]; then
    echo -e "${GREEN}✓${NC} Binary is executable"
else
    echo -e "${RED}✗${NC} Binary is not executable"
    chmod +x bin/automation
    echo -e "${GREEN}✓${NC} Made binary executable"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ Build Complete!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Next steps:"
echo "  1. Configure .env file with credentials"
echo "  2. Run: ./bin/automation"
echo "  3. Or use: ./start.sh to start both backend and frontend"
echo ""

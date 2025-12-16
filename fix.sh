#!/bin/bash

# Auto-fix script - Installs build tools and compiles the project
# Run this if setup.sh failed

echo "🔧 Auto-Fix: Installing Build Tools & Compiling"
echo "================================================"
echo ""

# Check if running with sudo capability
if [ "$EUID" -ne 0 ] && ! sudo -n true 2>/dev/null; then
    echo "⚠️  This script needs sudo to install build tools"
    echo "   It will prompt for your password if needed"
    echo ""
fi

# Install build-essential (includes gcc, make, etc.)
echo "📦 Installing build-essential..."
if sudo apt-get update -qq && sudo apt-get install -y build-essential; then
    echo "✅ Build tools installed"
else
    echo "❌ Failed to install build tools"
    echo "   Try manually: sudo apt-get install build-essential"
    exit 1
fi

echo ""
echo "🔨 Building Go backend..."
echo ""

# Make scripts executable
chmod +x *.sh 2>/dev/null

# Run the build script
if [ -f "build.sh" ]; then
    ./build.sh
else
    # Fallback: Build directly
    mkdir -p bin data logs
    export CGO_ENABLED=1
    go build -v -o bin/automation ./cmd/app
    
    if [ $? -eq 0 ]; then
        echo ""
        echo "✅ Build successful!"
    else
        echo ""
        echo "❌ Build failed. Check errors above."
        exit 1
    fi
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Fix Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Next steps:"
echo "  1. nano .env     # Add credentials"
echo "  2. ./start.sh    # Start the app"
echo ""

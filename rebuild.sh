#!/bin/bash
echo "🔨 Rebuilding after fixes..."
export CGO_ENABLED=1
go build -v -o bin/automation ./cmd/app
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ BUILD SUCCESSFUL!"
    echo ""
    ls -lh bin/automation
    echo ""
    echo "Next steps:"
    echo "  1. nano .env       # Add credentials"
    echo "  2. ./start.sh      # Start the app"
else
    echo ""
    echo "❌ Build still has errors"
    exit 1
fi

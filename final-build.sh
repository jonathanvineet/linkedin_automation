#!/bin/bash
echo "🔨 Rebuilding with latest fixes..."
echo ""
export CGO_ENABLED=1
go build -v -o bin/automation ./cmd/app
BUILD_RESULT=$?

if [ $BUILD_RESULT -eq 0 ]; then
    echo ""
    echo "✅ BUILD SUCCESSFUL!"
    echo ""
    ls -lh bin/automation
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🎉 Ready to go!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Next steps:"
    echo "  1. nano .env               # Add test credentials"
    echo "  2. ./start.sh              # Start the app"
    echo ""
    echo "Or just run the backend:"
    echo "  ./bin/automation"
    echo ""
else
    echo ""
    echo "❌ Build failed - see errors above"
    exit 1
fi

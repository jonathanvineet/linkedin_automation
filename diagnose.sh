#!/bin/bash

# Quick diagnostic script
echo "🔍 Diagnostics"
echo "=============="
echo ""

echo "Go version:"
go version
echo ""

echo "C Compiler:"
if command -v gcc &> /dev/null; then
    echo "✅ gcc found: $(gcc --version | head -n1)"
elif command -v clang &> /dev/null; then
    echo "✅ clang found: $(clang --version | head -n1)"
else
    echo "❌ No C compiler found (gcc or clang)"
    echo "   This is likely why the build failed!"
    echo "   SQLite3 requires CGO which needs a C compiler"
    echo ""
    echo "   Fix: Run ./fix.sh to install build tools"
fi
echo ""

echo "Go environment:"
echo "  GOPATH: $(go env GOPATH)"
echo "  GOROOT: $(go env GOROOT)"
echo "  CGO_ENABLED: $(go env CGO_ENABLED)"
echo ""

echo "Project structure:"
ls -la cmd/app/*.go 2>/dev/null && echo "  ✅ cmd/app/main.go exists" || echo "  ❌ cmd/app/main.go missing"
ls -d internal/*/ 2>/dev/null && echo "  ✅ internal/ packages exist" || echo "  ❌ internal/ packages missing"
echo ""

echo "Dependencies:"
if [ -f "go.mod" ]; then
    echo "  ✅ go.mod exists"
    echo "  Module: $(head -n1 go.mod)"
else
    echo "  ❌ go.mod missing"
fi

if [ -f "go.sum" ]; then
    echo "  ✅ go.sum exists"
else
    echo "  ⚠️  go.sum missing (run: go mod tidy)"
fi
echo ""

echo "Attempt simple build test:"
mkdir -p /tmp/test-build
CGO_ENABLED=1 go build -o /tmp/test-build/test ./cmd/app 2>&1 | head -n 20
BUILD_STATUS=$?

echo ""
if [ $BUILD_STATUS -eq 0 ]; then
    echo "✅ Build test successful!"
    echo "   Binary created at: /tmp/test-build/test"
    rm -f /tmp/test-build/test
else
    echo "❌ Build test failed (see errors above)"
    echo ""
    echo "Common fixes:"
    echo "  1. Missing C compiler: ./fix.sh"
    echo "  2. Missing dependencies: go mod tidy"
    echo "  3. Import path issues: check go.mod module name"
fi

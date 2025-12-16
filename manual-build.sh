#!/bin/bash
# Manual build script to see errors
echo "=== Building Go Backend ==="
cd /workspaces/linkedin_automation
go build -v -o bin/automation ./cmd/app
echo ""
echo "Build exit code: $?"

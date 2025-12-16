#!/bin/bash
# Direct build command
export CGO_ENABLED=1
go build -o bin/automation ./cmd/app && echo "✅ SUCCESS!" || echo "❌ FAILED"

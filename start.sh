#!/bin/bash

# LinkedIn Automation - Startup Script
# ⚠️  EDUCATIONAL USE ONLY - DO NOT USE IN PRODUCTION

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🤖 LinkedIn Automation Proof-of-Concept"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "⚠️  WARNING: This violates LinkedIn's Terms of Service"
echo "    Use ONLY for educational purposes on test accounts!"
echo ""

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "❌ Error: .env file not found"
    echo "   Please copy .env.example to .env and configure it:"
    echo "   cp .env.example .env"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo "   Please install Go 1.21+ from https://golang.org/dl/"
    exit 1
fi

# Check if Node is installed
if ! command -v node &> /dev/null; then
    echo "❌ Error: Node.js is not installed"
    echo "   Please install Node.js 18+ from https://nodejs.org/"
    exit 1
fi

echo "✅ Prerequisites check passed"
echo ""

# Build Go backend
echo "📦 Building Go backend..."
mkdir -p bin
go mod download
CGO_ENABLED=1 go build -o bin/automation ./cmd/app
echo "✅ Backend built successfully"
echo ""

# Install Node dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    npm install
    echo "✅ Dependencies installed"
    echo ""
fi

# Create necessary directories
mkdir -p data logs

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 Starting Services"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Backend API will start on: http://localhost:8090"
echo "Frontend UI will start on: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Start backend in background
echo "▶️  Starting Go API server..."
./bin/automation &
BACKEND_PID=$!

# Wait for backend to start
sleep 3

# Start frontend
echo "▶️  Starting React frontend..."
npm run dev &
FRONTEND_PID=$!

# Cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Shutting down services..."
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    echo "✅ Services stopped"
    exit 0
}

trap cleanup INT TERM

# Wait for processes
wait

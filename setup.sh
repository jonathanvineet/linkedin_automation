#!/bin/bash

# LinkedIn Automation PoC - Quick Setup Script
# This script automates the initial setup process

set -e  # Exit on any error

echo "🚀 LinkedIn Automation PoC - Quick Setup"
echo "========================================"
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check Go installation
echo -n "Checking Go installation... "
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓${NC} Found $GO_VERSION"
else
    echo -e "${RED}✗${NC} Go not found!"
    echo "Please install Go 1.21+ from https://golang.org/dl/"
    exit 1
fi

# Check Node.js installation
echo -n "Checking Node.js installation... "
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    echo -e "${GREEN}✓${NC} Found $NODE_VERSION"
else
    echo -e "${RED}✗${NC} Node.js not found!"
    echo "Please install Node.js 18+ from https://nodejs.org/"
    exit 1
fi

# Check Chrome/Chromium
echo -n "Checking Chrome/Chromium... "
if command -v google-chrome &> /dev/null || command -v chromium &> /dev/null || command -v chromium-browser &> /dev/null; then
    echo -e "${GREEN}✓${NC} Found"
else
    echo -e "${YELLOW}⚠${NC} Chrome/Chromium not found"
    echo "  Rod will attempt to download a suitable browser"
fi

echo ""
echo "📦 Installing Dependencies"
echo "=========================="

# Install Go dependencies
echo -n "Running go mod download... "
if go mod download 2>/dev/null; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC} Failed"
    exit 1
fi

echo -n "Running go mod tidy... "
if go mod tidy 2>/dev/null; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC} Failed"
    exit 1
fi

# Install Node dependencies
echo -n "Running npm install... "
if npm install --silent 2>/dev/null; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC} Failed"
    exit 1
fi

echo ""
echo "🔨 Building Application"
echo "======================="

# Build backend
echo -n "Building Go backend... "
# SQLite requires CGO
export CGO_ENABLED=1
if go build -o bin/automation ./cmd/app; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC} Failed"
    echo "Error: Go build failed. Run 'go build -o bin/automation ./cmd/app' manually to see errors"
    exit 1
fi

# Build frontend
echo -n "Building React frontend... "
if npm run build --silent 2>/dev/null; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC} Failed"
    exit 1
fi

echo ""
echo "📁 Creating Directories"
echo "======================="

# Create necessary directories
mkdir -p data logs bin

echo -e "${GREEN}✓${NC} Created: data/, logs/, bin/"

echo ""
echo "🔐 Environment Configuration"
echo "============================"

# Check if .env exists
if [ -f ".env" ]; then
    echo -e "${GREEN}✓${NC} .env file exists"
    
    # Check if credentials are set
    if grep -q "LINKEDIN_EMAIL=.*@.*" .env && grep -q "LINKEDIN_PASSWORD=..*" .env; then
        echo -e "${GREEN}✓${NC} Credentials configured"
    else
        echo -e "${YELLOW}⚠${NC} Credentials not set in .env"
        echo "  Please edit .env and add your test account credentials"
    fi
else
    echo -e "${YELLOW}⚠${NC} .env file not found"
    
    if [ -f ".env.example" ]; then
        echo -n "Creating .env from .env.example... "
        cp .env.example .env
        echo -e "${GREEN}✓${NC}"
        echo -e "${YELLOW}⚠${NC} Please edit .env and add your credentials"
    else
        echo -e "${RED}✗${NC} .env.example not found"
        echo "  Please create .env manually"
    fi
fi

echo ""
echo "🔧 Making Scripts Executable"
echo "============================="

chmod +x start.sh health-check.sh 2>/dev/null && echo -e "${GREEN}✓${NC} Scripts are executable" || echo -e "${YELLOW}⚠${NC} Could not set execute permissions (may need sudo)"

echo ""
echo "✅ Setup Complete!"
echo "=================="
echo ""
echo -e "${BLUE}Next Steps:${NC}"
echo "1. Configure .env with your LinkedIn test account credentials"
echo "2. Review config/config.yaml settings"
echo "3. Run: ./start.sh"
echo "4. Open: http://localhost:8080"
echo ""
echo -e "${YELLOW}⚠️  IMPORTANT:${NC}"
echo "• Use ONLY test accounts, never production credentials"
echo "• This is for educational purposes only"
echo "• Respect LinkedIn's Terms of Service"
echo "• Start with low rate limits (≤20 connections/day)"
echo ""
echo -e "${GREEN}Happy learning! 🎓${NC}"

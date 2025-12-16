# 🤖 LinkedIn Automation Proof-of-Concept

[![Educational Purpose](https://img.shields.io/badge/Purpose-Educational-yellow.svg)](https://github.com/jonathanvineet/linkedin-automation)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> ⚠️ **CRITICAL DISCLAIMER**: This project is a **technical proof-of-concept** for educational and research purposes only. Using automation tools on LinkedIn **violates their Terms of Service** and may result in account suspension or legal action. **DO NOT use this in production or on real accounts.**

## 📋 Table of Contents

- [Overview](#overview)
- [⚠️ Ethical & Legal Warning](#ethical--legal-warning)
- [Architecture](#architecture)
- [Features](#features)
- [Stealth Techniques](#stealth-techniques)
- [Installation](#installation)
- [Usage](#usage)

## 📖 Overview

This project demonstrates advanced browser automation techniques, anti-detection strategies, and realistic human behavior modeling using **Go** and **React/TypeScript**.

### Technologies

- **Backend**: Go 1.21+, Rod (Chrome DevTools Protocol), SQLite
- **Frontend**: React 18, TypeScript, Tailwind CSS, shadcn/ui
- **Automation**: Rod browser automation library
- **Stealth**: Custom human behavior simulation algorithms

## ⚠️ Ethical & Legal Warning

**This tool violates LinkedIn's Terms of Service.** Use only for educational purposes on test accounts. Never use on production accounts.

## 🚀 Quick Start

### Option 1: Automated Setup (Recommended)
```bash
chmod +x *.sh
./setup.sh
```

**If setup.sh fails at "Building Go backend":**
```bash
./fix.sh    # Installs build tools and compiles
```

**Or diagnose the issue:**
```bash
./diagnose.sh    # Shows what's wrong
```

See [BUILD_FIX.md](BUILD_FIX.md) for detailed troubleshooting.

### Option 2: Manual Installation

#### Prerequisites
- Go 1.21+ ([Download](https://golang.org/dl/))
- Node.js 18+ ([Download](https://nodejs.org/))
- Chrome/Chromium browser
- **C compiler** (gcc/clang) - Required for SQLite

**Install build tools (Ubuntu/Debian/Codespaces):**
```bash
sudo apt-get update
sudo apt-get install -y build-essential
```

#### Step 1: Clone Repository
```bash
git clone https://github.com/jonathanvineet/linkedin-automation.git
cd linkedin-automation
```

#### Step 2: Backend Setup
```bash
# Install Go dependencies
go mod download

# Build the application (CGO required for SQLite)
mkdir -p bin
CGO_ENABLED=1 go build -o bin/automation ./cmd/app
```

### Step 3: Frontend Setup
```bash
# Install Node dependencies
npm install

# Build frontend
npm run build
```

### Step 4: Configure Environment
```bash
cp .env.example .env
# Edit .env with your LinkedIn credentials (test account only!)
```

## 🎮 Usage

**Terminal 1** - Start Go backend:
```bash
./bin/automation
```

**Terminal 2** - Start React frontend:
```bash
npm run dev
```

Access the dashboard at `http://localhost:8080`

## ✨ Features

- ✅ Human-like mouse movement (Bézier curves)
- ✅ Realistic typing with errors
- ✅ Browser fingerprint masking
- ✅ Context-aware timing & delays
- ✅ Persona-based behavior patterns
- ✅ Activity scheduling & rate limiting
- ✅ SQLite state persistence

## 📁 Project Structure

```
linkedin-automation/
├── cmd/app/main.go              # Go API server
├── internal/                    # Go modules
│   ├── browser/                 # Browser automation
│   ├── stealth/                 # Stealth techniques
│   ├── behavior/                # Persona system
│   └── ...
├── src/                         # React frontend
├── config/config.yaml           # Configuration
└── README.md
```

## 📜 License

MIT License - Educational use only. See [LICENSE](LICENSE) for details.

---

*⚠️ Remember: This is a proof-of-concept for learning. Do not use on real LinkedIn accounts.*

# 🚨 BUILD FAILED? FIX IT NOW!

The setup script failed because SQLite requires a C compiler (gcc).

## ⚡ QUICK FIX (Choose One)

### Option 1: Automatic Fix (Easiest)
```bash
chmod +x quickfix.sh
./quickfix.sh
```
This installs build tools and compiles everything for you.

### Option 2: Manual Fix
```bash
# Install build tools
sudo apt-get update
sudo apt-get install -y build-essential

# Build manually
mkdir -p bin
export CGO_ENABLED=1
go build -o bin/automation ./cmd/app
```

### Option 3: Detailed Diagnostics
```bash
chmod +x diagnose.sh
./diagnose.sh
```
Shows exactly what's wrong and how to fix it.

## After Fixing

1. **Configure credentials:**
   ```bash
   nano .env
   # Add your TEST LinkedIn account
   ```

2. **Start the app:**
   ```bash
   ./start.sh
   ```

## Need More Help?

- [BUILD_FIX.md](BUILD_FIX.md) - Detailed troubleshooting guide
- [QUICK_START.md](QUICK_START.md) - Full quick start guide
- [GETTING_STARTED.md](GETTING_STARTED.md) - Complete setup guide

---

**TL;DR: Run `./quickfix.sh` and you're done!** ✅

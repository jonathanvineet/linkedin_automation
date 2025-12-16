# 🔧 BUILD TROUBLESHOOTING

## Quick Fix

If setup.sh failed at "Building Go backend", run:

```bash
chmod +x build.sh
./build.sh
```

This will show you the actual error and guide you through fixes.

## Common Issues

### 1. Missing C Compiler (Most Likely Issue)

**Error:** "gcc: command not found" or similar

**Why:** SQLite3 requires CGO which needs a C compiler

**Fix:**
```bash
# On Ubuntu/Debian/Codespaces
sudo apt-get update
sudo apt-get install -y build-essential

# On macOS
xcode-select --install

# On Alpine Linux
apk add gcc musl-dev

# Then retry:
./build.sh
```

### 2. CGO_ENABLED Not Set

**Error:** Build works but binary crashes with "cannot find -lsqlite3"

**Fix:**
```bash
export CGO_ENABLED=1
go build -o bin/automation ./cmd/app
```

### 3. Import Path Mismatch

**Error:** "package github.com/jonathanvineet/linkedin-automation/internal/... is not in GOROOT"

**Fix:** Ensure go.mod has the correct module name:
```bash
head -1 go.mod
# Should show: module github.com/jonathanvineet/linkedin-automation

# If different, update all imports or fix go.mod
```

### 4. Missing Dependencies

**Error:** "cannot find package"

**Fix:**
```bash
go mod download
go mod tidy
go build -o bin/automation ./cmd/app
```

### 5. Permission Denied

**Error:** "permission denied" when running binary

**Fix:**
```bash
chmod +x bin/automation
```

## Manual Build Steps

If automated scripts fail, build manually:

```bash
# 1. Create directories
mkdir -p bin data logs

# 2. Install C compiler (if needed)
sudo apt-get install -y build-essential

# 3. Download dependencies
go mod download
go mod tidy

# 4. Build with CGO
export CGO_ENABLED=1
go build -v -o bin/automation ./cmd/app

# 5. Verify
./bin/automation --help  # Should run without errors
```

## Testing the Build

```bash
# Check if binary exists
ls -lh bin/automation

# Check if it's executable
file bin/automation

# Try running it (will fail without .env but should not crash)
./bin/automation
# Should show: "No .env file found" or start successfully
```

## Platform-Specific Notes

### GitHub Codespaces / Devcontainer
```bash
# Install build tools
sudo apt-get update && sudo apt-get install -y build-essential

# Then build
./build.sh
```

### WSL2 / Windows
```bash
# Ensure you have gcc
sudo apt install build-essential

# Build
./build.sh
```

### macOS
```bash
# Install Xcode command line tools
xcode-select --install

# Build
./build.sh
```

## Verification

After successful build:

```bash
# Should see:
$ ls -lh bin/
total 45M
-rwxr-xr-x 1 user user 45M Dec 16 10:30 automation

# Test run:
$ ./bin/automation
INFO[0000] 🚀 LinkedIn Automation PoC starting...
# (will wait for proper .env config)
```

## Still Having Issues?

1. **Check Go version:**
   ```bash
   go version  # Should be 1.21 or higher
   ```

2. **Check detailed build output:**
   ```bash
   CGO_ENABLED=1 go build -v -x -o bin/automation ./cmd/app
   ```

3. **Check C compiler:**
   ```bash
   which gcc || which clang
   gcc --version || clang --version
   ```

4. **Verify file structure:**
   ```bash
   ls -R cmd/ internal/
   # Should show all Go source files
   ```

5. **Check for syntax errors:**
   ```bash
   go vet ./...
   ```

## Success Indicators

✅ `bin/automation` binary exists and is ~40-50MB  
✅ Binary is executable (green in `ls -l`)  
✅ Running `./bin/automation` attempts to start (even if fails due to missing .env)  
✅ No "cannot find package" errors  
✅ No "gcc not found" errors  

---

**Once built successfully, proceed with:**
```bash
# 1. Configure credentials
nano .env

# 2. Start the application
./start.sh
```

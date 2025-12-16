# 🚀 Quick Reference Guide

One-page reference for getting started quickly with the LinkedIn Automation PoC.

## 📦 Installation (One Command)

```bash
./setup.sh
```

This will:
- ✅ Check system requirements (Go, Node.js, Chrome)
- ✅ Install all dependencies
- ✅ Build backend and frontend
- ✅ Create necessary directories
- ✅ Set up configuration files

## ⚙️ Configuration (2 Minutes)

1. **Edit `.env` file:**
```bash
nano .env
```

2. **Add credentials (TEST ACCOUNT ONLY!):**
```env
LINKEDIN_EMAIL=your-test-account@email.com
LINKEDIN_PASSWORD=your-test-password
API_PORT=8090
HEADLESS=false
```

3. **Adjust rate limits in `config/config.yaml`:**
```yaml
rate_limits:
  connections_per_day: 20    # Conservative limit
  messages_per_day: 10       # Safe limit
```

## 🏃 Running (Choose One Method)

### Method 1: Use Start Script (Recommended)
```bash
./start.sh
```

### Method 2: Use Makefile
```bash
make run
```

### Method 3: Manual Start
```bash
# Terminal 1 - Backend
./bin/automation

# Terminal 2 - Frontend
npm run dev
```

## 🌐 Access the Application

Open in your browser:
```
http://localhost:8080
```

API endpoint:
```
http://localhost:8090/api
```

## 🧪 Testing

### Quick System Test
```bash
./test.sh
```

### Health Check
```bash
./health-check.sh
```

### Manual API Test
```bash
curl http://localhost:8090/api/status
```

## 🎯 Common Tasks

### Start Automation
1. Open UI at `http://localhost:8080`
2. Click **Start** button
3. Watch browser authenticate
4. Monitor activity log

### Perform Search
```bash
curl -X POST http://localhost:8090/api/search \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": "Software Engineer",
    "location": "San Francisco",
    "max_results": 10
  }'
```

### Send Connection Request
```bash
curl -X POST http://localhost:8090/api/connect \
  -H "Content-Type: application/json" \
  -d '{
    "profile_url": "https://www.linkedin.com/in/testprofile",
    "note": "Hello! I would like to connect."
  }'
```

### Change Persona
```bash
curl -X POST http://localhost:8090/api/persona \
  -H "Content-Type: application/json" \
  -d '{"persona_type": "founder"}'
```

Available personas: `recruiter`, `founder`, `sales`

### Check Stats
```bash
curl http://localhost:8090/api/stats
```

### View Activity Log
```bash
curl http://localhost:8090/api/activity
```

## 🛑 Stopping

### From UI
Click the **Stop** button in the dashboard

### From Terminal
```bash
# Find the process
ps aux | grep automation

# Kill it
pkill -f automation

# Or use Ctrl+C in the terminal running it
```

## 📊 Monitoring

### View Logs
```bash
# Real-time logs
tail -f logs/app.log

# Pretty print JSON logs
tail -f logs/app.log | jq '.'

# Filter by level
cat logs/app.log | jq 'select(.level=="ERROR")'
```

### Check Database
```bash
# View all tables
sqlite3 data/automation.db ".tables"

# Count connections sent
sqlite3 data/automation.db "SELECT COUNT(*) FROM connection_requests;"

# View recent activity
sqlite3 data/automation.db "SELECT * FROM activity_logs ORDER BY timestamp DESC LIMIT 10;"
```

## 🔧 Troubleshooting

### Backend Won't Start
```bash
# Check if port is in use
lsof -i :8090

# View detailed logs
./bin/automation 2>&1 | tee debug.log
```

### Frontend Won't Start
```bash
# Check if port is in use
lsof -i :8080

# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

### API Not Responding
```bash
# Test connectivity
curl -v http://localhost:8090/api/status

# Check backend is running
ps aux | grep automation

# Review logs
tail -n 50 logs/app.log
```

### Authentication Fails
1. Verify credentials in `.env` are correct
2. Check if LinkedIn requires CAPTCHA (use `HEADLESS=false` to see)
3. Try with a different test account
4. Review logs for specific error messages

### Rate Limit Hit
```bash
# Check remaining quota
curl http://localhost:8090/api/stats

# Wait for cooldown or increase limits in config.yaml
```

## 🧹 Cleanup

### Remove Data (Keep Code)
```bash
make clean
```

### Complete Cleanup
```bash
rm -rf bin/ dist/ data/ logs/ node_modules/
rm -f *.log go.sum package-lock.json
```

## 📚 Documentation Files

- **README.md** - Main project overview
- **GETTING_STARTED.md** - Detailed setup guide
- **ARCHITECTURE.md** - System design and structure
- **VERIFICATION.md** - Complete testing checklist
- **CONTRIBUTING.md** - How to contribute
- **PROJECT_SUMMARY.md** - Quick project overview
- **QUICK_START.md** - This file!

## 🔑 Key Files

| File | Purpose |
|------|---------|
| `.env` | Credentials and configuration |
| `config/config.yaml` | Rate limits and behavior settings |
| `cmd/app/main.go` | Backend entry point |
| `src/pages/Dashboard.tsx` | Frontend main page |
| `data/automation.db` | SQLite database |
| `logs/app.log` | Application logs |

## ⚡ Quick Commands Cheat Sheet

```bash
# Setup
./setup.sh                  # One-time setup

# Build
make build                  # Build everything
go build -o bin/automation ./cmd/app  # Backend only
npm run build              # Frontend only

# Run
./start.sh                 # Start both services
make dev                   # Development mode
./bin/automation           # Backend only
npm run dev                # Frontend only

# Test
./test.sh                  # Full system test
./health-check.sh          # Quick health check
make test                  # Go unit tests

# Monitor
tail -f logs/app.log       # View logs
./health-check.sh          # Check status

# Clean
make clean                 # Remove build artifacts
rm -rf data/ logs/         # Remove data
```

## 🎓 Learning Path

1. **Start Here:** Read GETTING_STARTED.md
2. **Understand System:** Review ARCHITECTURE.md
3. **Run It:** Follow this QUICK_START.md
4. **Explore Code:** Start with `cmd/app/main.go`
5. **Test Features:** Use VERIFICATION.md checklist
6. **Contribute:** Read CONTRIBUTING.md

## ⚠️ Important Reminders

- ✅ **Use TEST accounts only** - Never use real credentials
- ✅ **Respect rate limits** - Start with ≤20 connections/day
- ✅ **Educational purpose** - This is for learning, not production
- ✅ **LinkedIn ToS** - Understand this violates LinkedIn's Terms
- ✅ **Business hours** - Configure realistic activity times
- ✅ **Monitor activity** - Watch logs and stats regularly

## 🆘 Getting Help

1. Check logs: `tail -f logs/app.log`
2. Run health check: `./health-check.sh`
3. Review troubleshooting section above
4. Read GETTING_STARTED.md for details
5. Check VERIFICATION.md for testing steps

## 🎯 Success Criteria

You're ready when:
- ✅ Setup script completes without errors
- ✅ Backend starts and responds to API calls
- ✅ Frontend loads at http://localhost:8080
- ✅ Health check passes all tests
- ✅ Can start/stop automation from UI
- ✅ Activity log shows actions
- ✅ Database records data

---

**Ready to start? Run: `./setup.sh`** 🚀

**Questions? Check: GETTING_STARTED.md** 📖

**Need details? Read: ARCHITECTURE.md** 🏗️

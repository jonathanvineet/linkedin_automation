# ✅ Verification Checklist

Use this checklist to verify your LinkedIn Automation setup is working correctly.

## 📋 Pre-Launch Checklist

### Environment Setup
- [ ] Go 1.21+ installed (`go version`)
- [ ] Node.js 18+ installed (`node --version`)
- [ ] Chrome/Chromium installed
- [ ] `.env` file created and configured
- [ ] LinkedIn credentials added (TEST ACCOUNT ONLY!)

### Dependencies
- [ ] Go modules downloaded (`go mod download`)
- [ ] Node packages installed (`npm install`)
- [ ] No dependency errors

### Configuration
- [ ] `config/config.yaml` exists
- [ ] `.env` has correct API_PORT (8090)
- [ ] `.env` has valid credentials
- [ ] `.env` has reasonable rate limits (≤20 connections/day)

## 🔨 Build Verification

### Backend Build
```bash
go build -o bin/automation ./cmd/app
```
- [ ] Build completes without errors
- [ ] `bin/automation` binary exists
- [ ] Binary is executable

### Frontend Build
```bash
npm run build
```
- [ ] Build completes without errors
- [ ] `dist/` directory created
- [ ] No TypeScript errors

## 🚀 Runtime Verification

### Start Backend
```bash
./bin/automation
```

**Expected Output:**
```
INFO[0000] 🚀 LinkedIn Automation PoC starting...
WARN[0000] ⚠️  EDUCATIONAL USE ONLY - DO NOT USE IN PRODUCTION
INFO[0001] Logger initialized successfully
INFO[0002] API server starting port=8090
```

Checklist:
- [ ] Server starts without errors
- [ ] Listens on port 8090
- [ ] Logs directory created
- [ ] Database directory created

### Start Frontend
```bash
npm run dev
```

**Expected Output:**
```
VITE v5.x.x ready in xxx ms

➜  Local:   http://localhost:8080/
➜  Network: use --host to expose
```

Checklist:
- [ ] Dev server starts
- [ ] Port 8080 is available
- [ ] No compilation errors

## 🌐 API Health Checks

### Status Endpoint
```bash
curl http://localhost:8090/api/status
```

**Expected Response:**
```json
{
  "running": false,
  "logged_in": false,
  "persona": "none",
  "stealth": true,
  "timestamp": "2024-12-16T10:30:00Z"
}
```

- [ ] Returns 200 OK
- [ ] JSON is valid
- [ ] Fields are present

### Stats Endpoint
```bash
curl http://localhost:8090/api/stats
```

**Expected Response:**
```json
{
  "connections_sent": 0,
  "messages_sent": 0,
  "cooldown_seconds": 0,
  "daily_limit": {
    "connections": 20,
    "messages": 10
  }
}
```

- [ ] Returns 200 OK
- [ ] Stats are zero initially
- [ ] Limits match config

### Activity Endpoint
```bash
curl http://localhost:8090/api/activity
```

**Expected Response:**
```json
[]
```

- [ ] Returns 200 OK or empty array
- [ ] JSON array format

## 🖥️ Frontend UI Checks

### Open Dashboard
Navigate to: `http://localhost:8080`

**Status Cards:**
- [ ] Browser Session shows "Idle"
- [ ] Logged In shows "No"
- [ ] Persona Active shows persona name
- [ ] Stealth Mode shows "Enabled"

**Daily Counters:**
- [ ] Connections Sent: 0
- [ ] Messages Sent: 0
- [ ] Cooldown Time: "Ready"

**Persona Panel:**
- [ ] Three persona buttons visible (Recruiter/Founder/Sales)
- [ ] Sliders work for typing/precision/error
- [ ] Behavior summary displays

**Stealth Techniques:**
- [ ] 4+ techniques listed
- [ ] Switches show enabled/disabled
- [ ] Icons display correctly

**Activity Log:**
- [ ] Empty or shows initialization
- [ ] Scrollable area
- [ ] Timestamp format correct

**Control Buttons:**
- [ ] Start button visible
- [ ] Stop button disabled initially
- [ ] No JavaScript errors in console

## 🔐 Authentication Test

### Start Automation
Click **Start** button in UI

**Expected Behavior:**
- [ ] Button shows loading state
- [ ] Browser window opens (if headless=false)
- [ ] Navigates to LinkedIn login
- [ ] Credentials entered automatically
- [ ] Login button clicked
- [ ] Status changes to "Active"
- [ ] Activity log updates

**Check Logs:**
```bash
tail -f logs/app.log
```
- [ ] Login process logged
- [ ] No critical errors
- [ ] Session initialized

**Verify Status:**
- [ ] Browser Session: "Active"
- [ ] Logged In: "Yes"
- [ ] System status indicator green

## 🔍 Search Test

### API Search
```bash
curl -X POST http://localhost:8090/api/search \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": "Software Engineer",
    "location": "San Francisco",
    "max_results": 5
  }'
```

**Expected:**
- [ ] Returns 200 OK or results
- [ ] Results array with profiles
- [ ] Profile URLs present
- [ ] Activity log shows search action

## 🔗 Connection Test (Optional)

⚠️ **Only test with dummy profiles or test accounts!**

```bash
curl -X POST http://localhost:8090/api/connect \
  -H "Content-Type: application/json" \
  -d '{
    "profile_url": "https://www.linkedin.com/in/testprofile",
    "note": "This is a test connection request."
  }'
```

**Expected:**
- [ ] Browser navigates to profile
- [ ] Human-like delays observed
- [ ] Connect button clicked
- [ ] Note added (if available)
- [ ] Send button clicked
- [ ] Activity logged
- [ ] Connection count increases

## 🛑 Stop Test

### Stop Automation
Click **Stop** button in UI

**Expected:**
- [ ] Browser closes gracefully
- [ ] Status changes to "Idle"
- [ ] Logged In becomes "No"
- [ ] Activity logged
- [ ] No errors in logs

## 💾 Data Persistence

### Database Check
```bash
sqlite3 data/automation.db "SELECT name FROM sqlite_master WHERE type='table';"
```

**Expected Tables:**
- [ ] connection_requests
- [ ] messages
- [ ] activity_logs
- [ ] session_data

### Query Data
```bash
# Connections
sqlite3 data/automation.db "SELECT COUNT(*) FROM connection_requests;"

# Messages
sqlite3 data/automation.db "SELECT COUNT(*) FROM messages;"

# Activity
sqlite3 data/automation.db "SELECT COUNT(*) FROM activity_logs;"
```

- [ ] Queries execute successfully
- [ ] Counts match expectations
- [ ] No database errors

## 📝 Logging Verification

### Log File
```bash
ls -lh logs/app.log
```
- [ ] File exists
- [ ] Growing in size
- [ ] Readable permissions

### Log Format
```bash
cat logs/app.log | jq '.'
```
- [ ] Valid JSON format
- [ ] Timestamps present
- [ ] Log levels correct
- [ ] Structured fields

## 🏥 Health Check Script

### Run Health Check
```bash
./health-check.sh
```

**Expected:**
- [ ] API connectivity confirmed
- [ ] Status displayed
- [ ] Stats shown
- [ ] Activity logs listed
- [ ] Database info shown
- [ ] Recent errors checked
- [ ] No critical issues

## 🔄 Rate Limiting Test

### Send Multiple Connections
(Use test profiles only!)

```bash
for i in {1..3}; do
  curl -X POST http://localhost:8090/api/connect \
    -H "Content-Type: application/json" \
    -d '{"profile_url": "https://linkedin.com/in/test'$i'", "note": "Test"}'
  sleep 2
done
```

**Expected:**
- [ ] First request succeeds
- [ ] Subsequent requests delayed
- [ ] Cooldown enforced
- [ ] Rate limit respected
- [ ] Daily limit tracked

### Check Stats After
```bash
curl http://localhost:8090/api/stats
```
- [ ] Connection count increased
- [ ] Cooldown timer active
- [ ] Daily limit decremented

## 🎭 Persona Switch Test

### Change Persona
```bash
curl -X POST http://localhost:8090/api/persona \
  -H "Content-Type: application/json" \
  -d '{"persona_type": "founder"}'
```

**Expected:**
- [ ] Returns success
- [ ] Status endpoint shows new persona
- [ ] UI updates persona name
- [ ] Activity logged

## 🚨 Error Handling

### Invalid API Calls

**Missing Credentials:**
```bash
# Stop and clear env
unset LINKEDIN_EMAIL
./bin/automation
```
- [ ] Fails gracefully
- [ ] Clear error message
- [ ] No crash

**Invalid Profile URL:**
```bash
curl -X POST http://localhost:8090/api/connect \
  -H "Content-Type: application/json" \
  -d '{"profile_url": "invalid"}'
```
- [ ] Returns error response
- [ ] Appropriate error message
- [ ] No server crash

## 📊 Performance Check

### Response Times
- [ ] `/api/status` < 50ms
- [ ] `/api/stats` < 100ms
- [ ] `/api/activity` < 200ms

### Resource Usage
```bash
# While running
ps aux | grep automation
```
- [ ] Memory < 500 MB
- [ ] CPU < 20%
- [ ] No memory leaks

## 🎯 Final Verification

### Complete System Test
1. [ ] Start both services
2. [ ] Open UI in browser
3. [ ] Click Start
4. [ ] Verify login succeeds
5. [ ] Perform search
6. [ ] Send test connection (optional)
7. [ ] Check activity log updates
8. [ ] Verify stats increase
9. [ ] Click Stop
10. [ ] Verify graceful shutdown

### Documentation Review
- [ ] README.md is clear
- [ ] GETTING_STARTED.md is accurate
- [ ] ARCHITECTURE.md explains design
- [ ] All code comments present

### Safety Checks
- [ ] Warnings displayed prominently
- [ ] Test account credentials only
- [ ] Rate limits set conservatively
- [ ] Business hours enforced
- [ ] No production use enabled

---

## ✅ All Clear!

If all items are checked, your LinkedIn Automation PoC is:
- ✅ Properly installed
- ✅ Correctly configured
- ✅ Fully functional
- ✅ Ready for educational use

---

## ⚠️ If Issues Found

1. **Check Logs**: `tail -f logs/app.log`
2. **Run Health Check**: `./health-check.sh`
3. **Verify Environment**: Check `.env` settings
4. **Review Documentation**: GETTING_STARTED.md
5. **Open Issue**: GitHub issues if needed

---

**Remember: This is for learning only. Use responsibly!** 🎓

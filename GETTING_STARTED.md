# Getting Started Guide

## 🚀 Quick Start

### 1. Clone and Setup
```bash
git clone https://github.com/jonathanvineet/linkedin-automation.git
cd linkedin-automation
```

### 2. Configure Environment
```bash
# Copy environment template
cp .env.example .env

# Edit with your settings (TEST ACCOUNT ONLY!)
nano .env
```

**Required Configuration:**
```env
LINKEDIN_EMAIL=test@example.com      # Test account only!
LINKEDIN_PASSWORD=test_password
API_PORT=8090
LOG_LEVEL=info
```

### 3. Install Dependencies
```bash
# Go dependencies
go mod download

# Node dependencies
npm install
```

### 4. Start Application

**Option A: Using startup script**
```bash
./start.sh
```

**Option B: Manual start (two terminals)**

Terminal 1 - Backend:
```bash
go build -o bin/automation ./cmd/app
./bin/automation
```

Terminal 2 - Frontend:
```bash
npm run dev
```

### 5. Access Dashboard

Open your browser to: **http://localhost:8080**

Click **"Start"** button to begin automation

---

## 📝 Basic Usage

### Dashboard Controls

1. **Start Button**: Initializes browser session and logs into LinkedIn
2. **Stop Button**: Closes browser and stops automation
3. **Persona Selector**: Choose behavior profile (Recruiter/Founder/Sales)
4. **Activity Log**: Real-time action tracking

### API Endpoints

**Get Status**
```bash
curl http://localhost:8090/api/status
```

**Start Automation**
```bash
curl -X POST http://localhost:8090/api/start
```

**Search People**
```bash
curl -X POST http://localhost:8090/api/search \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": "Software Engineer",
    "location": "San Francisco",
    "max_results": 10
  }'
```

**Send Connection Request**
```bash
curl -X POST http://localhost:8090/api/connect \
  -H "Content-Type: application/json" \
  -d '{
    "profile_url": "https://linkedin.com/in/johndoe",
    "note": "Hi! Would love to connect."
  }'
```

---

## 🔧 Configuration Options

### config/config.yaml

**Rate Limits** (Keep these LOW for safety):
```yaml
rate_limits:
  daily_connections: 20     # Max connections per day
  daily_messages: 10        # Max messages per day
  min_action_delay_seconds: 30
  max_action_delay_seconds: 180
```

**Business Hours**:
```yaml
business_hours:
  enabled: true
  start_hour: 9             # 9 AM
  end_hour: 17              # 5 PM
  days: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
```

**Stealth Techniques**:
```yaml
stealth:
  enabled: true
  techniques:
    mouse_curves: true
    typing_variance: true
    scroll_randomization: true
    fingerprint_masking: true
```

---

## 🐛 Troubleshooting

### Backend won't start
```bash
# Check if port is already in use
lsof -i :8090

# Try different port
API_PORT=8091 ./bin/automation
```

### Login fails
- Verify credentials in `.env`
- Check if LinkedIn requires 2FA (complete manually in browser)
- Reduce `HEADLESS_MODE=false` to see what's happening

### Browser doesn't launch
```bash
# Install Chrome/Chromium
# Ubuntu/Debian:
sudo apt-get install chromium-browser

# macOS:
brew install --cask google-chrome
```

### Rate limit errors
```bash
# Check current stats
curl http://localhost:8090/api/stats

# Adjust limits in .env
DAILY_CONNECTION_LIMIT=5
DAILY_MESSAGE_LIMIT=3
```

---

## 📊 Monitoring

**View Logs**:
```bash
# Application logs
tail -f logs/app.log

# Structured JSON output
cat logs/app.log | jq '.'
```

**Database Inspection**:
```bash
# View sent connections
sqlite3 data/automation.db "SELECT * FROM connection_requests;"

# View activity log
sqlite3 data/automation.db "SELECT * FROM activity_logs ORDER BY timestamp DESC LIMIT 10;"
```

---

## ⚠️  Safety Tips

1. **Use test accounts ONLY** - Never use your real LinkedIn account
2. **Keep limits low** - Start with 5-10 connections per day
3. **Monitor closely** - Watch for unusual behavior
4. **Business hours only** - Don't automate outside 9-5
5. **No production use** - This violates LinkedIn ToS

---

## 🛑 Emergency Stop

**Ctrl+C** in terminal will stop all services gracefully

Or force kill:
```bash
pkill automation
pkill node
```

---

## 📚 Next Steps

- Read the full [README.md](README.md) for architecture details
- Explore [internal/stealth/](internal/stealth/) for behavior algorithms
- Check [config/config.yaml](config/config.yaml) for customization
- Review code in [cmd/app/main.go](cmd/app/main.go) for API endpoints

---

**Remember**: This is for learning only. Respect platform terms of service! 🎓

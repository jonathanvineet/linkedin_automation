# LinkedIn Automation - Architecture Document

## System Overview

This application demonstrates a sophisticated LinkedIn automation system built with Go and React, featuring advanced anti-detection techniques and human behavior modeling.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     USER INTERFACE                       │
│              (React + TypeScript + Tailwind)             │
│                                                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐│
│  │Dashboard │  │ Persona  │  │ Activity │  │ Stealth ││
│  │  Panel   │  │  Config  │  │   Log    │  │  Panel  ││
│  └──────────┘  └──────────┘  └──────────┘  └─────────┘│
└─────────────────────┬───────────────────────────────────┘
                      │ HTTP/REST API
                      │ (JSON over HTTP)
┌─────────────────────▼───────────────────────────────────┐
│                    API LAYER                             │
│            (Gorilla Mux + CORS Handler)                  │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │  /api/status    /api/start     /api/stop          │ │
│  │  /api/search    /api/connect   /api/message       │ │
│  │  /api/stats     /api/activity  /api/persona       │ │
│  └────────────────────────────────────────────────────┘ │
└──────────────┬───────────────┬──────────────┬───────────┘
               │               │              │
    ┌──────────▼────┐  ┌───────▼─────┐  ┌────▼────────┐
    │   Browser     │  │  Behavior   │  │   State     │
    │   Session     │  │   Engine    │  │  Manager    │
    │   (Rod)       │  │             │  │  (SQLite)   │
    └───────┬───────┘  └──────┬──────┘  └─────────────┘
            │                 │
    ┌───────▼─────────────────▼─────┐
    │     Stealth Module             │
    │  - Mouse Movement              │
    │  - Typing Simulation           │
    │  - Timing & Delays             │
    │  - Scroll Behavior             │
    │  - Fingerprint Masking         │
    └───────┬────────────────────────┘
            │
    ┌───────▼────────────────────────┐
    │   Chrome Browser (CDP)         │
    │   - Headless/Headed Mode       │
    │   - Stealth Extensions         │
    │   - Session Persistence        │
    └────────────────────────────────┘
```

## Component Breakdown

### 1. Frontend Layer (React/TypeScript)

**Purpose**: User interface and visualization

**Key Components**:
- `Dashboard.tsx`: Main control panel
- `PersonaPanel.tsx`: Behavior configuration
- `ActivityLog.tsx`: Real-time action monitoring
- `StealthTechniques.tsx`: Detection evasion status
- `AutomationFlow.tsx`: Workflow visualization

**State Management**:
- React hooks (`useState`, `useEffect`)
- TanStack Query for API state
- Real-time polling (5-second intervals)

**Communication**:
- REST API client (`lib/api.ts`)
- JSON request/response format
- CORS-enabled

### 2. API Server (Go + Gorilla Mux)

**Purpose**: RESTful API gateway

**Endpoints**:
| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/status` | System health check |
| POST | `/api/start` | Initialize automation |
| POST | `/api/stop` | Shutdown automation |
| GET | `/api/stats` | Retrieve metrics |
| GET | `/api/activity` | Fetch activity logs |
| POST | `/api/persona` | Change behavior profile |
| POST | `/api/search` | Execute LinkedIn search |
| POST | `/api/connect` | Send connection request |
| POST | `/api/message` | Send follow-up message |

**Middleware**:
- CORS handler (allows frontend origins)
- JSON content-type enforcement
- Error handling & logging

### 3. Browser Automation (Rod Library)

**Purpose**: Chrome DevTools Protocol (CDP) control

**Capabilities**:
- Launch Chrome/Chromium
- Navigate to URLs
- Find and interact with elements
- Execute JavaScript
- Manage cookies & sessions
- Take screenshots (debugging)

**Stealth Integration**:
- `go-rod/stealth` package
- Custom fingerprint masking
- User agent rotation
- Automation flag removal

### 4. Authentication Module

**Purpose**: LinkedIn login and session management

**Flow**:
1. Navigate to `https://linkedin.com/login`
2. Type username with human-like delays
3. Type password with typo simulation
4. Click login button
5. Detect security challenges (2FA)
6. Verify login success
7. Save session cookies

**Error Handling**:
- Invalid credentials detection
- Security challenge detection
- Session expiration handling

### 5. Behavior Engine

**Components**:

**Persona System** (`internal/behavior/persona.go`):
```
┌─────────────────────────────────────────┐
│              Persona                    │
├─────────────────────────────────────────┤
│  - Name                                 │
│  - Typing Speed (WPM)                   │
│  - Mouse Precision (%)                  │
│  - Error Rate (%)                       │
│  - Attention Span (seconds)             │
│  - Break Frequency (actions)            │
│  - Scroll Impatience (low/med/high)     │
└─────────────────────────────────────────┘
         │
         ├─→ Recruiter: Methodical, careful
         ├─→ Founder: Fast, impatient
         └─→ Sales: Balanced, personalized
```

**Decision Engine** (`internal/behavior/decision_engine.go`):
```
Input: Current Context
  ├─ Time of Day
  ├─ Page Complexity
  ├─ Actions Today
  └─ Persona Traits
         │
         ▼
   Decision Making
  ├─ Should hover first?
  ├─ How long to think?
  ├─ Should hesitate?
  ├─ Should scroll?
  └─ Should re-read?
         │
         ▼
    Output: Actions
```

### 6. Stealth Techniques

**Mouse Movement** (`internal/stealth/mouse.go`):
```
Algorithm: Quadratic Bézier Curve

P(t) = (1-t)²P₀ + 2(1-t)tP₁ + t²P₂

Where:
  P₀ = Start point
  P₁ = Control point (randomized)
  P₂ = End point
  t  = Time parameter [0,1]

Features:
  - Variable speed
  - Overshoot + correction
  - Micro-adjustments
  - Idle wandering
```

**Typing Simulation** (`internal/stealth/typing.go`):
```
Input: Text to type

For each character:
  1. Calculate delay from WPM
     delay = 60000 / (WPM * 5 chars/word)
  
  2. Add variance (±40%)
     actual_delay = delay ± (delay * 0.4)
  
  3. Check for typo (based on error rate)
     if random() < error_rate:
       - Type wrong character (nearby key)
       - Pause (realize mistake)
       - Backspace
       - Type correct character
  
  4. Occasional think pauses
     Every 20 chars: 300-700ms pause

Output: Sequence of typing events
```

**Timing Jitter** (`internal/stealth/timing.go`):
```
Factors affecting delays:

Base Delay
  ├─ Minimum delay (configured)
  └─ Maximum delay (configured)

Context Modifiers:
  ├─ Time of Day Factor (1.0-1.3x)
  │   └─ Slower during lunch/evening
  ├─ Fatigue Factor (1.0-2.0x)
  │   └─ Increases with actions/day
  ├─ Page Complexity (1.0-2.0x)
  │   └─ Longer for complex pages
  └─ Persona Impatience (0.8-1.2x)
      └─ Based on persona traits

Final Delay = Base × Modifiers × Random(0.8-1.2)
```

### 7. State Persistence

**Database Schema** (SQLite):
```sql
connection_requests
  ├─ id              INTEGER PRIMARY KEY
  ├─ profile_url     TEXT UNIQUE NOT NULL
  ├─ profile_name    TEXT
  ├─ note            TEXT
  ├─ status          TEXT (sent/accepted/rejected)
  ├─ sent_at         DATETIME
  └─ accepted_at     DATETIME

messages
  ├─ id              INTEGER PRIMARY KEY
  ├─ profile_url     TEXT NOT NULL
  ├─ content         TEXT NOT NULL
  └─ sent_at         DATETIME

activity_logs
  ├─ id              INTEGER PRIMARY KEY
  ├─ timestamp       DATETIME
  ├─ action          TEXT NOT NULL
  ├─ type            TEXT (info/success/warning/error)
  └─ details         TEXT

session_data
  ├─ key             TEXT PRIMARY KEY
  ├─ value           TEXT
  └─ updated_at      DATETIME
```

**Indexes**:
- `idx_connection_status` ON `connection_requests(status)`
- `idx_activity_timestamp` ON `activity_logs(timestamp DESC)`

### 8. Logging System

**Structured Logging** (JSON format):
```json
{
  "level": "info",
  "msg": "Automation action executed",
  "action": "connection_sent",
  "persona": "recruiter",
  "delay_ms": 45000,
  "target": "https://linkedin.com/in/johndoe",
  "timestamp": "2024-12-16T14:32:15Z",
  "metadata": {
    "think_time_ms": 2300,
    "hover_duration_ms": 450,
    "typing_errors": 2
  }
}
```

**Log Levels**:
- `DEBUG`: Detailed execution flow
- `INFO`: Normal operations
- `WARN`: Suspicious behavior, rate limits
- `ERROR`: Failures, exceptions

## Data Flow

### Connection Request Flow

```
1. User clicks "Send Connection" in UI
   │
2. Frontend sends POST /api/connect
   {
     profile_url: "...",
     note: "..."
   }
   │
3. API validates request & checks rate limits
   │
4. Behavior engine calculates delays
   │
5. Browser navigates to profile
   │
6. Stealth: Simulate reading profile
   - Scroll down
   - Pause
   - Scroll up (re-read)
   │
7. Stealth: Mouse movement to "Connect" button
   - Generate Bézier curve
   - Move with variable speed
   - Micro-corrections
   │
8. Decision: Should hover first? → Yes (60% chance)
   - Hover for 200-800ms
   │
9. Click "Connect" button
   │
10. Find "Add note" button, click
    │
11. Type personalized note
    - Character by character
    - Random delays
    - Occasional typos + backspace
    │
12. Click "Send" button
    │
13. Save to database
    │
14. Update activity log
    │
15. Response to frontend
    {
      success: true,
      message: "Connection request sent"
    }
```

## Security & Anti-Detection

### Detection Vectors & Countermeasures

| Detection Method | Our Countermeasure |
|------------------|-------------------|
| `navigator.webdriver` | Removed via stealth library |
| Canvas fingerprinting | Noise injection |
| WebGL fingerprinting | Parameter randomization |
| Uniform timing | Jitter + context-aware delays |
| Straight mouse lines | Bézier curves + overshoots |
| Perfect typing | Errors + backspaces |
| Automation patterns | Persona-based variation |
| Rapid actions | Rate limiting + cooldowns |
| 24/7 activity | Business hours only |
| Same viewport size | Randomized dimensions |
| Same user agent | Rotating from real pool |

## Performance Characteristics

### Latency
- API response time: < 50ms (local)
- Browser action time: 3-8 seconds (with human delays)
- Full connection flow: 15-30 seconds

### Throughput
- Max connections/day: 20 (configurable, keep low!)
- Max messages/day: 10 (configurable)
- Actions per hour: 5-10 (with breaks)

### Resource Usage
- Memory: ~100-200 MB (Go) + ~500 MB (Chrome)
- CPU: < 5% idle, 10-20% during actions
- Disk: ~10 MB (database + logs)

## Scalability Considerations

**Current Design**: Single-user, single-session

**Limitations**:
- One LinkedIn account per instance
- No concurrent sessions
- Local database (SQLite)

**Future Enhancements** (Educational):
- Multi-account support (requires session management)
- Distributed execution (worker pool)
- PostgreSQL/MySQL for scale
- Message queuing (RabbitMQ/Redis)
- Kubernetes deployment

## Testing Strategy

### Unit Tests
- Persona behavior calculations
- Timing jitter algorithms
- Mouse curve generation
- Database operations

### Integration Tests
- API endpoint responses
- Browser session lifecycle
- State persistence

### Manual Tests
- Visual inspection (headless=false)
- Log analysis
- Detection evasion validation

## Deployment

**Development**:
```bash
./start.sh
```

**Production** (Not recommended!):
```bash
# Build
go build -o automation ./cmd/app

# Run with environment
API_PORT=8090 ./automation
```

**Docker** (Future):
```dockerfile
FROM golang:1.21-alpine
# ... build steps
CMD ["./automation"]
```

---

## Maintenance & Updates

### Regular Updates Needed

1. **LinkedIn UI Changes**
   - Selectors in `connect/request.go`
   - Search result parsing
   - Login flow elements

2. **Detection Techniques**
   - New fingerprinting methods
   - Updated stealth scripts
   - Timing pattern adjustments

3. **Dependencies**
   - Go modules: `go get -u`
   - Node packages: `npm update`
   - Rod library updates

### Monitoring

**Health Checks**:
- `/api/status` endpoint
- Log file growth
- Database size
- Error rates

**Alerts** (Manual):
- Login failures
- Rate limit violations
- Unusual delays
- Detection warnings

---

*This architecture prioritizes learning and demonstration over production readiness.*

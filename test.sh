#!/bin/bash

# LinkedIn Automation PoC - Automated Testing Script
# Tests the complete system to ensure everything works

set -e

echo "🧪 LinkedIn Automation PoC - System Test"
echo "========================================="
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
API_URL="http://localhost:8090/api"
BACKEND_PID=""
FRONTEND_PID=""
TEST_FAILED=0

# Cleanup function
cleanup() {
    echo ""
    echo "🧹 Cleaning up..."
    
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
        echo "Stopped backend (PID: $BACKEND_PID)"
    fi
    
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
        echo "Stopped frontend (PID: $FRONTEND_PID)"
    fi
    
    if [ $TEST_FAILED -eq 0 ]; then
        echo -e "${GREEN}✅ All tests passed!${NC}"
        exit 0
    else
        echo -e "${RED}❌ Some tests failed${NC}"
        exit 1
    fi
}

trap cleanup EXIT INT TERM

# Test function
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_code=${5:-200}
    
    echo -n "Testing $name... "
    
    if [ "$method" == "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$API_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$API_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$http_code" == "$expected_code" ]; then
        echo -e "${GREEN}✓${NC} (HTTP $http_code)"
        return 0
    else
        echo -e "${RED}✗${NC} (Expected $expected_code, got $http_code)"
        echo "Response: $body"
        TEST_FAILED=1
        return 1
    fi
}

echo "📋 Phase 1: Build Verification"
echo "==============================="

# Check if binaries exist
echo -n "Checking backend binary... "
if [ -f "bin/automation" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC} Not found"
    echo "Run: go build -o bin/automation ./cmd/app"
    exit 1
fi

echo -n "Checking frontend build... "
if [ -d "dist" ]; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${YELLOW}⚠${NC} Not found (not required for dev mode)"
fi

echo ""
echo "🚀 Phase 2: Starting Services"
echo "=============================="

# Start backend
echo -n "Starting backend server... "
./bin/automation > /dev/null 2>&1 &
BACKEND_PID=$!
sleep 3

if kill -0 $BACKEND_PID 2>/dev/null; then
    echo -e "${GREEN}✓${NC} (PID: $BACKEND_PID)"
else
    echo -e "${RED}✗${NC} Failed to start"
    exit 1
fi

# Check if backend is responding
echo -n "Waiting for backend to be ready... "
for i in {1..10}; do
    if curl -s "$API_URL/status" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        break
    fi
    sleep 1
    if [ $i -eq 10 ]; then
        echo -e "${RED}✗${NC} Timeout"
        exit 1
    fi
done

echo ""
echo "🔍 Phase 3: API Endpoint Tests"
echo "==============================="

# Test status endpoint
test_endpoint "GET /api/status" "GET" "/status"

# Test stats endpoint
test_endpoint "GET /api/stats" "GET" "/stats"

# Test activity endpoint
test_endpoint "GET /api/activity" "GET" "/activity"

# Test persona change
test_endpoint "POST /api/persona" "POST" "/persona" '{"persona_type":"founder"}'

# Test invalid persona (should fail)
test_endpoint "POST /api/persona (invalid)" "POST" "/persona" '{"persona_type":"invalid"}' "400"

# Test search without auth (should fail gracefully)
echo -n "Testing POST /api/search (without auth)... "
response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/search" \
    -H "Content-Type: application/json" \
    -d '{"keywords":"test","location":"test","max_results":5}')
http_code=$(echo "$response" | tail -n1)
if [ "$http_code" == "200" ] || [ "$http_code" == "400" ] || [ "$http_code" == "401" ]; then
    echo -e "${GREEN}✓${NC} (HTTP $http_code - expected behavior)"
else
    echo -e "${YELLOW}⚠${NC} (HTTP $http_code - unexpected)"
fi

echo ""
echo "📊 Phase 4: Data Validation"
echo "============================"

# Validate status response structure
echo -n "Validating status response... "
status_json=$(curl -s "$API_URL/status")
if echo "$status_json" | jq -e '.running, .logged_in, .persona, .stealth, .timestamp' > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC} Invalid JSON structure"
    TEST_FAILED=1
fi

# Validate stats response structure
echo -n "Validating stats response... "
stats_json=$(curl -s "$API_URL/stats")
if echo "$stats_json" | jq -e '.connections_sent, .messages_sent, .cooldown_seconds, .daily_limit' > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC} Invalid JSON structure"
    TEST_FAILED=1
fi

# Validate activity response structure
echo -n "Validating activity response... "
activity_json=$(curl -s "$API_URL/activity")
if echo "$activity_json" | jq -e 'type == "array"' > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC} Should return an array"
    TEST_FAILED=1
fi

echo ""
echo "💾 Phase 5: Database Check"
echo "=========================="

if [ -f "data/automation.db" ]; then
    echo -n "Checking database tables... "
    tables=$(sqlite3 data/automation.db "SELECT name FROM sqlite_master WHERE type='table';" 2>/dev/null || echo "")
    
    required_tables=("connection_requests" "messages" "activity_logs" "session_data")
    all_found=true
    
    for table in "${required_tables[@]}"; do
        if ! echo "$tables" | grep -q "$table"; then
            all_found=false
            break
        fi
    done
    
    if $all_found; then
        echo -e "${GREEN}✓${NC} All tables present"
    else
        echo -e "${YELLOW}⚠${NC} Some tables missing (created on first use)"
    fi
else
    echo -e "${YELLOW}⚠${NC} Database not yet created (normal on first run)"
fi

echo ""
echo "📝 Phase 6: Configuration Validation"
echo "====================================="

# Check .env file
echo -n "Checking .env file... "
if [ -f ".env" ]; then
    if grep -q "LINKEDIN_EMAIL=" .env && grep -q "LINKEDIN_PASSWORD=" .env; then
        echo -e "${GREEN}✓${NC} Configured"
    else
        echo -e "${YELLOW}⚠${NC} Missing credentials"
    fi
else
    echo -e "${RED}✗${NC} Not found"
fi

# Check config.yaml
echo -n "Checking config.yaml... "
if [ -f "config/config.yaml" ]; then
    echo -e "${GREEN}✓${NC} Found"
else
    echo -e "${RED}✗${NC} Not found"
    TEST_FAILED=1
fi

echo ""
echo "⚡ Phase 7: Performance Tests"
echo "============================="

# Measure response times
endpoints=("/status" "/stats" "/activity")
for endpoint in "${endpoints[@]}"; do
    echo -n "Response time for $endpoint... "
    start_time=$(date +%s%N)
    curl -s "$API_URL$endpoint" > /dev/null
    end_time=$(date +%s%N)
    duration=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds
    
    if [ $duration -lt 200 ]; then
        echo -e "${GREEN}✓${NC} ${duration}ms"
    else
        echo -e "${YELLOW}⚠${NC} ${duration}ms (slower than expected)"
    fi
done

echo ""
echo "🔄 Phase 8: State Management Tests"
echo "==================================="

# Test persona persistence
echo -n "Testing persona change persistence... "
curl -s -X POST "$API_URL/persona" \
    -H "Content-Type: application/json" \
    -d '{"persona_type":"recruiter"}' > /dev/null

sleep 1

status=$(curl -s "$API_URL/status")
persona=$(echo "$status" | jq -r '.persona')

if [ "$persona" == "recruiter" ]; then
    echo -e "${GREEN}✓${NC} (Persona: $persona)"
else
    echo -e "${RED}✗${NC} (Expected 'recruiter', got '$persona')"
    TEST_FAILED=1
fi

echo ""
echo "🎯 Phase 9: Error Handling Tests"
echo "================================="

# Test invalid JSON
echo -n "Testing invalid JSON handling... "
response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/search" \
    -H "Content-Type: application/json" \
    -d 'invalid json')
http_code=$(echo "$response" | tail -n1)
if [ "$http_code" == "400" ]; then
    echo -e "${GREEN}✓${NC} (Correctly rejected)"
else
    echo -e "${YELLOW}⚠${NC} (Got HTTP $http_code)"
fi

# Test missing required fields
echo -n "Testing missing fields handling... "
response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/connect" \
    -H "Content-Type: application/json" \
    -d '{}')
http_code=$(echo "$response" | tail -n1)
if [ "$http_code" == "400" ]; then
    echo -e "${GREEN}✓${NC} (Correctly rejected)"
else
    echo -e "${YELLOW}⚠${NC} (Got HTTP $http_code)"
fi

echo ""
echo "📈 Test Summary"
echo "==============="
echo ""

if [ $TEST_FAILED -eq 0 ]; then
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✅ ALL TESTS PASSED!${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo "Your LinkedIn Automation PoC is working correctly!"
    echo ""
    echo "Next steps:"
    echo "1. Configure LinkedIn credentials in .env"
    echo "2. Review rate limits in config/config.yaml"
    echo "3. Start the frontend: npm run dev"
    echo "4. Open: http://localhost:8080"
else
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${RED}❌ SOME TESTS FAILED${NC}"
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo "Please review the errors above and:"
    echo "1. Check logs/app.log for details"
    echo "2. Verify .env configuration"
    echo "3. Ensure all dependencies are installed"
fi

echo ""

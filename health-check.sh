#!/bin/bash

# Health Check Script for LinkedIn Automation

API_URL="${API_URL:-http://localhost:8090/api}"

echo "🏥 LinkedIn Automation - Health Check"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if API is responding
echo "📡 Checking API connectivity..."
if curl -s -f "${API_URL}/status" > /dev/null 2>&1; then
    echo "✅ API is responding"
else
    echo "❌ API is not responding"
    echo "   Make sure the backend is running: ./bin/automation"
    exit 1
fi

# Get system status
echo ""
echo "📊 System Status:"
STATUS=$(curl -s "${API_URL}/status")
echo "$STATUS" | jq '.'

# Get statistics
echo ""
echo "📈 Statistics:"
STATS=$(curl -s "${API_URL}/stats")
echo "$STATS" | jq '.'

# Check activity log
echo ""
echo "📝 Recent Activity (last 5):"
ACTIVITY=$(curl -s "${API_URL}/activity")
echo "$ACTIVITY" | jq '.[0:5]'

# Check database
if [ -f "data/automation.db" ]; then
    echo ""
    echo "💾 Database Status:"
    echo "   Size: $(du -h data/automation.db | cut -f1)"
    echo "   Connections: $(sqlite3 data/automation.db 'SELECT COUNT(*) FROM connection_requests;' 2>/dev/null || echo 'N/A')"
    echo "   Messages: $(sqlite3 data/automation.db 'SELECT COUNT(*) FROM messages;' 2>/dev/null || echo 'N/A')"
    echo "   Activity Logs: $(sqlite3 data/automation.db 'SELECT COUNT(*) FROM activity_logs;' 2>/dev/null || echo 'N/A')"
fi

# Check logs
if [ -f "logs/app.log" ]; then
    echo ""
    echo "📋 Log File:"
    echo "   Size: $(du -h logs/app.log | cut -f1)"
    echo "   Last 3 errors:"
    tail -n 100 logs/app.log | grep -i "error" | tail -n 3 || echo "   No recent errors"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Health check complete"

#!/bin/bash
# Auto-executable wrapper - Makes quickfix.sh executable and runs it
chmod +x /workspaces/linkedin_automation/quickfix.sh 2>/dev/null || true
chmod +x /workspaces/linkedin_automation/*.sh 2>/dev/null || true
exec /workspaces/linkedin_automation/quickfix.sh

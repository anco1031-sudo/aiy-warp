#!/bin/bash
# Launch Brave with remote debugging for Aiy
if pgrep -x brave > /dev/null; then
  if curl -s http://127.0.0.1:9222/json/version > /dev/null 2>&1; then
    echo "✓ Brave already running with debugging on port 9222"
    exit 0
  else
    echo "Brave running but no debugging — restarting with port 9222..."
    pkill brave
    sleep 2
  fi
fi
echo "Launching Brave with --remote-debugging-port=9222 ..."
DISPLAY=:0 /opt/brave-bin/brave --remote-debugging-port=9222 > /tmp/brave_debug.log 2>&1 & disown
sleep 3
curl -s http://127.0.0.1:9222/json/version | head -n 5
echo "Done"

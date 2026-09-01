---
name: brave-ctl
description: Control Brave with remote debugging for Aiy's MCP automation. Launch, check, and manage Brave's debugging port.
---

# Brave Control Skill — for Aiy

This skill lets Aiy launch and manage Brave with remote debugging so `brave-devtools` MCP can connect to Louis's real Brave profile (shared cookies/login).

## When to use
- Louis closed Brave and Aiy needs to reopen it
- Need to check if Brave is running with debugging port
- Need to switch between isolated vs shared mode

## Commands

### 1. Launch Brave with debugging (shared profile — keeps login)
```bash
DISPLAY=:0 /opt/brave-bin/brave --remote-debugging-port=9222 & disown
# or
brave --remote-debugging-port=9222 & disown
```
Check:
```bash
curl -s http://127.0.0.1:9222/json/version | head
```

### 2. Launch Brave isolated (clean, no login — for testing)
```bash
# Handled by MCP directly: npx brave-mcp --headless --isolated
```

### 3. Check status
```bash
curl -s http://127.0.0.1:9222/json/list | python3 -m json.tool | head
ps aux | grep brave | grep -v grep
opencode mcp list | grep brave
```

### 4. Kill and relaunch
```bash
pkill brave; sleep 1; DISPLAY=:0 brave --remote-debugging-port=9222 & disown
```

## Persistent setup (already done for Louis)
- User desktop file: `~/.local/share/applications/brave-browser.desktop`
- Exec is patched to `brave --remote-debugging-port=9222 %U`
- Every click on Brave icon now auto-opens with port 9222 (shared profile)
- No manual command needed after this

## MCP Config (in ~/.config/opencode/opencode.jsonc)
```json
"brave-devtools": {
  "type": "local",
  "command": ["npx", "-y", "brave-mcp@latest", "--browserUrl", "http://127.0.0.1:9222"],
  "enabled": true
}
```
For isolated mode, switch to:
```json
"command": ["npx", "-y", "brave-mcp@latest", "--executable-path", "/usr/bin/brave", "--headless", "--isolated"]
```

## Notes
- Shared mode = Aiy sees all Louis's tabs (privacy note)
- Isolated mode = clean, no cookies, safe for testing
- Louis's Brave profile is at `~/.config/BraveSoftware/Brave-Browser/Default`

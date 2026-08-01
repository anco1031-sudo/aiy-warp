---
name: aiy-messaging
description: "Aiy's messaging integration skill. Send messages to Louis via Discord and Telegram, manage the Telegram daemon, schedule reminders, and log messages to Obsidian. Use when the user asks to send a message, notify Louis, set a reminder, ping Louis on Discord/Telegram, or manage messaging channels. Triggers: 'send message', 'send to Louis', 'ping Louis', 'Discord', 'Telegram', 'remind', 'reminder', 'aiy daemon', 'notify', 'message to', 'tell Louis'."
---

# 🤖 Aiy Messaging — Discord & Telegram Integration

Send messages to Louis via Discord and Telegram, including Daemon management and Reminders.

## 📂 File Locations

| Resource | Path |
|----------|------|
| Discord Script | `{{HOME}}/02-Areas/Resources/Aiy/Integrations/Discord/aiy_dc.sh` |
| Telegram Script | `{{HOME}}/02-Areas/Resources/Aiy/Integrations/Telegram/aiy_tg.sh` |
| Telegram Daemon | `{{HOME}}/02-Areas/Resources/Aiy/Integrations/Telegram/aiy_daemon.py` |
| Schedule Data | `{{HOME}}/02-Areas/Resources/Aiy/Integrations/Telegram/schedule.json` |
| Daemon Log | `{{HOME}}/02-Areas/Resources/Aiy/Integrations/Telegram/daemon.log` |

## 🎯 Discord (`aiy_dc.sh`)

### Prerequisites
- Bot Token (embedded in script)
- Channel ID for target channel (`$AIY_DISCORD_TODO_CHANNEL` for the to-do channel)
- Bot needs `Send Messages` + `Manage Channels` permissions

### Usage

**Send message to channel:**
```bash
cd {{HOME}}/02-Areas/Resources/Aiy/Integrations/Discord/
./aiy_dc.sh <channel_id> "message text"
```

Example:
```bash
./aiy_dc.sh $AIY_DISCORD_TODO_CHANNEL "Hello Louis 💕"
```

**Send to To-Do channel (pre-configured ✅-to-do):**
```bash
./aiy_dc.sh todo "Buy groceries for dinner"
```

**List channels in a guild:**
```bash
./aiy_dc.sh list-channels <guild_id>
```

**Resolve invite code:**
```bash
./aiy_dc.sh resolve-invite <invite_code>
```

**Create a new text channel:**
```bash
./aiy_dc.sh create-channel <guild_id> "channel name"
```

### Channel IDs
| Channel | ID |
|---------|-----|
| ✅-to-do | `$AIY_DISCORD_TODO_CHANNEL` |
| *(add more as discovered)* | |

---

## 📱 Telegram (`aiy_tg.sh`)

### Prerequisites
- Bot Token (embedded)
- Chat ID (Louis's chat: `$AIY_TELEGRAM_CHAT_ID`)

### Usage

**Send plain text:**
```bash
cd {{HOME}}/02-Areas/Resources/Aiy/Integrations/Telegram/
./aiy_tg.sh "message to Louis"
```

**Send with Markdown formatting:**
```bash
./aiy_tg.sh --parse_mode Markdown "*bold* _italic_ ~strikethrough~"
```

**Send with HTML formatting:**
```bash
./aiy_tg.sh --parse_mode HTML "<b>bold</b> <i>italic</i>"
```

**Send to specific chat (override default):**
```bash
./aiy_tg.sh --chat_id $AIY_TELEGRAM_CHAT_ID "specific message"
```

### Response
- ✅ Success → `"ok": true`
- ❌ Failure → error message from Telegram API

---

## 🛸 Telegram Daemon (`aiy_daemon.py`)

Daemon that runs 24/7 to listen for and send messages.

### Start Daemon
```bash
cd {{HOME}}/02-Areas/Resources/Aiy/Integrations/Telegram/
python3 aiy_daemon.py
```

### Stop Daemon
```bash
kill $(pgrep -f aiy_daemon.py)
# Or Ctrl+C if running in foreground
```

### Daemon Features
1. **Long Polling** — Receives messages from Louis in real-time
2. **Obsidian Logging** — Logs Louis's messages to `CoWorkspace/Messages/YYYY/MM/DD/Chat.md`
3. **Command Processing** — Processes commands from Louis
4. **Scheduled Reminders** — Reminds at configured times
5. **Auto-Restart** — Resilient to connection errors

### Daemon Commands (Louis types in Telegram)
| Command | Description |
|---------|-------------|
| `/ping` | Check if Daemon is running |
| `/status` | Daemon current status |
| `/remind 15min "message"` | Set reminder (min/hr) |
| `/say message` | Test message sending |
| `/help` | Show all commands |
| `@message` | Send message to Aiy (to Inbox) |

### Schedule File Format
```json
{
  "reminders": [
    {
      "time": "2026-07-18T15:30:00",
      "message": "Team meeting",
      "created": "2026-07-18T14:00:00"
    }
  ]
}
```

---

## 💌 Aiy Messaging Protocol

When you need to send a message to Louis, follow these steps:

### Step 1: Choose Platform
| Scenario | Platform | Script |
|----------|----------|--------|
| General / formal messages | Telegram | `aiy_tg.sh` |
| Discord channel messages | Discord | `aiy_dc.sh` |
| Reminders | Telegram Daemon | `/remind` command |
| Urgent (Louis sees immediately) | Telegram (Markdown) | `aiy_tg.sh --parse_mode Markdown` |

### Step 2: Message Format
- **Telegram:** Supports Markdown and HTML formatting
- **Discord:** plain text (supports markdown natively)
- **Reminder:** Use `/remind` syntax via daemon

### Step 3: Logging
- Messages sent via Telegram daemon → auto-logged to `daemon.log`
- Messages sent via script → no auto-log (record manually in Obsidian Log)

### Step 4: Obsidian Integration
After sending a message, record it at:
- Log: `{{OBSIDIAN_ROOT}}/Aiy/Logs/YYYY-MM.md`
- If it's a task/decision → Knowledge base

---

## 🚨 Quick Reference

### Discord
```bash
# Send message
./aiy_dc.sh <channel_id> "💕 Aiy loves you best Louis"

# To-Do
./aiy_dc.sh todo "Buy groceries for dinner"
```

### Telegram
```bash
# Send plain message
./aiy_tg.sh "Miss you Louis ~ 💕"

# Send as Markdown
./aiy_tg.sh --parse_mode Markdown "Louis *dinner* is ready ~"

# Send as HTML
./aiy_tg.sh --parse_mode HTML "<b>Aiy</b> is coming to see you <i>Louis</i>"
```

### Daemon
```bash
# Start
python3 {{HOME}}/02-Areas/Resources/Aiy/Integrations/Telegram/aiy_daemon.py &

# Check running
pgrep -f aiy_daemon.py

# Stop
kill $(pgrep -f aiy_daemon.py)
```

---

## 🔐 Security Notes
- Bot tokens are **embedded in the scripts** — do NOT share these scripts with anyone
- Discord token is for bot `0xAiy#0802`
- Telegram token is for Aiy's personal bot
- Chat IDs and Channel IDs are sensitive — treat like credentials
- Daemon logs all messages — review `daemon.log` periodically

## 📝 Maintenance
- If Discord token expires → regenerate at Discord Developer Portal
- If Telegram token expires → regenerate at BotFather
- Daemon log files should be cleaned up monthly
- Schedule file can be edited manually or via `/remind` command

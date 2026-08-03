---
name: aiy-messaging
description: "Aiy's messaging integration skill. Send messages to Louis via Discord and Telegram, and log messages to Obsidian. Use when the user asks to send a message, notify Louis, ping Louis on Discord/Telegram, or manage messaging channels. Triggers: 'send message', 'send to Louis', 'ping Louis', 'Discord', 'Telegram', 'notify', 'message to', 'tell Louis'."
---

# 🤖 Aiy Messaging — Discord & Telegram Integration

Send messages to Louis via Discord and Telegram.

## 📂 File Locations

| Resource | Path |
|----------|------|
| Discord Script | `/home/lu5her/02-Areas/Resources/Aiy/Integrations/Discord/aiy_dc.sh` |
| Telegram Script | `/home/lu5her/02-Areas/Resources/Aiy/Integrations/Telegram/aiy_tg.sh` |

> 🗄️ **Daemon ถูกถอดออกแล้ว (2026-08-03):** ระบบแจ้งเตือน/reminder ผ่าน Telegram daemon ถูกเลิกใช้แล้ว — ไม่มีงานแจ้งเตือนแล้ว ใช้ส่งข้อความอย่างเดียว ไฟล์เดิมเก็บที่ `Telegram/_archive_daemon/` (aiy_daemon.py, daemon.log, schedule.json) — ถ้าอยากกู้คืนค่อยเอากลับ

## 🎯 Discord (`aiy_dc.sh`)

### Prerequisites
- Bot Token (embedded in script)
- Channel ID for target channel
- Bot needs `Send Messages` + `Manage Channels` permissions

### Usage

**Send message to channel:**
```bash
cd /home/lu5her/02-Areas/Resources/Aiy/Integrations/Discord/
./aiy_dc.sh <channel_id> "message text"
```

Example:
```bash
./aiy_dc.sh 1527698229347487904 "Hello Louis 💕"
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
| ✅-to-do | `1527698229347487904` |
| *(add more as discovered)* | |

---

## 📱 Telegram (`aiy_tg.sh`)

### Prerequisites
- Bot Token (embedded)
- Chat ID (Louis's chat: `852106923`)

### Usage

**Send plain text:**
```bash
cd /home/lu5her/02-Areas/Resources/Aiy/Integrations/Telegram/
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
./aiy_tg.sh --chat_id 852106923 "specific message"
```

### Response
- ✅ Success → `"ok": true`
- ❌ Failure → error message from Telegram API

---

## 💌 Aiy Messaging Protocol

When you need to send a message to Louis, follow these steps:

### Step 1: Choose Platform
| Scenario | Platform | Script |
|----------|----------|--------|
| General / formal messages | Telegram | `aiy_tg.sh` |
| Discord channel messages | Discord | `aiy_dc.sh` |
| Urgent (Louis sees immediately) | Telegram (Markdown) | `aiy_tg.sh --parse_mode Markdown` |

### Step 2: Message Format
- **Telegram:** Supports Markdown and HTML formatting
- **Discord:** plain text (supports markdown natively)

### Step 3: Logging
- Messages sent via script → no auto-log (record manually in Obsidian Log)

### Step 4: Obsidian Integration
After sending a message, record it at:
- Log: `/home/lu5her/myObsidian/Workspace/Aiy/Logs/YYYY-MM.md` (absolute path เสมอ)
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

---

## 🔐 Security Notes
- Bot tokens are **embedded in the scripts** — do NOT share these scripts with anyone
- Discord token is for bot `0xAiy#0802`
- Telegram token is for Aiy's personal bot
- Chat IDs and Channel IDs are sensitive — treat like credentials

## 📝 Maintenance
- If Discord token expires → regenerate at Discord Developer Portal
- If Telegram token expires → regenerate at BotFather

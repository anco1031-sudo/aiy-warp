# 🏗️ CONSTRUCTION — โครงสร้าง aiy-warp & วิธี Init เครื่องใหม่

> **Source of Truth:** repo นี้ (`01-Projects/aiy-warp`) คือ canonical kit — ทุกอย่างที่เกี่ยวกับอัยรวมศูนย์ที่นี่ แล้ว `warp` ไปเครื่องอื่น

## 📦 โครงสร้างปัจจุบัน (v0.3 — รวมศูนย์ 24 ส.ค. 2026)

```
aiy-warp/
├── agents/                 18 agents (aiy, lin, an, pao, mint, fah, cloud, jean, cher, ploy, jewel, sai, pleng, kwan, fon, june, bee, nam)
├── skills/                 14 skills
│   ├── aiy/                persona หลัก
│   ├── aiy-messaging/      Discord/Telegram
│   ├── affiliate-poster/   ของมันต้องมี
│   ├── obsidian/           vault CRUD
│   ├── teach-protocol/     INGEST→SKILL-FORGE
│   ├── handoff-protocol/   ใบงานมาตรฐาน
│   ├── trading-pipeline/   Kwan team
│   ├── warp-sync/          sync live → repo
│   ├── sync-copilot-prompts/
│   ├── media-ocr/          🆕 Tesseract tha+eng (Aetox 3.1/3.3)
│   ├── media-transcribe/   🆕 ffmpeg+Whisper (Aetox 3.2)
│   ├── pdf-toolkit/        🆕 pdftotext+pymupdf (Aetox 3.4)
│   ├── office-suite/       🆕 openpyxl/docx/pptx (Aetox 4)
│   └── browser-control/    🆕 Playwright MCP (Aetox 5)
├── scripts/                🆕 รวมศูนย์ helper
│   ├── warp-init.sh        เช็ค dependencies ทั้งหมด
│   ├── media/ocr.py, transcribe.py, video_ocr.py
│   ├── pdf/read.py
│   ├── office/sheet_write.py, doc_write.py, slide_write.py
│   ├── browser/control.py
│   └── utils/model-check.ts
├── config/                 🆕 ต้นฉบับ live config (sanitized)
│   ├── opencode.json       provider 9router (apiKey = PLACEHOLDER)
│   ├── opencode.jsonc      เต็ม (MCP playwright/chrome-devtools)
│   ├── AGENTS.md           global rules
│   ├── model-mapping.json  role→model
│   ├── oh-my-openagent.json
│   ├── tui.json
│   └── git-team/git-team.sh
├── templates/              Obsidian templates
├── playbooks/              Team Charter, ADRs
├── install/                Go CLI (aiy-warp binary)
│   ├── bundles.yaml        org chart (warp_version: 1)
│   └── src/                Go source
├── scaffold/               PARA 00-05
└── docs/                   WARP.md + CONSTRUCTION.md (ไฟล์นี้)
```

## 🔄 ย้ายอะไรไปไหน — ตอน `warp install` ที่เครื่องใหม่

| ใน repo (`aiy-warp/`) | ไปที่ (`$HOME`) | วิธี |
|---|---|---|
| `agents/*.md` | `~/.config/opencode/agents/` | `aiy-warp install opencode` |
| `skills/*/SKILL.md` | `~/.config/opencode/skills/*/` | `install` |
| `scripts/utils/model-check.ts` | `~/.config/opencode/scripts/` | 🆕 `install` (เพิ่มใน bundles.yaml) |
| `scripts/media/*.py`, `pdf/*.py`, `office/*.py`, `browser/*.py` | `~/.config/opencode/scripts/` + ติดไปกับ skill | เรียกตรงจาก repo ก็ได้ (`python3 scripts/media/ocr.py`) |
| `config/opencode.json` | `~/.config/opencode/opencode.json` | copy — **ต้องเติม apiKey จริงเอง** (PLACEHOLDER) |
| `config/opencode.jsonc` | `~/.config/opencode/opencode.jsonc` | copy |
| `config/AGENTS.md` | `~/.config/opencode/AGENTS.md` | copy |
| `config/model-mapping.json` | `~/.config/opencode/model-mapping.json` | copy |
| `config/git-team/git-team.sh` | `~/.config/opencode/git-team/` | copy + `chmod +x` |
| `scaffold/` | `~/00-Inbox` ... `~/05-Other` | `aiy-warp init-workspace` |
| `scripts/warp-init.sh` | รันเช็ค deps | `bash scripts/warp-init.sh` |

> **Private soul ไม่ย้าย:** `~/myObsidian/Workspace/Aiy/{Knowledge,Logs,Memory}` อยู่ใน private vault — ไม่ commit เข้า public warp

## 🚀 Quick Start — เครื่องใหม่ (copy-paste ได้เลย)

```bash
# 1. Clone kit
git clone https://github.com/anco1031-sudo/aiy-warp.git ~/01-Projects/aiy-warp
cd ~/01-Projects/aiy-warp

# 2. Build CLI (ต้องมี Go 1.22+)
cd install/src && go build -o aiy-warp . && cd ../..

# 3. ติดตั้ง identities + skills + scripts + config
./install/src/aiy-warp install opencode -y
# ถ้า CLI ยังไม่รองรับ scripts/config ใหม่ — fallback manual:
mkdir -p ~/.config/opencode/scripts ~/.config/opencode/git-team
cp -r scripts/utils/* ~/.config/opencode/scripts/
cp -r scripts/media scripts/pdf scripts/office scripts/browser ~/.config/opencode/scripts/ 2>/dev/null || true
cp config/AGENTS.md config/model-mapping.json config/tui.json ~/.config/opencode/
cp config/git-team/git-team.sh ~/.config/opencode/git-team/ && chmod +x ~/.config/opencode/git-team/git-team.sh
# opencode.json — อย่าทับ apiKey ถ้ามีของเดิม:
cp -n config/opencode.json ~/.config/opencode/opencode.json  # -n = ไม่ทับถ้ามีแล้ว
cp -n config/opencode.jsonc ~/.config/opencode/opencode.jsonc
# แล้วเติม apiKey จริงใน ~/.config/opencode/opencode.json (แทน PLACEHOLDER)

# 4. Scaffold workspace (ถ้ายังไม่มี PARA)
./install/src/aiy-warp init-workspace

# 5. เช็ค dependencies (poppler, tesseract, ffmpeg, whisper, office libs, playwright)
bash scripts/warp-init.sh
# ขาดอะไรติดตั้งตามที่สคริปต์บอก:
# sudo apt install tesseract-ocr tesseract-ocr-tha poppler-utils ffmpeg
# pip install openai-whisper openpyxl python-docx python-pptx playwright && playwright install chromium

# 6. Verify
./install/src/aiy-warp doctor
./install/src/aiy-warp list
# ควรเห็น 18 agents + 14 skills + doctor exit 0

# 7. ทดสอบ media (ถ้าลง deps แล้ว)
python3 scripts/media/ocr.py /path/to/image.jpg --langs tha+eng
python3 scripts/pdf/read.py /path/to/doc.pdf
```

## 🔧 Dependencies ทั้งหมด (อ้างอิง Aetox guide)

| ตัว | ติดตั้ง | เช็ค |
|---|---|---|
| `tesseract 5.5+` + `tha+eng` | `sudo apt install tesseract-ocr tesseract-ocr-tha` | `tesseract --list-langs` |
| `poppler` (`pdftotext`) | `sudo apt install poppler-utils` | `pdftotext -v` |
| `ffmpeg` | `sudo apt install ffmpeg` | `ffmpeg -version` |
| `openai-whisper` | `pip install openai-whisper` | `python3 -c "import whisper"` |
| `openpyxl, python-docx, python-pptx` | `pip install openpyxl python-docx python-pptx` | `scripts/warp-init.sh` |
| `playwright` | `pip install playwright && playwright install chromium` | `python3 -c "import playwright"` |
| `Go 1.22+` | สำหรับ build CLI | `go version` |

## 📝 วิธี Sync หลังแก้ที่เครื่องแม่ (warp-sync)

```bash
# ที่เครื่องแม่ (MASTER) หลังแก้ agents/skills/config/scripts
cd ~/01-Projects/aiy-warp
# 1. ซิงค์ live → repo (ตาม warp-sync skill)
# 2. รัน gates
git grep -lP "[\x{4E00}-\x{9FFF}]" -- . || echo "✅ no CJK"
./install/src/aiy-warp list && ./install/src/aiy-warp doctor
# 3. commit + push
git add -A && git commit -m "🔄 sync: ..." && git push
# เครื่องอื่นก็ git pull + ./install/src/aiy-warp install opencode -y
```

## ⚠️ ข้อระวัง
- `config/opencode.json` ใน repo มี `apiKey = PLACEHOLDER` — **ห้าม commit key จริง** — ไปเติมที่เครื่องใหม่เอง
- อย่ารัน `install opencode` บนเครื่องแม่ (live ใหม่กว่า + เสี่ยงทับ) — ใช้ `warp-sync` แทน
- `affiliate-poster` มี `{{HOME}}` placeholder — installer ต้องแทนค่าตอน install (P1 debt เดิม)

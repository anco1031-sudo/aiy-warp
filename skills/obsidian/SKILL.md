---
name: obsidian
description: "Universal Obsidian vault operations — standard CRUD, frontmatter, templates, naming conventions, linking, tags, MOC management, and vault navigation for ALL agents. Use when any agent needs to read/write/search/manage notes in the Obsidian vault. Triggers: 'note', 'obsidian', 'vault', 'daily log', 'knowledge entry', 'MOC', 'frontmatter', 'wiki link', 'dataview', 'templater', 'file.note', 'บันทึก', 'โน้ต', 'templater', 'frontmatter', 'metadata'."
---

# Obsidian Universal Vault Operations

> มาตรฐานกลางสำหรับทุก agent ในการทำงานกับ Obsidian vault — อ่าน, เขียน, ค้นหา, จัดการ notes, templates, links, tags

## 📍 Vault Location

```
VAULT_ROOT = /home/lu5her/myObsidian/Aiy_Workspace
```

**หมายเหตุ:** ปัจจุบัน vault ถูก mount แยกเป็น `Aiy_Workspace/` — ถ้ามีการ sync หรือย้าย paths ในอนาคต ให้อัปเดตที่ `Knowledge-Obsidian-Vault.md` และแจ้ง Aiy

---

## 🔧 Two Operation Modes

| Mode | Method | When to Use |
|------|--------|-------------|
| **MCP** | `read_mcp_resource`, `list_mcp_resources` | When Obsidian MCP server is connected and tools are available |
| **Direct** | `read`, `write`, `edit`, `glob`, `grep`, `bash` | **Default** — always works, MCP or not |

### Priority Rule
1. ถ้า `read_mcp_resource` หรือ `list_mcp_resources` สามารถใช้กับ Obsidian server ได้ → **ใช้ MCP mode**
2. ถ้าไม่ — **ใช้ Direct mode** ทันที (tools พื้นฐานมีครบ)

---

## 📁 Vault Structure — มาตรฐาน

```
Aiy_Workspace/
├── Aiy/                       # Aiy's personal workspace
│   ├── Handoff/               # Session handoff files
│   ├── Integrations/          # Integration notes
│   ├── Knowledge/             # Knowledge entries
│   │   ├── INDEX.md           # MOC for Knowledge
│   │   └── Knowledge-*.md     # Individual entries
│   └── Logs/                  # Monthly logs
│       └── YYYY-MM.md         # e.g. 2026-07.md
│
├── Content_Team/              # Cher / Ploy / Jewel workspace
├── CoWorkspace/               # Cross-team collaboration
│   ├── Handoff/
│   ├── Kanban/
│   ├── Messages/
│   └── Standup/
│
├── Legal_Team/                # Jean workspace
├── Product_Team/              # Lin / An / Pao / Mint / Fah / Cloud
├── Security_Team/             # Sai / Pleng
├── Trader_Team/               # Kwan / Fon / June / Bee / Nam
│
├── Meta/                      # Templates
│   ├── Template_Daily_Log.md
│   ├── Template_Knowledge.md
│   └── ...
│
├── Shared/                    # Cross-team shared resources
│   ├── Knowledge_INDEX.md     # CENTRAL Knowledge INDEX (all teams)
│   ├── Team_Charter.md
│   └── Architecture_Decisions.md
│
└── Dashboard.md               # Team Status Dashboard (Dataview-driven)
```

### Vault Root Files
| File | Purpose |
|------|---------|
| `Dashboard.md` | Dataview-powered team dashboard — overview, recent logs, kanban |
| `Shared/Knowledge_INDEX.md` | Central INDEX — รวมทุก Knowledge entry จากทุกทีม |

---

## 🏷️ Frontmatter Standard — YAML Metadata

Every note MUST have frontmatter. มาตรฐาน:

### Log Entry
```yaml
---
date: YYYY-MM-DD
type: log
created: Ddd DD,MMM YY HH:mm
updated: Ddd DD,MMM YY HH:mm
author: AgentName
---
```

### Knowledge Entry
```yaml
---
date: YYYY-MM-DD
type: knowledge
created: Ddd DD,MMM YY HH:mm
updated: Ddd DD,MMM YY HH:mm
tags: [tag1, tag2]
---
```

### Dashboard / MOC / Index
```yaml
---
created: Ddd DD,MMM YY HH:mm
updated: Ddd DD,MMM YY HH:mm
type: index
---
```

### Project / Deliverable Note
```yaml
---
title: Project Name
status: draft | active | completed | archived
created: YYYY-MM-DD
updated: YYYY-MM-DD
tags: [project, category]
owner: AgentName
---
```

### General Note
```yaml
---
created: YYYY-MM-DD
updated: YYYY-MM-DD
tags: [tag1, tag2]
---
```

### Date Format
```
Ddd DD,MMM YY HH:mm
Example: Sun 26,Jul 26 22:15
```

### Frontmatter Rules
- `created` = date note was first created (never change)
- `updated` = last modified date (update every edit)
- `type` = one of: `log`, `knowledge`, `index`, `dashboard`, `project`, `note`
- `tags` = lowercase, kebab-case, array format

---

## 📝 Standard Operations

### 1. Create Note

```markdown
write
filePath: "VAULT_ROOT/{team_folder}/{filename}.md"
content: |-
  ---
  created: {current datetime}
  updated: {current datetime}
  tags: []
  ---

  # {Title}

  {content}
```

**Naming Convention:**
| Note Type | Format | Example |
|-----------|--------|---------|
| Daily Log | `YYYY-MM.md` | `2026-07.md` |
| Knowledge | `Knowledge-{Topic}.md` | `Knowledge-ICT-SMC-Roadmap.md` |
| Knowledge INDEX | `INDEX.md` | `Aiy/Knowledge/INDEX.md` |
| Handoff | `Handoff-{from}-{to}-{topic}.md` | `Handoff-Aiy-Kwan-ADVANC.md` |
| General | `{slug}.md` | `meeting-notes.md` |
| Dashboard | `Dashboard.md` | `Dashboard.md` |
| Standup | `Standup.md` | `CoWorkspace/Standup/YYYY/MM/DD/Standup.md` |

### 2. Read Note

```markdown
read
filePath: "VAULT_ROOT/{path}/{filename}.md"
```

**Tip:** For large files, use `limit` and `offset` parameters.

### 3. Update Note (Edit)

```markdown
edit
filePath: "VAULT_ROOT/{path}/{filename}.md"
oldString: "{exact text to replace}"
newString: "{new text}"
```

### 4. Search Notes

**By filename:**
```markdown
glob
pattern: "**/*{query}*.md"
path: "VAULT_ROOT"
```

**By content:**
```markdown
grep
pattern: "{regex pattern}"
include: "*.md"
path: "VAULT_ROOT"
```

**By tags (frontmatter):**
```markdown
grep
pattern: "tags:.*{tag}"
include: "*.md"
path: "VAULT_ROOT"
```

---

## 🔗 Linking Standard (Wiki Links)

Obsidian ใช้ `[[wiki links]]` — ทุก agent ต้องใช้มาตรฐานนี้:

### Internal Links
```markdown
[[Knowledge-Socratic-Session|Socratic Session]]
```
- `[[filename|Display Name]]`
- ใช้ **filename ไม่รวม .md**
- ถ้าชื่อเหมือนกัน ไม่ต้องมี pipe

### Link to Headings
```markdown
[[Knowledge-Socratic-Session#Protocol]]
```

### Link to Blocks
```markdown
[[Knowledge-Socratic-Session#^block-id]]
```

### Embedding
```markdown
![[Logs/2026-07.md#07-30]]
```

### Tagging
```yaml
# ใน frontmatter
tags: [trading, strategy, ICT]

# หรือ inline ในเนื้อหา
#trading #ICT #SMC
```

**Tag Convention:** lowercase, kebab-case, 1-3 words max

---

## 📋 Templates — ใช้ Template เมื่อ

| Situation | Template | Location |
|-----------|----------|----------|
| Create daily log entry | Append to `Logs/YYYY-MM.md` ตาม format `Template_Daily_Log.md` | `Meta/Template_Daily_Log.md` |
| Create knowledge entry | `Template_Knowledge.md` | `Meta/Template_Knowledge.md` |
| Create handoff | CoWorkspace Handoff format | `CoWorkspace/Handoff/YYYY/MM/DD/` |
| Create standup | CoWorkspace Standup format | `CoWorkspace/Standup/YYYY/MM/DD/` |

---

## 🗺️ MOC (Map of Content) / INDEX Management

ทุกทีมมี INDEX.md ของตัวเอง + รวมใน Central INDEX

### Local INDEX (ต่อทีม)
```
Aiy/Knowledge/INDEX.md           — Aiy's knowledge MOC
Shared/Knowledge_INDEX.md         — Central INDEX (รวมทุกทีม)
```

### เมื่อสร้าง Knowledge Entry ใหม่ → ต้อง:
1. **Create entry** → `{team}/Knowledge/Knowledge-{topic}.md`
2. **Update Local INDEX** → เพิ่มแถวใน `{team}/Knowledge/INDEX.md`
3. **Update Central INDEX** → เพิ่มแถวใน `Shared/Knowledge_INDEX.md`

### INDEX Table Format
```markdown
| [[Link\|Display]] | Owner | Tags | Summary |
|-------------------|-------|------|---------|
```

---

## ⚡ Workflows

### Workflow A: Daily Logging
```
Context: Agent finishes a task and needs to log it
Steps:
1. Read current monthly log → Logs/YYYY-MM.md
2. Append new ## MM-DD section (or append under existing date)
3. Format: ตาม Template_Daily_Log.md
4. Include: work items, deliverables, reflection, knowledge entries
```

### Workflow B: Creating Knowledge Entry
```
Context: Agent discovers reusable knowledge (policy, architecture, strategy)
Steps:
1. Check if entry already exists (grep / glob)
2. Read Template_Knowledge.md for format
3. Create file at {team}/Knowledge/Knowledge-{topic}.md
4. Update LOCAL INDEX.md
5. Update Shared/Knowledge_INDEX.md (Central INDEX)
```

### Workflow C: Session Handoff
```
Context: Session is ending or switching context
Steps:
1. Read Handoff/SESSION.md for current state
2. Update with latest progress, decisions, pending items
3. File: Aiy/Handoff/SESSION.md
```

### Workflow D: Cross-Team Communication
```
Context: Need to send message or handoff between teams
Steps:
1. Create handoff at CoWorkspace/Handoff/YYYY/MM/DD/Handoff-{from}-{to}-{topic}.md
2. Update Kanban at CoWorkspace/Kanban/ if applicable
3. Notify target agent via task delegation or standup
```

---

## ⛔ File Creation Restrictions

| Location | Allowed For |
|----------|-------------|
| `Aiy/Knowledge/` | Knowledge entries only |
| `Aiy/Logs/` | Monthly log files only |
| `Meta/` | **READ-ONLY** — templates managed by Aiy only |
| `Shared/` | Cross-team resources — consult Aiy before writing |
| `CoWorkspace/Handoff/` | Handoff documents |
| `CoWorkspace/Standup/` | Standup updates |
| `CoWorkspace/Messages/` | Team messages |
| `CoWorkspace/Kanban/` | Kanban items |
| `{TeamName}/` | Each team's own workspace |

**Golden Rule:** DON'T create files outside your team's workspace. ถ้าต้อง cross boundary → use Handoff protocol.

---

## 📐 Dataview Queries (สำหรับ Dashboard)

Dashboard.md ใช้ Dataview — agents สามารถอ้างอิง queries ได้:

```dataview
TABLE created AS "Date", author AS "Agent"
FROM "Aiy_Workspace"
WHERE contains(file.path, "/Logs/")
SORT file.name DESC
LIMIT 10
```

```dataview
TABLE created AS "Date" 
FROM "Aiy_Workspace"
WHERE contains(file.path, "/Knowledge/") AND !contains(file.path, "Knowledge_base")
SORT file.ctime DESC
LIMIT 10
```

> **Note:** Dataview queries ทำงานใน Obsidian app เท่านั้น — agents เห็นคือ rendered markdown logic เท่านั้น

---

## 🔍 Quick Reference Commands

```bash
# ค้นหาไฟล์
glob pattern="**/*Keyword*.md" path="/home/lu5her/myObsidian/Aiy_Workspace"

# ค้นหาเนื้อหา
grep pattern="keyword" include="*.md" path="/home/lu5her/myObsidian/Aiy_Workspace"

# ค้นหาจาก tag ใน frontmatter
grep pattern="tags:.*trading" include="*.md" path="/home/lu5her/myObsidian/Aiy_Workspace"

# ดู team structure
read filePath="/home/lu5her/myObsidian/Aiy_Workspace"

# อ่าน template
read filePath="/home/lu5her/myObsidian/Aiy_Workspace/Meta/Template_Knowledge.md"

# ตรวจสอบ frontmatter (บรรทัดแรก)
read filePath="/home/lu5her/myObsidian/Aiy_Workspace/{path}/file.md" limit=5
```

---

## 🧠 Knowledge Entry Conditional Creation Rule

สรุปจาก `Template_Knowledge.md` — **CREATE ONLY IF:**
- Policy/decision ที่เปลี่ยนวิธีทำงาน
- Architecture/technical reference ที่ควรเก็บไว้
- Reusable strategy/playbook
- Lesson learned ที่ป้องกันการ重复犯错 (ซ้ำ)

**SKIP IF:**
- One-time setup/config (keep in Log only)
- Quick tip ที่ covered ใน Log แล้ว
- Duplicate of existing knowledge

---

*Last updated: 2026-07-30*
*Maintainer: Aiy (อัย)*

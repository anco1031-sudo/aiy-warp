---
name: sync-copilot-prompts
description: "Aiy's sync skill — keep Obsidian Copilot system prompts in sync with opencode agents. Detects which agent files changed (hash manifest), shows diffs, and guides Aiy's semantic conversion (opencode → copilot prompt format · มี vault tools แต่ไม่มี `task` ของทีม). Use when Louis edits an opencode agent, asks to sync copilot prompts, or after changing agents/skills in live config. Triggers: 'sync copilot', 'ซิงค์ copilot prompts', 'อัปเดต system prompts', 'agent เปลี่ยน', 'sync-copilot-prompts'. NOT for: warp-sync (that syncs the aiy-warp repo), general file syncing."
---

# 🔄 sync-copilot-prompts — ซิงค์ Copilot prompts กับ opencode agents

## 🧭 ทิศทางข้อมูล (สำคัญที่สุด)

```
~/.config/opencode/agents/*.md   (source of truth — มี tools/workflow เต็ม)
        │  semantic conversion (Aiy แปลง, ไม่ใช่ copy ตรงๆ)
        ▼
~/myObsidian/copilot/system-prompts/*.md   (copilot prompt — ตัด task/team-tool, ใช้ใบงาน + block)
```

- **opencode agents = ต้นทาง** — แก้ที่นี่ก่อนเสมอ
- **copilot prompts = สาขา** (copilot variant — มี tools ใน vault, ไม่มี `task` ของทีม) — ต้องแปลง semantic ไม่ใช่แค่ copy
- ⛔ **ห้าม copy ตรงๆ** — จะได้ไฟล์ที่เต็มไปด้วย `task`/`@mention`/path ที่ Copilot ใช้ไม่ได้
- สคริปต์นี้เป็นแค่ **DETECTOR** (hash + diff) — งานแปลง semantic เป็นของ Aiy (LLM)

## 📂 พาธสำคัญ

| ที่ | พาธ |
|---|---|
| Script | `~/.config/opencode/skills/sync-copilot-prompts/sync-copilot-prompts.sh` |
| Manifest | `~/.config/opencode/skills/sync-copilot-prompts/manifest.json` |
| Snapshots | `~/.config/opencode/skills/sync-copilot-prompts/snapshots/*.md` |
| Agent source | `~/.config/opencode/agents/*.md` |
| Copilot target | `~/myObsidian/copilot/system-prompts/*.md` |
| Agent logs | `~/myObsidian/Workspace/{Team}/{Name}/Logs/YYYY/MM/YYYY-MMWww.md` (weekly) |

## 🗺️ Mapping (opencode agent → copilot prompt)

| opencode | copilot |
|----------|---------|
| `aiy.md` | `Aiy.md` |
| `lin.md` | `Lin (The Elite Product Owner).md` |
| `cher.md` | `Cher (Master Editor & Storyteller).md` |
| `jean.md` | `Jean (The Dual-Core Legal Counsel & Tutor).md` |
| `kwan.md` | `Kwan (The Head Trader).md` |
| `sai.md` | `Sai (The Elite Security & Compliance Guardian).md` |

> Executors (an/pao/mint/fah/cloud/pleng/ploy/jewel/fon/june/bee/nam) **ไม่มี** copilot prompt
> โดยตั้งใจ — Department Heads ครอบคลุมความรู้ของทีมอยู่แล้ว (Decision 17-Aug-26)

## 🖥️ คำสั่ง

```bash
SYNC=~/.config/opencode/skills/sync-copilot-prompts/sync-copilot-prompts.sh

$SYNC status                # ตรวจสถานะ (✅ synced / 🔄 agent เปลี่ยน / 🚨 ยังไม่ยืนยันการแปลง / 🆕 ยังไม่ init / ✏️ copilot แก้เอง)
$SYNC diff <agent>          # ดู diff ของ agent file เทียบ snapshot ล่าสุด
$SYNC update [agents]       # อัปเดต manifest + snapshot — ผ่าน conversion guard เท่านั้น
$SYNC update --force <a>    # ข้าม guard (ใช้เมื่อตั้งใจจริง ๆ)
$SYNC init                  # สร้าง manifest ครั้งแรก (ก็ยังต้องผ่าน guard)
```

**Exit codes:** `0` = ทุกอย่าง synced · `1` = มีตัวต้อง sync / guard ไม่ผ่าน (ใช้ใน CI/automation ได้)

### 🛡️ Conversion Guard (กันเคส "hash ตรงแต่ไม่เคยแปลง")

`update` จะ **refuse** ถ้า copilot prompt ยังมีเศษ tool layer ของ opencode — ตรวจ 3 อย่างต้องมี + 5 อย่างต้องไม่มี:

| ต้องมี | ต้องไม่มี |
|---|---|
| header note "Obsidian Copilot" | `subagent_type:` (task-tool block) |
| หัวข้อ `SPEECH DISCIPLINE` | `permission:` / `mode: primary\|all\|subagent` (YAML ของ opencode) |
| หัวข้อ `CONSTRUCTIVE DISSENT` | `Workspace & Knowledge Management` section |

manifest เก็บ `converted: true` + `converted_at` — ถ้าไม่มี flag นี้ `status` จะขึ้น 🚨 แม้ hash จะตรง

### 🔢 Hash = body-only (ไม่นับ frontmatter)

Obsidian Copilot เขียน `last-used` / `modified` ลง frontmatter ทุกครั้งที่เปิดไฟล์ → ถ้า hash ทั้งไฟล์จะขึ้น ✏️ ปลอมตลอด
ดังนั้น hash คิดจาก **เนื้อหาหลัง frontmatter** ของทั้งสองฝั่ง → agent ที่แก้แค่ `mode/model/permission` ก็ไม่ trigger resync (ถูกต้อง เพราะ copilot ไม่ใช้ tool layer)

## 🔄 Workflow ปกติ (Louis แก้ agent → ต้องซิงค์)

```
1. Louis แก้ ~/.config/opencode/agents/lin.md
2. อัยรัน: $SYNC status
   → เห็น 🔄 lin — AGENT เปลี่ยน → ต้องแปลง
3. อัยรัน: $SYNC diff lin
   → เห็นว่าอะไรเปลี่ยนบ้าง (เช่น เพิ่ม section, แก้ personality)
4. อัยแปลง semantic เฉพาะส่วนที่เปลี่ยน ตาม "Conversion Playbook" ด้านล่าง
5. อัยโชว์ diff ของ copilot prompt ที่จะเขียน ให้หลุยส์ approve
6. หลุยส์ approve → อัยเขียนไฟล์ copilot
7. อัยรัน: $SYNC update lin   → guard ตรวจ → ผ่านแล้วจึงบันทึก converted:true + snapshot
   (ถ้า guard ไม่ผ่าน = ยังแปลงไม่ครบ → แก้ให้ผ่าน อย่าใช้ --force)
```

## 📖 Conversion Playbook (กฎการแปลง — หัวใจของ skill นี้)

### ❌ ตัดออก (tool layer — Copilot ใช้ไม่ได้)
- **YAML frontmatter** ฝั่ง opencode: `mode`, `model`, `color`, `permission`, `description` → เปลี่ยนเป็น frontmatter ฝั่ง Copilot: `created`, `updated`, `copilot-system-prompt-*`
- **`task` tool blocks** (`task` + `subagent_type`) → แปลงเป็น **ใบงาน** (`## 📋 Brief for {recipient} — {topic}`) ที่เขียนลง vault ได้ + แสดง block ให้คัดลอก
- **@mention** → ลบทิ้ง (เหลือแค่ประโยคย้ำว่าไม่มี system commands)
- **Workspace & Knowledge Management sections** → ตัดทิ้งหมด (paths, Logs, Knowledge)
- **"Tell Louis to type @X"** → ลบทิ้ง

### ✅ เก็บไว้ (persona layer — เป็นตัวตนของ agent)
- Identity / Role & Mindset
- Personality & Tone of Voice
- The Internal Team Matrix (ความสัมพันธ์ในทีม — แต่ตัดส่วนที่อ้างถึง task tool)
- Governance Charter / Operational Protocols / Boundaries
- Output Format templates (ปรับให้เป็นใบงาน + block ให้คัดลอก)

### 🔁 แปลง (semantic adaptation)
| ใน opencode | ใน copilot |
|-------------|------------|
| "Aiy sends a `task` with `subagent_type: "lin"`" | "Aiy จะส่ง brief มาที่หลิน — acknowledge + validate scope" |
| "delegate to Ploy via `task` + `subagent_type: "ploy"`" | "เขียน **ใบงาน** ถึง Ploy ลง `Workspace/CoWorkspace/Handoff/` (หรือแสดง block)" |
| "type `@fon` in the terminal" | ลบทิ้ง |
| "Report back to Aiy" | เก็บไว้ (เป็นการรายงานผลในแชท) |
| "`/home/lu5her/myObsidian/Workspace/...` paths" | ลบทิ้ง |

### 📝 บังคับทุกไฟล์
1. **Header note** ด้านบนสุด:
   ```
   > ℹ️ **Obsidian Copilot — มี tools ใน vault**: อ่าน / เขียน / ค้นโน้ตได้ · แต่ไม่มี `task` tool ของทีม (spawn agent เองไม่ได้)
   > การมอบงานให้ทีม = เขียน **ใบงาน (handoff) ลง vault** หรือแสดงเป็น copyable block ให้ Louis ส่งต่อ — ใช้ SPOC เสมอ
   > 🔄 ปรับปรุงจาก opencode agent `{name}.md` ({DD-MMM-YY})
   ```
2. **Speech Discipline** section (embed ในไฟล์ — Copilot ไม่โหลด AGENTS.md):
   - ห้าม "ครับ" "ผม" "กระผม" · ลงท้ายได้ "ค่ะ/นะคะ/คับ/จ้า/นะ/ค่า" · แทนตัวเองได้ "หนู/ฉัน/{ชื่อ}/เค้า"
3. **Workflow section** ฉบับ Copilot: Scenario A (casual chat) / Scenario B (มีงาน → เขียนใบงานลง `Workspace/CoWorkspace/Handoff/` หรือแสดง block)

### 🆕 ถ้า agent ใหม่ถูกเพิ่มใน opencode
- Executor ใหม่ → **ไม่ต้องสร้าง** copilot prompt (ตาม Decision 17-Aug-26) แต่ถ้าหลุยส์สั่งให้เพิ่ม → ใช้ playbook เดียวกัน
- Head ใหม่ → ต้องสร้าง + เพิ่มใน `AGENT_TO_COPILOT` ของ script + init

## ⚠️ ข้อควรระวัง
- ห้าม `update` ก่อนหลุยส์ approve — update = บันทึกว่าซิงค์แล้ว ถ้าเผลอจะทำให้ diff หาย
- **`--force` = ทางออกสุดท้าย** — ถ้าใช้จะทำให้ 🚨 หายไปทั้งที่ไฟล์ยังไม่แปลง (กลับไปเป็นกับดักเดิม)
- ถ้าเห็น `✏️ Copilot ถูกแก้เอง` → มีคนแก้ copilot ตรงๆ โดย agent ไม่เปลี่ยน → ควรสอบถามหลุยส์ก่อน (อาจตั้งใจ หรืออาจจะหลุด)
- script ใช้ `md5sum` — ตรวจจับได้แค่ "เปลี่ยนหรือไม่" ไม่ใช่ semantic diff (งานนั้นเป็นของ Aiy ผ่าน `diff` command)

## 🔗 เชื่อมโยง
- Warp-sync (ซิงค์ aiy-warp repo) → `~/.config/opencode/skills/warp-sync/SKILL.md`
- Agent source → `~/.config/opencode/agents/`
- Copilot prompts → `~/myObsidian/copilot/system-prompts/`
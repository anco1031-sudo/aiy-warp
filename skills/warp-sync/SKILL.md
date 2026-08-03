---
name: warp-sync
description: "Aiy's warp-kit sync skill — sync the aiy-warp repo (canonical kit) to match this master machine's live config, then commit + push so other machines can warp the latest. Use when Louis asks to update/sync the repo, push changes to aiy-warp, keep the kit in sync, or after making changes to agents/skills/playbooks in live config. Triggers: 'warp sync', 'sync repo', 'update repo', 'อัพเดท repo', 'ซิงค์ repo', 'sync kit', 'push to aiy-warp', 'sync live to repo', 'อัปเดต aiy-warp'. NOT for: installing the kit on this machine (this is the master — see Direction), general git commits, or Obsidian vault sync."
---

# 🔄 warp-sync — ซิงค์ aiy-warp repo ให้ทันเครื่องหลัก

## 🧭 ทิศทางข้อมูล (สำคัญที่สุด)

```
live (เครื่องแม่)  →  repo (canonical)  →  push → เครื่องอื่น warp
```

- **เครื่องนี้คือ MASTER** — live config (`~/.config/opencode/`) คือต้นทางของทุกการปรับปรุง
- **repo (`01-Projects/aiy-warp/`) คือ canonical kit** ที่เครื่องอื่นจะ pull ไป warp
- ⛔ **ห้ามรัน `install opencode` บนเครื่องนี้** — ไม่มีประโยชน์ (live ใหม่กว่าแล้ว) + เสี่ยงพัง (affiliate-poster มี `{{HOME}}` ที่ installer ยังไม่แทนค่า — P1 debt)
- `install` มีไว้สำหรับ **เครื่องใหม่เท่านั้น**: `git clone` → `./install/src/aiy-warp install opencode`

## 📂 พาธสำคัญ

| ที่ | พาธ |
|---|---|
| Live agents | `/home/lu5her/.config/opencode/agents/*.md` |
| Live skills | `/home/lu5her/.config/opencode/skills/*/SKILL.md` |
| Repo | `/home/lu5her/01-Projects/aiy-warp/` |
| Repo agents | `.../agents/*.md` · Repo skills | `.../skills/*/SKILL.md` |
| bundles.yaml | `.../install/bundles.yaml` |
| warp.config.example | `.../install/warp.config.example` |

## 📋 ขอบเขตการซิงค์ (checklist)

### 1. `agents/*.md` — **content merge เท่านั้น! ห้าม copy ทับ**
- repo canonical มี frontmatter พิเศษ: `warp_version, name, display_name, role, color, department, rank, reports_to, model_hint, personality, directives, platform_targets, skills` — **ต้องเก็บไว้**
- live มี body content ที่ repo ขาด → **merge content เข้า canonical** (เช่น PATH DISCIPLINE section, absolute paths)
- ตรวจทุก agent ว่ามี section/กติกาใหม่ใน live ที่ repo ไม่มี: เทียบ marker เช่น `PATH DISCIPLINE`, `Maintain these files (absolute` — ถ้า live มี repo ไม่มี → เพิ่มเข้า canonical
- frontmatter `skills:` ใน agents/*.md ต้องตรงกับ bundles.yaml

### 2. `skills/*/SKILL.md` — copy live → repo (มีข้อยกเว้น)
- ปกติ: `cp` ตรงๆ ได้เลย
- **`affiliate-poster`: ต้อง restore `{{HOME}}` parameterization** — replace `/home/lu5her/01-Projects/affiliate-content` → `{{HOME}}/01-Projects/affiliate-content` (canonical เก็บแบบ portable; installer ยังไม่แทนค่า P0)

### 3. `install/bundles.yaml` — ไล่ตาม skill list
- agent aiy: `skills: [aiy-messaging, obsidian, affiliate-poster, teach-protocol, handoff-protocol, trading-pipeline, warp-sync]`
- section `skills:` ต้องมีทุก skill ที่อยู่ใน `skills/` dir (เคยเจอ `handoff-protocol` หาย → doctor บอก 5 resolve ทั้งที่ kit มี 6)

### 4. `playbooks/*.md` + `install/warp.config.example`
- playbooks ที่แก้ใน workspace (เช่น Team_Charter) → ซิงค์เข้า repo
- warp.config.example: ถ้า path เปลี่ยน (เช่น obsidian_root) → แก้ตาม

## 🚫 ห้ามแตะ (ไม่ใช่ความผิดพลาด)

- `playbooks/Architecture_Decisions.md`, `install/DESIGN-NOTE*.md` — **บันทึกประวัติ** อ้างชื่อเก่าได้ (ตั้งใจ)
- `install/src/internal/**/*_test.go` — test fixture ค่าไม่ต้องตามยุค
- `install/src/internal/condense/condense.go` noise regex — ยังครอบคลุม (`myObsidian`)
- ไฟล์ใน `~/.cache/opencode/skills/` — plugin-managed (security-research etc.)

## ✅ Quality Gates (บังคับก่อน commit)

```bash
cd /home/lu5her/01-Projects/aiy-warp
```

1. **CJK scan** — ต้องเป็น 0 ใน repo:
   ```bash
   git grep -lP "[\x{4E00}-\x{9FFF}]" -- . || echo "✅ no CJK"
   ```
   (เจอ → แก้เป็นภาษาไทย — เคยเจอ "重复犯错" ปนใน obsidian + Team_Charter)
2. **Aiy_Workspace scan** — เหลือได้แค่บันทึกประวัติ rename:
   ```bash
   git grep -n "Aiy_Workspace" -- skills/ agents/ playbooks/ install/ | grep -v "rename จาก"
   ```
3. **Skill ตรงกัน** — ทุก skill ต้อง identical กับ live (ยกเว้น affiliate-poster {{HOME}}):
   ```bash
   for d in aiy-messaging obsidian handoff-protocol teach-protocol trading-pipeline warp-sync; do
     diff -q "/home/lu5her/.config/opencode/skills/$d/SKILL.md" "skills/$d/SKILL.md" || echo "DIFF: $d"
   done
   ```
4. **warp ตรวจ** — ต้องผ่าน:
   ```bash
   ./install/src/aiy-warp list        # 18 agents + ทุก skills resolve
   ./install/src/aiy-warp doctor      # PASS 4/5/6 (drift warning = ปกติ)
   ```

## 🚀 ขั้นตอน sync ฉบับเต็ม

```bash
cd /home/lu5her/01-Projects/aiy-warp

# 1. หา drift: เทียบ live vs repo (markers + diff)
# 2. ซิงค์ skills (copy + {{HOME}} restore สำหรับ affiliate-poster)
# 3. merge content เข้า agents/*.md (เก็บ canonical frontmatter)
# 4. อัปเดต bundles.yaml + warp.config.example ตามที่เปลี่ยน
# 5. รัน Quality Gates 4 ข้อ
# 6. commit + push:
git add -A
git commit -m "🔄 sync: warp kit ล่าสุดจากเครื่อง — <สรุปสั้นๆ สิ่งที่เปลี่ยน>"
git push origin main
# 7. Log ลง Obsidian: Workspace/Aiy/Logs/YYYY-MM.md (ระบุ commit hash)
```

## 💡 หมายเหตุ

- Commit style: emoji + ภาษาไทย/อังกฤษสั้นๆ ตามประวัติ repo (`🔄 sync:`, `🐛 fix:`, `🇹🇭 remove:`, `🎓 jean:`)
- ถ้าสงสัยว่าอะไร drift → ใช้ `diff -rq` ไล่ทีละ folder ก่อนซิงค์
- หลัง commit บอก Louis ว่า repo พร้อมให้เครื่องอื่น warp แล้ว

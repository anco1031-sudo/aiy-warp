---
status: active
deadline: 2026-08-15
goal: "CLI v1 (`aiy warp`) ที่ install/export ทีม 18 คนได้ บน opencode"
owner: lin
priority: high
tags: [project, warp, cli]
---

# 🚀 aiy-warp — Warp Kit CLI

> โปรเจกต์พกพาทีม Aiy ทั้ง 18 คน ไปเครื่อง/platform ใดก็ได้ — ตาม template `Template_Project.md`

## 🎯 Goal

`aiy warp install opencode` บนเครื่องใหม่ → ทีม 18 คนตอบได้ผ่าน @mention, `doctor` = 0

## 📌 Milestones

- [x] M0: สร้าง repo + WARP architecture + spec — 01 ส.ค.
- [x] M1: P0 build (install/export/doctor/list) — 01 ส.ค. ✅ (ปิด P0 แล้ว)
- [x] M2: P1 (web chat export, init-workspace, migrate) — 01 ส.ค. ✅ (ปิด P1 แล้ว)
- [ ] M3: P2 (sync daemon, claude adapter) — อนาคต

## ✅ Definition of Done (P0)

- [x] `aiy warp install opencode` ผ่านเครื่องใหม่ (sandbox HOME) — integration suite 30/30 ✅
- [x] `aiy warp doctor opencode` = 0 — verified in sandbox + E2E เครื่องจริง ✅
- [x] Redaction gate ทำงาน (exit 5 เมื่อเจอ secret + IDs) — fake-secret + Telegram ID test ✅
- [x] Fah test ผ่านครบ — unit 27/27 + integration 30/30 ✅
- [x] Pao code review + fixes F1-F9 — symlink-safe, frontmatter no-drop, Telegram gate, PGP regex ✅
- [x] E2E จริง: install + doctor บนเครื่องนี้ — ผ่าน (แล้วลบของออก คืนสภาพ 100%) ✅

## ✅ Definition of Done (P1)

- [x] Condensed persona engine (`internal/condense`) — single + `--collapse` (team → conductor) ✅
- [x] `export chatgpt|gemini|web` renderers — shapes match reference; redaction gate enforced ✅
- [x] `init-workspace` — PARA 6 dirs + READMEs, dry-run, noop exit 4 ✅
- [x] `migrate` — 18/18 agents stamped canonical frontmatter, `.bak` backup, idempotent (2nd run exit 4), **body zero data loss** ✅
- [x] Skill parameterization — aiy-messaging + affiliate-poster → `{{…}}`/`$AIY_*` placeholders, **0 residual real IDs** ✅
- [x] NITs F10-F12 — hostPath cross-platform, rune-safe clip, fsync-before-rename, merge fallback warn ✅
- [x] Aiy review round: FullName double-wrap + integration fixture bugs → Pao fix — unit 67 + integration **50/50** ✅
- [x] Version `0.2.0-p1`; docs updated (README, CLI_SPEC, WARP.md, DESIGN-NOTE-P1) ✅

## 📝 Notes

- **Owner:** lin (cascade: an → pao → fah)
- **Deliverables:** `install/src/` (Go module), `install/bundles.yaml`,
  `install/warp.config.example`, `install/DESIGN-NOTE.md`, `install/README.md`,
  `install/test/integration.sh`, `.gitignore` (+warp entries)
- **Spec interpretation (flag ให้ Aiy):** redaction gate = credential VALUES
  (entropy-based) ไม่ใช่ word-scan — มิฉะนั้น P0 exit criterion (doctor=0)
  เป็นไปไม่ได้กับเนื้อหา kit ปัจจุบัน (prose "token" + IDs) ดู `DESIGN-NOTE.md §2`
- **ไม่ commit** — Aiy review แล้ว commit เอง → ✅ committed `2ff6e2b` (P0 ปิด)
- Deadline อัปเดตได้ — ถ้า P0 เลื่อน ให้แก้ frontmatter (git เห็น history)
- **Governance lesson:** P0 นี้ Lin build เองแทน delegate → Pao review (fresh eyes) จับ MAJOR bugs ได้จริง — ต่อไป PO ต้อง cascade เสมอ
- **P1 backlog:** init-workspace, export chatgpt/gemini/web, migrate, skill parameterization, NITs F10-F12

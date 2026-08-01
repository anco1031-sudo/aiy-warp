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
- [x] M1: P0 build (install/export/doctor/list) — 01 ส.ค. ✅ (build+test เสร็จ, รอ Aiy review/commit)
- [ ] M2: P1 (web chat export, init-workspace) — due ส.ค. ปลาย
- [ ] M3: P2 (sync daemon, claude adapter) — อนาคต

## ✅ Definition of Done (P0)

- [x] `aiy warp install opencode` ผ่านเครื่องใหม่ (sandbox HOME) — integration suite 30/30 ✅
- [x] `aiy warp doctor opencode` = 0 — verified in sandbox ✅
- [x] Redaction gate ทำงาน (exit 5 เมื่อเจอ secret) — fake-secret test ✅
- [x] Fah test ผ่านครบ — unit + integration 30/30 ✅
- [ ] E2E จริง: @mention ทั้ง 18 คนตอบบนเครื่องใหม่ — เหลือ Louis/Aiy ยืนยัน

## 📝 Notes

- **Owner:** lin (cascade: an → pao → fah)
- **Deliverables:** `install/src/` (Go module), `install/bundles.yaml`,
  `install/warp.config.example`, `install/DESIGN-NOTE.md`, `install/README-P0.md`,
  `install/test/integration.sh`, `.gitignore` (+warp entries)
- **Spec interpretation (flag ให้ Aiy):** redaction gate = credential VALUES
  (entropy-based) ไม่ใช่ word-scan — มิฉะนั้น P0 exit criterion (doctor=0)
  เป็นไปไม่ได้กับเนื้อหา kit ปัจจุบัน (prose "token" + IDs) ดู `DESIGN-NOTE.md §2`
- **ไม่ commit** — Aiy review แล้ว commit เอง
- Deadline อัปเดตได้ — ถ้า P0 เลื่อน ให้แก้ frontmatter (git เห็น history)

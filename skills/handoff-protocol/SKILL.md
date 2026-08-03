---
name: handoff-protocol
description: "Team's standardized work-order (ใบงาน) workflow. Create, dispatch, review, and deliver handoffs using the unified Template_Handoff standard with from/to/status frontmatter and the Aiy Light-Review 3-check gate. Use when Louis assigns work, a head delegates to an executor, an executor returns results, a handoff needs review before reaching Louis, or anyone asks about work order flow. Triggers: 'ส่งงาน', 'รับงาน', 'ใบงาน', 'handoff', 'dispatch', 'delegate', 'ส่งต่อ', 'review งาน', 'ส่งมอบ', 'work order', 'สั่งงาน', 'มอบหมาย'."
---

# 🤝 Handoff Protocol — มาตรฐานใบงานทั้งทีม

Work-order (ใบงาน) workflow มาตรฐานเดียว: สร้าง → dispatch → execute → review → deliver

## 📂 ที่อยู่

| Resource | Path |
|----------|------|
| Template (ต้นแบบใบงาน) | `Workspace/CoWorkspace/Meta/Template_Handoff.md` |
| ใบงานที่ส่งแล้ว | `Workspace/CoWorkspace/Handoff/YYYY/MM/DD/Handoff-{from}-{to}-{topic}.md` |
| ประกาศทีม | `Workspace/CoWorkspace/Messages/YYYY/MM/DD/Chat.md` |
| Knowledge อ้างอิง | `Workspace/Aiy/Knowledge/Knowledge-Handoff-Standard.md` |

## 🔄 Workflow (5 ขั้น)

```
Louis ──(ใบสั่งงาน)──▶ Aiy ──(brief)──▶ Executor
                                         │
Louis ◀──(ส่งต่อ)── Aiy ◀──(verdict)── Head  ◀──(ผลงาน)──┘
                          🔎 review #2     🔎 review #1
```

| ขั้น | ใคร | ทำอะไร |
|---|---|---|
| 1. สร้างใบงาน | Aiy | copy template → ใส่ from/to/topic/status → ส่ง brief decision-complete |
| 2. Execute | Executor | ทำตาม brief → ส่งผลงาน + evidence (commit/path/test) |
| 3. Review #1 | Department Head | validate scope + quality gate → verdict accept/reject + เหตุผล |
| 4. Review #2 | Aiy | **Aiy Light-Review 3 ข้อ** → ผ่าน = ส่งต่อ / ไม่ผ่าน = ส่งกลับแก้ |
| 5. Deliver | Aiy | ส่งต่อ Louis พร้อมสรุป + flag |

## 🔎 Aiy Light-Review — 3 ข้อ (บังคับก่อนถึง Louis)

> **ไม่ผ่าน 3 ข้อ = ส่งกลับแก้ — ไม่ส่งต่อ Louis**

- [ ] **Deliverables ครบ** — ตรงตาม brief / ใบสั่งงานทุกข้อ?
- [ ] **Verdict ชัดเจน** — accept / reject + เหตุผล (งานที่ head review)?
- [ ] **Flag ครบ** — มีอะไรที่ Louis ควรเห็นก่อน (risk / deviation / blocker)?

**Verdict:** ✅ PASS / ❌ REJECT — {เหตุผลสั้น ๆ}

## 📋 Template ใบงาน (ย่อ)

```yaml
---
from: "{sender} ({ชื่อไทย})"
to: "{receiver} ({ชื่อไทย}) — {title}"
date: "YYYY-MM-DD"
topic: "{topic}"
status: "📥 from Louis → 🏗 executing → 🔎 head review → ✅ delivered"
created: "YYYY-MM-DDTHH:MM"
updated: "YYYY-MM-DDTHH:MM"
---
```

Sections: 📌 สิ่งที่ส่งต่อ · ✅ เสร็จแล้ว (พร้อม evidence) · 📋 สิ่งที่ต้องทำต่อ · 🔎 Aiy Light-Review · 🔗 อ้างอิง · 💬 หมายเหตุ

## 🌍 Portability

- **opencode:** Aiy dispatch ผ่าน `task` tool (brief decision-complete) หรือ spawn head session (`opencode run --agent <head>`) สำหรับงานที่ต้องการ autonomy ทั้ง department
- **Web chat (ChatGPT/Gemini/web):** paste ใบงานเป็น markdown block — คนอ่านคนเดียวก็ทำตาม checklist ได้
- **ทุก platform:** ใช้ template เดียวกัน → หน้าตาใบงานเหมือนกันหมด

## 🪟 Tmux Window Dispatch (lessons 03-Aug-26)

**เมื่อเปิด window ใหม่ส่งงานให้ head (`tmux new-window -t aiy -n <head>`):**

1. **สั่ง opencode แบบเปล่า ไม่มี argument** — `clear && opencode --agent <head>` (กด Enter)
2. **รอให้ TUI โหลดเสร็จ** (เห็น prompt "Ask anything...") แล้วค่อยส่งงาน
3. **ส่ง prompt เป็น ASCII ล้วน + ชี้ path ใบงาน** — ห้าม Thai/emoji/space ใน command line (quote พัง → opencode ตีความ brief เป็น directory → fail ตั้งแต่เริ่ม)
   ```bash
   tmux send-keys -t "aiy:<head>.1" "Read the handoff brief at /path/to/Handoff-xxx.md then execute it per that brief. Start now." Enter
   ```
4. **Monitor งาน** — capture-pane ดูความคืบหน้า + รอ completion marker (ไฟล์ touch) / commit

**ปิด window แบบ graceful (งานเสร็จแล้ว):**

1. `tmux send-keys -t "aiy:<head>.1" "/exit" Enter` — **ต้องกด Enter ตาม!** (/exit ค้างใน input ถ้าไม่ Enter)
2. รอ pane กลับเป็น shell (เช็ค `tmux list-panes` → command เปลี่ยนจาก opencode เป็น zsh)
3. `tmux kill-window -t aiy:<head>` — session save อัตโนมัติแล้ว ไม่ต้องกังวล
4. **⚠️ งานเสร็จ → ถาม Louis เรื่อง window นั้นก่อนปิดเสมอ** (Louis สั่ง 03-Aug-26): อย่าปิดเองโดยไม่ถาม — "ปิด window cher ไหม?"

## ⚠️ Rules

1. **คุย ≠ ส่งงาน (Aiy ↔ Head Protocol 2026-08-02):** คุย/สอบถาม/ติดตาม → `task` + `task_id` (subagent — ส่งต่อไม่ได้ ไม่จำเป็นต้องใช้) · **งานจริงที่ต้องใช้ทีม → `opencode run --agent <head> "brief"`** (primary — หัวหน้า cascade ลูกทีมเองได้) · window `<head>` ใช้สำหรับหลุยส์ดูสด/คุยเองเท่านั้น (ห้าม script ส่ง keys เข้า TUI)
2. **Executor (ลูกทีม) เป็น leaf worker** — spawn งานต่อไม่ได้ (subagent constraint, opencode #8114) — ถ้างานต้อง cascade ให้ส่งผ่าน head
3. **Work order ต้อง decision-complete** — strategic intent, scope, deliverables, acceptance criteria — หัวหน้าไม่ต้องเดา
4. **`--dry-run` ก่อน destructive ops** — ใช้กับงานที่แตะของจริง
5. **Sandbox HOME เสมอ** ตอนทดสอบ — ห้ามแตะ `~/.config/` จริง
6. **Respect redaction gate** (exit 5) — ไม่ bypass `--allow-identifiers` เว้น Louis อนุญาต

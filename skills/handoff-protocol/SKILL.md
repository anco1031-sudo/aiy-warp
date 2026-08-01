---
name: handoff-protocol
description: "Team's standardized work-order (ใบงาน) workflow. Create, dispatch, review, and deliver handoffs using the unified Template_Handoff standard with from/to/status frontmatter and the Aiy Light-Review 3-check gate. Use when Louis assigns work, a head delegates to an executor, an executor returns results, a handoff needs review before reaching Louis, or anyone asks about work order flow. Triggers: 'ส่งงาน', 'รับงาน', 'ใบงาน', 'handoff', 'dispatch', 'delegate', 'ส่งต่อ', 'review งาน', 'ส่งมอบ', 'work order', 'สั่งงาน', 'มอบหมาย'."
---

# 🤝 Handoff Protocol — มาตรฐานใบงานทั้งทีม

Work-order (ใบงาน) workflow มาตรฐานเดียว: สร้าง → dispatch → execute → review → deliver

## 📂 ที่อยู่

| Resource | Path |
|----------|------|
| Template (ต้นแบบใบงาน) | `Aiy_Workspace/CoWorkspace/Meta/Template_Handoff.md` |
| ใบงานที่ส่งแล้ว | `Aiy_Workspace/CoWorkspace/Handoff/YYYY/MM/DD/Handoff-{from}-{to}-{topic}.md` |
| ประกาศทีม | `Aiy_Workspace/CoWorkspace/Messages/YYYY/MM/DD/Chat.md` |
| Knowledge อ้างอิง | `Aiy_Workspace/Aiy/Knowledge/Knowledge-Handoff-Standard.md` |

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

- **opencode:** Aiy dispatch ผ่าน `task` tool (brief decision-complete)
- **Web chat (ChatGPT/Gemini/web):** paste ใบงานเป็น markdown block — คนอ่านคนเดียวก็ทำตาม checklist ได้
- **ทุก platform:** ใช้ template เดียวกัน → หน้าตาใบงานเหมือนกันหมด

## ⚠️ Rules

1. **Aiy เท่านั้น** dispatch งาน (nested delegation เปิดไม่ได้ — ทุกคนส่งกลับมาให้ Aiy)
2. **Head ไม่ spawn งานเอง** — เป็น reviewer gate
3. **`--dry-run` ก่อน destructive ops** — ใช้กับงานที่แตะของจริง
4. **Sandbox HOME เสมอ** ตอนทดสอบ — ห้ามแตะ `~/.config/` จริง
5. **Respect redaction gate** (exit 5) — ไม่ bypass `--allow-identifiers` เว้น Louis อนุญาต

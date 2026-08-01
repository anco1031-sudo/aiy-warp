---
created: Sun 12,Jul 26 17:21
updated: Thu 16,Jul 26 21:43
---

# 🏗️ Architecture Decision Records (ADR)

นี่คือพื้นที่บันทึกการตัดสินใจทางสถาปัตยกรรมที่สำคัญของทีม

## ADR Template
- **สถานะ**: เสนอ / อนุมัติ / ปรับปรุง / ยกเลิก
- **วันที่**: 
- **ผู้เสนอ**: 
- **เหตุผล**: 
- **ผลกระทบ**: 

---

## ADR-001: Direct Subagent Delegation (ยกเลิก Proxy Pattern)

- **สถานะ**: อนุมัติ ✅
- **วันที่**: 2026-07-12
- **ผู้เสนอ**: Aiy (อัย)
- **อนุมัติโดย**: Louis (หลุยส์)
- **เหตุผล**: 
  - เดิมใช้ `task general` + role-play เป็น proxy เพื่อเรียก subagent (เพราะตอนนั้น agent definitions ยังไม่เสร็จ)
  - หลังจากมี agent definitions ครบ 18 ตัว → เปลี่ยนเป็น `task` + `subagent_type` โดยตรง
  - ได้ identity จริงของแต่ละคน, prompt สั้นลง, permission/model แยกกันได้
- **ผลกระทบ**:
  - ✅ ไม่ต้อง role-play ทุกครั้งที่เรียก
  - ✅ Agent มี system prompt จริง ๆ ใช้
  - ✅ Permission/model แยกกันได้
  - ✅ Prompt สั้นลงมาก
  - ✅ Louis @mention ได้โดยตรง
  - 🔄 ต้องอัปเดต agent definitions ทุกตัว (aiy, lin, kwan, cher, sai — เสร็จแล้ว)
- **ไฟล์ที่แก้ไข**:
  - `~/.config/opencode/agents/aiy.md` — Routing table + workflow
  - `~/.config/opencode/agents/lin.md` — Delegation pattern
  - `~/.config/opencode/agents/kwan.md` — Delegation pattern
  - `~/.config/opencode/agents/cher.md` — Delegation pattern
  - `~/.config/opencode/agents/sai.md` — Delegation pattern

## ADR-002: ทีม Trader (Kwan, Fon, June, Bee, Nam)

- **สถานะ**: อนุมัติ ✅
- **วันที่**: 2026-07-12
- **ผู้เสนอ**: Louis (หลุยส์)
- **เหตุผล**: เพิ่มทีมเทรดหญิงเก่ง 5 คน เพื่อดูแลกลยุทธ์การลงทุน
- **ผลกระทบ**:
  - ✅ สมาชิกทีมรวม 18 คน (5 Department Heads + 13 Executors)
  - ✅ Obsidian system prompts ที่ `copilot/system-prompts/`
  - ✅ OpenCode agent definitions ที่ `~/.config/opencode/agents/`
  - ✅ Knowledge_base สำหรับกฎการทำงานของแต่ละคน

## ADR-003: Team Directory Restructure

- **สถานะ**: อนุมัติ ✅
- **วันที่**: 2026-07-12
- **ผู้เสนอ**: Louis (หลุยส์)
- **เหตุผล**: จัดระเบียบ `Aiy_Workspace` ตามทีมให้เป็นระบบเดียวกันกับ Trader_Team
- **โครงสร้างใหม่**:
  ```
  Aiy_Workspace/
  ├── Aiy/                  — Strategic Orchestrator
  ├── Product_Team/         — Lin, An, Pao, Mint, Fah, Cloud
  ├── Security_Team/        — Sai, Pleng
  ├── Legal_Team/           — Jean
  ├── Content_Team/         — Cher, Ploy, Jewel
  ├── Trader_Team/          — Kwan, Fon, June, Bee, Nam
  ├── CoWorkspace/
  ├── Meta/
  ├── Shared/
  └── Dashboard.md
  ```
- **ผลกระทบ**:
  - ✅ workspace เป็นระเบียบ ดูง่าย
  - ✅ path references ใน agent definitions 12 ไฟล์ถูกอัปเดตแล้ว
  - ✅ Louis หา folder น้อง ๆ ได้ง่ายขึ้น
  - 🔄 ต้องแจ้งทุกคนว่าทีมตัวเองย้ายไปไหน (ทำแล้วใน Index + Team_Charter)

## ADR-004: External Specialists Integration (oh-my-openagent Plugin)

- **สถานะ**: อนุมัติ ✅
- **วันที่**: 2026-07-14
- **ผู้เสนอ**: Louis (หลุยส์)
- **ดำเนินการโดย**: Aiy (อัย)
- **เหตุผล**: 
  - ติดตั้ง oh-my-openagent plugin เพื่อเพิ่มความสามารถด้าน software engineering
  - Agents ของ plugin (Sisyphus, Oracle, etc.) เป็น Plugin Agents — ไม่ใช่ native subagent_type
  - ไม่สามารถเรียกผ่าน `task` + `subagent_type` ได้โดยตรง
  - ต้องใช้ @mention หรือ magic word `ultrawork`/`ulw` แทน
- **ผลกระทบ**:
  - ✅ Core Team (18 คน) — เรียกผ่าน `task` + `subagent_type` ได้ปกติ
  - ✅ External Specialists (11 คน) — เรียกผ่าน @mention หรือ `ultrawork`
  - ✅ Aiy รู้จักทั้งสองทีม และ route งานให้ถูกต้อง
  - ✅ Team Charter อัปเดต — เพิ่ม External Specialists section
  - ✅ Knowledge Base สร้าง — `Knowledge-oh-my-openagent-team.md`
  - ✅ ตัวตน Aiy & น้อง ๆ ไม่ได้รับผลกระทบ — default_agent ยังเป็น "aiy"
  - ⚠️ External Specialists ใช้ tokens แยกต่างหาก — ต้อง monitor resource usage

## ADR-005: Token Optimization — Workspace Restructure

- **สถานะ**: อนุมัติ ✅
- **วันที่**: 2026-07-16
- **ผู้เสนอ**: Louis (หลุยส์)
- **ดำเนินการโดย**: Aiy (อัย)
- **เหตุผล**: 
  - ลด token waste จาก boilerplate, frontmatter ซ้ำซ้อน, และ content overlap ระหว่าง Log/Reflection/Knowledge
  - ประหยัด ~76% ของ token ที่ใช้กับการบันทึกข้อมูล (จาก ~1,920 → ~463 token/วัน/คน)
  - ถ้า rollout ทุกคนในทีม (18 คน) → ประหยัด ~26,000 token/วัน หรือ ~780,000 token/เดือน
- **การเปลี่ยนแปลง**:
  | Before | After | เหตุผล |
  |---|---|---|
  | `Daily_logs/YYYY/MM/DD/Log.md` | `Logs/YYYY-MM.md` (rolling monthly) | ลด frontmatter, folder nesting |
  | `Reflection/YYYY/MM/DD/Reflection.md` | merged into Log section | content overlap 40% |
  | `Knowledge_base/YYYY/MM/DD/Knowledge-{name}.md` | `Knowledge/Knowledge-{name}.md` (flat) | ลด path nesting, conditional creation |
  | Frontmatter 10 บรรทัด | Frontmatter 2 บรรทัด (`date:` + `type:`) | ไร้ boilerplate |
  | Emoji-heavy headers | Plain text headers | ลด token |
  | Path annotation ทุกไฟล์ | ไม่มี | path = ที่อยู่ไฟล์ |
- **ผลกระทบ**:
  - ✅ Token savings ~76% ต่อการบันทึก 1 วัน
  - ✅ Team-wide rollout → setup `Logs/` + `Knowledge/` สำหรับทุกคน
  - ✅ Meta templates อัปเดตครบ
  - ✅ Old files เก็บไว้เป็น history (ไมเกรต)
  - 🔄 Dept Heads ต้อง cascade กฎใหม่ให้ลูกทีม
- **ไฟล์ที่เกี่ยวข้อง**:
  - `Meta/Template_Daily_Log.md` — unified rolling-monthly format
  - `Meta/Template_Reflection.md` — archived
  - `Meta/Template_Knowledge.md` — stripped + conditional rule
  - `Meta/Index.md` — updated structure
  - `Shared/Team_Charter.md` — rules 8-10 updated
  - `Shared/Architecture_Decisions.md` — ADR-005 (นี้)

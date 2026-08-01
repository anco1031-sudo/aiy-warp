---
template: project
created: YYYY-MM-DD
---

# 🚀 {Project Name}

---
status: active          # active | paused | done | archived
deadline: YYYY-MM-DD    # เป้าหมายส่งมอบ (เปลี่ยนแล้ว git diff เห็น history)
goal: "เป้าหมายของโปรเจกต์ 1-2 บรรทัด"
owner: {lin|kwan|cher|aiy|...}   # ใครรับผิดชอบหลัก (routing ให้ agent ได้)
priority: med            # high | med | low
tags: [project, {category}]
---

## 🎯 Goal

> อะไรคือความสำเร็จของโปรเจกต์นี้ — เขียนให้ชัดเจน วัดผลได้

## 📌 Milestones

- [ ] M1: {คำอธิบาย}  — due {date}
- [ ] M2: {คำอธิบาย}  — due {date}
- [ ] M3: {คำอธิบาย}  — due {date}

## ✅ Definition of Done

- [ ] {เงื่อนไขการส่งมอบที่วัดผลได้}
- [ ] {เงื่อนไขที่ 2}
- [ ] {ทดสอบ/ตรวจสอบผ่าน}

## 📝 Notes

- {สิ่งที่ต้องรู้ / ข้อจำกัด / การตัดสินใจที่สำคัญ}

---

## 🤖 สำหรับทีม agents (อ่านอัตโนมัติ)

- **Deadline อยู่ใน frontmatter** — อ่านได้ด้วยเครื่อง (`status`, `deadline`, `owner`, `priority`)
- เมื่อ deadline เหลือ ≤ 3 วัน → แจ้งเตือน Louis ผ่าน Discord/Telegram (ใช้ skill `aiy-messaging`)
- อัปเดต `status` + `deadline` ทุกครั้งที่มีการเปลี่ยนแปลง — git จะบันทึก history การเลื่อน
- `owner` คือตัวบอกว่าใครรับผิดชอบ — ถ้าไม่แน่ใจให้ถาม Aiy ก่อน

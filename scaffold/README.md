# 🗂️ PARA Workspace Template (Filesystem Root)

> **โครงสร้างพื้นฐาน PARA ที่ root ของ filesystem** — ตามที่เครื่องนี้ใช้อยู่ (`~/00-Inbox`, `~/01-Projects`, …)
> ใช้ร่วมกับ `aiy warp init-workspace` เพื่อสร้าง "บ้านระบบ" บนเครื่องใหม่ก่อนติดตั้งทีม

## โครงสร้าง

```
$HOME/
├── 00-Inbox/        ← กล่องจดหมาย — ทุกอย่างที่ยังไม่ได้จัดเรียง (ของดิบ)
├── 01-Projects/     ← โปรเจกต์ที่กำลังทำ — มีเป้าหมาย + deadline ชัดเจน
├── 02-Areas/        ← พื้นที่รับผิดชอบระยะยาว — มาตรฐานที่ต้องรักษา (ไม่มีวันจบ)
├── 03-Resources/    ← ความรู้/แหล่งอ้างอิง — แชร์ได้, สาธารณะได้
├── 04-Archives/     ← เก็บถาวร — จบแล้ว/ปิดแล้ว
└── 05-Other/        ← อื่น ๆ ที่ไม่เข้า PARA (เครื่องมือ, misc)
```

## หลักการ PARA (จาก note ของ Louis)

| ประเภท | คำถามที่ใช้จัดเรียง | ลักษณะ |
|---|---|---|
| **Project** | มี deadline ไหม? มีเป้าหมายไหม? | ระยะสั้น, จบได้, แบ่งย่อยได้ |
| **Area** | เป็นบทบาท/หน้าที่ที่ต้องรักษามาตรฐานไหม? | ไม่มีวันสิ้นสุด |
| **Resource** | เป็นความรู้/ความสนใจไหม? | แชร์เป็น public ได้ |
| **Archive** | จบแล้ว/ไม่ได้ใช้แล้วไหม? | เก็บไว้เฉย ๆ |

**Flow:** ทุกอย่างเข้า `00-Inbox` ก่อน → จัดเรียงตามคำถามข้างบน → เข้า Project/Area/Resource → จบแล้วย้ายไป Archive

## วิธีใช้

```bash
# ผ่าน CLI (P1 — ยังไม่ build):
aiy warp init-workspace --home $HOME

# หรือ manual:
mkdir -p ~/00-Inbox ~/01-Projects ~/02-Areas ~/03-Resources ~/04-Archives ~/05-Other
cp -r scaffold/PARA-template/*/README.md ~/00-Inbox/  # copy คำอธิบายแต่ละโฟลเดอร์
```

## หมายเหตุ

- โฟลเดอร์นี้เป็น **template tree** — CLI จะ copy ไปสร้างที่ `$HOME` (หรือ `--home` ที่กำหนด)
- แต่ละโฟลเดอร์ย่อยมี `README.md` อธิบายการใช้งาน — ทำหน้าที่เป็นทั้งคำอธิบายและ placeholder (.gitkeep)
- เชื่อมโยงกับ AGENTS.md convention: project deliverables → `~/01-Projects/{ProjectName}/`

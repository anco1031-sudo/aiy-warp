---
name: office-suite
description: "Aiy's Office Suite — generate Excel/Word/PowerPoint from JSON via openpyxl/python-docx/python-pptx (Aetox 4 pattern). Use when Louis asks create Excel, Word, PowerPoint, sheet, doc, slide, report. Triggers: 'excel', 'sheet', 'word', 'doc', 'pptx', 'slide', 'สร้างรายงาน', 'office'."
---

# 📊 office-suite — สร้าง Excel/Word/PPTX จาก JSON

> Pattern จาก Aetox 4 — Model ส่ง JSON → python สร้าง OOXML

## เมื่อไหร่ใช้
- สร้างรายงาน Excel/Word/PPTX แบบไม่ต้องรู้ OOXML
- Agent สร้างไฟล์ให้ Louis ดาวน์โหลด

## สคริปต์
| สคริปต์ | ทำอะไร | JSON schema |
|---|---|---|
| `scripts/office/sheet_write.py` | Excel | `{"path":"report.xlsx","sheets":[{"name":"Sales","columns":[],"formats":[],"rows":[]}]}` |
| `scripts/office/doc_write.py` | Word | `{"path":"report.docx","blocks":[{"type":"heading","text":""}]}` |
| `scripts/office/slide_write.py` | PPTX | `{"path":"deck.pptx","slides":[{"title":"","content":""}]}` |

## วิธีใช้
```bash
# Excel
python3 scripts/office/sheet_write.py /tmp/input.json --out report.xlsx
echo '{"path":"report.xlsx","sheets":[{"name":"Sales","columns":["Month","Revenue"],"rows":[["Jan",100000]]}]}' | python3 scripts/office/sheet_write.py - --out report.xlsx

# Word — blocks: heading/paragraph/bullets/numbered/table/lineitems/image
python3 scripts/office/doc_write.py /tmp/doc.json --out report.docx

# PPTX
python3 scripts/office/slide_write.py /tmp/slides.json --out deck.pptx
```

## Cell formats (sheet)
`formats: ["", "money", "percent"]` → `#,##0.00` / `0.00%`, รองรับ formula `=SUM(B2:B13)`

## Dependencies
`pip install openpyxl python-docx python-pptx` — รัน `scripts/warp-init.sh --check-deps` เพื่อเช็ค

## ทิป
- ส่ง path เป็น absolute เสมอ (`/home/lu5her/01-Projects/.../report.xlsx`)
- สร้างแล้วบอก Louis ว่าไฟล์อยู่ที่ไหน + ให้ปุ่มเปิดไฟล์

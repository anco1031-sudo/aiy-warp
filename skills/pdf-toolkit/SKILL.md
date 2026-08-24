---
name: pdf-toolkit
description: "Aiy's PDF Toolkit — extract text/layout from PDFs via poppler pdftotext + pymupdf/pdfplumber fallback. Wraps scripts/pdf/read.py. Use when Louis asks read PDF, extract PDF text, parse PDF. Triggers: 'pdf', 'read pdf', 'extract pdf', 'อ่าน pdf'."
---

# 📄 pdf-toolkit — อ่าน PDF (poppler)

> Pattern จาก Aetox 3.4 — `pdftotext -layout` + fallback

## เมื่อไหร่ใช้
- อ่านเนื้อหา PDF แบบคง layout
- ดึงข้อความจากเอกสารสแกนที่ไม่ใช่รูป (text-based PDF)

## สคริปต์
`scripts/pdf/read.py` — ลอง `pdftotext -layout` ก่อน, ถ้าไม่มีหรือไม่สำเร็จ fallback เป็น `pymupdf` → `pdfplumber`

## วิธีใช้
```bash
python3 scripts/pdf/read.py /path/to/doc.pdf
python3 scripts/pdf/read.py /path/to/doc.pdf --json
```

## Dependencies
`poppler-utils` (`pdftotext`) — `sudo apt install poppler-utils` / `brew install poppler`
Fallback มีแล้ว: `pymupdf`, `pdfplumber`, `pillow`

## ทิป
- ถ้า PDF เป็นสแกนรูป → ใช้ `media-ocr` แทน (ต้อง OCR)
- อยากได้ layout เป๊ะ → `pdftotext -layout` ดีกว่า fallback

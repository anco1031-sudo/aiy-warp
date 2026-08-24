---
name: media-ocr
description: "Aiy's Media OCR skill — extract text from images & video frames via Tesseract (tha+eng). Wraps scripts/media/ocr.py + video_ocr.py. Use when Louis asks OCR, extract text from image/screenshot, video OCR, read Thai text in picture. Triggers: 'ocr', 'extract text', 'read image', 'video ocr', 'สแกนรูป', 'อ่านข้อความในรูป'."
---

# 🔍 media-ocr — OCR รูป & วิดีโอ (Tesseract tha+eng)

> Pattern จาก Aetox 3.1 + 3.3 — shell out to `tesseract`, parse TSV confidence

## เมื่อไหร่ใช้
- Louis ส่งรูป/สกรีนช็อตแล้วอยากได้ข้อความข้างใน
- วิดีโอที่มีตัวอักษรบนจอ อยากดึงออกมาทุก 5 วิ
- เอกสารสแกนภาษาไทย+อังกฤษ

## สคริปต์
| สคริปต์ | ทำอะไร |
|---|---|
| `scripts/media/ocr.py` | `tesseract <img> -l tha+eng txt tsv` → text + confidence |
| `scripts/media/video_ocr.py` | `ffmpeg -vf fps=1/5` → วนเรียก ocr.py ทีละ frame |

## วิธีใช้

```bash
# รูปเดียว
python3 scripts/media/ocr.py /path/to/image.jpg --langs tha+eng
python3 scripts/media/ocr.py /path/to/image.jpg --json

# วิดีโอ (ทุก 5 วิ)
python3 scripts/media/video_ocr.py /path/to/video.mp4 --interval 5 --json
```

## Output
```json
{"text": "สวัสดีครับ", "confidence": 92.3, "warning": ""}
# confidence <70% → warning: text may not match image (Tesseract อาจเดาภาษาผิด)
```

## Dependencies
`tesseract 5.5+` + `tessdata tha+eng` → `tesseract --list-langs` ต้องเห็น `tha` `eng`
ติดตั้ง: `sudo apt install tesseract-ocr tesseract-ocr-tha` / `brew install tesseract tesseract-lang`

## ทิปสำหรับ agent
- ถ้า confidence ต่ำ ให้บอก Louis ว่าผลอาจไม่แม่น
- รูปเบลอ → แนะนำให้ crop/เพิ่ม contrast ก่อน ocr

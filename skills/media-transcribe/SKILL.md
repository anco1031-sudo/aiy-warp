---
name: media-transcribe
description: "Aiy's Audio Transcription skill — transcribe audio/video via ffmpeg + Whisper (offline). Wraps scripts/media/transcribe.py. Use when Louis asks transcribe audio,ถอดเสียง, speech-to-text, audio to text, video to text. Triggers: 'transcribe', 'ถอดเสียง', 'audio', 'whisper', 'speech to text'."
---

# 🎙️ media-transcribe — ถอดเสียง Audio/Video (Whisper)

> Pattern จาก Aetox 3.2 — ffmpeg → 16kHz mono WAV → whisper

## เมื่อไหร่ใช้
- ไฟล์เสียง mp3/wav/m4a อยากได้ transcript
- วิดีโอ mp4 อยากได้เสียงพูดเป็น `[m:ss] text`

## สคริปต์
`scripts/media/transcribe.py` — ขั้นตอน: `ffmpeg -vn -ar 16000 -ac 1` → `whisper.load_model("base").transcribe()`

## วิธีใช้
```bash
python3 scripts/media/transcribe.py /path/to/audio.mp3 --model base
python3 scripts/media/transcribe.py /path/to/video.mp4 --model base --json
# model: tiny/base/small/medium/large — base สมดุลเร็ว/แม่น
```

## Output
```
[0:00] สวัสดีครับ วันนี้เราจะ...
[0:07] เนื้อหาต่อไปคือ...
```

## Dependencies
`ffmpeg` + `pip install openai-whisper` — whisper โหลด model ครั้งแรกจะดาวน์โหลด ~150MB

## ทิป
- ไฟล์ยาว → ใช้ `base` หรือ `small` พอ, `large` ช้าแต่แม่นกว่า
- ถ้าเสียงไทยปนอังกฤษ whisper auto-detect ได้ ไม่ต้องระบุภาษา

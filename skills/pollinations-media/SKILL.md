---
name: pollinations-media
description: "สร้างสื่อจริงด้วย Pollinations API — ภาพ (text-to-image), แก้/แต่งภาพ, วิดีโอ, เสียงพูด (TTS), เพลง/SFX, โมเดล 3D และถอดเสียงเป็นข้อความ แล้วเซฟไฟล์ลงเครื่อง. Use whenever Louis wants any media produced: 'สร้างภาพ', 'วาดรูป', 'ทำภาพประกอบ', 'ทำภาพปก', 'แก้ภาพ', 'แต่งภาพ', 'เปลี่ยนสไตล์ภาพ', 'สร้างวิดีโอ', 'ทำคลิป', 'ทำเสียงพูด', 'พากย์เสียง', 'ทำเพลง', 'ทำเสียงประกอบ', 'สร้างโมเดล 3D', 'ถอดเสียง', 'transcribe', 'generate an image/video/audio', 'text to speech', 'tts'. Also use when he needs a specific model (flux / nano-banana / gpt-image / seedream / ideogram / veo / seedance), a specific size or aspect, reproducible seeds, or a batch saved to disk — even without saying the word 'media' (e.g. 'ขอภาพประกอบโพสต์นี้ 3 แบบ', 'ทำคลิปสั้น 5 วิ')."
---

# 🎬 Pollinations Media — ภาพ · วิดีโอ · เสียง · 3D จาก prompt เดียว

> ทุกอย่างที่ skill นี้สร้าง **เป็นไฟล์จริงบนเครื่อง** ไม่ใช่ลิงก์ลอย ๆ
> หลังสร้างเสร็จ **ต้องเปิดดู/ฟังไฟล์จริงก่อนส่งให้หลุยส์เสมอ** (`read` กับภาพ) — ถ้าไม่ตรงโจทย์ ให้แก้ prompt ทำใหม่ อย่าส่งชิ้นแรกแบบไม่ตรวจ

## ⚡ เริ่มใน 3 บรรทัด

```bash
M=~/.config/opencode/scripts/pollinations_media.py
python3 $M doctor                       # key + สิทธิ์ + โฟลเดอร์ + connectivity (ครั้งแรกของ session)
python3 $M image "a bowl of ramen, steam rising, moody side light, 50mm" -a 4:3
```

ผลลัพธ์: `saved : /home/lu5her/00-Inbox/Aiy/images/2026-11-19-1249-...jpg`
→ เปิดดูด้วย `read` → ส่งให้หลุยส์พร้อม path + model + seed

## 🔑 Key & สิทธิ์ (ตั้งไว้แล้ว 19-Nov-26)

| ที่ | ค่า |
|---|---|
| config | `~/.config/opencode/secrets/pollinations.json` (`api_key`, `default_model`) |
| ทางเลือก | env `POLLINATIONS_API_KEY` หรือ `--key sk_...` |
| ตรวจ | `python3 $M doctor` · `python3 $M models -t image` |

- **ไม่มี key = สร้างไม่ได้** (Pollinations ปิด anonymous tier แล้ว — ได้ 401)
- ห้าม copy key ลง skill/log/commit (ไฟล์ config อยู่นอก `skills/` โดยเจตนา เพราะ `skills/` ถูกซิงค์เข้า repo)

### 💰 Quest balance vs Paid balance — ตัวตัดสินว่าอะไรใช้ได้ (ตรวจจริง 19-Nov-26)

บัญชีนี้มี **Quest (ฟรี) เท่านั้น** → model พรีเมียมตอบ `402 ... requires the paid balance`
**ตัวเลขจริงจาก `/models`:** ภาพ 84 ตัว (Quest 18) · วิดีโอ 25 ตัว (Quest 0) · 3D 4 ตัว (Quest 0) · เสียง 43 ตัว (Quest 9)

| ✅ ใช้ได้ด้วย Quest | ❌ ต้องเติม paid balance |
|---|---|
| **ภาพ (18):** `flux.1-schnell` · `flux.1.1-pro` · `flux.2-flex` · `flux.2-klein-4b` · `flux.1-kontext-pro` · `z-image-turbo` · `mai-image-2.6`/`2.7` (+flash, ต้อง ≥768px) · `community/*` (sdxl-lightning-free, flux-2-dev-free, lucid-origin, phoenix-1.0, grok-imagine-2-quest) | **ภาพ:** `gemini-3.1-flash-image` (nano-banana) · `gpt-image-*` · `seedream-*` · `qwen-image-*` · `recraft-*` · `ideogram-*` · `flux-3-*` |
| **เสียงพูด:** `openai/tts-1` · `openai/tts-1-hd` | **เพลง/SFX:** ทุกตัว (`elevenlabs/*`, `stable-audio-*`, `lyria-*`) |
| **ถอดเสียง:** `openai/gpt-transcribe` · `whisper-large-v3` | **วิดีโอ:** ทุกตัว (25/25) |
| | **3D:** ทุกตัว (4/4) |

- ตารางนี้อัปเดตได้เองจาก field `paid_only` ใน API: `python3 $M models -t image --quest` (ไม่ต้องยิงทดสอบ ไม่เสียเงิน) · `probe -n 12` ใช้เมื่อสงสัยว่าตารางล้าสมัย (ตัวที่พังไม่เสียเงิน)
- เจอ `402` อย่าถือว่าบั๊ก: เปลี่ยนไปตัวที่อยู่ในคอลัมน์ซ้าย หรือบอกหลุยส์ว่า "งานนี้ต้องเติม paid balance ~$X"

### 🧾 Pollen ไม่ unlimited — โควต้าคืออะไร (ตรวจจริง 19-Nov-26)

`pollen` = หน่วยเงินของ Pollinations · บัญชีมี **2 bucket** (ยืนยันจาก OpenAPI: `accountBalance: {total, tier, paid}`)

| bucket | ที่มา | ใช้กับ |
|---|---|---|
| **`tier`** (Quest/ฟรี) | เควสต์ (`/quests/catalog`) — เช่น `first_api_key` +$0.25 · `use_video_model` +$1 · `invite_friends` +$2 · `merged_pr` +$5 | model ที่ `paid_only: false` |
| **`paid`** | เติมเงิน (`/top-up`) | model ที่ `paid_only: true` (วิดีโอ/เพลง/3D/nano-banana/gpt-image) |

- **ฟรี ≠ ไม่จำกัด** — บัญชีนี้เริ่มด้วย ~$0.25 จากเควสต์ และใช้หมดไปแล้วในวันเดียว (27 คำขอ = $0.2502 · `flux.1.1-pro` กินไป $0.12 คนเดียว)
- ดูยอดใช้จริงได้เสมอ: `python3 $M usage` (อ่าน `/account/key/usage` — ไม่ต้องมีสิทธิ์พิเศษ)
- **ยอดคงเหลือ** ต้องมีสิทธิ์ `account:usage` → `python3 $M balance` · ถ้ายังไม่ให้สิทธิ์ ให้อ่านตัวเลขจากข้อความ `402` แทน (มันบอก "available balance" ตรง ๆ)
- เติมฟรีเพิ่ม: ทำเควสต์ที่ `https://enter.pollinations.ai/quests` (API อ่านรายการได้ที่ `/quests/catalog` แต่ **กดรับต้องทำในหน้าเว็บ**)
- เควสต์ที่ได้เงินแต่ยังไม่รับ และทำได้เลย: `use_image_model` / `use_text_model` / `use_audio_model` (+$0.25 each) · `join_discord` (+$1) · `use_video_model` (+$1 — ต้องมี paid ก่อน) · `invite_friends` (+$2) · `reported_merged_issue` (+$2)

**ประมาณการราคาต่อภาพ ( Quest):** `flux.1-schnell` $0.002 · `z-image-turbo` $0.004 · `flux.2-klein-4b` $0.005 · `flux.1.1-pro` $0.030 → **งบ $0.25 = schnell ~125 ภาพ หรือ flux.1.1-pro แค่ ~8 ภาพ**

### 🆓 นโยบายของหลุยส์: "เน้นฟรีไว้ก่อน" (19-Nov-26)**ทุกคำสั่งรันด้วย `--free` เป็นค่าเริ่มต้นของงานปกติ** — CLI จะเช็ค `paid_only` ก่อนยิง และปฏิเสธทันทีถ้า model ต้องเติมเงิน (ไม่เสีย request ไม่ต้องรอ 402)

```bash
python3 $M image "prompt" --free            # ✅ ผ่านถ้าเป็น Quest
python3 $M image "prompt" --free -m openai/gpt-image-1-mini
# ✖ ปฏิเสธก่อนยิง + บอกทางเลือกรุ่นฟรี
```

- จะใช้ model paid ต้อง**เอา `--free` ออกโดยเจตนา** + บอกราคาหลุยส์ก่อน
- งานที่ต้องใช้ของแพงจริง ๆ (ตัวหนังสือในภาพ, วิดีโอ, เพลง, 3D) → ถามก่อนทุกครั้ง ไม่เดา

## 🧠 เลือก model ให้ตรงงาน

| งาน | ใช้ | หมายเหตุ |
|---|---|---|
| ร่างไอเดีย/ลอง composition | `flux.1-schnell` | ถูกสุด (0.002/token) เร็ว วนได้ 5-10 แบบ |
| ภาพประกอบคุณภาพ (Quest) | `flux.1.1-pro` | คมกว่า schnell ชัดเจน — ตัวหลักของงานจริง |
| ภาพมีตัวหนังสือ/โลโก้ | `ideogram-*`, `gemini-3-pro-image` | paid — ถ้าจำเป็นต้องขออนุมัติหลุยส์ |
| แก้ภาพ/คงหน้าคน/เปลี่ยนแสง | `flux.2-klein-4b` (default ของ `edit`) · `flux.1-kontext-pro` · `flux.2-flex` | Quest ทั้งหมด — รับ `--image` |
| ภาพสินค้า/mockup เป๊ะตามสั่ง | `gpt-image-2` | paid |
| ภาพแนวอนิเมะ/ภาพวาด | `z-image-turbo`, `sdxl-lightning-free` | Quest |
| วิดีโอ | `gemini-omni-1.2-flash`, `veo-3.1-fast` (มีเสียง), `p-video` | **paid ทั้งหมด** — เช็คราคาก่อน: `-d 5` ที่ 0.08/s ≈ $0.40 |
| เสียงพูด | `tts-1` (เร็ว/ถูก) · `tts-1-hd` (คม) | Quest · ภาษาไทยใช้ได้ |
| เพลง/SFX | `elevenlabs/music-v2.5`, `stable-audio-3` | paid |
| 3D | `microsoft/trellis-2` (รับภาพอ้างอิง), `tripo3d/p2` | paid · trellis ต้องมี `-i` |

ดูของสด: `python3 $M models -t image --quest` · `-t video` · `-t audio --voices` · `-q flux --json`

## 🖋️ เขียน prompt (สูตรของอัย)

```
[ประธาน + การกระทำ] , [ฉาก] , [แสง] , [มุมกล้อง/เลนส์] , [สไตล์/อารมณ์] , [รายละเอียดห้ามขาด]
```

- ✅ `"a Thai-Chinese tea table, ceramic gaiwan steaming, soft window light from the left, 50mm shallow depth of field, warm muted film tones"`
- ❌ `"beautiful tea"` — กว้างเกิน model จะเดาเอง
- **แสงคุมอารมณ์มากที่สุด** — ระบุทุกครั้ง: `soft window light` / `golden hour backlight` / `hard flash` / `overcast flat light`
- อยากได้ชุดภาพเดียวกัน → ตรึง `--seed` เดิม เปลี่ยนแค่ส่วนที่ต้องการ
- **วิดีโอ:** เพิ่มการเคลื่อนกล้อง + จังหวะ → `slow dolly in`, `static locked-off shot`, `slow pan left, gentle parallax` (prompt นิ่ง ๆ = ภาพนิ่งทั้งคลิป)
- **TTS:** เขียนธรรมชาติ มีวรรคตอน = หยุดหายใจถูกที่ · ใช้ `--instructions "warm, playful"` กับ voice ที่รองรับ
- ตัวอักษรไทยในภาพมักเพี้ยน → ใช้ Ideogram/Nano Banana Pro หรือสร้างภาพเปล่าแล้วเติมข้อความด้วย `document-forge`

## 🧾 คำสั่งทั้งหมด

```bash
# ภาพ
python3 $M image "prompt" -a 16:9 --model black-forest-labs/flux.1.1-pro --seed 7 --free
python3 $M image "prompt" -s 1280x720 --transparent --quality hd --safe privacy --free
# แก้/แต่งภาพ (img2img — รับไฟล์ในเครื่องหรือ URL, หลายไฟล์ได้)
python3 $M edit "same cup, golden hour, warm amber tones" -i base.jpg
python3 $M edit "place this mug on a wooden desk, product photo" -i mug.png -i desk.jpg
# วิดีโอ (paid)
python3 $M video "a paper boat drifting on a calm river, slow dolly in" -d 5 -r 720p --model google/veo-3.1-fast
python3 $M video "..." -i first-frame.jpg --audio            # image→video / มีเสียง
# เสียง
python3 $M speech "สวัสดีค่ะหลุยส์" --voice nova -o out.mp3
python3 $M speech "..." --model openai/tts-1-hd --instructions "warm, playful"
python3 $M music "calm lo-fi piano with rain" -d 20 --instrumental     # paid
# 3D + ถอดเสียง
python3 $M 3d "ceramic cup" -i cup.jpg -r medium                        # paid
python3 $M transcribe audio.mp3 --language th -o talk.txt
# สำรวจ / ตรวจ
python3 $M models -t image --quest            # เฉพาะ model ที่ไม่ต้องเติมเงิน
python3 $M models -t audio --voices          # รายชื่อเสียงของแต่ละ model
python3 $M models -t video -q veo --json
python3 $M probe -n 12                       # ยิงทดสอบจริงว่า key ใช้ model ไหนได้
python3 $M balance ; python3 $M doctor
python3 $M usage                             # ใช้ pollen ไปเท่าไร แยกตามวัน/model
```

ทุกคำสั่งรับ `--json` (ผลลัพธ์เครื่องอ่าน) และ `-o <ไฟล์/โฟลเดอร์>`

## 📂 ที่เก็บไฟล์

| ชนิด | ค่าเริ่มต้น |
|---|---|
| ภาพ | `~/00-Inbox/Aiy/images/YYYY-MM-DD-HHMM-<slug>-<model>.<ext>` |
| วิดีโอ | `~/00-Inbox/Aiy/videos/` |
| เสียง | `~/00-Inbox/Aiy/audio/` |
| 3D | `~/00-Inbox/Aiy/models3d/` |
| transcript | `~/00-Inbox/Aiy/transcripts/` |

- นามสกุลไฟล์ตรวจจาก magic bytes จริง (server มักบอก `image/png` ทั้งที่ bytes เป็น JPEG) → ไฟล์ที่ได้เปิดได้เสมอ
- งานที่เป็น deliverable ของโปรเจกต์ → ย้ายไป `/home/lu5her/01-Projects/{Project}/` แล้วบอก path ใน Log

## ✅ Quality Gate ก่อนส่งงาน

1. **ภาพ:** `read` ไฟล์จริง — ตรงโจทย์ / นิ้ว-หน้าไม่พัง / ตัวหนังสืออ่านออก / composition ไม่โล่ง
2. **เสียง:** เช็คว่าเป็นไฟล์เสียงจริง (`file`), ความยาวพอ, ไม่ตัดกลางคำ
3. **วิดีโอ/3D:** เช็คขนาดไฟล์ > 100 KB และ `file` ยืนยัน container
4. ส่งพร้อม: path · model · seed (ถ้าต้องแก้ต่อ) · 1 บรรทัดว่าสื่ออะไร
5. ไม่ผ่าน → บอกตรง ๆ ว่าพลาดตรงไหน แล้วเสนอรอบถัดไป (อย่าป่อยของที่รู้ว่าไม่ดี)

## 💰 วินัยต้นทุน

0. **`--free` เป็นค่าตั้งต้นของทุกงานปกติ** — งานศิลป์/ภาพประกอบ/เสียงพูด ไม่มีเหตุผลต้องจ่าย · ใช้ของแพงเฉพาะเมื่อฟรีทำไม่ได้จริงและหลุยส์อนุมัติแล้ว
1. ร่างถูกก่อนเสมอ (schnell) → พอใจแล้วค่อยเรนเดอร์ `flux.1.1-pro`
2. วิดีโอ/เพลง/3D = paid: **คำนวณก่อนยิง** (ราคา/วินาที × จำนวนวินาที) และบอกหลุยส์ทุกครั้ง
3. ก่อนสร้างใหม่ ให้ `ls -t ~/00-Inbox/Aiy/{images,videos,audio}` ดูของเก่าก่อน
4. `402` = เงินไม่พอ ไม่ใช่บั๊ก → ลด model หรือแจ้งเติมเงิน

## 🔧 แก้ปัญหา

| อาการ | ทางแก้ |
|---|---|
| `401 UNAUTHORIZED` | key ว่าง/ผิด → `doctor` แล้วแก้ `secrets/pollinations.json` |
| `402 ... requires the paid balance` | model อยู่นอก Quest → เปลี่ยนเป็นตัวในคอลัมน์ซ้าย หรือเติมเงิน |
| `403 FORBIDDEN` (account:usage) | key ไม่มีสิทธิ์อ่าน balance — ไม่กระทบการสร้าง; ขอ grant ที่หน้า keys ถ้าต้องการตัวเลข |
| `400` + size | model บังคับขนาด (mai ต้อง ≥768px/ด้าน, flux-2 ต้องคูณ 16) → ลอง `-a 1:1` |
| `422` | policy บล็อก → ปรับถ้อยคำให้เป็นงานทั่วไป/ศิลปะ |
| `429` | retry อัตโนมัติ 3 รอบแล้ว; ยังพังให้เว้น 1-2 นาที |
| ภาพไม่เหมือน prompt | prompt กว้างไป → เติมแสง/เลนส์/สไตล์ หรือเปลี่ยน model |
| ไฟล์เปิดไม่ได้ | ควรถูกแก้ด้วย magic-byte sniffing แล้ว — ถ้ายังเจอ ให้รายงานเป็นบั๊ก |
| timeout (วิดีโอ/3D) | retry request เดิมได้เลย เซิร์ฟเวอร์ cache ผลไว้ ไม่เสียเงินซ้ำ |

## 🔗 เทียบกับ tool ในตัว

`experimental_image_generate` เหมาะกับ "ขอภาพเร็ว ๆ 1 ใบ" — skill นี้เหมาะเมื่อต้อง **เลือก model, คุม seed ทำชุดภาพเดียวกัน, แก้ภาพเดิม, สร้างวิดีโอ/เสียง/3D, หรือเซฟลง path ที่กำหนด**

---

— อัย (Aiy) · Strategic Orchestrator · Thu 19, Nov 26

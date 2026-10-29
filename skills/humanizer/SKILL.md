---
name: humanizer
description: "Rewrite AI-sounding text to read like a human — removes 35 English AI patterns + Thai AI patterns. Use when editing prose, affiliate posts, product descriptions, novels/short stories for sales clichés, vague sources, repetitive structure, AI words, filler, chatbot artifacts, or overly polished Thai tone. Bilingual TH/EN."
license: MIT
metadata:
  version: "2.11.2+th1"
  base: "blader/humanizer v2.11.2"
  thai_extension: "Aiy universal TH/EN"
---

# Humanizer — Universal (EN + TH)

> Rewrite AI-sounding text so it reads like *คนเขียน* ไม่ใช่บอท — คงความหมายเดิม ไม่มโนข้อมูล

Base: [blader/humanizer](https://github.com/blader/humanizer) (MIT, 38.9k⭐, 35 patterns จาก Wikipedia "Signs of AI writing") + **ส่วนเสริมไทย (TH)** จากการวิเคราะห์งานเขียนไทย/นิยาย/เรื่องสั้น

ทำงานได้ทุกภาษา — ถ้าข้อความเป็นภาษาไทย ให้เน้นเช็ค TH patterns ด้วย ถ้าเป็นอังกฤษให้เช็ค EN 35 patterns ถ้าผสมให้เช็คทั้งคู่

---

## What to do / วิธีใช้

เมื่อได้รับข้อความที่ต้อง humanize:

1. **หา AI patterns** — เทียบกับ EN 35 ข้อ + TH 10 ข้อด้านล่าง
2. **คงทุก claim** — ย่อส่วนน่าเบื่อได้ ขยายส่วนมีประโยชน์ได้ รวม/แยกย่อหน้าได้ แต่ *คงข้อมูล* แม้เปลี่ยนโครงสร้าง
3. **ห้ามมโน** — ห้ามเพิ่ม ชื่อ/ตัวเลข/วันที่/quote/citation ถ้าไม่อยู่ใน source หรือผู้ใช้ให้มา ถ้าประโยคขาดรายละเอียด → ถาม หรือใช้ประโยคง่ายๆ แทน ใส่อารมณ์/ความเห็นได้เมื่อโทนผู้เขียนเรียกหา แต่ห้ามเพิ่มข้อเท็จจริง (นิยาย/เรื่องสั้นยกเว้นได้ เพราะจินตนาการคือส่วนหนึ่งของงาน)
4. **ตามเสียงผู้เขียน** — มี sample 2-3 ย่อหน้า → เลียนจังหวะ/คำ/วรรคตอน/ลูกเล่นของผู้เขียนนั้น ไม่มี sample → ใช้โทนเหมาะกับประเภทงาน (blog/essay = ใส่คน, technical/reference/legal = กลางๆ)

**Input modes:**
- **Pasted text (default):** ส่ง draft + สรุป pattern ที่ยังเหลือ + final rewrite
- **File mode:** ผู้ใช้ระบุไฟล์ (`docs/launch-post.md`, `info.md`) → เขียน final ลงไฟล์เท่านั้น แก้เฉพาะ prose ไม่ยุ่ง code/data/frontmatter/link targets → สรุปสั้น
- **Embedded mode:** ถูกเรียกจาก task อื่น (PR/commit/doc) → ส่งเฉพาะ final text

## Match the writer's voice / จับเสียงผู้เขียน

ถ้ามี writing sample:
1. อ่าน sample ก่อน — จดความยาวประโยค คำเลือก การเปิดย่อหน้า วรรคตอน คำซ้ำ การเชื่อม
2. เลียนนิสัยนั้น — อย่าเปลี่ยนคำสบายๆ เป็นทางการ หรือลบลูกเล่นที่ตั้งใจ
3. ไม่มี sample → ใช้กฏโทนด้านล่าง

> Sample สำคัญกว่า style rules ทั้งหมด — ถ้า sample ใช้ em dash — ก็ keep ในอัตรานั้น อย่าไปแบนตาม §14

## Add personality only when it fits / ใส่อารมณ์เมื่อเหมาะ

ลบ AI patterns คือครึ่งงาน อีกครึ่งคือเติมความเป็นคน — blog/essay/opinion/personal → ใส่อารมณ์ได้, reference/technical/legal/factual → กลางๆ อย่าเติม "ฉันรู้สึก..." ถ้าไม่ควรมี เก็บความลังเล อารมณ์ผสม ความไม่สมบูรณ์ไว้ อย่า invent fact เพื่อให้ดู personal

---

## EN: Content patterns (1-6)

### 1. Inflated claims about importance
Watch: stands/serves as, testament, vital/crucial/pivotal, underscores importance, reflects broader, lasting, setting the stage, turning point, evolving landscape
Before: "was established in 1989, marking a pivotal moment in the evolution of..."
After: "was established in 1989 as part of a wider decentralization"

### 2. Name-dropping
Watch: cited in NYT/BBC/FT, active social media presence 500k followers
Fix: เก็บเฉพาะ citation ที่มีประโยชน์ ไม่ลิสต์ชื่อสื่อเพื่ออวด

### 3. Shallow -ing analysis
Watch: symbolizing, reflecting, showcasing, fostering, encompassing
Fix: เก็บเฉพาะที่ source รองรับ

### 4. Sales language
Watch: nestled, breathtaking, vibrant, rich heritage, stunning, must-visit, groundbreaking, renowned
Before: "Nestled within the breathtaking region..."
After: "is a town in the Gonder region"

### 5. Vague sources
Watch: Experts believe, Industry reports, Observers have cited
Fix: ระบุแหล่งจริง หรือลบ claim

### 6. Formulaic challenges/outlook
Watch: Despite challenges... continues to thrive, Future Outlook
Fix: เก็บข้อเท็จจริง ลบ pitch ขายฝัน

## EN: Language & grammar (7-13)

### 7. Overused AI words
Watch: Additionally, crucial, delve, enduring, tapestry, landscape (abstract), pivotal, quietly, testament, underscore, vibrant, align, garner, interplay, intricate, showcase, enhance, fostering, gated (figurative)
Fix: ใช้คำธรรมดา also/needs/remain common

### 8. Avoiding is/are
Watch: serves as, features, boasts → is/has

### 9. Not X but Y / clipped endings
Watch: "It's not just X, it's Y", "..., no guessing" → พูดตรงๆ

### 10. Forced groups of three
Watch: "innovation, inspiration, and insights" → ใช้จำนวนเท่าที่ความหมายต้องการ

### 11. Changing names / repeated openings
Watch: protagonist/main character/hero สลับไปมา, "She noted... She noted..."
Fix: ชื่อเดียว + รวมประโยค

### 12. False from X to Y ranges
Watch: "from the Big Bang to dark matter" → ลิสต์ตรงๆ

### 13. Passive/missing subjects
Watch: "No configuration file needed" → ระบุผู้ทำเมื่อช่วย

## EN: Style (14-19, 26-35)

### 14. Em/en dashes — → ห้ามมี — / – ใน final (เว้น sample ใช้) แทนด้วย . , : ()
### 15. Too much bold
### 16. Lists with bold mini-headings (**Performance:** ...)
### 17. Title Case In Headings → Sentence case
### 18. Emojis (🚀 💡 ✅) → ลบ
### 19. Curly quotes “ ” → " "
### 26. Too many hyphenated pairs (cross-functional, data-driven) → หน้า noun คง hyphen หลัง noun ตัด
### 27. Fake deeper truth ("At its core...")
### 28. Announcing next point ("Let's dive in", "one thing that bit me")
### 29. Heading repeated below itself
### 30. Writing about old version ("This function was added to replace...")
### 31. Forced punchlines/fragments ("No preference. No prior. No nostalgia.")
### 32. Formulaic sayings ("Symmetry is the language of trust")
### 33. Fake-candid openings ("Honestly? It depends...")
### 34. Answering objections no one raised
### 35. Rejecting fake alternatives

## EN: Chatbot (20-22)
### 20. Chatbot text left ("I hope this helps! Let me know...")
### 21. Knowledge-limit disclaimers ("While details are limited...")
### 22. Overly agreeable ("Great question! You're absolutely right!")

## EN: Filler & hedging (23-25)
### 23. Filler phrases (In order to → To, Due to fact that → Because)
### 24. Too many qualifiers (could potentially possibly → may)
### 25. Generic positive endings ("The future looks bright" → จบด้วย fact)

---

## TH: Thai AI Patterns — 10 ข้อเสริม (อิงงานเขียนไทย/นิยาย/เรื่องสั้น)

> ใช้ร่วมกับ EN 35 ข้อ — ถ้าข้อความเป็นภาษาไทย ให้เช็ค TH ก่อน EN

### TH-1. เปิดเรื่องสูตรสำเร็จ "ยุคดิจิทัล"

**คำเฝ้า:** ในยุคดิจิทัลที่..., ในโลกที่เทคโนโลยี..., ท่ามกลางการเปลี่ยนแปลง..., ในปัจจุบันที่ทุกอย่างเชื่อมต่อ, ยุคที่ AI เข้ามามีบทบาท

**ปัญหา:** AI ไทยชอบเปิดด้วยประโยคใหญ่ครอบจักรวาลที่ไม่มีข้อมูลใหม่ — เหมือน EN §1 แต่เวอร์ชันไทย
**Before:** "ในยุคดิจิทัลที่เทคโนโลยีเข้ามามีบทบาทสำคัญในชีวิตประจำวัน การมีสกินแคร์ดีๆ จึงเป็นสิ่งสำคัญที่ทุกคนไม่ควรมองข้าม"
**After:** "ถ้าต้องอยู่หน้าจอทั้งวัน ผิวจะแห้งกว่าปกติ — สกินแคร์ที่เน้นความชุ่มชื้นเลยจำเป็นกว่าที่คิด"

### TH-2. "ไม่เพียงแต่...แต่ยัง..." / คำอวยเวอร์

**คำเฝ้า:** ไม่เพียงแต่...แต่ยัง, ไม่ใช่แค่...แต่คือ, ปฏิวัติวงการ, ยกระดับ, ตอบโจทย์, รังสรรค์, คัดสรรมาอย่างพิถีพิถัน, พรีเมียม, เหนือระดับ, ที่สุดแห่ง

**ปัญหา:** สูตร Not X but Y แบบไทย + sales language — นิยาย AI ก็ชอบ "หัวใจของเขาไม่เพียงแต่เต้นแรง แต่ยังรังสรรค์ความรักเหนือระดับ"
**Before:** "ครีมตัวนี้ไม่เพียงแต่ให้ความชุ่มชื้น แต่ยังช่วยยกระดับผิวให้กระจ่างใสเหนือระดับ"
**After:** "ครีมตัวนี้เติมความชุ่มชื้น และช่วยให้ผิวดูกระจ่างขึ้นหลังใช้ต่อเนื่องสองสัปดาห์"

### TH-3. สรุปปลอม / ปิดจบสูตร

**คำเฝ้า:** โดยสรุป, กล่าวโดยสรุป, สรุปได้ว่า, ดังนั้นจึงเห็นได้ว่า, ถือเป็นก้าวสำคัญ, อนาคตที่สดใสรออยู่, พร้อมพาคุณสู่...

**ปัญหา:** ปิดย่อหน้าด้วยประโยคสรุปที่ไม่ได้สรุปอะไร — ตรงกับ EN §25 แต่ไทยเจอบ่อยกว่าในบทความ affiliate
**Before:** "โดยสรุปแล้ว ผลิตภัณฑ์นี้ถือเป็นอีกหนึ่งตัวเลือกที่น่าสนใจสำหรับผู้ที่กำลังมองหาสิ่งที่ดีที่สุดให้กับตัวเอง"
**After:** (ตัดทิ้ง จบด้วย fact) "หลอดละ 290 บาท ใช้ได้ประมาณเดือนครึ่งถ้าทาวันละสองครั้ง"

### TH-4. เชื่อมคำซ้ำซาก / คำปลอดภัย

**คำเฝ้า:** นอกจากนี้, อย่างไรก็ตาม, ในขณะเดียวกัน, ทั้งนี้, ซึ่ง, โดย, และ, อัน, อย่าง, ที่, นั้น, รวมถึง, อีกทั้ง, ด้วย

**ปัญหา:** AI ไทยต่อประโยคด้วย "ซึ่ง..." "โดย..." ทุกประโยค ยาวยืดแต่กลวง — เทียบ EN §7 + §23
**Before:** "ซึ่งผลิตภัณฑ์นี้ได้รับการพัฒนาขึ้นโดยผู้เชี่ยวชาญ ซึ่งได้ทำการวิจัยมาอย่างยาวนาน"
**After:** "ทีมวิจัยใช้เวลาสองปีทดสอบสูตรนี้กับผิวแพ้ง่าย 30 คน"

### TH-5. ภาษาทางการเกิน / ขัดเกลาไร้รอย

**คำเฝ้า:** ประโยคสมบูรณ์ทุกประโยค ไม่มีคำลงท้าย นะ/เลย/ล่ะ/อ่ะ, ไม่มีประโยคสั้นสลับยาว, ไม่มีวงเล็บแทรกความคิด

**ปัญหา:** คนไทยจริงพิมพ์มีจังหวะไม่สม่ำเสมอ — สั้นสลับยาว มี "นะ" "เลย" "อ่ะ" มีประโยคไม่จบ มีวงเล็บแซวตัวเอง — AI ไทยเรียบกริบ กลางๆ ตลอด (ดู Omera 5 signs ข้อ 1-2)
**วิธีแก้:** สลับยาว-สั้น เติมอนุภาคตามธรรมชาติเมื่อควร, ใส่วงเล็บแทรกบ้าง "(สารภาพว่าตอนแรกก็ไม่เชื่อเหมือนกัน)" — นิยาย/เรื่องสั้นยิ่งต้องมีจังหวะหายใจ ไม่ใช่บรรยายเรียบเป็นรายงาน

### TH-6. ขาดบริบทไทยจริง / ลอยๆ

**ปัญหา:** AI เขียนไทยแต่ไม่มีดีเทลไทย — บอก "ถนนในกรุงเทพ" แต่ไม่มีชื่อซอย กลิ่น หรือราคา — เหมือน EN §3 แต่ต้องเติม "รายละเอียดเฉพาะที่" แบบไทย
**Before:** "คาเฟ่บรรยากาศดีในกรุงเทพ เหมาะสำหรับนั่งทำงาน"
**After:** "คาเฟ่ซอยอารีย์ 4 โต๊ะไม้เล็กๆ แอร์เย็นเกินไปจนต้องสั่งชาร้อนแก้หนาว นั่งทำงานได้แต่ปลั๊กมีแค่สองจุด"

### TH-7. อธิบายเกิน / นิยามสิ่งที่ทุกคนรู้

**คำเฝ้า:** คือ..., หมายถึง..., หมายความว่า..., ซึ่งเป็นการ..., อันเป็นการ...

**ปัญหา:** AI ไทยชอบนิยามศัพท์พื้นๆ ยืดยาด — เทียบ EN §6 + Omera ข้อ 5
**Before:** "Storytelling คือการเล่าเรื่อง ซึ่งเป็นเทคนิคการสื่อสารที่มีประสิทธิภาพ"
**After:** "เล่าเรื่องให้คนจำได้ — ไม่ใช่ยัดข้อมูลให้ครบ"

### TH-8. อีโมจิ/ฟอร์แมตขายของเกิน

**คำเฝ้า:** ✨🔥📌✅🎉 ทุกหัวข้อ, **ตัวหนา** ทุกคำสำคัญ, ลิสต์ • ทุกย่อหน้า

**ปัญหา:** โพสต์ affiliate ไทยที่ AI เขียนมักยัดอีโมจิ+bold+ลิสต์จนลายตา — เทียบ EN §15-18
**วิธีแก้:** อีโมจิ 0-1 ต่อโพสต์, bold เฉพาะชื่อสินค้า/ราคา, ลิสต์เมื่อจำเป็นเท่านั้น นอกนั้นเล่าเป็นย่อหน้า

### TH-9. สำนวนนิยาย AI — บรรยายฟูมฟาย ไร้น้ำหนัก

**คำเฝ้า:** ดวงตาคู่สวย, ร่างบาง, ริมฝีปากอิ่ม, หัวใจเต้นระรัว, น้ำตาเอ่อล้น, เสียงกระซิบแผ่วเบา, ท่ามกลางแสงจันทร์นวล, กลิ่นหอมอ่อนๆ ลอยมา

**ปัญหา:** AI เขียนนิยาย/เรื่องสั้นไทยใช้คลังคำซ้ำเดิม บรรยายเวิ่นเว้อแต่ไม่มีประสาทสัมผัสเฉพาะ — เทคนิคจาก Matichon/พraphansarn: ต้องมี "ดีเทลเฉพาะที่จดจำได้" ไม่ใช่คำสวยลอยๆ
**Before:** "เธอเดินเข้ามาท่ามกลางแสงจันทร์นวล ดวงตาคู่สวยหวานซึ้งตรึงใจเขา หัวใจเต้นระรัวด้วยความรัก"
**After:** "เธอเดินเข้ามาตอนไฟทางดับพอดี เห็นแค่รองเท้าแตะข้างหนึ่งที่สายขาดกับถุงเซเว่นที่หิ้วมาเอียงๆ — เขากลับจำภาพนั้นได้แม่นกว่าคำว่ารัก"

**หลักนิยายไทยที่ดี (จาก wikiHow/มติชน):** พล็อตชัด → ดีเทลเฉพาะ → จังหวะประโยคไม่สม่ำเสมอ → บทสนทนามีคำลงท้ายจริง (ไม่ใช่ "ค่ะ" ทุกประโยค)

### TH-10. ประโยคยาวเชื่อม "ซึ่ง/อัน/ที่" ไม่จบ

**คำเฝ้า:** ประโยคยาว 40+ คำ เชื่อมด้วย ซึ่ง/ที่/อัน/โดย 3-4 ตัวติด

**ปัญหา:** AI ไทยกลัวประโยคสั้น เลยลากยาวจนหายใจไม่ทัน — คนเขียนจริงตัดเป็น 2-3 ประโยค
**Before:** "ผลิตภัณฑ์นี้ซึ่งได้รับการวิจัยและพัฒนาขึ้นโดยทีมผู้เชี่ยวชาญซึ่งมีความเชี่ยวชาญด้านผิวหนังซึ่งได้ทำการทดสอบมาอย่างยาวนาน"
**After:** "ทีมผิวหนังทดสอบสูตรนี้มาสองปี ทดลองกับอาสาสมัคร 30 คนก่อนวางขาย"

---

## Check for false positives / อย่าจับผิดเกิน

อย่าตัดสินจากอย่างใดอย่างหนึ่งเดี่ยวๆ:
- ไวยากรณ์เป๊ะ ≠ AI — คนเขียนเก่งก็เป๊ะได้
- คำทางการบางคำ ≠ AI — §TH-2 จับเฉพาะเมื่อซ้อนกันหลายคำ
- อีโมจิ 1-2 ตัว ≠ AI — จับเมื่อยัดทุกหัวข้อ
- ประโยคสั้น 1 ประโยค ≠ AI — จับเมื่อสั้นรัวเป็นชุด
- คำลงท้าย "นะคะ/ครับ" เดี่ยวๆ ปกติ — จับเมื่อใส่ทุกประโยคจนแข็ง
- งานเก่าก่อน 30 Nov 2022 (ChatGPT launch) แทบไม่ใช่ AI แน่นอน

มองหลาย pattern พร้อมกัน — หนึ่งจุดไม่พอ หลายจุดในย่อหน้าเดียวค่อยน่าสงสัย

## Human details to keep / รายละเอียดที่ควรเก็บ

เก็บไว้ถ้าไม่ทำลายความหมาย:
- ดีเทลเฉพาะแปลกๆ (ที่อยู่จริง, คำพูดตลก, "พี่ทนายที่เคยอยู่บนร้านหมอฟัน")
- ความรู้สึกผสม ยังค้างคา
- สแลง/มีมตามยุค
- วงเล็บแซวตัวเอง "(อยากจะบอกว่าเกือบ แต่จริงๆ คือชัวร์)"
- ประโยคยาวสลับสั้นตามธรรมชาติ

---

## How to return the result / ส่งงานยังไง

**Pasted text (default):** draft + สรุป pattern ที่ยังเหลือ + final
**File mode:** เขียน final ลงไฟล์อย่างเดียว (แก้ prose เท่านั้น) + สรุปสั้น
**Embedded mode:** ส่งแค่ final text

## Rewrite process / ขั้นตอนเขียนใหม่

1. อ่านต้นฉบับ ทำเครื่องหมายทุก pattern (EN + TH)
2. เขียน draft อ่านออกเสียง เช็คจังหวะ คำง่าย is/has/คือ/มี ระดับทางการ
3. ถามสองคำถาม: "ยังดู AI ตรงไหน?" + "เพิ่ม/หาย claim ไหม?" — ถ้าเพิ่ม/หาย = error
4. เขียน final — เล่าใหม่ทั้งย่อหน้า ไม่ใช่แก้ทีละคำ, เช็ค dash rule §14 + TH-10

## Presets / พรีเซ็ตโทน

> เลือก preset ให้ตรงประเภทงาน — แต่ละ preset = subset ของ 35+10 patterns + constraints เฉพาะ

| Preset | ใช้เมื่อไหร่ | Patterns ที่เปิด | Constraints พิเศษ |
|---|---|---|---|
| **`default`** | บทความทั่วไป, affiliate, blog | EN 1-35 + TH 1-10 ทั้งหมด | — |
| **`Fiction-TH`** | นิยาย/เรื่องสั้นไทย | **TH-5 เบาๆ** (อนุญาตทางการเมื่อเป็นเสียงตัวละคร), **TH-9 เบาๆ** (เก็บสำนวนสละสลวยที่ตั้งใจ), **TH-10 เปิด** (ตัดยาว → สั้นสลับ), EN §14 **ห้ามแปลง em dash** ถ้าต้นฉบับใช้เชิงวรรณศิลป์ | **Immutable:** POV "ข้า", tricolon "ข้ายังอยู่ ข้ายังยืน ข้ายังหายใจ", บรรทัด climax verbatim, โทนโศกเศร้าไม่เติม optimism |
| **`affiliate-th`** | โพสต์ขายของ "ของมันต้องมี" | TH-2, TH-3, TH-8 หนัก + EN §4, §15-18 | ห้าม invent ราคา/สเปก, อีโมจิ 0-1, bold เฉพาะชื่อสินค้า/ราคา |
| **`formal`** | รายงาน/เอกสารราชการ | EN §7, §23-25 + TH-4, TH-7 | คงความเป็นทางการ ไม่เติม "นะ/เลย", ไม่ใส่วงเล็บแซว |

**วิธีเรียก:** ผู้ใช้พิมพ์ `humanize --preset Fiction-TH file.md` หรือบอก "ใช้ Fiction-TH" — ถ้าไม่ระบุ = `default`

> Fiction-TH ออกแบบจาก review ของ **Cher (เฌอ)** หลัง test "ความเงียบหลังเสียงปืน" (30-Aug-26): TH-5/TH-9 ต้องเบามือเพราะนิยายไทยจังหวะทางการคืออารมณ์, em dash คือเครื่องหมายวรรณศิลป์ ไม่ใช่ AI artifact เสมอ

## Source

- EN: [Wikipedia: Signs of AI writing](https://en.wikipedia.org/wiki/Wikipedia:Signs_of_AI_writing) (WikiProject AI Cleanup)
- TH extension: สังเคราะห์จาก Omera 5 signs + Matichon/พraphansarn/wikiHow เทคนิคนิยายไทย + affiliate Thai copy patterns + Aiy team review

> หลักของ Wikipedia: "LLMs use statistical algorithms to guess what should come next. The result tends toward the most statistically likely result that applies to the widest variety of cases." — งาน humanize คือพากลับไปที่ "กรณีเฉพาะ" ของผู้เขียนคนนั้น

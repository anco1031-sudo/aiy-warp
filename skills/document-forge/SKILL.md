---
name: document-forge
description: "Universal Document Factory — สร้างไฟล์จริง Word (.docx), Excel (.xlsx), PowerPoint (.pptx), PDF ด้วย python-docx/openpyxl/python-pptx/reportlab + TH Sarabun + ครุฑราชการ. ใช้เมื่อต้องสร้างเอกสารจริงส่งให้ Louis ไม่ใช่แค่ร่างข้อความ. Triggers: สร้างไฟล์, word, excel, slide, pdf, เอกสาร, รายงาน, ทะเบียน, สไลด์, ส่งออกไฟล์"
---

# Document Forge — Universal Document Factory

> สร้างไฟล์จริง ไม่ใช่แค่ข้อความ — Word / Excel / PowerPoint / PDF พร้อม TH Sarabun และฟอร์มราชการ 100%

## 📍 Output Rule (บังคับทุก agent)

**ไฟล์ deliverable ทุกชิ้นที่สร้างให้ Louis ต้องวางใน `~/00-Inbox/{Team}/` เท่านั้น — ห้ามวางใน Workspace หรือ /tmp แล้วจบ**

| ทีม | Output Folder (absolute) |
|---|---|
| Legal_Team (Jean/Rin) | `/home/lu5her/00-Inbox/Legal/` |
| Product_Team (Lin/An/Pao/Mint/Fah/Cloud) | `/home/lu5her/00-Inbox/Product/` |
| Content_Team (Cher/Ploy/Jewel) | `/home/lu5her/00-Inbox/Content/` |
| Security_Team (Sai/Pleng) | `/home/lu5her/00-Inbox/Security/` |
| Trader_Team (Kwan/Fon/June/Bee/Nam) | `/home/lu5her/00-Inbox/Trader/` |
| Aiy (Orchestrator) | `/home/lu5her/00-Inbox/Aiy/` |
| ทั่วไป/ไม่ระบุทีม | `/home/lu5her/00-Inbox/` |

> ถ้าโฟลเดอร์ยังไม่มี → `mkdir -p` ก่อนเขียนเสมอ — อย่าเดาว่ามีแล้ว

## 🔧 Stack ที่มีอยู่แล้ว (ไม่ต้องติดตั้งใหม่)

```bash
python3 -c "import docx; print(docx.__version__)"   # 1.2.0 ✅
python3 -c "import openpyxl; print(openpyxl.__version__)" # 3.1.5 ✅
python3 -c "import pptx; print(pptx.__version__)"   # 1.0.2 ✅
python3 -c "import reportlab; print(reportlab.Version)" # 5.0.1 ✅
# pandas มีถ้าต้อง, pandoc/libreoffice ถ้าต้องแปลง PDF
```

ติดตั้งเพิ่มถ้าขาด: `pip install python-docx openpyxl python-pptx reportlab pandas` หรือ `uv pip install ...`

## 📝 Template 1: Word หนังสือราชการ (.docx) — TH Sarabun + ครุฑ

```python
from docx import Document
from docx.shared import Pt, Inches, RGBColor
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.oxml.ns import qn
import os

output = "/home/lu5her/00-Inbox/Legal/หนังสือภายนอก-ขอความอนุเคราะห์-25690902.docx"
os.makedirs(os.path.dirname(output), exist_ok=True)

doc = Document()

# ตั้งฟอนต์ default เป็น TH Sarabun New
style = doc.styles['Normal']
font = style.font
font.name = 'TH Sarabun New'
font.size = Pt(16)
# แก้ fallback สำหรับภาษาไทย
r = style.element.get_or_add_rPr()
rFonts = r.get_or_add_rFonts()
rFonts.set(qn('w:ascii'), 'TH Sarabun New')
rFonts.set(qn('w:hAnsi'), 'TH Sarabun New')
rFonts.set(qn('w:cs'), 'TH Sarabun New')

# ระยะขอบ A4 มาตรฐานราชการ
for section in doc.sections:
    section.top_margin = Inches(0.75)
    section.bottom_margin = Inches(0.75)
    section.left_margin = Inches(1.0)
    section.right_margin = Inches(0.75)
    section.page_height = Inches(11.69)
    section.page_width = Inches(8.27)

# ครุฑ (placeholder — แทนด้วยรูปจริงถ้ามีไฟล์ครุฑ.png)
p = doc.add_paragraph()
p.alignment = WD_ALIGN_PARAGRAPH.CENTER
run = p.add_run('[ครุฑพ่าห์ — ใส่รูปครุฑสูง 1.5cm ตรงนี้]')
run.font.size = Pt(10)
run.font.color.rgb = RGBColor(0x99, 0x99, 0x99)

# ที่ / วันที่
p = doc.add_paragraph()
p.add_run('ที่ นร ๐๑๐๑/ว ๑๒๓').bold = True
p = doc.add_paragraph('วันที่ ๒ กันยายน พ.ศ. ๒๕๖๙')

# เรื่อง / เรียน
doc.add_paragraph('เรื่อง    ขอความอนุเคราะห์........................')
doc.add_paragraph('เรียน    ............................................')
doc.add_paragraph('อ้างถึง  ............................................')
doc.add_paragraph('สิ่งที่ส่งมาด้วย  ....................................')

# เนื้อหา
doc.add_paragraph('ตามที่ ............................................ นั้น')
p = doc.add_paragraph('ในการนี้ ............................................')
p.add_run(' จึงเรียนมาเพื่อโปรดพิจารณาดำเนินการต่อไป').bold = False

# ลงท้าย
p = doc.add_paragraph()
p.alignment = WD_ALIGN_PARAGRAPH.RIGHT
p.add_run('ขอแสดงความนับถือ').bold = False

p = doc.add_paragraph()
p.alignment = WD_ALIGN_PARAGRAPH.RIGHT
p.add_run('(นาย...................................)')
p = doc.add_paragraph()
p.alignment = WD_ALIGN_PARAGRAPH.RIGHT
p.add_run('ตำแหน่ง...................................')

# ส่วนราชการเจ้าของเรื่อง
doc.add_paragraph('ส่วนราชการเจ้าของเรื่อง: ........................ โทร. .................... โทรสาร ....................')

doc.save(output)
print(f"Saved: {output}")
```

## 📊 Template 2: Excel ทะเบียนหนังสือ (.xlsx)

```python
import openpyxl
from openpyxl.styles import Font, PatternFill, Alignment, Border, Side
import os

output = "/home/lu5her/00-Inbox/Legal/ทะเบียนหนังสือรับ-ส่ง-256909.xlsx"
os.makedirs(os.path.dirname(output), exist_ok=True)

wb = openpyxl.Workbook()
ws = wb.active
ws.title = "ทะเบียนหนังสือ"

# หัวตาราง
headers = ["ลำดับ", "เลขที่หนังสือ", "ลงวันที่", "จาก", "ถึง", "เรื่อง", "ชั้นความเร็ว", "ชั้นความลับ", "ผู้รับผิดชอบ", "หมายเหตุ"]
header_fill = PatternFill(start_color="37474F", end_color="37474F", fill_type="solid")
header_font = Font(name="TH Sarabun New", size=12, bold=True, color="FFFFFF")

for col, h in enumerate(headers, 1):
    c = ws.cell(row=1, column=col, value=h)
    c.fill = header_fill
    c.font = header_font
    c.alignment = Alignment(horizontal="center", vertical="center", wrap_text=True)

# ตัวอย่างข้อมูล
rows = [
    [1, "นร 0101/ว123", "02 ก.ย. 2569", "อบต.บึงทวาย", "อำเภอ", "ขอความอนุเคราะห์", "ปกติ", "-", "ริน", ""],
]
thin = Side(style="thin", color="BDBDBD")
border = Border(left=thin, right=thin, top=thin, bottom=thin)
for r in rows:
    ws.append(r)
    for c in ws[ws.max_row]:
        c.font = Font(name="TH Sarabun New", size=11)
        c.border = border
        c.alignment = Alignment(vertical="center")

# ปรับความกว้าง
widths = [8, 18, 14, 18, 14, 30, 12, 12, 14, 18]
for i, w in enumerate(widths, 1):
    ws.column_dimensions[openpyxl.utils.get_column_letter(i)].width = w

ws.freeze_panes = "A2"
ws.auto_filter.ref = ws.dimensions
wb.save(output)
print(f"Saved: {output}")
```

## 🎨 Template 3: PowerPoint สไลด์ราชการ (.pptx)

```python
from pptx import Presentation
from pptx.util import Inches, Pt
from pptx.dml.color import RGBColor
from pptx.enum.text import PP_ALIGN
import os

output = "/home/lu5her/00-Inbox/Legal/สไลด์สรุปผู้บริหาร-25690902.pptx"
os.makedirs(os.path.dirname(output), exist_ok=True)

prs = Presentation()
prs.slide_width = Inches(13.33)
prs.slide_height = Inches(7.5)

# สไลด์ปก
slide = prs.slides.add_slide(prs.slide_layouts[6]) # blank
bg = slide.background.fill
bg.solid()
bg.fore_color.rgb = RGBColor(0x37, 0x47, 0x4F)
# ชื่อเรื่อง
txBox = slide.shapes.add_textbox(Inches(0.5), Inches(2.5), Inches(12.33), Inches(1))
tf = txBox.text_frame
tf.word_wrap = True
p = tf.paragraphs[0]
p.text = "รายงานสรุปผู้บริหาร"
p.font.size = Pt(40)
p.font.bold = True
p.font.color.rgb = RGBColor(0xFF, 0xFF, 0xFF)
p.font.name = "TH Sarabun New"
p.alignment = PP_ALIGN.CENTER

prs.save(output)
print(f"Saved: {output}")
```

## 📄 Template 4: PDF ด้วย reportlab (ถ้าต้องการ PDF ตรง ไม่ผ่าน Word)

```python
from reportlab.lib.pagesizes import A4
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.platypus import SimpleDocTemplate, Paragraph
from reportlab.lib.styles import ParagraphStyle
# ต้องมีฟอนต์ THSarabunNew.ttf ใน /usr/share/fonts หรือ project fonts/
```

## ✅ Checklist ก่อนส่ง Louis

- [ ] ไฟล์อยู่ใน `~/00-Inbox/{Team}/` ตามตาราง Output Rule แล้ว
- [ ] ชื่อไฟล์สื่อความหมาย + วันที่ `YYYYMMDD` หรือ `2569MMDD` (เช่น `หนังสือภายนอก-ขอความอนุเคราะห์-25690902.docx`)
- [ ] ฟอนต์ TH Sarabun New 16pt ทั้งไฟล์ (Word/Excel/PDF)
- [ ] ระยะขอบ A4 ถูกต้อง, ครุฑ placeholder หรือรูปจริงครบ
- [ ] ทดสอบเปิดไฟล์ได้ (docx/xlsx/pptx เปิดด้วย LibreOffice/Word ได้, pdf เปิดได้)
- [ ] Log ใน `Workspace/{Team}/{Name}/Logs/YYYY-MMWww.md` ว่าไฟล์อยู่ที่ไหน

## 🔗 ใช้ร่วมกับ Rin อย่างไร

Rin เป็นผู้เชี่ยวชาญฟอร์มราชการ — ทุกครั้งที่ Louis สั่งสร้างไฟล์ราชการ:
1. Rin เลือกประเภทหนังสือ (6 ชนิด) + ครุฑ + คำขึ้นต้น-ลงท้าย ให้ถูกต้อง
2. Rin เรียก Template จาก skill นี้ (Word/Excel/Slide) แล้วเติมเนื้อหาจริง
3. Rin เซฟลง `~/00-Inbox/Legal/` ตาม Output Rule แล้วรายงาน path ให้ Louis

> ห้ามสร้างไฟล์แล้วทิ้งไว้ใน /tmp หรือ Workspace แล้วจบ — ต้องอยู่ใน 00-Inbox เท่านั้น Louis ถึงจะหาเจอ

#!/usr/bin/env python3
"""
office/doc_write.py — Word generation from JSON blocks (Aetox 4.3)
Usage: python3 scripts/office/doc_write.py <input.json> [--out report.docx]
Blocks: heading/paragraph/bullets/numbered/table/lineitems/image
Deps: python-docx (pip install python-docx)
"""
import json, sys, argparse
def write_doc(data: dict, out_path: str=None) -> str:
    try: from docx import Document
    except ImportError: return "ERROR: pip install python-docx"
    path = out_path or data.get("path","report.docx")
    doc=Document()
    for b in data.get("blocks",[]):
        k=b.get("type","")
        if k=="heading": doc.add_heading(b["text"], level=b.get("level",1))
        elif k=="paragraph": doc.add_paragraph(b["text"])
        elif k=="bullets":
            for it in b.get("items",[]): doc.add_paragraph(it, style="List Bullet")
        elif k=="numbered":
            for it in b.get("items",[]): doc.add_paragraph(it, style="List Number")
        elif k=="table":
            cols=b.get("columns",[]); rows=b.get("rows",[])
            tbl=doc.add_table(rows=len(rows)+1, cols=len(cols))
            for i,c in enumerate(cols): tbl.cell(0,i).text=c
            for ri, row in enumerate(rows,1):
                for ci, cell in enumerate(row): tbl.cell(ri,ci).text=str(cell)
        elif k=="image":
            try: doc.add_picture(b["image"])
            except Exception as e: doc.add_paragraph(f"[image missing: {e}]")
    doc.save(path); return path
if __name__=="__main__":
    ap=argparse.ArgumentParser()
    ap.add_argument("json_file")
    ap.add_argument("--out", default=None)
    args=ap.parse_args()
    with open(args.json_file) as f: data=json.load(f)
    print(write_doc(data, args.out))

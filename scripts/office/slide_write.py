#!/usr/bin/env python3
"""
office/slide_write.py — PPTX generation from JSON (Aetox 4.4)
Usage: python3 scripts/office/slide_write.py <input.json> [--out deck.pptx]
JSON: {"path":"deck.pptx","slides":[{"title":"","content":""}]}
Deps: python-pptx (pip install python-pptx)
"""
import json, sys, argparse
def write_pptx(data: dict, out_path: str=None) -> str:
    try: from pptx import Presentation
    except ImportError: return "ERROR: pip install python-pptx"
    path = out_path or data.get("path","deck.pptx")
    prs=Presentation()
    for s in data.get("slides",[]):
        slide=prs.slides.add_slide(prs.slide_layouts[1])
        slide.shapes.title.text=s.get("title","")
        try: slide.placeholders[1].text_frame.text=s.get("content","")
        except: pass
    prs.save(path); return path
if __name__=="__main__":
    ap=argparse.ArgumentParser()
    ap.add_argument("json_file")
    ap.add_argument("--out", default=None)
    args=ap.parse_args()
    with open(args.json_file) as f: data=json.load(f)
    print(write_pptx(data, args.out))

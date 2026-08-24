#!/usr/bin/env python3
"""
pdf/read.py — PDF text extraction via poppler (Aetox 3.4) + fallback pymupdf/pdfplumber
Usage: python3 scripts/pdf/read.py <file.pdf> [--raw] [--json]
Returns: text to stdout
Deps: poppler (pdftotext), fallback: pymupdf / pdfplumber
"""
import subprocess, sys, os, json, argparse

def read_poppler(pdf: str) -> str | None:
    try:
        r = subprocess.run(["pdftotext","-layout",pdf,"-"], capture_output=True, text=True)
        if r.returncode==0 and r.stdout.strip():
            return r.stdout.strip()
    except FileNotFoundError:
        pass
    return None

def read_fallback(pdf: str) -> str:
    try:
        import fitz  # pymupdf
        doc=fitz.open(pdf)
        return "\n".join([p.get_text() for p in doc]).strip()
    except: pass
    try:
        import pdfplumber
        out=[]
        with pdfplumber.open(pdf) as pdf_obj:
            for p in pdf_obj.pages:
                t=p.extract_text()
                if t: out.append(t)
        if out: return "\n".join(out).strip()
    except: pass
    return ""

def read_pdf(pdf: str) -> str:
    if not os.path.exists(pdf): return f"file not found: {pdf}"
    txt = read_poppler(pdf)
    if txt: return txt
    return read_fallback(pdf)

if __name__=="__main__":
    ap=argparse.ArgumentParser()
    ap.add_argument("pdf")
    ap.add_argument("--json", action="store_true")
    args=ap.parse_args()
    txt=read_pdf(args.pdf)
    if args.json: print(json.dumps({"text": txt[:20000]}, ensure_ascii=False, indent=2))
    else: print(txt)

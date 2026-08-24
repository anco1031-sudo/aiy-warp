#!/usr/bin/env python3
"""
media/ocr.py — Image OCR via Tesseract (Aetox 3.1 pattern)
Usage: python3 scripts/media/ocr.py <image_path> [--langs tha+eng] [--tsv]
Returns: {"text": str, "confidence": float, "warning": str}
Deps: tesseract 5.5+, tha+eng traineddata
"""
import subprocess, os, sys, tempfile, json, argparse

def parse_tsv(tsv_path: str) -> float:
    total, count = 0, 0
    try:
        with open(tsv_path, encoding="utf-8") as f:
            for i, line in enumerate(f):
                if i == 0: continue
                cols = line.rstrip("\n").split("\t")
                if len(cols) >= 12 and cols[11].strip() == "5":  # block level? fallback to col 10
                    pass
                if len(cols) >= 11 and cols[10].strip().lstrip("-").replace(".","",1).isdigit():
                    try:
                        conf = float(cols[10])
                        if conf >= 0:
                            total += conf; count += 1
                    except: continue
    except FileNotFoundError:
        return -1
    return total / count if count else -1

def ocr_image(image_path: str, langs: str = "tha+eng") -> dict:
    if not os.path.exists(image_path):
        return {"text": "", "confidence": -1, "warning": f"file not found: {image_path}"}
    with tempfile.TemporaryDirectory() as tmpdir:
        base = os.path.join(tmpdir, "page")
        cmd = ["tesseract", image_path, base, "-l", langs, "txt", "tsv"]
        try:
            subprocess.run(cmd, check=True, capture_output=True)
        except subprocess.CalledProcessError as e:
            return {"text": "", "confidence": -1, "warning": e.stderr.decode()[:500] if e.stderr else str(e)}
        text = ""
        try:
            with open(f"{base}.txt", encoding="utf-8") as f: text = f.read().strip()
        except: pass
        conf = parse_tsv(f"{base}.tsv")
        warn = f"(confidence {conf:.0f}% — text may not match image)" if 0 <= conf < 70 else ""
        return {"text": text, "confidence": round(conf,1), "warning": warn}

if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("image", help="path to image")
    ap.add_argument("--langs", default="tha+eng")
    ap.add_argument("--json", action="store_true", help="output json")
    args = ap.parse_args()
    r = ocr_image(args.image, args.langs)
    if args.json: print(json.dumps(r, ensure_ascii=False, indent=2))
    else:
        print(r["text"])
        if r["warning"]: print(f"\n⚠️ {r['warning']}", file=sys.stderr)
        print(f"\n[confidence: {r['confidence']}%]", file=sys.stderr)

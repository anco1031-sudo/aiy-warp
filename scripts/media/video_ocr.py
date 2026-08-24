#!/usr/bin/env python3
"""
media/video_ocr.py — Video OCR via ffmpeg frames + Tesseract (Aetox 3.3)
Usage: python3 scripts/media/video_ocr.py <video.mp4> [--interval 5] [--langs tha+eng] [--json]
Returns: [{"time":"[m:ss]","text":str}, ...]
Deps: ffmpeg, tesseract
"""
import subprocess, os, sys, tempfile, glob, json, argparse
from ocr import ocr_image

def video_ocr(video_path: str, interval: int = 5, langs: str = "tha+eng") -> list:
    if not os.path.exists(video_path):
        return [{"error": f"file not found: {video_path}"}]
    with tempfile.TemporaryDirectory() as tmpdir:
        pat = os.path.join(tmpdir, "frame_%04d.jpg")
        try:
            subprocess.run(["ffmpeg","-y","-i",video_path,"-vf",f"fps=1/{interval}", pat],
                           capture_output=True, check=False)
        except Exception as e:
            return [{"error": str(e)}]
        frames = sorted(glob.glob(os.path.join(tmpdir, "frame_*.jpg")))
        if not frames:
            return [{"error": "no frames extracted — check ffmpeg/video"}]
        out=[]
        for i, fr in enumerate(frames):
            r = ocr_image(fr, langs)
            txt = r["text"].strip()
            if txt:
                m,s = divmod(i*interval, 60)
                out.append({"time": f"[{m}:{s:02d}]", "text": txt})
        return out

if __name__ == "__main__":
    ap=argparse.ArgumentParser()
    ap.add_argument("video")
    ap.add_argument("--interval", type=int, default=5)
    ap.add_argument("--langs", default="tha+eng")
    ap.add_argument("--json", action="store_true")
    args=ap.parse_args()
    sys.path.insert(0, os.path.dirname(__file__))
    from ocr import ocr_image
    res=video_ocr(args.video, args.interval, args.langs)
    if args.json: print(json.dumps(res, ensure_ascii=False, indent=2))
    else:
        for r in res:
            if "error" in r: print(f"❌ {r['error']}", file=sys.stderr); sys.exit(1)
            print(f"{r['time']} {r['text']}")

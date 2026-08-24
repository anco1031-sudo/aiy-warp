#!/usr/bin/env python3
"""
media/transcribe.py — Audio/Video transcription via ffmpeg + Whisper (Aetox 3.2)
Usage: python3 scripts/media/transcribe.py <file.mp3|mp4|wav> [--model base] [--json]
Returns: [{"time": "[m:ss]", "text": str}, ...]
Deps: ffmpeg, openai-whisper
"""
import subprocess, os, sys, tempfile, json, argparse

def transcribe(file_path: str, model_name: str = "base") -> list:
    if not os.path.exists(file_path):
        return [{"error": f"file not found: {file_path}"}]
    with tempfile.TemporaryDirectory() as tmpdir:
        wav = os.path.join(tmpdir, "audio.wav")
        try:
            subprocess.run(["ffmpeg","-y","-i",file_path,"-vn","-ar","16000","-ac","1","-c:a","pcm_s16le",wav],
                           check=True, capture_output=True)
        except subprocess.CalledProcessError as e:
            return [{"error": e.stderr.decode()[:800] if e.stderr else str(e)}]
        try:
            import whisper
        except ImportError:
            return [{"error": "openai-whisper not installed: pip install openai-whisper"}]
        model = whisper.load_model(model_name)
        result = model.transcribe(wav, language=None)
        out=[]
        for seg in result.get("segments", []):
            m,s = divmod(int(seg["start"]), 60)
            out.append({"time": f"[{m}:{s:02d}]", "text": seg["text"].strip()})
        if not out and result.get("text"):
            out.append({"time":"[0:00]","text": result["text"].strip()})
        return out

if __name__ == "__main__":
    ap=argparse.ArgumentParser()
    ap.add_argument("file")
    ap.add_argument("--model", default="base")
    ap.add_argument("--json", action="store_true")
    args=ap.parse_args()
    res=transcribe(args.file, args.model)
    if args.json: print(json.dumps(res, ensure_ascii=False, indent=2))
    else:
        for r in res:
            if "error" in r: print(f"❌ {r['error']}", file=sys.stderr); sys.exit(1)
            print(f"{r['time']} {r['text']}")

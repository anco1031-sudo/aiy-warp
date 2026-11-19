#!/usr/bin/env python3
"""
Pollinations Media CLI — image · edit · video · speech · music/sfx · 3D · transcribe.
Stdlib only (no pip installs).

  media.py image  "prompt" [-a 16:9] [--model M] [--seed N] [--out PATH]
  media.py edit   "prompt" -i ref.png [--model M]
  media.py video  "prompt" [--duration 5] [--resolution 720p] [--image first.png]
  media.py speech "สวัสดีค่ะ" [--voice nova] [--instructions "warm tone"]
  media.py music  "lo-fi piano, rain" [--duration 20] [--instrumental]
  media.py 3d     "a wooden chair" [--image chair.png] [--resolution low]
  media.py transcribe audio.mp3 [--language th]
  media.py models [--type video|audio|3d|image|text|embedding] [-q flux]
  media.py balance | doctor

Key: --key  >  $POLLINATIONS_API_KEY  >  ~/.config/opencode/secrets/pollinations.json
"""

import argparse
import base64
import json
import mimetypes
import os
import re
import ssl
import sys
import time
import unicodedata
import urllib.error
import urllib.parse
import urllib.request
import uuid
from datetime import datetime
from pathlib import Path

BASE = "https://gen.pollinations.ai"
CONFIG_PATH = Path.home() / ".config/opencode/secrets/pollinations.json"
INBOX = Path.home() / "00-Inbox/Aiy"
DEFAULT_IMAGE_MODEL = "black-forest-labs/flux.1-schnell"
DEFAULT_EDIT_MODEL = "black-forest-labs/flux.2-klein-4b"
DEFAULT_VIDEO_MODEL = "prunaai/p-video"
DEFAULT_SPEECH_MODEL = "openai/tts-1"
DEFAULT_MUSIC_MODEL = "elevenlabs/music-v2.5"
DEFAULT_3D_MODEL = "microsoft/trellis-2"

KINDS = {
    "image": (INBOX / "images", ".png"),
    "video": (INBOX / "videos", ".mp4"),
    "audio": (INBOX / "audio", ".mp3"),
    "3d": (INBOX / "models3d", ".glb"),
    "text": (INBOX / "transcripts", ".txt"),
}

ASPECTS = {
    "1:1": (1024, 1024),
    "16:9": (1344, 768),
    "9:16": (768, 1344),
    "4:3": (1152, 864),
    "3:4": (864, 1152),
    "3:2": (1216, 832),
    "2:3": (832, 1216),
    "21:9": (1536, 656),
}

# ------------------------------------------------------------------ config/key


def load_config():
    try:
        with open(CONFIG_PATH, "r", encoding="utf-8") as fh:
            return json.load(fh)
    except FileNotFoundError:
        return {}
    except json.JSONDecodeError as exc:
        warn(f"config ไม่ใช่ JSON ที่ถูกต้อง ({exc}) — ข้ามไฟล์นี้: {CONFIG_PATH}")
        return {}


def resolve_key(cli_key=None):
    for candidate in (cli_key, os.environ.get("POLLINATIONS_API_KEY"), load_config().get("api_key")):
        if candidate and str(candidate).strip():
            return str(candidate).strip()
    return ""


def require_key(cli_key=None):
    key = resolve_key(cli_key)
    if not key:
        die(
            "ต้องมี API key ก่อนค่ะ (API key required)\n"
            f"  ใส่ที่: {CONFIG_PATH}  →  {{\"api_key\": \"sk_...\"}}\n"
            "  ขอ key: https://enter.pollinations.ai/keys"
        )
    return key


def auth_headers(key):
    return {"Authorization": f"Bearer {key}"} if key else {}


# ------------------------------------------------------------------ http layer


def _request(url, data=None, headers=None, method=None, timeout=300, form=None, raw_body=None):
    """Return (status, headers, body_bytes).

    `form`     = (fields dict, file tuple) → multipart/form-data
    `raw_body` = bytes → sent as-is (caller sets Content-Type)
    """
    hdrs = {"User-Agent": "aiy-pollinations-media/2.0"}
    hdrs.update(headers or {})
    body = None
    if form is not None:
        body, ctype = encode_multipart(*form)
        hdrs["Content-Type"] = ctype
    elif raw_body is not None:
        body = raw_body
    elif data is not None:
        body = json.dumps(data).encode("utf-8")
        hdrs["Content-Type"] = "application/json"
    req = urllib.request.Request(url, data=body, headers=hdrs, method=method)
    ctx = ssl.create_default_context()
    with urllib.request.urlopen(req, timeout=timeout, context=ctx) as resp:
        return resp.status, dict(resp.headers), resp.read()


def encode_multipart(fields, file_field):
    """Minimal multipart/form-data encoder. file_field = (name, filename, bytes)."""
    boundary = f"----aiy{uuid.uuid4().hex}"
    parts = []
    for name, value in fields.items():
        if value is None:
            continue
        parts.append(
            f'--{boundary}\r\nContent-Disposition: form-data; name="{name}"\r\n\r\n{value}\r\n'.encode()
        )
    name, filename, blob = file_field
    mime = mimetypes.guess_type(filename)[0] or "application/octet-stream"
    parts.append(
        (
            f'--{boundary}\r\nContent-Disposition: form-data; name="{name}"; filename="{filename}"\r\n'
            f"Content-Type: {mime}\r\n\r\n"
        ).encode()
        + blob
        + b"\r\n"
    )
    parts.append(f"--{boundary}--\r\n".encode())
    return b"".join(parts), f"multipart/form-data; boundary={boundary}"


def api_error_message(exc):
    status = exc.code
    raw = b""
    try:
        raw = exc.read()
    except Exception:
        pass
    code, msg = "", ""
    try:
        payload = json.loads(raw.decode("utf-8", "replace"))
        err = payload.get("error") or {}
        code = err.get("code", "")
        msg = err.get("message", "")
    except Exception:
        msg = raw.decode("utf-8", "replace")[:300]

    hints = {
        400: "คำขอไม่ถูกต้อง — เช็คขนาด/ชื่อ model (bad request: check size or model)",
        401: "key ไม่ผ่าน — เช็คค่า api_key ใน secrets/pollinations.json (invalid key)",
        402: "Pollen ไม่พอ หรือ key ไม่มี budget — ใช้ model ถูกกว่า/เติมเงิน (insufficient balance)",
        403: "key ไม่มีสิทธิ์ (permission) สำหรับ model นี้ — ขอ grant ที่หน้า keys",
        404: "ไม่พบ endpoint/ไฟล์ (not found)",
        422: "เนื้อหาถูก policy บล็อก (content policy)",
        429: "ยิงถี่เกิน limit — รอแล้วลองใหม่ (rate limited)",
        500: "เซิร์ฟเวอร์ขัดข้อง (server error)",
        502: "upstream ล้ม — ลองใหม่ (upstream failed)",
        503: "บริการไม่พร้อม — ลองใหม่ (service unavailable)",
    }
    detail = f" [{code}] {msg}" if (code or msg) else ""
    return f"HTTP {status} — {hints.get(status, 'unknown error')}{detail}"


def with_retry(fn, attempts=3, base_delay=3.0, retry_status=(429, 500, 502, 503)):
    last = None
    for i in range(attempts):
        try:
            return fn()
        except urllib.error.HTTPError as exc:
            last = exc
            if exc.code not in retry_status or i == attempts - 1:
                raise
        except (urllib.error.URLError, TimeoutError, OSError) as exc:
            last = exc
            if i == attempts - 1:
                raise
        time.sleep(base_delay * (2 ** i))
    raise last


def maybe_json_error(body, ctype):
    """Binary endpoints answer JSON on failure — surface it instead of saving garbage.

    Image endpoints legitimately answer JSON too, so only bail out when the payload
    actually carries an error.
    """
    if "json" not in (ctype or "").lower():
        return
    try:
        payload = json.loads(body.decode("utf-8", "replace"))
    except json.JSONDecodeError:
        die(f"unexpected response: {body[:200]!r}")
    if isinstance(payload, dict) and (payload.get("error") or payload.get("success") is False):
        err = payload.get("error") or {}
        message = err.get("message") if isinstance(err, dict) else str(err)
        code = err.get("code", "ERROR") if isinstance(err, dict) else "ERROR"
        die(f"{code}: {message or body[:200].decode('utf-8', 'replace')}")


# ------------------------------------------------------------------ file naming


def slugify(text, maxlen=40):
    text = unicodedata.normalize("NFKC", str(text)).lower()
    text = re.sub(r"[^0-9a-z\u0e00-\u0e7f]+", "-", text).strip("-")
    return (text[:maxlen].rstrip("-") or "media")


def ext_for(ctype, kind):
    ct = (ctype or "").lower()
    table = [
        ("svg", ".svg"), ("png", ".png"), ("webp", ".webp"), ("gif", ".gif"),
        ("jpeg", ".jpg"), ("jpg", ".jpg"),
        ("mp4", ".mp4"), ("webm", ".webm"),
        ("mpeg", ".mp3"), ("mp3", ".mp3"), ("wav", ".wav"), ("ogg", ".ogg"), ("opus", ".opus"),
        ("gltf-binary", ".glb"), ("glb", ".glb"), ("octet-stream", KINDS[kind][1]),
    ]
    for needle, ext in table:
        if needle in ct:
            return ext
    return KINDS[kind][1]


def sniff_ext(blob):
    """Detect the real container from magic bytes — servers often label base64 as image/png
    even when the bytes are JPEG, and a wrong extension breaks downstream tools."""
    head = blob[:16]
    if head.startswith(b"\xff\xd8\xff"):
        return ".jpg"
    if head.startswith(b"\x89PNG"):
        return ".png"
    if head.startswith(b"RIFF") and blob[8:12] == b"WEBP":
        return ".webp"
    if head.startswith(b"GIF8"):
        return ".gif"
    if blob[4:8] == b"ftyp":
        return ".mp4"
    if head.startswith(b"glTF"):
        return ".glb"
    if head.startswith(b"ID3") or head[:2] in (b"\xff\xfb", b"\xff\xf3"):
        return ".mp3"
    if head.startswith(b"RIFF") and blob[8:12] == b"WAVE":
        return ".wav"
    if head.startswith(b"OggS"):
        return ".ogg"
    if head.lstrip()[:5].lower() in (b"<svg ", b"<?xml"):
        return ".svg"
    return None


def resolve_out(out, prompt, model, kind):
    default_dir, default_ext = KINDS[kind]
    if out:
        p = Path(out).expanduser()
        if str(out).endswith("/") or p.is_dir():
            stamp = datetime.now().strftime("%Y-%m-%d-%H%M")
            return p / f"{stamp}-{slugify(prompt)}{default_ext}"
        p.parent.mkdir(parents=True, exist_ok=True)
        return p
    default_dir.mkdir(parents=True, exist_ok=True)
    stamp = datetime.now().strftime("%Y-%m-%d-%H%M")
    short = (model or kind).split("/")[-1].replace(".", "").replace(":", "-")[:24]
    return default_dir / f"{stamp}-{slugify(prompt)}-{short}{default_ext}"


def save_media(path, blob, ctype, kind):
    path = Path(path)
    real = sniff_ext(blob) or ext_for(ctype, kind)
    if path.suffix.lower() != real:
        path = path.with_suffix(real)
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_bytes(blob)
    return path


def data_uri(path):
    p = Path(path).expanduser()
    if not p.is_file():
        die(f"ไม่พบไฟล์: {path} (file not found)")
    mime = mimetypes.guess_type(p.name)[0] or "image/png"
    return f"data:{mime};base64,{base64.b64encode(p.read_bytes()).decode()}"


def parse_size(value):
    m = re.fullmatch(r"\s*(\d{2,5})\s*[xX]\s*(\d{2,5})\s*", value or "")
    if not m:
        die(f"ขนาดต้องเป็นแบบ 1024x1024 (got: {value})")
    return int(m.group(1)), int(m.group(2))


def size_args(args):
    if args.size:
        return parse_size(args.size)
    return ASPECTS.get(args.aspect or "1:1", ASPECTS["1:1"])


# ------------------------------------------------------------------ reporting


def report(args, *, model, path, blob, headers, kind, extra=None):
    info = {
        "file": str(path),
        "kind": kind,
        "model": model,
        "bytes": len(blob),
        # JSON APIs report application/json even when the payload is an image —
        # the saved extension is the truthful signal.
        "content_type": mimetypes.guess_type(str(path))[0] or headers.get("Content-Type", ""),
    }
    for key in ("X-Image-Dimensions", "X-Seed", "X-Video-Duration", "X-Media-Id"):
        if headers.get(key):
            info[key.lower().replace("x-", "").replace("-", "_")] = headers[key]
    if extra:
        info.update(extra)
    if args.json:
        print(json.dumps(info, indent=2))
        return
    print(f"saved  : {path}  ({len(blob) / 1024:.0f} KB)")
    print(f"model  : {model}")
    if info.get("image_dimensions"):
        print(f"size   : {info['image_dimensions']}")
    if info.get("seed"):
        print(f"seed   : {info['seed']}  (ใช้ --seed เดิมเพื่อแก้ต่อเนื่อง)")
    for key, value in (extra or {}).items():
        print(f"{key:7s}: {value}")


# ------------------------------------------------------------------ image / edit


def cmd_image(args):
    key = require_key(args.key)
    model = args.model or load_config().get("default_model") or DEFAULT_IMAGE_MODEL
    guard_free(model, "image", args)
    w, h = size_args(args)
    payload = {
        "model": model,
        "prompt": args.prompt,
        "size": f"{w}x{h}",
        "response_format": "b64_json",
    }
    if args.seed is not None:
        payload["seed"] = args.seed
    if args.transparent:
        payload["background"] = "transparent"
    if args.quality:
        payload["quality"] = args.quality
    if args.safe:
        payload["safe"] = args.safe
    if args.image:
        payload["image"] = [data_uri(p) if not str(p).startswith(("http", "data:")) else p for p in args.image]

    out = resolve_out(args.out, args.prompt, model, "image")
    headers, blob, ctype = post_image(f"{BASE}/v1/images/generations", payload, key)
    path = save_media(out, blob, ctype, "image")
    report(args, model=model, path=path, blob=blob, headers=headers, kind="image")


def cmd_edit(args):
    """Image-to-image through /v1/images/generations (the `image` extension).

    One code path for every image-input model — gpt-image/flux-kontext (paid tier)
    and flux.2-klein / community models (Quest tier) all accept it.
    """
    key = require_key(args.key)
    model = args.model or DEFAULT_EDIT_MODEL
    guard_free(model, "image", args)
    images = [
        data_uri(p) if not str(p).startswith(("http", "data:")) else p for p in args.image
    ]
    payload = {
        "prompt": args.prompt,
        "model": model,
        "image": images[0] if len(images) == 1 else images,
        "response_format": "b64_json",
    }
    if args.size:
        payload["size"] = args.size
    elif args.aspect:
        w, h = ASPECTS.get(args.aspect, ASPECTS["1:1"])
        payload["size"] = f"{w}x{h}"
    if args.seed is not None:
        payload["seed"] = args.seed

    out = resolve_out(args.out, args.prompt, model, "image")
    headers, blob, ctype = post_image(f"{BASE}/v1/images/generations", payload, key)
    path = save_media(out, blob, ctype, "image")
    report(args, model=model, path=path, blob=blob, headers=headers, kind="image")


def post_image(url, payload, key):
    def call():
        _, headers, body = _request(url, data=payload, headers=auth_headers(key))
        ctype = headers.get("Content-Type", "")
        maybe_json_error(body, ctype)
        data = json.loads(body.decode("utf-8", "replace"))
        item = (data.get("data") or [{}])[0]
        if item.get("b64_json"):
            return headers, base64.b64decode(item["b64_json"]), "image/png"
        if item.get("url"):
            _, h2, blob = _request(item["url"], timeout=180)
            return headers, blob, h2.get("Content-Type", "image/png")
        die(f"no image in response: {json.dumps(data)[:300]}")

    try:
        return with_retry(call)
    except urllib.error.HTTPError as exc:
        die(api_error_message(exc))
    except Exception as exc:
        die(f"image request failed: {exc}")


# ------------------------------------------------------------------ video / 3d / audio


def binary_get(url, key, kind, timeout=900):
    def call():
        _, headers, body = _request(url, headers=auth_headers(key), timeout=timeout)
        ctype = headers.get("Content-Type", "")
        maybe_json_error(body, ctype)
        return headers, body, ctype

    try:
        return with_retry(call, attempts=2, base_delay=5.0)
    except urllib.error.HTTPError as exc:
        die(api_error_message(exc))
    except Exception as exc:
        die(f"{kind} request failed: {exc}")


def cmd_video(args):
    key = require_key(args.key)
    model = args.model or DEFAULT_VIDEO_MODEL
    guard_free(model, "video", args)
    params = {"model": model}
    if args.duration:
        params["duration"] = args.duration
    if args.resolution:
        params["resolution"] = args.resolution
    if args.aspect:
        params["aspectRatio"] = args.aspect
    if args.seed is not None:
        params["seed"] = args.seed
    if args.audio:
        params["audio"] = "true"
    if args.safe:
        params["safe"] = args.safe
    refs = []
    for p in args.image or []:
        refs.append(p if str(p).startswith(("http", "data:")) else upload(p, key))
    if refs:
        params["image"] = ",".join(refs)

    url = f"{BASE}/video/{urllib.parse.quote(args.prompt, safe='')}?{urllib.parse.urlencode(params)}"
    out = resolve_out(args.out, args.prompt, model, "video")
    headers, blob, ctype = binary_get(url, key, "video")
    if len(blob) < 5000:
        die(f"ไฟล์วิดีโอเล็กผิดปกติ ({len(blob)} bytes) — ตรวจ prompt/model")
    path = save_media(out, blob, ctype, "video")
    report(args, model=model, path=path, blob=blob, headers=headers, kind="video",
           extra={"duration": args.duration} if args.duration else None)


def cmd_3d(args):
    key = require_key(args.key)
    model = args.model or DEFAULT_3D_MODEL
    guard_free(model, "3d", args)
    params = {"model": model}
    if args.resolution:
        params["resolution"] = args.resolution
    if args.seed is not None:
        params["seed"] = args.seed
    if args.image:
        params["image"] = upload(args.image[0], key) if not str(args.image[0]).startswith(("http", "data:")) else args.image[0]
    elif model.startswith(("microsoft/", "nvidia/")):
        die(f"{model} ต้องมีภาพอ้างอิง — ใส่ -i ภาพ.png (this model needs a reference image)")

    url = f"{BASE}/3d/{urllib.parse.quote(args.prompt or 'model', safe='')}?{urllib.parse.urlencode(params)}"
    out = resolve_out(args.out, args.prompt or "model", model, "3d")
    headers, blob, ctype = binary_get(url, key, "3d")
    if len(blob) < 1000:
        die(f"ไฟล์ 3D เล็กผิดปกติ ({len(blob)} bytes)")
    path = save_media(out, blob, ctype, "3d")
    report(args, model=model, path=path, blob=blob, headers=headers, kind="3d")


def cmd_speech(args):
    key = require_key(args.key)
    model = args.model or DEFAULT_SPEECH_MODEL
    guard_free(model, "audio", args)
    payload = {
        "model": model,
        "input": args.text,
        "voice": args.voice or "nova",
        "response_format": args.format or "mp3",
    }
    if args.instructions:
        payload["instructions"] = args.instructions
    if args.speed:
        payload["speed"] = args.speed

    out = resolve_out(args.out, args.text, model, "audio")

    def call():
        _, headers, body = _request(f"{BASE}/v1/audio/speech", data=payload, headers=auth_headers(key))
        ctype = headers.get("Content-Type", "")
        maybe_json_error(body, ctype)
        return headers, body, ctype

    try:
        headers, blob, ctype = with_retry(call)
    except urllib.error.HTTPError as exc:
        die(api_error_message(exc))
    except Exception as exc:
        die(f"speech failed: {exc}")

    path = save_media(out, blob, ctype, "audio")
    report(args, model=model, path=path, blob=blob, headers=headers, kind="audio",
           extra={"voice": payload["voice"]})


def cmd_music(args):
    key = require_key(args.key)
    model = args.model or DEFAULT_MUSIC_MODEL
    guard_free(model, "audio", args)
    params = {"model": model}
    if args.duration:
        params["duration"] = args.duration
    if args.instrumental:
        params["instrumental"] = "true"
    if args.steps:
        params["steps"] = args.steps
    if args.negative:
        params["negative_prompt"] = args.negative
    if args.seed is not None:
        params["seed"] = args.seed
    if args.loop:
        params["loop"] = "true"
    if args.influence is not None:
        params["prompt_influence"] = args.influence

    url = f"{BASE}/audio/{urllib.parse.quote(args.prompt, safe='')}?{urllib.parse.urlencode(params)}"
    out = resolve_out(args.out, args.prompt, model, "audio")
    headers, blob, ctype = binary_get(url, key, "audio", timeout=600)
    if len(blob) < 1000:
        die(f"ไฟล์เสียงเล็กผิดปกติ ({len(blob)} bytes) — ตรวจ prompt/model")
    path = save_media(out, blob, ctype, "audio")
    report(args, model=model, path=path, blob=blob, headers=headers, kind="audio",
           extra={"duration": args.duration} if args.duration else None)


def cmd_transcribe(args):
    key = require_key(args.key)
    src = Path(args.file).expanduser()
    if not src.is_file():
        die(f"ไม่พบไฟล์เสียง: {args.file}")
    model = args.model or "openai/gpt-transcribe"
    fields = {"model": model, "response_format": args.format or "json"}
    if args.language:
        fields["language"] = args.language
    if args.prompt:
        fields["prompt"] = args.prompt

    def call():
        _, headers, body = _request(
            f"{BASE}/v1/audio/transcriptions",
            headers=auth_headers(key),
            form=(fields, ("file", src.name, src.read_bytes())),
            timeout=600,
        )
        maybe_json_error(body, headers.get("Content-Type", ""))
        return headers, body

    try:
        headers, body = with_retry(call)
    except urllib.error.HTTPError as exc:
        die(api_error_message(exc))
    except Exception as exc:
        die(f"transcription failed: {exc}")

    text = body.decode("utf-8", "replace")
    if (args.format or "json") == "json":
        try:
            data = json.loads(text)
            text = data.get("text", text)
        except json.JSONDecodeError:
            pass

    if args.out:
        path = Path(args.out).expanduser()
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(text, encoding="utf-8")
        print(f"saved  : {path}")
    print(f"model  : {model}")
    print("---")
    print(text.strip())


def upload(path, key):
    """Upload a local file to media.pollinations.ai and return its public URL."""
    src = Path(path).expanduser()
    if not src.is_file():
        die(f"ไม่พบไฟล์: {path}")
    mime = mimetypes.guess_type(src.name)[0] or "application/octet-stream"
    try:
        _, headers, body = _request(
            "https://media.pollinations.ai/upload",
            headers={**auth_headers(key), "Content-Type": mime, "X-File-Name": src.name},
            raw_body=src.read_bytes(),
            method="POST",
            timeout=300,
        )
    except urllib.error.HTTPError as exc:
        die(f"อัปโหลดไม่สำเร็จ: {api_error_message(exc)}")
    try:
        return json.loads(body.decode())["url"]
    except Exception:
        die(f"อัปโหลดไม่สำเร็จ: {body[:200]!r}")


# ------------------------------------------------------------------ discovery


def fetch_models(kind, key):
    """Full public catalog from /models — it carries `category`, `pricing` and `paid_only`,
    which is the authoritative Quest-vs-paid signal. Authenticated per-type endpoints filter
    by reliability/permissions and hide whole categories, so they are only a fallback."""
    headers = {"Accept": "application/json"}
    try:
        _, _, body = _request(f"{BASE}/models", headers=headers, timeout=90)
        models = json.loads(body.decode("utf-8", "replace"))
    except Exception:
        endpoint = {"image": "/image/models", "audio": "/audio/models", "text": "/text/models"}.get(
            kind, "/models"
        )
        _, _, body = _request(f"{BASE}{endpoint}", headers=headers, timeout=90)
        models = json.loads(body.decode("utf-8", "replace"))
    return [m for m in models if m.get("category") == kind]


_TIER_CACHE = {}


def model_is_paid(model, kind):
    """True when the model needs a paid top-up. Unknown/offline → False (never block a run)."""
    if kind not in _TIER_CACHE:
        try:
            _TIER_CACHE[kind] = {m["name"]: bool(m.get("paid_only")) for m in fetch_models(kind, "")}
        except Exception:
            _TIER_CACHE[kind] = {}
    return _TIER_CACHE[kind].get(model, False)


def guard_free(model, kind, args):
    """--free: refuse paid models before spending a request (default posture for this account)."""
    if getattr(args, "free", False) and model_is_paid(model, kind):
        die(
            f"{model} ต้องใช้ paid balance — โหมด --free ปฏิเสธก่อนยิง\n"
            f"  ใช้ตัวฟรีแทน: python3 $M models -t {kind} --quest\n"
            f"  (ถ้าจำเป็นต้องใช้จริง: เอา --free ออก แล้วบอกหลุยส์ว่าราคาเท่าไร)"
        )


def cmd_probe(args):
    """Report which models this key can actually run — Quest balance vs paid top-up.

    Failures are free (402 is rejected before generation); a successful probe costs
    the model's minimum, so the default is a small 512×512 image.
    """
    key = require_key(args.key)
    models = fetch_models("image", key)
    if args.query:
        q = args.query.lower()
        models = [m for m in models if q in m["name"].lower() or q in m.get("title", "").lower()]
    models = models[: args.limit or 12]

    print(f"probing {len(models)} image models at {args.size} (failures are free)\n")
    for m in models:
        name = m["name"]
        payload = {
            "model": name,
            "prompt": "a red apple on a white table",
            "size": args.size,
            "response_format": "b64_json",
        }
        try:
            _request(f"{BASE}/v1/images/generations", data=payload, headers=auth_headers(key), timeout=180)
            verdict = "QUEST ok"
        except urllib.error.HTTPError as exc:
            raw = ""
            try:
                raw = exc.read().decode("utf-8", "replace")
            except Exception:
                pass
            if "paid balance" in raw:
                verdict = "PAID only (top-up)"
            elif exc.code == 400:
                verdict = "BAD-ARGS (size/params)"
            else:
                verdict = f"HTTP {exc.code}"
        except Exception as exc:
            verdict = f"ERR {type(exc).__name__}"
        print(f"{name:48s} {verdict}")


def cmd_models(args):
    key = resolve_key(args.key)
    kind = args.type or "image"
    models = fetch_models(kind, key)

    if args.quest:
        models = [m for m in models if not m.get("paid_only")]
    if args.query:
        q = args.query.lower()
        models = [
            m for m in models
            if q in m.get("name", "").lower()
            or q in m.get("title", "").lower()
            or any(q in a.lower() for a in m.get("aliases", []))
        ]
    total = len(models)
    quest_total = sum(1 for m in models if not m.get("paid_only"))
    if args.limit:
        models = models[: args.limit]

    if args.json:
        print(json.dumps(models, indent=2))
        return

    print(f"{total} {kind} models — {quest_total} ใช้ Quest ได้ / {total - quest_total} ต้อง paid\n")
    for m in models:
        tier = "PAID " if m.get("paid_only") else "QUEST"
        pricing = {k: v for k, v in m.get("pricing", {}).items() if k != "currency"}
        bits = [tier, m.get("name", ""), m.get("title", "")[:22]]
        if m.get("input_modalities"):
            bits.append("in=" + ",".join(m["input_modalities"]))
        if pricing:
            bits.append(str(pricing)[:48])
        if m.get("min_duration") or m.get("max_duration"):
            bits.append(f"dur {m.get('min_duration')}-{m.get('max_duration')}s")
        if m.get("resolutions"):
            bits.append("res=" + ",".join(m["resolutions"]))
        print(" · ".join(str(b) for b in bits))
        if args.voices and m.get("voices"):
            print("   voices:", ", ".join(v if isinstance(v, str) else v.get("name", "?") for v in m["voices"]))


def cmd_balance(args):
    key = require_key(args.key)
    try:
        _, _, body = _request(f"{BASE}/account/balance", headers=auth_headers(key), timeout=30)
    except urllib.error.HTTPError as exc:
        die(
            api_error_message(exc)
            + "\n  → key ยังไม่มีสิทธิ์ account:usage · เปิดที่ "
              "https://enter.pollinations.ai/grant?id=dX0iiIfxsx81OHLGw0BKe4bduVajsrmd&permission=usage\n"
              "  → ระหว่างนี้ดูยอดใช้จริงได้ด้วย: python3 $M usage"
        )
    print(json.dumps(json.loads(body.decode()), indent=2))


def cmd_usage(args):
    """Spend history for this key — the only quota signal that works without extra permissions.

    Pollen is NOT unlimited: free requests draw from the `tier` bucket (earned via quests),
    premium models need the `paid` bucket. Each row here is one billed request.
    """
    key = require_key(args.key)
    try:
        _, _, body = _request(f"{BASE}/account/key/usage", headers=auth_headers(key), timeout=30)
    except urllib.error.HTTPError as exc:
        die(api_error_message(exc))
    rows = json.loads(body.decode("utf-8", "replace")).get("usage") or []

    if args.json:
        print(json.dumps(rows, indent=2))
        return
    if not rows:
        print("ยังไม่มีประวัติการใช้ (no usage yet)")
        return

    by_day, by_model, by_bucket = {}, {}, {}
    for r in rows:
        day = (r.get("timestamp") or "")[:10]
        cost = float(r.get("cost_usd") or 0)
        by_day[day] = by_day.get(day, 0) + cost
        label = f"{r.get('type', '?')} · {r.get('model') or '-'}"
        by_model[label] = by_model.get(label, 0) + cost
        b = r.get("meter_source") or "?"
        by_bucket[b] = by_bucket.get(b, 0) + cost

    total = sum(by_day.values())
    print(f"{len(rows)} requests · ใช้ไป ${total:.4f} pollen")
    print("bucket:", ", ".join(f"{k}=${v:.4f}" for k, v in by_bucket.items()))
    print("\nตามวัน (ใหม่ → เก่า):")
    for day in sorted(by_day, reverse=True)[: args.days]:
        n = sum(1 for r in rows if (r.get("timestamp") or "").startswith(day))
        print(f"  {day}  {n:3d} req  ${by_day[day]:.4f}")
    print("\nตาม model (มาก → น้อย):")
    for label, cost in sorted(by_model.items(), key=lambda kv: -kv[1])[: args.limit]:
        print(f"  ${cost:.4f}  {label}")
    print(
        "\n⚠️ ยอดคงเหลือที่เหลือจริงต้องใช้สิทธิ์ account:usage "
        "(`$M balance`) — ถ้าไม่ให้สิทธิ์ ให้ดูตัวเลขจากข้อความ 402 แทน"
    )


def cmd_doctor(args):
    key = resolve_key(args.key)
    cfg = load_config()
    print(f"config      : {CONFIG_PATH} {'(exists)' if CONFIG_PATH.exists() else '(missing)'}")
    if key:
        masked = f"{key[:6]}…{key[-4:]}" if len(key) > 12 else "set"
        print(f"api key     : {masked}")
        try:
            _, _, body = _request(f"{BASE}/account/key", headers=auth_headers(key), timeout=30)
            info = json.loads(body.decode())
            print(f"key status  : valid={info.get('valid')} type={info.get('type')} name={info.get('name')}")
        except urllib.error.HTTPError as exc:
            print(f"key status  : {api_error_message(exc)}")
    else:
        print("api key     : NOT SET → ใส่ที่ secrets/pollinations.json (ไม่มี key = สร้างไม่ได้)")

    for label, path in (("images", INBOX / "images"), ("videos", INBOX / "videos"),
                        ("audio", INBOX / "audio"), ("models3d", INBOX / "models3d")):
        print(f"{label:12s}: {path} {'✓' if path.exists() else '(จะสร้างให้)'}")
    try:
        _, _, body = _request(f"{BASE}/models", timeout=30)
        total = len(json.loads(body.decode()))
        print(f"connectivity: OK ({total} models reachable)")
    except Exception as exc:
        print(f"connectivity: FAILED — {exc}")


# ------------------------------------------------------------------ plumbing


def warn(message):
    print(f"! {message}", file=sys.stderr)


def die(message):
    print(f"✖ {message}", file=sys.stderr)
    sys.exit(1)


def build_parser():
    ap = argparse.ArgumentParser(
        prog="media.py",
        description="Pollinations media: image · edit · video · speech · music · 3D · transcribe",
    )
    ap.add_argument("--key", help="API key (default: env POLLINATIONS_API_KEY or config file)")
    ap.add_argument("--json", action="store_true", help="machine-readable result")
    ap.add_argument("--free", action="store_true",
                    help="ปฏิเสธ model ที่ต้อง paid balance ก่อนยิง (โหมดประหยัดของบัญชีนี้)")
    sub = ap.add_subparsers(dest="cmd", required=True)

    def add_json(p):
        """Allow --json / --free after the subcommand too (SUPPRESS keeps the global default)."""
        p.add_argument("--json", action="store_true", default=argparse.SUPPRESS,
                       help="machine-readable result")
        p.add_argument("--free", action="store_true", default=argparse.SUPPRESS,
                       help="ปฏิเสธ model ที่ต้อง paid balance ก่อนยิง")

    def add_common(p, size=True):
        p.add_argument("--model", "-m")
        p.add_argument("--seed", type=int)
        p.add_argument("--out", "-o")
        p.add_argument("--aspect", "-a", choices=sorted(ASPECTS))
        if size:
            p.add_argument("--size", "-s", help="WxH e.g. 1024x1024")

    i = sub.add_parser("image", aliases=["img"], help="text → image")
    i.add_argument("prompt")
    add_common(i)
    i.add_argument("--transparent", action="store_true")
    i.add_argument("--quality")
    i.add_argument("--safe")
    i.add_argument("--image", action="append", help="reference image (local path or URL)")
    i.set_defaults(func=cmd_image)

    e = sub.add_parser("edit", help="edit / restyle an image")
    e.add_argument("prompt")
    e.add_argument("--image", "-i", action="append", required=True)
    e.add_argument("--model", "-m")
    e.add_argument("--size", "-s")
    e.add_argument("--aspect", "-a", choices=sorted(ASPECTS))
    e.add_argument("--seed", type=int)
    e.add_argument("--out", "-o")
    e.set_defaults(func=cmd_edit)

    v = sub.add_parser("video", help="text → video (mp4)")
    v.add_argument("prompt")
    add_common(v, size=False)
    v.add_argument("--duration", "-d", type=int)
    v.add_argument("--resolution", "-r")
    v.add_argument("--image", action="append", help="start frame (local path or URL)")
    v.add_argument("--audio", action="store_true", help="generate audio track (veo etc.)")
    v.add_argument("--safe")
    v.set_defaults(func=cmd_video)

    s = sub.add_parser("speech", help="text → speech (tts)")
    s.add_argument("text")
    s.add_argument("--model", "-m")
    s.add_argument("--voice", "-v")
    s.add_argument("--format", "-f", help="mp3 | wav | opus …")
    s.add_argument("--instructions", help="tone/emotion hint (Gemini/Qwen voices)")
    s.add_argument("--speed", type=float)
    s.add_argument("--out", "-o")
    s.set_defaults(func=cmd_speech)

    mu = sub.add_parser("music", help="prompt → music / sound effect")
    mu.add_argument("prompt")
    mu.add_argument("--model", "-m")
    mu.add_argument("--duration", "-d", type=int)
    mu.add_argument("--instrumental", action="store_true")
    mu.add_argument("--negative")
    mu.add_argument("--steps", type=int)
    mu.add_argument("--loop", action="store_true")
    mu.add_argument("--influence", type=float)
    mu.add_argument("--seed", type=int)
    mu.add_argument("--out", "-o")
    mu.set_defaults(func=cmd_music)

    t = sub.add_parser("3d", help="prompt/image → 3D model (glb)")
    t.add_argument("prompt", nargs="?", default="")
    t.add_argument("--model", "-m")
    t.add_argument("--image", "-i", action="append")
    t.add_argument("--resolution", "-r", choices=["low", "medium", "high"])
    t.add_argument("--seed", type=int)
    t.add_argument("--out", "-o")
    t.set_defaults(func=cmd_3d)

    tr = sub.add_parser("transcribe", help="audio → text")
    tr.add_argument("file")
    tr.add_argument("--model", "-m")
    tr.add_argument("--language", "-l")
    tr.add_argument("--prompt")
    tr.add_argument("--format", "-f", help="json | text | srt | vtt")
    tr.add_argument("--out", "-o")
    tr.set_defaults(func=cmd_transcribe)

    m = sub.add_parser("models", help="list models with pricing")
    m.add_argument("--type", "-t", default="image",
                   choices=["image", "video", "audio", "3d", "text", "embedding"])
    m.add_argument("--query", "-q")
    m.add_argument("--limit", "-n", type=int, default=0)
    m.add_argument("--quest", action="store_true", help="แสดงเฉพาะ model ที่ใช้ Quest balance ได้")
    m.add_argument("--voices", action="store_true", help="show voice lists (audio models)")
    m.set_defaults(func=cmd_models)

    p = sub.add_parser("probe", help="which image models this key can run (Quest vs paid)")
    p.add_argument("--query", "-q")
    p.add_argument("--limit", "-n", type=int, default=12)
    p.add_argument("--size", "-s", default="512x512")
    p.set_defaults(func=cmd_probe)

    b = sub.add_parser("balance", help="show remaining balance (needs account:usage grant)")
    b.set_defaults(func=cmd_balance)

    u = sub.add_parser("usage", help="spend history — ดูว่าใช้ pollen ไปเท่าไร")
    u.add_argument("--days", type=int, default=7)
    u.add_argument("--limit", "-n", type=int, default=8)
    u.set_defaults(func=cmd_usage)

    d = sub.add_parser("doctor", help="check key, permissions, folders, connectivity")
    d.set_defaults(func=cmd_doctor)

    # --json must work before OR after the subcommand (aliases share the same parser object)
    seen = set()
    for parser in sub.choices.values():
        if id(parser) in seen:
            continue
        seen.add(id(parser))
        add_json(parser)
    return ap


def main():
    args = build_parser().parse_args()
    args.func(args)


if __name__ == "__main__":
    main()

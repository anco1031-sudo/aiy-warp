#!/usr/bin/env python3
"""
check-chars.py — ตรวจ "อักษรต่างภาษาตกค้าง" ในงานภาษาไทย (Pre-Send Quality Gate ข้อ 1)

ที่มา: เคสจริง 29-Oct-26 — อักษรจีน/รัสเซียหลุดกลางประโยคไทย 4 ครั้งใน 1 วัน
      และ `grep "[๐-๙]"` ให้ผลผิดเพราะ bracket range จับที่ระดับ byte
      → ตัวนี้ไล่ที่ระดับ code point (แม่นกว่า grep ทุกกรณี)

Usage:
  python3 check-chars.py <path...>              # ไฟล์หรือโฟลเดอร์ (ไล่ .md/.txt/.py/.sh/.yaml/.json)
  python3 check-chars.py <path> --skip-code     # ข้ามเนื้อหาใน ``` code fence ``` (เช่นเอกสารสอน)
  python3 check-chars.py <path> --quiet         # เฉพาะสรุป (ไม่พิมพ์รายบรรทัด)
  python3 check-chars.py <path> --allow "อักษร"  # ยกเว้นคำที่ตั้งใจอ้างอิง

Exit codes: 0 = สะอาด · 1 = เจออักษรต่างภาษา · 2 = usage error
"""
from __future__ import annotations

import argparse
import os
import sys

# (ชื่อเรียก, ช่วง code point)
SCRIPTS: list[tuple[str, list[tuple[int, int]]]] = [
    ("จีน (Han)", [(0x3400, 0x4DBF), (0x4E00, 0x9FFF), (0xF900, 0xFAFF)]),
    ("ญี่ปุ่น (Kana)", [(0x3040, 0x309F), (0x30A0, 0x30FF)]),
    ("เกาหลี (Hangul)", [(0x1100, 0x11FF), (0xAC00, 0xD7AF)]),
    ("รัสเซีย (Cyrillic)", [(0x0400, 0x04FF)]),
    ("อาหรับ (Arabic)", [(0x0600, 0x06FF), (0x0750, 0x077F)]),
    ("ฮีบรู (Hebrew)", [(0x0590, 0x05FF)]),
    ("กรีก (Greek)", [(0x0370, 0x03FF)]),
    ("เดวานาครี (Devanagari)", [(0x0900, 0x097F)]),
]

EXTS = {".md", ".txt", ".py", ".sh", ".yaml", ".yml", ".json", ".csv", ".html"}
SKIP_DIRS = {".git", "node_modules", "__pycache__", ".obsidian", "snapshots"}


def script_of(ch: str) -> str | None:
    cp = ord(ch)
    for name, ranges in SCRIPTS:
        for lo, hi in ranges:
            if lo <= cp <= hi:
                return name
    return None


def strip_code_fences(text: str) -> str:
    """แทนเนื้อหาใน code fence ด้วยช่องว่าง (คงเลขบรรทัดไว้)"""
    out: list[str] = []
    in_fence = False
    for line in text.splitlines(keepends=True):
        if line.lstrip().startswith("```"):
            in_fence = not in_fence
            out.append(line)
            continue
        out.append(("" if in_fence else line))
    return "".join(out)


def iter_files(paths: list[str]):
    for p in paths:
        if os.path.isfile(p):
            yield p
        elif os.path.isdir(p):
            for root, dirs, files in os.walk(p):
                dirs[:] = [d for d in dirs if d not in SKIP_DIRS]
                for f in sorted(files):
                    if os.path.splitext(f)[1].lower() in EXTS:
                        yield os.path.join(root, f)
        else:
            print(f"⚠️  ไม่พบ path: {p}", file=sys.stderr)


def main() -> int:
    ap = argparse.ArgumentParser(add_help=True)
    ap.add_argument("paths", nargs="+")
    ap.add_argument("--skip-code", action="store_true", help="ข้ามเนื้อหาใน ``` code fence ```")
    ap.add_argument("--quiet", action="store_true", help="พิมพ์เฉพาะสรุป")
    ap.add_argument("--allow", action="append", default=[], help="คำ/ตัวอักษรที่ตั้งใจใช้ (ยกเว้น)")
    ap.add_argument("--allow-script", action="append", default=[], help="ข้ามทั้งตระกูลอักษร เช่น \"กรีก (Greek)\" (สัญลักษณ์คณิตศาสตร์)")
    ap.add_argument("--max-per-file", type=int, default=8)
    args = ap.parse_args()

    allow = set("".join(args.allow))
    allow_scripts = set(args.allow_script)
    total_files = 0
    dirty_files: list[tuple[str, int]] = []
    grand_total = 0

    for path in iter_files(args.paths):
        total_files += 1
        try:
            text = open(path, encoding="utf-8", errors="strict").read()
        except (UnicodeDecodeError, OSError):
            continue
        if args.skip_code:
            text = strip_code_fences(text)
        if not text:
            continue

        hits: list[tuple[int, str, str, str]] = []  # line, script, char, snippet
        for lineno, line in enumerate(text.splitlines(), 1):
            for pos, ch in enumerate(line):
                if ch in allow:
                    continue
                name = script_of(ch)
                if name and name in allow_scripts:
                    continue
                if name:
                    start = max(0, pos - 18)
                    end = min(len(line), pos + 18)
                    snippet = line[start:end].replace("\t", " ")
                    hits.append((lineno, name, ch, snippet))

        if hits:
            dirty_files.append((path, len(hits)))
            grand_total += len(hits)
            if not args.quiet:
                print(f"\n❌ {path}  ({len(hits)} ตัว)")
                for lineno, name, ch, snippet in hits[: args.max_per_file]:
                    print(f"   บรรทัด {lineno} · {name} · U+{ord(ch):04X} · …{snippet}…")
                if len(hits) > args.max_per_file:
                    print(f"   … อีก {len(hits) - args.max_per_file} ตัว")

    print("\n" + "=" * 58)
    if grand_total == 0:
        print(f"✅ สะอาด — สแกน {total_files} ไฟล์ ไม่พบอักษรต่างภาษาตกค้าง")
        return 0
    print(f"❌ พบ {grand_total} ตัว ใน {len(dirty_files)} ไฟล์ (สแกนทั้งหมด {total_files} ไฟล์)")
    for path, n in dirty_files[:20]:
        print(f"   - {path} ({n})")
    if len(dirty_files) > 20:
        print(f"   … อีก {len(dirty_files) - 20} ไฟล์")
    print("\n💡 แก้: ลบอักษรนั้น หรือถ้าตั้งใจอ้างอิง → ใส่ --allow หรือเขียนเป็น code point (เช่น U+7BC4)")
    return 1


if __name__ == "__main__":
    sys.exit(main())

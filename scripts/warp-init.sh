#!/usr/bin/env bash
# warp-init.sh — ตรวจสอบ dependencies สำหรับ aiy-warp media/office/browser
set -euo pipefail
echo "🔍 Aiy Warp — Dependency Check"
echo "================================"
check() { command -v "$1" >/dev/null 2>&1 && echo "✅ $1 $( $1 --version 2>&1 | head -n1)" || echo "❌ $1 missing"; }
check tesseract
check pdftotext
check ffmpeg
echo "--- python ---"
python3 -c "import whisper; print('✅ whisper', whisper.__version__)" 2>&1 || echo "❌ whisper (pip install openai-whisper)"
python3 -c "import pymupdf; print('✅ pymupdf', pymupdf.__version__)" 2>&1 || echo "⚠️ pymupdf missing"
python3 -c "import pdfplumber; print('✅ pdfplumber')" 2>&1 || echo "⚠️ pdfplumber missing"
python3 -c "import openpyxl; print('✅ openpyxl', openpyxl.__version__)" 2>&1 || echo "❌ openpyxl (pip install openpyxl)"
python3 -c "import docx; print('✅ python-docx')" 2>&1 || echo "❌ python-docx (pip install python-docx)"
python3 -c "import pptx; print('✅ python-pptx')" 2>&1 || echo "❌ python-pptx (pip install python-pptx)"
python3 -c "import playwright; print('✅ playwright')" 2>&1 || echo "❌ playwright (pip install playwright && playwright install chromium)"
echo "================================"
echo "ทิป: pip install openpyxl python-docx python-pptx playwright && playwright install chromium"

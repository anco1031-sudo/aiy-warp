#!/usr/bin/env bash
# warp-init.sh — ตรวจสอบ dependencies สำหรับ aiy-warp media/office/browser
set -euo pipefail
echo "🔍 Aiy Warp — Dependency Check"
echo "================================"
check() { command -v "$1" >/dev/null 2>&1 && echo "✅ $1 $( $1 --version 2>&1 | head -n1)" || echo "❌ $1 missing"; }
check tesseract
check pdftotext
check ffmpeg
echo "--- node / browser ---"
check playwright
check chromium
if pacman -Qs playwright >/dev/null 2>&1; then echo "✅ playwright (yay/node) $(playwright --version 2>&1 | head -n1)"; else echo "❌ playwright node (yay -S playwright)"; fi
echo "--- python ---"
python3 -c "import whisper; print('✅ whisper', whisper.__version__)" 2>&1 || echo "❌ whisper (pip install openai-whisper)"
python3 -c "import pymupdf; print('✅ pymupdf', pymupdf.__version__)" 2>&1 || echo "⚠️ pymupdf missing"
python3 -c "import pdfplumber; print('✅ pdfplumber')" 2>&1 || echo "⚠️ pdfplumber missing"
python3 -c "import openpyxl; print('✅ openpyxl', openpyxl.__version__)" 2>&1 || echo "❌ openpyxl (pip install openpyxl)"
python3 -c "import docx; print('✅ python-docx')" 2>&1 || echo "❌ python-docx (pip install python-docx)"
python3 -c "import pptx; print('✅ python-pptx')" 2>&1 || echo "❌ python-pptx (pip install python-pptx)"
python3 -c "import playwright; print('✅ playwright (python)')" 2>&1 || echo "ℹ️ playwright python ไม่บังคับ — ใช้ Node playwright ผ่าน MCP ได้แล้ว (yay ติดตั้งแล้ว ✅)"
echo "================================"
echo "MCP browser: ใช้ yay playwright + chromium ได้เลย (config/opencode.jsonc เปิดไว้แล้ว)"
echo "Python playwright (optional): pip install playwright && playwright install chromium"

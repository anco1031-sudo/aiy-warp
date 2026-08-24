---
name: browser-control
description: "Aiy's Browser Control — automate browser via Playwright MCP (preferred) + standalone scripts/browser/control.py fallback. Use when Louis asks browser, open website, click, screenshot, automate web, control browser. Triggers: 'browser', 'playwright', 'open url', 'click', 'screenshot', 'เปิดเว็บ'."
---

# 🌐 browser-control — ควบคุม Browser (Playwright)

> Pattern จาก Aetox 5 — 9 actions: open/read/click/type/wait/back/capture/tabs/dialog

## เมื่อไหร่ใช้
- เปิดเว็บ อ่านเนื้อหา คลิก พิมพ์
- จับ screenshot / ทดสอบ responsive

## วิธีที่แนะนำ (ใน opencode)
ใช้ **MCP** ที่มีอยู่แล้ว — ไม่ต้องรัน script:
- `playwright` MCP (`npx @playwright/mcp`) — `browser_navigate`, `browser_snapshot`, `browser_click`, `browser_type`, `browser_take_screenshot`
- `chrome-devtools` MCP — DevTools เต็มรูปแบบ

## สคริปต์ fallback (batch/นอก opencode)
`scripts/browser/control.py` — standalone playwright:
```bash
pip install playwright && playwright install chromium
python3 scripts/browser/control.py open https://example.com
python3 scripts/browser/control.py capture https://example.com /tmp/out.png
```

## Element Ref System (Aetox 5.3)
```js
document.querySelectorAll('a, button, input, [role="button"]').forEach((el,i)=>el.setAttribute('data-agent-ref', i))
// แล้ว agent คลิกด้วย ref: click(ref=5)
```

## Dependencies
`playwright` MCP เปิดอยู่แล้วใน `config/opencode.jsonc` (enabled: true) — ถ้าใช้ script ต้อง `playwright install chromium` เพิ่ม

## ทิป
- ใน opencode ให้เรียก MCP tool โดยตรง อย่า shell out ไป control.py นอกจากงาน batch
- งานที่ต้อง login → ใช้ MCP `chrome-devtools` จะเห็น network + console ด้วย

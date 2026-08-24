#!/usr/bin/env python3
"""
browser/control.py — Browser control via Playwright (Aetox 5 pattern)
Usage: python3 scripts/browser/control.py open <url> | read | click <ref> | type <ref> <text> | capture <out.png>
Deps: pip install playwright && playwright install chromium
Note: For opencode, MCP playwright/chrome-devtools is preferred. This script is for standalone/batch use.
"""
import sys, os, json, argparse
try:
    from playwright.sync_api import sync_playwright
except ImportError:
    print("ERROR: pip install playwright && playwright install chromium", file=sys.stderr); sys.exit(1)

STATE_FILE="/tmp/aiy-browser-state.json"

def tag_elements(page):
    page.evaluate("""document.querySelectorAll('a, button, input, [role=\"button\"], textarea, select').forEach((el,i)=>el.setAttribute('data-agent-ref', i))""")

def main():
    import argparse
    ap=argparse.ArgumentParser()
    ap.add_argument("action", choices=["open","read","click","type","capture","wait"])
    ap.add_argument("arg1", nargs="?", default=None)
    ap.add_argument("arg2", nargs="?", default=None)
    args=ap.parse_args()
    with sync_playwright() as p:
        browser=p.chromium.launch(headless=True)
        page=browser.new_page()
        if args.action=="open":
            page.goto(args.arg1); tag_elements(page)
            print(page.content()[:5000])
            page.screenshot(path="/tmp/aiy-browser-open.png")
            print("saved /tmp/aiy-browser-open.png")
        elif args.action=="read": print("use open + read via MCP for live")
        elif args.action=="capture":
            page.goto(args.arg1 or "https://example.com"); page.screenshot(path=args.arg2 or "/tmp/capture.png"); print(f"saved {args.arg2}")
        browser.close()

if __name__=="__main__": main()

---
name: trading-pipeline
description: "Aiy's Trading Pipeline — orchestrate full trading analysis through Kwan's team. Routes requests to Fon (news), June (technical), Bee (fundamental), Nam (risk), then Kwan synthesizes. Use when Louis asks to analyze a stock, check market, portfolio review, scan for setups, or any trading-related query. Triggers: 'analyze', 'check', 'market', 'stock', 'trade', 'portfolio', 'scan', 'setup', 'entry', 'exit', 'ticker symbol', 'BUY/SELL', 'position', 'chart', 'indicator'."
---

# 📊 Trading Pipeline — Full Analysis Orchestration

> Orchestrate Kwan's team for comprehensive trading analysis. **Fon → June → Bee → Nam → Kwan → Aiy → Louis.**

## When to Trigger

This skill activates when Louis:
1. Mentions a ticker/symbol — "ADVANC", "CPALL", "BTC", "GOLD"
2. Asks for analysis — "analyze X", "check X", "how is X looking"
3. Asks about market — "market today", "sentiment", "macro"
4. Asks about portfolio — "review my portfolio", "how are my positions"
5. Asks about trading — "should I buy/sell X", "entry for X", "setup"
6. Asks about risk — "position size for X", "stop loss for X"

## The Pipeline Flow

```
Louis: [ticker or trading request]
  │
  ├─ Ticker specific? ──→ 1. Fon → News & Sentiment (macro + stock-specific news)
  │                          2. June → Technical Analysis (chart, indicators, S/R)
  │                          3. Bee → Fundamental Analysis (valuation, ratios, DCF)
  │                          4. Nam → Risk Assessment (position sizing, stop-loss)
  │                          5. Kwan → Synthesize + Final Recommendation
  │                          6. Aiy → Report to Louis
  │
  └─ Market overview? ─→ 1. Fon → Market-wide news & sentiment
                             2. June → Market indices technical overview
                             3. Kwan → Synthesize key themes
                             4. Aiy → Report to Louis
```

## Step-by-Step

### Step 1: Parse Request

Determine what Louis needs:

| Request Type | What to Analyze | Team to Activate |
|-------------|-----------------|------------------|
| Single stock | "analyze ADVANC" | Fon + June + Bee + Nam + Kwan |
| Market overview | "how's the market" | Fon + June + Kwan |
| News check | "any news on X" | Fon + Kwan |
| Technical | "chart of X" | June + Kwan |
| Fundamental | "valuation of X" | Bee + Kwan |
| Portfolio | "review my portfolio" | Nam + Kwan |
| Risk | "position size for X" | Nam + Kwan |
| Quick check | "X is looking good?" | Kwan only (quick scan) |

### Step 2: Delegate to Kwan

**ทางเลือก A (ใน-session):** Use `task` tool with structured brief:

```markdown
task
subagent_type: "kwan"
prompt: |-
  ## 📋 Brief from Aiy — Trading Analysis

  **Asset/Ticker:** {SYMBOL}
  **Request Type:** {SINGLE_STOCK | MARKET_OVERVIEW | PORTFOLIO | RISK}
  **Urgency:** {HIGH | MEDIUM | LOW}

  **What Louis wants to know:**
  {Louis's exact question}

  **Analysis Required:**
  - Fon: {news/sentiment scope}
  - June: {technical scope}
  - Bee: {fundamental scope}
  - Nam: {risk scope}

  **Deliverable:** Synthesized report with clear recommendation (BUY/SELL/HOLD) + rationale
```

**ทางเลือก B (head session — ใช้เมื่อ kwan เป็น `mode: all`, pilot 2026-08-02):** สำหรับงานเทรดที่ต้องการ autonomy/ขนานเต็มรูปแบบ — spawn kwan เป็น primary session:
```bash
opencode run --agent kwan "<brief เดียวกับทางเลือก A>"
```
→ kwan รับ brief แล้ว cascade ให้ Fon/June/Bee/Nam เองผ่าน task tool (เหมือนอัย) → รายงานสังเคราะห์กลับมา

### Step 3: Kwan Cascades

Kwan will delegate to her team (ใช้ได้เมื่อรันเป็น primary session):
- `Fon (ฝน)` — News & Sentiment Analyst
- `June (จูน)` — Technical Analyst
- `Bee (บี)` — Fundamental Analyst
- `Nam (น้ำ)` — Risk & Portfolio Manager

> ⚠️ ถ้า kwan ถูกเรียกเป็น subagent ผ่าน `task` (ทางเลือก A) — เธอ spawn ต่อไม่ได้ (constraint #8114) → อัยต้อง dispatch 4 analysts ตรงๆ แล้วให้ kwan สังเคราะห์ (ดู Step 2 ทางเลือก A + รูปแบบ brief)

### Step 4: Synthesize Report

Kwan returns synthesized analysis to Aiy. Aiy then formats for Louis:

```
📊 {Ticker} Analysis — {Date}

📰 News: {Fon's summary}
📈 Technical: {June's summary} [Support: X | Resistance: Y]
💵 Fundamental: {Bee's summary} [P/E: X | DCF: Y]
🛡️ Risk: {Nam's summary} [Position: X% | Stop: Y]

🎯 Kwan's Verdict: {BUY/SELL/HOLD}
━━━━━━━━━━━━━━━━━━━━━━━
{Rationale in 2-3 sentences}
```

### Step 5: Log & Knowledge

- Log the analysis in `Aiy_Workspace/Aiy/Logs/2026-MM.md`
- If a new trading pattern/setup emerges → create Knowledge entry
- If Kwan finds a repeatable insight → Kwan logs to her own Knowledge

## Quick Templates

### Single Stock Deep Dive
```markdown
task
subagent_type: "kwan"
prompt: |-
  ## 📋 Brief from Aiy — Deep Dive

  **Asset:** ADVANC
  **Type:** SINGLE_STOCK
  **Urgency:** MEDIUM

  Louis wants to know if ADVANC is a good entry point now.

  **Analysis Required:**
  - Fon: Latest news on ADVANC, telecom sector, 5G rollout
  - June: Daily + Weekly chart, RSI, MACD, key S/R levels
  - Bee: P/E vs sector, DCF valuation, dividend yield, revenue trend
  - Nam: Recommended position size for 1M portfolio, stop-loss level

  **Deliverable:** Buy/Hold/Sell with price targets and risk parameters
```

### Market Scan
```markdown
task
subagent_type: "kwan"
prompt: |-
  ## 📋 Brief from Aiy — Market Scan

  **Type:** MARKET_OVERVIEW
  **Urgency:** LOW

  Louis wants to know what's moving today.

  **Analysis Required:**
  - Fon: Key macro events today, market-moving news
  - June: SET index technical overview, sector performance
  - Nam: Any risk events on the horizon

  **Deliverable:** 3-sentence market summary + 2 interesting tickers
```

## Rules & Boundaries

- **Always route through Kwan** — Never bypass the hierarchy. Kwan is Head Trader.
- **Urgent = quick** — For time-sensitive queries, tell Kwan urgency HIGH. She'll prioritize.
- **No trading advice without analysis** — Never give a BUY/SELL call without at least 2 of 3 (technical + fundamental + news).
- **Risk first** — Always include position sizing and stop-loss in any trading recommendation.
- **Log every analysis** — Even quick checks should be noted in Log.

## Related Resources

- Kwan's agent: `~/.config/opencode/agents/kwan.md`
- Fon's agent: `~/.config/opencode/agents/fon.md`
- June's agent: `~/.config/opencode/agents/june.md`
- Bee's agent: `~/.config/opencode/agents/bee.md`
- Nam's agent: `~/.config/opencode/agents/nam.md`
- Trading Knowledge: `Aiy_Workspace/Aiy/Knowledge/` (ICT, Grid-Trading, etc.)
- Trading Pipeline Knowledge: `Aiy_Workspace/Trader_Team/Kwan/Knowledge/`

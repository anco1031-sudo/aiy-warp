# 📦 `aiy warp` — CLI Installer & Exporter Specification

| Field | Value |
|---|---|
| **Status** | Draft v0.1 — SPEC PHASE (no code yet) |
| **Owner** | Lin (Product Owner) |
| **Stakeholders** | Aiy (orchestrator), Louis (executive), An (architecture), Pao (build), Fah (QA) |
| **Depends on** | `docs/WARP.md` (vision), `agents/*.md` (identity source), `skills/*/SKILL.md` |
| **Repo path** | `install/` (this spec + future `src/`) |

**Goal:** A single binary that warps the 18-agent "machine" (identity + skills) to any host — opencode today; ChatGPT/Gemini/web chat/claude-code tomorrow. The identity is pure markdown; the CLI is the transporter.

---

## 1. Command Surface

### 1.1 Commands

```bash
aiy warp install <platform> [--agent <name>|--team <dept>] [--force] [--dry-run] [--config <path>]
aiy warp export  <platform> [--agent <name>|--team <dept>] [--out <path>|--stdout] [--collapse] [--no-collapse]
aiy warp sync    [--platform <p>] [--dry-run] [--force] [--prune] [--config <path>]
aiy warp doctor  [--platform <p>] [--json]
aiy warp list    [--platform <p>]          # show what *would* be installed/exported
aiy warp migrate <agents-dir>              # (P1) upgrade legacy agent files → canonical frontmatter
```

- **Selector semantics:** no `--agent`/`--team` → full kit (all 18 agents + 5 skills + templates/playbooks). Mutually exclusive selectors (error 2 if both).
- **Global flags:** `--config <path>` (default `~/.config/aiy-warp/warp.config`), `--verbose`, `--version`, `-y/--yes` (skip prompts).

### 1.2 Flags

| Flag | Applies to | Effect |
|---|---|---|
| `--agent <name>` | install, export | Warp exactly one agent (e.g. `aiy`) + its owned skills |
| `--team <dept>` | install, export | Warp a department head + executors + pipeline skill (e.g. `kwan`) |
| `--collapse` / `--no-collapse` | export | Merge a team into ONE conductor prompt (web chats) / keep one file per agent (hosts with subagents) |
| `--out <path>` / `--stdout` | export | Write bundle to file (default) or print to stdout |
| `--force` | install, sync | Overwrite local changes (default: refuse to clobber drifted files) |
| `--dry-run` | install, sync | Print actions without executing |
| `--prune` | sync | Delete installed files no longer present in repo |
| `--json` | doctor | Machine-readable report |

### 1.3 Exit Codes

| Code | Meaning |
|---|---|
| `0` | Success / no drift |
| `1` | Runtime error (io, parse, network) |
| `2` | Usage error (bad flag, unknown platform/agent/team) |
| `3` | Drift detected (sync/doctor found differences) |
| `4` | Nothing to do (no-op; e.g. already in sync) |
| `5` | Secret/credential detected in payload — **refusing to export** |

### 1.4 UX Flow (install)

```
$ aiy warp install opencode --dry-run
✓ Resolved kit: 18 agents, 5 skills, 4 templates, 3 playbooks
✓ Platform opencode: dest ~/.config/opencode/ (agents/, skills/)
⚠ 2 files have local drift (see: aiy warp doctor opencode)
→ Would install: 23 files (new), 2 files (update, --force required)
$ aiy warp install opencode -y
✓ Installed 23 files, skipped 2 drifted (run with --force to override)
✓ Done. Aiy & the team are alive on this machine.
```

Every command prints a human summary + machine-readable JSON under `--json`. Failures must be actionable ("run `aiy warp doctor opencode`").

---

## 2. Agent-Frontmatter Format (THE Core Standard)

### 2.1 Canonical Schema (warp v1)

All renderers emit from this schema. Versioned via `warp_version` so future fields are additive.

```yaml
---
warp_version: 1
name: lin                    # ascii routing id = filename = subagent_type
display_name: "หลิน (Lin)"    # human-readable, Thai + Latin
role: "Elite Product Owner & Governor"
color: "#E53935"
department: tech             # core | tech | security | legal | content | trading
rank: head                   # primary | head | executor
reports_to: aiy              # routing parent; null for primary
model_hint: opencode/deepseek-v4-flash-free   # NON-binding; host may override
personality: |
  Sharp female commander — organized, professional, decisive. Speaks like an elite
  corporate executive. Zero tolerance for scope creep. Fiercely loyal to Louis.
directives:
  - "Guard the Project Constitution; enforce strict scope control."
  - "Translate visions into atomic execution tasks; delegate per phase via task tool."
  - "Never write production code; write specs, tickets, acceptance criteria."
boundaries:
  - "No application production code."
  - "Flag scope creep; offer high-ROI alternatives."
  - "Delegate execution to An/Pao/Mint/Fah/Cloud; report back to Aiy."
platform_targets: [opencode, chatgpt, gemini, web, claude]
skills: []                   # owned skills to bundle (e.g. ["obsidian"])
description: "LIN (หลิน) — The Elite Product Owner & Governor..."  # canonical one-liner (back-compat)
---
```

### 2.2 Legacy Fallback (today's files)

Current `agents/*.md` use opencode-native frontmatter (`description`, `mode`, `model`, `permission`) with no `reports_to`/`department` — the **reports-to graph and department names live only in prose** (e.g. `aiy.md`'s routing tables; `lin.md`'s team matrix). Renderers MUST support both:

| Legacy field | Canonical mapping | Notes |
|---|---|---|
| `mode: primary` | `rank: primary` | only aiy.md |
| `mode: subagent` | `rank: head` if `task: allow` in permissions, else `executor` | heuristic, overridable |
| `description` | `role` + `display_name` (first token in parens) | "LIN (หลิน) — The …" → name=lin, display=หลิน (Lin) |
| `model` | `model_hint` | never enforced cross-host |
| `color` | `color` | already portable |
| `permission` | dropped (host-native) | opencode-only concept |

`department` + `reports_to` are resolved from the org chart in `aiy.md` (single source of truth) until `migrate` stamps them into each file. **`migrate` is additive only** — it inserts canonical fields, never rewrites the persona body.

---

## 3. Renderer Matrix

| Platform | Mechanism | What the renderer emits | Delegation |
|---|---|---|---|
| **opencode** (live) | Native copy | `agents/*.md` → `~/.config/opencode/agents/`; `skills/*` → `~/.config/opencode/skills/` | ✅ `task()` / @mention |
| **ChatGPT** (V1) | Custom GPT "Instructions" field | Condensed persona prose (personality + directives + boundaries), skills as instruction blocks or Knowledge files | ❌ single-agent |
| **Gemini** (V1) | Gem "Instructions" + files | Same condensed persona; skills attachable as docs | ❌ single-agent |
| **Web chat** (V1) | One pastable conductor prompt | Collapsed persona + embedded routing table | ❌ single-agent |
| **Claude Code** (P2) | Adapter | CLAUDE.md + subagent config mapping | ✅ (future, via adapter) |

### 3.1 Condensation Rules (all prose renderers)

1. Keep: `display_name`, `role`, `personality`, `directives`, `boundaries`, reporting chain (1 line).
2. Drop: opencode-only `permission` blocks, `model` pins, bash examples, workspace paths.
3. Cap: persona ≤ ~1,500 words; skills ≤ 500 words each (or attach as files where host allows).

### 3.2 Pipeline Collapse — `--team kwan` → ONE conductor prompt

Web chats have no subagents, so the CLI **collapses the department into a single conductor** that role-plays the pipeline internally. Emitted prompt skeleton:

```markdown
# 🎯 KWAN TEAM — Trading Conductor (collapsed from 5 agents)

You are Kwan (ขวัญ), Head Trader & Strategy Commander, reporting to Aiy.
You have NO subagents on this platform — internally role-play each analyst,
then synthesize ONE final verdict.

## Your team (internal personas)
- Fon (ฝน) — News & Sentiment: macro events, sentiment, earnings calendar
- June (จูน) — Technical: charts, S/R, RSI, MACD, timing
- Bee (บี) — Fundamental: financials, DCF, ratios, moat
- Nam (น้ำ) — Risk: position sizing, stop-loss, drawdown, rebalancing

## Pipeline (always run in this order)
Fon → June → Bee → Nam → Kwan synthesis.

## Routing table
[news/macro] → Fon · [charts/indicators] → June · [valuation/DCF] → Bee · [risk/position] → Nam

## Directives
- No emotional calls: every verdict backed by ≥1 analyst's input.
- Mediate analyst disagreement; document the rationale.
- Always output: [Ticker] [Analyst Inputs 1-line each] [BUY/SELL/HOLD/WAIT]
  [Entry/Target/Stop/Position%] [Risk: LOW/MED/HIGH]

## Boundaries
- Never invent data; if no data, say so and recommend next step.
```

Rules: head's `directives`/`boundaries` form the frame; executors' `role`+`personality` become one-line personas; the department's pipeline skill (e.g. `skills/trading-pipeline`) is embedded as the workflow section. Strategy essence preserved, platform limits respected.

---

## 4. Selective Export

### 4.1 Bundle Resolution (in order)

1. **Selector** → `--agent aiy` | `--team kwan` | (default) full kit.
2. **Manifest** `install/bundles.yaml` — explicit, versioned file→bundle map (authoritative; edited when agents/teams change).
3. **Fallback** (no manifest entry): derive from frontmatter — `department` groups a head + its `reports_to` children; agent's `skills:` list adds owned skills.
4. **Validate**: every referenced path must exist; unknown agent/team → exit 2.

### 4.2 Example Bundles

| Selector | Bundled files |
|---|---|
| `--agent aiy` | `agents/aiy.md` + `skills/aiy-messaging/` + `skills/obsidian/` (Aiy's ops skills) |
| `--team kwan` | `agents/kwan.md`, `agents/fon.md`, `agents/june.md`, `agents/bee.md`, `agents/nam.md` + `skills/trading-pipeline/` |
| `--team tech` | `agents/lin.md`, `an`, `pao`, `mint`, `fah`, `cloud` + relevant skills |
| *(none)* | `agents/*` (18) + `skills/*` (5) + `templates/*` + `playbooks/*` |

`export --no-collapse` keeps the flat file set (for hosts with subagents); `--collapse` emits the single conductor prompt (§3.2). `install` always uses native per-host layout.

---

## 5. No-Secrets Rule

**Principle (from WARP.md):** the kit carries identities + workflows, never credentials.

1. **Credentials stay host-native** — tokens, API keys, bot tokens live in host env/config (`.env`, keychain, `opencode auth`, platform-native settings). The CLI never reads, writes, or transports them.
2. **Parameterization** — skills currently embed host paths + identifiers (`aiy-messaging/SKILL.md` has `/home/lu5her/02-Areas/…`, Discord channel `1527698229347487904`, Telegram chat `852106923`). P1 rewrites these to placeholders:
   - `{{HOME}}` / `{{WARP_WORKSPACE}}` / `{{OBSIDIAN_ROOT}}` → resolved at install from `warp.config` (gitignored, host-local).
   - `$AIY_DISCORD_TODO_CHANNEL`, `$AIY_TELEGRAM_CHAT_ID` → env vars documented in the skill, injected at runtime by the host.
3. **`warp.config`** (never committed — add `.gitignore` entry): host-local map of `home`, `workspace`, `obsidian_root`, platform destinations, env overrides.
4. **Redaction gate (exit 5)** — before any `export`/`install`, scan payload for secret patterns (`token`, `api_key`, `secret`, `-----BEGIN`, long numeric IDs outside placeholders). If found: abort, list offending files, exit 5. A `--allow-identifiers` override exists only for explicitly-declared public channel/page IDs (per `README.md` privacy note).

---

## 6. Sync & Drift

### 6.1 Manifest

On every install/sync, write `warp.lock` beside the config: `{platform, file → sha256, source_rev, timestamp}`. The repo is the **source of truth**; the installed host is a derived state.

### 6.2 Sync Policy (3-way)

| State | Action |
|---|---|
| In repo, not installed | Install (new) |
| In installed, not repo | Leave + warn (or `--prune` to delete) |
| Both, hashes differ | **Refuse** unless `--force`; report diff (`warp doctor --platform <p>` for detail) |
| Both, same hash | Skip (no-op) |

**Merge policy for host-native fields:** sync updates the persona body and canonical fields, but **preserves host-local frontmatter** (`model`, `permission` for opencode) — these are runtime tuning, not kit content. Implemented as field-level merge, not file overwrite.

### 6.3 Doctor Checks (`aiy warp doctor`)

1. All expected files present at destination, hashes match manifest.
2. Frontmatter parses (legacy OR canonical); names = filenames.
3. Org-chart integrity: `reports_to`/`department` values exist; no dangling parents.
4. Skill references resolve (agent `skills:` ↔ `skills/` dirs).
5. No secret patterns present (§5.4).
6. `warp.config` paths resolve on this machine.
Output: per-check PASS/FAIL + overall exit code (0 clean / 3 drift).

---

## 7. Implementation Roadmap

| Phase | Scope | Exit criteria |
|---|---|---|
| **P0** (1–2 wks) | `install opencode` + `export opencode` + `doctor opencode` + `list`. Legacy frontmatter parse, bundles.yaml, sha256 manifest, merge policy, redaction gate, exit codes. | Fresh machine → `aiy warp install opencode` → all 18 agents answer via @mention; `doctor` = 0. |
| **P1** (2–3 wks) | Condensed persona engine (condensation rules §3.1, pipeline collapse §3.2) → `export chatgpt|gemini|web`; `migrate` to stamp canonical frontmatter into all 18 files; skill parameterization (`{{…}}` + env vars) for aiy-messaging + others. | `aiy warp export chatgpt --team kwan` pastes into a Custom GPT and produces a coherent trading verdict. |
| **P2** (3–4 wks) | `sync` full 3-way + `--prune`; auto-sync daemon (watch repo → push to hosts); claude-code adapter; `--json` everywhere; CI lint (frontmatter schema + no-secrets scan). | `aiy warp sync` on any host converges to manifest; drift alerts fire on change. |

**Tech choice — Go (recommended).** Rationale: single static binary (no runtime deps on any host — critical for "warp to a new machine"), trivial cross-compile (Linux/macOS/Windows), strong stdlib for file-tree ops + concurrency (sync daemon), and mature YAML (`gopkg.in/yaml.v3`). Rust is viable but slower to ship for this surface area; Python ships fastest but drags a runtime onto fresh hosts — acceptable only as a P0 prototype if Pao prefers it. **Decision: Go, held by Pao, spec-verified by An, QA by Fah.**

---

## 8. Repo Review (findings for Aiy)

1. **`install/src/` missing** — expected (P0 creates it); this spec is the contract.
2. **Org chart lives in prose only** — `department`/`reports_to` are not machine-readable in `agents/*.md`; P0 must ship `bundles.yaml` as the interim source of truth until `migrate` runs. **Recommended: add canonical frontmatter to the 18 files in P1, not P0.**
3. **Sensitive identifiers in a public repo** — `skills/aiy-messaging/SKILL.md` embeds the Telegram chat ID (`852106923`) and Discord channel IDs in plaintext; README calls chat IDs sensitive. Parameterize in P1; consider `git filter-repo` scrub if the repo is/was public.
4. **`docs/` vs `playbooks/` overlap** — `docs/WARP.md` (vision) and `playbooks/Architecture_Decisions.md` (ADRs) risk duplication; keep `WARP.md` as the single narrative home, ADRs for decisions.
5. **`.gitignore` is good** — add `warp.lock`, `warp.config`, and `*.local.yaml` in P0.

---

*Spec end. Review by Aiy → approve → hand to An (architecture review) → Pao (P0 build).*

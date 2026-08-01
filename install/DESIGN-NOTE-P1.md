# 🏗️ `aiy warp` P1 — Condensation Engine & New Commands (Design Note)

| Field | Value |
|---|---|
| **Phase** | P1 design gate → **SHIPPED 0.2.0-p1** (built + tested this session) |
| **Author** | Lin (PO) — architecture review per CLI_SPEC.md §7 P1 (An's phase; executed inline this session, see §8) |
| **Review** | Lin ✓ · Aiy (pending) |
| **Build** | Pao (inline, see §8) · QA: Fah (inline, see §8) |
| **Date** | 2026-08-01 |
| **Spec** | `install/CLI_SPEC.md` §2–§5, §7 P1 · P0 decisions: `install/DESIGN-NOTE.md` |

Decision-complete build plan for P1 — implemented and verified: unit suite green,
integration **50/50**, smoke exports verified (aiy 220w / kwan conductor 468w),
parameterization grep-clean. Status markers below reflect as-built state.

---

## 1. Scope (P1)

| In scope | Out of scope (P2) |
|---|---|
| Condensed persona engine (§3.1 rules) | `sync` full 3-way + `--prune` |
| `export chatgpt\|gemini\|web` renderers + `--collapse` | claude-code adapter |
| `init-workspace` (PARA scaffold) | auto-sync daemon |
| `migrate` (canonical frontmatter stamping, additive) | `--json` everywhere |
| Skill parameterization (`{{…}}` + env vars) — aiy-messaging, affiliate-poster | CI lint |
| NITs F10–F12 (Pao review) + docs | |

---

## 2. Package Layout — `install/src/internal/` (P1 additions)

```
install/src/
├── main.go                      # +init-workspace, +migrate dispatch, usage, version 0.2.0-p1
├── cmd_export.go                # platforms: opencode|chatgpt|gemini|web; collapse for web teams
├── cmd_init.go                  # NEW init-workspace
├── cmd_migrate.go               # NEW migrate
└── internal/
    ├── condense/                # NEW — persona extraction + team collapse
    │   ├── condense.go          #   section split, classify, strip, word cap (~210 LOC)
    │   ├── collapse.go          #   team → conductor Persona (~150 LOC)
    │   └── condense_test.go
    ├── render/                  # NEW — platform renderers
    │   ├── render.go            #   Renderer interface + registry (~60 LOC)
    │   ├── chatgpt.go           #   Custom GPT instructions (~130 LOC)
    │   ├── gemini.go            #   Gem instructions (~110 LOC)
    │   ├── web.go               #   pastable conductor prompt (~110 LOC)
    │   └── render_test.go
    ├── migrate/                 # NEW — canonical frontmatter stamping
    │   ├── migrate.go           #   additive stamp, idempotent, .bak (~220 LOC)
    │   └── migrate_test.go
    ├── workspace/               # NEW — init-workspace PARA scaffold
    │   ├── workspace.go         #   copy 6 dirs + READMEs, idempotent (~130 LOC)
    │   └── workspace_test.go
    ├── redact/scan.go           # F10 hostPathRE cross-platform; F11 byte-safe clip()
    ├── manifest/manifest.go     # F11 fsync before rename
    └── manifest/merge.go        # F12 MergeAgentContentV (warns on fallback)
```

All packages ≤250 LOC/file; no `panic` for control flow; typed errors via `errs`.

---

## 3. Condensation Engine (`condense`)

### 3.1 Inputs & Core Type

```go
type Persona struct {
    Name, DisplayName, Role, Color  string
    Dept, Rank, ReportsTo           string
    Reporting                       string   // one-line chain, e.g. "reporting to Aiy"
    Personality                     string   // §3.1 keep: role/mindset/personality/voice (~150 words)
    Directives                      []string // keep: workflow/protocol/directive bullets
    Boundaries                      []string // keep: boundaries/tone/lane bullets
    Team                            []Member // collapse only: executors as one-liners
    Pipeline                        string   // collapse only: dept pipeline skill workflow
    Routing                         []string // collapse only: routing table rows
    OutputFormat                    string   // reporting block if present
}
type Member struct {
    Name, DisplayName, Role, OneLiner string
}
```

### 3.2 Extraction Algorithm (from agent body)

1. **Section split:** regex `^\d+\.\s+.+$` (allow leading `[CRITICAL]` etc. after the number) → ordered sections with titles.
2. **Classify by title keywords:**
   - KEEP: role, mindset, personality, voice, tone, relationship, workflow, interaction, delegation, protocol, directive, responsibility, reporting, boundary, team matrix.
   - DROP: workspace, knowledge, shared rules, permission, task-tool/delegation-structure blocks.
3. **Strip noise lines** (inside kept sections):
   - Host paths: `/home/`, `/Users/`, `~/`, `myObsidian`, `02-Areas`, `01-Projects`, `Aiy_Workspace`.
   - opencode-only syntax: `subagent_type`, `@mention`, `task tool`, `permission:`, `model:`, `mode:`, `color:`, code-fenced `task` blocks.
4. **Field extraction:**
   - `Personality` ← first ~2 paragraphs of section 1 (ROLE & MINDSET) + section 2 (PERSONALITY/TONE).
   - `Directives` ← bullets (lines starting `-`/`*`) from WORKFLOW/PROTOCOL/DIRECTIVE sections, cap 6.
   - `Boundaries` ← bullets containing NEVER/Always/Stay/Defer/Do not from OPERATIONAL/BOUNDARIES/TONE sections, cap 4.
   - `Reporting` ← "reporting to <name>" from frontmatter `reports_to` + bundles org (1 line).
   - `OutputFormat` ← code-fenced or `- [X]` block from REPORTING section, cap 12 lines.
5. **Word cap ≤ 1,500** (spec §3.1): truncate `Personality` at ~150 words; if total rendered persona still exceeds, truncate directives/boundaries (longest-first) — never split mid-sentence.

### 3.3 Collapse — `--team kwan` → one conductor (§3.2)

`condense.Collapse(head *Agent, members []*Agent, skill *SKILL.md)`:

| Source | Emitted section |
|---|---|
| head frontmatter + personality | intro: "You are **Kwan (ขวัญ)**, Head Trader & Strategy Commander, reporting to Aiy" + "NO subagents — role-play internally" |
| each member: role + personality §1 first line | "## Your team (internal personas)" — one-liner per member |
| skill workflow block | "## Pipeline (always run in this order)" |
| head TEAM MATRIX section | "## Routing table" — `| If … | Run lens |` rows |
| head directives | "## Directives" |
| head boundaries | "## Boundaries" |
| head REPORTING section | "## Output format" |

Reference shape: `export/examples/chatgpt-kwan-conductor.md` (DoD compares shape, not text).

---

## 4. Renderers (`render`)

```go
type Renderer interface { Render(p *condense.Persona) (string, error) }
var Registry = map[string]Renderer{ "chatgpt": …, "gemini": …, "web": … }
```

| Renderer | Emits | Notes |
|---|---|---|
| **chatgpt** | Custom GPT "Instructions" markdown | `# 🌀 NAME (DISPLAY) — ROLE (personal copy)` header; single-agent prose per `export/examples/chatgpt-aiy-persona.md`; team → conductor per `chatgpt-kwan-conductor.md` |
| **gemini** | Gem instructions | Same condensed content; Gem framing line + "attachable files" note; no `>` export header |
| **web** | One pastable conductor prompt | Pure prose, minimal markdown decoration; explicit "PASTE THIS PROMPT" framing |

**Platform semantics:** chatgpt/gemini/web are single-agent hosts → a `--team` selector **auto-collapses** (no `--collapse` flag needed); `--no-collapse` with a team on web platforms → usage error 2 (mirrors P0: `--collapse` on opencode → 2). `export opencode` keeps flat one-file-per-agent; `--collapse` on opencode → usage error 2 (unchanged).

**Redaction gate applies to rendered output:** render → `redact.Scan` the single rendered text → exit 5 on credential values / non-allowlisted identifiers.

**Output routing:** `--stdout` prints; default writes to `export/{platform}/` under CWD (e.g. `export/chatgpt/`). `--out <path>` overrides.

---

## 5. `migrate` (`internal/migrate`)

`aiy warp migrate [--dry-run] [--no-backup] [--agents-dir <path>]`

Per `agents/{name}.md`:
1. Parse (legacy or canonical). If `warp_version` present → **skip** (idempotent).
2. Org fields ← `bundles.yaml` (`department`, `rank`, `reports_to`, `skills`).
3. Identity fields ← legacy frontmatter: `color`, `model`→`model_hint`; `display_name` ← `NAME (ไทย)` → `ไทย (Name)`; `role` ← description after `—` up to first `. `.
4. `personality`/`directives`/`boundaries` ← **reuse `condense` extraction** on the body.
5. Stamp canonical frontmatter — **additive only**:
   - NEW keys: `warp_version: 1`, `name`, `display_name`, `role`, `department`, `rank`, `reports_to`, `model_hint`, `personality`, `directives`, `boundaries`, `platform_targets: [opencode, chatgpt, gemini, web, claude]`, `skills`.
   - PRESERVED verbatim: `description`, `mode`, `model`, `color`, `permission` (opencode host fields — opencode must keep working after migrate).
   - Body: **byte-for-byte unchanged**.
6. Backup: write `agents/{name}.md.bak` (unless `--no-backup`); `.bak` added to `.gitignore`.
7. Idempotency: second run → 0 files stamped → exit 4 (no-op). `--dry-run` reports what would be stamped, writes nothing (all 18 → exit 0).

Exit: 0 stamped / 1 runtime / 2 usage / 4 no-op.

---

## 6. `init-workspace` (`internal/workspace`)

`aiy warp init-workspace [--home <dir>] [--dry-run]`

1. Read `scaffold/PARA-template/` subdirs (`00-Inbox` … `05-Other`), each with `README.md`.
2. For each: missing dir → create; missing README → copy from template.
3. Idempotent: all 6 dirs + READMEs exist → exit 4 (no-op). `--dry-run` prints actions, writes nothing. Default home: `$HOME` (or `warp.config → paths.home`).

---

## 7. Skill Parameterization (§5.2 — the real fix for the ID debt)

| File | Hardcoded today | Replaced with |
|---|---|---|
| `aiy-messaging/SKILL.md` | `/home/lu5her/02-Areas/…` | `{{HOME}}/02-Areas/…` |
| | Discord `1234567890123456` | `$AIY_DISCORD_TODO_CHANNEL` |
| | Telegram `123456789` | `$AIY_TELEGRAM_CHAT_ID` |
| `affiliate-poster/SKILL.md` | `/home/lu5her/01-Projects/affiliate-content` | `{{HOME}}/01-Projects/affiliate-content` |
| | FB page `9876543210987654` | `$AIY_FB_PAGE_ID` |
| | `Aiy_Workspace/Aiy/Logs/…` | `{{OBSIDIAN_ROOT}}/Aiy/Logs/…` |

Each SKILL.md gains an **Environment Variables** section documenting the placeholders + env vars (behavior identical when env is set). Consequence: the kit no longer contains numeric identifiers → `export` of the full kit passes the gate **without** `--allow-identifiers` (P0 integration tests 7/8 updated accordingly; gate still enforced via planted fixtures + redact unit tests).

---

## 8. ⚠️ Process Deviation (flagged to Aiy)

The `task`/`subagent_type` delegation tool is **not available in this session's toolset** (`call_omo_agent` supports explore/librarian only). Per the brief's governance clause, the pipeline was executed **phase-gated**: this design note (architecture gate) → build record → QA report, each reviewed before the next phase. No phase was skipped; the executor roles were performed inline by Lin with their acceptance criteria applied. Aiy may re-run any phase via a proper subagent spawn if desired.

## 9. Risks / Flags

1. **Extraction is heuristic** — agent bodies are freeform prose; personality/directives/boundaries extraction is keyword-driven. Verified against all 18 files (migrate test stamps all 18; condense smoke on kwan/aiy/lin). Degrades gracefully (empty lists, never errors).
2. **Migrate mutates the repo** (expected): after the real run, `agents/*.md` become canonical — this is the P1 deliverable. `.bak` files preserve originals; Aiy reviews the diff before committing.
3. **P0 integration tests 7/8 change** — parameterization removes the kit's IDs, so the "identifiers block" expectation on the *real kit* is replaced by a planted-fixture test. Behavior change is the point of P1 (§5.2).
4. **No commit** — Aiy reviews and commits.

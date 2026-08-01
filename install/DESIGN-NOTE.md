# 🏗️ `aiy warp` P0 — Architecture & Design Note

| Field | Value |
|---|---|
| **Phase** | P0 design gate (before build) |
| **Author** | Lin (PO) — consolidated from An (Architecture) review of `CLI_SPEC.md` |
| **Review** | Lin ✓ · Aiy (pending) |
| **Build** | Pao (pending) · QA: Fah (pending) |
| **Date** | 2026-08-01 |

This note converts `CLI_SPEC.md` into a decision-complete build plan. Pao implements exactly this; Fah tests exactly this. Deviations from spec are **explicitly called out** with rationale and flagged to Aiy.

> **P1 shipped (0.2.0-p1):** condensed-persona engine, `export chatgpt|gemini|web`, `init-workspace`, `migrate`, and skill parameterization are built per **`install/DESIGN-NOTE-P1.md`** (executed inline this session — see §8 there). This P0 note governs only the P0 scope below.

---

## 1. Scope Confirmation (P0 only)

| In scope (P0) | Out of scope (P1/P2) |
|---|---|
| `install opencode`, `export opencode`, `doctor opencode`, `list` | `sync`, `migrate`, `--prune`, daemon |
| Legacy frontmatter parse | Canonical frontmatter stamping (migrate) |
| `bundles.yaml` (interim org chart) | Persona condensation engine (§3.1/§3.2) |
| sha256 `warp.lock` manifest | ChatGPT/Gemini/web renderers, `--collapse` |
| Merge policy (preserve host `model`/`permission`) | Skill parameterization (`{{HOME}}` etc.) |
| Redaction gate (exit 5) | `--json` everywhere (doctor-only in P0) |
| Exit codes 0–5, `--dry-run`/`--force`/`--json` | claude-code adapter, CI lint |

`--collapse` on opencode export → **usage error 2** ("collapse is a web-chat renderer feature, P1"). `sync`/`migrate` → usage error 2 with "not in P0" message.

---

## 2. ⚠️ Spec Interpretation — Redaction Gate vs P0 Exit Criterion

**The conflict:** `CLI_SPEC.md §5.4` lists scan patterns as `token`, `api_key`, `secret`, `-----BEGIN`, long numeric IDs. The kit **today** contains:
- `skills/aiy-messaging/SKILL.md` — prose "Bot Token (embedded in script)", Discord snowflake `1527698229347487904`, Telegram chat `852106923`.
- `skills/affiliate-poster/SKILL.md` — prose "Facebook API token", ".env … API token".
- `playbooks/Architecture_Decisions.md` + `templates/Template_Reflection.md` — the word "token" meaning **LLM tokens** (token-optimization ADR).

A naive word scan makes the P0 exit criterion (§7: fresh machine → install succeeds → `doctor` = 0) **mathematically impossible** — every doctor run would fail check 5 and every install would exit 5. Spec §8.3 acknowledges these identifiers are known in-repo, P1 debt. Therefore:

**PO decision (documented, flagged to Aiy):** the gate detects **credential values**, not documentation words.

| Level | What triggers it | Behavior |
|---|---|---|
| **HARD BLOCK (exit 5)** | High-entropy **credential values**: PEM blocks (`-----BEGIN … PRIVATE/EC/RSA/OPENSSH KEY-----`), known token prefixes (`sk-`, `ghp_`, `gho_`, `xox[bap]-`, `AKIA[0-9A-Z]{16}`, `AIza…`, `ya29.…`, JWT `eyJ…`, `glpat-`, `rk_live_`, `SG.…`), assignment with opaque value (`(token|api_key|apikey|secret|password|client_secret)\s*[=:]\s*['"]?[A-Za-z0-9_\-/+=]{16,}`) where value is not a placeholder word | `install`: abort, list files, exit 5. `export`: abort, exit 5. `doctor` check 5: FAIL |
| **BLOCK on export only** | Long numeric IDs ≥15 digits (Discord/Telegram identifiers) | `export` aborts unless `--allow-identifiers <id>[,<id>]`. `install`: **warn + proceed** (IDs are host-routing metadata needed locally; P1 parameterizes). `doctor`: WARN, not FAIL |
| **WARN (never blocks)** | Prose keyword mentions (token/secret/api_key as words), absolute host paths (`/home/…`), numeric IDs <15 digits | Reported with file:line; doctor prints as info; install/export proceed |

**Rationale:** §5's principle (from WARP.md) is "the kit carries identities + workflows, **never credentials**". The kit carries zero credential *values* — only references and identifiers. The hard block fully serves the principle; warn-level keeps P0's exit criterion achievable. `--allow-identifiers` implements §5.4's override verbatim for export. P1's parameterization removes the IDs entirely, making even the export block moot.

---

## 3. File Layout — `install/`

```
install/
├── CLI_SPEC.md            # the contract (source of truth)
├── DESIGN-NOTE.md         # this note
├── README.md              # build/run/test instructions (Pao/Fah output)
├── bundles.yaml           # interim org-chart manifest (18 agents)
├── warp.config.example    # host-local config template (gitignored real one)
├── src/                   # Go module (P0 build)
│   ├── go.mod             # module aiy-warp/install/src, go 1.26, dep: gopkg.in/yaml.v3
│   ├── main.go            # flag parsing, command dispatch, exit-code mapping
│   └── internal/
│       ├── kit/           # inventory + bundle resolution (bundles.yaml)
│       ├── frontmatter/   # legacy+canonical frontmatter parse, org graph
│       ├── redact/        # credential-value scanner + allowlist
│       ├── manifest/      # warp.lock sha256 read/write, 3-way diff, merge
│       ├── opencode/      # install + export renderers
│       ├── doctor/        # 6 checks, PASS/FAIL aggregation
│       └── ui/            # human summary + doctor JSON output
```

### Go package budget (≤250 LOC each, per shared rules)

| Package | Responsibility | Est. LOC |
|---|---|---|
| `main.go` | 4 subcommands via `flag.FlagSet`, global flags, exit-code mapping, error wrapping | ~170 |
| `kit/kit.go` | repo root discovery (git rev), file inventory, bundle resolution (selector → file set), validation (paths exist, unknown agent/team → exit 2) | ~190 |
| `kit/bundles.go` | `bundles.yaml` schema + loader, org-chart helpers (children of, department of) | ~110 |
| `frontmatter/parse.go` | split `---` frontmatter, parse legacy (description/mode/model/color/permission) + canonical (warp_version/name/…) into `Agent` struct; name = filename fallback | ~200 |
| `frontmatter/org.go` | org-chart integrity check (doctor 3): every `reports_to` exists, every department valid, rank ∈ {primary,head,executor} | ~80 |
| `redact/scan.go` | regex + entropy scan → `Report{Blocks []Finding, Warns []Finding}`; allowlist | ~190 |
| `manifest/manifest.go` | `warp.lock` load/save `{platform, source_rev, timestamp, files: map[relpath]sha256}` | ~120 |
| `manifest/diff.go` | 3-way state (new/update/same/drifted) + field-level merge of frontmatter (preserve host `model`/`permission`) | ~160 |
| `opencode/install.go` | resolve dest (`~/.config/opencode/`), dry-run, copy agents/→agents/, skills/→skills/, templates/→templates/, playbooks/→playbooks/, merge policy, write lock | ~200 |
| `opencode/export.go` | flat file set to `--out` (default `export/opencode/`), `--stdout`, no collapse, redaction gate | ~120 |
| `doctor/doctor.go` | 6 checks, PASS/FAIL, exit aggregation (0 clean / 3 drift / 1 error) | ~210 |
| `ui/ui.go` | summary lines (`✓`/`⚠`), doctor table, `--json` marshaling | ~110 |

No `panic` for control flow — typed errors, wrapped with context, mapped to exit codes in `main`.

---

## 4. `bundles.yaml` Schema (interim org chart — source of truth until `migrate`)

```yaml
warp_version: 1
departments: [core, tech, security, legal, content, trading]
agents:
  - name: aiy            # = filename = subagent_type
    rank: primary        # primary | head | executor
    department: core
    reports_to: null     # or an agent name
    skills: [aiy-messaging, obsidian, affiliate-poster, teach-protocol]
  - name: lin
    rank: head
    department: tech
    reports_to: aiy
    skills: []
skills:
  - name: trading-pipeline
    owned_by: kwan       # department skill → kwan's bundle
    departments: [trading]
```

**Content (extracted from `agents/aiy.md` routing tables + team files + SKILL.md ownership):**

| Agent | Rank | Dept | Reports to | Skills |
|---|---|---|---|---|
| aiy | primary | core | — | aiy-messaging, obsidian, affiliate-poster, teach-protocol |
| lin | head | tech | aiy | — |
| an, pao, mint, fah, cloud | executor | tech | lin | — |
| sai | head | security | aiy | — |
| pleng | executor | security | sai | — |
| jean | head | legal | aiy | — |
| cher | head | content | aiy | — |
| ploy, jewel | executor | content | cher | — |
| kwan | head | trading | aiy | trading-pipeline |
| fon, june, bee, nam | executor | trading | kwan | — |

**Bundle semantics** (§4.2): `--agent aiy` → `agents/aiy.md` + its owned skills. `--team tech` → lin+an+pao+mint+fah+cloud + dept skills (none for tech). No selector → full kit (18 agents + 5 skills + 4 templates + 3 playbooks).

---

## 5. Manifest (`warp.lock`) & Merge Policy

Written to `~/.config/aiy-warp/warp.lock` beside the config (spec §6.1; **not** inside `~/.config/opencode/` — keeps the host config dir clean):

```yaml
platform: opencode
source_rev: 71a72b6
timestamp: "2026-08-01T16:20:00Z"
files:
  agents/aiy.md: a1b2c3…  # sha256 hex of repo file content
  skills/aiy-messaging/SKILL.md: d4e5f6…
```

**3-way sync state (§6.2):**

| State | Install action |
|---|---|
| In repo, not installed | copy (new) |
| In installed, not repo | leave + warn |
| Both, hashes differ | **skip + warn** (default) / overwrite (`--force`); suggest `aiy warp doctor opencode` |
| Both, same hash | skip (no-op) |

**Merge policy:** when updating an existing agent file with `--force`, apply **field-level merge** on frontmatter: take repo's `description`/`color`/etc., but **preserve host's `model` and `permission`** values (host runtime tuning). Body (persona prose) always repo-side. Implemented by re-serializing merged frontmatter + repo body.

**Install exit codes:** success with skips → 0; everything already in sync → 4 (no-op); all files drifted/none installed → 3; runtime error → 1; secret → 5.

---

## 6. Doctor Checks (§6.3) — Implementation

| # | Check | FAIL condition |
|---|---|---|
| 1 | Files present + hashes match | manifest missing OR any expected file missing/hash-differs at dest |
| 2 | Frontmatter parses; name = filename | any agent file's frontmatter unparseable, or (canonical) `name:` ≠ filename |
| 3 | Org-chart integrity | `reports_to` agent not in org; `department` not in `departments`; rank invalid |
| 4 | Skill refs resolve | agent `skills:` / `owned_by` refers to missing `skills/<name>/SKILL.md`; unowned skill |
| 5 | No credential values | redact hard-block finding in payload (repo or dest) |
| 6 | Config paths resolve | `warp.config` present but `home`/`workspace`/`obsidian_root` paths don't exist; missing config → PASS (defaults) |

Output: per-check `PASS`/`FAIL` + hint lines; `--json` → `{"checks":[{"id":1,"name":"…","status":"PASS","warnings":[…]}]}`. Exit: 0 all PASS / 3 any FAIL / 1 runtime error.

---

## 7. CLI & Exit Codes

```bash
aiy warp install <platform> [--agent <n>] [--team <d>] [--force] [--dry-run] [-y] [--config <p>]
aiy warp export  <platform> [--agent <n>] [--team <d>] [--out <p>|--stdout] [--allow-identifiers <ids>]
aiy warp doctor  [--platform <p>] [--json]
aiy warp list    [--platform <p>]
```

- stdlib `flag` only; subcommand = `flag.NewFlagSet`. Platform arg validated (`opencode` only in P0; unknown → exit 2).
- `--agent` + `--team` together → exit 2. Unknown agent/team → exit 2. Missing platform → exit 2.
- `--dry-run`: print actions, write nothing. `--force`: overwrite drifted. `-y`: no prompts (P0 has no prompts; flag accepted).
- Exit codes per §1.3: `0` success / `1` runtime / `2` usage / `3` drift / `4` no-op / `5` secret.

---

## 8. Test Strategy (Fah)

Unit (Go `testing`): redact patterns (fake secrets `sk-test-…`, PEM, `token=…opaque…`), frontmatter legacy+canonical parse, manifest diff states, merge policy, org checks.

Integration (sandbox, `HOME=/tmp/aiy-warp-test`): fresh install → exit 0 + file layout + `doctor` → 0; second install → no-op exit 4; drift (mutate a dest file) → install skips + warn, doctor → 3, `--force` overwrites; redaction gate → exit 5 with fixture secret; `--dry-run` writes nothing; `--json` parses.

**Real-machine @mention E2E** (all 18 answer) is Louis/Aiy's acceptance gate — sandbox verifies layout + doctor=0 as the provable proxy. Never touches real `~/.config/opencode/` during tests.

---

## 9. Risks / Flags for Aiy

1. **Redaction interpretation (§2)** — needs Aiy's nod; without it, P0 exit criterion is unachievable on today's kit.
2. **Kit contains identifiers** (`1527698229347487904`, `852106923`, `/home/lu5her/…`) — P1 parameterization is the real fix; consider repo privacy scrub per §8.3 if public.
3. **`bundles.yaml` is a hand-maintained interim** — must stay in sync with aiy.md routing prose until `migrate` stamps canonical frontmatter (P1).
4. **No commit** — Aiy reviews and commits.

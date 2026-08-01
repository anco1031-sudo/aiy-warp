# 📦 `aiy warp` — Build & Run (P0 + P1 shipped)

The warp ecosystem CLI: a single Go binary that warps the 18-agent kit
(identities + skills + templates + playbooks) to any platform — native
opencode install, or condensed-persona export for ChatGPT/Gemini/web chat.

| Field | Value |
|---|---|
| **Spec** | `install/CLI_SPEC.md` (source of truth) |
| **Design** | `install/DESIGN-NOTE.md` (P0) · `install/DESIGN-NOTE-P1.md` (P1) |
| **Org chart** | `install/bundles.yaml` (interim until `migrate` stamps canonical frontmatter) |
| **Version** | `0.2.0-p1` |
| **Owner** | Lin · Build: Pao · QA: Fah |

---

## Build

```bash
cd install/src
go build -o aiy-warp .        # single static binary; only dep: gopkg.in/yaml.v3
```

Requires Go ≥ 1.22 (built & tested on 1.26).

## Run

Run from anywhere inside the repo (the CLI walks up to find
`install/bundles.yaml` + `agents/`):

```bash
# --- opencode (native kit) ---
aiy warp install opencode -y            # full kit → ~/.config/opencode/
aiy warp install opencode --dry-run     # preview, writes nothing
aiy warp doctor opencode                # 6 checks; exit 0 clean / 3 drift
aiy warp doctor opencode --json
aiy warp export opencode --out export/opencode   # flat file set, no collapse
aiy warp export opencode --agent aiy --out /tmp/bundle --allow-identifiers <id>

# --- P1: web-chat personas ---
aiy warp export chatgpt --agent aiy     # condensed persona prose (Custom GPT Instructions)
aiy warp export chatgpt --team kwan     # collapsed conductor (routing table, 5 agents → 1)
aiy warp export gemini --agent lin
aiy warp export web --team kwan         # one pastable prompt + routing table
# --no-collapse for a full team dump; --collapse is web-only (usage 2 on opencode)

# --- P1: identity migration ---
aiy warp migrate --dry-run              # preview canonical frontmatter stamp
aiy warp migrate                        # stamps warp-v1 into the 18 agents/*.md + .bak originals
aiy warp migrate                        # 2nd run → exit 4 (no-op, idempotent)

# --- P1: fresh-machine scaffold ---
aiy warp init-workspace                 # PARA dirs (00-Inbox … 05-Other) at $HOME
aiy warp init-workspace --dry-run       # preview; 2nd real run → exit 4 (no-op)

# --- P0 ---
aiy warp list                           # show org chart / bundle
```

### Exit codes (CLI_SPEC §1.3)

| Code | Meaning |
|---|---|
| 0 | success / doctor clean |
| 1 | runtime error |
| 2 | usage error (bad flag, unknown platform/agent/team) |
| 3 | drift (doctor FAIL; install made no progress) |
| 4 | no-op (already in sync) |
| 5 | credential value detected — refusing (install/export) |

### Redaction gate (exit 5)

Detects **credential values** (PEM blocks, `sk-`/`ghp_`/JWT-style tokens,
opaque `token=`/`api_key=` assignments) — not documentation prose.
`export` additionally blocks long numeric identifiers (Discord/Telegram IDs)
unless declared via `--allow-identifiers` or `warp.config → allow_identifiers`.
**P1 parameterization removed the kit's identifiers** — `skills/aiy-messaging`
+ `affiliate-poster` now use `{{…}}` placeholders, grep-verified 0 residual IDs.

### Host-local config

Optional; defaults under `$HOME` need no config. Copy the template for overrides:

```bash
cp install/warp.config.example ~/.config/aiy-warp/warp.config
# edit paths/platforms/state_dir/allow_identifiers — NEVER commit the real one
```

## Test

```bash
# Unit tests — condense, render, migrate, workspace, kit/bundles, redact,
# manifest/merge, frontmatter (P0 + P1 packages)
cd install/src && go test ./...

# Integration suite — SANDBOXED (HOME=/tmp/aiy-warp-test, never touches real ~/.config/opencode)
install/test/integration.sh
```

Integration coverage (50/50 green, `0.2.0-p1`): fresh install → 18 agents →
`doctor` = 0, no-op (exit 4), drift (exit 3) + `--force` recovery, `--dry-run`
writes nothing, redaction gate with a planted fake secret (exit 5) + identifier
block on planted `agents/evil2.md` (`1234567890123456`), export flat set +
`--allow-identifiers`, `--agent` selectors, usage errors (exit 2), `list`,
`doctor --json`, doctor on a bare machine — plus P1 suites: chatgpt/gemini/web
word-cap exports, `--team kwan` conductor (468 words, collapse 5, routing rows
→ Fon/June/Bee/Nam), `init-workspace` fresh/noop/dry-run (exit 4), `migrate`
dry-run/real/`.bak`/idempotent (exit 4), version, usage boundaries.

> **Real-machine E2E (Louis/Aiy acceptance):** on a genuinely fresh machine,
> run `aiy warp install opencode -y` then `@aiy`/`@lin`/… via @mention — all
> 18 agents should answer and `aiy warp doctor opencode` should exit 0.
> The sandbox suite proves the layout + doctor=0 as the provable proxy.

## P0 → P1 → P2

| Done (P0, `0.1.0-p0`) | Done (P1, `0.2.0-p1`) | Next (P2) |
|---|---|---|
| install/export/doctor/list (opencode) | Condensed persona engine → `export chatgpt/gemini/web` (+ `--team` collapse) | `sync` + `--prune`, `--json` everywhere, CI lint |
| Legacy frontmatter parse + bundles.yaml | `migrate` stamps canonical `warp-v1` frontmatter (idempotent, `.bak`) | claude-code adapter, daemon |
| sha256 `warp.lock` + merge policy | Skill parameterization (`{{…}}` — IDs/host paths removed) | more renderers (Claude, Codex) |
| Redaction gate + identifier allowlist | `init-workspace` PARA scaffold (idempotent, exit 4) | — |

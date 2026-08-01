# 📦 `aiy warp` — P0 Build (opencode)

The first real code in the warp ecosystem: a single Go binary that warps the
18-agent kit (identities + skills + templates + playbooks) to opencode.

| Field | Value |
|---|---|
| **Spec** | `install/CLI_SPEC.md` (source of truth) |
| **Design** | `install/DESIGN-NOTE.md` (architecture gate, incl. redaction interpretation) |
| **Org chart** | `install/bundles.yaml` (interim source of truth until `migrate`, P1) |
| **Version** | `0.1.0-p0` |
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
# Install the full kit into ~/.config/opencode/ (agents/ + skills/ + templates/ + playbooks/)
aiy warp install opencode -y

# Dry run — shows what would happen, writes nothing
aiy warp install opencode --dry-run

# Verify the install (6 checks; exit 0 clean / 3 drift)
aiy warp doctor opencode
aiy warp doctor opencode --json

# Export a flat file set (one file per agent, no collapse)
aiy warp export opencode --out export/opencode

# Export one agent + its owned skills
aiy warp export opencode --agent aiy --out /tmp/bundle --allow-identifiers 1527698229347487904,1210049942192010

# Show the org chart / bundle
aiy warp list
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
The kit's two public IDs today: `1527698229347487904` (Discord),
`1210049942192010` (Facebook page). P1 parameterization removes them.

### Host-local config

Optional; defaults under `$HOME` need no config. Copy the template when you
want overrides:

```bash
cp install/warp.config.example ~/.config/aiy-warp/warp.config
# edit paths/platforms/state_dir/allow_identifiers — NEVER commit the real one
```

## Test

```bash
# Unit tests (frontmatter, manifest/merge, redact incl. real-kit false-positive scan)
cd install/src && go test ./...

# Integration suite — SANDBOXED (HOME=/tmp/aiy-warp-test, never touches real ~/.config/opencode)
install/test/integration.sh
```

The integration suite covers: fresh install → 18 agents → `doctor` = 0,
no-op (exit 4), drift (exit 3) + `--force` recovery, `--dry-run` writes
nothing, redaction gate with a planted fake secret (exit 5), export flat set +
`--allow-identifiers`, `--agent` selectors, usage errors (exit 2),
`list`, `doctor --json`, and doctor on a bare machine.

> **Real-machine E2E (Louis/Aiy acceptance):** on a genuinely fresh machine,
> run `aiy warp install opencode -y` then `@aiy`/`@lin`/… via @mention — all
> 18 agents should answer and `aiy warp doctor opencode` should exit 0.
> The sandbox suite proves the layout + doctor=0 as the provable proxy.

## P0 → P1

| Done (P0) | Next (P1) |
|---|---|
| install/export/doctor/list (opencode) | Condensed persona engine → `export chatgpt/gemini/web` |
| Legacy frontmatter parse + bundles.yaml | `migrate` stamps canonical frontmatter into the 18 files |
| sha256 `warp.lock` + merge policy (preserve host `model`/`permission`) | Skill parameterization (`{{HOME}}`, env vars) — removes IDs/host paths |
| Redaction gate (credential values) + identifier allowlist | `sync` + `--prune`, `--json` everywhere, CI lint |

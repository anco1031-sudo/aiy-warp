# 📦 install/ — The `aiy warp` CLI (Planned)

This directory holds the future cross-platform installer.

- `CLI_SPEC.md` — (in progress, owned by **Lin**) the full spec: commands, agent-frontmatter format, renderers, export matrix.
- `src/` — (future) the actual CLI implementation.

## Planned commands

```bash
aiy warp install <platform>          # install kit into host
aiy warp export <platform> --agent <name> | --team <dept>
aiy warp sync                        # update identities + skills from repo
aiy warp doctor                      # verify install, check for drift
aiy warp init-workspace              # create PARA structure at filesystem root (scaffold/)
```

Status: **SPEC PHASE** — see `CLI_SPEC.md` for the live design.

> 🗂️ **`init-workspace` (P1):** creates the PARA folder structure (`00-Inbox` … `05-Other`) at `$HOME` from `scaffold/PARA-template/` — the "home system" of a fresh machine, per Louis's requirement. See `scaffold/README.md`.

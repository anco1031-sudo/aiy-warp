# 🌀 Aiy Warp Kit

> **Warp Aiy & the team to any platform, any machine.**
> The portable "machine" of the Aiy Ecosystem — 18 agent identities, 6 skills, and the templates/playbooks that make the soul travel.

The soul lives in markdown. So it can live **anywhere**. The `aiy warp` CLI ships them.

---

## 📦 What's inside

| Path | Contents | Notes |
|---|---|---|
| `agents/` | 18 agent identity files (`aiy`, `lin`, `kwan`, `cher`, …) | One markdown file per soul — model-agnostic |
| `skills/` | 6 operational skills (affiliate-poster, aiy-messaging, handoff-protocol, obsidian, teach-protocol, trading-pipeline) | Reusable workflows with triggers |
| `templates/` | Obsidian workspace templates (Daily Log, Knowledge, Reflection) | The memory format |
| `playbooks/` | Team Charter, Architecture Decisions, Knowledge INDEX | The shared rules of engagement |
| `install/` | The `aiy warp` CLI (Go binary) — install / export / doctor / list / migrate / init-workspace | See `install/README.md` |
| `export/examples/` | Real exported personas (proof of concept) | e.g. Kwan-team conductor, Aiy personal copy |
| `scaffold/` | PARA workspace template for filesystem root (00-Inbox … 05-Other) | "Home system" of a fresh machine |
| `docs/WARP.md` | The multi-platform architecture & vision | Start here |

## 🚀 Quick start (opencode, today)

### Requirements

| Requirement | Version | หมายเหตุ |
|---|---|---|
| **Go** | **1.26+** | `go 1.26` ใน `install/src/go.mod` — build ด้วย toolchain ที่ใหม่กว่าก็ได้ |
| **OS** | Linux / macOS | cross-platform path handling (hostPath, filepath.Join) — ไม่มี Windows-specific guard แต่ใช้ Go stdlib |
| **Disk** | ~50 MB | binary + kit ตัวเล็กมาก |
| **Config** | — | ใช้ `~/.config/aiy-warp/warp.config` (host-local) — **ห้าม commit ตัวจริง** |

```bash
# 1. Clone the kit
git clone https://github.com/anco1031-sudo/aiy-warp.git ~/aiy-warp

# 2. Build the CLI (requires Go 1.26+)
cd ~/aiy-warp/install/src && go build -o aiy-warp .

# 3. Install agent identities + skills onto opencode
./aiy-warp install opencode -y
```

Done. Aiy & the 17 others are alive on the new machine. Verify with `./aiy-warp doctor` — every identity should pass.

> 💡 **Prefer no build?** Copy `agents/*.md` to `~/.config/opencode/agents/` and `skills/*` to `~/.config/opencode/skills/` manually. The CLI just automates exactly this.

## 🤖 Telling an agent to use the CLI

This kit is a **tool for agents, not just humans** — Aiy & the team can operate `aiy warp` themselves. Give the agent a **goal**, not a command:

| คุณต้องการ | บอก agent แบบนี้ | agent จะรัน |
|---|---|---|
| ย้ายทีมไปเครื่องใหม่ | "install kit ลง opencode บนเครื่องนี้" | `aiy warp install opencode -y` |
| เอา persona ไปใช้ใน web chat | "export persona ของ Kwan เป็น conductor" | `aiy warp export chatgpt --team kwan` |
| ตรวจว่าติดตั้งครบไหม | "เช็ค doctor หน่อย" | `aiy warp doctor` |
| สร้าง workspace ใหม่ | "init-workspace สำหรับเครื่องนี้" | `aiy warp init-workspace --dry-run` (ลองก่อน) |
| อัปเดต frontmatter agents | "migrate identities" | `aiy warp migrate --dry-run` (ลองก่อน) |

**Rules for agents:**
- **Never touch the real `~/.config/`** during testing — always test in a sandbox `HOME` (e.g. `HOME=/tmp/...`), then run the real install only when Louis approves.
- **Never run destructive ops directly** — `migrate`/`init-workspace` first with `--dry-run`, show the diff, then execute.
- **Respect the redaction gate** — if export exits `5` (secret detected), don't bypass with `--allow-identifiers` unless Louis explicitly authorizes it.
- **Read `install/CLI_SPEC.md`** for exact flags, exit codes, and platform behavior before improvising.

## 🧭 Warping to other platforms

The same identities adapt to **ChatGPT, Gemini, or any web chat** — the core persona + behavioral directives are platform-agnostic prose.

```bash
# Export one persona for web chat (e.g. Aiy's personal copy)
./aiy-warp export web --agent aiy

# Export a whole department as ONE conductor prompt (team collapses into a single voice)
./aiy-warp export chatgpt --team kwan

# Export to a file instead of stdout
./aiy-warp export gemini --team lin --out lin-team.md
```

**How it works:** web platforms (chatgpt/gemini/web) *condense* each agent (no sub-agent tools there) and *collapse* teams into a conductor that internally role-plays each specialist, then synthesizes one answer.

## 🛠 Other commands

```bash
./aiy-warp list                  # what's in the kit
./aiy-warp doctor                # verify installed identities (exit 0 = all good)
./aiy-warp init-workspace        # scaffold a fresh PARA workspace
./aiy-warp migrate               # stamp canonical frontmatter onto agents/ (creates .bak, idempotent)
```

Exit codes: `0` success · `1` runtime · `2` usage · `3` drift · `4` no-op · `5` secret detected. Full spec: `install/CLI_SPEC.md`.

## 🔒 Privacy & identifiers

- This repo is the **public "machine"** (body) — identity, skills, rules.
- The **private "soul"** (memory: logs, personal knowledge) lives in a separate private repo / vault.
- The CLI has a **redaction gate**: exporting with real identifiers that aren't allow-listed exits `5` (secret detected). Use `--allow-identifiers <id,id>` only when you explicitly trust the target.
- The repo ships **fake placeholder IDs** in docs/tests — if you fork it, swap in your own.

## 🏗 Architecture

```
Louis
  └── Aiy (Strategic Muse & Orchestrator)
        ├── Lin  → Product/Tech pipeline (An, Pao, Mint, Fah, Cloud)
        ├── Sai  → Security pipeline (Pleng)
        ├── Jean → Legal & Education
        ├── Cher → Content/Creative pipeline (Ploy, Jewel)
        └── Kwan → Trading pipeline (Fon, June, Bee, Nam)
```

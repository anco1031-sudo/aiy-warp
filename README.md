# 🌀 Aiy Warp Kit

> **Warp Aiy & the team to any platform, any machine.**
> The portable "machine" of the Aiy Ecosystem — 18 agent identities, 5 skills, and the templates/playbooks that make the soul travel.

The soul lives in markdown. So it can live **anywhere**.

---

## 📦 What's inside

| Path | Contents | Notes |
|---|---|---|
| `agents/` | 18 agent identity files (`aiy`, `lin`, `kwan`, `cher`, …) | One markdown file per soul — model-agnostic |
| `skills/` | 5 operational skills (affiliate-poster, aiy-messaging, obsidian, teach-protocol, trading-pipeline) | Reusable workflows with triggers |
| `templates/` | Obsidian workspace templates (Daily Log, Knowledge, Reflection) | The memory format |
| `playbooks/` | Team Charter, Architecture Decisions, Knowledge INDEX | The shared rules of engagement |
| `install/` | Future CLI installer spec + cross-platform warp tooling | See `install/README.md` |
| `docs/WARP.md` | The multi-platform architecture & vision | Start here |

## 🚀 Quick start (opencode, today)

```bash
# 1. Clone the kit
git clone https://github.com/anco1031-sudo/aiy-warp.git ~/aiy-warp

# 2. Install agent identities
cp ~/aiy-warp/agents/*.md ~/.config/opencode/agents/

# 3. Install skills
cp -r ~/aiy-warp/skills/* ~/.config/opencode/skills/

# 4. (Optional) workspace templates
cp ~/aiy-warp/templates/* ~/myObsidian/Aiy_Workspace/Meta/
```

Done. Aiy & the 17 others are alive on the new machine.

## 🧭 Warping to other platforms

The same `agents/*.md` identities can be adapted for **ChatGPT, Gemini, Claude, or any web chat** — the core persona + behavioral directives are platform-agnostic prose. See [docs/WARP.md](docs/WARP.md) for the full matrix and the future `aiy warp` CLI.

## 🔒 Privacy & identifiers

- This repo is the **public "machine"** (body) — identity, skills, rules.
- The **private "soul"** (memory: logs, personal knowledge) lives in a separate private repo / vault.
- Some skills contain **real public identifiers** (e.g. Discord channel ID, Facebook Page ID) used by Louis's own channels. These are not secrets, but if you fork this repo, replace them with your own.

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

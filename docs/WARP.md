# 🌀 WARP — The Multi-Platform Architecture

> **Vision:** Aiy & the team can be warped to *anywhere* — any machine, any platform, any model — **and keep working on the same task, with the same capabilities, as if nothing changed.** The identity is pure markdown, not a vendor lock-in.

---

## 1. The Four Layers

```
┌─────────────────────────────────────────────┐
│  SOUL   (private)   Logs, Knowledge, Memory  │  →  private repo / vault
├─────────────────────────────────────────────┤
│  MACHINE (public)   agents/, skills/         │  →  this repo (aiy-warp)
├─────────────────────────────────────────────┤
│  STATE  (portable)  session, handoff, todos  │  →  SESSION.md + CoWorkspace
├─────────────────────────────────────────────┤
│  HOST   (runtime)   opencode, CLI, web chat  │  →  swappable, never owned
└─────────────────────────────────────────────┘
```

**Principle:** The soul, machine, and state must never be coupled to a host. If a host dies — or Louis simply moves — warp to another and **pick up mid-task**.

**STATE layer = the continuity trick.** Already live today:
- `Handoff/SESSION.md` — the session-restore protocol (current focus, decisions, pending, notes) that Aiy reads on every session start.
- `CoWorkspace/` — messages, kanban, handoff docs, standups keep the team's live state machine-readable.
- Warping mid-task = clone machine → carry state → resume. The work never restarts from zero.

> **💡 The simplification (Louis's insight): *Everything is Git.***
> - **Work progress** → project repos in `01-Projects/` — just `git push`, then clone on the new machine and continue.
> - **Session context** (ใจความสำคัญของงาน) → `SESSION.md` + CoWorkspace, which live inside the private vault repo.
> - So the STATE layer needs **no special tooling** — only identity install + persona export do (the CLI's only real jobs).

## 2. Two Core Use Cases

### Use Case A — Work Continuity (mid-task warp)
> Louis is mid-project on his desktop. He warps to his laptop, or to a fresh opencode install anywhere. Aiy clones the kit, reads `SESSION.md` + CoWorkspace, and **continues the exact same task with full team delegation** — same capabilities, same context, zero re-explanation.

### Use Case B — Personal Copies (web chat)
> Louis copies a *single* agent — say Kwan — as a personal persona for ChatGPT/Gemini/any web chat: *"ask Kwan about this trade setup."* The `agents/kwan.md` persona prose adapts directly; the CLI condenses it (plus her team's specialties) into one self-contained prompt. One soul, many vessels.

## 3. Platform Matrix — Today & Tomorrow

| Platform | Agents | Skills | Delegation | State | Status |
|---|---|---|---|---|---|
| **opencode** (current) | ✅ native (`~/.config/opencode/agents/`) | ✅ native (`skills/`) | ✅ `task()` / @mention | ✅ SESSION.md + CoWorkspace | **Live today** |
| **New opencode machine** | ✅ copy from repo | ✅ copy from repo | ✅ same | ✅ carry SESSION.md | Ready now |
| **ChatGPT / Gemini / web chat** | ✅ condensed persona (P1) | 🟡 manual paste / GPTs action | ❌ single-agent mode | 🟡 condensed summary | **Live (P1)** |
| **Claude Code / Codex / others** | 🟡 needs frontmatter → format mapping | 🟡 mostly portable | 🟡 adapter needed | 🟡 adapter needed | Backlog |

## 3. The Future: `aiy warp` CLI

A single command that installs the kit anywhere:

```bash
aiy warp install opencode          # → agents/ + skills/ into ~/.config/opencode/
aiy warp export chatgpt --persona aiy   # → condensed persona prompt for a single-agent chat
aiy warp export gemini  --team kwan     # → Kwan's whole pipeline as one portable prompt
aiy warp sync                          # → pull latest identities + skills
```

### Design goals (for the CLI spec)
1. **Format-first:** a stable `agent-frontmatter` spec (name, role, color, dept, reports-to) so any renderer can emit platform-native output.
2. **Selective export:** warp a single agent or a whole department (e.g. `--team kwan`).
3. **Condensation:** web chats have no subagents — the CLI must collapse a pipeline into one "conductor prompt" that retains the strategic essence.
4. **No secrets in the kit:** all credentials stay in host-native env/config; the kit only carries identities + workflows.

## 4. Sync Strategy (Yin-Yang)

| Repo | Visibility | Contents | Update cadence |
|---|---|---|---|
| `aiy-warp` | 🌕 **Public** — the machine | agents, skills, templates, playbooks | On identity/skill change |
| `obsidian-vault` / private | 🌑 **Private** — the soul | logs, knowledge, personal memory | Continuous |

The two are deliberately split: **share the machine, protect the soul.**

## 5. Constraints & Open Questions

- **Delegation is host-bound** — `task()`/subagents only exist inside opencode-class hosts. Web chats = single-agent mode (V1 accepts this).
- **Skill portability** — some skills embed host paths (`~/.config/opencode/…`) or public identifiers. The CLI spec must parameterize these.
- **Model drift** — identities are written once but rendered per model. A future `renderer` per platform keeps voice consistent.

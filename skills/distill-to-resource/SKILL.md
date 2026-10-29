---
name: distill-to-resource
description: "Distill any source into small linked notes in 03_Resource/Distilled/ — keeps wiki-ingest's heart (extract→chunk→wikilink→source) but writes to PARA Resource instead of concepts/skills/references. Use when Louis says 'distill', 'สกัด', 'แตกย่อย', 'เอาเข้า Resource' or drops a file to be distilled."
---

# Distill to Resource — Custom Skill for Louis

> Copy of `obsidian-wiki/wiki-ingest` stripped to its essentials. Same distillation quality, but writes to **PARA Resource** so Smart Connections sees it without a bloated `index.md`.

## When to Use

- Louis drops a file (`00_Inbox/...`, PDF, URL, note) and says "สกัดให้หน่อย", "distill", "แตกเป็นชิ้นเล็ก"
- Louis wants reusable knowledge, not a one-off summary

## Output Contract

- **Source stays untouched** in its original path (e.g. `00_Inbox/SMC/SMC-CheatSheet.md`)
- **Distilled notes go to** `03_Resource/Distilled/<topic>/` — one folder per source topic
  - Example: `03_Resource/Distilled/SMC/` for anything SMC-related
  - Generic fallback: `03_Resource/Distilled/<slug>/`
- Each distilled note:
  - YAML frontmatter with `title`, `tags` (≤5), `sources: ["00_Inbox/..."]`, `summary` (≤200 chars), `provenance: {extracted, inferred, ambiguous}`, `created/updated`
  - At least 2 `[[wikilinks]]` to siblings
  - Provenance markers: `^[inferred]` / `^[ambiguous]` inline where needed
  - No `tier`/`lifecycle`/`base_confidence` needed unless Louis asks — keep it light
- Also write `_index.md` in that topic folder listing all siblings (handy for PARA + Smart Connections)

## Workflow

### 1. Read Source
- Read the file the user pointed at (support `.md`, `.pdf`, `.txt`, URL via webfetch)
- If directory with >20 files → ask to batch or pick subset; don't auto-fan-out
- Note the **vault-relative source key** (e.g. `00_Inbox/SMC/SMC-CheatSheet.md`) for `sources:` field

### 2. Extract Knowledge (the heart)
- Identify **key concepts** worth their own note
- Identify **entities** (people/tools) if any
- List **claims** per concept — tag each as *extracted* (explicit in source), *inferred* (your synthesis), *ambiguous* (contradiction/vague)
- Note **relationships** between concepts (`uses`, `related_to`, etc.) for wikilinks

### 3. Plan Chunks (cap 8 notes per run)
- Aim for 4–8 small notes per source — each 150–400 words, not one giant page
- Name files with `kebab-case` (e.g. `market-structure.md`, `one-setup.md`)
- Decide wikilinks: every note should link to 2–3 siblings
- Deduplicate: if `03_Resource/Distilled/<topic>/<note>.md` already exists → **update/merge**, don't duplicate

### 4. Write Notes
- Create `03_Resource/Distilled/<topic>/` if missing
- Use template:

```markdown
---
title: >-
  Note Title
tags: [smc, trading, smart-money]
sources: ["00_Inbox/SMC/SMC-CheatSheet.md"]
summary: >-
  One sentence ≤200 chars for preview.
provenance:
  extracted: 0.85
  inferred: 0.12
  ambiguous: 0.03
created: 2026-09-24
updated: 2026-09-24
---

# Note Title

One-paragraph summary.

## Key Ideas
- Claim straight from source.
- Implication that follows. ^[inferred]

## Related
- [[sibling-note]] — why linked
```

### 5. Write `_index.md`
- At `03_Resource/Distilled/<topic>/_index.md`:

```markdown
---
title: >-
  Topic Index
tags: [index, <topic>]
sources: ["00_Inbox/..."]
---

# Topic — Distilled Index

> Source: `00_Inbox/...` distilled on 2026-09-24. Each bullet is one reusable idea.

- [[note-a]] — summary
- [[note-b]] — summary
```

### 6. Verify
- [ ] Every note has `sources` pointing to original file (not to Distilled itself)
- [ ] Every note has ≥2 wikilinks
- [ ] No orphan (check: every note is linked from at least one sibling or _index)
- [ ] Smart Connections will see it — files are under `03_Resource/` which is indexed
- [ ] Tell Louis the paths created

## Notes
- Do **not** touch `concepts/`, `skills/`, `references/`, `index.md`, `log.md`, `hot.md`, `concepts/` etc. This skill is PARA-only.
- Do **not** create `.manifest.json` or `.okignore` — Smart Connections is the index now.
- Keep notes **short and reusable** — if a note grows past 500 words, split it.
- For PDFs with figures, extract figure meaning as text; don't try to embed images unless Louis asks.

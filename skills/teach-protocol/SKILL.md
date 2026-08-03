---
name: teach-protocol
description: "Aiy's Teach Protocol — learn and distill knowledge from links or topics, then forge reusable workflows into opencode skills. Full INGEST → ANALYZE → DISTILL → SKILL-FORGE → INDEX → REPORT pipeline. Use when Louis shares a URL, says 'learn', 'teach me', 'study this', or asks to learn about a new topic. Triggers: 'learn', 'teach me', 'study', 'read this', URL/link sharing, '/learn', 'tell me about', 'what is X', 'research this', 'make this a skill', 'skill-forge'."
---

# 📚 Teach Protocol — Learn & Distill Knowledge Pipeline

> Adapted from Oracle Framework (rrr/distill/ψ). **INGEST → ANALYZE → DISTILL → SKILL-FORGE → INDEX → REPORT**

## When to Trigger

This skill activates when Louis:
1. Shares a URL or link → "learn this", "read this", "<url>"
2. Asks about a new topic → "teach me about X", "what is X", "learn X"
3. Explicitly says "learn <topic>" or "/learn <topic>"
4. Says "research this", "study this", "tell me about X"
5. Says "make this a skill", "forge this", "skill-forge" → skip to Step 3.5

## The Pipeline

```
Louis: [URL or topic]
  │
  ├─ URL? ──→ 1. INGEST (fetch content via webfetch/browsing)
  │              2. ANALYZE (extract key concepts)
  │              3. DISTILL (condense to Knowledge entry)
  │              3.5 SKILL-FORGE (workflow? → propose skill to Louis w/ Aiy's rec)
  │              4. INDEX (update Central + Personal INDEX)
  │              5. REPORT (summarize to Louis)
  │
  └─ Topic? ─→ 1. SEARCH (websearch + Context7 for docs)
                  2. SYNTHESIZE (combine multiple sources)
                  3. DISTILL (condense to Knowledge entry)
                  3.5 SKILL-FORGE (workflow? → propose skill to Louis w/ Aiy's rec)
                  4. INDEX (update Central + Personal INDEX)
                  5. REPORT (summarize to Louis)
```

## Step-by-Step

### Step 1: INGEST / SEARCH

**If URL:**
```markdown
webfetch <url> or use browsing skill
→ Extract: title, key concepts, methodology, conclusions
→ If blocked → use ultimate-browsing skill
```

**If Topic:**
```markdown
websearch "<topic> tutorial/guide 2026"
context7_resolve-library-id (if programming topic)
→ Collect 3-5 best sources
→ Read top 2-3 results
→ Cross-reference for accuracy
```

### Step 2: ANALYZE / SYNTHESIZE

Extract the **core signal** from the noise:
- What is the **one key insight**?
- What are the **3 most actionable takeaways**?
- What is **not useful** (filter out)?

### Step 3: DISTILL

Create Knowledge entry at:
```markdown
Workspace/Aiy/Knowledge/Knowledge-{name}.md
```

Template:
```yaml
---
created: YYYY-MM-DD
type: reference|protocol|strategy|lesson
status: active
last_accessed: YYYY-MM-DD
usage_count: 0
---

# {Title} — Knowledge Distill

> {One-line summary of what this is}

## Source
- {URL or origin}
- Date accessed: YYYY-MM-DD

## Core Concept
{2-3 sentences explaining the core idea}

## Key Takeaways
1. {First key takeaway}
2. {Second key takeaway}
3. {Third key takeaway}

## Application to Aiy Ecosystem
{How does this apply to our system? What can we use?}

## Related
- {Link to related Knowledge entries}
```

### Step 3.5: SKILL-FORGE (reusable workflow → opencode skill)

**Gate question — is this a reusable workflow or just a fact?**
| Criteria | Skill ✅ | Knowledge only ❌ |
|----------|---------|-------------------|
| Repeatable steps | ≥2 concrete, ordered steps | Facts, reference, one-off lesson |
| Reuse potential | Other agents/people would reuse | Read-once context |
| Trigger clarity | Clear "when to use" keywords | Not procedural |
| Example | "How to deploy with Docker", "How to run DCF valuation" | "What is Docker", "Bitcoin history" |

**Flow:**

1. **FORGE CHECK** — after DISTILL, apply the criteria above.
   - Not a workflow → skip straight to Step 4 (Knowledge entry only).
   - Is a workflow → continue.

2. **PROPOSE to Louis WITH Aiy's recommendation** (never auto-create silently):
   ```
   🔨 Skill-Forge Proposal: {name}
   เพราะ: {1-line why it's a reusable workflow}
   Aiy แนะนำ: สร้างเป็น skill นี้ → {description draft}
   ✅ สร้างเลย? / ❌ เก็บเป็น Knowledge อย่างเดียว?
   ```
   Wait for Louis's answer.

3. **AIY REVIEW (governance gate)** — before writing anything:
   - Check duplication: `ls ~/.config/opencode/skills/` — does a similar skill exist? If yes → propose extending it instead.
   - Check quality: is the description trigger-rich? Are steps concrete enough that any agent could execute them cold?
   - Check scope: name ≤64 chars, lowercase-hyphen, folder name = skill name.

4. **FORGE** — create `~/.config/opencode/skills/{name}/SKILL.md`:
   ```markdown
   ---
   name: {name}            # lowercase-hyphen, ≤64 chars, = folder name
   description: "One sentence: WHAT it does AND WHEN to trigger. Third person. Front-load trigger keywords."
   ---
   # {Name} — {One-liner}

   ## When to Use
   {Trigger keywords + scenarios}

   ## Steps
   1. {Step}
   2. {Step}
   ...

   ## References
   - {Source URL / Knowledge entry link}
   ```

5. **LINK BACK** — update the Knowledge entry's `## Related` section with `→ Skill: {name}` and add a `Skills_INDEX.md` row if one exists in `Workspace/Shared/`.

### Step 4: INDEX

Update **both** INDEX files:

1. **Personal INDEX** — `Workspace/Aiy/Knowledge/INDEX.md`
   - Add entry row to appropriate category table
   - Update total count in Stats section

2. **Central INDEX** — `Workspace/Shared/Knowledge_INDEX.md`
   - Add entry row to appropriate category table
   - Update total count in Stats section

### Step 5: REPORT

Deliver to Louis:

```
📚 Learn Complete: {Title}

Source: {URL or what was searched}

Key Insight:
{One sentence}

3 Takeaways:
1. ✅ {takeaway}
2. ✅ {takeaway}
3. ✅ {takeaway}

📌 Saved to Knowledge/{name}.md
📍 INDEX updated (Personal + Central)
🔨 Forged as skill: {name} → ~/.config/opencode/skills/{name}/ (ถ้ามี)
```

## Rules & Boundaries

- **Quality over quantity** — If the content is shallow or low-quality, say so. Don't create Knowledge entry for junk.
- **Filter aggressively** — Not everything needs a Knowledge entry. Only distill if it has reusable value.
- **One entry per concept** — If similar knowledge exists, append to existing entry instead of creating new.
- **One skill per workflow** — If a similar skill exists in `~/.config/opencode/skills/`, extend it instead of creating a new one.
- **Never forge silently** — Skills are team-shared capability; always ask Louis WITH Aiy's recommendation attached, then Aiy reviews before writing.
- **No skill for facts** — The FORGE CHECK gate exists to keep the skill library lean. Facts/reference → Knowledge entry only.
- **Log the learning** — Always add a note in daily Log: "Learned: {topic}" (+ "Forged skill: {name}" when applicable).

## Relationship to Socratic Session

| Protocol | When | Output |
|----------|------|--------|
| Teach Protocol | Louis wants to learn something NEW | Knowledge entry + INDEX update (+ skill ถ้าเป็น workflow & Louis อนุมัติ) |
| Socratic Session | Louis wants to DECIDE something | Deeper thinking + synthesis |

Teach = input. Socratic = distillation of Louis's own wisdom.

## Skill-Forge Governance Summary

```
DISTILL เสร็จ → FORGE CHECK → workflow? ──ไม่ใช่──→ Knowledge entry อย่างเดียว
                                    │
                                    └─ใช่─→ เสนอ Louis (แนบข้อเสนออัย) → อนุมัติ?
                                                      │
                                                      ├─ไม่อนุมัติ─→ Knowledge entry อย่างเดียว
                                                      └─อนุมัติ─→ อัย review (duplicate/quality)
                                                                  → Forge SKILL.md → Link กลับ Knowledge
```

## Oracle Framework Reference
- Adapted from: [Oracle Framework](https://github.com/Soul-Brews-Studio/oracle-framework)
- Original concepts: rrr (Recursive Reflection & Refinement), distill pipeline, ψ (Collective Resonance)
- Knowledge entry: `Workspace/Aiy/Knowledge/Knowledge-Teach-Protocol.md`

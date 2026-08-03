---
name: affiliate-poster
description: "Aiy's Affiliate Content Pipeline — manage product research, content creation, and Facebook auto-posting for the 'ของมันต้องมี' page. Use ONLY when Louis mentions the affiliate/Facebook content pipeline or its assets: 'ของมันต้องมี' page, Facebook post, Shopee products, info.md drafts, post scheduling. NOT for general 'content/product/draft' talk (those are ambiguous). Triggers: 'ของมันต้องมี', 'facebook post', 'schedule post', 'affiliate', 'Shopee', 'info.md', 'publish post', 'โพสต์'."
---

# 🛍️ Affiliate Content Pipeline — "ของมันต้องมี"

> Full lifecycle: Product Research → info.md → Images → Content Draft → Facebook Post

## Project Structure

```
{{HOME}}/01-Projects/affiliate-content/
├── .env                          # Facebook API token (chmod 600, source to use)
├── README.md                     # Workflow docs + product status
├── NEW-content-drafts.md         # Cher's content drafts (002-006)
├── fb_poster.py                  # Core Facebook Graph API v20.0 client
├── content_pipeline.py           # Product publishing orchestrator
└── products/
    ├── 002-3in1-wireless-charger/
    │   ├── info.md               # Product specs, pros/cons, price
    │   └── images/               # Product photos (ringke-3in1-*.webp)
    ├── 003-satin-plus-bedset/
    ├── 004-vacuum-sealer/
    ├── 005-minimal-desk-fan/
    └── 006-cable-management-box/
```

## Quick Reference

```bash
# Pipeline status — see all products + readiness
cd {{HOME}}/01-Projects/affiliate-content
python3 content_pipeline.py status

# Preview a post without publishing
python3 content_pipeline.py publish 002 --dry-run

# Publish immediately
source .env && python3 content_pipeline.py publish 002

# Schedule for later
source .env && python3 content_pipeline.py publish 002 --schedule 2026-07-28T10:00

# Quick text-only post
source .env && python3 fb_poster.py publish -m "ข้อความ"

# Quick photo post
source .env && python3 fb_poster.py publish -m "ข้อความ" -p images/*.jpg

# List / delete posts
source .env && python3 fb_poster.py list
source .env && python3 fb_poster.py delete <post_id>
```

## The Pipeline Flow

```
Louis: "check this product" / "post product X"
  │
  ├─ New product research ──→ 1. Search Shopee / web for product
  │                             2. Create products/{NNN}-{slug}/ folder
  │                             3. Write info.md (specs, pros, cons, price)
  │                             4. Download product images to images/
  │                             5. Cher writes content draft to NEW-content-drafts.md
  │                             6. Update README + summary table
  │
  └─ Existing product post ──→ 1. Verify product folder + images + draft ready
                                2. Run --dry-run to preview
                                3. Get Louis approval + Shopee affiliate link
                                4. Update draft with affiliate link
                                5. Remove markdown formatting (Facebook no render)
                                6. source .env && python3 content_pipeline.py publish {NNN}
```

## Step-by-Step

### Step 1: New Product Research

When Louis shares a Shopee link or product idea:

```
1. Fetch the Shopee page:
   - Try web search for product details + images
   - Shopee API blocked (error 90309999 + CAPTCHA)
   - Alternative: search retailers / official brand sites / Google Images

2. Create product folder:
   mkdir -p products/{NNN}-{slug}/images

3. Write info.md with template:
   - # Product Name
   - ## ข้อมูลสินค้า (brand, model, price range, colors)
   - ## Specs (table if multiple specs)
   - ## จุดเด่น (bullet points ✅)
   - ## ข้อเสีย (bullet points ❌)
   - ## เหมาะกับ (target audience)
   - ## แหล่งอ้างอิง (reference URLs)

4. Download images to images/ folder
```

### Step 2: Content Draft

Delegate to Cher for content writing or write directly:

```
Update NEW-content-drafts.md with section:
- ## {NNN}. {emoji} {Product Name}
- ### Headline (one line, attention-grabbing, Thai)
- ### Body (400-800 chars, casual tone, no markdown bold)
- ### Call to Action (with [ลิงก์สินค้า] placeholder or actual URL)
- ### Hashtags (#Thai #tags #related)

Rules:
- NO markdown formatting (Facebook doesn't render **bold**, etc.)
- NO foreign language text (pure Thai, or English brand names only)
- Keep emojis ✅
- Include honest pros ✅ and cons ❌
- Call to Action should encourage comments or clicks
```

### Step 3: Pre-Publish Checklist

Before posting, verify:

- [ ] Product folder exists with info.md + images
- [ ] Content draft in NEW-content-drafts.md
- [ ] Draft has NO markdown formatting (**, __, etc.)
- [ ] Draft has NO non-Thai language text
- [ ] Shopee affiliate link is included (not placeholder)
- [ ] Run `python3 content_pipeline.py publish {NNN} --dry-run` to preview
- [ ] Louis has approved the preview

### Step 4: Publish

```bash
source {{HOME}}/01-Projects/affiliate-content/.env
python3 {{HOME}}/01-Projects/affiliate-content/content_pipeline.py publish {NNN}
```

The pipeline will:
1. Auto-acquire page token from user token (/me/accounts)
2. Upload all images with published=false
3. Create feed post with attached_media (album)
4. Return post ID

### Step 5: Log & Knowledge

- Log the post in `/home/lu5her/myObsidian/Workspace/Aiy/Logs/2026-MM.md` (absolute path เสมอ)
- If a new product category/pattern emerges → create Knowledge entry
- Update README product status table

## Facebook API Notes

- **API Version:** v20.0 (v21.0 returns #10 error)
- **Token:** User token (long-lived) stored in `.env`
- **Page Token:** Auto-exchanged from user token via `/me/accounts`
- **Permissions needed:** `pages_manage_posts`, `pages_read_engagement`
- **Multi-photo posts:** Upload each photo with `published=false`, then create feed post with `attached_media` JSON array
- **Scheduled posts:** Pass `published=false` + `scheduled_publish_time` (Unix timestamp)

## Known Limitations

| Issue | Workaround |
|-------|-----------|
| Shopee API blocked (90309999) | Search via web / retailers / brand sites |
| Shopee CAPTCHA | Can't fetch product pages directly |
| Facebook v21.0 #10 error | Use v20.0 |
| No Chrome/Playwright | Use curl + web search instead |

## File Descriptions

| File | Purpose |
|------|---------|
| `fb_poster.py` | Core Facebook Graph API module (publish, list, delete, debug) |
| `content_pipeline.py` | Product orchestrator (reads folders + drafts, coordinates posting) |
| `NEW-content-drafts.md` | Cher's content drafts for all products |
| `products/{NNN}-{slug}/info.md` | Product specs and details |
| `products/{NNN}-{slug}/images/` | Product photos (jpg, png, webp) |
| `.env` | Facebook token (chmod 600, never commit) |

## Related Resources

- Project: `{{HOME}}/01-Projects/affiliate-content/`
- Facebook Page: "ของมันต้องมี" (ID: 1210049942192010)
- Agent: Aiy (aiy.md) — Strategic Orchestrator
- Cher (cher.md) — Content & Creative Pipeline

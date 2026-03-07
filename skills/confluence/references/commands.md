# orbit Confluence Commands Reference

Complete reference for all `orbit confluence` commands with flags and examples.

## Table of Contents

- [Global Flags](#global-flags)
- [page — View Page](#page)
- [children — List Child Pages](#children)
- [create — Create Page](#create)
- [update — Update Page](#update)
- [publish — Publish Directory](#publish)
- [set-width — Set Page Width](#set-width)

---

## Global Flags

These flags are available on all confluence subcommands:

| Flag | Description |
|------|-------------|
| `-p, --profile` | Profile to use (overrides default) |
| `-o, --output` | Output format: table, json, yaml (default: table) |
| `--service` | Confluence service name (if profile has multiple) |
| `--config` | Config file path (default: ~/.config/orbit/config.yaml) |

---

## page

View a Confluence page's details including title, version, space, and URL.

```bash
orbit -p profile confluence page <page-id> [flags]
```

**Examples:**

```bash
orbit -p paybook confluence page 473676972036
orbit -p paybook confluence page 473676972036 -o json
```

**Output (table format):**

```
ID:      473676972036
Title:   AI Development Process
Version: 3
Space:   FO
URL:     https://paybook.atlassian.net/wiki/spaces/FO/pages/473676972036
```

**Output (json format):** Full page object including `body.storage.value` with the page content in Confluence storage format.

---

## children

List child pages of a given page.

```bash
orbit -p profile confluence children <page-id> [flags]
```

**Examples:**

```bash
orbit -p paybook confluence children 473677299713
orbit -p paybook confluence children 473677299713 -o json
```

**Output:**

```
ID              VERSION  TITLE
--              -------  -----
473676972036    3        AI Development Process
473677103107    2        Foundations
473676873739    1        Organization
```

---

## create

Create a new Confluence page. Automatically sets wide width on the created page.

```bash
orbit -p profile confluence create [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--title` | string | Page title (**required**) |
| `--space` | string | Space key (**required**) |
| `--parent` | string | Parent page ID |
| `-b, --body` | string | Page body in storage format (XHTML) |
| `-f, --file` | string | Markdown file to convert and upload |

**Examples:**

```bash
# From markdown file
orbit -p paybook confluence create --space FO --parent 473677299713 \
  --title "Sprint Ceremonies" --file docs/sprint-ceremonies.md

# With inline content
orbit -p paybook confluence create --space FO --parent 473677299713 \
  --title "Quick Note" --body "<p>Remember to update the docs</p>"
```

**Notes:**
- When `--file` is provided, the markdown is automatically converted to Confluence storage format
- The `--body` flag expects raw Confluence storage format (XHTML), not markdown
- Created pages are automatically set to full-width layout

---

## update

Update an existing Confluence page. Automatically increments the version number.

```bash
orbit -p profile confluence update <page-id> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--title` | string | New title (keeps current if empty) |
| `-b, --body` | string | New body in storage format (XHTML) |
| `-f, --file` | string | Markdown file to convert and upload |

**Examples:**

```bash
# Update content from markdown
orbit -p paybook confluence update 473676972036 --file docs/overview.md

# Update title and content
orbit -p paybook confluence update 473676972036 \
  --title "Updated Overview" --file docs/overview.md

# Update with inline content
orbit -p paybook confluence update 473676972036 \
  --body "<p>Updated content</p>"
```

**Notes:**
- The command automatically fetches the current page version and increments it
- If `--title` is not provided, the existing title is preserved
- When using `--file`, the title from the markdown file is NOT used — provide `--title` explicitly if you want to change it

---

## publish

Recursively publish a directory of markdown files to Confluence, preserving folder hierarchy.

```bash
orbit -p profile confluence publish <directory> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--space` | string | Confluence space key (**required**) |
| `--parent` | string | Parent page ID (**required**) |
| `--dry-run` | bool | Preview without creating pages |

**Examples:**

```bash
# Preview the publish plan
orbit -p paybook confluence publish ./docs --space FO --parent 473677299713 --dry-run

# Publish for real
orbit -p paybook confluence publish ./docs --space FO --parent 473677299713
```

**Directory Structure Rules:**

```
docs/
├── INDEX.md          → becomes parent page for "docs"
├── overview.md       → child page under docs
├── quick-ref.md      → child page under docs
├── foundations/
│   ├── INDEX.md      → becomes parent page for "foundations" (under docs)
│   ├── concepts.md   → child page under foundations
│   └── setup.md      → child page under foundations
└── workflow/
    ├── INDEX.md      → becomes parent page for "workflow" (under docs)
    └── sprints.md    → child page under workflow
```

**Title Resolution Order:**
1. YAML frontmatter `title:` field
2. First `# heading` in the file
3. Filename converted to title case (e.g., `quick-ref.md` → "Quick Ref")

**Frontmatter Fields:**

The `publish` command recognizes the following YAML frontmatter fields:

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Page title (overrides filename) |
| `confluence_page_id` | string | Existing page ID for updates |
| `confluence_url` | string | Page URL (set after publish) |
| `confluence_labels` | list | Labels/tags applied to the page after create/update |
| `confluence_properties` | map | Page Properties macro prepended to page content |

### Labels

Add `confluence_labels` to frontmatter to tag pages with Confluence labels. Labels are applied via a POST to `/content/{id}/label` after each page is created or updated.

```yaml
---
title: "My Policy"
confluence_labels:
  - ai-process
  - foundation
---
```

### Page Properties (details macro)

Add `confluence_properties` to frontmatter to generate a Page Properties macro (`ac:structured-macro ac:name="details"`) prepended to the page content. This creates a structured metadata block that can be queried by Page Properties Report macros on other pages.

```yaml
---
title: "Compounding Engineering & System Evolution"
confluence_properties:
  id: status
  fields:
    Owner: AI Tooling Guild
    Classification: Internal
    Status: "{status:Green|approved}"
    Reviewed on: 2026-03-06
    Approved on: 2026-03-06
---
```

**Sub-keys:**

| Key | Description |
|-----|-------------|
| `id` | Macro ID used by `detailssummary` reports to target specific property blocks |
| `fields` | Ordered key-value pairs rendered as a two-column table inside the macro |

**Field value formats:**

| Format | Rendered As |
|--------|-------------|
| `{status:Color\|Text}` | Status badge macro (e.g., `{status:Green\|approved}`) |
| `YYYY-MM-DD` | Confluence `<time>` date macro |
| Plain text | Plain text |

### Page Properties Report (detailssummary macro)

Use either directive format in markdown to generate a dynamic table that pulls Page Properties from child or labeled pages.

**HTML comment format (full control):**
```markdown
<!-- confluence:properties-report cql="label = 'ai-process' and space = currentSpace()" firstcolumn="Document" headings="Status, Classification, Reviewed on" -->
```

**Shorthand format:**
```markdown
{properties-report: label="ai-process", columns="Status, Classification, Reviewed on"}
```

**Parameters:**

| Parameter | Description |
|-----------|-------------|
| `cql` | CQL query to find pages (e.g., `label = "policy" and space = currentSpace()`) |
| `label` | Shorthand for CQL: generates `label = "value" and space = currentSpace()` |
| `firstcolumn` | Name of the first column (defaults to "Title") |
| `columns` / `headings` | Comma-separated list of property names to show as columns |
| `sortBy` | Optional column to sort by |

**Notes:**
- Subdirectories are processed before sibling files
- Hidden directories (starting with `.`) are skipped
- Only `.md` files are processed
- Each `INDEX.md` creates a parent page; other `.md` files create child pages under it
- Labels are applied after page creation/update via the Confluence REST API
- Page Properties macros are prepended to the converted page body before upload
- Properties Report directives are converted inline during markdown-to-storage-format conversion
- `<!-- confluence:ignore-start -->` / `<!-- confluence:ignore-end -->` blocks are stripped — content between the markers is completely skipped during conversion, useful for static markdown tables that should not appear on Confluence

---

## set-width

Set the content width of one or more pages.

```bash
orbit -p profile confluence set-width <page-id...> [flags]
```

**Flags:**

| Flag | Type | Description |
|------|------|-------------|
| `--width` | string | Page width: wide or fixed (default: wide) |
| `--recursive` | bool | Apply to all child pages recursively |

**Examples:**

```bash
# Set single page to wide
orbit -p paybook confluence set-width 473676972036

# Set multiple pages
orbit -p paybook confluence set-width 473676972036 473677103107

# Set page and all descendants to wide
orbit -p paybook confluence set-width 473677299713 --recursive

# Set to fixed width
orbit -p paybook confluence set-width 473676972036 --width fixed
```

**Notes:**
- Sets both draft and published appearance properties
- "wide" maps to `full-width` appearance, "fixed" maps to `fixed`
- With `--recursive`, traverses all child pages depth-first

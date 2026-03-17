# Formatting Reference

Detailed rules for each formatting check. The SKILL.md references this file for specifics.

## Table of Contents

1. [Frontmatter Fields](#frontmatter-fields)
2. [Heading Rules](#heading-rules)
3. [Table of Contents Generation](#table-of-contents-generation)
4. [Cross-References](#cross-references)
5. [INDEX.md Templates](#indexmd-templates)
6. [Confluence Properties](#confluence-properties)
7. [Markdown Elements Conversion Table](#markdown-elements-conversion-table)
8. [Common Anti-Patterns](#common-anti-patterns)
   - [7. Redundant Inline Metadata](#7-redundant-inline-metadata)

---

## Frontmatter Fields

### Standard Fields (always populate all)

| Field | Type | Source | Example |
|-------|------|--------|---------|
| `title` | String | First `# heading`, or filename | `"API Reference"` |
| `subtitle` | String | Document type + project context | `"Reference — Project Docs v1.0"` |
| `date` | String | Today's date or file mtime | `"2026-03-09"` |
| `author` | String | Document owner / team | `"Engineering Team"` |
| `confluence_ignore` | Boolean | When `true`, skip publishing (delete if previously published) | `false` |
| `confluence_labels` | Array | Confluence page labels for search | `[api, reference, foundation]` |
| `confluence_properties` | Array | Page Properties macro metadata | See [Confluence Properties](#confluence-properties) |

All six fields should be present in every document. The frontmatter is the **single source of truth** for document metadata — information here should not be duplicated in the document body.

### Auto-Managed Fields (do not edit)

| Field | Purpose |
|-------|---------|
| `confluence_page_id` | Set by `orbit confluence publish` after first sync |
| `confluence_url` | Set by `orbit confluence publish` — direct link |

### Deriving the Title

Order of precedence:

1. **Existing `title` in frontmatter** — keep it
2. **First `# heading`** — extract text, strip inline formatting
3. **Filename** — convert to title case: `api-reference.md` → `"Api Reference"`

When extracting from a `# heading`, strip markdown formatting:
- `# **Bold Title**` → `"Bold Title"`
- `# Title with `code`` → `"Title with code"`
- `# Title [with link](url)` → `"Title with link"`

### Deriving Labels

Suggest labels based on:
- Directory name (e.g., files in `workflow/` get label `workflow`)
- Common project labels if a pattern exists across sibling files
- Content type: `reference`, `guide`, `process`, `architecture`

---

## Heading Rules

### The H1 Rule

Confluence displays the page title (from frontmatter `title:` or the first `# heading`) separately. The markdown converter skips the first `#` heading to avoid duplication. This means:

- A file should have exactly **one** `#` heading (or zero, if `title` is in frontmatter)
- All content sections should start at `##`

### Fixing Multiple H1s

Before:
```markdown
# Introduction
Some content...
# Architecture
More content...
# Deployment
Even more...
```

After:
```markdown
# Introduction
Some content...
## Architecture
More content...
## Deployment
Even more...
```

The first `#` stays as the page title source. All subsequent `#` become `##`, and their subsections shift accordingly.

### Fixing Heading Gaps

Before:
```markdown
## Overview
#### Details
###### Deep detail
```

After:
```markdown
## Overview
### Details
#### Deep detail
```

Each level should increment by one. Never skip from `##` to `####`.

---

## Table of Contents Generation

Add a TOC section when a document has 3 or more `##` headings. Wrap the TOC section with `<!-- confluence:toc-start -->` / `<!-- confluence:toc-end -->` directives so it gets replaced by the Confluence TOC macro on publish. The markdown TOC remains visible in GitHub/local renderers.

### Format

```markdown
<!-- confluence:toc-start -->

## Table of Contents

1. [Section Name](#section-name)
2. [Another Section](#another-section)
3. [Cross-References](#cross-references)
4. [Revision History](#revision-history)

<!-- confluence:toc-end -->
```

Place the TOC immediately after the title and optional summary blockquote, before the first content section. The numbered list with anchor links serves as a clickable TOC in markdown renderers.

### TOC Directive Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `maxLevel` / `maxHeadingLevel` | `3` | Maximum heading depth to include |
| `minLevel` / `minHeadingLevel` | `2` | Minimum heading depth to include |
| `style` | `none` | List style (`none`, `circle`, `disc`, `square`, `decimal`) |
| `outline` | `false` | Show section numbers (set `true` to enable numbered outline) |
| `printable` | `true` | Whether the TOC appears in print view |

Parameters can be quoted or unquoted: `maxLevel=2` or `maxLevel="2"` both work.

---

## Cross-References

Add a Cross-References section when the file is part of a multi-document set (directory with siblings).

### Format

```markdown
## Cross-References

| Document | Relevance |
|----------|-----------|
| [Agent Roles](./agent-roles.md) | Defines the agents referenced in this workflow |
| [PIV Loops](./piv-loops.md) | The execution cycle this process follows |
```

### Rules

- Use relative paths: `./sibling.md` for same directory, `../other-dir/file.md` for cross-directory
- Include 3-5 most relevant related documents, not every sibling
- Write a brief relevance note (one phrase) explaining the connection
- Place before Revision History, after all content sections

---

## INDEX.md Templates

### Root INDEX.md

The root `INDEX.md` is the landing page for the entire documentation set.

```markdown
---
title: Project Documentation
date: "2026-03-09"
author: "Engineering Team"
confluence_labels:
  - documentation
  - index
---

# Project Documentation

> Central index for all project documentation.

## Sections

| Section | Description |
|---------|-------------|
| [Foundations](./foundations/INDEX.md) | Core principles and philosophy |
| [Workflow](./workflow/INDEX.md) | Day-to-day engineering workflows |
| [Process](./process/INDEX.md) | Sprint ceremonies, metrics, tracking |
| [Reference](./reference/INDEX.md) | Quick-reference guides and cheat sheets |
```

### Category INDEX.md

```markdown
---
title: Workflow
subtitle: "How engineers work day-to-day"
date: "2026-03-09"
confluence_labels:
  - workflow
  - index
---

# Workflow

> Guides for the team's engineering workflow with AI agents.

| Document | Description |
|----------|-------------|
| [Agent Roles](./agent-roles.md) | The 5 agent roles and their responsibilities |
| [PIV Loops](./piv-loops.md) | Plan, Implement, Validate execution cycle |
| [Parallelization](./parallelization.md) | Running multiple agents concurrently |
```

### Key Patterns

- Keep INDEX.md files short — they're navigation pages, not content pages
- Use a table (not a list) to link child pages — tables render cleanly in Confluence
- Include a one-line description for each linked page
- List all `.md` files in the directory (except INDEX.md itself)
- List all subdirectories that have their own INDEX.md

---

## Confluence Properties

The `confluence_properties` frontmatter field renders as a Page Properties macro at the top of the Confluence page. This enables Page Properties Report macros on parent pages to aggregate metadata across child pages.

### Single Properties Block

```yaml
confluence_properties:
  - id: status
    fields:
      Owner: "Engineering Team"
      Classification: "Internal"
      Status: "{status:Green|Published}"
      Reviewed on: "2026-03-09"
```

### Field Value Types

| Type | Example | Confluence Rendering |
|------|---------|---------------------|
| Plain text | `"Engineering Team"` | Text |
| Status badge | `"{status:Green|Published}"` | Colored status lozenge |
| Date | `"2026-03-09"` | Date macro (when in properties) |
| Link | `"[Team Page](https://...)"` | Clickable link |

### Properties Report

On an INDEX.md page, you can add a directive that generates a table aggregating properties from all child pages:

```markdown
<!-- confluence:properties-report cql="label = 'workflow'" firstcolumn="Title" headings="Owner,Status,Reviewed on" -->
```

This is useful when all pages in a section share the same properties schema — the parent page automatically shows a summary table.

---

## Markdown Elements Conversion Table

Elements that convert well to Confluence and how to use them:

| Markdown | Confluence Result | Notes |
|----------|-------------------|-------|
| `## Heading` | `<h2>` | Use `##` as top section level |
| `**bold**` | Bold text | Standard |
| `*italic*` | Italic text | Standard |
| `` `inline code` `` | Monospace span | Use for field names, commands, paths |
| ` ```lang ` code block | Code macro with syntax highlighting | Always include language hint |
| `> blockquote` | Info panel macro | Use for callouts and notes |
| `> [!WARNING]` | Warning panel macro | GitHub-style alerts supported |
| `- list item` | Bullet list | Max 2 levels of nesting |
| `1. list item` | Numbered list | Good for procedures |
| Pipe tables | HTML tables | Keep cells concise |
| `[text](url)` | Link | Relative `.md` links resolved on publish |
| `![alt](path)` | Image macro | Relative paths resolved |
| `---` | Horizontal rule | Section separators |
| `<!-- confluence:toc-start -->` | TOC macro | Wrap TOC section with start/end directives |
| `- [ ] task` | Status macro (incomplete) | Task list items |
| `- [x] task` | Status macro (complete) | Checked task items |
| `<details><summary>` | Expand macro | Collapsible sections |

---

## Common Anti-Patterns

Issues to detect and fix:

### 1. No Frontmatter

**Before:**
```markdown
# My Document
Content starts here...
```

**After:**
```markdown
---
title: My Document
date: "2026-03-09"
---

# My Document
Content starts here...
```

### 2. Title Mismatch

The `title` frontmatter and `#` heading should match. If they differ, the frontmatter `title` wins for the Confluence page title, and the `#` heading is skipped — potentially losing information.

**Fix:** Make them match, or remove the `#` heading and rely solely on frontmatter `title`.

### 3. Spaces in Filenames

**Before:** `API Reference Guide.md`
**After:** `api-reference-guide.md`

Update all internal links that reference the renamed file.

### 4. Missing INDEX.md

A directory without INDEX.md will not have a proper parent page in Confluence. The files will still publish but without a navigation landing page.

### 5. Deeply Nested Lists

Lists nested 3+ levels render poorly in Confluence.

**Before:**
```markdown
- Level 1
  - Level 2
    - Level 3
      - Level 4
```

**After:**
```markdown
- Level 1
  - Level 2: Level 3 — Level 4
```

Or restructure as separate sections.

### 6. Raw HTML

Most raw HTML is stripped or mangled during conversion. Replace with markdown equivalents.

| HTML | Markdown Replacement |
|------|---------------------|
| `<br>` | Blank line |
| `<b>text</b>` | `**text**` |
| `<i>text</i>` | `*text*` |
| `<a href="url">text</a>` | `[text](url)` |
| `<img src="url">` | `![alt](url)` |
| `<table>` complex tables | Pipe table syntax |

Exception: `<!-- confluence: -->` HTML comments are intentional directives and should be preserved.

### 7. Redundant Inline Metadata

Document metadata that duplicates frontmatter or `confluence_properties` must be removed from the body. The frontmatter is the single source of truth.

**Before:**
```markdown
---
title: "Programmatic Tool Execution"
subtitle: "ADR — EPAP v2.0"
date: "2026-03-07"
author: "Platform Team"
confluence_properties:
  - id: status
    fields:
      Owner: "Platform Team"
      Classification: Internal
      Status: "{status:Yellow|Draft}"
      Version: v1.0
      Reviewed on: "2026-03-07"
---

# Programmatic Tool Execution

**Document Version**: 1.0
**Status**: Draft (Proposal)
**Classification**: Internal
**Last Updated**: March 7, 2026
**Author**: Platform Team

---

## Overview
Content starts here...
```

**After:**
```markdown
---
title: "Programmatic Tool Execution"
subtitle: "ADR — EPAP v2.0"
date: "2026-03-07"
author: "Platform Team"
confluence_properties:
  - id: status
    fields:
      Owner: "Platform Team"
      Classification: Internal
      Status: "{status:Yellow|Draft}"
      Version: v1.0
      Reviewed on: "2026-03-07"
---

# Programmatic Tool Execution

> Architecture Decision Record for code-based tool orchestration.

## Overview
Content starts here...
```

**Detection patterns** — look for these near the top of documents (typically between the `#` heading and the first `##` section):

| Pattern | Example |
|---------|---------|
| Bold key-value | `**Status**: Draft` |
| Plain key-value | `Version: 1.0` |
| Metadata table | `\| Field \| Value \|` with rows like `\| Owner \| Team \|` |
| Blockquote metadata | `> **Author**: Platform Team` |
| Italic key-value | `*Last Updated*: March 7, 2026` |

**Absorption rule:** If the inline metadata contains information NOT yet in the frontmatter, absorb it into the appropriate field first:

| Inline Key | Maps To |
|------------|---------|
| `Document Version`, `Version` | `confluence_properties.fields.Version` |
| `Status` | `confluence_properties.fields.Status` (wrap with `{status:Color\|Text}`) |
| `Classification` | `confluence_properties.fields.Classification` |
| `Last Updated`, `Date`, `Updated` | `date` (frontmatter) |
| `Author`, `Owner`, `Maintained by` | `author` (frontmatter) and `confluence_properties.fields.Owner` |
| `Reviewed`, `Reviewed on` | `confluence_properties.fields.Reviewed on` |

### 8. Orphaned Links

After renaming files or reorganizing directories, check that all relative links (`[text](./path.md)`) still point to valid files. Broken links in Confluence show as dead links.

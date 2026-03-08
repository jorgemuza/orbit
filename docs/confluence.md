# Confluence Command Reference

Command reference for `orbit confluence` -- manage Confluence pages from the terminal.

## Table of Contents

- [Global Flags](#global-flags)
- [page](#page) -- View a Confluence page
- [children](#children) -- List child pages
- [create](#create) -- Create a new page
- [update](#update) -- Update an existing page
- [delete](#delete) -- Delete a page
- [publish](#publish) -- Publish a directory of markdown files
- [set-width](#set-width) -- Set page width
- [Notes](#notes)

---

## Global Flags

| Flag | Description |
|------|-------------|
| `-p, --profile` | Profile to use (applies to all `orbit` commands) |
| `--service` | Confluence service name, when a profile has multiple Confluence services configured |

---

## page

View a Confluence page. Outputs the page title, space, version number, URL, and body content.

```
orbit confluence page [page-id] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id` | The numeric ID of the Confluence page |

### Examples

```bash
# View a page by ID
orbit confluence page 123456789 -p myprofile

# View a page using a specific service
orbit confluence page 123456789 -p myprofile --service wiki-prod
```

---

## children

List the child pages of a given Confluence page.

```
orbit confluence children [page-id] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id` | The numeric ID of the parent Confluence page |

### Examples

```bash
# List all child pages
orbit confluence children 123456789 -p myprofile
```

---

## create

Create a new Confluence page. Pages created via `orbit` are automatically set to wide width.

```
orbit confluence create [flags]
```

### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--title` | Yes | Title for the new page |
| `--space` | Yes | Space key where the page will be created |
| `--parent` | No | Parent page ID to nest the new page under |
| `-b, --body` | No | Page body content as a string (Confluence storage format / XHTML) |
| `-f, --file` | No | Path to a markdown file to convert and upload as the page body |

### Examples

```bash
# Create a page with inline body content
orbit confluence create --title "Sprint Retrospective" --space DEV --parent 123456789 -p myprofile

# Create a page from a markdown file
orbit confluence create --title "Architecture Overview" --space ENG --parent 123456789 -f ./docs/architecture.md -p myprofile

# Create a top-level page in a space (no parent)
orbit confluence create --title "Project Home" --space PROJ -b "<p>Welcome to the project.</p>" -p myprofile
```

---

## update

Update an existing Confluence page. You can change its title, body, or both.

```
orbit confluence update [page-id] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id` | The numeric ID of the page to update |

### Flags

| Flag | Description |
|------|-------------|
| `--title` | New title for the page |
| `-b, --body` | New body content as a string |
| `-f, --file` | Path to a markdown file to convert and use as the new body |

### Examples

```bash
# Update a page title
orbit confluence update 123456789 --title "Updated Title" -p myprofile

# Update a page body from a markdown file
orbit confluence update 123456789 -f ./docs/revised-design.md -p myprofile

# Update both title and body
orbit confluence update 123456789 --title "Q2 Plan (Final)" -b "<p>Finalized plan.</p>" -p myprofile
```

---

## delete

Delete a Confluence page. The page is moved to the Confluence trash.

```
orbit confluence delete [page-id] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id` | The numeric ID of the page to delete |

### Examples

```bash
# Delete a page (moves to trash)
orbit confluence delete 123456789 -p myprofile
```

---

## publish

Publish an entire directory of markdown files to Confluence. This command creates or updates a tree of pages that mirrors the directory structure.

```
orbit confluence publish [directory] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `directory` | Path to the directory containing markdown files |

### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--space` | Yes | Confluence space key |
| `--parent` | Yes | Parent page ID under which the pages will be published |
| `--dry-run` | No | Preview what would be published without making any changes |

### Directory Structure Rules

- **`INDEX.md`** files become the parent page for their directory.
- Other **`.md`** files become child pages under that parent.
- **Subdirectories** are processed recursively, maintaining the hierarchy.

### Page Title Resolution

Page titles are determined in the following order of precedence:

1. YAML frontmatter `title:` field
2. First `# heading` in the file
3. Filename (without `.md` extension)

### Frontmatter Tracking

After a successful publish, each markdown file has `confluence_page_id` and `confluence_url` added to its YAML frontmatter. These fields are used on subsequent runs to update existing pages rather than creating duplicates.

### Examples

```bash
# Publish a docs directory
orbit confluence publish ./docs --space ENG --parent 123456789 -p myprofile

# Preview what would be published (no changes made)
orbit confluence publish ./docs --space ENG --parent 123456789 --dry-run -p myprofile
```

Given this directory:

```
docs/
  INDEX.md
  getting-started.md
  api/
    INDEX.md
    endpoints.md
    authentication.md
```

The resulting Confluence page tree would be:

```
Parent Page (123456789)
  +-- docs (from INDEX.md)
        +-- Getting Started (from getting-started.md)
        +-- api (from api/INDEX.md)
              +-- Endpoints (from endpoints.md)
              +-- Authentication (from authentication.md)
```

---

## set-width

Set the display width of one or more Confluence pages. Accepts multiple page IDs and can optionally apply the change recursively to all child pages.

```
orbit confluence set-width [page-id...] [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `page-id...` | One or more page IDs to update |

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--width` | `wide` | Width setting: `wide` or `fixed` |
| `--recursive` | `false` | Apply the width setting to all descendant pages |

### Examples

```bash
# Set a single page to wide width
orbit confluence set-width 123456789 -p myprofile

# Set multiple pages to fixed width
orbit confluence set-width 111111111 222222222 --width fixed -p myprofile

# Set a page and all its children to wide width
orbit confluence set-width 123456789 --recursive -p myprofile
```

---

## Notes

### Confluence API Variants

Orbit automatically selects the correct API path based on the configured instance type:

| Variant | API Prefix |
|---------|------------|
| Cloud | `/wiki/rest/api` |
| Server / Data Center | `/rest/api` |

### Markdown Conversion

When a markdown file is provided via `--file` / `-f` (or through `publish`), Orbit converts it to Confluence storage format (XHTML). The converter handles:

- Headings (`#`, `##`, etc.)
- Ordered and unordered lists
- Tables
- Fenced code blocks
- Blockquotes
- Inline formatting (bold, italic, code, links)
- Images
- Table of Contents sections
- Metadata lines

Additional conversion rules:

- The first `# heading` is **skipped** because Confluence already displays the page title.
- YAML frontmatter is **stripped** during conversion.

# Jira Wiki Markup Reference

Jira Server and Data Center use **wiki markup** for rich text formatting in descriptions, comments, and text fields. This is NOT Markdown — using Markdown syntax will render as plain text.

## Headings

```
h1. Largest heading
h2. Section heading
h3. Subsection heading
h4. Small heading
h5. Smallest heading
h6. Tiny heading
```

Use `h2.` for main sections in epic/story descriptions. Always leave a blank line after the heading.

## Text Effects

| Markup | Result |
|--------|--------|
| `*bold*` | **bold** |
| `_italic_` | *italic* |
| `-strikethrough-` | ~~strikethrough~~ |
| `+underline+` | underline |
| `{{monospaced}}` | `monospaced` |
| `^superscript^` | superscript |
| `~subscript~` | subscript |
| `??citation??` | citation |

## Lists

### Bullet Lists

```
* Item one
* Item two
** Nested item
** Another nested
*** Deep nested
* Back to top level
```

### Numbered Lists

```
# First item
# Second item
## Sub-item a
## Sub-item b
# Third item
```

### Mixed Lists

```
* Bullet item
*# Numbered sub-item
*# Another numbered sub-item
* Another bullet
```

## Links

```
[Link text|https://example.com]
[PROJ-123]                          (auto-links to Jira issue)
[Link to anchor|#anchor-name]
[mailto:user@example.com]
```

## Tables

```
||Header 1||Header 2||Header 3||
|Cell 1|Cell 2|Cell 3|
|Cell 4|Cell 5|Cell 6|
```

Tables must use `||` for header cells and `|` for regular cells. Each row must be on its own line.

## Code and Preformatted Text

### Inline Code

```
{{my_variable}} or {{some.method()}}
```

### Code Blocks

```
{code:python}
def hello():
    print("Hello World")
{code}
```

```
{code:bash}
#!/bin/bash
echo "Hello"
{code}
```

### Preformatted (No Syntax Highlighting)

```
{noformat}
This is preformatted text.
  Whitespace is preserved.
{noformat}
```

## Panels and Boxes

```
{panel:title=My Panel Title}
Content inside the panel.
{panel}
```

```
{info:title=Information}
This is an info box.
{info}
```

```
{warning:title=Warning}
This is a warning box.
{warning}
```

```
{note:title=Note}
This is a note box.
{note}
```

## Images

```
!image.png!                         (attached image)
!image.png|width=300!               (with size)
!https://example.com/image.png!     (external URL)
```

## Horizontal Rule

```
----
```

## Quotes

```
{quote}
This is a block quote.
Multiple lines are supported.
{quote}
```

## Line Breaks

A single newline in wiki markup is generally rendered as a space. Use `\\` for a forced line break or leave a blank line for a paragraph break.

```
Line one\\
Line two (forced break)

Line three (paragraph break above)
```

## Colors

```
{color:red}This text is red.{color}
{color:#0000ff}This text is blue.{color}
```

## Anchors

```
{anchor:my-section}
[Jump to section|#my-section]
```

## Emoticons

```
:)  :D  :(  :P  ;)  (y)  (n)  (i)  (/)  (x)  (!)  (+)  (-)  (?)  (on)  (off)  (*)
```

---

## Epic Description Template

When creating epics, use this structure for well-formatted descriptions:

```
h2. Value Statement

[One paragraph describing the value this epic delivers and who benefits.]

h2. Dependencies

* [Dependency description, e.g., "Depends on: Epic Name (E1) — reason"]
* None — can start immediately

h2. User Stories

* *Story 1:* As a [role], I want to [action] so that [benefit].
** [Acceptance criterion 1]
** [Acceptance criterion 2]
** [Acceptance criterion 3]
* *Story 2:* As a [role], I want to [action] so that [benefit].
** [Acceptance criterion 1]
** [Acceptance criterion 2]

h2. Functional Requirements

||Req ID||Description||Priority||
|FR-001|[Requirement description]|Must|
|FR-002|[Requirement description]|Should|
|FR-003|[Requirement description]|Could|
```

## Story Description Template

```
h2. Summary

[Brief summary of what this story delivers.]

h2. Acceptance Criteria

* [Criterion 1]
* [Criterion 2]
* [Criterion 3]

h2. Technical Notes

* [Any implementation details, constraints, or considerations]
```

## Common Mistakes

| Wrong (Markdown) | Right (Wiki Markup) |
|-------------------|---------------------|
| `## Heading` | `h2. Heading` |
| `**bold**` | `*bold*` |
| `*italic*` | `_italic_` |
| `` `code` `` | `{{code}}` |
| `- bullet` | `* bullet` |
| `1. numbered` | `# numbered` |
| `[text](url)` | `[text\|url]` |
| `\| h1 \| h2 \|` | `\|\|h1\|\|h2\|\|` |
| ` ``` code ``` ` | `{code}...{code}` |

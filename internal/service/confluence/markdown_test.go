package confluence

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// GitHub Alerts
// ---------------------------------------------------------------------------

func TestGitHubAlertNote(t *testing.T) {
	md := "> [!NOTE]\n> This is a note"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="note"`) {
		t.Errorf("expected note macro, got: %s", got)
	}
	if !strings.Contains(got, "This is a note") {
		t.Errorf("expected note content, got: %s", got)
	}
}

func TestGitHubAlertWarning(t *testing.T) {
	md := "> [!WARNING]\n> Be careful here"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="warning"`) {
		t.Errorf("expected warning macro, got: %s", got)
	}
}

func TestGitHubAlertTip(t *testing.T) {
	md := "> [!TIP]\n> Use this shortcut"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="tip"`) {
		t.Errorf("expected tip macro, got: %s", got)
	}
}

func TestGitHubAlertCaution(t *testing.T) {
	md := "> [!CAUTION]\n> This might break things"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="warning"`) {
		t.Errorf("expected warning macro for caution, got: %s", got)
	}
}

func TestGitHubAlertImportant(t *testing.T) {
	md := "> [!IMPORTANT]\n> Read this first"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="info"`) {
		t.Errorf("expected info macro for important, got: %s", got)
	}
	if !strings.Contains(got, `title">Important`) {
		t.Errorf("expected Important title, got: %s", got)
	}
}

func TestPlainBlockquote(t *testing.T) {
	md := "> Just a regular quote"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="info"`) {
		t.Errorf("expected info macro for plain blockquote, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// Nested Lists
// ---------------------------------------------------------------------------

func TestNestedUnorderedList(t *testing.T) {
	md := "- item 1\n  - nested 1\n  - nested 2\n- item 2"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "<ul><li>item 1<ul><li>nested 1</li><li>nested 2</li></ul></li><li>item 2</li></ul>") {
		t.Errorf("unexpected nested list output: %s", got)
	}
}

func TestNestedOrderedList(t *testing.T) {
	md := "1. first\n   1. sub-first\n2. second"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "<ol>") {
		t.Errorf("expected ol tag, got: %s", got)
	}
	if !strings.Contains(got, "sub-first") {
		t.Errorf("expected nested content, got: %s", got)
	}
}

func TestDeeplyNestedList(t *testing.T) {
	md := "- level 0\n  - level 1\n    - level 2\n- back to 0"
	got := MarkdownToStorage(md)
	// Should have nested ul tags
	ulCount := strings.Count(got, "<ul>")
	if ulCount < 3 {
		t.Errorf("expected 3 nested ul tags, got %d in: %s", ulCount, got)
	}
}

// ---------------------------------------------------------------------------
// Task Lists
// ---------------------------------------------------------------------------

func TestTaskListUnchecked(t *testing.T) {
	md := "- [ ] Todo item"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "TODO") {
		t.Errorf("expected TODO status, got: %s", got)
	}
	if !strings.Contains(got, "Grey") {
		t.Errorf("expected Grey colour, got: %s", got)
	}
}

func TestTaskListChecked(t *testing.T) {
	md := "- [x] Done item"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "DONE") {
		t.Errorf("expected DONE status, got: %s", got)
	}
	if !strings.Contains(got, "Green") {
		t.Errorf("expected Green colour, got: %s", got)
	}
}

func TestTaskListMixed(t *testing.T) {
	md := "- [x] Done\n- [ ] Todo\n- Regular item"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "DONE") || !strings.Contains(got, "TODO") {
		t.Errorf("expected both DONE and TODO, got: %s", got)
	}
	if !strings.Contains(got, "Regular item") {
		t.Errorf("expected regular item, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// Table Alignment
// ---------------------------------------------------------------------------

func TestTableLeftAlign(t *testing.T) {
	md := "| Name | Age |\n|:---|---|\n| Alice | 30 |"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `text-align: left`) {
		t.Errorf("expected left alignment, got: %s", got)
	}
}

func TestTableCenterAlign(t *testing.T) {
	md := "| Name | Age |\n|:---:|:---:|\n| Alice | 30 |"
	got := MarkdownToStorage(md)
	if strings.Count(got, `text-align: center`) < 2 {
		t.Errorf("expected center alignment on both columns, got: %s", got)
	}
}

func TestTableRightAlign(t *testing.T) {
	md := "| Name | Amount |\n|---|---:|\n| Alice | $100 |"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `text-align: right`) {
		t.Errorf("expected right alignment, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// Image Alt Text
// ---------------------------------------------------------------------------

func TestImageWithAltText(t *testing.T) {
	md := "![My screenshot](https://example.com/img.png)"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:alt="My screenshot"`) {
		t.Errorf("expected alt text preserved, got: %s", got)
	}
	if !strings.Contains(got, `ri:value="https://example.com/img.png"`) {
		t.Errorf("expected image URL, got: %s", got)
	}
}

func TestImageWithoutAltText(t *testing.T) {
	md := "![](https://example.com/img.png)"
	got := MarkdownToStorage(md)
	if strings.Contains(got, `ac:alt`) {
		t.Errorf("expected no alt attribute for empty alt, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// Code Block Title
// ---------------------------------------------------------------------------

func TestCodeBlockWithTitle(t *testing.T) {
	md := "```python title=\"example.py\"\nprint('hello')\n```"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="language">python`) {
		t.Errorf("expected python language, got: %s", got)
	}
	if !strings.Contains(got, `ac:name="title">example.py`) {
		t.Errorf("expected title parameter, got: %s", got)
	}
}

func TestCodeBlockWithoutTitle(t *testing.T) {
	md := "```go\nfmt.Println()\n```"
	got := MarkdownToStorage(md)
	if strings.Contains(got, `ac:name="title"`) {
		t.Errorf("expected no title parameter, got: %s", got)
	}
	if !strings.Contains(got, `ac:name="language">go`) {
		t.Errorf("expected go language, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// Collapsible Sections
// ---------------------------------------------------------------------------

func TestDetailsSummary(t *testing.T) {
	md := "<details>\n<summary>Click me</summary>\n\nSome hidden content\n</details>"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="expand"`) {
		t.Errorf("expected expand macro, got: %s", got)
	}
	if !strings.Contains(got, `title">Click me`) {
		t.Errorf("expected title 'Click me', got: %s", got)
	}
	if !strings.Contains(got, "Some hidden content") {
		t.Errorf("expected hidden content, got: %s", got)
	}
}

func TestDetailsPreservesH1(t *testing.T) {
	md := "<details>\n<summary>Section</summary>\n\n# Inner Heading\n\nContent\n</details>"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "<h1>Inner Heading</h1>") {
		t.Errorf("H1 inside details should be preserved, got: %s", got)
	}
}

func TestDetailsWithoutSummary(t *testing.T) {
	md := "<details>\nSome content\n</details>"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `title">Click to expand`) {
		t.Errorf("expected default title, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// MkDocs Admonitions
// ---------------------------------------------------------------------------

func TestMkDocsNote(t *testing.T) {
	md := "!!! note \"Important Note\"\n    This is the content\n    More content"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="note"`) {
		t.Errorf("expected note macro, got: %s", got)
	}
	if !strings.Contains(got, `title">Important Note`) {
		t.Errorf("expected title, got: %s", got)
	}
	if !strings.Contains(got, "This is the content") {
		t.Errorf("expected content, got: %s", got)
	}
}

func TestMkDocsWarning(t *testing.T) {
	md := "!!! warning\n    Be careful"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="warning"`) {
		t.Errorf("expected warning macro, got: %s", got)
	}
}

func TestMkDocsTip(t *testing.T) {
	md := "!!! tip\n    Helpful hint"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="tip"`) {
		t.Errorf("expected tip macro, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// Horizontal Rules
// ---------------------------------------------------------------------------

func TestHorizontalRule(t *testing.T) {
	md := "Before\n\n---\n\nAfter"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "<hr />") {
		t.Errorf("expected hr tag, got: %s", got)
	}
}

func TestHorizontalRuleVariants(t *testing.T) {
	for _, rule := range []string{"---", "***", "___"} {
		md := "Before\n\n" + rule + "\n\nAfter"
		got := MarkdownToStorage(md)
		if !strings.Contains(got, "<hr />") {
			t.Errorf("expected hr tag for %q, got: %s", rule, got)
		}
	}
}

// ---------------------------------------------------------------------------
// Status Badges
// ---------------------------------------------------------------------------

func TestStatusBadge(t *testing.T) {
	md := "Status: {status:Green|DONE}"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="status"`) {
		t.Errorf("expected status macro, got: %s", got)
	}
	if !strings.Contains(got, `title">DONE`) {
		t.Errorf("expected DONE title, got: %s", got)
	}
	if !strings.Contains(got, `colour">Green`) {
		t.Errorf("expected Green colour, got: %s", got)
	}
}

func TestStatusBadgeInline(t *testing.T) {
	md := "Task is {status:Red|FAILED} please fix"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `colour">Red`) {
		t.Errorf("expected Red colour, got: %s", got)
	}
	if !strings.Contains(got, "please fix") {
		t.Errorf("expected surrounding text, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// Existing features (regression tests)
// ---------------------------------------------------------------------------

func TestHeadingSkipsFirstH1(t *testing.T) {
	md := "# Page Title\n\n## Section"
	got := MarkdownToStorage(md)
	if strings.Contains(got, "<h1>") {
		t.Errorf("first h1 should be skipped, got: %s", got)
	}
	if !strings.Contains(got, "<h2>Section</h2>") {
		t.Errorf("expected h2, got: %s", got)
	}
}

func TestCodeBlock(t *testing.T) {
	md := "```python\nprint('hello')\n```"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="code"`) {
		t.Errorf("expected code macro, got: %s", got)
	}
	if !strings.Contains(got, "CDATA[print('hello')]") {
		t.Errorf("expected code content, got: %s", got)
	}
}

func TestSimpleTable(t *testing.T) {
	md := "| A | B |\n|---|---|\n| 1 | 2 |"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "<table") {
		t.Errorf("expected table, got: %s", got)
	}
	if !strings.Contains(got, "<th>") {
		t.Errorf("expected th, got: %s", got)
	}
	if !strings.Contains(got, "<td>") {
		t.Errorf("expected td, got: %s", got)
	}
}

func TestBoldAndItalic(t *testing.T) {
	md := "This is **bold** and *italic*"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "<strong>bold</strong>") {
		t.Errorf("expected bold, got: %s", got)
	}
	if !strings.Contains(got, "<em>italic</em>") {
		t.Errorf("expected italic, got: %s", got)
	}
}

func TestInlineCode(t *testing.T) {
	md := "Use `fmt.Println()` here"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "<code>fmt.Println()</code>") {
		t.Errorf("expected inline code, got: %s", got)
	}
}

func TestStrikethrough(t *testing.T) {
	md := "This is ~~deleted~~ text"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "<del>deleted</del>") {
		t.Errorf("expected strikethrough, got: %s", got)
	}
}

func TestExternalLink(t *testing.T) {
	md := "Visit [Google](https://google.com)"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `<a href="https://google.com">Google</a>`) {
		t.Errorf("expected link, got: %s", got)
	}
}

func TestFrontmatterStripping(t *testing.T) {
	md := "---\ntitle: Test\n---\n\n## Section"
	got := MarkdownToStorage(md)
	if strings.Contains(got, "title: Test") {
		t.Errorf("frontmatter should be stripped, got: %s", got)
	}
	if !strings.Contains(got, "<h2>Section</h2>") {
		t.Errorf("expected section heading, got: %s", got)
	}
}

func TestTocReplacement(t *testing.T) {
	md := "## Table of Contents\n\n- [Section](#section)\n\n## Section\n\nContent"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="toc"`) {
		t.Errorf("expected toc macro, got: %s", got)
	}
}

func TestFlatUnorderedList(t *testing.T) {
	md := "- one\n- two\n- three"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "<ul>") {
		t.Errorf("expected ul, got: %s", got)
	}
	if strings.Count(got, "<li>") != 3 {
		t.Errorf("expected 3 list items, got: %s", got)
	}
}

func TestFlatOrderedList(t *testing.T) {
	md := "1. one\n2. two\n3. three"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "<ol>") {
		t.Errorf("expected ol, got: %s", got)
	}
	if strings.Count(got, "<li>") != 3 {
		t.Errorf("expected 3 list items, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// parseCodeFence
// ---------------------------------------------------------------------------

func TestParseCodeFence(t *testing.T) {
	tests := []struct {
		fence     string
		wantLang  string
		wantTitle string
	}{
		{"```", "", ""},
		{"```python", "python", ""},
		{"```python title=\"example.py\"", "python", "example.py"},
		{"```go", "go", ""},
		{"```bash title=\"install.sh\"", "bash", "install.sh"},
	}
	for _, tt := range tests {
		lang, title := parseCodeFence(tt.fence)
		if lang != tt.wantLang {
			t.Errorf("parseCodeFence(%q): lang = %q, want %q", tt.fence, lang, tt.wantLang)
		}
		if title != tt.wantTitle {
			t.Errorf("parseCodeFence(%q): title = %q, want %q", tt.fence, title, tt.wantTitle)
		}
	}
}

// ---------------------------------------------------------------------------
// parseTableAlignments
// ---------------------------------------------------------------------------

func TestParseTableAlignments(t *testing.T) {
	tests := []struct {
		sep  string
		want []string
	}{
		{"|---|---|", []string{"", ""}},
		{"|:---|---|", []string{"left", ""}},
		{"|---:|---|", []string{"right", ""}},
		{"|:---:|:---:|", []string{"center", "center"}},
		{"|:---|:---:|---:|", []string{"left", "center", "right"}},
	}
	for _, tt := range tests {
		got := parseTableAlignments(tt.sep)
		if len(got) != len(tt.want) {
			t.Errorf("parseTableAlignments(%q): len = %d, want %d", tt.sep, len(got), len(tt.want))
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("parseTableAlignments(%q)[%d] = %q, want %q", tt.sep, i, got[i], tt.want[i])
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Page Properties macro
// ---------------------------------------------------------------------------

func TestBuildPagePropertiesMacro(t *testing.T) {
	props := []PageProperty{
		{Key: "Owner", Value: "AI Tooling Guild"},
		{Key: "Classification", Value: "Internal"},
		{Key: "Status", Value: "{status:Green|approved}"},
		{Key: "Reviewed on", Value: "2026-03-06"},
	}
	got := BuildPagePropertiesMacro("status", props)

	checks := []struct {
		desc     string
		contains string
	}{
		{"details macro", `ac:name="details"`},
		{"id param", `ac:name="id">status`},
		{"owner label", `<strong>Owner</strong>`},
		{"owner value", `AI Tooling Guild`},
		{"classification", `Internal`},
		{"status macro", `ac:name="status"`},
		{"status title", `ac:name="title">approved`},
		{"status colour", `ac:name="colour">Green`},
		{"date tag", `<time datetime="2026-03-06" />`},
		{"table wrapper", `<table data-layout="align-start">`},
	}

	for _, c := range checks {
		if !strings.Contains(got, c.contains) {
			t.Errorf("BuildPagePropertiesMacro: should contain %q (%s)\ngot: %s", c.contains, c.desc, got)
		}
	}
}

func TestBuildPagePropertiesMacroNoID(t *testing.T) {
	props := []PageProperty{
		{Key: "Owner", Value: "Team"},
	}
	got := BuildPagePropertiesMacro("", props)
	if strings.Contains(got, `ac:name="id"`) {
		t.Errorf("should not contain id param when empty, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// Properties Report (HTML comment)
// ---------------------------------------------------------------------------

func TestPropertiesReportHTMLComment(t *testing.T) {
	md := `## Report

<!-- confluence:properties-report cql="label = 'ai-process' and space = currentSpace()" firstcolumn="Document" headings="Status, Classification, Reviewed on" -->

## Next Section`

	got := MarkdownToStorage(md)

	checks := []struct {
		desc     string
		contains string
	}{
		{"detailssummary macro", `ac:name="detailssummary"`},
		{"firstcolumn", `ac:name="firstcolumn">Document`},
		{"headings", `ac:name="headings">Status, Classification, Reviewed on`},
		{"cql", `ac:name="cql">label = 'ai-process' and space = currentSpace()`},
	}

	for _, c := range checks {
		if !strings.Contains(got, c.contains) {
			t.Errorf("PropertiesReportHTMLComment: should contain %q (%s)\ngot: %s", c.contains, c.desc, got)
		}
	}
}

// ---------------------------------------------------------------------------
// Properties Report (shorthand)
// ---------------------------------------------------------------------------

func TestPropertiesReportShorthand(t *testing.T) {
	md := `## Report

{properties-report: label="ai-process", columns="Status, Classification, Reviewed on"}

## Next Section`

	got := MarkdownToStorage(md)

	checks := []struct {
		desc     string
		contains string
	}{
		{"detailssummary macro", `ac:name="detailssummary"`},
		{"default firstcolumn", `ac:name="firstcolumn">Title`},
		{"headings from columns", `ac:name="headings">Status, Classification, Reviewed on`},
		{"cql from label", `ac:name="cql">label = "ai-process" and space = currentSpace()`},
	}

	for _, c := range checks {
		if !strings.Contains(got, c.contains) {
			t.Errorf("PropertiesReportShorthand: should contain %q (%s)\ngot: %s", c.contains, c.desc, got)
		}
	}
}

func TestPropertiesReportShorthandWithSortBy(t *testing.T) {
	md := `{properties-report: label="docs", columns="Status", sortBy="title"}`
	got := MarkdownToStorage(md)
	if !strings.Contains(got, `ac:name="sortBy">title`) {
		t.Errorf("expected sortBy parameter, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// Ignore blocks
// ---------------------------------------------------------------------------

func TestIgnoreBlock(t *testing.T) {
	md := "# Title\n\nVisible paragraph\n\n<!-- confluence:ignore-start -->\n\n| Col1 | Col2 |\n|------|------|\n| a | b |\n\n<!-- confluence:ignore-end -->\n\nAfter ignore"
	got := MarkdownToStorage(md)
	if strings.Contains(got, "<table") {
		t.Errorf("Table inside ignore block should be skipped, got: %s", got)
	}
	if !strings.Contains(got, "Visible paragraph") {
		t.Errorf("Content before ignore block should be preserved, got: %s", got)
	}
	if !strings.Contains(got, "After ignore") {
		t.Errorf("Content after ignore block should be preserved, got: %s", got)
	}
}

func TestIgnoreBlockDoesNotAffectCodeBlocks(t *testing.T) {
	md := "# Title\n\n```\n<!-- confluence:ignore-start -->\ncode\n<!-- confluence:ignore-end -->\n```\n\nVisible"
	got := MarkdownToStorage(md)
	if !strings.Contains(got, "confluence:ignore-start") {
		t.Errorf("Ignore directives inside code blocks should be literal, got: %s", got)
	}
	if !strings.Contains(got, "Visible") {
		t.Errorf("Content after code block should be preserved, got: %s", got)
	}
}

// ---------------------------------------------------------------------------
// Complex integration tests
// ---------------------------------------------------------------------------

func TestFullDocument(t *testing.T) {
	md := `---
title: Test Doc
---

# Test Document

## Introduction

This is a **test** document with *formatting*.

> [!NOTE]
> Pay attention to this

- Item 1
  - Nested A
  - Nested B
- Item 2

| Name | Age |
|:---|---:|
| Alice | 30 |
| Bob | 25 |

` + "```go" + ` title="main.go"
func main() {}
` + "```" + `

<details>
<summary>Advanced</summary>

More info here

</details>

---

## Conclusion

All done.
`

	got := MarkdownToStorage(md)

	checks := []struct {
		desc     string
		contains string
	}{
		{"no frontmatter", "title: Test Doc"},
		{"h1 skipped", "<h1>"},
		{"h2 present", "<h2>Introduction</h2>"},
		{"bold", "<strong>test</strong>"},
		{"italic", "<em>formatting</em>"},
		{"note alert", `ac:name="note"`},
		{"nested list", "<ul><li>Nested A"},
		{"table", "<table"},
		{"left align", `text-align: left`},
		{"right align", `text-align: right`},
		{"code block", `ac:name="code"`},
		{"code title", `title">main.go`},
		{"expand macro", `ac:name="expand"`},
		{"hr", "<hr />"},
		{"conclusion", "<h2>Conclusion</h2>"},
	}

	for _, c := range checks {
		if c.desc == "no frontmatter" || c.desc == "h1 skipped" {
			if strings.Contains(got, c.contains) {
				t.Errorf("full doc: should NOT contain %q (%s)", c.contains, c.desc)
			}
		} else {
			if !strings.Contains(got, c.contains) {
				t.Errorf("full doc: should contain %q (%s)\ngot: %s", c.contains, c.desc, got)
			}
		}
	}
}

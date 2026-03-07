package confluence

import (
	"fmt"
	"regexp"
	"strings"
)

// MarkdownToStorage converts Markdown to Confluence storage format (XHTML).
// The first h1 heading is skipped because Confluence already displays the page title.
// Any "Table of Contents" section is replaced with the Confluence toc macro.
func MarkdownToStorage(md string) string {
	return MarkdownToStorageWithLinks(md, nil)
}

// MarkdownToStorageWithLinks converts Markdown to Confluence storage format.
// linkMap maps relative .md paths to exact Confluence page titles for resolving internal links.
//
// Supported features:
//   - Headings (h1-h6), with first h1 skipped (Confluence shows page title)
//   - Code blocks with language, title, and collapse support
//   - Tables with column alignment (left/center/right)
//   - Blockquotes with GitHub Alerts ([!NOTE], [!WARNING], [!TIP], [!CAUTION], [!IMPORTANT])
//   - Nested ordered and unordered lists
//   - Task lists (checkboxes: - [ ] and - [x])
//   - MkDocs admonitions (!!! note "title")
//   - Collapsible sections (<details>/<summary> → expand macro)
//   - Inline formatting: bold, italic, code, strikethrough, links, images (with alt text)
//   - Status badges: {status:Color|Text}
//   - Horizontal rules
//   - YAML frontmatter stripping
//   - Metadata block conversion (**Key**: Value → table)
//   - Table of Contents auto-replacement
func MarkdownToStorageWithLinks(md string, linkMap map[string]string) string {
	md = stripFrontmatter(md)
	md = convertMetadataBlock(md)

	lines := strings.Split(md, "\n")
	var out []string

	// State
	inCodeBlock := false
	codeLang := ""
	codeTitle := ""
	var codeLines []string
	firstH1Skipped := false
	inTocSection := false
	inIgnoreBlock := false
	inTable := false
	var tableAlignments []string

	inline := func(text string) string {
		return convertInline(text, linkMap)
	}

	flushTable := func() {
		if inTable {
			out = append(out, "</tbody></table>")
			inTable = false
			tableAlignments = nil
		}
	}

	i := 0
	for i < len(lines) {
		line := lines[i]

			trimmed := strings.TrimSpace(line)

		// --- Code blocks ---
		if strings.HasPrefix(trimmed, "```") {
			if inCodeBlock {
				// Close code block
				code := escapeCDATA(strings.Join(codeLines, "\n"))
				out = append(out, buildCodeMacro(codeLang, codeTitle, code))
				inCodeBlock = false
				codeLang = ""
				codeTitle = ""
				codeLines = nil
				i++
				continue
			}
			// Open code block
			flushTable()
			inCodeBlock = true
			codeLang, codeTitle = parseCodeFence(trimmed)
			codeLines = nil
			i++
			continue
		}
		if inCodeBlock {
			codeLines = append(codeLines, line)
			i++
			continue
		}

		// --- Ignore blocks (confluence:ignore-start / confluence:ignore-end) ---
		if strings.Contains(trimmed, "confluence:ignore-start") {
			flushTable()
			inIgnoreBlock = true
			i++
			continue
		}
		if strings.Contains(trimmed, "confluence:ignore-end") {
			inIgnoreBlock = false
			i++
			continue
		}
		if inIgnoreBlock {
			i++
			continue
		}

		// --- Collapsible sections (details/summary) ---
		if strings.HasPrefix(trimmed, "<details>") || strings.HasPrefix(trimmed, "<details ") {
			flushTable()
			title, content, end := collectDetails(lines, i)
			if end > i {
				innerHTML := convertBodyWithLinks(content, linkMap)
				if title == "" {
					title = "Click to expand"
				}
				out = append(out, fmt.Sprintf(`<ac:structured-macro ac:name="expand"><ac:parameter ac:name="title">%s</ac:parameter><ac:rich-text-body>%s</ac:rich-text-body></ac:structured-macro>`, title, innerHTML))
				i = end
				continue
			}
		}

		// --- Blank line ---
		if trimmed == "" {
			flushTable()
			i++
			continue
		}

		// --- Horizontal rule ---
		if isHorizontalRule(trimmed) {
			if !inTocSection {
				flushTable()
				out = append(out, "<hr />")
			} else {
				inTocSection = false
			}
			i++
			continue
		}

		// --- Headings ---
		if heading, level := parseHeading(trimmed); heading != "" {
			flushTable()

			// Skip the first h1 — Confluence already shows the page title
			if level == 1 && !firstH1Skipped {
				firstH1Skipped = true
				i++
				continue
			}

			// Replace "Table of Contents" heading with Confluence toc macro
			headingLower := strings.ToLower(heading)
			if headingLower == "table of contents" || headingLower == "contents" || headingLower == "toc" {
				out = append(out, `<ac:structured-macro ac:name="toc"><ac:parameter ac:name="printable">true</ac:parameter><ac:parameter ac:name="style">none</ac:parameter><ac:parameter ac:name="maxLevel">3</ac:parameter><ac:parameter ac:name="minLevel">2</ac:parameter></ac:structured-macro>`)
				inTocSection = true
				i++
				continue
			}

			inTocSection = false
			out = append(out, fmt.Sprintf("<h%d>%s</h%d>", level, inline(heading), level))
			i++
			continue
		}

		// --- Skip lines inside TOC section ---
		if inTocSection {
			i++
			continue
		}

		// --- MkDocs admonitions (!!! type "title") ---
		if strings.HasPrefix(trimmed, "!!! ") {
			flushTable()
			adType, adTitle, adContent, end := collectAdmonition(lines, i)
			macroName := admonitionToMacro(adType)
			body := inline(strings.Join(adContent, "<br />"))
			out = append(out, buildPanelMacro(macroName, adTitle, body))
			i = end
			continue
		}

		// --- Tables ---
		if isTableRow(trimmed) {
			if isTableSeparator(trimmed) {
				tableAlignments = parseTableAlignments(trimmed)
				i++
				continue
			}
			cells := parseTableRow(trimmed)
			if !inTable {
				// First row is header
				out = append(out, `<table data-layout="full-width"><thead><tr>`)
				for ci, cell := range cells {
					align := alignAttr(tableAlignments, ci)
					out = append(out, fmt.Sprintf("<th%s><strong>%s</strong></th>", align, inline(cell)))
				}
				out = append(out, "</tr></thead><tbody>")
				inTable = true
			} else {
				out = append(out, "<tr>")
				for ci, cell := range cells {
					align := alignAttr(tableAlignments, ci)
					out = append(out, fmt.Sprintf("<td%s>%s</td>", align, inline(cell)))
				}
				out = append(out, "</tr>")
			}
			i++
			continue
		}

		flushTable()

		// --- Blockquotes (with GitHub Alerts) ---
		if strings.HasPrefix(trimmed, "> ") || trimmed == ">" {
			macroName, title, quoteContent, end := collectBlockquote(lines, i)
			body := inline(strings.Join(quoteContent, "<br />"))
			out = append(out, buildPanelMacro(macroName, title, body))
			i = end
			continue
		}

		// --- Lists (nested, with task list support) ---
		if isListItem(line) {
			items, end := collectListItems(lines, i)
			out = append(out, buildNestedList(items, inline))
			i = end
			continue
		}

		// --- Properties Report (HTML comment) ---
		if m := rePropertiesReportComment.FindStringSubmatch(trimmed); m != nil {
			params := parseDirectiveParams(m[1])
			out = append(out, buildPropertiesReportMacro(params))
			i++
			continue
		}

		// --- Properties Report (shorthand) ---
		if m := rePropertiesReportShort.FindStringSubmatch(trimmed); m != nil {
			params := parseDirectiveParams(m[1])
			out = append(out, buildPropertiesReportMacro(params))
			i++
			continue
		}

		// --- HTML/Confluence pass-through ---
		if strings.HasPrefix(trimmed, "<table") || strings.HasPrefix(trimmed, "<div") || strings.HasPrefix(trimmed, "<ac:") {
			out = append(out, trimmed)
			i++
			continue
		}

		// --- Regular paragraph ---
		out = append(out, fmt.Sprintf("<p>%s</p>", inline(trimmed)))
		i++
	}

	flushTable()

	return strings.Join(out, "\n")
}

// convertBodyWithLinks converts markdown body content without skipping the
// first H1 heading. Used for nested content (e.g. inside <details> blocks)
// where headings should be preserved as-is.
func convertBodyWithLinks(md string, linkMap map[string]string) string {
	// Prepend a dummy H1 that will be skipped, so the real content is preserved.
	return MarkdownToStorageWithLinks("# _\n\n"+md, linkMap)
}

// buildPanelMacro builds a Confluence panel macro (info, note, warning, tip)
// with an optional title.
func buildPanelMacro(macroName, title, body string) string {
	if title != "" {
		return fmt.Sprintf(`<ac:structured-macro ac:name="%s"><ac:parameter ac:name="title">%s</ac:parameter><ac:rich-text-body><p>%s</p></ac:rich-text-body></ac:structured-macro>`, macroName, title, body)
	}
	return fmt.Sprintf(`<ac:structured-macro ac:name="%s"><ac:rich-text-body><p>%s</p></ac:rich-text-body></ac:structured-macro>`, macroName, body)
}

// ---------------------------------------------------------------------------
// Code blocks
// ---------------------------------------------------------------------------

// parseCodeFence extracts language and title from a code fence line.
// Supports: ```python, ```python title="Example", ```mermaid
func parseCodeFence(fence string) (lang, title string) {
	rest := strings.TrimPrefix(fence, "```")
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return "", ""
	}

	// Check for title="..." attribute
	if idx := strings.Index(rest, " "); idx != -1 {
		lang = rest[:idx]
		attrs := rest[idx+1:]
		if m := reCodeTitle.FindStringSubmatch(attrs); m != nil {
			title = m[1]
		}
	} else {
		lang = rest
	}
	return lang, title
}

func buildCodeMacro(lang, title, code string) string {
	var sb strings.Builder
	sb.WriteString(`<ac:structured-macro ac:name="code">`)
	if lang != "" {
		sb.WriteString(fmt.Sprintf(`<ac:parameter ac:name="language">%s</ac:parameter>`, lang))
	}
	if title != "" {
		sb.WriteString(fmt.Sprintf(`<ac:parameter ac:name="title">%s</ac:parameter>`, title))
	}
	sb.WriteString(fmt.Sprintf(`<ac:plain-text-body><![CDATA[%s]]></ac:plain-text-body>`, code))
	sb.WriteString(`</ac:structured-macro>`)
	return sb.String()
}

// ---------------------------------------------------------------------------
// Collapsible sections (<details>/<summary>)
// ---------------------------------------------------------------------------

func collectDetails(lines []string, startIdx int) (title, content string, endIdx int) {
	i := startIdx + 1
	var contentLines []string
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "</details>" {
			return title, strings.Join(contentLines, "\n"), i + 1
		}
		if m := reSummary.FindStringSubmatch(trimmed); m != nil {
			title = m[1]
			i++
			continue
		}
		contentLines = append(contentLines, lines[i])
		i++
	}
	// No closing tag found — return unchanged
	return "", "", startIdx
}

// ---------------------------------------------------------------------------
// Blockquotes with GitHub Alerts
// ---------------------------------------------------------------------------

var reGitHubAlert = regexp.MustCompile(`^\[!(NOTE|WARNING|TIP|CAUTION|IMPORTANT)\]$`)

func collectBlockquote(lines []string, startIdx int) (macroName, title string, content []string, endIdx int) {
	macroName = "info" // default for plain blockquotes
	i := startIdx
	var quoteLines []string

	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(trimmed, "> ") && trimmed != ">" {
			break
		}
		text := strings.TrimPrefix(trimmed, "> ")
		text = strings.TrimPrefix(text, ">") // bare ">"
		quoteLines = append(quoteLines, text)
		i++
	}

	// Check first line for GitHub Alert syntax
	if len(quoteLines) > 0 {
		if m := reGitHubAlert.FindStringSubmatch(strings.TrimSpace(quoteLines[0])); m != nil {
			alertType := strings.ToLower(m[1])
			macroName = admonitionToMacro(alertType)
			if alertType == "important" {
				title = "Important"
			}
			quoteLines = quoteLines[1:]
			// Trim leading blank line after alert marker
			if len(quoteLines) > 0 && strings.TrimSpace(quoteLines[0]) == "" {
				quoteLines = quoteLines[1:]
			}
		}
	}

	return macroName, title, quoteLines, i
}

// ---------------------------------------------------------------------------
// MkDocs admonitions (!!! type "title")
// ---------------------------------------------------------------------------

var reAdmonition = regexp.MustCompile(`^!!!\s+(\w+)(?:\s+"([^"]*)")?$`)

func collectAdmonition(lines []string, startIdx int) (adType, title string, content []string, endIdx int) {
	m := reAdmonition.FindStringSubmatch(strings.TrimSpace(lines[startIdx]))
	if m == nil {
		return "info", "", nil, startIdx + 1
	}
	adType = m[1]
	title = m[2]

	i := startIdx + 1
	for i < len(lines) {
		line := lines[i]
		// Admonition content must be indented (4 spaces or 1 tab)
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			content = append(content, strings.TrimSpace(line))
			i++
			continue
		}
		// Blank lines within indented block are OK
		if strings.TrimSpace(line) == "" && i+1 < len(lines) {
			next := lines[i+1]
			if strings.HasPrefix(next, "    ") || strings.HasPrefix(next, "\t") {
				content = append(content, "")
				i++
				continue
			}
		}
		break
	}
	return adType, title, content, i
}

func admonitionToMacro(adType string) string {
	switch adType {
	case "note", "abstract", "summary", "tldr":
		return "note"
	case "tip", "hint", "success", "check", "done":
		return "tip"
	case "warning", "caution", "attention", "danger", "failure", "fail", "missing", "bug":
		return "warning"
	default:
		return "info"
	}
}

// ---------------------------------------------------------------------------
// Lists (nested, with task list support)
// ---------------------------------------------------------------------------

type listItem struct {
	depth   int
	tag     string // "ul" or "ol"
	content string
	isTask  bool
	checked bool
}

var reOrderedItem = regexp.MustCompile(`^\d+\.\s(.*)`)

func isListItem(line string) bool {
	stripped := strings.TrimLeft(line, " \t")
	if strings.HasPrefix(stripped, "- ") || strings.HasPrefix(stripped, "* ") {
		return true
	}
	return reOrderedItem.MatchString(stripped)
}

func collectListItems(lines []string, startIdx int) ([]listItem, int) {
	var items []listItem
	i := startIdx

	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Blank line — check if list continues after it
		if trimmed == "" {
			j := i + 1
			for j < len(lines) && strings.TrimSpace(lines[j]) == "" {
				j++
			}
			if j < len(lines) && isListItem(lines[j]) {
				i = j
				continue
			}
			break
		}

		if !isListItem(line) {
			break
		}

		stripped := strings.TrimLeft(line, " \t")
		leadingSpaces := len(line) - len(stripped)
		depth := leadingSpaces / 2

		item := listItem{depth: depth}

		switch {
		case strings.HasPrefix(stripped, "- [x] ") || strings.HasPrefix(stripped, "- [X] "):
			item.tag = "ul"
			item.isTask = true
			item.checked = true
			item.content = stripped[6:]
		case strings.HasPrefix(stripped, "- [ ] "):
			item.tag = "ul"
			item.isTask = true
			item.checked = false
			item.content = stripped[6:]
		case strings.HasPrefix(stripped, "- "):
			item.tag = "ul"
			item.content = stripped[2:]
		case strings.HasPrefix(stripped, "* "):
			item.tag = "ul"
			item.content = stripped[2:]
		default:
			m := reOrderedItem.FindStringSubmatch(stripped)
			if m != nil {
				item.tag = "ol"
				item.content = m[1]
			}
		}

		items = append(items, item)
		i++
	}

	return items, i
}

func buildNestedList(items []listItem, inlineFn func(string) string) string {
	if len(items) == 0 {
		return ""
	}
	var out strings.Builder
	renderListLevel(&out, items, 0, items[0].depth, inlineFn)
	return out.String()
}

func renderListLevel(out *strings.Builder, items []listItem, startIdx, baseDepth int, inlineFn func(string) string) int {
	if startIdx >= len(items) {
		return startIdx
	}

	tag := items[startIdx].tag
	out.WriteString("<" + tag + ">")

	i := startIdx
	for i < len(items) {
		item := items[i]
		if item.depth < baseDepth {
			break
		}

		content := inlineFn(item.content)
		if item.isTask {
			if item.checked {
				content = `<ac:structured-macro ac:name="status"><ac:parameter ac:name="title">DONE</ac:parameter><ac:parameter ac:name="colour">Green</ac:parameter></ac:structured-macro> ` + content
			} else {
				content = `<ac:structured-macro ac:name="status"><ac:parameter ac:name="title">TODO</ac:parameter><ac:parameter ac:name="colour">Grey</ac:parameter></ac:structured-macro> ` + content
			}
		}

		out.WriteString("<li>" + content)
		i++

		// If next item is deeper, render nested list inside this li
		if i < len(items) && items[i].depth > baseDepth {
			i = renderListLevel(out, items, i, items[i].depth, inlineFn)
		}

		out.WriteString("</li>")
	}

	out.WriteString("</" + tag + ">")
	return i
}

// ---------------------------------------------------------------------------
// Tables
// ---------------------------------------------------------------------------

func parseTableAlignments(separator string) []string {
	separator = strings.TrimPrefix(separator, "|")
	separator = strings.TrimSuffix(separator, "|")
	parts := strings.Split(separator, "|")
	aligns := make([]string, len(parts))
	for i, p := range parts {
		p = strings.TrimSpace(p)
		left := strings.HasPrefix(p, ":")
		right := strings.HasSuffix(p, ":")
		switch {
		case left && right:
			aligns[i] = "center"
		case right:
			aligns[i] = "right"
		case left:
			aligns[i] = "left"
		default:
			aligns[i] = ""
		}
	}
	return aligns
}

func alignAttr(alignments []string, colIdx int) string {
	if colIdx < len(alignments) && alignments[colIdx] != "" {
		return fmt.Sprintf(` style="text-align: %s"`, alignments[colIdx])
	}
	return ""
}

func isTableRow(line string) bool {
	return strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|")
}

func isTableSeparator(line string) bool {
	cleaned := strings.ReplaceAll(line, "|", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, ":", "")
	cleaned = strings.TrimSpace(cleaned)
	return cleaned == ""
}

func parseTableRow(line string) []string {
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	parts := strings.Split(line, "|")
	cells := make([]string, len(parts))
	for i, p := range parts {
		cells[i] = strings.TrimSpace(p)
	}
	return cells
}

// ---------------------------------------------------------------------------
// Headings
// ---------------------------------------------------------------------------

func parseHeading(line string) (string, int) {
	for level := 6; level >= 1; level-- {
		prefix := strings.Repeat("#", level) + " "
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(line[len(prefix):]), level
		}
	}
	return "", 0
}

// ---------------------------------------------------------------------------
// Horizontal rules
// ---------------------------------------------------------------------------

func isHorizontalRule(trimmed string) bool {
	return trimmed == "---" || trimmed == "***" || trimmed == "___" ||
		trimmed == "- - -" || trimmed == "* * *"
}

// ---------------------------------------------------------------------------
// Inline formatting
// ---------------------------------------------------------------------------

var (
	reBold      = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reItalicAlt = regexp.MustCompile(`\*([^*]+?)\*`)
	reCode      = regexp.MustCompile("`([^`]+)`")
	reStrike    = regexp.MustCompile(`~~(.+?)~~`)
	reLink      = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reImage     = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	reMetadata  = regexp.MustCompile(`^\*\*(.+?)\*\*:\s*(.+)$`)
	reCodeTitle = regexp.MustCompile(`title="([^"]+)"`)
	reSummary   = regexp.MustCompile(`<summary>(.+?)</summary>`)
	reStatus    = regexp.MustCompile(`\{status:(\w+)\|([^}]+)\}`)
)

func convertInline(text string, linkMap map[string]string) string {
	// Images with alt text preservation
	text = reImage.ReplaceAllStringFunc(text, func(s string) string {
		m := reImage.FindStringSubmatch(s)
		if m == nil {
			return s
		}
		alt, src := m[1], m[2]
		if alt != "" {
			return fmt.Sprintf(`<ac:image ac:alt="%s"><ri:url ri:value="%s" /></ac:image>`, alt, src)
		}
		return fmt.Sprintf(`<ac:image><ri:url ri:value="%s" /></ac:image>`, src)
	})

	// Links
	text = reLink.ReplaceAllStringFunc(text, func(s string) string {
		m := reLink.FindStringSubmatch(s)
		if m == nil {
			return s
		}
		label, href := m[1], m[2]
		if strings.HasPrefix(href, "#") {
			return label
		}
		if strings.HasSuffix(href, ".md") || strings.HasPrefix(href, "./") || strings.HasPrefix(href, "../") {
			// Try to resolve relative .md link to a Confluence page link
			if linkMap != nil {
				// Try the href as-is, then normalized variants
				candidates := []string{href}
				// Strip ./ prefix
				if strings.HasPrefix(href, "./") {
					candidates = append(candidates, href[2:])
				}
				// Strip .md and add it back (normalize path)
				cleaned := strings.TrimSuffix(href, ".md")
				cleaned = strings.TrimPrefix(cleaned, "./")
				cleaned = strings.TrimPrefix(cleaned, "../")
				candidates = append(candidates, cleaned+".md")
				// Just the filename
				parts := strings.Split(href, "/")
				candidates = append(candidates, parts[len(parts)-1])

				for _, candidate := range candidates {
					if pageTitle, ok := linkMap[candidate]; ok {
						return fmt.Sprintf(`<ac:link><ri:page ri:content-title="%s" /><ac:plain-text-link-body><![CDATA[%s]]></ac:plain-text-link-body></ac:link>`, pageTitle, label)
					}
				}
			}
			return label
		}
		return fmt.Sprintf(`<a href="%s">%s</a>`, href, label)
	})

	// Bold before italic
	text = reBold.ReplaceAllString(text, "<strong>$1</strong>")
	// Code
	text = reCode.ReplaceAllString(text, "<code>$1</code>")
	// Strikethrough
	text = reStrike.ReplaceAllString(text, "<del>$1</del>")

	// Status badges: {status:Green|DONE}
	text = reStatus.ReplaceAllStringFunc(text, func(s string) string {
		m := reStatus.FindStringSubmatch(s)
		if m == nil {
			return s
		}
		colour, title := m[1], m[2]
		return fmt.Sprintf(`<ac:structured-macro ac:name="status"><ac:parameter ac:name="title">%s</ac:parameter><ac:parameter ac:name="colour">%s</ac:parameter></ac:structured-macro>`, title, colour)
	})

	// Italic (using * — careful not to conflict with bold already converted)
	text = reItalicAlt.ReplaceAllStringFunc(text, func(s string) string {
		if strings.Contains(s, "<strong>") || strings.Contains(s, "</strong>") {
			return s
		}
		m := reItalicAlt.FindStringSubmatch(s)
		if m == nil {
			return s
		}
		return "<em>" + m[1] + "</em>"
	})

	return text
}

// ---------------------------------------------------------------------------
// Frontmatter & metadata
// ---------------------------------------------------------------------------

func stripFrontmatter(md string) string {
	if !strings.HasPrefix(md, "---") {
		return md
	}
	rest := md[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return md
	}
	return strings.TrimSpace(rest[idx+4:])
}

func convertMetadataBlock(md string) string {
	lines := strings.Split(md, "\n")
	var metaRows [][]string
	metaStart := -1
	metaEnd := 0

	i := 0
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "> ") {
			i++
			continue
		}
		break
	}

	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			if len(metaRows) > 0 {
				i++
				continue
			}
			break
		}
		m := reMetadata.FindStringSubmatch(trimmed)
		if m == nil {
			break
		}
		if metaStart == -1 {
			metaStart = i
		}
		metaRows = append(metaRows, []string{m[1], m[2]})
		metaEnd = i + 1
		i++
	}

	if len(metaRows) == 0 {
		return md
	}

	var tbl strings.Builder
	tbl.WriteString(`<table data-layout="full-width"><tbody>`)
	for _, row := range metaRows {
		tbl.WriteString(fmt.Sprintf(`<tr><td style="background-color:#f4f5f7"><strong>%s</strong></td><td>%s</td></tr>`, row[0], row[1]))
	}
	tbl.WriteString(`</tbody></table>`)

	before := lines[:metaStart]
	after := lines[metaEnd:]
	for len(after) > 0 && strings.TrimSpace(after[0]) == "" {
		after = after[1:]
	}

	var result []string
	result = append(result, before...)
	result = append(result, tbl.String())
	result = append(result, after...)
	return strings.Join(result, "\n")
}

// ---------------------------------------------------------------------------
// Page Properties macro (details)
// ---------------------------------------------------------------------------

// PageProperty represents a single key-value property for the Page Properties macro.
type PageProperty struct {
	Key   string
	Value string
}

var reDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// BuildPagePropertiesMacro generates a Confluence "details" (Page Properties) macro.
// The id parameter sets the macro's ID for Page Properties Report queries.
// Properties are rendered as a two-column table inside the macro.
func BuildPagePropertiesMacro(id string, properties []PageProperty) string {
	var sb strings.Builder
	sb.WriteString(`<ac:structured-macro ac:name="details">`)
	if id != "" {
		sb.WriteString(fmt.Sprintf(`<ac:parameter ac:name="id">%s</ac:parameter>`, id))
	}
	sb.WriteString(`<ac:rich-text-body><table data-layout="align-start"><tbody>`)
	for _, prop := range properties {
		sb.WriteString(fmt.Sprintf(`<tr><th><p><strong>%s</strong></p></th><td><p>%s</p></td></tr>`, prop.Key, renderPropertyValue(prop.Value)))
	}
	sb.WriteString(`</tbody></table></ac:rich-text-body></ac:structured-macro>`)
	return sb.String()
}

// renderPropertyValue converts a property value string to Confluence storage format.
// It handles status badges ({status:Color|Text}) and dates (YYYY-MM-DD).
func renderPropertyValue(value string) string {
	// Status badge
	if m := reStatus.FindStringSubmatch(value); m != nil {
		colour, title := m[1], m[2]
		return fmt.Sprintf(`<ac:structured-macro ac:name="status"><ac:parameter ac:name="title">%s</ac:parameter><ac:parameter ac:name="colour">%s</ac:parameter></ac:structured-macro>`, title, colour)
	}
	// Date
	if reDate.MatchString(value) {
		return fmt.Sprintf(`<time datetime="%s" />`, value)
	}
	return value
}

// ---------------------------------------------------------------------------
// Properties Report directive (detailssummary)
// ---------------------------------------------------------------------------

var (
	rePropertiesReportComment = regexp.MustCompile(`<!--\s*confluence:properties-report\s+(.+?)\s*-->`)
	rePropertiesReportShort   = regexp.MustCompile(`^\{properties-report:\s*(.+?)\}$`)
	reDirectiveParam          = regexp.MustCompile(`(\w+)\s*=\s*"([^"]*)"`)
)

// buildPropertiesReportMacro generates a Confluence "detailssummary" macro from parameters.
func buildPropertiesReportMacro(params map[string]string) string {
	var sb strings.Builder
	sb.WriteString(`<ac:structured-macro ac:name="detailssummary">`)

	firstcolumn := params["firstcolumn"]
	if firstcolumn == "" {
		firstcolumn = "Title"
	}
	sb.WriteString(fmt.Sprintf(`<ac:parameter ac:name="firstcolumn">%s</ac:parameter>`, firstcolumn))

	if headings, ok := params["headings"]; ok && headings != "" {
		sb.WriteString(fmt.Sprintf(`<ac:parameter ac:name="headings">%s</ac:parameter>`, headings))
	} else if columns, ok := params["columns"]; ok && columns != "" {
		sb.WriteString(fmt.Sprintf(`<ac:parameter ac:name="headings">%s</ac:parameter>`, columns))
	}

	if cql, ok := params["cql"]; ok && cql != "" {
		sb.WriteString(fmt.Sprintf(`<ac:parameter ac:name="cql">%s</ac:parameter>`, cql))
	} else if label, ok := params["label"]; ok && label != "" {
		sb.WriteString(fmt.Sprintf(`<ac:parameter ac:name="cql">label = "%s" and space = currentSpace()</ac:parameter>`, label))
	}

	if sortBy, ok := params["sortBy"]; ok && sortBy != "" {
		sb.WriteString(fmt.Sprintf(`<ac:parameter ac:name="sortBy">%s</ac:parameter>`, sortBy))
	}

	sb.WriteString(`</ac:structured-macro>`)
	return sb.String()
}

// parseDirectiveParams extracts key="value" pairs from a directive string.
func parseDirectiveParams(s string) map[string]string {
	params := make(map[string]string)
	for _, m := range reDirectiveParam.FindAllStringSubmatch(s, -1) {
		params[m[1]] = m[2]
	}
	return params
}

// ---------------------------------------------------------------------------
// Utilities
// ---------------------------------------------------------------------------

func escapeCDATA(s string) string {
	return strings.ReplaceAll(s, "]]>", "]]]]><![CDATA[>")
}

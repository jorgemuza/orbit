package confluence

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	conflsvc "github.com/jorgemuza/aidlc-cli/internal/service/confluence"
	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish [directory]",
	Short: "Publish a directory of markdown files to Confluence",
	Long: `Publish recursively converts markdown files to Confluence storage format
and creates pages preserving the directory structure. INDEX.md files become
the parent page for their directory. Other .md files become child pages.`,
	Args: cobra.ExactArgs(1),
	Example: `  aidlc confluence publish ./docs --space FO --parent 12345
  aidlc confluence publish ./docs --space FO --parent 12345 --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		space, _ := cmd.Flags().GetString("space")
		parentID, _ := cmd.Flags().GetString("parent")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			fmt.Println("DRY RUN — no pages will be created")
			fmt.Println()
		}

		var client *conflsvc.Client
		if !dryRun {
			var err error
			client, err = resolveConfluenceClient(cmd)
			if err != nil {
				return err
			}
		}

		return publishDir(client, dir, space, parentID, dryRun, 0)
	},
}

func init() {
	publishCmd.Flags().String("space", "", "Confluence space key (required)")
	publishCmd.Flags().String("parent", "", "parent page ID (required)")
	publishCmd.Flags().Bool("dry-run", false, "preview without creating pages")
	_ = publishCmd.MarkFlagRequired("space")
	_ = publishCmd.MarkFlagRequired("parent")
}

func publishDir(client *conflsvc.Client, dir, space, parentID string, dryRun bool, depth int) error {
	indent := strings.Repeat("  ", depth)

	// Check for INDEX.md — it becomes the parent page for this directory
	indexPath := filepath.Join(dir, "INDEX.md")
	currentParentID := parentID

	if _, err := os.Stat(indexPath); err == nil {
		title := titleFromMarkdown(indexPath)
		if dryRun {
			fmt.Printf("%s📄 CREATE page: %q (from INDEX.md) under parent %s\n", indent, title, currentParentID)
			// Use a placeholder ID for dry run
			currentParentID = fmt.Sprintf("<%s>", title)
		} else {
			content, err := readAndConvertMarkdown(indexPath)
			if err != nil {
				return fmt.Errorf("converting %s: %w", indexPath, err)
			}
			page, err := client.CreatePage(space, currentParentID, title, content)
			if err != nil {
				return fmt.Errorf("creating page from %s: %w", indexPath, err)
			}
			currentParentID = page.ID
			fmt.Printf("%s✅ Created: %s — %s (ID: %s)\n", indent, filepath.Base(dir), title, page.ID)
		}
	}

	// Process subdirectories first (they have their own INDEX.md)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory %s: %w", dir, err)
	}

	// Subdirectories
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		subdir := filepath.Join(dir, entry.Name())
		if err := publishDir(client, subdir, space, currentParentID, dryRun, depth+1); err != nil {
			return err
		}
	}

	// Markdown files (skip INDEX.md — already processed)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || entry.Name() == "INDEX.md" {
			continue
		}
		filePath := filepath.Join(dir, entry.Name())
		title := titleFromMarkdown(filePath)

		if dryRun {
			fmt.Printf("%s  📄 CREATE page: %q (from %s) under parent %s\n", indent, title, entry.Name(), currentParentID)
		} else {
			content, err := readAndConvertMarkdown(filePath)
			if err != nil {
				return fmt.Errorf("converting %s: %w", filePath, err)
			}
			page, err := client.CreatePage(space, currentParentID, title, content)
			if err != nil {
				return fmt.Errorf("creating page from %s: %w", filePath, err)
			}
			fmt.Printf("%s  ✅ Created: %s — %s (ID: %s)\n", indent, entry.Name(), title, page.ID)
		}
	}

	return nil
}

// titleFromMarkdown extracts a title from a markdown file.
// Uses the frontmatter title, or the first # heading, or the filename.
func titleFromMarkdown(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return filenameToTitle(path)
	}
	content := string(data)

	// Check frontmatter for title
	if strings.HasPrefix(content, "---") {
		rest := content[3:]
		idx := strings.Index(rest, "\n---")
		if idx != -1 {
			frontmatter := rest[:idx]
			for _, line := range strings.Split(frontmatter, "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "title:") {
					title := strings.TrimSpace(strings.TrimPrefix(line, "title:"))
					title = strings.Trim(title, `"'`)
					if title != "" {
						return title
					}
				}
			}
		}
	}

	// Check for first # heading
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}

	return filenameToTitle(path)
}

func filenameToTitle(path string) string {
	name := strings.TrimSuffix(filepath.Base(path), ".md")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	return strings.Title(name)
}

// readAndConvertMarkdown reads a markdown file and converts to Confluence storage format.
func readAndConvertMarkdown(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading file %s: %w", path, err)
	}
	return conflsvc.MarkdownToStorage(string(data)), nil
}
